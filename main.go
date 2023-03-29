package main

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:           "stalk",
	Short:         "watch/wait a file for change",
	Long:          "watch/wait a file for change",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
}

func init() {
	RootCommand.SetVersionTemplate("{{.Version}}\n")
}

func main() {
	if err := RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
