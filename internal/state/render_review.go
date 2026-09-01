package state

import "fmt"

// renderReview writes code-review-report.md with the exact headings
// review-contract.md requires. `## Overall Status` still carries the
// bolded literal string: the markdown pipeline greps for it, and the typed
// verdict replaces that parse only for runs under the executor.
func renderReview(review ReviewState) string {
	doc := &document{}
	doc.frontmatter(review.frontmatter())
	doc.title("Code Review: " + review.Feature)
	doc.section("## Overall Status", []string{"**" + review.Verdict.literal() + "**"})
	doc.section("## Design Narrative", []string{review.DesignNarrative})
	doc.section("## Design Score", scoreLines(review.DesignScore))
	doc.section("## Security Surface", bullets(review.SecuritySurface))
	doc.section("## Performance Surface", bullets(review.PerformanceSurface))
	doc.section("## Test Design Review", bullets(review.TestDesignReview))
	doc.section("## Verification of Developer Self-Review", bullets(review.SelfReviewCheck))
	doc.section("## Feedback for the Developer", findingLines(review.Findings))
	return doc.String()
}

func scoreLines(score DesignScore) []string {
	lines := make([]string, 0, 4)
	for _, dimension := range score.dimensions() {
		lines = append(lines, fmt.Sprintf("- **%s** %d/5", dimension.Name, dimension.Score))
	}
	return lines
}

func findingLines(findings []Finding) []string {
	lines := make([]string, 0, len(findings)*5)
	for index, finding := range findings {
		lines = append(lines,
			fmt.Sprintf("### %d. %s%s", index+1, finding.Operation, blockingSuffix(finding.Blocking)),
			"- **File**: `"+finding.File+"`",
			"- **Smell**: "+finding.Smell,
			"- **Instruction**: "+finding.Instruction)
	}
	return lines
}

func blockingSuffix(blocking bool) string {
	if blocking {
		return ""
	}
	return " (non-blocking)"
}

func (r ReviewState) frontmatter() frontmatter {
	return newFrontmatter(r.Feature, "", r.Retrieval, findingFiles(r.Findings))
}

// findingFiles derives files_touched from what the review actually points
// at — the closest thing this document has to a file list.
func findingFiles(findings []Finding) []string {
	files := make([]string, 0, len(findings))
	for _, finding := range findings {
		if finding.File != "" {
			files = append(files, finding.File)
		}
	}
	return files
}
