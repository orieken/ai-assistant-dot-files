package server

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/orieken/loom/shared/mcp/internal/domain"
)

// Tools returns every domain.Tool the handler exposes to MCP clients.
func (h *Handler) Tools() []domain.Tool {
	return h.tools
}

// RegisterTools registers every tool returned by h.Tools() with the MCP server.
func (h *Handler) RegisterTools(s *server.MCPServer) error {
	h.logger.Info("Registering tools")

	for _, t := range h.Tools() {
		s.AddTool(mcpToolDefinition(t), mcpToolHandler(t))
	}

	h.logger.Info("Tools registered successfully")
	return nil
}
