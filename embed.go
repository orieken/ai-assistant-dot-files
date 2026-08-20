// Package loom exposes the framework content embedded in the loom binary.
package loom

import (
	"embed"

	mcpassets "github.com/orieken/ai-assistant-dotfiles/mcp"
)

// FrameworkFS holds all shared framework content baked in at compile time.
// It lives at the module root because go:embed cannot traverse parent directories.
//
//go:embed all:shared/agents all:shared/skills all:shared/rules all:shared/configs all:shared/contracts all:shared/schemas shared/ARCHITECTURE_RULES.md shared/DOMAIN_DICTIONARY.md all:templates/claude-feature-team
var FrameworkFS embed.FS

// MCPFS holds the MCP reference module embedded by its own nested Go module.
var MCPFS = mcpassets.SourceFS
