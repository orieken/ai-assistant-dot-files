package tools_test

import (
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// TestPublicAPIImportsOnlyStdlib is the D.2 fitness function (and carries
// forward M0.3's): the public embedding package must never import internal
// packages, mcp-go, jsonschema, or any other non-stdlib package — consumers
// must not inherit loom's dependency pins, and exported signatures cannot
// reference what is never imported. Standard-library import paths never
// contain a dot in their first segment; anything else fails.
func TestPublicAPIImportsOnlyStdlib(t *testing.T) {
	for _, path := range packageImports(t) {
		if isStdlibPath(path) {
			continue
		}
		t.Errorf("public tools package imports non-stdlib package %q — keep the embedding API dependency-free", path)
	}
}

func packageImports(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	var imports []string
	for _, entry := range entries {
		if !isNonTestGoFile(entry.Name()) {
			continue
		}
		imports = append(imports, fileImports(t, entry.Name())...)
	}
	return imports
}

func isNonTestGoFile(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

func fileImports(t *testing.T, name string) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		imports = append(imports, strings.Trim(spec.Path.Value, `"`))
	}
	return imports
}

func isStdlibPath(path string) bool {
	firstSegment, _, _ := strings.Cut(path, "/")
	return !strings.Contains(firstSegment, ".")
}
