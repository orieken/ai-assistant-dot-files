package cmd

import "github.com/spf13/cobra"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the installed framework",
	RunE:  runStub,
}

func init() { rootCmd.AddCommand(updateCmd) }
