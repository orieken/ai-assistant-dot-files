package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var testTools = []contextTool{
	{
		name:        "tokei",
		tier:        tierHigh,
		description: "accurate per-language line counts",
		binary:      "tokei",
		installCmds: []string{"brew install tokei"},
		kiName:      "tokei-token-budget",
	},
	{
		name:        "git-churn",
		tier:        tierHigh,
		description: "built-in git risk signal",
		binary:      "",
		kiName:      "git-churn-risk-signal",
	},
	{
		name:        "ctx",
		tier:        tierHigh,
		description: "session history search",
		binary:      "ctx",
		manualNote:  "curl -fsSL https://ctx.rs/install | sh",
		postNote:    "after install, run: ctx setup",
		kiName:      "ctx-session-history-search",
	},
	{
		name:        "repomix",
		tier:        tierMedium,
		description: "codebase packing",
		binary:      "repomix",
		installCmds: []string{"npm install -g repomix"},
		kiName:      "repomix-codebase-packing",
	},
}

func lookPathFound(_ string) (string, error)            { return "/usr/bin/found", nil }
func lookPathMissing(_ string) (string, error)          { return "", errors.New("not found") }
func lookPathOnlyTokei(name string) (string, error) {
	if name == "tokei" {
		return "/usr/bin/tokei", nil
	}
	return "", errors.New("not found")
}

func TestToolsStatusAllInstalled(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsStatus(toolsStatusFlags{tier: "all"}, testTools, lookPathFound, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "All context tools are installed.") {
		t.Errorf("expected all-installed summary, got:\n%s", out)
	}
	if strings.Contains(out, "✗") {
		t.Errorf("expected no failures, got:\n%s", out)
	}
}

func TestToolsStatusNoneInstalled(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsStatus(toolsStatusFlags{tier: "all"}, testTools, lookPathMissing, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// git-churn has no binary so is always ✓; tokei, ctx, repomix are missing = 3
	if !strings.Contains(out, "3 tools not installed") {
		t.Errorf("expected 3-missing summary, got:\n%s", out)
	}
	if !strings.Contains(out, "brew install tokei") {
		t.Errorf("expected tokei install hint, got:\n%s", out)
	}
	if !strings.Contains(out, "curl -fsSL https://ctx.rs/install | sh") {
		t.Errorf("expected ctx manual note, got:\n%s", out)
	}
}

func TestToolsStatusBuiltInAlwaysInstalled(t *testing.T) {
	var buf bytes.Buffer
	_ = executeToolsStatus(toolsStatusFlags{tier: "all"}, testTools, lookPathMissing, &buf)
	out := buf.String()
	// git-churn line should show ✓ even when lookPath returns not-found for everything
	if !strings.Contains(out, "✓ git-churn") {
		t.Errorf("expected git-churn to show as installed, got:\n%s", out)
	}
}

func TestToolsStatusFilterHighTier(t *testing.T) {
	var buf bytes.Buffer
	_ = executeToolsStatus(toolsStatusFlags{tier: "high"}, testTools, lookPathFound, &buf)
	out := buf.String()
	if strings.Contains(out, "repomix") {
		t.Errorf("expected repomix (medium) to be filtered out, got:\n%s", out)
	}
}

func TestToolsStatusFilterMediumTier(t *testing.T) {
	var buf bytes.Buffer
	_ = executeToolsStatus(toolsStatusFlags{tier: "medium"}, testTools, lookPathFound, &buf)
	out := buf.String()
	if strings.Contains(out, "tokei") {
		t.Errorf("expected tokei (high) to be filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "repomix") {
		t.Errorf("expected repomix (medium) in output, got:\n%s", out)
	}
}

func TestToolsStatusInvalidTier(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsStatus(toolsStatusFlags{tier: "extreme"}, testTools, lookPathFound, &buf)
	if err == nil {
		t.Fatal("expected error for unknown tier")
	}
	if !strings.Contains(err.Error(), "unknown tier") {
		t.Errorf("unexpected error: %v", err)
	}
}
