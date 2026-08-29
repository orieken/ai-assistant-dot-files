package cmd

import "github.com/spf13/cobra"

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Model Context Protocol server commands",
	Long:  "Serve the framework's MCP tools from the loom binary.",
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
