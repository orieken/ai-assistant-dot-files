package cmd

import "github.com/spf13/cobra"

type uninstallFlags struct {
	target   string
	platform string
	isDryRun bool
}

var uninstallArgs uninstallFlags

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the framework from AI platforms",
	Args:  cobra.NoArgs,
	RunE:  runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVar(&uninstallArgs.target, "target", ".", "project root to uninstall from")
	uninstallCmd.Flags().StringVar(&uninstallArgs.platform, "platform", "", "uninstall one AI platform")
	uninstallCmd.Flags().BoolVar(&uninstallArgs.isDryRun, "dry-run", false, "print actions without removing files")
}
