# Framework Gap Audit — 2026-07-31

Follow-up to `framework-gap-audit-2026-07-25.md`. That audit spawned 10 Phase-10 epics (42, 44, 46,
48–51, 53, 55, 57) — **all 10 shipped** as of 2026-07-31. This audit re-baselines: what's verified
green, what broke since, and what remains open.

## 1. Verified Green (checked this audit, not assumed)

| Surface | Result |
|---|---|
| `scripts/health-check.sh` (local, macOS) | 259 passed, 0 warned, 0 failed |
| `scripts/check-parity.sh` (local) | 32 passed, 0 failed — all 10 platforms in sync |
| `scripts/test-agents.sh` | all fixtures pass; 32/38 agents covered, 6 specialists deferred (documented in `tests/agents/README.md`) |
| `scripts/check-inventory-drift.sh` | no drift in authoritative prose docs |
| `scripts/test-install.sh` | PASS (inside ci-check container) |
| `docs/human-tasks.md` Outstanding | empty |
| `.roomodes` | valid YAML, 38 custom modes (verified with local python3 + PyYAML) |

## 2. Findings

### F1 — BUG: `ci-check.sh` container lacks PyYAML → false CI failures (quick fix)

`scripts/ci-check.sh` provisions its ubuntu:24.04 container with **only `python3`** (line ~50), but
`scripts/check-parity.sh`'s Roo Code section (lines 240–276) does `import yaml`. The `ImportError`
is swallowed by `2>/dev/null`, so the check misreports:

```
DRIFT .roomodes — file exists but customModes array is empty or YAML is invalid
DRIFT .roomodes agent roster — 38 agents missing
```

This cascades: `check-parity.sh` FAILs → `health-check.sh --verbose` (which shells out to it) also
FAILs → `ci-check.sh` final result is **2 passed, 2 failed** even though the repo is actually clean.
The tool whose entire purpose is "match GitHub Actions ubuntu-latest" diverges from it — real GH
runners happen to have PyYAML preinstalled, so real CI passes while the local pre-push check fails.

Fix (three parts, all small):
1. `ci-check.sh`: add `python3-yaml` to the container's `apt-get install`.
2. `check-parity.sh`: distinguish "yaml module unavailable" (SKIP/WARN with install hint) from
   "YAML actually invalid" (FAIL) — stop swallowing ImportError as drift.
3. `.github/workflows/framework-ci.yml`: the parity job silently depends on runner-preinstalled
   PyYAML; add an explicit `pip install pyyaml` (or `apt-get install python3-yaml`) step so a
   future runner-image change can't break CI mysteriously.

### F2 — DOC DRIFT: AOS Phase 2 shipped but not closed out

Phase 2 (governance skeleton) fully landed 2026-07-25 — commits `28d398d`→`48102a4`, changelog has a
v3.1.0 entry, `governance-pairs.md` exists, health-check detects counter agents. But:
- `docs/aos/prompts/README.md` still lists `phase-2-governance.md` without the `— DONE` marker
  Phase 1 got.
- **No `v3.1.0` git tag exists** (only v1.0.0/v2.0.0/v3.0.0). Phase 1's tag went through a
  human-verification task in `docs/human-tasks.md`; Phase 2's equivalent entry was never added, and
  the Outstanding section is empty — the tag step fell through the cracks.

Fix: mark phase-2 DONE in the AOS prompts README, and add a human-task entry for v3.1.0
verification + tag (tagging is a human step per repo convention).

### F3 — DOC INACCURACY: prompts README claims audit reasons that don't exist

`docs/prompts/README.md` says the un-drafted epics (43, 45, 47, 52, 54, 56, 58) have "reasons
documented inline in the audit (superseded, resolved, or scope-clarification needed)" — but
`framework-gap-audit-2026-07-25.md` contains no such inline annotations; those checklist items are
simply unchecked. Either add the per-epic disposition notes to the 07-25 audit (this audit's §3
drafts them) or correct the README claim.

### F4 — HOUSEKEEPING: untracked working-tree debt

