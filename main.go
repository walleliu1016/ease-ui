package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		// CLI init handled in Task 21
		fmt.Println("init placeholder")
		return
	}
	fmt.Println("ease-ui", version)
}
