// Package cli implements the ease-ui init subcommand.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/akke/ease-ui/internal/app"
	"github.com/akke/ease-ui/internal/auth"
	"golang.org/x/term"
)

func RunInit() error {
	if auth.Exists() {
		fmt.Fprintln(os.Stderr, "ease-ui: already initialized at", auth.Path())
		fmt.Fprintln(os.Stderr, "delete the file to re-init.")
		return nil
	}

	fmt.Print("New password: ")
	pw1, err := readPassword()
	if err != nil {
		return err
	}
	fmt.Println()
	if len(pw1) < 4 {
		return fmt.Errorf("password too short (min 4 chars)")
	}
	fmt.Print("Confirm password: ")
	pw2, err := readPassword()
	if err != nil {
		return err
	}
	fmt.Println()
	if pw1 != pw2 {
		return fmt.Errorf("passwords do not match")
	}

	a, err := app.New(app.Options{})
	if err != nil {
		return err
	}
	if err := a.SetPassword(pw1); err != nil {
		return err
	}
	fmt.Println("Initialized. You can now launch the GUI with: ease-ui")
	return nil
}

func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		return string(b), err
	}
	// Fallback: read line
	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	return strings.TrimRight(line, "\r\n"), err
}

var _ = syscall.Stdin
