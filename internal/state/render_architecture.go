package state

import "fmt"

// renderArchitecture writes architecture-notes.md with the exact headings
// architecture-contract.md requires, in contract order.
func renderArchitecture(architecture ArchitectureState) string {
	doc := &document{}
	doc.title("Architecture Notes: " + architecture.Feature)
	doc.section("## Structural Decisions", decisionLines(architecture.StructuralDecisions))
	doc.section("## Component Placement", placementLines(architecture.ComponentPlacement))
	doc.section("## Bounded Context", contextLines(architecture.BoundedContext))
	doc.section("## Stability Design", stabilityLines(architecture.StabilityDesign))
	doc.section("## Observability Design", observabilityLines(architecture.ObservabilityDesign))
	doc.section("## Layer Boundary Checks", boundaryLines(architecture.LayerBoundaryChecks))
	doc.section("## Anti-Pattern Check", antiPatternLines(architecture.AntiPatternChecks))
	doc.section("## Fitness Functions", fitnessLines(architecture.FitnessFunctions))
	doc.section("## Refactoring Opportunities", refactoringLines(architecture.RefactoringOpportunities))
	doc.section("## Developer Handoff Notes", bullets(architecture.DeveloperHandoffNotes))
	doc.section("## Open Architectural Questions", bullets(architecture.OpenQuestions))
	return doc.String()
}

// decisionLines renders each decision with the **Fitness Function** line
// architecture-contract.md's validation rule looks for, or the explicit
// judgment-only marker that is its documented exception.
func decisionLines(decisions []StructuralDecision) []string {
	lines := make([]string, 0, len(decisions)*4)
	for _, decision := range decisions {
		lines = append(lines,
			"### "+decision.Decision,
			"- **Rationale**: "+decision.Rationale)
		lines = append(lines, tradeOffLines(decision.TradeOffs)...)
		lines = append(lines, fitnessLine(decision))
	}
	return lines
}

func tradeOffLines(tradeOffs []string) []string {
	if len(tradeOffs) == 0 {
		return nil
	}
	lines := []string{"- **Trade-offs**:"}
	for _, tradeOff := range tradeOffs {
		lines = append(lines, "  - "+tradeOff)
	}
	return lines
}

func fitnessLine(decision StructuralDecision) string {
	if decision.Fitness == nil {
		return "- **Fitness Function**: judgment-only — " + decision.Rationale
	}
	return fmt.Sprintf("- **Fitness Function**: %s — verified by %s",
		decision.Fitness.Property, decision.Fitness.Verification)
}

func placementLines(placements []ComponentPlacement) []string {
	lines := make([]string, 0, len(placements))
	for _, placement := range placements {
		lines = append(lines, fmt.Sprintf("- `%s` → %s layer, `%s`%s",
			placement.Component, placement.Layer, placement.Package, reasonSuffix(placement.Reason)))
	}
	return lines
}

func reasonSuffix(reason string) string {
	if reason == "" {
		return ""
	}
	return " — " + reason
}

func stabilityLines(patterns []StabilityPattern) []string {
	lines := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		lines = append(lines, fmt.Sprintf("- **%s** on %s%s",
			pattern.Pattern, pattern.AppliesTo, reasonSuffix(pattern.Rationale)))
	}
	return lines
}

func observabilityLines(signals []ObservabilitySignal) []string {
	lines := make([]string, 0, len(signals))
	for _, signal := range signals {
		lines = append(lines, fmt.Sprintf("- %s emitted from %s%s",
			signal.Signal, signal.EmittedFrom, cardinalitySuffix(signal.Cardinality)))
	}
	return lines
}

func cardinalitySuffix(cardinality string) string {
	if cardinality == "" {
		return ""
	}
	return " (cardinality: " + cardinality + ")"
}

func boundaryLines(checks []BoundaryCheck) []string {
	lines := make([]string, 0, len(checks))
	for _, check := range checks {
		lines = append(lines, fmt.Sprintf("- [%s] %s%s", check.Verdict, check.Rule, reasonSuffix(check.Notes)))
	}
	return lines
}

func antiPatternLines(checks []AntiPatternCheck) []string {
	lines := make([]string, 0, len(checks))
	for _, check := range checks {
		lines = append(lines, fmt.Sprintf("- %s: %s%s", check.Pattern, foundLabel(check.Found), reasonSuffix(check.Notes)))
	}
	return lines
}

func foundLabel(found bool) string {
	if found {
		return "FOUND"
	}
	return "not found"
}

func refactoringLines(opportunities []RefactoringOpportunity) []string {
	lines := make([]string, 0, len(opportunities))
	for _, opportunity := range opportunities {
		lines = append(lines, fmt.Sprintf("- %s on `%s`%s",
			opportunity.Operation, opportunity.Target, reasonSuffix(opportunity.Reason)))
	}
	return lines
}
