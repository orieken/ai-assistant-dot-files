# Epic 78 — The Executor Owns Pipeline State: Integrity in Go, Not in Prose (L2.12)

Source: epic 77 completion, 2026-08-29. Operationalizes roadmap item **L2.12**
(BUILD-ROADMAP.md, KERNEL — Checkpointing). With M0.4 (epic 76) and L2.13 (epic 77, commits
`868c281`, `53880b6`, `a971b2f`) landed, the executor runs stages and halts at gates. This epic
makes its state file trustworthy: artifact digests are computed *and verified* in Go, so a stage
whose artifact changed on disk is not treated as complete. L2.12 unblocks **L2.14** (approvals
bound to digests — "any edit resets the gate") and **L2.15** (resume as a real capability).

Phase D extends the same guarantee to the markdown pipeline without making it delegate execution:
the model keeps routing, but records every checkpoint through `loom state` so no model ever
computes its own integrity hash again.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — item L2.12 (the spec), plus L2.14 and L2.15 for the boundary
   of what NOT to build. Respect Blocked-by: L2.12 needs only M0.4 (done).
2. `docs/adrs/ADR-006-loom-executes-pipelines.md` — Accepted. State ownership moving below the
   model is a named consequence.
3. `internal/orchestrator/` — `state.go` (schema-versioned `run-state.json`, atomic writes,
   `ArtifactSHA256`), `executor.go` (`prepareState`, run loop, `IsStageCompleted` as the resume
   predicate), `gate.go` (approvals, `WAITING_APPROVAL`). Read `executor_test.go`, `gate_test.go`,
   and `cmd/loom/cmd/run_test.go` for the established test patterns (mock provider harness;
   build-the-binary integration tests).
4. `shared/skills/deliver-feature/SKILL.md` — "Checkpointing & Pipeline State" (the prose that
   asks the *model* to compute `sha256` and to distrust a mismatched artifact) and "Rollback".
5. `shared/skills/resume-pipeline/SKILL.md` — Mode 1 recomputes checksums by instruction; Modes 2
   and 3 (`--from-phase`, rollback) are the parts this epic deliberately does not implement.
