package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	if err := runApp(); err != nil {
		fmt.Fprintln(os.Stderr, "ease-ui:", err)
		os.Exit(1)
	}
}
