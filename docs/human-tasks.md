# Human-Only Tasks

Actions that require a human — not a fireable handoff prompt. Parallel to `docs/prompts/` (agent-doable) and `docs/aos/prompts/` (AOS-migration-specific).

Every entry: what to do, why it's here (i.e., why not a prompt), estimated time, where the source-of-truth lives if there's a paired doc.

---

## Outstanding

*(none)*

---

## Decisions embedded inside queued handoffs

These aren't standalone tasks — they're decision points that a subagent will halt at when it fires the corresponding prompt. Listed here so you know what to expect when a subagent asks.

*(none in this repository)*

---

## External — saturday-mcp

These decision points belong to the external `saturday-mcp` repository and have not been verified during
this repository's documentation cleanup.

| Prompt | Decision the subagent will halt on |
|---|---|
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

### 2. AOS Phase 2 v3.1.0 verification + tag — 2026-07-31

- **Verification results** (all checks green):
  1. `bash scripts/health-check.sh --verbose` — **259 passed, 0 warned, 0 failed.**
  2. `bash scripts/check-parity.sh` — **PASS.** All platform configs in sync, no DRIFT.
  3. Backward-compat confirmed: `validate-artifact` is purely structural with no hooks dependency; projects without `.claude/hooks/` see zero behavior change from v3.0.0.
  4. `git push origin main && git push origin v3.1.0` — pushed 7 commits + tag.
- **Tag applied**: `v3.1.0` at `594696d`.
- **Pushed to origin**: yes.
- **Note**: Push required GitHub secret-scanning bypass for `tests/agents/privacy-auditor/` fixture files (synthetic Stripe test key used as audit target). Bypassed as test fixture false-positive.

### 1. AOS Phase 1 v3.0.0 verification + tag — 2026-07-30

- **Verification results** (all four checks green):
  1. `bash scripts/health-check.sh --verbose` — **214 passed, 0 warned, 0 failed.** (The two pre-existing WARNs from `framework-hygiene-sweep.md` Items 7+8 were cleared before verification.)
  2. `bash scripts/check-parity.sh` — **PASS.** All platform configs in sync with `shared/` canonical source, no DRIFT.
  3. `deliver-feature` against a scratch project cloned from this repo — **PASS.** Artifacts landed in expected paths with expected required sections.
  4. `.claude/telemetry/events.jsonl` NOT created by scratch run — **PASS.** The opt-in guarantee holds; nothing in v3.0 emits telemetry by default.
- **Tag applied**: `v3.0.0` at HEAD.
- **Pushed to origin**: yes.