6. Audit finding context (`docs/roadmaps/architectural-audit-2026-08-29.md`): the failure mode
   being closed is "a model computing and recording its own integrity hashes is not integrity."

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| One state file | `run-state.json` is the single durable record for both pipelines. The executor writes it during `loom run`; the markdown pipeline writes it **through `loom state` subcommands** (Phase D) rather than by hand. `pipeline-state.json` is not read, written, or migrated by any Go code — it remains the fallback the prose uses when the `loom` binary is absent | Two hand-written owners for one fact is the defect. Splitting *who routes* (the model, for now) from *who owns integrity* (Go) closes the L2.12 hole for the markdown pipeline without first needing conditional routing (L3.1), contract validation (L2.11), or policy (L2.16) |
| Recorded-by provenance | `RunState` carries `CreatedBy` (`executor` \| `markdown`). `loom run` refuses to resume a state file created by the markdown pipeline, and vice versa | The two pipelines route differently — the markdown one skips conditional agents, the executor's linear plan does not. A shared file is fine; a shared *resume* would silently mis-order work |
| Stage sequence | Every stage record carries `Sequence` — a monotonic counter assigned when the record is first created and **preserved across re-runs of that stage**. Verification cascades by sequence, not plan order, so both writers share one rule | The markdown pipeline has no plan to order by: it routes conditionally, skips agents, and loops developer↔code-reviewer. Recording order is the only ordering that exists for it. For the executor's linear plan, sequence order equals plan order, so this is behavior-preserving. It also answers "where did this run get stuck" directly from the state file |
| `CreatedBy` lands in Phase D | The field is introduced with the Phase D schema bump (v3 → v4) together with `Sequence`, not retrofitted into Phase A's v3 | One bump for the whole markdown-pipeline surface beats two. Phase A's v3 is already released to nobody but this repo's own runs |
| `loom state approve` | Phase D ships `loom state approve --spec <x> --gate <name>`, recording an approval into the same `approvals` map `loom run` uses. It **records** the approval; it does not enforce the gate — the markdown pipeline can still proceed without calling it | The prose gates currently leave no durable trace at all, so a recorded approval is strictly better than nothing for audit. But the epic must say plainly that this is bookkeeping, not enforcement: enforcement needs the action itself running under the executor |
| `loom state` failure mode | The markdown pipeline is instructed to call `loom state`, and can still decline to — this is prompt-discipline about *recording*, not about integrity. But a skipped call leaves a **missing** record, where hand-written JSON left a plausible-but-wrong hash | Missing state is visible and fails loudly on the next verify; a forged digest is invisible. Degrading the failure mode is the whole point of the middle path |
| Verification point | `Executor.prepareState` verifies every `COMPLETED` stage record before the run loop starts: recompute the artifact's SHA-256 and compare to the recorded one. Missing file, unreadable file, or digest mismatch all count as mismatch | Verification must happen where the resume decision is made. Doing it lazily per stage would let a later stage consume an artifact never checked |
| What a mismatch does | The mismatched stage is demoted to `StageStatusStale` (not completed → it re-runs), **and every stage completed after it in plan order is demoted too**, because they consumed content that no longer exists | This is the `deliver-feature` Rollback rule ("every agent after it") expressed as code. Demoting only the edited stage would leave downstream artifacts derived from vanished input silently trusted |
| Stale is a status, not a deletion | `StageStatusStale` records `previousStatus`, the recorded digest, and the digest found, so the CLI can report exactly what changed. Records are never dropped | The state file is also the audit trail; deleting evidence of an edit defeats the point of detecting it |
| Stages with no artifact | A `COMPLETED` record with an empty `ArtifactPath` verifies trivially (nothing to compare) and stays completed | The executor cannot verify what a stage never produced; inventing a failure there would make every artifact-less stage unresumable |
| Gate interaction | Demoting stages does **not** revoke any approval. An approved gate stays approved even when the stage behind it goes stale | "Any edit resets the gate" is **L2.14** and needs approval-to-digest binding, which this epic deliberately does not build. Doing half of it here would ship a rule nobody can rely on |
| Run identity | `RunState` gains `FeatureName`, `SpecPath`, and `StartedAt`, set once at run creation and verified on load: a state file whose `SpecPath` differs from the current invocation is refused, exactly as a mismatched `PlanName` already is | These are the only `pipeline-state.json` fields with a real consumer today. Step ordinals (`currentPhase`, `lastCompletedStep`) are deliberately NOT adopted — stable stage IDs replace them (L2.15) |
| Contract/skip metadata | `contractStatus`, `contractRetries`, and `SKIPPED` entries are **not** adopted | They belong to conditional routing and contract validation, which are L3.1 and L2.11. Adding fields the executor cannot populate would be state theatre |
| State schema | Adding `StageStatusStale`, the stale detail fields, and the run identity fields bumps `StateSchemaVersion` to 3; no migration code — a v1/v2 file is refused with the existing clear message | Same reasoning as epic 77: runs are hours old, migration machinery is not yet earning its keep |
| CLI surface | No new flags. `loom run --resume` reports invalidated stages on stderr before the run loop starts (`stage %q was COMPLETED but its artifact changed — re-running (and N later stage(s))`). A fresh run is unaffected | Verification is not opt-in; a flag to skip it would recreate the hole this epic closes |
| Prose relationship | `deliver-feature/SKILL.md` and `resume-pipeline/SKILL.md` instruct: call `loom state` when the binary is present; fall back to the existing hand-written `pipeline-state.json` procedure when it is not. Do not delete the prose procedure | The framework must keep working for installs without the Go binary (Cursor, Windsurf, and every non-Claude-Code target). A hard `loom` dependency in the flagship workflow is a portability regression, not an upgrade |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` when `shared/` content changed. Coverage must stay ≥ the CI ratchet
  floor (**52.4%** as of epic 77; raise the floor in `framework-ci.yml` if measured coverage rises).
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — Digest verification in the library — UNBLOCKED

1. `internal/orchestrator/integrity.go` (+ edits to `state.go`, `executor.go`):
   - `StageStatusStale` + the stale detail fields on `StageRecord` (`previousStatus`,
     `recordedSha256`/`foundSha256` or equivalent — keep the names boring and the shape flat).
   - `RunState.FeatureName`, `SpecPath`, `StartedAt`; schema version 3.
   - A verification pass invoked from `prepareState`: walk the plan in order, recompute digests for
     `COMPLETED` records that have an artifact path, demote on mismatch, then cascade the demotion
     to every later completed stage. Return the list of demoted stage IDs so the caller can report
     them; persist once, atomically.
   - `SpecPath` mismatch refused like `PlanName` mismatch already is.
2. Tests (extend the existing harness):
   - Hand-edit an artifact between two runs: that stage re-runs, and the provider invocation list
     proves it. **This is the L2.12 done-when test — name it accordingly.**
   - Cascade: editing an early artifact demotes the later completed stages too, and they re-run.
   - Deleting an artifact counts as a mismatch.
   - A stage with no artifact stays completed.
   - Editing an artifact does **not** remove an approval (the L2.14 boundary, asserted explicitly
     so the next epic starts from a known state).
   - A v2 state file is refused with the schema message.
   - Unchanged artifacts do not demote anything — no false positives on a clean resume.

**Done when**: hand-editing an artifact causes the executor to refuse to treat that stage as
complete, in Go, with no prompt involved. **Commit**
(`feat(orchestrator): verify artifact digests on resume — stale stages re-run`), report, PAUSE.

## Phase B — CLI reporting + run identity — BLOCKED BY Phase A

1. `cmd/loom/cmd/run.go` (+ a small helper file if it keeps functions under the complexity cap):
   populate `FeatureName`/`SpecPath`/`StartedAt` on fresh runs; print the invalidated-stage report
   on stderr before the run loop; keep exit codes as they are (0 success, 3 waiting on a gate).
2. Integration tests through the real binary (extending `run_test.go`): run to the first gate,
   hand-edit `analyst.md` in the workspace, resume with `--approve confirm-design` — the analyst
   stage re-runs, the report names it, and the run still halts at the next gate. A run resumed
   against a different `--spec` is refused.
3. Verify the gate interplay explicitly: a stale stage behind an approved gate does not re-prompt
   for approval (L2.14 is not being implemented here — prove the current behavior, don't change it).

**Done when**: the real binary detects an out-of-band edit, says which stages it invalidated, and
re-runs exactly those. **Commit** (`feat(loom): report and re-run stages invalidated by edits`),
report, PAUSE.

## Phase C — Prose alignment + docs — BLOCKED BY Phase B

1. `shared/skills/deliver-feature/SKILL.md` ("Checkpointing & Pipeline State") and
   `shared/skills/resume-pipeline/SKILL.md` (Mode 1): add the honest scope note per the design
   table — under `loom run`, `run-state.json` is executor-owned and digests are verified in Go;
   the prose procedure governs the markdown pipeline until Phase D wires it to `loom state`. Do
   not weaken or delete the prose.
2. `cmd/loom/README.md` "Running pipelines": document digest verification on resume and what a
   stale stage means; keep the NOT-yet list accurate (L2.14 gate-reset, L2.15 `--from-phase` and
   rollback, L3.1 routing, L3.8 telemetry).
3. `docs/ARCHITECTURE.md` §4 ("The Pipeline's Own Observability") must stay honest about which
   numbers are measured and which are estimated — Phase E changes that answer, so leave §4's
   trace/timing prose for Phase E rather than half-updating it here.
4. `README.md` + `docs/ARCHITECTURE.md`: update the L2.12 status markers — `pipeline-state.json`
   ownership has moved for `loom run`; the markdown pipeline's copy is still prompt-discipline.
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` (shared/ changed) plus the full Go
   gate.

