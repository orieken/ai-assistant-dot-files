package analyzers

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Statuses reported by ArtifactValidationResult.Status.
const (
	// ArtifactStatusPass means every required section is present.
	ArtifactStatusPass = "PASS"
	// ArtifactStatusFail means at least one required section is missing.
	ArtifactStatusFail = "FAIL"
)

// ViolationMissingSection is the Type of a violation for an absent required heading.
const ViolationMissingSection = "missing_section"

// retrievalFrontmatterFields are the WARN-level retrieval fields every
// pipeline artifact's YAML frontmatter should carry (see any
// shared/contracts/*-contract.md "Retrieval Frontmatter (WARN)" section).
var retrievalFrontmatterFields = []string{
	"feature", "bounded_context", "domain_terms", "files_touched",
	"issue_refs", "linked_adrs", "linked_kis",
}

// requiredSectionPattern extracts the backticked heading from a contract's
// "## Required Sections" bullet line, e.g. "- `## Summary`".
var requiredSectionPattern = regexp.MustCompile("^-\\s+`(#+ [^`]+)`\\s*$")

// ArtifactViolation is one typed structural violation.
type ArtifactViolation struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

// ArtifactValidationResult is the outcome of validating one pipeline artifact
// against its inter-agent contract.
type ArtifactValidationResult struct {
	ArtifactPath     string              `json:"artifactPath"`
	ContractPath     string              `json:"contractPath"`
	Status           string              `json:"status"`
	RequiredSections []string            `json:"requiredSections"`
	PresentSections  []string            `json:"presentSections,omitempty"`
	Violations       []ArtifactViolation `json:"violations,omitempty"`
	Warnings         []string            `json:"warnings,omitempty"`
}

// ArtifactContractAnalyzer validates pipeline artifacts against the
// structural half of their contracts: required heading presence (exact text
// and level) plus WARN-level retrieval frontmatter checks. Contract-specific
// prose content rules and qualitative judgment stay with the
// validate-artifact skill and the human checkpoints.
type ArtifactContractAnalyzer struct{}

// NewArtifactContractAnalyzer returns a stateless analyzer.
func NewArtifactContractAnalyzer() *ArtifactContractAnalyzer {
	return &ArtifactContractAnalyzer{}
}

// Validate checks the artifact at artifactPath against the contract at
// contractPath and returns typed violations. It fails with an error (not a
// FAIL result) when either file is unreadable or the contract declares no
// "## Required Sections" list — the three frontmatter contracts are out of
// this analyzer's scope.
func (a *ArtifactContractAnalyzer) Validate(artifactPath, contractPath string) (*ArtifactValidationResult, error) {
	required, err := requiredSections(contractPath)
	if err != nil {
		return nil, err
	}
	artifact, err := os.ReadFile(artifactPath)
	if err != nil {
		return nil, fmt.Errorf("read artifact: %w", err)
	}
	return buildValidationResult(artifactPath, contractPath, required, string(artifact)), nil
}

func buildValidationResult(artifactPath, contractPath string, required []string, artifact string) *ArtifactValidationResult {
	present, violations := checkSections(required, artifactHeadings(artifact))
	result := &ArtifactValidationResult{
		ArtifactPath:     artifactPath,
		ContractPath:     contractPath,
		Status:           statusFor(violations),
		RequiredSections: required,
		PresentSections:  present,
		Violations:       violations,
		Warnings:         frontmatterWarnings(artifact),
	}
	return result
}

func statusFor(violations []ArtifactViolation) string {
	if len(violations) > 0 {
		return ArtifactStatusFail
	}
	return ArtifactStatusPass
}

func checkSections(required []string, headings map[string]bool) ([]string, []ArtifactViolation) {
	var present []string
	var violations []ArtifactViolation
	for _, heading := range required {
		if headings[heading] {
			present = append(present, heading)
			continue
		}
		violations = append(violations, ArtifactViolation{Type: ViolationMissingSection, Detail: heading})
	}
	return present, violations
}

// requiredSections parses the backticked heading bullets under the contract's
// "## Required Sections" heading, stopping at the next `## ` heading.
func requiredSections(contractPath string) ([]string, error) {
	contract, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, fmt.Errorf("read contract: %w", err)
	}
	sections := collectRequiredSections(string(contract))
	if len(sections) == 0 {
		return nil, fmt.Errorf("contract %s declares no \"## Required Sections\" list — frontmatter contracts are validated by the validate-artifact skill, not this tool", contractPath)
	}
	return sections, nil
}

func collectRequiredSections(contract string) []string {
	var sections []string
	inList := false
	for _, line := range strings.Split(contract, "\n") {
		trimmed := strings.TrimRight(line, " \t")
		inList = requiredListState(inList, trimmed)
		if match := requiredSectionPattern.FindStringSubmatch(trimmed); inList && match != nil {
			sections = append(sections, match[1])
		}
	}
	return sections
}

func requiredListState(inList bool, line string) bool {
	if strings.HasPrefix(line, "## Required Sections") {
		return true
	}
	if strings.HasPrefix(line, "## ") {
		return false
	}
	return inList
}

// artifactHeadings collects every markdown heading line in the artifact,
// ignoring lines inside fenced code blocks.
func artifactHeadings(artifact string) map[string]bool {
	headings := map[string]bool{}
	inFence := false
	for _, line := range strings.Split(artifact, "\n") {
		trimmed := strings.TrimRight(line, " \t")
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if !inFence && strings.HasPrefix(trimmed, "#") {
			headings[trimmed] = true
		}
	}
	return headings
}

// frontmatterWarnings applies the WARN-level retrieval frontmatter checks:
// one warning when the block is absent, one per missing field when present.
func frontmatterWarnings(artifact string) []string {
	block, ok := frontmatterBlock(artifact)
	if !ok {
		return []string{"missing retrieval frontmatter block — add a YAML `---` block at the top of the file"}
	}
	var warnings []string
	for _, field := range retrievalFrontmatterFields {
		if !frontmatterHasField(block, field) {
			warnings = append(warnings, fmt.Sprintf("retrieval frontmatter field `%s` is missing", field))
		}
	}
	return warnings
}

func frontmatterBlock(artifact string) (string, bool) {
	lines := strings.Split(artifact, "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], " \t") != "---" {
		return "", false
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], " \t") == "---" {
			return strings.Join(lines[1:i], "\n"), true
		}
	}
	return "", false
}

func frontmatterHasField(block, field string) bool {
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(line, field+":") {
			return true
		}
	}
	return false
}
