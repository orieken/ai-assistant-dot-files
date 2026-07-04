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
