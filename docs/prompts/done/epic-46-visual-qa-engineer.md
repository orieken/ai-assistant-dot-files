# Epic 46 — Visual QA & Screenshot Diffing Agent

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 2.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

Framework has 36 pipeline agents but nothing dedicated to visual regression / screenshot diffing. Saturday-family repos have `Saturday.ML` / `saturday-ml` packages that do screenshot diffing + click-heatmap generation — but no framework-level `visual-qa-engineer` agent orchestrates them per-feature.

Investigation needed first: read `saturday-monorepo/packages/saturday-ml/` (if present) and any related patterns to confirm what visual-QA capability actually exists across the Saturday ecosystem before designing the agent.

## Scope

**Phase A — Investigate** (one commit): map current visual-QA capability across saturday-monorepo. Answer: is this an agent that WRAPS existing Saturday.ML capability, or does it prescribe new tooling?

Draft:
- `shared/agents/visual-qa-engineer.md` frontmatter shape + Process outline
- `shared/contracts/visual-qa-report-contract.md` — what a visual-qa-report.md must contain
- `shared/templates/visual-qa-report.template.md` — output template
- Pipeline positioning — where does it fit relative to `qa-engineer`? (Probably parallel; runs on UI features)

Commit as: `docs(agents): investigate visual-qa-engineer design (Epic 46 Phase A)`.

**Pause for user approval before Phase B.**

**Phase B — Implementation** (multiple commits, one per file):

1. `shared/agents/visual-qa-engineer.md` — the agent, marked conditional "if UI"
2. `shared/contracts/visual-qa-report-contract.md` — required sections
3. `shared/templates/visual-qa-report.template.md` — the template
4. Update `shared/skills/deliver-feature/SKILL.md` — inject visual-qa-engineer into the review phase for UI-touching features (alongside accessibility-engineer)
5. Update `docs/patterns/deliver-feature-workflow.md` diagram to include the new agent
6. Update `shared/skills/validate-artifact/SKILL.md` Contract Mapping table

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- If `Saturday.ML` is not the right foundation (missing capability, wrong shape) — halt, describe what the agent would need but doesn't have.
- If pipeline positioning creates a race with accessibility-engineer or qa-engineer — halt, propose ordering.
- If Phase A reveals visual-QA already exists as a skill (not agent) — halt, propose demoting Epic 46 to "extend existing skill" or promoting the skill to an agent.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A findings:
  - Saturday.ML capabilities: <list>
  - Pipeline positioning proposal: <parallel with accessibility-engineer | after qa-engineer | ...>
  - Agent-vs-skill choice: <agent, rationale>

Phase B commits (if approved):
  <sha> <message>
  ...

Verification: deliver-feature-workflow.md diagram updated, contract validates, template referenced correctly.
```

Go.
