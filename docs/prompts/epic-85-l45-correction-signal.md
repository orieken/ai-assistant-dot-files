# Epic 85 — Capture the Human-Correction Signal (L4.5)

Source: session discussion, 2026-09-01, following epic 84. Implements roadmap item **L4.5**, whose
three blockers — **L2.13** (gates as interrupts), **L2.14** (approvals bind digests), **L3.8**
(telemetry) — all shipped. It **blocks L4.4** (prompt registry with eval-gated promotion).

The roadmap calls `edited_then_approved` "the single highest-value learning signal in the design",
notes that it is precisely specified, and observes that **it is emitted by nothing**. The framework
describes its own reward signal and never collects it. This epic collects it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — **L4.5**, **L4.4** (what consumes this), **L3.5** (episodic
   memory, where the roadmap says this should eventually live), **L3.9** (owns `events.jsonl`).
2. `internal/orchestrator/gate.go` — `Approval`, `RecordApproval`, `InvalidateApproval`,
   `ArtifactDigests`, `WaitingGate`. L2.14's machinery is the foundation and is already half of
   this.
3. `internal/orchestrator/approval_binding.go` — where invalidation is detected, at verification
   time rather than at the barrier. Epic 80 moved it there for a reason worth not undoing.
4. `internal/orchestrator/loop.go` — `retainIteration`, `IterationsDirName`. **Follow this pattern**
   rather than inventing a second one: it already copies an artifact aside, digests the copy, and
   records it in state so an earlier round is verifiable rather than merely remembered.
5. `internal/orchestrator/timeline.go` — the `EventKind` enum and `Event` struct this extends.
6. `shared/telemetry/event-schema.md` — the `gate_decision` event and its three outcomes. Read the
   Executor note: it already says L2.14 made this detection real and that the *emitter* is L3.9.

## Second discovery (Phase A, 2026-09-01) — where a correction can actually live

Phase A's own test turned out to be vacuous, and finding out why changed Phase B's design.

**A human edit to a tracked artifact does not survive under `loom run`.** Editing a completed
artifact makes its stage STALE (L2.12), so the executor re-runs the stage and the agent's fresh
output overwrites what the human wrote. Verified directly: invocations go from `[analyst]` to
`[analyst analyst]` and the file reverts. The executor treats an edit as drift to correct, not as a
correction to keep.

**And for the seven typed artifacts, the file a human would naturally edit is not tracked at all.**
The artifact is `state/<stage>.json`; `analysis.md` is a *rendered view*, deliberately not
digest-tracked since epic 79 so that editing it cannot corrupt a run. So editing `analysis.md` is
silently ignored and overwritten at the next render, while hand-editing the JSON produces a signal
about an edit the executor then discards. The natural action yields nothing; the unnatural one
yields something worthless.

**The resolution: the correction signal comes from the rendered view.** The property that makes the
view safe is exactly what makes it the right channel — a human can write "here is what this should
have said" into `analysis.md`, the executor will never read it, and the divergence from what was
rendered is a pure corrective signal with nothing to discard and nothing to guard against. For the
eight artifacts that are still markdown, that markdown *is* the tracked artifact, so the capture
happens at detection time before the staleness cascade re-runs the stage. Both cases end up as one
comparison against a retained baseline.

Two things this design must state plainly rather than let someone discover: a view edit is
**advisory by construction** — it changes nothing the pipeline reads — and the human's text survives
only in the retained copy and the diff, because the next render overwrites the view.

## The discovery that shapes this epic

**The schema and L2.14 describe two different sequences, and only one of them is covered.**

- The schema's definition is *gate presented → human edits the artifact → human approves*:
  "checksum changed between gate-presented and gate-approved."
- L2.14's mechanism is *approval granted → artifact edited → approval invalidated → re-approved*.

L2.14 cannot see the schema's sequence at all, because at the moment of the edit **no approval
exists yet to invalidate**. And the schema's sequence is the common one: a human halts at a gate,
reads the artifact, fixes it, and approves. Building only on L2.14's path would collect the rarer
case and miss the signal the design was actually describing.

