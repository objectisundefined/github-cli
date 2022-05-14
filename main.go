package main

import (
	"fmt"

	"github.com/objectisundefined/github-cli/internal/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
	}
}
