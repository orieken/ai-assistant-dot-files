package server

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
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
		s.AddTool(mcpTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tool.Execute(ctx, req)
		})
	}

	h.logger.Info("Tools registered successfully")
	return nil
}
