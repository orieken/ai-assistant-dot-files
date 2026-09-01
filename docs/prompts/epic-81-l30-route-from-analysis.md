# Epic 81 — Compute the Route Before the Design Gate (L3.0)

Source: session discussion, 2026-08-31. Operationalizes roadmap item **L3.0** (BUILD-ROADMAP.md,
KERNEL — Dynamic Routing), a new item scoped out of L3.1. `loom run` executes all fourteen stages
unconditionally; the markdown pipeline's six conditional stages are prose an LLM evaluates about an
artifact it just read, leaving nothing durable that says why a stage was skipped. Neither pipeline
can answer "is devops running on this feature, and if not, why not?" before it gets there.

**The shape.** A fixed prologue runs, then the route is computed once from typed analysis and
recorded as an artifact — before the design gate, so the human approves *what will run* along with
the design, and L2.14 makes editing the route reset that approval.

```
context-engineer → analyst → [ROUTE] → confirm-design → …the routed stages…
```

**Why not before delivery starts.** Almost every condition depends on the analysis: crossings and
migrations for the architect, latency SLAs for performance, schema changes for data, UI surface for
accessibility, infra tasks for devops. A route computed from the spec alone would be guessing at
most of its own decisions. After the analyst, they are all typed facts.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — item **L3.0** (the spec) including its "Known limit"; **L3.1**
   and **L3.2** for what this deliberately does not become; **L2.17** for the cycle problem.
2. `internal/state/analysis_state.go` — `RequiresArchitect()` and its four derived signals plus the
   `ArchitecturalFlags` override, with table-driven tests. This is the pattern every other predicate
   follows; do not invent a second style.
3. `internal/orchestrator/plan.go` — `Stage` (`Gate`, `StateKind`, `Consumes`) and
   `DefaultDeliverFeaturePlan`. Routing joins these as plan data, not as a new mechanism.
4. `internal/orchestrator/approval_binding.go` + `integrity.go` — why a completed artifact before a
   gate is automatically bound to that gate's approval. The route inherits this; it is the reason
   the re-plan point sits before `confirm-design`.
5. `shared/skills/deliver-feature/SKILL.md` steps 12–29 — the six conditional stages in prose, the
   conditions this epic replaces with predicates.
6. `shared/contracts/analysis-contract.md` — what the analyst actually promises, i.e. the only facts
   a predicate may read.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Where the route is computed | A **re-plan point** after `analyst` completes and before the `developer` stage's gate. The prologue (`context-engineer`, `analyst`) is never routed | Every useful condition is analyst output. Routing earlier means guessing; routing later means the human approves a design without knowing what it will cause to run |
| What the route is | A typed artifact (`internal/state/route.go`), one row per stage: stage ID, included or skipped, and the reason. Rendered to markdown like any other typed state | An artifact gets digest recording, staleness, and gate binding for free. A decision that lives only in memory cannot be approved, audited, or edited |
| Gate binding | The route completes before `confirm-design`, so L2.14 binds it automatically. Editing the route invalidates that approval and the run halts until re-approved | This is the point of the whole shape: forcing a stage back in becomes a supported, attributed, gate-bound act rather than a workaround. **Write no new code for this** — verify it falls out, with a test |
| Skippability | An allow-list in plan data (`Stage.Skippable`). `code-reviewer` and `security-reviewer` are **never** auto-skippable; a human may still skip them by editing the route, which the gate makes them re-approve | The costs are asymmetric: running devops unnecessarily wastes an invocation, skipping a security review on an auth change does not. Automation gets the cheap half of that trade |
| Predicates | Methods on `AnalysisState`, in the shape of `RequiresArchitect()`: derived signals OR'd with an explicit override, each with table-driven tests including the false-negative case | One style, one place, already tested. A predicate that cannot be written from the analysis contract's fields is a signal the contract is missing something — escalate rather than reaching into the workspace |
| Skipped stages in state | `StageStatusSkipped` with the reason, written at the re-plan point — not when the stage would have run | "Visible before the developer starts" is the requirement. A skip recorded lazily answers the question too late to be worth anything |
| The markdown pipeline | `deliver-feature` steps 12–29 keep their conditionals, reworded to name the same signals the predicates use. No `loom` dependency is introduced into those steps | Same split as every epic since 77: the executor gets enforcement, the markdown pipeline gets an honest description. Wiring it to `loom route` is a later decision, not a silent one |
| Re-planning mid-run | Out of scope. The route is computed once. If the developer touches a UI file the analysis never mentioned, the route was wrong and stays wrong for that run | Re-planning is a cycle, and `Plan` cannot express one — that is **L2.17**'s mechanism. Building half a loop here would produce the worst version of both |
| No capability registry | Predicates are hardcoded Go reading typed analysis. The planner selecting over declared agent capabilities is **L3.1**, needing **L3.2** | Generality is worth paying for when something needs it. Nothing does yet, and hardcoded predicates are testable today |
| State schema | Adding `StageStatusSkipped` and the route artifact bumps `StateSchemaVersion` to 6; no migration code | Same policy as every prior bump |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**58.9%** as of epic 80). Note that integration tests driving
  the built binary are invisible to the coverage tool — if the floor is threatened, add in-process
  tests, never lower the floor.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — Predicates and the typed route — UNBLOCKED

