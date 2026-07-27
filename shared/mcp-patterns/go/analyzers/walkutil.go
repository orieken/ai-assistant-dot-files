//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic walkutil helpers.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

// Package analyzers - walkutil.go
//
// Shared filepath.Walk helpers used by every analyzer that scans a
// project tree.
package analyzers

import (
	"path/filepath"
	"strings"
)

// SkippedDirNames lists directory basenames the walk should never descend into.
var SkippedDirNames = map[string]struct{}{
	"node_modules": {},
	".git":         {},
	"dist":         {},
	"build":        {},
	"vendor":       {},
	".venv":        {},
}

// SkipUninterestingDir returns filepath.SkipDir for directory names the
// walk should skip: any name in SkippedDirNames, plus any hidden
// dot-directory below the root.
func SkipUninterestingDir(root, current, name string) error {
	if _, skip := SkippedDirNames[name]; skip {
		return filepath.SkipDir
	}
	if strings.HasPrefix(name, ".") && current != root {
		return filepath.SkipDir
	}
	return nil
}
