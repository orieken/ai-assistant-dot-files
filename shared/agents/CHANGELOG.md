# Agent Changelog

Tracks version bumps for every agent in `shared/agents/`. Every prompt edit that changes agent *behavior*
(not just a typo or formatting fix) requires a version bump here, enforced by the pre-commit hook in
`scripts/hooks/pre-commit` (see that file's header comment for how to enable it — it's opt-in, not wired up
automatically for you).

## Versioning
Semantic-ish, not strict SemVer:
- **Patch** (1.0.x): wording/clarity fixes that don't change behavior.
- **Minor** (1.x.0): new process step, new output section, expanded guardrail — additive, backward compatible.
- **Major** (x.0.0): changed output contract (update the matching file in `shared/contracts/` too if one
  exists), removed/renamed a process step, or changed tool access.

## How to add an entry
When you bump an agent's `version:` frontmatter field, add a row under a new dated heading here in the same
commit — the pre-commit hook checks for exactly this.

---

## 2026-07-06 — Cursor native skills/agents compatibility (Epic 30, Phase 1)

Cursor shipped native Agent Skills (`.cursor/skills/*/SKILL.md`) and subagent (`.cursor/agents/*.md`)
support using the same open standard this repo's `shared/skills/`/`shared/agents/` already follow
(confirmed against `cursor.com/docs/skills` and `cursor.com/docs/subagents`). Scoping that integration
surfaced two prerequisite issues affecting all 24 agents, fixed together here since both touch the
same files in the same pass:

1. **`model: sonnet` → `model: inherit`**: every agent hardcoded a specific model regardless of what
   the user's own session was running. Both Claude Code and Cursor subagents default to `inherit` when
   the field is omitted and accept the literal keyword explicitly — confirmed via Cursor's own docs and
   a live Claude Code frontmatter check this session. `inherit` lets each subagent match whatever model
   the operator already chose for their session instead of forcing Sonnet unconditionally.
2. **Frontmatter preamble relocated**: 23 of 24 agents (every one except `documentation-manager`) had
   a "Read `.claude/rules/design-principles.md`..." instruction *before* their opening `---`, which is
   invisible to `health-check.sh`'s lenient grep-anywhere frontmatter check and tolerated by Claude
   Code's loader, but would very likely break Cursor's stricter parser once agents are symlinked
   directly (a planned later phase of this epic). Moved into the body as the agent's own first
   instruction, using canonical `shared/rules/` paths instead of the Claude-Code-only `.claude/rules/`
   prefix, so every file now starts with `---` on line 1.

| Agent | Version | Change |
|---|---|---|
| accessibility-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| analyst | 1.1.0 -> 1.1.1 | Patch: preamble relocated, model: inherit |
| api-test-generator | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| architect | 1.1.1 -> 1.1.2 | Patch: preamble relocated, model: inherit |
| chaos-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| code-reviewer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| context-engineer | 2.1.1 -> 2.1.2 | Patch: preamble relocated, model: inherit |
| data-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| dependency-auditor | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| developer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| devops-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| documentation-manager | 2.0.0 -> 2.0.1 | Patch: model: inherit (already had correct frontmatter placement from its 2026-07-05 rewrite -- no preamble to relocate) |
| dx-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| finops-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| modernization-supervisor | 1.0.0 -> 1.0.1 | Patch: preamble relocated (2-file variant: design-principles.md + ARCHITECTURE_RULES.md), model: inherit |
| performance-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| product-owner | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| qa-engineer | 1.1.0 -> 1.1.1 | Patch: preamble relocated, model: inherit |
| release-manager | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| security-reviewer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| spec-writer | 1.1.0 -> 1.1.1 | Patch: preamble relocated, model: inherit |
| sre-engineer | 1.0.0 -> 1.0.1 | Patch: preamble relocated, model: inherit |
| tech-writer | 1.1.0 -> 1.1.1 | Patch: preamble relocated, model: inherit |
| test-driven-developer | 1.0.0 -> 1.0.1 | Patch: preamble relocated (2-file variant: design-principles.md + ARCHITECTURE_RULES.md), model: inherit |

---

## 2026-07-05 — External audit fixes (api-generator portability, context-engineer casing tolerance)

| Agent | Version | Change |
|---|---|---|
| context-engineer | 2.1.0 -> 2.1.1 | **Patch**: an external audit found that step 6's Prior-Deliveries grep for `**Owning Context**` exact-matches, but the archived `docs/features/context-engineering-framework/analysis.md` uses `**Owning context**` (lowercase c) — a real historical drift that would silently miss that feature's retrospective. Fixed by making the lookup case-insensitive (documented in both the agent and `shared/skills/context-engineer/SKILL.md` twin) instead of retroactively editing the archived doc, since the feature archive is treated as an immutable historical record elsewhere in this framework. No output format change. |