- `.idea/` (JetBrains IDE state) is untracked and not gitignored — add to `.gitignore`.
- `docs/blog-posts/memory_engineering_prompts/` (10 prompt files) + redundant
  `memory_engineering_prompts.zip` are untracked. Content inspection: these are handoff prompts for
  memory-engineering work that **already shipped** (memory registry → `shared/memory-registry.json`,
  memory-engineer skill, promote-memory, query-memory, LightRAG stanza all exist). Precedent:
  `framework-hygiene-sweep.md` deleted the redundant AOS zip. Recommended: delete the zip; either
  delete the extracted directory or move it under an archive/done location if the provenance is
  worth keeping. (Human call — don't auto-delete.)

### F5 — CARRIED-OVER ROADMAP GAPS (from 07-25 audit, still open)

Dispositions as of this audit:

| Epic | Status | Disposition |
|---|---|---|
| 43 — MCP tool packaging (`shared/mcp/` exporter) | Open, partially mitigated | `shared/mcp-patterns/` now exists (patterns, not an exporter). Standalone epic still valid if packaging framework skills as MCP servers for Cursor/Claude Desktop is wanted. |
| 45 — `refactor-engineer` agent + contract | Open | No pipeline agent for codemods/structural refactors; `modernization-supervisor` + `refactor-to-pattern` remain the partial cover. |
| 47 — `ship-feature` release/PR skill | Open — was audit priority #3 | Still no automated branch/commit/PR/release-tag skill. `release-manager` agent covers planning, not execution. Highest-value remaining standalone epic. |
| 52 — Semantic vector RAG (LightRAG) | Subsumed by AOS Phase 3 Op 3.1 | Blocked on the RAG-backend ADR human decision (`docs/human-tasks.md` decision table). Don't draft separately. |
| 54 — Telemetry event-loop wiring | Subsumed by AOS Phase 3 | Runtime wiring belongs to the orchestration layer Phase 3 builds. |
| 56 — Producer-auditor loop | Subsumed by AOS Phase 3 | The trinity-native workflow refactor (Op set 3.x) is exactly this. |
| 58 — documentation-auditor automation | Open, small | `documentation-auditor` agent + fixtures exist; wiring it into a repeatable health-check/scheduler workflow is not done. Small epic or a scheduler-skill config example. |

### F6 — AOS Phases 3 and 4 not yet executed (known, not drift)

`phase-3-runtime.md` (v3.2.0: `shared/rag/`, orchestration, learning/forgetting hooks, trinity
workflows) and `phase-4-policy.md` (v3.3.0: policies layer) prompts exist and are the biggest
remaining roadmap items.

**Correction (post-audit follow-up, same day):** this audit initially repeated
`docs/human-tasks.md`'s claim that Phase 3 halts at a RAG-backend ADR decision. That claim was
stale — ADR-002 (Accepted 2026-07-22, commit `b55c0dd`) already decided the retrieval strategy
(graduated corpus-aware: LLM-as-retriever / BM25 sqlite-fts5 / sqlite-vec; LightRAG deferred to
v2; pgvector rejected) and the same commit amended Phase 3 to match. **Phase 3 is fireable today.**
Phase 4 still depends on Phase 3 plus the approval-gate classification from
`automate-deliver-feature`.

## 3. Actionable TODO Checklist

Quick fixes (single sitting, no design needed):
- [ ] **[F1] Fix ci-check.sh PyYAML gap** — add `python3-yaml` to container install; make
      check-parity.sh SKIP (not FAIL) when the yaml module is missing; pin PyYAML in
      `framework-ci.yml`.
- [ ] **[F2] Close out AOS Phase 2** — mark `phase-2-governance.md` DONE in
      `docs/aos/prompts/README.md`; add v3.1.0 verify+tag entry to `docs/human-tasks.md`.
- [ ] **[F3] Fix prompts README claim** — add per-epic disposition notes (per §F5 table) to the
      07-25 audit, or reword the README.
- [ ] **[F4a] Gitignore `.idea/`.**
- [ ] **[F4b] Human decision: delete `memory_engineering_prompts.zip` (+ keep/archive/delete the
      extracted directory).**

Standalone epics (draft as `docs/prompts/epic-NN-*.md` handoffs when picked up):
- [ ] **[Epic 47] `ship-feature` release/PR skill** — branch creation, conventional-commit
      formatting, PR description compiling spec/retrospective links, release tagging. Respect
      Approval Gates #1/#2. (Priority #1 of the remaining standalone epics — was #3 in 07-25.)
- [ ] **[Epic 58] documentation-auditor automation** — scheduler-driven or health-check-adjacent
      repeatable invocation; smallest remaining epic.
- [ ] **[Epic 45] `refactor-engineer` agent + refactoring contract.**
- [ ] **[Epic 43] MCP exporter (`shared/mcp/`)** — package framework skills as standalone MCP
      servers; builds on `shared/mcp-patterns/`.

AOS phases (prompts already exist under `docs/aos/prompts/` — fire, don't draft):
- [ ] **AOS Phase 3** (`phase-3-runtime.md`, absorbs Epics 52, 54, 56) — ~~blocked on RAG-backend
      ADR~~ **unblocked**: ADR-002 already decided the retrieval strategy (see §F6 correction).
- [ ] **AOS Phase 4** (`phase-4-policy.md`) — after Phase 3 + the approval-gate classification
      decision from `automate-deliver-feature`.

## 3b. Structural Gap Review (same-day follow-up) — Epics 61–68

A second pass the same day asked a different question: not "what's on the backlog" but "what is
the framework structurally missing." Eight gaps, in priority order, each drafted as a handoff in
`docs/prompts/epic-6N-*.md`:

| Epic | Gap | Why it matters |
|---|---|---|
| **61** | No automated prompt-regression eval harness | Golden-file fixtures only run manually in live sessions (`tests/agents/README.md`); CI validates committed outputs only. A model change silently alters behavior across all 38 agents. Everything exists (fixtures, patterns, rubrics, agent-eval grading logic) except a headless runner. |
| **62** | Human gate decisions + token spend absent from telemetry | `shared/telemetry/event-schema.md` has no gate-approval/rejection/edit events — human corrections, the richest feedback signal, evaporate; `extract-lessons` mines what agents wrote, not what humans fixed. `pipeline-trace.json` has duration/status/iterations/estimated budget but no actual token spend, so cost-optimizer/finops-engineer have no data. |
| **63** | No parallel-delivery story | `.claude/feature-workspace/` is a singleton; `deliver-feature` halts if another run's state file exists. One feature in flight per project. Touches checkpointing, resume-pipeline, rollback history — needs a design pause. |
| **64** | Conventions prescribe linter configs the framework doesn't ship | Every `<language>-conventions.md` names a tool + cap (eslint 6, detekt 6, SwiftLint 6, clippy 6…) but no `shared/configs/` exists — installed projects hand-author fitness functions from prose, defeating guardrail #7. |
| **65** | Framework itself never threat-modeled | KIs load as trusted agent context (ADR-003 sync = instruction-injection vector via a compromised org memory repo), hooks execute on events, prompts distribute by symlink. `threat-model` skill exists but was never pointed inward. |
| **66** | No capability-inventory lifecycle | `forgetting-engine` expires KIs only. Observed sprawl: `analyze-complexity` vs `complexity-check` near-duplicates; agent+skill name collisions. No deprecation/merge mechanism for 38 agents + ~65 skills. |
| **67** | Production never teaches the system; no bugfix-weight path | All learning loops mine pipeline artifacts; `on-call`/`five-whys` outputs aren't wired to memory. Everything is feature-shaped — small fixes bypass the framework or pay full ceremony. |
| **68** | No framework-version marker in installs | `install.sh` writes no version record; install-vs-upstream drift detection is manual archaeology. |

Checklist:
- [ ] **[Epic 61] Prompt-regression eval harness** — headless fixture runs + rubric grading + regression diff.
- [ ] **[Epic 62] Gate-decision + cost telemetry** — event-schema extension + pipeline-trace spend fields.
- [ ] **[Epic 63] Parallel-delivery workspace isolation** — workspace-per-feature or worktree; design pause first.
- [ ] **[Epic 64] Ship `shared/configs/` linter configs** — configs matching convention caps + drift cross-check.
- [ ] **[Epic 65] Framework threat model** — STRIDE over memory/hooks/sync/install; provenance + "memory is data" rule.
- [ ] **[Epic 66] Capability-inventory lifecycle** — deprecation convention; resolve observed duplicates first.
- [ ] **[Epic 67] Production feedback loop + `deliver-bugfix`** — incident → candidate KI flow; lightweight fix pipeline.
- [ ] **[Epic 68] Install version marker** — `.framework-version` written at install, read by health-check/updater.

## 4. Priority Recommendation

1. **F1** — the local CI-parity tool currently cries wolf on every run; fix before it trains you to
   ignore red.
2. **F2 + F3 + F4a** — 15 minutes of drift cleanup, keeps the "docs match reality" property the
   framework sells.
3. **Epic 47 (`ship-feature`)** — closes the last unautomated stretch of the delivery pipeline.
4. ~~RAG-backend ADR decision~~ **Fire AOS Phase 3** — the ADR blocker turned out to be stale
   (ADR-002 already decided it); the largest remaining body of work is ready to execute.
