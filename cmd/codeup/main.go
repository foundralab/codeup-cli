package main

import (
	"fmt"
	"os"

	"github.com/foundralab/codeup-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
