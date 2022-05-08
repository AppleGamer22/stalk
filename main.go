package main

import (
	"os"

	"github.com/AppleGamer22/stalk/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
