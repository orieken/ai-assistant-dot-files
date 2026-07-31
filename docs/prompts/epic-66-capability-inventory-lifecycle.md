# Epic 66 — Capability-Inventory Lifecycle (dedupe / deprecate skills and agents)

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #6). The gap:
`forgetting-engine` expires stale KIs, but nothing expires, merges, or deprecates stale skills
and agents — and sprawl is already observable.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

- Observed duplication (verify, don't assume it's the whole list):
  `shared/skills/analyze-complexity/` vs `shared/skills/complexity-check/` — overlapping
  "complexity" triggers, same job (one adds fail thresholds, one is "on-demand"). Two skills
  matching one keyword degrades trigger precision for every router.
  Also agent+skill name collisions: `spec-writer` and `context-engineer` exist as both — these
  MAY be an intentional wrapper pattern (skill invokes agent); confirm before flagging.
- Inventory scale: 38 agents, ~65 skills, 10 platforms' generated configs that all enumerate
  them. Removing/renaming anything ripples through `generate-configs.sh`, `.roomodes`, rosters,
  `check-inventory-drift.sh` counts, and fixture directories.
- Existing machinery to reuse: `forgetting-engine` skill (staleness → proposed expiration →
  human approves), agent/skill frontmatter schemas (`shared/schemas/*.schema.json`),
  `shared/agents/CHANGELOG.md` versioning discipline, health-check frontmatter validation.

## Scope

**Op 1 — Deprecation convention (one commit).**
Add optional frontmatter fields to both schemas: `status: active | deprecated` (absent =
active) and `superseded_by: <name>`. Update the frontmatter contracts, JSON schemas, and
health-check: a deprecated entry with no `superseded_by` is a WARN; a `superseded_by` pointing
at a nonexistent capability is a FAIL. Deprecated entries stay installed (backward compat) but
generators annotate them in rosters ("deprecated — use X").
Commit: `feat(schemas): capability deprecation convention (Epic 66 Op 1)`

**Op 2 — Inventory-lifecycle audit flow (one commit).**
Extend `forgetting-engine`'s SKILL.md scope to cover the capability inventory (or add a sibling
section — read it first and match its structure): scan for trigger-keyword collisions,
description overlap, and agent+skill name pairs; produce merge/deprecate PROPOSALS for human
review, never auto-apply. Same opposing-force discipline as the existing memory pair.
Commit: `feat(skills): inventory lifecycle audit in forgetting-engine (Epic 66 Op 2)`

**Op 3 — Resolve the observed duplicate (one commit, after proposing in-report and getting
approval).**
Run the Op 2 flow for real on `analyze-complexity` vs `complexity-check`. Expected outcome: merge
into one (probably keep `analyze-complexity`'s thresholds + `complexity-check`'s on-demand
framing), deprecate the other via Op 1's convention. Regenerate configs, re-run
`check-parity.sh` + `check-inventory-drift.sh` (counts change or deprecation annotations appear —
update prose counts accordingly).
Commit: `refactor(skills): merge complexity skills per lifecycle flow (Epic 66 Op 3)`

**Op 4 — Rule on the agent+skill name pairs (one commit).**
Investigate `spec-writer` and `context-engineer` pairs: if intentional wrapper pattern, document
it as a named convention (in `docs/patterns/` or the skills README) so future collisions are
deliberate; if accidental, propose resolution via the Op 2 flow.
Commit: `docs(patterns): agent+skill pairing convention ruling (Epic 66 Op 4)`

After every op: `bash scripts/health-check.sh` + `bash scripts/check-parity.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- Op 3: if merging would break any documented external reference to the deprecated skill name
  (README examples, blog drafts, platform configs that can't express deprecation), halt with the
  reference list.
- If the trigger-collision scan surfaces MORE than the known duplicate pair, don't resolve them
  all in this epic — list them as proposals in the report; Op 3 handles exactly one, as the
  pattern-proving case.
- If schemas can't take optional fields without breaking existing validation, halt on Op 1.

## Report (under 150 words)

```
Commits: <sha> x4
Deprecation convention: <fields, WARN/FAIL semantics>
Collision scan results: <n> candidate pairs (resolved: 1, proposed: <list>)
complexity merge: <survivor, deprecated, roster annotation confirmed>
Agent+skill pairs ruling: <intentional-documented | accidental-proposed>
parity + inventory-drift: <pass>
```

Go.
