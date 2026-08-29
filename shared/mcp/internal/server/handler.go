package server

import (
	"github.com/orieken/loom/shared/mcp/internal/domain"
	"github.com/orieken/loom/shared/mcp/internal/logging"
)

// Handler owns the tool registry and orchestrates MCP registration.
// Constructed via New; tools are wired in tool_provider.go.
type Handler struct {
	logger   *logging.Logger
	registry *domain.Registry
}

// New constructs a Handler wired with all framework M1 tools.
func New(logger *logging.Logger) *Handler {
	return &Handler{
		logger:   logger,
		registry: buildFrameworkRegistry(logger),
	}
}

// Registry returns the handler's tool registry.
func (h *Handler) Registry() *domain.Registry {
	return h.registry
}
