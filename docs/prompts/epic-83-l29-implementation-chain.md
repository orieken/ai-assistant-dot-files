# Epic 83 — Type the Implementation Chain (L2.9, second cut)

Source: session discussion, 2026-08-31, following epic 82. Continues roadmap item **L2.9**, whose
first cut (epic 79) typed the `analyst → architect` hop and whose review artifact was typed by
**L2.17** (epic 82). Eleven of the fifteen pipeline artifacts still pass as markdown, and three
items are waiting on more of them: **L2.10** (retire summarization from the data path), **L2.11**
(semantic validation needs typed state to hang rules on), and **L2.18** (bounded contract retries).

This cut types the **implementation chain**: what the developer produces and what the two stages
that consume it produce — `implementation-notes`, `security-report`, `qa-report`. The review loop
already reads a typed verdict; typing the developer's own output closes its other half.

It also folds in two call sites of **L2.10**, which are reachable now and were not when that item
was written: `qa-engineer` and `tech-writer` both call `summarize-artifact` on `analysis.md`, an
LLM summarising an artifact that has been typed since epic 79.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — **L2.9**, **L2.10**, **L2.11**, **L2.18**.
2. `internal/state/` — `analysis_state.go`, `review_state.go`, `projection.go`, `render*.go`,
   `retrieval.go`. Four artifacts have been through this pattern; follow it rather than inventing a
   fifth style. Note especially how `ReviewState` models a verdict the executor reads as a field.
3. `shared/contracts/implementation-contract.md`, `security-contract.md`, `qa-contract.md` — the
   headings and, importantly, the **content rules** each declares: a Critical/High security finding
   must carry a non-empty `**Fix applied**` line, and `## Test Results` must show `Failed: 0`.
4. `shared/agents/qa-engineer.md` step 2 and `tech-writer.md` step 1 — the two `summarize-artifact`
   call sites this epic replaces, and their stated reason (Context Decay).
5. `shared/skills/summarize-artifact/SKILL.md` — what stays (human-facing prose, the step 37a
   retrieval surrogate) and what goes.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Scope | `implementation-notes`, `security-report`, `qa-report`. Nothing else | One coherent hop: the developer's output and its two consumers. The end-of-pipeline reports (docs, devops, observability) stay markdown — nothing evaluates a condition over them, so typing them buys nothing yet |
| Contract content rules become typed invariants | `Failed: 0` becomes a numeric field validated as zero-or-fail; a Critical/High finding without a `FixApplied` becomes a validation error naming the finding | These are the two places where a contract already asserts something semantic. Typed state is what lets them be checked rather than grepped — and it is a preview of L2.11 done where the contract already asked for it |
| The severity ladder is typed | `Severity` enum (`CRITICAL`…`INFO`) on security findings, matching the agent's own vocabulary | The pipeline's "block on Critical findings" rule reads this. A string comparison against prose is the defect L2.9 exists to remove |
| Two `summarize-artifact` sites are replaced | `qa-engineer` (acceptance criteria + edge cases) and `tech-writer` (feature intent + scope) receive projections of `AnalysisState` instead | Both read an artifact that is already typed, so the deterministic replacement exists today. **This is part of L2.10's scope, taken deliberately** — the epic must say so, and L2.10's roadmap entry must be updated to reflect what remains |
| What stays with `summarize-artifact` | The step 37a retrieval surrogate, and any human-facing summary | The surrogate feeds `memory-registry`'s retrieval tier and has its own consumers; converting it is a different problem with a different done-when |
| Untyped consumers keep working | Rendered markdown views under the contract filenames, as before | `sre-engineer`, `devops-engineer`, and `accessibility-engineer` still read `implementation-notes.md`. The typed hop must stay invisible to them |
| Agent prompt files | `qa-engineer.md` and `tech-writer.md` **are** edited here, unlike every prior cut — their steps 2 and 1 name `summarize-artifact` explicitly | Those two steps describe *how the agent gets its input*, which is exactly what changes. Bump versions and add `shared/agents/CHANGELOG.md` entries; the pre-commit hook requires it. The developer/security/qa output instructions are still appended by the provider, not written into the agent files |
| Multiple upstreams | `Stage.Consumes` becomes `[]string`, and `StageInput.UpstreamState` becomes a map keyed by upstream stage ID. Each upstream is projected separately and delivered as its own labelled block | qa-engineer reads three upstreams per its contract — implementation notes, the security report, and the analysis's acceptance criteria. Merging them into one flat object would lose provenance and let same-named fields (`feature`, `filesModified`) collide silently |
| Projections are keyed by consumer stage and upstream kind | `(consuming stage ID, upstream Kind)`, not by the consumer's own state kind | What a stage *reads* and what it *writes* vary independently. `tech-writer` produces markdown in this cut but still needs projections, and keying by upstream kind alone would give the architect and the tech-writer the same slice of the analysis when they demonstrably need different fields |
| State schema | Adding three kinds bumps `SchemaVersion` in `internal/state` to 2; run-state's own version is unaffected | The two version numbers are independent, and only state documents change shape |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**60.9%** as of epic 82). Integration tests driving the built
  binary are invisible to the coverage tool — if the floor is threatened, add in-process tests,
  never lower the floor.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — Three typed artifacts — UNBLOCKED