Also (no agent version bump, non-agent fixes from the same external audit):
- `scripts/api-generator/index.ts`: removed hardcoded personal machine paths (`/Users/oscarrieken/...`) for
  Go/TS client output dirs — now CLI args or `API_GENERATOR_GO_DIR`/`API_GENERATOR_TS_DIR` env vars.
- `scripts/api-generator/package.json` + new `scripts/api-generator/README.md`: marked the tool explicitly
  experimental/unsupported (it has no tests and isn't wired into `scripts/ci-check.sh` or CI) instead of
  leaving that unstated.

---

## 2026-07-05 — documentation-manager narrowed to ad-hoc-session counterpart of promote-memory

| Agent | Version | Change |
|---|---|---|
| documentation-manager | 1.0.0 -> 2.0.0 | **Major**: changed output contract entirely. Previously wrote directly to `ARCHITECTURE.md`/`RUNBOOKS.md`/`GOTCHAS.md`/`ONBOARDING.md` with no review step -- an undocumented overlap with `memory-engineer`/`promote-memory`/`extract-lessons` (added later, in the Memory Engineering epic) that `docs/AGENT_REFERENCE.md` flagged explicitly. Redesigned as the ad-hoc-session counterpart to `promote-memory`: now produces Candidate Records (Source/Type/Evidence/Tags/Expiration) via the same Memory Contract, requires explicit human approval before any KI/ADR/rule/living-doc edit, and retires `GOTCHAS.md` as a target (gotchas are Knowledge Items now, via `create-ki`). Still manual/on-demand, not hooked to auto-run after every session. |

---

## 2026-07-04 — Memory Engineering epic (v2 scope, split from AOS/v3 prototyping)

| Agent | Version | Change |
|---|---|---|
| context-engineer | 2.0.0 -> 2.1.0 | New: Proactive RAG step now checks whether the task's question is KI/ADR-shaped (invoke `search-ki`, unchanged default) or broader (invoke the new `query-memory` skill instead, which also covers the feature archive and DOMAIN_DICTIONARY.md). Additive — existing behavior and output format unchanged. Applied identically to both the agent and its `shared/skills/context-engineer/SKILL.md` twin in the same edit, to avoid repeating the twin-drift bugs found across three independent audits this session |

---

## 2026-07-03 — Cross-agent audit fixes (independent review via docs/runbooks/self-audit-prompt.md)

| Agent | Version | Change |
|---|---|---|
| spec-writer | 1.0.0 -> 1.1.0 | Twin drift fix: the agent's Critique Report used emoji verdicts (`READY ✅ \| NEEDS WORK ⚠️`, `✅/⚠️` per row) while `shared/skills/spec-writer/SKILL.md` used plain text (`READY \| NEEDS WORK`, `PASS/FAIL`). Standardized both to plain text — matches the `PASS/FAIL` vocabulary used everywhere else in the framework (`validate-artifact`, contracts) and avoids emoji-rendering inconsistency across the 6 target platforms |
| architect | 1.1.0 -> 1.1.1 | Patch: removed a self-contradictory parenthetical ("read at step 3 as per instructions") on what is actually step 2 of its own process list — wording fix, no behavior change |

Also (no version bump — pure renames/config changes, not agent behavior changes):
- `modernization-swarm.md` -> `modernization-supervisor.md`, `test-driven-development-agent.md` ->
  `test-driven-developer.md`: filenames now match their own `name:` frontmatter field, like every other
  agent in `shared/agents/`.
- `context-engineer`'s skill twin (`shared/skills/context-engineer/SKILL.md`) had its Prune Recommendations
  bullet format aligned to match the agent's (proper `- [ ]` instead of backtick-wrapped `` `[ ]` ``, plus
  the reason column the skill was missing).

---

## 2026-07-03 — Context Engineering audit: contract + agent/skill heading realignment

| Agent | Version | Change |
|---|---|---|
| context-engineer | 1.4.0 -> 2.0.0 | **Major**: `shared/agents/context-engineer.md`'s Output Format headings (`## Scope & Boundaries`, `## Relevant Knowledge Items (KIs) & ADRs`, `## Pinpoint Files to Open...`, `## Pruning Checklist...`) had drifted from its own "standalone twin" in `shared/skills/context-engineer/SKILL.md` (`## 1. Scope and Boundaries` ... `## 7. Token Budget`) and was missing a `## 3. Global Rules and Constraints` section entirely. Realigned the agent's headings to match the skill's numbered format exactly, and added the missing section. This was found while adding `shared/contracts/context-manifest-contract.md` (see below) — the contract would have failed every real run against the agent's old headings. New contract added: `context-manifest.md` now gets the same `validate-artifact` structural gate every other pipeline artifact already had; wired into `deliver-feature` as new step 7 (renumbering all subsequent steps by one). |

