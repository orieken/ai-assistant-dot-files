package cmd

import "github.com/spf13/cobra"

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the installed framework health",
	RunE:  runStub,
}
