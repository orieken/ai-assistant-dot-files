package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0-dev"

var rootCmd = &cobra.Command{
	Use:     "loom",
	Short:   "AI assistant framework installer",
	Long:    "Loom installs agents, skills, and rules for Claude Code, Cursor, Windsurf, and other AI platforms.",
	Version: version,
}

// Execute runs the loom command-line interface.
func Execute() {
	rootCmd.AddCommand(installCmd, uninstallCmd, updateCmd, healthCmd, versionCmd)
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	cobra.CheckErr(rootCmd.Execute())
}

func runStub(_ *cobra.Command, _ []string) error {
	fmt.Println("not yet implemented")
	return nil
}
