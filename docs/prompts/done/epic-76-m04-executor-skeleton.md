# Epic 76 — Executor Skeleton: `loom run` Owns the Run Loop (M0.1 + M0.2 remainder + M0.4)

Source: epic 75 Phase E blocker analysis, 2026-08-29. Operationalizes roadmap items **M0.1**,
the unfinished half of **M0.2**, and **M0.4** (BUILD-ROADMAP.md, Milestone 0). M0.4 is the
highest-leverage single item on the roadmap: with L2.13 it unblocks roughly two-thirds of
everything else (L2.9, L2.12, L2.13, L2.14, L3.1, L3.3, L3.8, L4.1 — and through those, the
remaining three epic-75 Phase E MCP tools).

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — items M0.1, M0.2, M0.4, and (for design alignment only,
   NOT implementation scope) L2.12–L2.15. Respect Blocked-by: M0.4 needs M0.1 and M0.2 done.
2. `docs/roadmaps/architectural-audit-2026-08-29.md` — audit findings H9 (CI gaps) and the
   kernel-absence analysis.
3. `docs/prompts/epic-75-distribution-adoption.md` — the fixed strategy sentence that constrains
   this epic: **the orchestration kernel is `loom run` acting as an MCP client — the pipeline
   itself is never a tool someone else calls.**
4. Current blocker evidence (verified 2026-08-29): no kernel ADR exists (`docs/adrs/` tops out at
   ADR-005, install-script paths); `framework-ci.yml` has `go build`/`go test` and the embedding
   example but **no golangci-lint job, no SHA-pinned actions, no coverage ratchet**, and its own
   comment defers workflow hardening to M0.2; `.golangci.yml` exists at the root (gocyclo 7).

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| What loom is | loom **executes** pipelines (ADR to be accepted in Phase A) | Every Milestone-1+ item assumes it; drafting the ADR makes the drift a decision |
| Executor scope | Minimal linear run loop: load plan → execute stages in order → persist state → stop. No routing, no parallelism, no policy, no gates | Those are L2.13/L3.x items that plug in later; skeleton first |
| Default plan | The existing linear `deliver-feature` agent sequence, hardcoded as the built-in plan | Behavior preserved while the substrate changes underneath (M0.4 fix text) |
| Stage identity | Stable string IDs (agent names), never ordinals | L2.15 explicitly calls hand-numbered step lists a defect; don't build a new one |
| State file | Executor-owned, atomic write (temp + `os.Rename`), schema version field, real SHA-256 computed in Go | Aligns with L2.12's shape so that item becomes "move ownership", not "rewrite" — but full L2.12 (skill rewrite, tamper refusal semantics) stays out of scope |
| Provider abstraction | `Provider` interface defined in `internal/orchestrator` (the consumer), implemented in `internal/provider` | go-conventions.md: interfaces at the consumer |
| First providers | (1) `claude` — spawns `claude -p` headless with the stage's agent prompt; (2) `mock` — deterministic canned responses for tests | One real, one testable; anything fancier (API-direct, other hosts) is a later item |
| Provider timeouts | Explicit per-stage timeout from the plan, enforced with `context.WithTimeout` | architecture-guardrails.md #5; no unbounded subprocess waits |
| Resume | `loom run --resume` re-loads persisted state and continues from the first stage not marked complete | M0.4's done-when; full L2.15 semantics (rollback, --from-phase) out of scope |
| MCP relationship | The executor may *call* MCP tools (client role) in later items; it never *exposes* the pipeline as an MCP tool | Epic-75 fixed strategy |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, errors handled explicitly, interfaces at
  the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, plus
  `scripts/health-check.sh` if `shared/` content changed.
- Update docs in the same phase that changes behavior (`README.md`, `cmd/loom/README.md`,
  `docs/ARCHITECTURE.md`).
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — ADR: loom executes pipelines (roadmap M0.1) — UNBLOCKED, ENDS IN A HUMAN GATE

1. Write `docs/adrs/ADR-006-loom-executes-pipelines.md` (follow the existing ADR format): the
   question ("does loom execute pipelines, or only validate and distribute content a host
   executes?"), the decision (**execute** — a minimal Go executor owning run loop, state, and
   eventually gates), the alternative seriously stated (validate-and-distribute is a legitimate
   product; name what gets deleted from the roadmap if chosen), and consequences (Milestones 1–4
   become real; host platforms become one provider among several, not the runtime).
2. Audit `README.md` and `docs/ARCHITECTURE.md` for present-tense descriptions of unimplemented
   subsystems (orchestration, telemetry, policy evaluation, retrieval tiers). Rewrite those to
   honest status markers ("specified, ships with LX.Y") — do not delete the descriptions.
3. Status the ADR **Proposed**, not Accepted. Accepting it is the human's call.

**Done when** (M0.1): the Proposed ADR exists and README/ARCHITECTURE no longer claim
unimplemented subsystems in the present tense. **Commit**
(`docs(adrs): propose ADR-006 — loom executes pipelines`), report, **PAUSE — do not start
Phase B until the human accepts (or redirects) the ADR. If the human chooses
validate-and-distribute, this epic ends here.**

## Phase B — Close the M0.2 remainder (roadmap M0.2) — BLOCKED BY Phase A acceptance

Verify first what already exists (go build/test job and the embedding-example build are in CI;
`.golangci.yml` exists) and implement only the delta:

1. Add a `golangci-lint` job to `framework-ci.yml` (pin the linter version; v2 config, go1.27 —
   see the repo's established lint setup).
2. Harden the workflow per `iac-conventions.md`: explicit `permissions:` block on every job,
   every action pinned to a full commit SHA (mirror what `loom-release.yml` already does).
3. Add a coverage step with a ratchet: measure today's real combined coverage for the Go
   packages, fail CI below that number, record the number and the ratchet policy in the workflow
   comment. Do NOT pretend 85% — start at the honest measured value.
4. Fix `scripts/test-agents.sh` so a missing fixture is a FAIL, not a SKIP that exits 0.
5. Prove the gate: in a scratch branch or locally, verify CI-equivalent commands fail on (a) a
   deliberate compile error, (b) a deliberate complexity-9 function, (c) a deleted test fixture.
   Record the evidence in the report; do not commit the deliberate breakage.

**Done when** (M0.2): the three deliberate-breakage checks each fail the gate. **Commit**
(`chore(ci): lint job, SHA pinning, coverage ratchet, fixture-fail — close M0.2`), report, PAUSE.

## Phase C — Executor core (roadmap M0.4, part 1) — BLOCKED BY Phases A+B

1. `internal/orchestrator/`:
   - `plan.go` — `Plan{Name string; Stages []Stage}`, `Stage{ID string; Agent string; Timeout
     time.Duration}`; a `DefaultDeliverFeaturePlan()` constructor encoding the current linear
     agent sequence from `shared/skills/deliver-feature/SKILL.md` (IDs = agent names).
   - `executor.go` — the run loop: for each stage not already complete in state, invoke the
     provider with the stage context, persist state after every transition (stage started,
     completed, failed), stop on first failure. Honors `context.Context` cancellation (SIGINT)
     by persisting a clean checkpoint before exit.
   - `state.go` — run state: schema version, plan name, per-stage status + timestamps + SHA-256
     of the stage's output artifact (computed in Go). Atomic persistence: write temp file in the
     same directory, `os.Rename`. Location: `.claude/feature-workspace/<feature>/run-state.json`
     (note in the code comment: L2.12 will migrate `pipeline-state.json` semantics here).
   - `provider.go` — the consumer-side `Provider` interface: something shaped like
     `Invoke(ctx context.Context, stage Stage, input StageInput) (StageOutput, error)`. Keep it
     small; resist adding methods for later items.
2. `internal/provider/mock/` — deterministic mock provider for tests (scripted outputs,
   scripted failures, scripted hangs to exercise timeout).
3. Table-driven tests: happy path over ≥3 stages; failure stops the run and persists FAILED;
   timeout fires; SIGINT (context cancel) mid-stage persists a resumable checkpoint; resume skips
   completed stages and re-runs the interrupted one; state file survives a crash between temp
   write and rename (i.e., rename atomicity is actually relied on).
4. No `cmd/` wiring yet — this phase is the library, fully tested against the mock provider.

**Done when**: executor runs the default plan against the mock provider with resume-after-cancel
proven in tests. **Commit** (`feat(orchestrator): executor skeleton — plan, run loop, durable
state`), report, PAUSE.

## Phase D — `loom run` + claude provider (roadmap M0.4, part 2) — BLOCKED BY Phase C

1. `internal/provider/claude/` — invokes `claude -p` as a subprocess: builds the stage prompt
   from the agent's `shared/agents/<agent>.md` definition plus the feature spec path, explicit
   `exec.CommandContext` timeout, captures output to the stage's artifact path. Structured JSON
   logs to stderr. If the `claude` binary is absent, fail the stage with a clear remediation
   message — never silently fall back to the mock.
2. `cmd/loom/cmd/run.go` — `loom run --spec docs/features/<x>/spec.md [--resume]
   [--provider claude|mock] [--plan deliver-feature]`. Default provider `claude`, default plan
   the built-in one. `--resume` requires existing state and refuses to start fresh over it.
3. SIGINT handling at the command level: first Ctrl-C cancels the context (graceful checkpoint),
   second kills.
4. Docs: `cmd/loom/README.md` gains a "Running pipelines" section — explicitly labeled
   experimental/skeleton, listing what it does NOT yet do (no gates, no retries, no parallelism,
   no policy — pointers to L2.13/L3.x).
5. Integration test with the mock provider through the real CLI (build the binary, run it, kill
   it mid-stage, resume it — same pattern as `mcp_serve_test.go`). A claude-provider test may be
   skipped when the binary is absent, but the subprocess plumbing (arg construction, timeout,
   missing-binary error) must be unit-tested.

**Done when** (M0.4): `loom run --spec features/<x>.md` executes at least three real stages
end-to-end, writes state, and resumes correctly after SIGINT — demonstrated with the mock
provider in CI and manually with the claude provider (record the manual run's evidence in the
report). **Commit** (`feat(loom): loom run — execute the delivery plan with durable state`),
report, PAUSE.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Gates as process interrupts (L2.13) and approval tokens (L2.14)
- `pipeline-state.json` migration and tamper-refusal semantics (L2.12)
- `--from-phase` / rollback resume modes (L2.15)
- Typed inter-stage state and projections (L2.9)
- Telemetry emission (L3.8), retries/backoff, parallel branches, policy evaluation
- Any new MCP tool (that's epic 75 Phase E, which this epic unblocks)

## Report format (end of every phase)

```
## Epic 76 Phase <X> Report
- Roadmap item: <M0.n> — <title>
- Blockers verified: <list, with evidence (commit SHAs / files)>
- Commits: <sha> <subject>
- Build/lint/test: go build PASS|FAIL · golangci-lint PASS|FAIL · go test PASS|FAIL (counts)
- Done-when criterion: <restate it> — MET | NOT MET (why)
- Escalations / open questions: <list or "none">
- Next phase blocked by: <what must land first>
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
