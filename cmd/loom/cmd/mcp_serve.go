package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/loom/shared/mcp/register"
	"github.com/spf13/cobra"
)

type mcpServeFlags struct {
	logFile string
}

var mcpServeArgs mcpServeFlags

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the framework MCP tools over stdio",
	Long: `Serve the six framework MCP tools over stdio transport.

Structured JSON logs go to stderr (or --log-file) — never stdout,
which carries the MCP wire protocol. The server blocks until the
client closes stdin or the process receives SIGINT.`,
	Args: cobra.NoArgs,
	RunE: runMCPServe,
}

func init() {
	mcpCmd.AddCommand(mcpServeCmd)
	mcpServeCmd.Flags().StringVar(&mcpServeArgs.logFile, "log-file", "", "append structured JSON logs to this file instead of stderr")
}

func runMCPServe(_ *cobra.Command, _ []string) error {
	logWriter, closeLog, err := openMCPLogWriter(mcpServeArgs.logFile)
	if err != nil {
		return err
	}
	defer closeLog()
	return serveMCP(logWriter)
}

func serveMCP(logWriter io.Writer) error {
	mcpServer := server.NewMCPServer("loom", version, server.WithToolCapabilities(true))
	// `loom mcp serve` IS the deprecation notice's recommended replacement, so
	// it legitimately keeps using the compat wrapper until D.2's embedding API
	// grows an in-binary adapter.
	if err := register.FrameworkTools(mcpServer, logWriter); err != nil { //nolint:staticcheck
		return fmt.Errorf("register framework tools: %w", err)
	}
	if err := server.ServeStdio(mcpServer); err != nil {
		return fmt.Errorf("mcp serve: %w", err)
	}
	return nil
}

func openMCPLogWriter(path string) (io.Writer, func(), error) {
	if path == "" {
		return os.Stderr, func() {}, nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}
	return file, func() { _ = file.Close() }, nil
}
