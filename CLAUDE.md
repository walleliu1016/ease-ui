# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## 全局约定

- 所有回复用简体中文（包括代码注释、commit message、PR 描述）。
- Go 代理：`GOPROXY=https://goproxy.cn,direct`。
- 每次 commit 前在仓库内设置 local git identity，不要依赖全局身份：
  ```bash
  git config user.name "<name>"
  git config user.email "<email>"
  ```
- 不要提交构建产物、诊断文件或运行时垃圾（如 `.claude/`、`build/bin/`、`*-err.log`、临时 `.cmd`、`.exe` 等）。

---

## 常用命令

```bash
# 后端测试（commit 前必须全绿）
go test ./...
go test ./internal/terminal -run TestNewExecCmd_StartUsesRawCmdLine -v

# 前端类型检查
cd frontend && npx vue-tsc --noEmit

# 前端 Vite dev server（无 Wails runtime，window.go 是 undefined）
cd frontend && npm run dev

# 全栈开发（推荐）
wails dev

# 生产构建
wails build -clean -trimpath -ldflags "-X main.version=v$(date +%Y%m%d)"

# 仅 Go 后端快速验证
# Go 1.26.3 + Wails v2.12.0 在当前环境生成 bindings 会失败，推荐加 --skipbindings
wails build -s --skipbindings

# 跨平台打包（必须在目标 OS 上执行）
wails build -platform darwin/universal
wails build -platform windows/amd64
wails build -platform linux/amd64
```

Go 工具链在 `~/go/bin`，PATH 经常没 export：
```bash
export PATH=$PATH:~/go/bin
```

---

## 代码风格

### Go
- 单一职责包：每个 `internal/<pkg>` 有明确边界；`internal/app` 是唯一对前端暴露的 Wails binding 层。
- 永远不要在 `internal/app` 之外的包直接 import 前端或 Wails runtime。
- 错误返回 `error`，不要 panic；Wails 的 panic 会被 runtime 吞掉，前端只会看到黑屏。
- 共享状态用 `sync.RWMutex`，不要用 channel 做状态同步。
- 测试用 `testify/assert` + `testify/require`，文件名 `<pkg>_test.go`。

### TypeScript
- `useWails.ts` 是唯一接触 `window.go` 的文件；其他文件必须 `import { X } from '../composables/useWails'`。
- 禁止直接 `window.go.app.App.X(...)`。
- Pinia 用 setup style；Vue 组件用 `<script setup lang="ts">`；路由用 hash mode。
- 样式用 `styles/theme.css` 的 CSS 变量，不要硬编码颜色。
- `Pinia ref<Record<K, V>>` 更新要用整体 spread：`state.value = { ...state.value, [id]: v }`。

---

## 架构要点（需要读多文件才能理解）

### 1. Wails binding 与 IPC
- `app.go` 配置 Wails、注册 `OnStartup/OnShutdown`、并通过 `Bind` 暴露 `internal/app` 的方法。
- 前端 `frontend/src/composables/useWails.ts` 是类型化的 IPC 转发层。
- 修改 `internal/app/*.go` 增加或删除 binding 方法后，必须跑 `wails generate` 或 `wails dev/build`，让 `frontend/wailsjs/go/app/App.js` 重新生成。**不要手改这个文件**。

### 2. Session 生命周期与 Owner
- `internal/app/session.go` 是核心 orchestrator：
  - `CreateSession`：用 `process.ModeAuto` 启动 Claude（不带 session flag），阻塞等待 `SessionStart` hook 返回真实 UUID（15s 超时）。
  - `AdoptSession`：对 Ease UI 启动前已存在的历史 session 做懒加载，第一次发消息时才起进程。
  - `SendMessage`：写 stream-json envelope。
  - `SwitchOwner`：在 `app` 与 `terminal` 之间切换写权限。
- `internal/session.Session` 维护 owner（`app` / `terminal`）和 `SwitchLock()`，切换 owner 时使用，不要与 `Send/RespondPermission` 的 `mu` 混淆。
- `process.Start` 三种 mode：
  - `ModeAuto`：Claude 自己生成 UUID（新建 session）。
  - `ModeNew`：`--session-id <sid>`，只用于全新 jsonl 不存在的情况。
  - `ModeResume`：`--resume <sid>`，jsonl 已存在时必须用它，否则 Claude 会 DEAD。

