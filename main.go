package main

import (
	"fmt"
	"os"

	"github.com/akke/ease-ui/internal/cli"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := cli.RunInit(); err != nil {
			fmt.Fprintln(os.Stderr, "ease-ui init:", err)
			os.Exit(1)
		}
		return
	}
	if err := runApp(); err != nil {
		fmt.Fprintln(os.Stderr, "ease-ui:", err)
		os.Exit(1)
	}
}
