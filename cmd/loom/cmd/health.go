package cmd

import "github.com/spf13/cobra"

type healthFlags struct{ target string }

var healthArgs healthFlags

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the installed framework health",
	Args:  cobra.NoArgs,
	RunE:  runHealth,
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.Flags().StringVar(&healthArgs.target, "target", ".", "project root to check")
}