---

## 2026-07-02 — Team Topologies alignment

| Agent | Version | Change |
|---|---|---|
| architect | 1.0.0 -> 1.1.0 | New "Team Topology Fit" sub-step under Strategic Domain Design: for any Context Crossing, invokes `team-topology-check` (new skill) against the new `TEAM_TOPOLOGY.md` registry to flag a stale Collaboration interaction mode or a bypassed Platform team — a Conway's-Law-shaped version of the existing Distributed Monolith anti-pattern check. New "Team Topology Fit" line added inside the already-required `## Bounded Context` section (no contract change needed — the heading itself is unchanged) and a new Anti-Pattern Check checklist item |

---

## 2026-07-02 — Epic 14 KI infrastructure

| Agent | Version | Change |
|---|---|---|
| context-engineer | 1.0.0 -> 1.1.0 | Step 5 (Proactive RAG) now invokes the `search-ki` skill instead of ad-hoc grepping `shared/knowledge/`, `.claude/knowledge/`, and `docs/adrs/` directly — additive, output format unchanged |

---

## 2026-07-02 — Epic 17 context decay and bounded-context pruning

| Agent | Version | Change |
|---|---|---|
| qa-engineer | 1.0.0 -> 1.1.0 | Step 2 now gets `analysis.md`'s acceptance criteria/edge cases via `summarize-artifact` instead of a full read (Context Decay — 2 phases old by this point) |
| tech-writer | 1.0.0 -> 1.1.0 | Step 1 now gets `analysis.md`'s scope via `summarize-artifact` instead of a full read (same reason) |
| context-engineer | 1.1.0 -> 1.2.0 | New step: auto-prune Pinpoint Files by bounded-context mapping (exclude other contexts' files unless the analysis explicitly flags a crossing) and by change surface (exclude infrastructure/migration files for UI-only tasks) |

---

## 2026-07-02 — Proactive self-invocation

| Agent | Version | Change |
|---|---|---|
| context-engineer | 1.2.0 -> 1.3.0 | Description now says "Use PROACTIVELY before starting any task that touches 3+ files, a new feature area, or unfamiliar code" instead of only firing on explicit request — closes the gap where context engineering only ever applied inside `deliver-feature`, never in ad-hoc sessions. Additive framing change, no process/output format change |

---

## 2026-07-02 — Cross-feature learning: same-bounded-context retrieval

| Agent | Version | Change |
|---|---|---|
| context-engineer | 1.3.0 -> 1.4.0 | New step: search `docs/features/*/analysis.md` for prior deliveries in the same Bounded Context (recency-independent) and surface their `retrospective.md` lessons in a new "Prior Deliveries in This Bounded Context" context-manifest.md section. Closes the gap where a same-area mistake from more than 3 deliveries ago was invisible to `analyst`'s recency-based feedback loop |
| analyst | 1.0.0 -> 1.1.0 | Step 5 (feedback loop) now treats context-manifest.md's "Prior Deliveries in This Bounded Context" as the primary, recency-independent same-area check, with the existing 3-most-recent-deliveries scan kept as a secondary check for general cross-cutting process trends |

---

## 2026-07-02 — Initial versioning rollout
All 24 agents in `shared/agents/` set to `1.0.0` — no prior version was tracked before this.

| Agent | Version | Change |
|---|---|---|
| accessibility-engineer | 1.0.0 | Initial version |
| analyst | 1.0.0 | Initial version |
| api-test-generator | 1.0.0 | Initial version |
| architect | 1.0.0 | Initial version |
| chaos-engineer | 1.0.0 | Initial version |
| code-reviewer | 1.0.0 | Initial version |
| context-engineer | 1.0.0 | Initial version |
| data-engineer | 1.0.0 | Initial version |
| dependency-auditor | 1.0.0 | Initial version |
| developer | 1.0.0 | Initial version |
| devops-engineer | 1.0.0 | Initial version |
| documentation-manager | 1.0.0 | Initial version |
| dx-engineer | 1.0.0 | Initial version |
| finops-engineer | 1.0.0 | Initial version |
| modernization-supervisor | 1.0.0 | Initial version |
| performance-engineer | 1.0.0 | Initial version |
| product-owner | 1.0.0 | Initial version |
| qa-engineer | 1.0.0 | Initial version |
| release-manager | 1.0.0 | Initial version |
| security-reviewer | 1.0.0 | Initial version |
| spec-writer | 1.0.0 | Initial version |
| sre-engineer | 1.0.0 | Initial version |
| tech-writer | 1.0.0 | Initial version |
| test-driven-developer | 1.0.0 | Initial version |
