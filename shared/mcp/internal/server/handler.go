package server

import (
	"github.com/orieken/loom/shared/mcp/internal/domain"
	"github.com/orieken/loom/shared/mcp/internal/logging"
)

// Handler owns the tool registry and orchestrates MCP registration.
// Constructed via New; tools are wired in tool_provider.go.
type Handler struct {
	logger *logging.Logger
	tools  []domain.Tool
}

// New constructs a Handler wired with all framework M1 tools.
func New(logger *logging.Logger) *Handler {
	return &Handler{
		logger: logger,
		tools:  buildFrameworkTools(logger),
	}
}
