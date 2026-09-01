# Epic 82 — The Review Loop Runs Under the Executor (L2.17)

Source: session discussion, 2026-08-31, following epic 81. Operationalizes roadmap item **L2.17**
(BUILD-ROADMAP.md, KERNEL). `deliver-feature` steps 18–21 describe an *iteration*: the code-reviewer
returns CHANGES REQUESTED, both artifacts are copied to `.history/`, and the pipeline repeats "until
APPROVED and structurally valid". Nothing executes any of it — the loop condition, the bound, the
backup, and the decision to stop are instructions a model is asked to follow about its own prior
output. `Plan` is a flat list with no notion of a cycle, so `loom run` invokes the developer exactly
once and the reviewer's verdict changes nothing.

This is the largest remaining piece of the pipeline that exists only as prose.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — item **L2.17**; **L2.16** and **L3.2** for the condition
   language this must *not* become; **L4.5**, which will mine what this records.
2. `shared/skills/deliver-feature/SKILL.md` steps 18–21 and the Rollback section — the loop as
   written, including the `.history/` backup and the "until APPROVED and structurally valid" bound
   that has no number.
3. `shared/contracts/review-contract.md` — `## Overall Status` must contain exactly one of
   `APPROVED` or `CHANGES REQUESTED`, "since the orchestrator's CHANGES REQUESTED loop parses this
   literal string". That parse is what this epic replaces.
4. `internal/state/` — `AnalysisState`, `ArchitectureState`, the projection in `projection.go`, and
   the renderers. A new typed artifact follows these conventions exactly.
5. `internal/orchestrator/plan.go`, `executor.go` (the `advance` loop), `router.go` (how epic 81
   added an executor-internal stage without a new control-flow concept), `state.go` (`Sequence`,
   which already distinguishes a re-run from a new step).

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| What repeats | A **declared span**: plan data names a `From` stage, a `To` stage, a named condition, and `MaxIterations`. `developer → code-reviewer` repeats as a unit | It is what the prose actually describes, and it generalises to the Tier B contract-retry loop (agent → validate-artifact) without a second mechanism. A stage that secretly re-invokes another hides the cycle from the executor, and a `goto` reintroduces exactly the positional jumping L2.15 retired |
| The condition | A **named** predicate over the loop's final typed artifact — `review-approved` — resolved in Go from a small registry. No expression language | A condition language whose evaluator is a prompt is the L2.16 defect; a condition language whose evaluator is CEL is L2.16's *solution* and needs the policy work. Named predicates are testable today and do not prejudge that decision |
| Ending on the bound | Exhausting `MaxIterations` with changes still requested **halts at a gate** (`confirm-unresolved-review`) carrying the outstanding findings. A human approves proceeding or stops the run | Exhaustion is not a crash; it is where automation has done what it can and a human has to look. Failing the run discards a complete run's work over a judgement a human may reasonably make differently; proceeding with a warning ships past a reviewer still asking for changes, which is the asymmetric-cost trade this workstream refuses everywhere else |
| Iteration artifacts | Every iteration is kept as its own artifact — `state/<stage>.<n>.json` — each digested and recorded with its iteration number. `StageRecord` gains `Iteration` and a list of the per-iteration artifacts and their digests; the latest also stays in `ArtifactPath`/`ArtifactSHA256` so existing machinery is untouched | "Why did this take four rounds?" is answerable only if round three still exists, and L4.5 will mine exactly this. `.history/` copies are not digest-tracked, so what earlier rounds said would be unverifiable — the property L2.12 exists to establish |
| Typing scope | The **review artifact only**: a `ReviewState` with a typed verdict enum and the findings the developer needs. `implementation-notes` stays markdown | The condition needs to be a field rather than a grep; nothing evaluates a condition over implementation notes, and the reviewer is an agent that can read markdown. Smallest typing that removes prose parsing from the data path |
| Feedback to the next iteration | The findings are projected into the developer's input for the next round, through the existing `ProjectFor` mechanism | The loop is pointless if the developer cannot see what to fix, and a projection is how L2.9 already moves fields between stages. No new channel |
| Gates inside a loop | A gate before the loop (`confirm-design` on `developer`) is approved once and does not re-halt each iteration. Verify this falls out of the existing approval binding rather than building it | An approval binds the artifacts complete when it was given; the developer's own output was not among them, so iterating does not invalidate it. If that turns out to be false, **STOP and report** — do not add a special case |
| Sequence | Iterating does not change a stage's `Sequence`; the iteration number is separate | L2.12 already established that a re-run is the same step of the run, not a new one. A loop is many re-runs |
| State schema | The new fields bump `StateSchemaVersion` to 7; no migration code | Same policy as every prior bump |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**60.1%** as of epic 81). Integration tests driving the built
  binary are invisible to the coverage tool — if the floor is threatened, add in-process tests,
  never lower the floor.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — The typed review verdict — UNBLOCKED

