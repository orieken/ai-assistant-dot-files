package orchestrator

// The durable record of a policy decision (roadmap L2.16).
//
// It carries what would have happened as well as what did, because those
// differ in this build and a record that hid the difference would be
// useless for deciding whether to close the gap. `Honoured` is false on
// every record written today; when a later change starts acting on
// decisions, the field is already there and the history is comparable
// across the boundary.

import (
	"strings"
	"time"

	"github.com/orieken/loom/internal/policy"
)

// PolicyRecord is one gate's policy decision.
type PolicyRecord struct {
	Gate string    `json:"gate"`
	At   time.Time `json:"at"`
	// Effect is what the matching policies asked for.
	Effect string `json:"effect"`
	// Honoured is whether the executor acted on it. Always false today:
	// this build records decisions and still halts for a human.
	Honoured bool `json:"honoured"`
	// Conflict names policies that disagreed, when they did.
	Conflict []string        `json:"conflict,omitempty"`
	Policies []PolicyOutcome `json:"policies"`
}

// PolicyOutcome is one policy's result within a decision.
type PolicyOutcome struct {
	Name    string `json:"name"`
	Action  string `json:"action"`
	Outcome string `json:"outcome"`
	// Missing names the condition fields that could not be answered, which
	// is why a policy did not match when it did not.
	Missing []string `json:"missing,omitempty"`
	Source  string   `json:"source,omitempty"`
}

func newPolicyRecord(decision policy.Decision) PolicyRecord {
	record := PolicyRecord{
		Gate: string(decision.Gate), At: time.Now().UTC(),
		Effect: string(decision.Effect), Honoured: false, Conflict: decision.Conflict,
	}
	for _, result := range decision.Policies {
		record.Policies = append(record.Policies, newPolicyOutcome(result))
	}
	return record
}

func newPolicyOutcome(result policy.Result) PolicyOutcome {
	outcome := PolicyOutcome{
		Name: result.Name, Action: string(result.Action),
		Outcome: string(result.Outcome), Source: result.Source,
	}
	for _, field := range result.Missing {
		outcome.Missing = append(outcome.Missing, string(field))
	}
	return outcome
}

// policyEvent puts the decision on the timeline as one line. `policy.evaluated`
// is the event approval-gates.md said was emitted for every decision — it
// finally is, which is what earns it a place in the vocabulary (L3.9's rule:
// only emitted types).
func policyEvent(decision policy.Decision) Event {
	return Event{
		Kind: EventPolicyEvaluated, Gate: string(decision.Gate),
		Reason: decision.Summary(), Correction: missingSummary(decision),
	}
}

// missingSummary lists the facts no policy could see, so the timeline line
// says why nothing matched rather than only that nothing did.
func missingSummary(decision policy.Decision) string {
	seen := map[string]bool{}
	names := make([]string, 0)
	for _, result := range decision.Policies {
		for _, field := range result.Missing {
			if !seen[string(field)] {
				seen[string(field)] = true
				names = append(names, string(field))
			}
		}
	}
	if len(names) == 0 {
		return ""
	}
	return "unknown: " + strings.Join(names, ",")
}
