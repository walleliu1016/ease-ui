# Ease UI Gateway 模式设计文档

**日期**: 2026-07-02
**作者**: Brainstorming session
**状态**: 待用户审阅（session ID 策略待最终确认）

## 1. 概述

将 Ease UI 从「hook + jsonl 实时监听」架构迁移为 **Gateway 模式**：

- ease-desktop 作为 Claude CLI 与 Anthropic API 之间的本地 HTTP 网关
- 通过改写 `ANTHROPIC_BASE_URL` 截获所有模型请求/响应
- 增量消息从网关解析 Anthropic SSE 实时推送到前端
- jsonl 退化为只读历史库，仅在首次加载、展开详情、断线兜底时读取
- 会话列表持久化到 SQLite，启动时加载到内存
- 移除「在终端打开」owner 切换、独立 hookserver、大部分 hook 配置 UI

**唯一保留的 hook**：`PermissionRequest` HTTP hook，由 gateway 自动配置，用户无感知。

## 2. 目标与非目标

### 2.1 目标

- 所有交互走 ease-desktop，不再支持切到外部终端
- 增量消息走网关 SSE 解析，实现 token 级流式渲染
- 会话列表持久化到 SQLite，启动秒开
- 去掉 jsonl fsnotify 实时监听，降低文件系统开销
- 去掉 SessionStart command shim，简化启动链路
- 统一 Gateway 内部处理 PermissionRequest

### 2.2 非目标

- 不替换 Claude CLI 为直接 Anthropic API 调用
- 不做远程/云端网关（v2 可扩展）
- 不改造 Claude CLI 的权限模型

## 3. 架构

```
┌─────────────────────────────────────────────────────────────────┐
│  Frontend (Vue 3 + Pinia)                                       │
│  useEventStream 订阅 session:{uuid} / permission:request        │
└──────────────────────┬──────────────────────────────────────────┘
                       │ Wails JSON-RPC / EventsEmit
┌──────────────────────▼──────────────────────────────────────────┐
│  internal/app                                                   │
│  Wails Bindings: SendMessage, CreateSession, RespondPermission, │
│  GetSessionMessages, ListSessions, ...                          │
└──────────────────────┬──────────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────────┐
│  internal/gateway                                               │
│  HTTP server:                                                   │
│  • /sessions/{id}/v1/messages  → 代理到 Anthropic API          │
│  • /hook                       → PermissionRequest 阻塞处理     │
│  • /api/send                   → 外部脚本写入                   │
│  解析 SSE，emit 增量事件                                         │
└──────────────────────┬──────────────────────────────────────────┘
                       │
┌──────────────────────┬──────────────────────────────────────────┐
│  internal/session    │  internal/store (SQLite)                  │
│  进程生命周期 / map    │  sessions 表：uuid, workdir, title, ...   │
│  pending permission  │  启动加载到内存                            │
└──────────────────────┴──────────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────────┐
│  internal/process                                               │
│  启动 claude --session-id <uuid> --input-format stream-json     │
│  --output-format stream-json                                    │
│  env: ANTHROPIC_BASE_URL=http://127.0.0.1:<port>/sessions/<uuid>│
└─────────────────────────────────────────────────────────────────┘
                       │
                       ▼
              Anthropic API (api.anthropic.com)
```

## 4. 关键设计决策

### 4.1 Session ID 策略（待最终确认）

**推荐方案**：ease 自生成 UUID v4，通过 `--session-id <uuid>` 传给 Claude。

理由：
- 满足「无 hook、无 jsonl 监听」约束
- `internal/process` 已存在 `ModeNew` 支持 `--session-id`
- jsonl 文件名即为该 uuid，resume 时可用
- SQLite 以 uuid 为主键，前后端一致

**风险/待验证**：
- Claude 对自生成 UUID v4 的接受度
- 已 ended session 用 `--session-id` 会立即退出（当前只用于全新会话）
- resume 行为是否与 Claude 生成 UUID 完全一致

