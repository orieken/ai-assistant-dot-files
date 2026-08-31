package state

// Rendering turns typed state back into the markdown a human reads. The
// view is an output, never the transport (roadmap L2.9) — nothing parses it
// back, and it is not digest-tracked, so editing it cannot corrupt a run.
//
// The filenames are the contract's, not the stage's, because the stages
// that are still untyped were told to read `analysis.md`: rendering under
// the contract name keeps the typed hop invisible to them.

import (
	"fmt"
	"strings"
)

// ViewFileName returns the markdown file a state kind renders to.
func ViewFileName(kind Kind) (string, bool) {
	switch kind {
	case KindAnalysis:
		return "analysis.md", true
	case KindArchitecture:
		return "architecture-notes.md", true
	default:
		return "", false
	}
}

// RenderView renders a validated payload as markdown, returning the file
// name it belongs in and the document body.
func RenderView(kind Kind, payload []byte) (string, string, error) {
	name, ok := ViewFileName(kind)
	if !ok {
		return "", "", fmt.Errorf("no markdown view defined for state kind %q", kind)
	}
	decoded, err := Decode(kind, payload)
	if err != nil {
		return "", "", err
	}
	body, err := renderDocument(kind, decoded)
	if err != nil {
		return "", "", err
	}
	return name, body, nil
}

func renderDocument(kind Kind, decoded Validatable) (string, error) {
	switch typed := decoded.(type) {
	case *AnalysisState:
		return renderAnalysis(*typed), nil
	case *ArchitectureState:
		return renderArchitecture(*typed), nil
	default:
		return "", fmt.Errorf("no renderer for state kind %q", kind)
	}
}

// document accumulates markdown. Its methods keep each renderer a flat list
// of sections instead of a nest of string concatenation.
type document struct {
	builder strings.Builder
}

// frontmatter writes the retrieval block validate-artifact checks for. All
// seven fields are always emitted, empty lists included: a missing field is
// a WARN, while an empty list is the fact that there are none.
func (d *document) frontmatter(meta frontmatter) {
	d.builder.WriteString("---\n")
	fmt.Fprintf(&d.builder, "feature: %q\n", meta.feature)
	fmt.Fprintf(&d.builder, "bounded_context: %q\n", meta.boundedContext)
	fmt.Fprintf(&d.builder, "domain_terms: %s\n", yamlList(meta.domainTerms))
	fmt.Fprintf(&d.builder, "files_touched: %s\n", yamlList(meta.filesTouched))
	fmt.Fprintf(&d.builder, "issue_refs: %s\n", yamlList(meta.issueRefs))
	fmt.Fprintf(&d.builder, "linked_adrs: %s\n", yamlList(meta.linkedADRs))
	fmt.Fprintf(&d.builder, "linked_kis: %s\n", yamlList(meta.linkedKIs))
	d.builder.WriteString("---\n\n")
}

// yamlList renders a flow-style sequence, which keeps the block compact and
// unambiguous for both empty and populated lists.
func yamlList(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("%q", value))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func (d *document) title(text string) {
	fmt.Fprintf(&d.builder, "# %s\n", text)
}

// section writes a heading followed by its lines, or the explicit "None"
// the contracts ask for when a section has no content — an empty section
// and a section that says there is nothing are different facts.
func (d *document) section(heading string, lines []string) {
	fmt.Fprintf(&d.builder, "\n%s\n\n", heading)
	if len(lines) == 0 {
		d.builder.WriteString("None\n")
		return
	}
	for _, line := range lines {
		d.builder.WriteString(line + "\n")
	}
}

func (d *document) String() string { return d.builder.String() }

func bullets(items []string) []string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, "- "+item)
	}
	return lines
}

func numbered(items []string) []string {
	lines := make([]string, 0, len(items))
	for i, item := range items {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, item))
	}
	return lines
}
