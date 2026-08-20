package cmd

import "github.com/spf13/cobra"

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the framework from AI platforms",
	RunE:  runStub,
}

func init() { rootCmd.AddCommand(uninstallCmd) }
