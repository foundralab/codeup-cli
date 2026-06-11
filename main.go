package main

import (
	"fmt"
	"os"

	"github.com/foundralab/codeup-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
