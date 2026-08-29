package cmd

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/spf13/cobra"
)

type toolsStatusFlags struct{ tier string }

var toolsStatusArgs toolsStatusFlags

var toolsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show which context tools are installed",
	Args:  cobra.NoArgs,
	RunE:  runToolsStatus,
}

func init() {
	toolsCmd.AddCommand(toolsStatusCmd)
	toolsStatusCmd.Flags().StringVar(&toolsStatusArgs.tier, "tier", "all", "filter by tier: high, medium, or all")
}

func runToolsStatus(command *cobra.Command, _ []string) error {
	return executeToolsStatus(toolsStatusArgs, contextTools, exec.LookPath, command.OutOrStdout())
}

func executeToolsStatus(flags toolsStatusFlags, tools []contextTool, lookPath func(string) (string, error), writer io.Writer) error {
	out := toolsOutput{writer: writer}
	out.heading("loom tools: context and memory optimization tools")

	tier, err := parseTier(flags.tier)
	if err != nil {
		return err
	}

	var lastTier toolTier
	missing := 0

	for _, tool := range tools {
		if tier != "" && tool.tier != tier {
			continue
		}
		lastTier = reportTierHeading(out, tool.tier, lastTier)
		if !reportToolRow(out, tool, lookPath) {
			missing++
		}
	}

	out.separator()
	out.statusSummary(missing)
	return nil
}

func reportTierHeading(out toolsOutput, current, last toolTier) toolTier {
	if current == last {
		return last
	}
	if last != "" {
		out.separator()
	}
	out.sectionHeading(current)
	return current
}

func reportToolRow(out toolsOutput, tool contextTool, lookPath func(string) (string, error)) bool {
	installed, installNote := toolInstallStatus(tool, lookPath)
	status := "✓"
	if !installed {
		status = "✗"
	}
	out.toolRow(tool.name, status, tool.description)
	if !installed && installNote != "" {
		out.installNote(installNote)
	}
	if !installed && tool.postNote != "" {
		out.postNote(tool.postNote)
	}
	return installed
}

// toolInstallStatus returns (installed bool, installHint string).
func toolInstallStatus(tool contextTool, lookPath func(string) (string, error)) (bool, string) {
	if tool.binary == "" {
		return true, ""
	}
	if _, err := lookPath(tool.binary); err == nil {
		return true, ""
	}
	if tool.manualNote != "" {
		return false, tool.manualNote
	}
	if len(tool.installCmds) > 0 {
		return false, tool.installCmds[0]
	}
	return false, ""
}

// parseTier converts a tier string flag to a toolTier, returning "" for "all".
func parseTier(s string) (toolTier, error) {
	switch s {
	case "all", "":
		return "", nil
	case string(tierHigh):
		return tierHigh, nil
	case string(tierMedium):
		return tierMedium, nil
	default:
		return "", fmt.Errorf("unknown tier %q — valid values: high, medium, all", s)
	}
}
