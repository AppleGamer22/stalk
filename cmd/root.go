package cmd

import (
	"github.com/spf13/cobra"
	// "github.com/fsnotify/fsnotify"
)

var RootCommand = &cobra.Command{
	Use:     "stalk",
	Short:   "watch/wait a file for change",
	Long:    "watch/wait a file for change",
	Version: Version,
}
