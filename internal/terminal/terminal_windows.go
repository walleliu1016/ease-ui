package terminal

import "fmt"

func (l *Launcher) buildArgs(workDir, cmd string) []string {
	return []string{"cmd", "/c", "start", "cmd", "/k", fmt.Sprintf("cd /d %s && %s", workDir, cmd)}
}
