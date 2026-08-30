package cmd

import (
	"strings"
	"testing"

	"github.com/orieken/loom/cmd/loom/internal/levels"
)

type fakeLevelProbe struct {
	installedBundles map[string]bool
	mcpAnswering     bool
	telemetryPresent bool
	policiesFound    bool
}

func (probe fakeLevelProbe) pathBundleInstalled(bundle levels.Bundle) (bool, string) {
	if probe.installedBundles[bundle.ID] {
		return true, "bundle \"" + bundle.ID + "\" installed"
	}
	return false, "bundle \"" + bundle.ID + "\" not fully installed"
}

func (probe fakeLevelProbe) mcpServerAnswering() (bool, string) {
	if probe.mcpAnswering {
		return true, "MCP server answering tools/list (6 tools)"
	}
	return false, "MCP server configured but not answering tools/list"
}

func (probe fakeLevelProbe) telemetryStreamPresent() (bool, string) {
	if probe.telemetryPresent {
		return true, "telemetry stream present"
	}
	return false, "telemetry stream absent"
}

func (probe fakeLevelProbe) policiesPresent() (bool, string) {
	if probe.policiesFound {
		return true, "policy files present"
	}
	return false, "no policy files"
}

// maturityTestProfile mirrors the real shared/levels.yaml shape: level 3
// carries a docs-only bundle that must never confer the level, and levels
// 2-4 carry roadmap-gated bundles.
func maturityTestProfile(landed ...string) levels.Profile {
	return levels.Profile{
		Version: 1,
		Landed:  landed,
		Levels: []levels.Level{
			{Level: 1, Name: "Foundational prompts", Bundles: []levels.Bundle{
				{ID: "core-rules", Paths: []string{"shared/rules/approval-gates.md"}},
				{ID: "agents", Paths: []string{"shared/agents"}},
			}},
			{Level: 2, Name: "Coordinated multi-agent", Bundles: []levels.Bundle{
				{ID: "mcp-registration", Requires: "D.1", Action: "mcp-config"},
				{ID: "workflows", Paths: []string{"shared/workflows"}},
				{ID: "executor", Requires: "M0.4", Action: "executor"},
			}},
			{Level: 3, Name: "Observed and governed", Bundles: []levels.Bundle{
				{ID: "telemetry-docs", DocsOnly: true, Paths: []string{"shared/telemetry"}},
				{ID: "telemetry-stream", Requires: "L3.9", Action: "telemetry"},
				{ID: "policy-engine", Requires: "L2.16", Action: "policies"},
			}},
			{Level: 4, Name: "Self-improving", Bundles: []levels.Bundle{
				{ID: "evaluation", DocsOnly: true, Paths: []string{"shared/evaluation"}},
				{ID: "eval-workspace", Requires: "L4.1", Paths: []string{"shared/eval-workspace"}},
			}},
		},
	}
}

func allBundlesInstalled() map[string]bool {
	return map[string]bool{
		"core-rules": true, "agents": true, "workflows": true,
		"telemetry-docs": true, "evaluation": true, "eval-workspace": true,
	}
}

func TestInferMaturityAcrossAllLevels(t *testing.T) {
	cases := []struct {
		name            string
		landed          []string
		probe           fakeLevelProbe
		wantLevel       int
		wantNextLevel   int
		wantUnreachable bool
		wantGaps        []string
	}{
		{
			name:          "nothing installed reports below level 1",
			landed:        []string{"D.1"},
			probe:         fakeLevelProbe{installedBundles: map[string]bool{}},
			wantLevel:     0,
			wantNextLevel: 1,
			wantGaps:      []string{`bundle "core-rules" not fully installed`, `bundle "agents" not fully installed`},
		},
		{
			name:          "level 2 docs present but server dead is still level 1",
			landed:        []string{"D.1"},
			probe:         fakeLevelProbe{installedBundles: allBundlesInstalled()},
			wantLevel:     1,
			wantNextLevel: 2,
			wantGaps: []string{
				"MCP server configured but not answering tools/list",
				"executor: gated on roadmap item M0.4 (not landed)",
			},
		},
		{
			name:            "level 2 attained and fully gated level 3 is unreachable",
			landed:          []string{"D.1"},
			probe:           fakeLevelProbe{installedBundles: allBundlesInstalled(), mcpAnswering: true},
			wantLevel:       2,
			wantNextLevel:   3,
			wantUnreachable: true,
			wantGaps: []string{
				"telemetry-stream: gated on roadmap item L3.9 (not landed)",
				"policy-engine: gated on roadmap item L2.16 (not landed)",
			},
		},
		{
			name:   "docs-only telemetry bundle never confers level 3",
			landed: []string{"D.1", "M0.4", "L3.9", "L2.16"},
			probe: fakeLevelProbe{
				installedBundles: allBundlesInstalled(), mcpAnswering: true, policiesFound: true,
			},
			wantLevel:     1,
			wantNextLevel: 2,
			wantGaps:      []string{`executor: no health check exists for action "executor" yet`},
		},
		{
			name:   "level 3 attained when telemetry and policies answer",
			landed: []string{"D.1", "L3.9", "L2.16"},
			probe: fakeLevelProbe{
				installedBundles: allBundlesInstalled(), mcpAnswering: true,
				telemetryPresent: true, policiesFound: true,
			},
			wantLevel:       3,
			wantNextLevel:   4,
			wantUnreachable: true,
			wantGaps:        []string{"eval-workspace: gated on roadmap item L4.1 (not landed)"},
		},
		{
			name:   "level 4 attained leaves no gap checklist",
			landed: []string{"D.1", "L3.9", "L2.16", "L4.1"},
			probe: fakeLevelProbe{
				installedBundles: allBundlesInstalled(), mcpAnswering: true,
				telemetryPresent: true, policiesFound: true,
			},
			wantLevel:     4,
			wantNextLevel: 0,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			profile := maturityTestProfile(testCase.landed...)
			report := inferMaturity(assessLevels(profile, testCase.probe))
			assertMaturityReport(t, report, testCase.wantLevel, testCase.wantNextLevel, testCase.wantUnreachable, testCase.wantGaps)
		})
	}
}

func assertMaturityReport(t *testing.T, report maturityReport, wantLevel, wantNextLevel int, wantUnreachable bool, wantGaps []string) {
	t.Helper()
	if report.level != wantLevel {
		t.Fatalf("level = %d, want %d (gaps: %v)", report.level, wantLevel, report.gaps)
	}
	if report.nextLevel != wantNextLevel {
		t.Fatalf("nextLevel = %d, want %d", report.nextLevel, wantNextLevel)
	}
	if report.nextIsUnreachable != wantUnreachable {
		t.Fatalf("nextIsUnreachable = %v, want %v", report.nextIsUnreachable, wantUnreachable)
	}
	assertGapsContain(t, report.gaps, wantGaps)
}

func assertGapsContain(t *testing.T, gaps, wanted []string) {
	t.Helper()
	joined := strings.Join(gaps, "\n")
	for _, want := range wanted {
		if !strings.Contains(joined, want) {
			t.Errorf("gaps %v missing %q", gaps, want)
		}
	}
}
