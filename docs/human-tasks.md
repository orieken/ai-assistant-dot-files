# Human-Only Tasks

Actions that require a human — not a fireable handoff prompt. Parallel to `docs/prompts/` (agent-doable) and `docs/aos/prompts/` (AOS-migration-specific).

Every entry: what to do, why it's here (i.e., why not a prompt), estimated time, where the source-of-truth lives if there's a paired doc.

---

## Outstanding

*(empty — no pending human-only tasks)*

---

## Decisions embedded inside queued handoffs

These aren't standalone tasks — they're decision points that a subagent will halt at when it fires the corresponding prompt. Listed here so you know what to expect when a subagent asks.

| Prompt | Decision the subagent will halt on |
|---|---|
| `docs/prompts/framework-hygiene-sweep.md` Item 6 | What to do with `docs/audits/perplex-audit.md` (commit, gitignore, or promote to lessons-learned) |
| `docs/prompts/capture-session-history.md` | Which of four options (A: run /retrospective, B: run /extract-lessons, C: manual archive, D: git log is the record) for capturing this session |
| `docs/prompts/automate-deliver-feature.md` § Deliverables | Confirming which gates from `shared/rules/approval-gates.md` are policy-eligible vs. always-human |
| `docs/aos/prompts/phase-2-governance.md` § Op 2.5 | Confirming producer/counter mapping when a "producer" role turns out to be a human, not an agent |
| `docs/aos/prompts/phase-3-runtime.md` § Op 3.1 | RAG backend ADR (LightRAG vs sqlite-vec vs pgvector — recommend LightRAG per design pack; needs ADR write-up before Phase 3 code starts) |
| `docs/aos/prompts/phase-4-policy.md` § Op 4.5 | Depends on the `automate-deliver-feature` design being done first — Phase 4 will halt at Op 4.5 if that gate-classification hasn't landed |
| `saturday-mcp/docs/prompts/mcp-expand-milestone-2.md` § Phase A | M2 tool selection + sqlite-vec dependency approval before M2 code starts |
| `saturday-mcp/docs/prompts/mcp-expand-milestone-3.md` § Phase A | Persona design decisions + composite workflow shape approval |
| Post-M3 in saturday-mcp | Rename decision (`saturday-mcp` → `framework-mcp`/`context-mcp`/`craftsmanship-mcp`) — deferred per mcp-expand plan scope decision until broad-scope reality is visible |

---

## Convention

When a new human-only task appears:

1. Add it under "Outstanding" with the same structure (What, Why human, Estimated time, Where the source-of-truth lives, Rollback if applicable).
2. When completed, move the entry to a new "Completed" section at the bottom of this file (with the date and commit/tag SHA if applicable) — don't delete, so the historical record persists.
3. If the task becomes agent-doable through a future capability (e.g., a subagent that can spawn a scratch project), demote it to a handoff prompt under `docs/prompts/` and remove from this list.

---

## Completed

### 1. AOS Phase 1 v3.0.0 verification + tag — 2026-07-30

- **Verification results** (all four checks green):
  1. `bash scripts/health-check.sh --verbose` — **214 passed, 0 warned, 0 failed.** (The two pre-existing WARNs from `framework-hygiene-sweep.md` Items 7+8 were cleared before verification.)
  2. `bash scripts/check-parity.sh` — **PASS.** All platform configs in sync with `shared/` canonical source, no DRIFT.
  3. `deliver-feature` against a scratch project cloned from this repo — **PASS.** Artifacts landed in expected paths with expected required sections.
  4. `.claude/telemetry/events.jsonl` NOT created by scratch run — **PASS.** The opt-in guarantee holds; nothing in v3.0 emits telemetry by default.
- **Tag applied**: `v3.0.0` at HEAD.
- **Pushed to origin**: yes.
