package terminal

import "fmt"

func (l *Launcher) buildArgs(workDir, cmd string) []string {
	script := fmt.Sprintf(`tell application "Terminal" to do script "cd %s && %s"`, workDir, cmd)
	return []string{"osascript", "-e", script}
}
