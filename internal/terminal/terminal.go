// Package terminal launches the system terminal with a claude -r command.
package terminal

import (
	"fmt"
	"os/exec"
)

func ResumeCommand(sessionID, binPath string) string {
	if binPath == "" {
		binPath = "claude"
	}
	return fmt.Sprintf("%s -r %s", binPath, sessionID)
}

type Launcher struct{}

func New() *Launcher { return &Launcher{} }

// Launch spawns the system terminal running claude -r in the given workdir.
func (l *Launcher) Launch(workDir, sessionID, binPath string) error {
	cmdStr := ResumeCommand(sessionID, binPath)
	args := l.buildArgs(workDir, cmdStr)
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Start()
}
