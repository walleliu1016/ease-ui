package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditor_LoadMissingReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	SetPath(filepath.Join(dir, "settings.json"))

	e := NewEditor()
	cfg, err := e.Load()
	require.NoError(t, err)
	assert.Nil(t, cfg.Hooks)
}

func TestEditor_SaveCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"hooks":{"PreToolUse":[]}}`), 0o644))

	SetPath(path)
	e := NewEditor()
	cfg := &Config{
		Hooks: map[string]any{
			"PreToolUse": []any{
				map[string]any{
					"hooks": []any{
						map[string]any{"type": "command", "command": "echo hi"},
					},
				},
			},
		},
	}
	require.NoError(t, e.Save(cfg))

	_, err := os.Stat(path + ".bak")
	assert.NoError(t, err)
}

func TestEditor_SavePreservesUnknownFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	original := `{"theme":"dark","hooks":{"PreToolUse":[{"hooks":[{"type":"command","command":"x"}]}]}}`
	require.NoError(t, os.WriteFile(path, []byte(original), 0o644))

	SetPath(path)
	e := NewEditor()
	cfg, err := e.Load()
	require.NoError(t, err)
	cfg.Hooks["PostToolUse"] = []any{
		map[string]any{
			"hooks": []any{
				map[string]any{"type": "command", "command": "y"},
			},
		},
	}
	require.NoError(t, e.Save(cfg))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"theme":"dark"`)
}

func TestEditor_SaveIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	require.NoError(t, os.WriteFile(path, []byte(`{}`), 0o644))

	SetPath(path)
	e := NewEditor()
	cfg := &Config{
		Hooks: map[string]any{
			"PreToolUse": []any{
				map[string]any{
					"hooks": []any{
						map[string]any{"type": "command", "command": "ls"},
					},
				},
			},
		},
	}
	require.NoError(t, e.Save(cfg))

	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err))
}
