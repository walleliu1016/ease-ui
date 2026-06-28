package app

import "github.com/akke/ease-ui/internal/hooks"

// GetHooksConfig 返回 settings.json 中 hooks 字段的原始 JSON，
// 前端负责展示和编辑所有 hook 类型。
func (a *App) GetHooksConfig() (map[string]any, error) {
	cfg, err := hooks.NewEditor().Load()
	if err != nil {
		return nil, err
	}
	if cfg.Hooks == nil {
		return map[string]any{}, nil
	}
	return cfg.Hooks, nil
}

// SaveHooksConfig 将 hooks JSON 写回 settings.json，保留其他字段不变。
func (a *App) SaveHooksConfig(hooksData map[string]any) error {
	cfg, err := hooks.NewEditor().Load()
	if err != nil {
		return err
	}
	cfg.Hooks = hooksData
	return hooks.NewEditor().Save(cfg)
}

// 清理旧类型（已经不需要 toHooksConfig / fromHooksConfig 转换了，
// 前端直接拿到 JSON 后自行解析）
