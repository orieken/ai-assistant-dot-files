// Package register exposes the framework's MCP tools for embedding in external
// MCP servers. Callers import this package and call FrameworkTools —
// the internal implementation packages stay encapsulated.
package register

import (
	"io"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/loom/shared/mcp/internal/logging"
	mcpserver "github.com/orieken/loom/shared/mcp/internal/server"
)

// FrameworkTools registers all framework tools on s:
//
//   - search_ki           — BM25 full-text search over Knowledge Items
//   - check_ubiquitous_language — validates domain term usage
//   - analyze_complexity  — cyclomatic complexity analysis
//   - check_accessibility — accessibility rule checks
//   - verify_dependencies — dependency boundary validation
//   - search_docs         — documentation search
//
// logWriter receives structured JSON log output; pass nil to use os.Stderr.
//
// Internally the tools live in a name-keyed registry with per-tool metadata
// (timeout budget, retry class, permission scope). A transport-free public
// embedding API over that registry ships with roadmap D.2 — until then this
// function remains the supported integration path.
func FrameworkTools(s *server.MCPServer, logWriter io.Writer) error {
	if logWriter == nil {
		logWriter = os.Stderr
	}
	handler := mcpserver.New(logging.NewLogger(logWriter))
	return handler.RegisterTools(s)
}
