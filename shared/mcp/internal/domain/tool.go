package domain

import (
	"context"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool is the framework's first-class abstraction for every capability exposed over MCP.
// Every M1 framework tool implements this interface so it can be registered uniformly
// by the server.Handler without knowing any tool's concrete type.
type Tool interface {
	Name() string
	Description() string
	InputSchema() mcp.ToolInputSchema
	OutputSchema() *jsonschema.Schema
	Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
