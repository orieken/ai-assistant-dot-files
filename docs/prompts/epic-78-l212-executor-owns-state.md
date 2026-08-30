# Epic 78 — The Executor Owns Pipeline State: Integrity in Go, Not in Prose (L2.12)

Source: epic 77 completion, 2026-08-29. Operationalizes roadmap item **L2.12**
(BUILD-ROADMAP.md, KERNEL — Checkpointing). With M0.4 (epic 76) and L2.13 (epic 77, commits
`868c281`, `53880b6`, `a971b2f`) landed, the executor runs stages and halts at gates. This epic
makes its state file trustworthy: artifact digests are computed *and verified* in Go, so a stage
whose artifact changed on disk is not treated as complete. L2.12 unblocks **L2.14** (approvals
bound to digests — "any edit resets the gate") and **L2.15** (resume as a real capability).

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
| One state file | `run-state.json` is the executor's single durable record. `pipeline-state.json` is **not** written, read, or migrated by the executor — it stays the markdown pipeline's prompt-owned file for runs that never touch `loom run` | Two owners for one fact is the defect. The executor owning its own file, and the prose owning its own, is honest; a shared file with two writers is not |
| Verification point | `Executor.prepareState` verifies every `COMPLETED` stage record before the run loop starts: recompute the artifact's SHA-256 and compare to the recorded one. Missing file, unreadable file, or digest mismatch all count as mismatch | Verification must happen where the resume decision is made. Doing it lazily per stage would let a later stage consume an artifact never checked |
| What a mismatch does | The mismatched stage is demoted to `StageStatusStale` (not completed → it re-runs), **and every stage completed after it in plan order is demoted too**, because they consumed content that no longer exists | This is the `deliver-feature` Rollback rule ("every agent after it") expressed as code. Demoting only the edited stage would leave downstream artifacts derived from vanished input silently trusted |
| Stale is a status, not a deletion | `StageStatusStale` records `previousStatus`, the recorded digest, and the digest found, so the CLI can report exactly what changed. Records are never dropped | The state file is also the audit trail; deleting evidence of an edit defeats the point of detecting it |
| Stages with no artifact | A `COMPLETED` record with an empty `ArtifactPath` verifies trivially (nothing to compare) and stays completed | The executor cannot verify what a stage never produced; inventing a failure there would make every artifact-less stage unresumable |
| Gate interaction | Demoting stages does **not** revoke any approval. An approved gate stays approved even when the stage behind it goes stale | "Any edit resets the gate" is **L2.14** and needs approval-to-digest binding, which this epic deliberately does not build. Doing half of it here would ship a rule nobody can rely on |
| Run identity | `RunState` gains `FeatureName`, `SpecPath`, and `StartedAt`, set once at run creation and verified on load: a state file whose `SpecPath` differs from the current invocation is refused, exactly as a mismatched `PlanName` already is | These are the only `pipeline-state.json` fields with a real consumer today. Step ordinals (`currentPhase`, `lastCompletedStep`) are deliberately NOT adopted — stable stage IDs replace them (L2.15) |
| Contract/skip metadata | `contractStatus`, `contractRetries`, and `SKIPPED` entries are **not** adopted | They belong to conditional routing and contract validation, which are L3.1 and L2.11. Adding fields the executor cannot populate would be state theatre |
| State schema | Adding `StageStatusStale`, the stale detail fields, and the run identity fields bumps `StateSchemaVersion` to 3; no migration code — a v1/v2 file is refused with the existing clear message | Same reasoning as epic 77: runs are hours old, migration machinery is not yet earning its keep |
| CLI surface | No new flags. `loom run --resume` reports invalidated stages on stderr before the run loop starts (`stage %q was COMPLETED but its artifact changed — re-running (and N later stage(s))`). A fresh run is unaffected | Verification is not opt-in; a flag to skip it would recreate the hole this epic closes |
| Prose relationship | `deliver-feature/SKILL.md` and `resume-pipeline/SKILL.md` gain a short scope note: for `loom run`, the executor owns `run-state.json` and verifies digests in code; the prose checkpointing/rollback procedure remains authoritative for the markdown pipeline. Do not delete the prose | The markdown pipeline is still how most runs happen. Deleting its checkpoint instructions would leave it with no state at all |

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
   the prose procedure governs the markdown pipeline only. Do not weaken or delete the prose.
2. `cmd/loom/README.md` "Running pipelines": document digest verification on resume and what a
   stale stage means; keep the NOT-yet list accurate (L2.14 gate-reset, L2.15 `--from-phase` and
   rollback, L3.1 routing, L3.8 telemetry).
3. `README.md` + `docs/ARCHITECTURE.md` §4: update the L2.12 status markers — `pipeline-state.json`
   ownership has moved for `loom run`; the markdown pipeline's copy is still prompt-discipline.
4. Run `scripts/health-check.sh` and `scripts/check-parity.sh` (shared/ changed) plus the full Go
   gate.

**Done when** (L2.12): the roadmap item's done-when holds end to end and the docs say exactly which
pipeline it holds for. **Commit** (`docs(state): executor-owned run state — align prose with
L2.12`), report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Binding approvals to artifact digests / "any edit resets the gate" (**L2.14** — the next epic)
- `--from-phase N`, rollback, and `.history/` restoration as executor operations (**L2.15**)
- Reading, writing, or migrating `pipeline-state.json` or `pipeline-trace.json`
- Contract-validation status, retry counts, or SKIPPED records in run state (L2.11, L3.1)
- Telemetry events for integrity failures (**L3.8/L3.9**)
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
