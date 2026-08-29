package cmd

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/spf13/cobra"
)

type toolsInstallFlags struct {
	tier string
	all  bool
}

var toolsInstallArgs toolsInstallFlags

var toolsInstallCmd = &cobra.Command{
	Use:   "install [tool]",
	Short: "Install opt-in context tools",
	Long: `Install opt-in context and memory optimization tools.

Without arguments, installs all high-tier tools that are automatable (brew/npm/pip).
Tools that require a manual install step (such as ctx) print the command to run instead.

Examples:
  loom tools install                   # install all high-tier tools
  loom tools install --tier medium     # install medium-tier tools
  loom tools install --all             # install all tiers
  loom tools install tokei             # install one tool by name`,
	Args: cobra.MaximumNArgs(1),
	RunE: runToolsInstall,
}

func init() {
	toolsCmd.AddCommand(toolsInstallCmd)
	toolsInstallCmd.Flags().StringVar(&toolsInstallArgs.tier, "tier", "high", "tier to install: high, medium, or all")
	toolsInstallCmd.Flags().BoolVar(&toolsInstallArgs.all, "all", false, "install all tiers (equivalent to --tier all)")
}

func runToolsInstall(command *cobra.Command, args []string) error {
	return executeToolsInstall(toolsInstallArgs, args, contextTools, exec.LookPath, runInstallCmd, command.OutOrStdout())
}

type lookPathFn func(string) (string, error)
type runCmdFn func(string) error

func executeToolsInstall(
	flags toolsInstallFlags,
	args []string,
	tools []contextTool,
	lookPath lookPathFn,
	runCmd runCmdFn,
	writer io.Writer,
) error {
	out := toolsOutput{writer: writer}

	candidates, err := selectTools(flags, args, tools)
	if err != nil {
		return err
	}

	out.installHeading(len(candidates))

	counts := map[installOutcome]int{}
	for _, tool := range candidates {
		counts[installOneTool(tool, lookPath, runCmd, out)]++
	}

	out.installSummary(counts[outcomeInstalled], counts[outcomeSkipped], counts[outcomeFailed], counts[outcomeManual])
	if counts[outcomeFailed] > 0 {
		return fmt.Errorf("%d tool install(s) failed", counts[outcomeFailed])
	}
	return nil
}

type installOutcome int

const (
	outcomeInstalled installOutcome = iota
	outcomeSkipped
	outcomeManual
	outcomeFailed
)

func installOneTool(tool contextTool, lookPath lookPathFn, runCmd runCmdFn, out toolsOutput) installOutcome {
	if tool.binary == "" {
		out.installBuiltIn(tool.name)
		return outcomeSkipped
	}
	if _, err := lookPath(tool.binary); err == nil {
		out.installSkip(tool.name)
		return outcomeSkipped
	}
	if tool.manualNote != "" {
		out.installManual(tool.name, tool.manualNote)
		reportPostNote(out, tool)
		return outcomeManual
	}
	if len(tool.installCmds) == 0 {
		out.installManual(tool.name, "no automated install available — see shared/knowledge/"+tool.kiName+".md")
		return outcomeManual
	}
	if err := runCmd(tool.installCmds[0]); err != nil {
		out.installFailed(tool.name, err)
		return outcomeFailed
	}
	out.installSuccess(tool.name)
	reportPostNote(out, tool)
	return outcomeInstalled
}

func reportPostNote(out toolsOutput, tool contextTool) {
	if tool.postNote != "" {
		out.installPost(tool.postNote)
	}
}

func selectTools(flags toolsInstallFlags, args []string, tools []contextTool) ([]contextTool, error) {
	if len(args) == 1 {
		return selectToolByName(args[0], tools)
	}

	tierStr := flags.tier
	if flags.all {
		tierStr = "all"
	}
	tier, err := parseTier(tierStr)
	if err != nil {
		return nil, err
	}

	var selected []contextTool
	for _, tool := range tools {
		if tier == "" || tool.tier == tier {
			selected = append(selected, tool)
		}
	}
	return selected, nil
}

func selectToolByName(name string, tools []contextTool) ([]contextTool, error) {
	for _, tool := range tools {
		if tool.name == name {
			return []contextTool{tool}, nil
		}
	}
	return nil, fmt.Errorf("unknown tool %q — run `loom tools status` to see available tools", name)
}

// runInstallCmd executes a single shell install command.
func runInstallCmd(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
