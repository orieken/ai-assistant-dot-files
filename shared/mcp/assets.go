// Package mcpassets exposes the MCP reference source for embedding in loom.
package mcpassets

import "embed"

// SourceFS holds the complete MCP scaffold without crossing a Go module boundary.
//
//go:embed README.md go.mod go.sum all:cmd all:internal all:register
var SourceFS embed.FS
