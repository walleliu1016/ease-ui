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
	assert.Empty(t, cfg.PreToolUse)
}

func TestEditor_SaveCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"PreToolUse":[]}`), 0o644))

	SetPath(path)
	e := NewEditor()
	cfg := &Config{}
	cfg.PreToolUse = []Hook{{Command: "echo hi", Type: HookTypeShell}}
	require.NoError(t, e.Save(cfg))

	_, err := os.Stat(path + ".bak")
	assert.NoError(t, err)
}

func TestEditor_SavePreservesUnknownFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	original := `{"theme":"dark","PreToolUse":[{"command":"x","type":"shell"}]}`
	require.NoError(t, os.WriteFile(path, []byte(original), 0o644))

	SetPath(path)
	e := NewEditor()
	cfg, err := e.Load()
	require.NoError(t, err)
	cfg.PostToolUse = []Hook{{Command: "y", Type: HookTypeShell}}
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
	cfg := &Config{}
	cfg.PreToolUse = []Hook{{Command: "ls", Type: HookTypeShell}}
	require.NoError(t, e.Save(cfg))

	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err))
}
