package app

import (
	"github.com/akke/ease-ui/internal/terminal"
)

func (a *App) OpenInTerminal(workDir, sessionID, _binPath string) error {
	// 先关闭当前进程（如果有），释放 stdin
	if s, ok := a.lookupSession(sessionID); ok {
		_ = s.Close()
	}
	a.appMu.RLock()
	bin := a.claudeBin
	if bin == "" {
		bin = a.settings.ClaudePath
	}
	a.appMu.RUnlock()
	return terminal.New().Launch(workDir, sessionID, bin)
}
