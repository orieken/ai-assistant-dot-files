package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func noopRunCmd(_ string) error          { return nil }
func failRunCmd(_ string) error          { return errors.New("install failed") }

func TestToolsInstallAllHighByDefault(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsInstall(
		toolsInstallFlags{tier: "high"},
		nil,
		testTools,
		lookPathMissing,
		noopRunCmd,
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// tokei should be installed (brew)
	if !strings.Contains(out, "✓ tokei") {
		t.Errorf("expected tokei installed, got:\n%s", out)
	}
	// git-churn is built-in
	if !strings.Contains(out, "git-churn") {
		t.Errorf("expected git-churn in output, got:\n%s", out)
	}
	// ctx should print manual note, not be installed
	if !strings.Contains(out, "curl -fsSL https://ctx.rs/install | sh") {
		t.Errorf("expected ctx manual note, got:\n%s", out)
	}
	// repomix is medium, should not appear
	if strings.Contains(out, "repomix") {
		t.Errorf("expected repomix (medium) excluded from high-tier install, got:\n%s", out)
	}
}

func TestToolsInstallSkipsAlreadyInstalled(t *testing.T) {
	var buf bytes.Buffer
	_ = executeToolsInstall(
		toolsInstallFlags{tier: "high"},
		nil,
		testTools,
		lookPathOnlyTokei,
		noopRunCmd,
		&buf,
	)
	out := buf.String()
	if !strings.Contains(out, "already installed") {
		t.Errorf("expected tokei to show as already installed, got:\n%s", out)
	}
}

func TestToolsInstallByName(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsInstall(
		toolsInstallFlags{tier: "high"},
		[]string{"repomix"},
		testTools,
		lookPathMissing,
		noopRunCmd,
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "✓ repomix") {
		t.Errorf("expected repomix installed by name, got:\n%s", out)
	}
	if strings.Contains(out, "tokei") {
		t.Errorf("expected only repomix in output, got:\n%s", out)
	}
}

func TestToolsInstallUnknownNameReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsInstall(
		toolsInstallFlags{tier: "high"},
		[]string{"nonexistent"},
		testTools,
		lookPathMissing,
		noopRunCmd,
		&buf,
	)
	if err == nil {
		t.Fatal("expected error for unknown tool name")
	}
	if !strings.Contains(err.Error(), "unknown tool") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToolsInstallReturnsErrorOnFailure(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsInstall(
		toolsInstallFlags{tier: "high"},
		[]string{"tokei"},
		testTools,
		lookPathMissing,
		failRunCmd,
		&buf,
	)
	if err == nil {
		t.Fatal("expected error when install command fails")
	}
	if !strings.Contains(err.Error(), "install(s) failed") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToolsInstallAllTiers(t *testing.T) {
	var buf bytes.Buffer
	err := executeToolsInstall(
		toolsInstallFlags{all: true},
		nil,
		testTools,
		lookPathMissing,
		noopRunCmd,
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "repomix") {
		t.Errorf("expected repomix in all-tiers install, got:\n%s", out)
	}
}
