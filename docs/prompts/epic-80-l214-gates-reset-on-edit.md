# Epic 80 — Any Edit Resets the Gate, In Code (L2.14)

Source: epic 79 completion, 2026-08-31. Operationalizes roadmap item **L2.14** (BUILD-ROADMAP.md,
KERNEL — Human-in-the-Loop). Every gate in `approval-gates.md` declares "Reset condition: any edit
to the pending artifact resets the gate" — eight times, enforced zero times. L2.12 made the
executor compute and verify digests; L2.13 made it halt at gates and record approvals. This epic
joins them: an approval binds to the digests it was given, and an edit after approval invalidates
it. It blocks **L4.5** (capturing the human-correction signal).

**What this closes.** Today L2.12 catches an edited artifact by demoting its stage so it re-runs —
but the *approval* survives, so the gated stage proceeds on a human decision made about content
that no longer exists. That is the specific hole; it is currently documented as a known gap in
`approval-gates.md` and `cmd/loom/README.md`, and this epic is what lets those notes come out.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — item L2.14 (the spec), plus L2.15 and L2.16 for the boundary
   of what NOT to build. Respect Blocked-by: L2.12 and L2.13 are both done.
2. `internal/orchestrator/gate.go` — `Approval`, `Executor.Approve`, `checkGate` (the pre-stage
   barrier this epic extends); `integrity.go` — digest verification and the staleness cascade;
   `timeline.go` — the event log a new event kind joins.
3. `shared/rules/approval-gates.md` — the eight prose gates, their identical reset condition, and
   the "Executor Enforcement (L2.13)" section whose **Not yet** paragraph this epic retires.
4. `shared/telemetry/event-schema.md` — the `gate_decision` event and its `edited_then_approved`
   outcome, including the "Edit detection heuristic" paragraph that describes a model comparing
   checksums it computed itself.
5. `cmd/loom/cmd/run_gates.go` and `state_run.go` — the two approval channels (TTY prompt,
   `--resume --approve`) and the markdown pipeline's `loom state approve`.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| What an approval binds to | The digests of **every completed stage** at the moment of approval, recorded as `Approval.ArtifactDigests: {stageID: sha256}` | A human approving `confirm-design` is approving the state of the run, not one file. Binding to a single "pending artifact" would let an edit to any earlier artifact slip through, and the prose's singular "the pending artifact" predates a pipeline with fifteen of them |
| When it is re-verified | In `checkGate`, immediately before the gated stage starts: recompute each bound digest and compare | The barrier is already the one place the executor decides whether work may proceed. Verifying anywhere else would leave a window between the check and the action |
| What invalidation does | The approval is **marked invalid**, not deleted — `InvalidatedAt`, `InvalidatedBy` (the stage whose artifact changed) — and the run halts at the gate again, needing a fresh approval | Deleting the record would destroy the audit trail of what was approved and when, which is the same reason epic 78 kept stale stage records instead of dropping them |
| Identical re-runs keep the approval | A stage that re-runs and produces a byte-identical artifact leaves its digest unchanged, so the approval stands | The rule is "any edit resets the gate", not "any re-run". Forcing re-approval for provably identical content would train humans to click through gates, which is the failure this whole workstream exists to prevent |
| No cryptographic token | The approval record *is* the token: a digest set in state the executor owns. No signing, no expiry | Signing matters when an approval crosses a trust boundary — a webhook or queue channel, which is explicitly later work. Adding crypto now would be ceremony over a file only this process writes |
| Stages completed after approval | Not bound. Only stages already complete when the human approved are part of what they saw | Binding future work to a past decision is incoherent; the next gate is what covers later stages |
| Markdown pipeline | `loom state approve` records the same digest set, and `loom state verify` reports an approval as INVALIDATED when a bound artifact changed. It **detects**; it cannot enforce | This is the `edited_then_approved` signal the telemetry spec describes, computed in Go instead of by the model that wrote the checksums. Enforcement needs the gated action running under the executor — say so plainly rather than implying the markdown pipeline gained a guarantee |
| Approving into a stale run | `--approve` is **refused up front** when verification is about to demote a bound stage: the error names the changed artifact and tells the human to resume first, then approve | Recording an approval that the same command's verification immediately invalidates wastes the human's decision and reads like a bug. Moving approval to *after* verification was considered and rejected: the demoted stage is then STALE rather than completed, so it would bind nothing, the stage would re-run, and the gate would pass content nobody saw — the exact hole this epic closes |
| Timeline | A new `gate.invalidated` event carrying the gate and the stage whose artifact changed | The timeline already records `gate.waiting` and `gate.approved`; an invalidation is the third thing that happens to a gate, and L4.5 will mine exactly this |
| Telemetry | `shared/telemetry/event-schema.md` gets a note that the executor now detects edited-then-approved structurally, and that its "edit detection heuristic" describes the markdown path only. Do **not** build the opt-in telemetry emitter | The `events.jsonl` emitter is L3.9. Writing it here would smuggle an unrelated roadmap item into a gate epic |
| State schema | `Approval` gaining digest and invalidation fields bumps `StateSchemaVersion` to 5; no migration code — a v4 file is refused with the existing clear message | Same policy as every prior bump, and for the same reason: nothing is deployed anywhere yet |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**58.6%** as of epic 79; raise the floor in `framework-ci.yml`
  if measured coverage rises).
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — Digest-bound approvals in the library — UNBLOCKED