### 3. stream-json 与 envelope（关键）
- 向 Claude CLI 写用户消息必须包装成 envelope：
  ```json
  {"type":"user","message":{"role":"user","content":"..."}}\n
  ```
- 裸文本会被 stream-json 解析器直接丢弃。
- `SendMessage` 和 `SwitchOwner` 切回 app 后都写 envelope。
- `internal/protocol` 按行解析 stdout 事件；`internal/jsonl` 读历史并 watch 变化。

### 4. Hooks
- `internal/hookserver` 内置 HTTP server，监听 `127.0.0.1:<port>`，端点 `/hook` 和 `/api/send`。
- `SessionStart` **必须**用 command 脚本：`~/.ease-app/session-start.sh`（macOS/Linux）或 `session-start.ps1`（Windows）。Claude 不支持 HTTP 类型的 `SessionStart` hook。
- 其余 12 类 hook 走 HTTP POST；`hookserver` 自动写 `~/.claude/settings.json`。
- command hook 的 JSON 字段是 `hook_event_name`，不是 HTTP hook 用的 `type`。
- 前端 `handleHookEvent` 收到 `SessionStart` 时，只有 `owner.value[sid] !== 'app'` 才标记为 `terminal`，避免覆盖 Ease UI 自己新建的 session。

### 5. 外部终端
- `internal/terminal` 负责跨平台打开系统终端并运行 `claude -r <sid>`。
- **Windows**：优先 `wt`，否则 `cmd /c start`。由于 Go 1.26 + Wails 的 `exec.Command` 引号转义会被 `cmd` 误读，Windows 实现绕过标准转义：
  - 用 `syscall.SysProcAttr.CmdLine` 直接传递原生命令行。
  - 用 `COMSPEC` 环境变量定位 `cmd.exe`，避免 Wails 进程 PATH 找不到 cmd。
- **macOS**：`osascript` 调 Terminal.app / iTerm2，注入 `echo $$ > pidfile && exec claude ...`。
- **Linux**：依次尝试 `gnome-terminal`、`konsole`、`xterm`，同样注入 pidfile。
- 切回 App 时 macOS/Linux 读 pidfile kill，Windows 用 `taskkill` 按命令行匹配兜底。

### 6. 单例锁
- `internal/single` 在 `main.go` 启动时抢独占锁。
- `wails dev` 的 hot-reload 会触发新进程，因抢不到锁而失败；开发时需要先退出当前实例。
- 可临时注释 `app.go` 里的 `single.Acquire()` 改善 dev 体验，但**不要提交**。

### 7. 启动诊断
- `frontend/index.html` 内联 diag 脚本 + `frontend/src/preload.ts` + `frontend/src/main.ts` 保留启动日志，生产 build 也保留，用于现场排查。
- 黑屏时按日志定位：
  - HB 心跳停掉 = JS event loop 卡死。
  - `auth.initPromise` 之后无响应 = XPC/TCC 权限问题（macOS）。

---

## 提交规范

- 一个 task 一个 commit，格式：`<type>: <subject>`。
  type：`feat` / `fix` / `refactor` / `test` / `docs` / `chore`。
- commit 前必须 `go test ./...` 和 `npx vue-tsc --noEmit` 全绿。
- 改 binding 后同步 commit 自动生成的 `frontend/wailsjs/go/app/App.js`。
- 改 `preload.ts` / `index.html` 的诊断代码时，在 commit message 里标 **临时**。

---

## 重要不变量（改之前必须确认）

- 绝不自己生成 session ID；Claude 通过 `SessionStart` hook 返回 UUID。
- jsonl 已存在的 sid 用 `--resume`；全新 sid 用 `--session-id`。
- 前端 `formatContent` 必须与 Go 端 `ContentBlock.ContentText()` 保持同步。
- 所有 `os/exec` 调用尽量带 context：`exec.CommandContext(ctx, ...)`。
- `//go:embed all:frontend/dist` 的 `all:` 前缀不能去掉，否则隐藏文件缺失会导致某些平台启动失败。

---

## 相关文档

- `README.md` —— 项目总览、完整数据流、目录结构、已知问题。
- `docs/superpowers/specs/2026-06-27-ease-ui-design.md` —— 设计决策。
- `docs/superpowers/plans/2026-06-27-ease-ui.md` —— 44 个 task 实施历史。
