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
