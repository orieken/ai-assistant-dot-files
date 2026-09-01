package state

import "fmt"

// renderSecurity writes security-report.md with the exact headings
// security-contract.md requires, including all six STRIDE sub-headings and
// the `**Fix applied**` line its validation rule looks for.
func renderSecurity(security SecurityState) string {
	doc := &document{}
	doc.frontmatter(security.frontmatter())
	doc.title("Security Review: " + security.Feature)
	doc.section("## Threat Model Summary", []string{security.ThreatModelSummary})
	doc.section("## Dependency Audit", bullets(security.DependencyAudit))
	doc.section("## STRIDE Analysis", nil)
	renderStride(doc, security.Stride)
	doc.section("## Findings", securityFindingLines(security.Findings))
	doc.section("## Files Modified", bullets(security.FilesModified))
	doc.section("## Security Checklist", checklistLines(security.SecurityChecklist))
	doc.section("## Notes for QA", bullets(security.NotesForQA))
	doc.section("## Notes for Tech Writer", bullets(security.NotesForTechWriter))
	return doc.String()
}

// renderStride emits every category's heading whether or not it has
// findings — the contract requires all six to be present.
func renderStride(doc *document, analyses []StrideAnalysis) {
	assessed := make(map[StrideCategory]string, len(analyses))
	for _, analysis := range analyses {
		assessed[analysis.Category] = analysis.Assessed
	}
	for _, category := range strideCategories() {
		doc.section("### "+strideHeading(category), []string{assessed[category]})
	}
}

// strideHeading renders the contract's exact sub-heading text.
func strideHeading(category StrideCategory) string {
	headings := map[StrideCategory]string{
		StrideSpoofing:              "Spoofing",
		StrideTampering:             "Tampering",
		StrideRepudiation:           "Repudiation",
		StrideInformationDisclosure: "Information Disclosure",
		StrideDenialOfService:       "Denial of Service",
		StrideElevationOfPrivilege:  "Elevation of Privilege",
	}
	return headings[category]
}

func securityFindingLines(findings []SecurityFinding) []string {
	lines := make([]string, 0, len(findings)*6)
	for _, finding := range findings {
		lines = append(lines,
			fmt.Sprintf("### [%s] — %s", finding.Severity, finding.Title),
			"**Location**: `"+finding.Location+"`",
			"**Threat category**: "+strideHeading(finding.Category),
			"**Description**: "+finding.Description,
			"**Fix applied**: "+fixOrRecommendation(finding),
			"**Verification**: "+finding.Verification)
	}
	return lines
}

// fixOrRecommendation keeps the template's wording for an unfixed finding.
// Validation already refuses this for CRITICAL and HIGH.
func fixOrRecommendation(finding SecurityFinding) string {
	if finding.FixApplied == "" {
		return "Recommendation only"
	}
	return finding.FixApplied
}

func (s SecurityState) frontmatter() frontmatter {
	return newFrontmatter(s.Feature, "", s.Retrieval, s.FilesModified)
}

// render satisfies the renderable interface.
func (s SecurityState) render() string { return renderSecurity(s) }
