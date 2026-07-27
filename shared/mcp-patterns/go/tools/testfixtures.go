//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for reusable test fixtures and helpers.
// Renamed from testfixtures_test.go to testfixtures.go so it can be imported by downstream tests.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

package tools

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	// Replace with your module's logging package path:
	// "<YOUR_MODULE>/internal/logging"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// buildRequest turns an argument map into an mcp.CallToolRequest.
func BuildRequest(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

// extractText pulls the text off the first Content entry of a tool result.
func ExtractText(t *testing.T, r *mcp.CallToolResult) string {
	t.Helper()
	if r == nil {
		t.Fatal("nil result")
	}
	if len(r.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := r.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", r.Content[0])
	}
	return tc.Text
}

// silentLogger returns a logger that discards its output.
func SilentLogger() *logging.Logger {
	return logging.NewLogger(&bytes.Buffer{})
}

// writeTSFile writes a .ts file at a project-relative path.
func WriteTSFile(t *testing.T, root, relPath, body string) {
	t.Helper()
	full := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", relPath, err)
	}
}

// writeFile writes body to an absolute path, creating parent directories.
func WriteFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
