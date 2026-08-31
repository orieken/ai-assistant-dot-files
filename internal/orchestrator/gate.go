package orchestrator

import (
	"errors"
	"fmt"
	"os/user"
	"strings"
	"time"
)

// ApprovalMethod records which CLI channel recorded an approval. Provider
// output is never a channel: nothing an agent returns can create an
// approval (roadmap L2.13 — enforcement lives below the model).
type ApprovalMethod string

// The only two approval channels that exist today. Webhook and queue
// channels are later roadmap items.
const (
	ApprovalMethodTTY  ApprovalMethod = "tty"
	ApprovalMethodFlag ApprovalMethod = "flag"
	// ApprovalMethodCLI records an approval entered with `loom state
	// approve` for a markdown-pipeline run. It is an audit record of a human
	// decision, not an enforced gate — the markdown pipeline can proceed
	// without it (roadmap L2.13 covers `loom run` only).
	ApprovalMethodCLI ApprovalMethod = "cli"
)

// Approval is the durable record that a named gate was unlocked by a human.
type Approval struct {
	ApprovedAt time.Time      `json:"approvedAt"`
	Method     ApprovalMethod `json:"method"`
	Approver   string         `json:"approver"`
}

// ErrWaitingApproval is the sentinel every gate halt wraps, so callers can
// branch with errors.Is without depending on the concrete error type.
var ErrWaitingApproval = errors.New("waiting for gate approval")

// WaitingApprovalError names the gate the run stopped at and the stage that
// gate guards. The CLI turns this into a prompt or a resume instruction.
type WaitingApprovalError struct {
	Gate  string
	Stage string
}

func (e *WaitingApprovalError) Error() string {
	return fmt.Sprintf("stage %q is gated by %q: %s", e.Stage, e.Gate, ErrWaitingApproval)
}

// Unwrap makes errors.Is(err, ErrWaitingApproval) true for every gate halt.
func (e *WaitingApprovalError) Unwrap() error { return ErrWaitingApproval }

// IsGateApproved reports whether the named gate has been unlocked for this
// run. One approval unlocks one gate for the lifetime of the run; digest
// binding and reset-on-edit are roadmap item L2.14.
func (s *RunState) IsGateApproved(gate string) bool {
	_, ok := s.Approvals[gate]
	return ok
}

// WaitingGate returns the gate this run is currently halted on, or "" when
// it is not waiting on anything.
func (s *RunState) WaitingGate() string {
	for _, record := range s.Stages {
		if record.Status == StageStatusWaitingApproval {
			return record.Gate
		}
	}
	return ""
}

// Approve records a human approval for the gate the run is currently
// waiting on and persists it atomically. It refuses any gate the run is not
// actually halted at, which is what stops a caller from pre-approving every
// gate in one command line and hollowing out the interrupt.
func (e *Executor) Approve(gate string, method ApprovalMethod) error {
	state, err := e.store.Load()
	if err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("no run state at %s — nothing to approve", e.store.Path())
	}
	if err := checkWaitingOn(state, gate); err != nil {
		return err
	}
	state.RecordApproval(gate, method)
	if err := e.store.Save(state); err != nil {
		return fmt.Errorf("persist approval for gate %q: %w", gate, err)
	}
	return nil
}

// RecordApproval writes an approval for a named gate. Callers that must
// enforce "the run is actually waiting on this gate" go through
// Executor.Approve; this is the plain record used by the markdown
// pipeline, which has no executor barrier to wait at.
func (s *RunState) RecordApproval(gate string, method ApprovalMethod) {
	s.Approvals[gate] = Approval{ApprovedAt: time.Now().UTC(), Method: method, Approver: currentApprover()}
}

func checkWaitingOn(state *RunState, gate string) error {
	waiting := state.WaitingGate()
	if waiting == "" {
		return fmt.Errorf("run is not waiting on any gate — cannot approve %q", gate)
	}
	if waiting != gate {
		return fmt.Errorf("run is waiting on gate %q, not %q", waiting, gate)
	}
	return nil
}

// currentApprover identifies the human at the CLI by OS username. An
// unresolvable user is recorded literally rather than failing the approval.
func currentApprover() string {
	current, err := user.Current()
	if err != nil || strings.TrimSpace(current.Username) == "" {
		return "unknown"
	}
	return current.Username
}
