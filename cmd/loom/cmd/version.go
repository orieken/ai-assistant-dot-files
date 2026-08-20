package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show loom version information",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println(version)
		return nil
	},
}

func init() { rootCmd.AddCommand(versionCmd) }
