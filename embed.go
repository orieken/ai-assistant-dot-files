// Package loom exposes the framework content embedded in the loom binary.
package loom

import "embed"

// FrameworkFS holds all shared framework content baked in at compile time.
// It lives at the module root because go:embed cannot traverse parent directories.
//
//go:embed all:shared/agents all:shared/skills all:shared/rules all:shared/configs all:shared/contracts all:shared/schemas
var FrameworkFS embed.FS