So the unit is the **presented baseline**: whatever the human was last shown. It is captured when a
gate halts, and refreshed whenever an approval is granted (approving means accepting what is there
now). A correction is a difference between the presented baseline and the content at the moment of
approval.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Where the signal lands | A new `EventKind` on `run-events.jsonl`, alongside `gate.approved` and `gate.invalidated` | The timeline is a Go-owned enum with no prose contradiction, needs nothing from L3.9, and sits exactly where this signal is produced. `events.jsonl` and the `gate_decision` schema stay **entirely L3.9's problem** — writing them here would build half of L3.9 under this item's number, which is the scope-absorption failure epic 83 had to write a roadmap correction for |
| How a diff becomes possible | Retain a copy of every gate-bound artifact when the gate is presented and when an approval is granted, following `retainIteration`'s pattern exactly — copy aside, digest the copy, record it in state | Digests alone answer *that* something changed and never *what*, and L4.4's prompt-variant generation needs to read the correction. L2.17 already established that retaining the artifact is what makes a question about an earlier state answerable |
| What is compared *(revised in Phase A)* | For a typed stage, the **rendered view** (`analysis.md`); for an untyped stage, the markdown artifact itself | The view is where a human's correction can exist without being destroyed or discarded, and it is what a person would open. Comparing only tracked artifacts would miss all seven typed stages and would describe an edit the executor immediately re-runs away |
| Where an untyped artifact is compared | At detection time, before the staleness cascade re-runs the stage | The re-run overwrites the human's text. Capturing after it would record the agent's second attempt as if it were the human's correction — a wrong number wearing the appearance of a measurement |
| What a view edit means | Advisory. It changes nothing the pipeline reads, and the next render overwrites it | Stated rather than discovered. The signal is the *record* of the correction, not the adoption of it — adopting it is L4.4's problem and possibly nobody's |
| What the baseline is | The **presented baseline**: set at gate halt, refreshed at each approval | See the discovery above. Anchoring only to approval would miss the schema's own primary sequence |
| Attribution | To the **producing stage and its agent**, not to the gate | The signal's value is "this agent's output needed human correction". A gate name says where it was caught, not who produced it. `run-state.json` already knows which stage wrote each artifact |
| Scope | `loom run` only. The docs state plainly that the markdown pipeline does not emit it | Same posture as every gate and state item before this. That pipeline's checksums are the model's own except where `loom state verify` runs, so a signal from it would vary in trustworthiness by whether the binary happened to be installed — worse than no signal, because it would be mined as if uniform |
| Diff format and bounds | Unified diff, written to a retained file. The event carries a **bounded summary** — stage, agent, lines added/removed, sections touched, and the diff's own path | An unbounded diff inline would put a whole artifact on one timeline line. The timeline stays greppable; the detail stays on disk |
| Episodic storage | **Not here.** The roadmap says this signal belongs in episodic memory; that store is **L3.5** and does not exist | This epic produces the signal and a durable place to read it from. L3.5 adopts it. The epic must say so and update L3.5's entry |
| Nothing new gates | Detection never blocks a run, and a failure to write the baseline or the diff never fails a stage | This is an observation of a human's action, not a control. A learning signal that can halt a delivery would be turned off |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer, max 4
  parameters (introduce a parameter object rather than a fifth).
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**63.2%** as of epic 84). Integration tests driving the built
  binary are invisible to the coverage tool — if the floor is threatened, add in-process tests,
  never lower the floor.
- Bumping `StateSchemaVersion` is expected in Phase A; note it in the report, because an existing
  run-state file is refused on load and anyone mid-run through the upgrade restarts.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — The presented baseline — UNBLOCKED

1. Retain gate-bound artifacts under a new directory beside `.iterations/`, following
   `retainIteration`: copy, digest the copy, record it in run state against the gate.
2. Capture points: when a gate halts (`waitingRecord` / the barrier path) and when an approval is
   recorded (both the TTY prompt and `--resume --approve`). Refreshing on approval is what makes the
   baseline mean "what the human last accepted".
