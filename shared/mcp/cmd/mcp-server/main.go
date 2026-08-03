package main

import (
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/logging"
	mcpserver "github.com/orieken/ai-assistant-dotfiles/mcp/internal/server"
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
