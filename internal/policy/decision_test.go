package policy_test

import (
	"strings"
	"testing"

	"github.com/orieken/loom/internal/policy"
)

func policyYAML(name, action, condition string) string {
	return "name: " + name + "\nversion: \"1.0\"\nmatcher:\n  gate: git-commit\ncondition:\n" +
		condition + "action:\n  type: " + action + "\n  reason: \"because\"\n"
}

func parseAll(t *testing.T, documents ...string) []policy.Policy {
	t.Helper()
	policies := make([]policy.Policy, 0, len(documents))
	for index, document := range documents {
		parsed, err := policy.Parse("p.policy.yaml", []byte(document))
		if err != nil {
			t.Fatalf("parse document %d: %v", index, err)
		}
		policies = append(policies, parsed)
	}
	return policies
}

// The rule the schema always stated and nothing ever applied.
func TestRequireHumanBeatsAutoApprove(t *testing.T) {
	policies := parseAll(t,
		policyYAML("approve-it", "auto-approve", "  testsPass: true\n"),
		policyYAML("stop-it", "require-human", "  testsPass: true\n"))

	decision := policy.Decide(policies, fullContext())

	if decision.Effect != policy.ActionRequireHuman {
		t.Errorf("effect = %q, want require-human when policies conflict", decision.Effect)
	}
	if len(decision.Conflict) != 2 {
		t.Errorf("conflict = %v, want both policies named", decision.Conflict)
	}
}

// The absence of a policy is not permission.
func TestNothingMatchingMeansRequireHuman(t *testing.T) {
	policies := parseAll(t, policyYAML("approve-it", "auto-approve",
		"  codeReviewer.verdict:\n    equals: \"CHANGES_REQUESTED\"\n"))

	decision := policy.Decide(policies, fullContext())

	if decision.Effect != policy.ActionRequireHuman {
		t.Errorf("effect = %q, want require-human when nothing matched", decision.Effect)
	}
	if decision.WouldAutoApprove() {
		t.Error("WouldAutoApprove is true with no matching policy")
	}
}

// An unknown fact must not become an auto-approval. This is the single
// most important property in the package: a policy that cannot be
// evaluated must never be treated as satisfied.
func TestAnUnknownConditionNeverAutoApproves(t *testing.T) {
	policies := parseAll(t, policyYAML("approve-it", "auto-approve",
		"  testsPass: true\n  diffLines:\n    lessThan: 500\n"))

	decision := policy.Decide(policies, fullContext())

	if decision.WouldAutoApprove() {
		t.Fatal("a policy with an unevaluable check would have auto-approved")
	}
	if len(decision.Policies) != 1 || decision.Policies[0].Outcome != policy.OutcomeUnknown {
		t.Fatalf("result = %+v, want a single UNKNOWN outcome", decision.Policies)
	}
	if len(decision.Policies[0].Missing) == 0 {
		t.Error("the decision does not name the fact it could not see")
	}
}

func TestAMatchingApprovalReportsWouldAutoApprove(t *testing.T) {
	policies := parseAll(t, policyYAML("approve-it", "auto-approve",
		"  testsPass: true\n  codeReviewer.verdict:\n    equals: \"APPROVED\"\n"))

	decision := policy.Decide(policies, fullContext())

	if !decision.WouldAutoApprove() {
		t.Fatalf("effect = %q, want auto-approve; results %+v", decision.Effect, decision.Policies)
	}
	// The summary must say the gate still stops, or a reader will assume
	// the opposite from the word "auto-approve".
	if !strings.Contains(decision.Summary(), "still stops for a human") {
		t.Errorf("summary %q does not say the run still halts", decision.Summary())
	}
}

// A gate no policy watches produces no decision at all, rather than an
// empty one that looks like a considered outcome.
func TestAnUnwatchedGateProducesNoResults(t *testing.T) {
	policies := parseAll(t, policyYAML("approve-it", "auto-approve", "  testsPass: true\n"))
	context := fullContext()
	context.Gate = policy.GateConfirmShip

	if decision := policy.Decide(policies, context); len(decision.Policies) != 0 {
		t.Errorf("results = %+v, want none for an unwatched gate", decision.Policies)
	}
}

// The executor's own gates must be targetable, or no policy could ever be
// evaluated by `loom run` — the two vocabularies had no overlap at all.
func TestExecutorGatesArePolicyEligible(t *testing.T) {
	for _, gate := range []policy.GateID{
		policy.GateConfirmDesign, policy.GateConfirmSecurity,
		policy.GateConfirmShip, policy.GateConfirmUnresolvedReview,
	} {
		t.Run(string(gate), func(t *testing.T) {
			if err := policy.CheckGate(gate); err != nil {
				t.Errorf("gate %q is not policy-eligible: %v", gate, err)
			}
		})
	}
}
