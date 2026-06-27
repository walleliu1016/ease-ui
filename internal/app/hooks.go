package app

import (
	"os"
	"path/filepath"

	"github.com/akke/ease-ui/internal/hooks"
)

// HooksConfig is the frontend-facing shape (no need to expose internal types).
type HooksConfig struct {
	PreToolUse        []HookEntry `json:"PreToolUse"`
	PermissionRequest []HookEntry `json:"PermissionRequest"`
	PostToolUse       []HookEntry `json:"PostToolUse"`
	Notification      []HookEntry `json:"Notification"`
	Stop              []HookEntry `json:"Stop"`
}

type HookEntry struct {
	Matcher string `json:"matcher,omitempty"`
	Command string `json:"command"`
	Type    string `json:"type"`
}

func (a *App) GetHooksConfig() (*HooksConfig, error) {
	cfg, err := hooks.NewEditor().Load()
	if err != nil {
		return nil, err
	}
	return toHooksConfig(cfg), nil
}

func (a *App) SaveHooksConfig(hc *HooksConfig) error {
	cfg := fromHooksConfig(hc)
	// Ensure parent dir exists; hooks.Editor.Save does not create it.
	if err := os.MkdirAll(filepath.Dir(hooks.Path()), 0o755); err != nil {
		return err
	}
	return hooks.NewEditor().Save(cfg)
}

func toHooksConfig(c *hooks.Config) *HooksConfig {
	return &HooksConfig{
		PreToolUse:        toEntries(c.PreToolUse),
		PermissionRequest: toEntries(c.PermissionRequest),
		PostToolUse:       toEntries(c.PostToolUse),
		Notification:      toEntries(c.Notification),
		Stop:              toEntries(c.Stop),
	}
}

func fromHooksConfig(hc *HooksConfig) *hooks.Config {
	return &hooks.Config{
		EventHooks: hooks.EventHooks{
			PreToolUse:        fromEntries(hc.PreToolUse),
			PermissionRequest: fromEntries(hc.PermissionRequest),
			PostToolUse:       fromEntries(hc.PostToolUse),
			Notification:      fromEntries(hc.Notification),
			Stop:              fromEntries(hc.Stop),
		},
	}
}

func toEntries(in []hooks.Hook) []HookEntry {
	out := make([]HookEntry, 0, len(in))
	for _, h := range in {
		out = append(out, HookEntry{Matcher: h.Matcher, Command: h.Command, Type: string(h.Type)})
	}
	return out
}

func fromEntries(in []HookEntry) []hooks.Hook {
	out := make([]hooks.Hook, 0, len(in))
	for _, h := range in {
		out = append(out, hooks.Hook{Matcher: h.Matcher, Command: h.Command, Type: hooks.HookType(h.Type)})
	}
	return out
}