3. A capture failure is logged and never fails the stage or the approval.
4. Tests: a halt retains exactly the artifacts the gate binds; a second approval refreshes them; a
   re-run producing byte-identical output leaves the baseline unchanged; retained copies carry their
   own digests; a run with no gates retains nothing.

**Done when**: at any gate, the artifact the human was last shown is recoverable from disk and
verifiable by digest. **Commit** (`feat(orchestrator): retain the artifact state a human was shown
at a gate`), report, PAUSE.

## Phase B — Detect the correction and emit it — BLOCKED BY Phase A

1. Extend Phase A's retention to cover the **rendered view** of every typed stage, not only the
   tracked artifact. This is the file a human edits, and it is the comparison target.
2. At approval, compare each retained view (typed) or artifact (untyped) against the presented
   baseline. Unchanged is a plain approval; changed is the correction case.
3. For untyped artifacts, the comparison must happen **before** verification demotes the stage and
   re-runs it — otherwise the recorded "correction" is the agent's second attempt.
4. A new `EventKind` on the timeline carrying the bounded summary from the design table, attributed
   to the **producing stage and agent**. Write the unified diff to a retained file beside the
   baseline and name it in the event.
5. Cover both sequences: edited-before-first-approval (the schema's) and
   edited-after-approval-then-re-approved (L2.14's). Both are the same comparison against the
   presented baseline; the tests must show that explicitly, because the epic exists partly because
   they were conflated.
6. Record the correction on run state too, so it survives without reading the timeline.
7. Tests: a view edit on a typed stage produces the signal and does **not** disturb the run; both
   sequences produce the event; an unchanged approval produces none; the attribution names the
   producing agent rather than the gate; an untyped edit is captured before the re-run rather than
   after; a diff that cannot be written leaves the approval intact.

**Done when**: editing an artifact at a gate produces a stored diff attributable to the producing
agent — L4.5's stated done-when, in full. **Commit** (`feat(orchestrator): emit the human-correction
signal at approval`), report, PAUSE.

## Phase C — Surface it, and draw the boundaries — BLOCKED BY Phase B

1. `loom state timeline` renders the correction legibly; `loom state show` reports which stages were
   corrected in this run. A signal nothing can read is not collected.
2. `shared/telemetry/event-schema.md`: record that the executor now detects and emits this on its
   own timeline, and that `gate_decision` on `events.jsonl` remains L3.9's. Do not add the type
   there.
3. `shared/skills/deliver-feature/SKILL.md` and `shared/skills/extract-lessons/SKILL.md`: state
   plainly which pipeline produces this signal and which does not, and what `extract-lessons` may
   now mine when a run was executed by `loom run`.
3a. Document the editing contract for a human at a gate: annotating a rendered view is how you say
   "this should have said X", it is recorded and attributed, and it is advisory — the pipeline does
   not adopt it and the next render overwrites the file. Editing a tracked artifact instead re-runs
   its stage (L2.12), which is a different act with a different consequence. `cmd/loom/README.md`
   carries this too.
4. Roadmap: L4.5 **SHIPPED**; update **L4.4** with what it inherits and what it still needs; update
   **L3.5** to record that it adopts this signal rather than inventing one; check L4.5's blockers
   line against reality the way epic 84 had to.
5. `cmd/loom/README.md`, `docs/ARCHITECTURE.md`, `shared/DOMAIN_DICTIONARY.md` if a term appears.
6. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when**: the signal is readable by a human without opening a JSON file, and the docs say
exactly which pipeline produces it. **Commit** (`docs(telemetry): the executor collects the
correction signal`), report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- `events.jsonl`, the `gate_decision` event type, or the event-type enum (**L3.9**)
- The episodic store this signal should eventually live in (**L3.5**)
- Prompt-variant generation or anything that consumes the signal to change an agent (**L4.4**)
- Emitting from the markdown pipeline
- Semantic interpretation of a diff — this epic records what changed, never why

## Report format (end of every phase)

```
## Epic 85 Phase <X> Report
- Roadmap item: L4.5 — the human-correction signal
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
