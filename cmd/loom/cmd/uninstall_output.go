package cmd

import (
	"fmt"
	"io"
)

type uninstallOutput struct {
	writer   io.Writer
	isDryRun bool
}

func (output uninstallOutput) heading(target string) {
	prefix := ""
	if output.isDryRun {
		prefix = "[dry-run] "
	}
	fmt.Fprintf(output.writer, "%sloom: uninstalling framework from %s\n\n", prefix, target)
}

func (output uninstallOutput) removed(path string) {
	verb := "removed"
	if output.isDryRun {
		verb = "would remove"
	}
	fmt.Fprintf(output.writer, "  ✓ %s %s\n", verb, path)
}

func (output uninstallOutput) missing(path string) {
	fmt.Fprintf(output.writer, "  ⚠ skipped %s (no longer exists)\n", path)
}

func (output uninstallOutput) summary(removed int, manifestAction string) {
	verb := "removed"
	if output.isDryRun {
		verb = "would be removed"
	}
	fmt.Fprintf(output.writer, "\n%d paths %s. Manifest %s.\n", removed, verb, manifestAction)
}
