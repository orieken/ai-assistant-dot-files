// Package mcpassets exposes the MCP reference source for embedding in loom.
package mcpassets

import "embed"

// SourceFS holds the complete MCP reference source. Since the module merge
// into github.com/orieken/loom, this tree is reference material — it no longer
// carries its own go.mod, so copies of it are not standalone-buildable.
// Use `loom mcp serve` (or import the register package) instead.
//
//go:embed README.md all:cmd all:internal all:register
var SourceFS embed.FS