1. `internal/state/route.go`:
   - `Route` — `[]RouteDecision{Stage, Included bool, Reason string}` plus `SchemaVersion`,
     `Validate()`, and a renderer, following the existing typed-state conventions exactly. Regenerate
     schemas (`go run ./cmd/gen-schemas`); the drift test must stay green.
   - Predicates on `AnalysisState` in the shape of `RequiresArchitect()`, one per conditional stage:
     performance-engineer, data-engineer, accessibility-engineer, visual-qa-engineer, devops-engineer.
     Read `deliver-feature/SKILL.md` steps 12–29 for what each condition means today, and
     `analysis-contract.md` for what may legally be read. **If a condition cannot be expressed from
     the contract's fields, stop and escalate** — do not read the workspace or invent a field.
   - `RouteFor(analysis, plan)` producing a decision for every stage, honouring `Skippable`.
2. Tests: each predicate table-driven, including the case where the derived signals say no and an
   explicit flag says yes; a feature with no infra work skips devops; `code-reviewer` and
   `security-reviewer` are never skipped even when a predicate would say so; the route validates,
   renders, and round-trips.

**Done when**: given a typed analysis, the route is computed deterministically and the never-skip
rule holds. **Commit** (`feat(state): typed route computed from the analysis`), report, PAUSE.

## Phase B — The re-plan point in the executor — BLOCKED BY Phase A

1. `internal/orchestrator/`: `Stage.Skippable`; `StageStatusSkipped` with reason; schema 6. After
   `analyst` completes, compute the route, write it as a typed artifact, and record every skipped
   stage immediately. The run loop skips a `SKIPPED` stage without invoking its provider.
2. Verify — do not build — that the route is bound to `confirm-design`: a test that edits the route
   artifact and asserts the approval is invalidated and the run halts. If it does not fall out,
   STOP and report rather than adding a special case.
3. Tests: a routed run invokes only the included stages; skips are in run state and on the timeline
   before the developer stage starts; a resumed run does not recompute the route (it is an artifact
   with a digest, like any other); an edited analysis demotes the analyst, which cascades to the
   route, which reruns it — assert the whole chain.

**Done when**: `loom run` skips `devops-engineer` on a feature with no infra work, by a recorded
decision visible before the developer stage. **Commit** (`feat(orchestrator): route the run after
the analysis`), report, PAUSE.

## Phase C — CLI surface and prose — BLOCKED BY Phase B

1. `loom state show` displays the route — included and skipped with reasons — and `loom run` prints
   it once at the re-plan point. A human should be able to answer "why isn't devops running?"
   without opening a JSON file.
2. `shared/skills/deliver-feature/SKILL.md` steps 12–29: reword each conditional to name the same
   signals the predicates use, so prose and code state one condition. Do not introduce a `loom`
   dependency into those steps.
3. `cmd/loom/README.md`, `README.md`, `docs/ARCHITECTURE.md`, `docs/roadmaps/BUILD-ROADMAP.md`
   (SHIPPED marker), `shared/DOMAIN_DICTIONARY.md` (**Route** as a term — §6 of design-principles
   requires it).
4. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when** (L3.0): the roadmap done-when holds and every doc describes routing the same way.
**Commit** (`docs(route): the run is routed from the analysis — align prose with L3.0`), report,
PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Re-planning mid-run, or any cycle in `Plan` (**L2.17**)
- A planner selecting over declared agent capabilities (**L3.1**) or the registry it needs (**L3.2**)
- Typing any artifact beyond the route itself (the L2.9 continuation)
- Parallelising the stages the route selects (**L3.3**)
- Policy-driven auto-approval of the route (**L2.16**)
- Making the markdown pipeline call `loom` for its conditionals

## Report format (end of every phase)

```
## Epic 81 Phase <X> Report
- Roadmap item: L3.0 — Route computed from the analysis
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
