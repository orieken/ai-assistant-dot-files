# Lessons Learned

Cross-delivery pattern extraction, produced by the `extract-lessons` skill
(`shared/skills/extract-lessons/SKILL.md`).

## Convention
```
docs/lessons-learned/lessons-<YYYY-MM-DD>.md
```

One file per extraction run (not a fixed cadence — run it whenever enough deliveries have accumulated to
look for patterns). Each file records:
- Recurring security findings (3+ occurrences) and whether a guardrail was proposed/approved/declined
- Recurring code-review rejections (3+ occurrences) and whether a prompt change was proposed/approved/declined
- Architecture decisions that recurred often enough to become a Knowledge Item
- KI usage analytics — which KIs in `shared/knowledge/`/`.claude/knowledge/` are actually being referenced

## Why promotions are gated, not automatic
"Auto-promote" in the sense the TODO originally used it does not mean silently editing `shared/rules/` or
an agent's prompt. Both of those are already governed by `.claude/rules/approval-gates.md` (Gate #7,
"Wiring a New Fitness Function") and by the versioning requirement in `shared/agents/CHANGELOG.md` (Epic 8).
`extract-lessons` drafts the proposed change and stops — a human approves before anything in `shared/rules/`
or `shared/agents/` actually changes. This directory is the record of what was noticed, not of what was
necessarily acted on.

## This is distinct from
- `docs/features/<name>/retrospective.md` — one delivery's narrative.
- `docs/pipeline-retrospectives/` — cross-delivery timing/iteration trends.
- `docs/agent-metrics/` — cross-delivery quality scores.

`lessons-learned` is the only one of the four that looks for *recurring* patterns worth promoting into a
permanent rule, prompt change, or Knowledge Item.
