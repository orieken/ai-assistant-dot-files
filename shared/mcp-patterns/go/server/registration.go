//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for loop-based MCP server registration pattern.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

package server

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	// Replace with your module's domain and logging packages:
	// "<YOUR_MODULE>/internal/domain"
	// "<YOUR_MODULE>/internal/logging"
	"github.com/orieken/saturday-mcp/internal/domain"
)

// Tools returns every domain.Tool the handler exposes to MCP clients.
func (h *Handler) Tools() []domain.Tool {
	return h.tools
}

// RegisterTools registers every tool returned by h.Tools() with the MCP server.
// Each tool's declarative OutputSchema is marshaled onto RawOutputSchema.
func (h *Handler) RegisterTools(s *server.MCPServer) error {
	h.logger.Info("Registering tools")

	for _, t := range h.Tools() {
		tool := t
		mcpTool := mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.InputSchema(),
		}
		if outSchema := tool.OutputSchema(); outSchema != nil {
			if raw, err := json.Marshal(outSchema); err == nil {
				mcpTool.RawOutputSchema = raw
			} else {
				h.logger.Warn("Failed to marshal output schema", "tool", tool.Name(), "error", err)
			}
		}
		s.AddTool(mcpTool, tool.Execute)
	}

	h.logger.Info("Tools registered successfully")
	return nil
}
