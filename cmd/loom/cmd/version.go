package cmd

import (
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"strings"

	"github.com/orieken/loom/cmd/loom/internal/platform"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show loom version information",
	Args:  cobra.NoArgs,
	RunE:  runVersion,
}

func init() { rootCmd.AddCommand(versionCmd) }

func runVersion(command *cobra.Command, _ []string) error {
	return writeVersion(command.OutOrStdout(), frameworkFS, version, commit, date)
}

func writeVersion(writer io.Writer, content platform.Content, binaryVersion, revision, builtAt string) error {
	frameworkVersion, err := readFrameworkVersion(content)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "loom %s (commit %s, built %s)\nframework %s (embedded)\n", binaryVersion, revision, builtAt, frameworkVersion); err != nil {
		return fmt.Errorf("write version output: %w", err)
	}
	return nil
}

func readFrameworkVersion(content platform.Content) (string, error) {
	if content == nil {
		return "embedded", nil
	}
	value, err := content.ReadFile("shared/VERSION")
	if errors.Is(err, iofs.ErrNotExist) {
		return "embedded", nil
	}
	if err != nil {
		return "", fmt.Errorf("read embedded framework version: %w", err)
	}
	return strings.TrimSpace(string(value)), nil
}
