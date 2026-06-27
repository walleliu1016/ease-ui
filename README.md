# Ease UI

跨平台 Claude 会话管理桌面 App (Wails v2 + Vue 3 + Go).

## 功能

- 本地密码登录（bcrypt + 锁定保护）
- 扫描并展示所有 Claude 会话
- 双向 CLI 包装：App 内发送 prompt、流式接收响应
- 拦截 `PermissionRequest`，App 内批准/拒绝
- 编辑 `~/.claude/settings.json` 中的 Hook 配置（带备份）
- 一键 `claude -r <session-id>` 打开系统终端

## 开发

```bash
# 后端测试
go test ./...

# 前端开发
cd frontend
pnpm install
pnpm dev

# 全栈开发（wails 会同时启 Vite + Go）
wails dev

# 打包（仅本机）
wails build -clean -trimpath -ldflags "-X main.version=v0.1.0"
```

## 首次使用

```bash
# 1. 初始化密码
./ease-ui init

# 2. 启动 GUI
./ease-ui
```

> v1 不在 App 内引导设置密码。删除 `~/.ease-ui/auth.json` 可重置。

## 项目结构

```
ease-ui/
├── main.go              # 入口（dispatch init vs GUI）
├── app.go               # Wails 启动 glue
├── internal/
│   ├── app/             # Wails 绑定（唯一对前端暴露的层）
│   ├── session/         # 会话领域 + 状态机
│   ├── process/         # Claude CLI 子进程
│   ├── protocol/        # stream-json 解析
│   ├── jsonl/           # jsonl 扫描 + 解析
│   ├── hooks/           # Hooks 处理 + 编辑器
│   ├── terminal/        # 系统终端调用
│   ├── auth/            # 本地密码鉴权
│   ├── settings/        # 用户设置
│   ├── events/          # 进程内 pub/sub
│   ├── log/             # 日志（lumberjack 轮转）
│   └── cli/             # ease-ui init 子命令
├── frontend/            # Vue 3 + TypeScript
└── docs/superpowers/    # 设计文档与实施计划
```

## 配置文件位置

| 文件 | 路径 |
|---|---|
| Ease 配置 | `~/.ease-ui/config.json` |
| 鉴权 | `~/.ease-ui/auth.json` |
| 日志 | `~/.ease-ui/logs/*.log` |
| Hook 配置（可写） | `~/.claude/settings.json` |
| 会话数据（只读） | `~/.claude/projects/*/<sid>.jsonl` |

## 打包

- macOS: `wails build -platform darwin/universal`（必须在 Mac 上）
- Windows: `wails build -platform windows/amd64`（Windows / mingw-w64）
- Linux: `wails build -platform linux/amd64`（Linux / Docker）

不做 GitHub Actions release workflow，每个平台手工打包。
