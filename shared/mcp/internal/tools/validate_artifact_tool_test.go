package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/loom/shared/mcp/internal/analyzers"
	"github.com/orieken/loom/shared/mcp/internal/domain"
)

const toolTestContract = "# Contract: analysis.md\n\n" +
	"## Required Sections (exact heading text and level)\n" +
	"- `## Summary`\n" +
	"- `## Definition of Done`\n"

func TestValidateArtifactToolInfersContractFromFilename(t *testing.T) {
	contractsDir := t.TempDir()
	WriteFile(t, filepath.Join(contractsDir, "analysis-contract.md"), toolTestContract)
	artifactPath := filepath.Join(t.TempDir(), "analysis.md")
	WriteFile(t, artifactPath, "## Summary\n\n## Definition of Done\n")

	result := executeValidateArtifact(t, contractsDir, map[string]any{"artifactPath": artifactPath})

	var parsed analyzers.ArtifactValidationResult
	if err := json.Unmarshal([]byte(ExtractText(t, result)), &parsed); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if parsed.Status != analyzers.ArtifactStatusPass {
		t.Errorf("status = %q, want PASS (violations %v)", parsed.Status, parsed.Violations)
	}
}

func TestValidateArtifactToolErrors(t *testing.T) {
	artifactPath := filepath.Join(t.TempDir(), "analysis.md")
	WriteFile(t, artifactPath, "## Summary\n")
	cases := []struct {
		name         string
		contractsDir string
		args         map[string]any
		wantSubstr   string
	}{
		{
			name:       "missing artifactPath",
			args:       map[string]any{},
			wantSubstr: "artifactPath is required",
		},
		{
			name:       "unknown artifact with no explicit contract",
			args:       map[string]any{"artifactPath": "/tmp/mystery.md"},
			wantSubstr: "no known contract",
		},
		{
			name:       "known artifact but contracts dir unknown",
			args:       map[string]any{"artifactPath": artifactPath},
			wantSubstr: "AI_ASSISTANT_DOTFILES_PATH",
		},
		{
			name:       "explicit contract path that does not exist",
			args:       map[string]any{"artifactPath": artifactPath, "contractPath": "/nonexistent/contract.md"},
			wantSubstr: "Artifact validation failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := executeValidateArtifact(t, tc.contractsDir, tc.args)
			assertErrorResult(t, result, tc.wantSubstr)
		})
	}
}

func executeValidateArtifact(t *testing.T, contractsDir string, args map[string]any) *domain.ToolResult {
	t.Helper()
	tool := NewValidateArtifactTool(SilentLogger(), analyzers.NewArtifactContractAnalyzer(), contractsDir)
	result, err := tool.Execute(context.Background(), BuildRequest(args))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	return result
}

func assertErrorResult(t *testing.T, result *domain.ToolResult, wantSubstr string) {
	t.Helper()
	if !result.IsError {
		t.Fatalf("want error result, got success: %s", ExtractText(t, result))
	}
	if text := ExtractText(t, result); !strings.Contains(text, wantSubstr) {
		t.Errorf("error text %q missing %q", text, wantSubstr)
	}
}
