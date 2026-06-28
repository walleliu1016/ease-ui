package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/akke/ease-ui/internal/jsonl"
	"github.com/akke/ease-ui/internal/process"
	"github.com/akke/ease-ui/internal/protocol"
	"github.com/akke/ease-ui/internal/session"
)

// switchSettleDelay 是 SwitchOwner 杀掉外部 claude 后等待 hook server
// 收到 SessionEnd、或者给进程预留退出的保守时间。太短可能切回后立刻
// 看到 "session 已被另一个 claude 持有" 之类的错误；太长则用户感知的
// 切换延迟变大。500ms 是经验值。
const switchSettleDelay = 500 * time.Millisecond

func (a *App) SetClaudeBinary(p string) {
	a.mu().Lock()
	defer a.mu().Unlock()
	a.claudeBin = p
}

func (a *App) ListSessions() ([]jsonl.SessionMeta, error) {
	if a.opts.ClaudeDir != "" {
		jsonl.SetRoot(filepath.Join(a.opts.ClaudeDir, "projects"))
	}
	return jsonl.ScanAll()
}

// GetSessionStates 返回持久化的会话状态（用于启动时恢复）。
func (a *App) GetSessionStates() map[string]string {
	out := map[string]string{}
	list, _ := a.ListSessions()
	for _, m := range list {
		s := a.inst.Get(m.ID)
		if s.State != "" {
			out[m.ID] = s.State
		}
	}
	return out
}

// CreateSession launches a new claude subprocess for the given workdir + prompt.
func (a *App) CreateSession(workDir, _prompt string) (string, error) {
	id, err := newID()
	if err != nil {
		return "", err
	}
	a.mu().RLock()
	bin := a.claudeBin
	if bin == "" {
		bin = a.settings.ClaudePath
	}
	a.mu().RUnlock()
	if bin == "" {
		bin = "claude"
	}

	proc, err := process.Start(workDir, id, bin)
	if err != nil {
		return "", err
	}

	s := session.New(id, workDir)
	s.SetProcessForTest(proc)
	a.registerSession(s)
	a.inst.Put(id, "idle")
	go a.pumpEvents(s, proc)
	return id, nil
}

func (a *App) SendMessage(sessionID, prompt string) error {
	s, ok := a.lookupSession(sessionID)
	if !ok {
		return errSessionNotFound
	}
	return s.Send(prompt)
}

func (a *App) RespondPermission(sessionID, reqID string, allow bool) error {
	s, ok := a.lookupSession(sessionID)
	if !ok {
		return errSessionNotFound
	}
	return s.RespondPermission(reqID, allow)
}

// GetSessionMessages reads the jsonl history for a session and returns
// all messages. Used when switching to a session that has no active process.
func (a *App) GetSessionMessages(sessionID, workDir string, offset, limit int) ([]jsonl.Message, error) {
	root := jsonl.Root()
	encodedDir := encodeProjectDirName(workDir)
	path := filepath.Join(root, encodedDir, sessionID+".jsonl")
	return jsonl.ParseFileRange(path, offset, limit)
}

// encodeProjectDirName converts "/Users/akke/foo" to "-Users-akke-foo"
func encodeProjectDirName(dir string) string {
	// "/Users/akke" → "-Users-akke"（ReplaceAll 会把首个 / 也替换成 -）
	return strings.ReplaceAll(dir, "/", "-")
}

func (a *App) CloseSession(sessionID string) error {
	s, ok := a.lookupSession(sessionID)
	if !ok {
		return errSessionNotFound
	}
	err := s.Close()
	if err == nil {
		a.inst.Put(sessionID, "done")
	}
	return err
}

func newID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

var errSessionNotFound = &appError{code: "E_SESSION_NOT_FOUND", msg: "session not found"}

type appError struct {
	code, msg string
}

func (e *appError) Error() string { return e.msg }

// --- internals ---

func (a *App) mu() *sync.RWMutex { return &a.appMu }

func (a *App) registerSession(s *session.Session) {
	a.appMu.Lock()
	defer a.appMu.Unlock()
	if a.sessions == nil {
		a.sessions = map[string]*session.Session{}
	}
	a.sessions[s.ID] = s
}

func (a *App) lookupSession(id string) (*session.Session, bool) {
	a.appMu.RLock()
	defer a.appMu.RUnlock()
	s, ok := a.sessions[id]
	return s, ok
}

func (a *App) pumpEvents(s *session.Session, p *process.Process) {
	topic := "session:" + s.ID
	for line := range p.Events() {
		// 广播原始行给前端，同时也发布到内部 bus
		a.bus.Publish(topic, line)
		if a.ctx != nil {
			wailsruntime.EventsEmit(a.ctx, topic, string(line))
		}

		// 解析事件并更新 session 状态
		evt, err := protocol.Parse(line)
		if err != nil {
			continue
		}
		switch evt.Type {
		case protocol.EvtPermissionReq:
			var req protocol.PermissionRequest
			if json.Unmarshal(evt.Data, &req) == nil {
				s.RegisterPermission(req.RequestID)
			}
		case protocol.EvtResult:
			s.SetIdle()
		}
	}
	// 子进程退出：清理状态
	s.SetIdle()
	if a.ctx != nil {
		wailsruntime.EventsEmit(a.ctx, topic, `{"type":"done"}`)
	}
}

