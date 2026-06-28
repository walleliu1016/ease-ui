package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsApp_UpdateAndReload(t *testing.T) {
	dir := t.TempDir()
	a, err := New(Options{ConfigDir: dir})
	require.NoError(t, err)

	cfg, err := a.GetSettings()
	require.NoError(t, err)
	cfg.AutoAllowBash = true
	require.NoError(t, a.UpdateSettings(cfg))

	cfg2, _ := a.GetSettings()
	assert.True(t, cfg2.AutoAllowBash)
}

func TestHooksApp_AddAndSave(t *testing.T) {
	dir := t.TempDir()
	a, err := New(Options{ConfigDir: dir})
	require.NoError(t, err)

	// 新接口：map[string]any，对应 Claude settings.json 原始 hooks 格式
	cfg, err := a.GetHooksConfig()
	require.NoError(t, err)

	cfg["PreToolUse"] = []any{
		map[string]any{
			"matcher": "Bash",
			"hooks": []any{
				map[string]any{
					"type":    "command",
					"command": "echo hi",
				},
			},
		},
	}
	require.NoError(t, a.SaveHooksConfig(cfg))

	loaded, _ := a.GetHooksConfig()
	pretool, ok := loaded["PreToolUse"].([]any)
	require.True(t, ok)
	require.Len(t, pretool, 1)
}
