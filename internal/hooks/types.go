// Package hooks provides settings.json hooks editing and permission handling.
package hooks

import "encoding/json"

// Config 是 settings.json 中 hooks 字段的完整 JSON 结构。
// 使用 map[string]any 保留 Claude 原生格式，前端负责展示和校验。
type Config struct {
	Hooks map[string]any `json:"hooks"`
}

// --- 以下为权限处理用类型（handler.go 使用）---

type Decision struct {
	Allow  bool
	Auto   bool
	Reason string
}

type PermissionRequest struct {
	RequestID string          `json:"request_id"`
	Tool      string          `json:"tool"`
	Args      json.RawMessage `json:"args"`
}
