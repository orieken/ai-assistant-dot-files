package analyzers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testContract = `# Contract: analysis.md

## Required Sections (exact heading text and level)
- ` + "`## Summary`" + `
- ` + "`### Acceptance Criteria`" + `
- ` + "`## Definition of Done`" + `

## Validation Rule
Presence of every heading above.
`

const compliantArtifact = `---
feature: "user-auth"
bounded_context: "identity"
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---
## Summary
Real content.

### Acceptance Criteria
- AC1

## Definition of Done
Done.
`

func TestValidateArtifactAgainstContract(t *testing.T) {
	cases := []struct {
		name            string
		artifact        string
		wantStatus      string
		wantMissing     []string
		wantWarnSubstrs []string
	}{
		{
			name:       "compliant artifact passes with no warnings",
			artifact:   compliantArtifact,
			wantStatus: ArtifactStatusPass,
		},
		{
			name:        "missing section fails",
			artifact:    "## Summary\n\n## Definition of Done\n",
			wantStatus:  ArtifactStatusFail,
			wantMissing: []string{"### Acceptance Criteria"},
			wantWarnSubstrs: []string{
				"missing retrieval frontmatter block",
			},
		},
		{
			name:        "wrong heading level fails",
			artifact:    "### Summary\n\n### Acceptance Criteria\n\n## Definition of Done\n",
			wantStatus:  ArtifactStatusFail,
			wantMissing: []string{"## Summary"},
			wantWarnSubstrs: []string{
				"missing retrieval frontmatter block",
			},
		},
		{
			name:        "heading inside a code fence does not count",
			artifact:    "## Summary\n```\n### Acceptance Criteria\n```\n## Definition of Done\n",
			wantStatus:  ArtifactStatusFail,
			wantMissing: []string{"### Acceptance Criteria"},
			wantWarnSubstrs: []string{
				"missing retrieval frontmatter block",
			},
		},
		{
			name: "incomplete frontmatter warns per missing field but still passes",
			artifact: "---\nfeature: \"user-auth\"\n---\n" +
				"## Summary\n\n### Acceptance Criteria\n\n## Definition of Done\n",
			wantStatus: ArtifactStatusPass,
			wantWarnSubstrs: []string{
				"`bounded_context` is missing",
				"`linked_kis` is missing",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateFixture(t, testContract, tc.artifact)
			assertValidation(t, result, tc.wantStatus, tc.wantMissing, tc.wantWarnSubstrs)
		})
	}
}

func TestValidateRejectsContractWithoutRequiredSections(t *testing.T) {
	_, err := validateFixtureErr(t, "# Contract: agent frontmatter\n\n## Required Fields\n| name |\n", "## Summary\n")
	if err == nil || !strings.Contains(err.Error(), "Required Sections") {
		t.Fatalf("want no-Required-Sections error, got %v", err)
	}
}

func TestValidateReportsUnreadableFiles(t *testing.T) {
	analyzer := NewArtifactContractAnalyzer()
	if _, err := analyzer.Validate("/nonexistent/a.md", "/nonexistent/c.md"); err == nil {
		t.Fatal("want error for unreadable contract, got nil")
	}
	contractPath := writeTempFile(t, "contract.md", testContract)
	if _, err := analyzer.Validate("/nonexistent/a.md", contractPath); err == nil {
		t.Fatal("want error for unreadable artifact, got nil")
	}
}

func validateFixture(t *testing.T, contract, artifact string) *ArtifactValidationResult {
	t.Helper()
	result, err := validateFixtureErr(t, contract, artifact)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	return result
}

func validateFixtureErr(t *testing.T, contract, artifact string) (*ArtifactValidationResult, error) {
	t.Helper()
	contractPath := writeTempFile(t, "contract.md", contract)
	artifactPath := writeTempFile(t, "artifact.md", artifact)
	return NewArtifactContractAnalyzer().Validate(artifactPath, contractPath)
}

func writeTempFile(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func assertValidation(t *testing.T, result *ArtifactValidationResult, wantStatus string, wantMissing, wantWarnSubstrs []string) {
	t.Helper()
	if result.Status != wantStatus {
		t.Errorf("status = %q, want %q (violations: %v)", result.Status, wantStatus, result.Violations)
	}
	assertMissingSections(t, result.Violations, wantMissing)
	assertWarningSubstrings(t, result.Warnings, wantWarnSubstrs)
}

func assertMissingSections(t *testing.T, violations []ArtifactViolation, wantMissing []string) {
	t.Helper()
	if len(violations) != len(wantMissing) {
		t.Fatalf("violations = %v, want %d missing sections %v", violations, len(wantMissing), wantMissing)
	}
	for i, want := range wantMissing {
		if violations[i].Type != ViolationMissingSection || violations[i].Detail != want {
			t.Errorf("violation[%d] = %+v, want missing_section %q", i, violations[i], want)
		}
	}
}

func assertWarningSubstrings(t *testing.T, warnings []string, wantSubstrs []string) {
	t.Helper()
	joined := strings.Join(warnings, "\n")
	for _, substr := range wantSubstrs {
		if !strings.Contains(joined, substr) {
			t.Errorf("warnings %v missing substring %q", warnings, substr)
		}
	}
	if len(wantSubstrs) == 0 && len(warnings) != 0 {
		t.Errorf("want no warnings, got %v", warnings)
	}
}
