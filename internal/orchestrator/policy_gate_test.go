package orchestrator_test

import (
	"context"
	"errors"
	"testing"

	"encoding/json"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/policy"
	"github.com/orieken/loom/internal/provider/mock"
	"github.com/orieken/loom/internal/state"
)

func approvingPolicy(t *testing.T, gate, condition string) []policy.Policy {
	t.Helper()
	raw := "name: approve-everything\nversion: \"1.0\"\nmatcher:\n  gate: " + gate +
		"\ncondition:\n" + condition + "action:\n  type: auto-approve\n  reason: \"looks fine\"\n"
	parsed, err := policy.Parse("approve.policy.yaml", []byte(raw))
	if err != nil {
		t.Fatalf("parse policy: %v", err)
	}
	return []policy.Policy{parsed}
}

// reviewedPlan produces a typed review before a gated stage, so a policy
// condition can actually be ANSWERED. Without typed state upstream every
// fact is unknown, and a test asserting "a matching policy does not open a
// gate" would pass while nothing matched — which is what the first draft of
// this file did.
func reviewedPlan() orchestrator.Plan {
	return orchestrator.Plan{
		Name: "reviewed-plan",
		Stages: []orchestrator.Stage{
			{ID: "code-reviewer", Agent: "code-reviewer", StateKind: string(state.KindReview), Timeout: 5 * time.Second},
			{ID: "developer", Agent: "developer", Gate: "confirm-design", Timeout: 5 * time.Second},
		},
	}
}

func reviewedHarness(t *testing.T) (*orchestrator.Executor, *orchestrator.StateStore, orchestrator.StageInput) {
	t.Helper()
	approved := mock.SampleReview(state.VerdictApproved)
	executor, _, store, input := newHarness(t, map[string]mock.Script{
		"code-reviewer": {Payload: mustJSON(t, approved)},
	})
	return executor, store, input
}

func mustJSON(t *testing.T, document any) []byte {
	t.Helper()
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	return raw
}

// verdictCondition is answerable: the review stage above produces it.
const verdictCondition = "  codeReviewer.verdict:\n    equals: \"APPROVED\"\n"

// The property this whole cut exists to preserve: evaluating is not
// approving. The policy below genuinely MATCHES — the review verdict is
// available and equals APPROVED — and the gate must still stop, because no
// run has yet demonstrated the evaluator decides what a human would.
func TestAMatchingPolicyDoesNotOpenAGate(t *testing.T) {
	executor, store, input := reviewedHarness(t)
	executor.WithPolicies(approvingPolicy(t, "confirm-design", verdictCondition))
	var reported []policy.Decision
	executor.OnPolicyDecision(func(decision policy.Decision) { reported = append(reported, decision) })

	err := executor.Run(context.Background(), reviewedPlan(), input)

	if !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v — a matching policy must not skip the gate", err)
	}
	// Precondition: the policy really did match. Without this the test
	// passes just as well when nothing matched, which proves nothing.
	if len(reported) != 1 || !reported[0].WouldAutoApprove() {
		t.Fatalf("precondition failed: policy did not match, decision = %+v", reported)
	}
	recorded := mustLoad(t, store)
	if len(recorded.Approvals) != 0 {
		t.Errorf("approvals = %+v, want none: a policy cannot approve anything in this build", recorded.Approvals)
	}
	if recorded.Stages["developer"].Status != orchestrator.StageStatusWaitingApproval {
		t.Errorf("developer status = %q, want WAITING_APPROVAL", recorded.Stages["developer"].Status)
	}
	if recorded.PolicyDecisions[0].Honoured {
		t.Error("the record claims the decision was honoured")
	}
}

// The decision is recorded even though it is not acted on — that record is
// the audit trail approval-gates.md asserted since v3.3 and never had.
func TestThePolicyDecisionIsRecordedAndReported(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})
	executor.WithPolicies(approvingPolicy(t, "confirm-design", "  testsPass: true\n"))
	var reported []policy.Decision
	executor.OnPolicyDecision(func(decision policy.Decision) { reported = append(reported, decision) })

	if err := executor.Run(context.Background(), gatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v", err)
	}

	decisions := mustLoad(t, store).PolicyDecisions
	if len(decisions) != 1 {
		t.Fatalf("recorded %d decisions, want one for confirm-design", len(decisions))
	}
	if decisions[0].Honoured {
		t.Error("the record claims the decision was honoured; nothing acts on decisions in this build")
	}
	if decisions[0].Gate != "confirm-design" {
		t.Errorf("record names gate %q", decisions[0].Gate)
	}
	if len(reported) != 1 {
		t.Errorf("the CLI callback fired %d times, want once", len(reported))
	}
}

// Without test results in run state, the policy's own condition cannot be
// answered — and the record must name the fact rather than reporting a
// bare "no match".
func TestAnUnansweredConditionNamesTheMissingFact(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})
	executor.WithPolicies(approvingPolicy(t, "confirm-design", "  testsPass: true\n"))

	if err := executor.Run(context.Background(), gatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v", err)
	}

	decisions := mustLoad(t, store).PolicyDecisions
	outcome := decisions[0].Policies[0]
	if outcome.Outcome != string(policy.OutcomeUnknown) {
		t.Fatalf("outcome = %q, want UNKNOWN — this plan produces no QA state", outcome.Outcome)
	}
	if len(outcome.Missing) == 0 || outcome.Missing[0] != string(policy.FieldTestsPass) {
		t.Errorf("missing = %v, want testsPass named", outcome.Missing)
	}
}

// The decision reaches the timeline, which is what makes it auditable
// without reading a JSON file.
func TestThePolicyDecisionReachesTheTimeline(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})
	executor.WithPolicies(approvingPolicy(t, "confirm-design", "  testsPass: true\n"))

	if err := executor.Run(context.Background(), gatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v", err)
	}

	found := false
	for _, event := range readTimeline(t, store) {
		if event.Kind == orchestrator.EventPolicyEvaluated {
			found = true
			if event.Gate != "confirm-design" {
				t.Errorf("event gate = %q", event.Gate)
			}
		}
	}
	if !found {
		t.Errorf("no %q event on the timeline", orchestrator.EventPolicyEvaluated)
	}
}

// A run with no policies must behave exactly as before, and record nothing.
func TestNoPoliciesChangesNothing(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})

	if err := executor.Run(context.Background(), gatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v", err)
	}
	if decisions := mustLoad(t, store).PolicyDecisions; len(decisions) != 0 {
		t.Errorf("recorded %+v with no policies loaded", decisions)
	}
}
