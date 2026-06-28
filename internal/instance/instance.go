// Package instance 持久化会话状态到 ~/.ease-app/instance.json。
// 检测到新增或变更时自动写入。
package instance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// State 保存单个会话的持久化状态。
type State struct {
	State    string    `json:"state"`     // idle, running, awaiting_permission, done
	LastSeen time.Time `json:"last_seen"`
}

// Store 管理 instance.json。
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]State
}

var store *Store

// Load 读取 ~/.ease-app/instance.json，不存在则返回空。
func Load() (*Store, error) {
	if store != nil {
		return store, nil
	}
	home, _ := os.UserHomeDir()
	p := filepath.Join(home, ".ease-app", "instance.json")

	s := &Store{path: p, data: map[string]State{}}
	data, err := os.ReadFile(p)
	if err == nil {
		_ = json.Unmarshal(data, &s.data)
	}
	store = s
	return s, nil
}

// Get 返回指定会话的状态（不存在返回空 State）。
func (s *Store) Get(sessionID string) State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[sessionID]
}

// Delete 从持久化中移除指定会话。
func (s *Store) Delete(sessionID string) {
	s.mu.Lock()
	delete(s.data, sessionID)
	s.mu.Unlock()
	_ = s.write()
}

// Put 更新状态并写入文件。
func (s *Store) Put(sessionID string, state string) {
	s.mu.Lock()
	prev, ok := s.data[sessionID]
	now := time.Now()
	if ok && prev.State == state {
		// 状态未变，仅更新 last_seen 频率限制（30s 内不重复写）
		if now.Sub(prev.LastSeen) < 30*time.Second {
			s.mu.Unlock()
			return
		}
	}
	s.data[sessionID] = State{State: state, LastSeen: now}
	s.mu.Unlock()
	_ = s.write()
}

// LastSeen 更新最近活动时间（不改变状态）。
func (s *Store) Touch(sessionID string) {
	s.mu.Lock()
	v, ok := s.data[sessionID]
	if !ok {
		s.mu.Unlock()
		return
	}
	v.LastSeen = time.Now()
	s.data[sessionID] = v
	s.mu.Unlock()
	_ = s.write()
}

func (s *Store) write() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p := s.path
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}