**Done when** (L2.12): the roadmap item's done-when holds end to end and the docs say exactly which
pipeline it holds for. **Commit** (`docs(state): executor-owned run state — align prose with
L2.12`), report, PAUSE — epic complete.

## Phase D — `loom state` for the markdown pipeline — BLOCKED BY Phase C

The middle path: the model keeps deciding what runs next; Go computes and owns every digest.

1. `cmd/loom/cmd/state.go` (+ small helper files to stay under the complexity cap) — three
   subcommands, all resolving the workspace exactly as `loom run` does (`featureSlug`,
   `.claude/feature-workspace/<feature>/`):
   - `loom state record --spec <spec> --stage <id> --artifact <file>` — creates the state file on
     first call (`CreatedBy: markdown`), computes the artifact's SHA-256 **in Go**, persists the
     stage as COMPLETED atomically. Recording a stage twice overwrites its record (a re-run after
     a CHANGES REQUESTED loop is normal, not an error).
   - `loom state verify --spec <spec>` — re-verifies every recorded digest, prints per-stage OK /
     MISMATCH / MISSING, exits non-zero if anything failed. This is what `resume-pipeline` Mode 1
     calls instead of asking the model to recompute checksums.
   - `loom state approve --spec <spec> --gate <name>` — records a gate approval into the same
     `approvals` map `loom run` writes (method `cli`, approver = OS username). Output and docs must
     say plainly that this **records** an approval for audit and does not enforce the gate.
   - `loom state show --spec <spec>` — prints the recorded stages in sequence order with their
     statuses, the recorded approvals, and which stage the run is sitting on (human-readable by
     default, `--json` for machine consumption). This is the "where did it get stuck" view.
   - No `--sha` or equivalent flag anywhere: a caller must never be able to supply a digest. The
     binary reads the file and hashes it, or it fails.
