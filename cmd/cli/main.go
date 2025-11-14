package main

import (
	"fmt"
	"os"

	"github.com/turnforge/weewar/cmd/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