1. `internal/state/review_state.go`: `ReviewState` modelling `review-contract.md` — a typed
   `Verdict` (`APPROVED` | `CHANGES_REQUESTED`), the four Design Score dimensions the contract
   already requires with their 1–5 ratings, and `Findings` (what the developer must address, each
   with enough structure to be actionable). Follow the existing typed-state conventions exactly:
   `Validate()`, generated schema, renderer under the contract's filename
   (`code-review-report.md`), retrieval frontmatter.
2. A projection from `ReviewState` into the developer's next iteration, carrying the findings and
   the verdict — through `ProjectFor`, not a new channel.
3. Tests: verdict round-trips; a report claiming APPROVED while carrying unresolved findings is
   *not* rejected (that judgement belongs to a human, per the contract's own note) but the
   condition reads the verdict field, never the prose; the renderer emits `## Overall Status` with
   the bolded literal the markdown pipeline still greps for.

**Done when**: the review verdict is a typed field, and the rendered view still satisfies the
existing contract. **Commit** (`feat(state): typed review verdict and findings`), report, PAUSE.

## Phase B — The bounded loop in the executor — BLOCKED BY Phase A

1. `internal/orchestrator/loop.go` + `plan.go`: `Loop{From, To, Condition, MaxIterations}` as plan
   data; the built-in plan declares `developer → code-reviewer` with `review-approved` and a bound.
   A named-condition registry resolved in Go. `StageRecord.Iteration` plus the per-iteration
   artifact list with digests; schema 7.
2. The run loop honours a declared span: on a false condition, re-enter at `From` with the
   iteration incremented and the previous artifacts retained. On the bound, halt at
   `confirm-unresolved-review` with the outstanding findings.
3. Verify — do not build — that a gate before the loop does not re-halt on each iteration.
4. Tests: an approving reviewer runs the loop once; a reviewer that requests changes twice then
   approves runs the developer three times, with three retained, digested review artifacts; the
   bound halts at the gate and approving it proceeds; every iteration appears on the timeline; an
   edit to an iteration artifact demotes the stage exactly as L2.12 does for any artifact.

**Done when**: a code-reviewer returning CHANGES REQUESTED causes `loom run` to re-invoke the
developer, and the run stops at a declared bound with every iteration visible in run state.
**Commit** (`feat(orchestrator): bounded developer/code-reviewer loop`), report, PAUSE.

## Phase C — CLI surface and prose — BLOCKED BY Phase B

1. `loom state show` reports the iteration count per looping stage; `loom run` says which round it
   is entering when it re-enters a loop.
2. `shared/skills/deliver-feature/SKILL.md` steps 18–21: state the bound (it currently has none),
   describe the executor's behaviour, and keep the prose procedure for the markdown pipeline.
   `shared/contracts/review-contract.md`: reference the schema and note that the literal-string
   parse is the markdown path.
3. `cmd/loom/README.md`, `README.md`, `docs/ARCHITECTURE.md`, BUILD-ROADMAP SHIPPED marker,
   `shared/DOMAIN_DICTIONARY.md` (**Review Loop**, and the health check will flag the term if the
   prose never uses it by name — it did in epic 81).
4. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when** (L2.17): the roadmap done-when holds — no prose instruction anywhere in the loop's
execution path — and the docs say which pipeline that is true for. **Commit** (`docs(loop): the
review loop is executed, not described — align prose with L2.17`), report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Wiring the Tier B contract-retry loop (agent → validate-artifact) onto this mechanism. The
  mechanism must generalise to it; using it there is a later item
- A condition expression language, CEL, or policy-driven bounds (**L2.16**)
- Typing `implementation-notes` or any other artifact (the L2.9 continuation)
- Rollback, `.history/` restoration, or `--from-phase` (**L2.15**)
- Mining the retained iterations for correction signal (**L4.5**)
- Making the markdown pipeline call `loom` for its loop

## Report format (end of every phase)

```
## Epic 82 Phase <X> Report
- Roadmap item: L2.17 — The review loop under the executor
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