1. `internal/state/implementation_state.go`, `security_state.go`, `qa_state.go`, each following the
   established pattern exactly: struct, `SchemaVersion`, `Validate()`, generated schema, renderer
   under the contract's filename, retrieval frontmatter, registration in `StageSchemas`/`Decode`/
   `RenderView`.
2. Type the content rules the contracts already assert: `QAState.TestResults.Failed` is numeric and
   a non-zero value fails validation; `SecurityFinding` carries a typed `Severity` and a
   `FixApplied`, and a CRITICAL or HIGH finding without one fails validation naming the finding.
3. Tests: round-trips; each renderer emits every contract heading verbatim; the two content rules
   reject what the contracts say they reject and accept what they allow; frontmatter is complete.

**Done when**: three artifacts are typed, and the two semantic rules the contracts already declare
are enforced by validation rather than by a grep. **Commit** (`feat(state): typed implementation,
security, and QA artifacts`), report, PAUSE.

## Phase B — Wiring and projections — BLOCKED BY Phase A

1. Declare the three stages' `StateKind` in the built-in plan, with `Consumes` where a projection
   applies. Provider instructions come from the schema as before — no agent prompt edits for these
   three.
2. Projections, keyed by `(consumer stage, upstream kind)` per the design table: `security-reviewer`
   and `qa-engineer` from `implementation-notes`; `qa-engineer` from `security-report`;
   `tech-writer` from `qa-report`. Field selection only. Re-key the two existing projections
   (`architect ← analysis`, `developer ← review`) onto the same table rather than leaving two
   addressing schemes.
3. Replace the two `summarize-artifact` call sites with projections of `AnalysisState`:
   `qa-engineer` receives acceptance criteria and edge cases; `tech-writer` receives feature intent
   and scope. Edit those two agent files, bump their versions, add CHANGELOG entries.
4. Tests: the typed stages exchange data with no markdown on the path; a security report with an
   unfixed CRITICAL fails its stage; the projections carry what each consumer's contract needs and
   omit the rest; untyped consumers still find their markdown views.

**Done when**: the implementation chain exchanges typed state, and no LLM call sits between the
analysis and the two agents that used to summarise it. **Commit** (`feat(orchestrator): typed
implementation chain and analysis projections`), report, PAUSE.

## Phase C — Contracts, docs, and the L2.10 boundary — BLOCKED BY Phase B

1. The three contracts gain their Typed State sections, as `analysis-contract.md` and
   `review-contract.md` have.
2. **Update L2.10's roadmap entry** to record what this epic took and what remains — an item whose
   scope has been partly absorbed must say so, or it will be re-done or dropped.
3. `shared/skills/summarize-artifact/SKILL.md`: state plainly that it is no longer on the
   inter-stage data path for typed artifacts, and what it is still for.
4. `cmd/loom/README.md`, `README.md`, `docs/ARCHITECTURE.md`, BUILD-ROADMAP status for L2.9's
   second cut, `shared/DOMAIN_DICTIONARY.md` if a new term is introduced.
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when**: the docs describe one typing story across seven artifacts, and L2.10's remaining
scope is accurate. **Commit** (`docs(state): the implementation chain is typed — align prose`),
report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Typing the end-of-pipeline reports (docs, devops, observability) or `context-manifest`
- The step 37a retrieval surrogate, or any other `summarize-artifact` call site
- Semantic validation beyond the two rules the contracts already declare (**L2.11**)
- Bounded contract retries or running `validate-artifact` under the executor (**L2.18**)
- Changing what any agent *produces* — only how `qa-engineer` and `tech-writer` receive input

## Report format (end of every phase)

```
## Epic 83 Phase <X> Report
- Roadmap item: L2.9 (second cut) — the implementation chain
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
