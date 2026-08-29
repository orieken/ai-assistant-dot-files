package tools

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/loom/shared/mcp/internal/domain"
	"github.com/orieken/loom/shared/mcp/internal/logging"
)

// BuildRequest turns an argument map into a domain.ToolRequest.
func BuildRequest(args map[string]any) domain.ToolRequest {
	return domain.ToolRequest{Args: args}
}

// ExtractText pulls the text off the first Content entry of a tool result.
func ExtractText(t *testing.T, r *domain.ToolResult) string {
	t.Helper()
	if r == nil {
		t.Fatal("nil result")
	}
	if len(r.Content) == 0 {
		t.Fatal("result has no content")
	}
	return r.Content[0].Text
}

// SilentLogger returns a logger that discards its output.
func SilentLogger() *logging.Logger {
	return logging.NewLogger(&bytes.Buffer{})
}

// WriteTSFile writes a .ts file at a project-relative path.
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

// WriteFile writes body to an absolute path, creating parent directories.
func WriteFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
