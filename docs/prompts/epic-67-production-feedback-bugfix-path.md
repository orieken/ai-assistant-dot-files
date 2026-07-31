# Epic 67 — Production Feedback Loop + Lightweight Bugfix Path

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #7). The gap: every learning
loop mines pipeline artifacts — production never teaches the system — and everything is
feature-shaped, so small fixes either bypass the framework (no learning) or pay full-pipeline
ceremony.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

- Learning today: `retrospective` → `promote-memory` (per-delivery, real rejection criteria),
  `extract-lessons` (cross-delivery mining), `learning-engine` (Phase 2 skill; Phase 3 wires its
  hook). ALL inputs are pipeline artifacts under `docs/features/`.
- Incident handling today: `on-call` skill (active incident response), `five-whys` skill
  (structured RCA), `engineering:incident-response` plugin skill (triage/comms/postmortem).
  None of their outputs feed `promote-memory`, KIs, lessons-learned, or rule/prompt changes. A
  production incident caused by a delivered feature never reaches the agent that introduced it.
- Work-item shapes: `spec-writer` already covers "features, bugs, spikes, or chores" as work
  items — the SPEC side handles bugs; the DELIVERY side has only `deliver-feature` (full
  pipeline) and `deliver-atdd`.
- `documentation-manager` is precedent for "capture knowledge from a session that never went
  through the pipeline" — the incident flow is its sibling.

## Scope

**Half 1 — Incident → memory pipeline (Ops 1–2):**

**Op 1**: Define the incident artifact + flow. `five-whys` and `on-call` gain a closing step:
persist an incident record (proposed home: `docs/incidents/<date>-<slug>.md` — confirm no
existing convention first) carrying: affected feature (link to `docs/features/<name>/` if it
exists), root cause, the five-whys chain, and a **Candidate Records** section using
`promote-memory`'s exact contract (so the same human-gated promotion machinery consumes it
unchanged — reuse, don't fork).
Commit: `feat(skills): incident records feed promote-memory (Epic 67 Op 1)`

**Op 2**: Close the loop to origin. When an incident links to a delivered feature,
`extract-lessons` treats the pair (feature artifacts + incident record) as first-class input:
"which pipeline stage should have caught this?" becomes a standing question whose answer is a
proposed rule/prompt change (human-gated, as extract-lessons already does). Update
`shared/memory-registry.json` to register `docs/incidents/` as a source (lexical).
Commit: `feat(skills): extract-lessons mines incident-feature pairs (Epic 67 Op 2)`

**Half 2 — `deliver-bugfix` (Ops 3–4):**

**Op 3**: `shared/skills/deliver-bugfix/SKILL.md` — a deliberately thin pipeline:
reproduce-first discipline (characterization test that fails, per the Feathers rule already in
CLAUDE.md) → developer fix → code-reviewer → qa run. Reuses existing agents and contracts;
artifacts land in `docs/features/<bug-slug>/` (same archive, so retrospective/extract-lessons
see bugfixes too — that's the point). Approval gates apply unchanged. Explicitly document what
it SKIPS vs deliver-feature (analyst/architect/threat-model by default) and the escalation
trigger ("if the fix touches >N files or a contract boundary, stop — this is a feature").
Commit: `feat(skills): deliver-bugfix lightweight pipeline (Epic 67 Op 3)`

**Op 4**: Wiring + docs: `new-feature`/`spec-writer` route bug-type work items to
`deliver-bugfix`; README + workflow pattern doc updated; inventory counts
(`check-inventory-drift.sh`) updated.
Commit: `docs: route bug work items to deliver-bugfix (Epic 67 Op 4)`

After every op: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If `promote-memory`'s Candidate Record contract can't represent an incident-sourced record
  without changes, halt — changing that contract ripples into documentation-manager too.
- If `deliver-bugfix` starts accreting phases until it's deliver-feature-lite-but-actually-heavy,
  halt and cut — the entire value is the weight difference.
- If an incidents-directory convention already exists somewhere (RUNBOOKS, on-call), adopt it —
  halt only if two conflicting conventions exist.

## Report (under 150 words)

```
Commits: <sha> x4
Incident record home: <path convention; pre-existing or new>
Candidate Records reuse: <unchanged | contract issue found>
deliver-bugfix phases: <list> (skipped vs feature: <list>)
Escalation trigger to full pipeline: <the rule>
health-check + inventory-drift: <pass>
```

Go.
