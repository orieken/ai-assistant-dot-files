// Deprecated: this standalone entrypoint is retained for one release cycle
// only. Use `loom mcp serve` instead — the brew-installed loom binary serves
// the same six framework tools over stdio. See shared/mcp/README.md.
package main

import (
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/orieken/loom/shared/mcp/internal/logging"
	mcpserver "github.com/orieken/loom/shared/mcp/internal/server"
)

func main() {
	logger := logging.NewLogger(os.Stderr)

	s := server.NewMCPServer(
		"ai-assistant-dotfiles-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	handler := mcpserver.New(logger)
	if err := handler.RegisterTools(s); err != nil {
		logger.Error("Failed to register tools", "error", err)
		os.Exit(1)
	}

	if err := server.ServeStdio(s); err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}
