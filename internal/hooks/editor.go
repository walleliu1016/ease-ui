package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type Editor struct{}

func NewEditor() *Editor { return &Editor{} }

var (
	pathMu sync.RWMutex
	path   = defaultSettingsPath()
)

func defaultSettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "settings.json")
}

func SetPath(p string) {
	pathMu.Lock()
	defer pathMu.Unlock()
	path = p
}

func Path() string {
	pathMu.RLock()
	defer pathMu.RUnlock()
	return path
}

func (e *Editor) Load() (*Config, error) {
	p := Path()
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config back to settings.json, preserving unknown top-level
// fields, rotating up to 5 .bak files, and using atomic write (tmp + rename).
func (e *Editor) Save(cfg *Config) error {
	p := Path()

	// Ensure parent dir exists so callers don't have to (mirrors settings.Save).
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}

	// Read existing raw to merge unknown fields
	existing := map[string]json.RawMessage{}
	if data, err := os.ReadFile(p); err == nil {
		_ = json.Unmarshal(data, &existing)
	}

	// Marshal new hook fields
	newData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	newMap := map[string]json.RawMessage{}
	_ = json.Unmarshal(newData, &newMap)

	// Merge: newMap wins on known keys
	for k, v := range newMap {
		existing[k] = v
	}

	out, err := json.Marshal(existing)
	if err != nil {
		return err
	}

	// Rotate backups
	if _, err := os.Stat(p); err == nil {
		_ = rotateBackups(p, 5)
	}

	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func rotateBackups(p string, max int) error {
	for i := max - 1; i >= 1; i-- {
		old := fmt.Sprintf("%s.bak.%d", p, i)
		next := fmt.Sprintf("%s.bak.%d", p, i+1)
		if _, err := os.Stat(old); err == nil {
			if err := os.Rename(old, next); err != nil {
				return err
			}
		}
	}
	if _, err := os.Stat(p + ".bak"); err == nil {
		if err := os.Rename(p+".bak", p+".bak.1"); err != nil {
			return err
		}
	}
	return os.Rename(p, p+".bak")
}

// ListBackups returns backup file paths newest first.
func ListBackups() ([]string, error) {
	p := Path()
	dir := filepath.Dir(p)
	base := filepath.Base(p)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		name := e.Name()
		if len(name) > len(base) && name[:len(base)] == base {
			out = append(out, filepath.Join(dir, name))
		}
	}
	sort.Strings(out)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}
