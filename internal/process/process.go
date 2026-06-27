// Package process wraps the Claude CLI as a managed subprocess.
package process

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/akke/ease-ui/internal/protocol"
)

type Process struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	events chan []byte
	done   chan struct{}
	mu     sync.Mutex
	closed bool
}

// Start launches the claude binary with the given workdir and session ID.
// Output is captured line-by-line into the Events channel as raw bytes
// (caller is responsible for protocol.Parse).
//
// If binPath is empty, "claude" is used and the standard Claude CLI flags
// (--cwd, --session-id, --output-format stream-json, --verbose) are appended.
// If binPath is an explicit path (e.g. /bin/cat in tests), no flags are added
// so the helper binary is invoked as-is.
func Start(workDir, sessionID, binPath string) (*Process, error) {
	var args []string
	if binPath == "" {
		binPath = "claude"
		args = []string{
			"--cwd", workDir,
			"--session-id", sessionID,
			"--output-format", "stream-json",
			"--verbose",
		}
	}
	cmd := exec.Command(binPath, args...)
	cmd.Dir = workDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	p := &Process{
		cmd:    cmd,
		stdin:  stdin,
		events: make(chan []byte, 256),
		done:   make(chan struct{}),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start %s: %w", binPath, err)
	}

	go p.readLoop(stdout)
	go p.waitLoop()

	return p, nil
}

func (p *Process) Events() <-chan []byte { return p.events }

func (p *Process) Write(s string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return fmt.Errorf("process: closed")
	}
	_, err := io.WriteString(p.stdin, s)
	return err
}

func (p *Process) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	_ = p.stdin.Close()
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
	<-p.done
	return nil
}

func (p *Process) readLoop(r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())
		select {
		case p.events <- line:
		case <-p.done:
			return
		}
	}
	close(p.events)
}

func (p *Process) waitLoop() {
	_ = p.cmd.Wait()
	close(p.done)
}

// SendParsed is a convenience that writes a line for Claude's user-message input.
func (p *Process) SendUserPrompt(prompt string) error {
	return p.Write(prompt + "\n")
}

var _ = protocol.EvtMessage