1. `internal/orchestrator/gate.go` (+ `state.go`, `timeline.go`):
   - `Approval.ArtifactDigests map[string]string`, `InvalidatedAt *time.Time`, `InvalidatedBy string`;
     schema version 5.
   - `Executor.Approve` records the digest of every completed stage's artifact at approval time.
     Stages with no artifact record no digest — there is nothing to bind to.
   - `checkGate` re-verifies a bound approval before starting the gated stage: any digest that no
     longer matches (or an artifact that vanished) invalidates the approval, persists the
     invalidation, emits `gate.invalidated`, and halts as if the gate were never approved.
   - `Approve` on a previously invalidated gate is a fresh approval, binding the current digests.
2. Tests (extend the existing harness):
   - Approve `confirm-design`, edit an upstream artifact, resume: the gated stage does **not** run
     and the halt names the artifact that changed. **This is the L2.14 done-when test — name it
     accordingly.**
   - A stage that re-runs to a byte-identical artifact keeps the approval and the run proceeds.
   - Re-approving after invalidation proceeds, and the state shows both the invalidated record and
     the new one.
   - An artifact deleted after approval invalidates it.
   - Editing an artifact completed *after* the approval does not invalidate that approval.
   - The interaction with L2.12 is explicit: an edit both demotes the stage and invalidates the
     approval, and the test asserts both.
   - A v4 state file is refused with the schema message.

**Done when**: editing an artifact between approval and execution causes the execution to be
refused. **Commit** (`feat(orchestrator): approvals bind to artifact digests — edits reset the
gate`), report, PAUSE.

## Phase B — CLI and the markdown pipeline — BLOCKED BY Phase A

1. `cmd/loom/cmd/run_gates.go`: when a gate is invalidated, say which artifact changed and that
   re-approval is required, then halt on the normal non-interactive path (exit code 3) or re-prompt
   on a TTY. The message must make the cause obvious — a human who edited a file on purpose should
   not have to guess why the run stopped.
2. `applyApproveFlag` refuses to record an approval when verification would immediately invalidate
   it, per the design table — name the artifact, say to resume first.
3. `loom state approve` records the same digest set; `loom state verify` reports each approval as
   OK or INVALIDATED (naming the changed artifact) and `loom state show` displays it. Exit non-zero
   on an invalidated approval, as `verify` already does for a stale stage.
4. Tests through the real binary: approve, edit, resume → refused with a message naming the
   artifact; re-approve → proceeds; the markdown path reports INVALIDATED without claiming to have
   prevented anything.

**Done when**: both channels refuse a stale approval and explain why. **Commit** (`feat(loom):
refuse a stale approval and say which artifact changed`), report, PAUSE.

## Phase C — Prose alignment + docs — BLOCKED BY Phase B

1. `shared/rules/approval-gates.md`: the "Executor Enforcement (L2.13)" section's **Not yet**
   paragraph comes out; the reset condition is now enforced for `loom run` and detected for the
   markdown pipeline. Keep the honest scope — the other prose gates are still prompt-discipline.
2. `shared/telemetry/event-schema.md`: note that the executor detects `edited_then_approved`
   structurally and that the checksum heuristic describes the markdown path; do not build the
   emitter (L3.9).
3. `cmd/loom/README.md`: gate-reset behavior; remove L2.14 from the NOT-yet list.
4. `README.md` + `docs/ARCHITECTURE.md`: L2.14 status markers.
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when** (L2.14): the roadmap done-when holds and no doc still describes reset-on-edit as
unenforced. **Commit** (`docs(gates): reset-on-edit is enforced — align prose with L2.14`), report,
PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- The opt-in telemetry emitter or `events.jsonl` (**L3.9**); the `gate_decision` event itself
- Policy-based auto-approval of executor gates (**L2.16**)
- Signed or expiring approval tokens, and any non-CLI approval channel
- `--from-phase`, rollback, or `.history` restoration (**L2.15**)
- Enforcing gates in the markdown pipeline (detection only — enforcement needs the gated action
  running under the executor)
- Migration code for v4 state files

## Report format (end of every phase)

```
## Epic 80 Phase <X> Report
- Roadmap item: L2.14 — Any edit resets the gate
- Blockers verified: <list, with evidence (commit SHAs / files)>
- Commits: <sha> <subject>
- Build/lint/test: go build PASS|FAIL · golangci-lint PASS|FAIL · go test PASS|FAIL (counts) · health-check PASS|FAIL|n/a
- Done-when criterion: <restate it> — MET | NOT MET (why)
- Escalations / open questions: <list or "none">
- Next phase blocked by: <what must land first>
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
