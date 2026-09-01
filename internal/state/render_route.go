package state

import "fmt"

// renderRoute writes route.md: which stages run, which do not, and why.
// This is the document a human reads at the design gate before approving
// what the run will do — and, because approving it binds its digest,
// editing it is what forces a re-approval.
func renderRoute(route Route) string {
	doc := &document{}
	doc.title("Delivery Route: " + route.Feature)
	doc.section("## Stages", routeLines(route.Decisions))
	doc.section("## Skipped", skippedLines(route.Decisions))
	doc.section("## How this was decided", []string{
		"Computed from `analysis.md` after the analyst stage, by predicates in `internal/state/` —",
		"not by a model re-reading the analysis (roadmap L3.0).",
		"",
		"Editing this file resets the design gate's approval: the run will halt until a human",
		"approves the route as it then stands.",
	})
	return doc.String()
}

func routeLines(decisions []RouteDecision) []string {
	lines := make([]string, 0, len(decisions)+2)
	lines = append(lines, "| Stage | Runs | Why |", "|---|---|---|")
	for _, decision := range decisions {
		lines = append(lines, fmt.Sprintf("| `%s` | %s | %s |", decision.Stage, runsLabel(decision.Included), decision.Reason))
	}
	return lines
}

func runsLabel(included bool) string {
	if included {
		return "yes"
	}
	return "**no**"
}

func skippedLines(decisions []RouteDecision) []string {
	lines := make([]string, 0)
	for _, decision := range decisions {
		if !decision.Included {
			lines = append(lines, fmt.Sprintf("- `%s` — %s", decision.Stage, decision.Reason))
		}
	}
	return lines
}
