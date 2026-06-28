package terminal

import (
	"fmt"
	"os/exec"
)

func (l *Launcher) buildArgs(workDir, cmd string) []string {
	// 优先用 Windows Terminal (wt)，其次用 cmd
	if _, err := exec.LookPath("wt"); err == nil {
		return []string{
			"wt", "-w", "0", "nt",
			"--title", "Claude",
			"--startingDirectory", workDir,
			"cmd", "/k", cmd,
		}
	}
	return []string{"cmd", "/c", "start", "Claude", "cmd", "/k",
		fmt.Sprintf("cd /d %s && %s", workDir, cmd)}
}
