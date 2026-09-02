package policy

// Resolving a single check against the context.
//
// Five of the nine declared fields have no source in run state today —
// diffLines, diffType, dryRunPass, fitnessFunction.allPass, and
// codeReviewer.behaviorChange. They resolve to UNKNOWN by falling through
// to the default, which is deliberate: the vocabulary describes what a
// policy may ask, and the executor answers what it can see. Sourcing the
// rest is its own roadmap item, because "shell to git for a diff size"
// carries questions (diff against what?) that do not belong in an
// evaluator.

import "path/filepath"

func evaluateCheck(check Check, context GateContext) Outcome {
	switch check.Field {
	case FieldReviewVerdict:
		return textOutcome(check, context.ReviewVerdict)
	case FieldSecurityCriticals:
		return numberOutcome(check, context.SecurityCriticals)
	case FieldTestsPass:
		return boolOutcome(check, context.TestsPass)
	case FieldFilePaths:
		return pathsOutcome(check, context)
	default:
		return OutcomeUnknown
	}
}

func textOutcome(check Check, value *string) Outcome {
	if value == nil {
		return OutcomeUnknown
	}
	if check.Operator != OpEquals {
		return OutcomeUnknown
	}
	return truth(*value == check.Text)
}

func numberOutcome(check Check, value *int) Outcome {
	if value == nil {
		return OutcomeUnknown
	}
	actual := float64(*value)
	switch check.Operator {
	case OpEquals:
		return truth(actual == check.Number)
	case OpLessThan:
		return truth(actual < check.Number)
	default:
		return OutcomeUnknown
	}
}

func boolOutcome(check Check, value *bool) Outcome {
	if value == nil {
		return OutcomeUnknown
	}
	switch check.Operator {
	case OpIsTrue, OpEquals:
		return truth(*value == check.Bool)
	default:
		return OutcomeUnknown
	}
}

// pathsOutcome applies a glob across every changed file. PathsKnown is what
// separates "this run changed no files" from "no implementation state was
// recorded" — allMatch over an empty set is vacuously true, and returning
// that for a run nobody can see would be a confident wrong answer.
func pathsOutcome(check Check, context GateContext) Outcome {
	if !context.PathsKnown {
		return OutcomeUnknown
	}
	matched := countMatches(check.Text, context.ChangedPaths)
	switch check.Operator {
	case OpAllMatch:
		return truth(matched == len(context.ChangedPaths))
	case OpNoneMatch:
		return truth(matched == 0)
	case OpAnyMatch:
		return truth(matched > 0)
	default:
		return OutcomeUnknown
	}
}

func countMatches(pattern string, paths []string) int {
	matched := 0
	for _, path := range paths {
		if matchesGlob(pattern, path) {
			matched++
		}
	}
	return matched
}

// matchesGlob supports the `**` the schema's examples use, which
// filepath.Match does not: `**/security/**` must match at any depth.
// A malformed pattern matches nothing rather than erroring — the pattern
// was accepted at load time, and a gate is not the place to discover a
// typo.
func matchesGlob(pattern, path string) bool {
	if ok, err := filepath.Match(pattern, path); err == nil && ok {
		return true
	}
	return matchesDoubleStar(pattern, path)
}

func truth(value bool) Outcome {
	if value {
		return OutcomeTrue
	}
	return OutcomeFalse
}