**备选方案**：保留最小 SessionStart command shim，由 gateway 自动管理。
- 优点：不打破现有原则
- 缺点：仍有一个不可见的 hook 存在

### 4.2 网关路径设计

每个 Claude 进程启动时设置：

```bash
ANTHROPIC_BASE_URL=http://127.0.0.1:<gw-port>/sessions/<uuid>
```

Claude 请求模型时会访问：

```
POST /sessions/<uuid>/v1/messages
```

gateway 从路径提取 uuid，转发到 `https://api.anthropic.com/v1/messages`。

对于外部 Claude 进程（如果也配置了 gateway），若 uuid 不在 SQLite 中，gateway 自动创建占位会话，workdir/title 后续从 jsonl 补全。

### 4.3 消息来源矩阵

| 场景 | 来源 | 说明 |
|------|------|------|
| 点击 session 首次加载 | jsonl | 只读，加载最近 N 条 |
| 流式对话中 | Gateway SSE 解析 | 唯一实时源，token 级增量 |
| 展开长消息 / tool_result 详情 | jsonl | 读完整内容 |
| done / error 兜底 | stream-json stdout | 仅作校验，不参与主渲染 |

## 5. 组件接口

### 5.1 internal/gateway

```go
type Server struct {
    addr       string
    upstream   *url.URL
    onEvent    func(sessionID string, event GatewayEvent)
    onPermissionRequest func(req PermissionRequest) (PermissionDecision, error)
    sessions   SessionResolver
}

type GatewayEvent struct {
    Type      string          // user | assistant_delta | tool_use | tool_result | done | error
    SessionID string
    Data      json.RawMessage
}

func New(cfg Config) (*Server, error)
func (s *Server) Start() error
func (s *Server) Stop() error
```

### 5.2 internal/store

```go
type Session struct {
    ID        string
    WorkDir   string
    Title     string
    Status    string // idle | running | awaiting_permission | done
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Store interface {
    List() ([]Session, error)
    Get(id string) (Session, error)
    Insert(s Session) error
    Update(s Session) error
    Delete(id string) error
}
```

### 5.3 internal/session

```go
type Manager struct {
    sessions map[string]*Session
    store    store.Store
    gateway  *gateway.Server
    mu       sync.RWMutex
}

func (m *Manager) Create(workDir, prompt string) (string, error)
func (m *Manager) Send(sessionID, prompt string) error
func (m *Manager) RespondPermission(reqID string, allow bool) error
func (m *Manager) GetMessages(sessionID string, offset, limit int) ([]Message, error)
```

## 6. 数据流

### 6.1 启动

1. `store.LoadSessions()` 读 SQLite 到内存
2. `jsonl.SyncOnce()` 扫描 `~/.claude/projects/`，补录 SQLite 中缺失的 jsonl
3. `gateway.Start()` 启动 HTTP 代理
4. 前端 `ListSessions()` 返回内存中的会话列表

### 6.2 新建会话

1. 生成 UUID v4（或等待 SessionStart，若选备选方案）
2. `session.InsertSQLite()` 持久化
3. `memorySessions[uuid] = session`
4. `process.Start(workDir, uuid, bin, ModeNew, env={ANTHROPIC_BASE_URL: gatewayURL(uuid)})`
5. `proc.Write(envelopeUserMessage(prompt))`
6. 返回 uuid 给前端

### 6.3 发送消息

1. 前端 `SendMessage(uuid, prompt)`
2. Go 端查找 session，写 stream-json envelope 到 stdin
3. Claude 构造 API 请求到 `/sessions/<uuid>/v1/messages`
4. Gateway 提取 uuid，emit `gateway:user` 事件
5. Gateway 转发到 Anthropic API

### 6.4 接收模型响应

