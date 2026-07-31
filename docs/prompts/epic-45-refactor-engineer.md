# Epic 45 — Automated Codemod & Refactoring Agent (`refactor-engineer`)

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 2; re-confirmed open by
`docs/audits/framework-gap-audit-2026-07-31.md` § F5.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

Partial cover exists but nothing owns large-scale structural refactoring:

- `modernization-supervisor` agent — a *coordinator* persona for parallel modernization agents;
  deferred from golden-file fixtures (non-deterministic orchestration, see
  `tests/agents/README.md`).
- `refactor-to-pattern` skill — surgical, single-target GoF/EIP pattern transitions.
- `shared/rules/design-principles.md` § 2 — the named Fowler refactoring-operation vocabulary the
  agent must speak.
- Michael Feathers discipline in root `CLAUDE.md` (docs copy): characterization tests first,
  NEVER refactor and add behavior in the same commit.
- `unit-tester` agent — already builds characterization-test safety nets; the natural upstream
  partner.

The 07-25 audit scoped this as "agent + contract for AST codemods, framework migrations, and
structural refactoring."

## Scope

**Phase A — Design (one commit, then PAUSE for user approval):**

Draft and commit as `docs(agents): investigate refactor-engineer design (Epic 45 Phase A)`:

- `shared/agents/refactor-engineer.md` frontmatter shape + Process outline. Follow the exemplar
  pattern used by the last agent addition (`docs/prompts/done/epic-46-visual-qa-engineer.md` and
  the shipped `shared/agents/visual-qa-engineer.md`). Decide and justify: `model_tier`, tool list
  (this agent MUTATES source — it is not a counter agent), pipeline position (probably standalone
  / on-demand, NOT a deliver-feature phase — say why).
- `shared/contracts/refactoring-contract.md` outline — a refactoring-notes.md must minimally
  carry: named refactoring operations applied (Fowler vocabulary), characterization-test evidence
  BEFORE the refactor, behavior-preservation evidence AFTER (same tests green), complexity
  delta, and an explicit "no behavior added" attestation.
- Relationship ruling: does `modernization-supervisor` delegate to refactor-engineer? Does
  refactor-engineer subsume `refactor-to-pattern` or call it? One paragraph each.

**Phase B — Implementation (after approval; one commit per file):**

1. `shared/agents/refactor-engineer.md` — versioned `1.0.0`, CHANGELOG entry in the same commit
   (`scripts/check-agent-versions-ci.sh` enforces this on PRs)
2. `shared/contracts/refactoring-contract.md`
3. `shared/templates/refactoring-notes.template.md`
4. `tests/agents/refactor-engineer/` — golden-file fixture (input + expected-patterns.txt +
   eval-rubric.md). NOT optional: `health-check.sh` FAILs any non-deferred agent lacking a
   fixture directory. Also extend `contract_for_agent()` in `scripts/test-agents.sh`.
5. Regenerate platform configs — agent count goes 38 → 39, which touches `.roomodes`, every
   Tier 2/3 roster, and `AGENTS.md`: `bash scripts/generate-configs.sh`, then
   `bash scripts/check-parity.sh` must pass.
6. Prose/count updates: README, `docs/AGENT_REFERENCE.md` —
   `bash scripts/check-inventory-drift.sh` must pass.
7. `shared/skills/validate-artifact/SKILL.md` Contract Mapping table + `shared/skills/list-agents`
   remain consistent (verify, update only if they hardcode rosters).

After every commit: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If Phase A concludes the right shape is a *skill* (or an extension of `refactor-to-pattern`)
  rather than an agent, halt and propose demotion — don't build an agent to satisfy the epic
  title.
- If the modernization-supervisor relationship turns circular (each claims to delegate to the
  other), halt with both options sketched.
- If fixture design can't produce deterministic expected-patterns (same trap that deferred
  modernization-supervisor), halt and propose deferred status with README justification instead.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A rulings:
  - Agent vs skill: <choice + one-line rationale>
  - Pipeline position: <standalone | phase X>
  - modernization-supervisor relationship: <ruling>
  - refactor-to-pattern relationship: <ruling>

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, check-parity <pass>, check-inventory-drift <pass>,
agent count 39 everywhere.
```

Go.