// SwitchOwner 切换 session 的写权限归属。
//
//   - target="app":    从外部终端切回 Ease UI 控制。会 kill 外部 claude
//     进程 (优先按 Session.ExtPIDFile 记录的 pidfile，其次按 pkill
//     -f 兜底)，等 hook server 收到 SessionEnd，再起新的 stream-json
//     进程接管。如果 prompt != ""，会作为新一轮 prompt 写入新进程。
//   - target="terminal": 从 App 控制切到外部终端。关 stream-json 进程，
//     调 OpenInTerminal 起外部 claude -r，Session 标记 Terminal-owned。
//
// 整个流程在 Session.switchMu 持锁下完成，避免与 Send/RespondPermission
// 抢同一把锁。
func (a *App) SwitchOwner(sessionID, target, prompt string) error {
	s, ok := a.lookupSession(sessionID)
	if !ok {
		return errSessionNotFound
	}
	s.SwitchLock().Lock()
	defer s.SwitchLock().Unlock()

	switch target {
	case "app":
		return a.switchToApp(s, prompt)
	case "terminal":
		return a.switchToTerminal(s)
	default:
		return &appError{code: "E_BAD_TARGET", msg: "SwitchOwner: target must be 'app' or 'terminal'"}
	}
}

func (a *App) switchToApp(s *session.Session, prompt string) error {
	if s.Owner() == session.OwnerApp {
		// 已经是 App-owned；prompt 不空则直写（envelope 化由调用方决定）
		if prompt != "" {
			return s.Send(prompt)
		}
		return nil
	}

	// 1) kill 外部 claude 进程
	_, pidfile := s.ExtPID()
	if pidfile != "" {
		if err := killByPIDFile(pidfile); err != nil {
			// pidfile 路径已失效（用户关了 iTerm 后 pid 已死）→ 兜底 pkill
			_ = pkillByPattern(s.ID)
		}
	} else {
		// Windows 或 Launch 未写 pidfile → 直接走 pkill
		_ = pkillByPattern(s.ID)
	}

	// 2) 等 hook server 收到 SessionEnd 给前端同步状态
	time.Sleep(switchSettleDelay)

	// 3) 关掉 session 旧 proc 引用 + 起新 stream-json 进程
	_ = s.Close()
	a.appMu.RLock()
	bin := a.claudeBin
	if bin == "" {
		bin = a.settings.ClaudePath
	}
	workdir := s.WorkDir
	a.appMu.RUnlock()

	proc, err := process.Start(workdir, s.ID, bin)
	if err != nil {
		return &appError{code: "E_START_STREAM", msg: "start stream-json: " + err.Error()}
	}
	s.SetProcessForTest(proc)
	s.SetOwner(session.OwnerApp)
	s.SetMode(session.ModeStream)
	s.SetExtPID(0, "")

	// 4) 重新订阅事件流
	go a.pumpEvents(s, proc)

	// 5) 写入新一轮 prompt（stream-json envelope）
	if prompt != "" {
		if err := proc.Write(envelopeUserMessage(prompt)); err != nil {
			return &appError{code: "E_WRITE_PROMPT", msg: "write prompt: " + err.Error()}
		}
	}
	return nil
}

func (a *App) switchToTerminal(s *session.Session) error {
	if s.Owner() == session.OwnerTerminal {
		return nil
	}
	// OpenInTerminal 内部会 s.Close() 再 Launch + 标记 session
	return a.OpenInTerminal(s.WorkDir, s.ID, "")
}

// envelopeUserMessage 构造 stream-json user message envelope。
// 跟 session.Session.Send 现在的裸文本行为不同——SwitchOwner 起的是
// 全新 stream-json 进程，写 envelope 跟 Claude CLI 严格规范对齐；
// Send 走 v1 的裸文本兼容性路径（Claude 对裸 user 文本宽容接受）。
func envelopeUserMessage(prompt string) string {
	body, _ := json.Marshal(map[string]any{
		"type": "user",
		"message": map[string]any{
			"role":    "user",
			"content": prompt,
		},
	})
	return string(body) + "\n"
}

// killByPIDFile 读 pidfile 拿 pid，kill 对应进程。pid 已死（用户关
// iTerm）时不报错。Linux/macOS 用 SIGTERM 让 claude 优雅退出。
func killByPIDFile(pidfile string) error {
	data, err := os.ReadFile(pidfile)
	if err != nil {
		return err
	}
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		return &appError{code: "E_BAD_PID", msg: "invalid pid in " + pidfile}
	}
	// Unix 上 os.FindProcess 永远返回 success，但 Kill 会失败如果进程
	// 已退出；Windows 上 FindProcess 会校验存在。统一用 error 表达。
	if err := exec.Command("kill", "-TERM", strconv.Itoa(pid)).Run(); err != nil {
		// SIGTERM 失败则 SIGKILL 兜底
		_ = exec.Command("kill", "-KILL", strconv.Itoa(pid)).Run()
	}
	_ = os.Remove(pidfile)
	return nil
}

// pkillByPattern 用 pkill/taskkill 按命令行匹配杀 claude -r 进程。
// macOS/Linux: pkill -f "claude.*-r.*<sid>"
// Windows:     taskkill /F /FI "WINDOWTITLE eq Claude"
func pkillByPattern(sid string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("taskkill", "/F", "/FI", "WINDOWTITLE eq Claude").Run()
	default:
		return exec.Command("pkill", "-f", "claude.*-r.*"+sid).Run()
	}
}
