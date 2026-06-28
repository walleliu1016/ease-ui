package terminal

import (
	"fmt"
	"os/exec"
)

func (l *Launcher) buildArgs(workDir, cmd string) []string {
	// 优先用 iTerm2，其次 Terminal.app，最后用 open 命令
	termApp := "Terminal"
	if _, err := exec.LookPath("iTerm"); err == nil {
		termApp = "iTerm"
	}
	script := fmt.Sprintf(`
		tell application "%s"
			activate
			do script "%s"
		end tell
	`, termApp, fmt.Sprintf("cd %s && %s", workDir, cmd))
	return []string{"osascript", "-e", script}
}
