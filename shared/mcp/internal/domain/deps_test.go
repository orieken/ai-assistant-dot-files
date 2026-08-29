package domain_test

import (
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// TestDomainImportsOnlyStdlib is the M0.3 fitness function: the domain package
// must never import third-party or framework packages (architecture-guardrails.md
// #1). Standard-library import paths never contain a dot in their first segment;
// anything else fails.
func TestDomainImportsOnlyStdlib(t *testing.T) {
	for _, path := range domainImports(t) {
		if isStdlibPath(path) {
			continue
		}
		t.Errorf("internal/domain imports non-stdlib package %q — move transport types into the server adapter layer", path)
	}
}

func domainImports(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read domain package dir: %v", err)
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