1. Anthropic 返回 SSE 流
2. Gateway 边转发给 Claude，边解析每个 event
3. 按 event 类型 emit：
   - `assistant_delta`：文本/thinking 增量
   - `tool_use`：工具调用开始
   - `tool_result`：从下一次请求体中提取
   - `done`：stream 结束且无 tool_use
4. 前端 `useEventStream` 接收并增量渲染

### 6.5 PermissionRequest

1. Claude 执行敏感工具前触发 `PermissionRequest` HTTP hook
2. POST 到 gateway `/hook`
3. Gateway 生成 reqID，存入 pending map，emit `permission:request`
4. 前端弹窗，用户决策后调用 `RespondPermission`
5. Gateway 返回 decision HTTP response
6. Claude 继续/停止

### 6.6 历史 / 展开

1. 用户点击展开或切换回旧 session
2. `GetSessionMessages(uuid, offset, limit)`
3. 读 `~/.claude/projects/<proj>/<uuid>.jsonl`
4. 返回完整消息列表

## 7. SQLite 表结构

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,        -- uuid
    workdir TEXT NOT NULL,
    title TEXT,
    status TEXT DEFAULT 'idle', -- idle | running | awaiting_permission | done
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE INDEX idx_sessions_updated ON sessions(updated_at DESC);
```

## 8. 事件类型

前端订阅的事件：

```ts
type GatewayEvent =
  | { type: 'user'; sessionId: string; message: Message }
  | { type: 'assistant_delta'; sessionId: string; delta: TextDelta }
  | { type: 'tool_use'; sessionId: string; tool: ToolUse }
  | { type: 'tool_result'; sessionId: string; result: ToolResult }
  | { type: 'done'; sessionId: string }
  | { type: 'error'; sessionId: string; error: string }
```

## 9. 错误处理

| 场景 | 级别 | 处理 |
|------|------|------|
| Gateway 启动端口被占用 | 致命 | 提示用户重启，写日志 |
| Claude 进程启动失败 | 可恢复 | toast，状态设为 error |
| jsonl 读取失败 | 可恢复 | 跳过，toast 提示 |
| SSE 解析失败 | 静默 | 记录日志，继续转发 |
| UUID 生成/传递验证失败 | 致命（开发期） | 明确错误，阻止发布 |
| PermissionRequest 超时 | 可恢复 | 默认 deny，toast |

## 10. 测试策略

### 10.1 Go 后端

- `internal/gateway`：mock Anthropic API，验证路径解析、SSE 转发、事件 emit
- `internal/store`：SQLite 读写、并发、迁移
- `internal/session`：创建/发送/权限/恢复流程
- `internal/process`：参数构造、环境变量注入

### 10.2 前端

- store 测试：gateway 事件处理、uuid 切换、jsonl 加载兜底
- 组件测试：流式渲染、tool_use 块、PermissionPanel

### 10.3 集成

- 真实 Claude CLI + gateway，验证 `--session-id <uuid>` 可行
- 验证 resume 行为
- 验证 PermissionRequest 拦截

## 11. 风险与待决策

1. **Session ID 自生成**：需验证 Claude 接受 UUID v4 的 `--session-id`，以及 resume 行为。
2. **jsonl 首次加载时机**：若用户点击 session 时 Claude 进程刚启动、jsonl 尚未写入，需 loading/空态处理。
3. **Gateway 与 Claude 进程绑定**：单 gateway 端口处理多个 session，依赖路径前缀区分，无状态。
4. **外部 Claude 进程**：若用户手动配置 `ANTHROPIC_BASE_URL` 指向 gateway，可能产生未知 session，需占位和补全逻辑。
5. **PermissionRequest 仍是 hook**：这是 Claude CLI 的硬约束，无法通过网关替代。

## 12. 后续计划

1. 用户确认 session ID 策略（自生成 UUID vs 保留 SessionStart shim）
2. 调用 `superpowers:writing-plans` 制定实施计划
3. 按组件拆分 task：gateway、store、session 改造、前端事件迁移
