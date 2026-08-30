package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/orieken/loom/cmd/loom/internal/manifest"
)

func TestExecuteHealthPassesForIntactInstall(t *testing.T) {
	target := t.TempDir()
	writeHealthFixture(t, target, []string{"alpha.md", "beta.md"})
	var output bytes.Buffer
	if err := executeHealth(healthRequest{target: target, display: target}, healthContent(), &output); err != nil {
		t.Fatalf("execute health: %v", err)
	}
	for _, want := range []string{"manifest found (v3.3.14)", "all 1 installed paths intact", "symlinks unbroken", "All checks passed."} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output %q does not contain %q", output.String(), want)
		}
	}
}

func TestExecuteHealthFailsForMissingInstalledPath(t *testing.T) {
	target := t.TempDir()
	records := []manifest.PlatformRecord{{Name: "claude-code", Paths: []string{".claude/agents"}}}
	writeHealthManifest(t, target, records)
	var output bytes.Buffer
	err := executeHealth(healthRequest{target: target, display: target}, healthContent(), &output)
	if !errors.Is(err, errHealthFailed) {
		t.Fatalf("error = %v", err)
	}
	if !strings.Contains(output.String(), "missing installed path .claude/agents") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestExecuteHealthFailsForBrokenSymlink(t *testing.T) {
	target := t.TempDir()
	link := filepath.Join(target, "agents")
	if err := os.Symlink(filepath.Join(target, "missing"), link); err != nil {
		t.Fatalf("create broken symlink: %v", err)
	}
	writeHealthManifest(t, target, []manifest.PlatformRecord{{Name: "claude-code", Paths: []string{"agents"}}})
	var output bytes.Buffer
	err := executeHealth(healthRequest{target: target, display: target}, healthContent(), &output)
	if !errors.Is(err, errHealthFailed) || !strings.Contains(output.String(), "broken symlink agents") {
		t.Fatalf("error = %v, output = %s", err, output.String())
	}
}

func TestExecuteHealthWarnsForAgentCountMismatch(t *testing.T) {
	target := t.TempDir()
	writeHealthFixture(t, target, []string{"alpha.md"})
	var output bytes.Buffer
	if err := executeHealth(healthRequest{target: target, display: target}, healthContent(), &output); err != nil {
		t.Fatalf("warning should not fail health: %v", err)
	}
	if !strings.Contains(output.String(), "embedded has 2, found 1") || !strings.Contains(output.String(), "1 warning") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

// TestExecuteHealthReportsLevelOneWithGaps covers the D.4 done-when
// criterion: a fresh level 1 install reports Level 1 plus a concrete
// checklist of level 2 gaps.
func TestExecuteHealthReportsLevelOneWithGaps(t *testing.T) {
	target := t.TempDir()
	writeHealthFixture(t, target, []string{"alpha.md", "beta.md"})
	var output bytes.Buffer
	if err := executeHealth(healthRequest{target: target, display: target}, healthContent(), &output); err != nil {
		t.Fatalf("execute health: %v", err)
	}
	wants := []string{
		"Level 1 — Foundational prompts",
		"gaps to Level 2",
		"MCP server not configured",
		`bundle "workflows" not fully installed`,
	}
	for _, want := range wants {
		if !strings.Contains(output.String(), want) {
			t.Errorf("output missing %q:\n%s", want, output.String())
		}
	}
}

func TestExecuteHealthFailsForMissingManifest(t *testing.T) {
	target := t.TempDir()
	var output bytes.Buffer
	err := executeHealth(healthRequest{target: target, display: target}, healthContent(), &output)
	if !errors.Is(err, errHealthFailed) || !strings.Contains(output.String(), "manifest missing") {
		t.Fatalf("error = %v, output = %s", err, output.String())
	}
}

// healthLevelsFixture is a minimal but valid shared/levels.yaml so the
// maturity assessment can run against test targets.
const healthLevelsFixture = `version: 1
landed: [D.1]
levels:
  - level: 1
    name: Foundational prompts
    bundles:
      - id: agents
        paths: [shared/agents]
  - level: 2
    name: Coordinated multi-agent
    bundles:
      - id: mcp-registration
        requires: D.1
        action: mcp-config
      - id: workflows
        paths: [shared/workflows]
  - level: 3
    name: Observed and governed
    bundles:
      - id: telemetry-stream
        requires: L3.9
        action: telemetry
  - level: 4
    name: Self-improving
    bundles:
      - id: eval-loop
        requires: L4.1
        action: evaluation-loop
`

func healthContent() fstest.MapFS {
	return fstest.MapFS{
		"shared/VERSION":         {Data: []byte("v3.3.14\n")},
		"shared/agents/alpha.md": {Data: []byte("alpha")},
		"shared/agents/beta.md":  {Data: []byte("beta")},
		"shared/levels.yaml":     {Data: []byte(healthLevelsFixture)},
	}
}

func writeHealthFixture(t *testing.T, target string, agents []string) {
	t.Helper()
	directory := filepath.Join(target, ".claude", "agents")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatalf("create agents directory: %v", err)
	}
	for _, name := range agents {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(name), 0o644); err != nil {
			t.Fatalf("write agent %s: %v", name, err)
		}
	}
	records := []manifest.PlatformRecord{{Name: "claude-code", Paths: []string{".claude/agents"}}}
	writeHealthManifest(t, target, records)
}

func writeHealthManifest(t *testing.T, target string, records []manifest.PlatformRecord) {
	t.Helper()
	installed := manifest.Manifest{Version: "v3.3.14", Platforms: records}
	if err := manifest.Write(target, installed); err != nil {
		t.Fatalf("write health manifest: %v", err)
	}
}
