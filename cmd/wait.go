package cmd

import (
	"github.com/spf13/cobra"
)

var waitCommand = &cobra.Command{
	Short: "wait a file for change",
	Long:  "wait a file for change",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
