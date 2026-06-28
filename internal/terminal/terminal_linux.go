package terminal

import (
	"fmt"
	"os/exec"
)

func (l *Launcher) buildArgs(workDir, cmd string) []string {
	// 依次尝试主流终端
	terms := []struct {
		bin  string
		args func() []string
	}{
		{"gnome-terminal", func() []string {
			return []string{"gnome-terminal", "--working-directory=" + workDir, "--", "bash", "-c", cmd + "; exec bash"}
		}},
		{"konsole", func() []string {
			return []string{"konsole", "--workdir", workDir, "-e", "bash", "-c", cmd + "; exec bash"}
		}},
		{"xterm", func() []string {
			return []string{"xterm", "-e", fmt.Sprintf("cd %s && %s; bash", workDir, cmd)}
		}},
	}

	for _, t := range terms {
		if _, err := exec.LookPath(t.bin); err == nil {
			return t.args()
		}
	}
	// 兜底：x-terminal-emulator
	return []string{"x-terminal-emulator", "-e", fmt.Sprintf("cd %s && %s; bash", workDir, cmd)}
}