2. Schema v4: add `RunState.CreatedBy` (`executor` | `markdown`) and `StageRecord.Sequence`
   (monotonic, assigned on first creation, preserved when a stage is re-recorded). Switch
   `verifyCompletedStages` from plan order to sequence order — behavior-preserving for the
   executor's linear plan, and the only ordering the markdown pipeline has. `loom run` refuses a
   `markdown` state file with a message pointing at the difference in routing; `loom state record`
   refuses an `executor` one.
3. Prose wiring: `deliver-feature/SKILL.md` "Checkpointing & Pipeline State" and
   `resume-pipeline/SKILL.md` Mode 1 instruct the `loom state` calls, with the existing
   hand-written procedure kept explicitly as the no-binary fallback.
4. Tests: subcommand-level tests through the real binary (record → verify OK → hand-edit →
   verify MISMATCH exit non-zero); cascade by sequence demotes later-recorded stages; sequence
   survives re-recording a stage (the CHANGES REQUESTED loop); the cross-pipeline refusals in both
   directions; `approve` records into the same map `loom run` reads; `--json` output parses.

**Done when**: a markdown-pipeline run records its checkpoints through `loom state`, and a
hand-edited artifact is caught by `loom state verify` — with no model computing a hash anywhere in
the path. **Commit** (`feat(loom): loom state — markdown pipeline records checkpoints in Go`),
report, PAUSE.

## Phase E — An owned event timeline — BLOCKED BY Phase D

State says *where a run is*; nothing says *what happened and when*, with timestamps anyone can
trust. `pipeline-trace.json` is written by the model, so its durations are estimates — the same
defect L2.12 closes for digests, unclosed for time.

1. `internal/orchestrator/timeline.go` — an append-only `run-events.jsonl` beside
   `run-state.json`, owned by Go and written by **both** `loom run` and `loom state`:
   - One JSON object per line, appended with `O_APPEND` single-line writes (append is the atomicity
     model here — never read-modify-write the whole file).
   - Event kinds, kept deliberately small: `stage.started`, `stage.completed`, `stage.failed`,
     `stage.interrupted`, `stage.stale`, `gate.waiting`, `gate.approved`, `run.started`,
     `run.completed`. Each carries `at` (UTC, from the clock, never from a model), the stage or
     gate name, the sequence, and the fields that kind needs (`staleReason`, `approvalMethod`, …).
   - Timestamps come from the process doing the work, so stage durations become derivable by
     subtraction — real measurement, not recall.
2. `loom state timeline --spec <spec>` (`--json`) renders the file as an ordered audit trail.
3. Honest boundary in the docs: this is a local audit log, **not** the OpenTelemetry emission of
   L3.8 and not a replacement for `pipeline-trace.json` (which keeps its model-written
   `budgetUtilization` and iteration counts until L3.8 lands). Say which numbers are measured and
   which are still estimated.
4. Tests: every executor transition appends exactly one event; an interrupted run's timeline is
   readable and ordered; `loom state` writes to the same file in the same shape; a corrupt trailing
   line does not break reading the events before it.

**Done when**: a completed run's timeline shows every stage transition and gate decision with
timestamps produced by the process, and stage durations can be computed from it without asking a
model. **Commit** (`feat(orchestrator): owned run event timeline for auditing`), report, PAUSE —
epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Binding approvals to artifact digests / "any edit resets the gate" (**L2.14** — the next epic)
- `--from-phase N`, rollback, and `.history/` restoration as executor operations (**L2.15**)
- Reading, writing, or migrating `pipeline-state.json` or `pipeline-trace.json` in Go (the prose
  keeps them as the no-binary fallback; no Go code touches either file)
- Making `deliver-feature` delegate execution to `loom run` (routing stays with the model until
  L3.1, L2.11, L2.16, L2.15, and L3.8 land — Phase D moves *state*, not orchestration)
- Contract-validation status, retry counts, or SKIPPED records in run state (L2.11, L3.1)
- OpenTelemetry emission, spans, or the `shared/telemetry/event-schema.md` wire format
  (**L3.8/L3.9**) — Phase E's timeline is a local audit log, not telemetry
- Making `loom state approve` *enforce* a prose gate (enforcement needs the gated action itself
  running under the executor)
- Migration code for v1/v2 state files

## Report format (end of every phase)

```
## Epic 78 Phase <X> Report
- Roadmap item: L2.12 — Executor owns pipeline state
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
