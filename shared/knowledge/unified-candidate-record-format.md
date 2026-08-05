---
name: unified-candidate-record-format
tags: [promote-memory, extract-lessons, on-call, five-whys, lessons, incidents, memory-pipeline]
domain: memory-pipeline
created: 2026-08-05
---

Every skill that produces lessons destined for the knowledge base uses the same Candidate Record
format — the one defined in `shared/skills/promote-memory/SKILL.md` — regardless of where the
lesson originates.

## The format

```markdown
### Candidate: [short title]
- **Source**: [skill that produced it, e.g. "on-call — 2026-08-05 login outage"]
- **Type**: KI | ADR-worthy | Rule-change-worthy | Lesson
- **Evidence**: [observation, pattern, or incident that warrants this candidate]
- **Tags**: [comma-separated, to be used as KI frontmatter tags]
- **Expiration condition**: [when this candidate would no longer apply]
- **Existing overlap checked**: [KIs/rules already covering this topic, if any]
```

## Sources that use this format

| Source skill | Where candidates appear |
|---|---|
| `promote-memory` | `.claude/feature-workspace/<feature>/lessons.md` |
| `extract-lessons` (Steps 1–5: code-review/retrospective patterns) | `docs/lessons-learned/` |
| `extract-lessons` (Step 6: incident-feature pair mining) | same output document |
| `on-call` (Step 6: incident record) | `docs/incidents/<date>-<slug>.md` |
| `five-whys` (Step 6: incident record) | same incident file |

## Why one format matters

The gate machinery in `promote-memory` evaluates a Candidate Record as a unit: it checks
evidence strength, overlap, and type before deciding KI / ADR / rule-change. If different source
skills used different formats, the gate would need source-specific logic to evaluate each one —
coupling the gatekeeper to the source, which violates the Open-Closed Principle for this
pipeline.

The unified format means: any source that produces Candidate Records in this format is
automatically compatible with the gated promotion machinery, today and for future sources.

## Practical implication

When adding a new lesson-sourcing skill (e.g., a customer-feedback analyzer), produce Candidate
Records in this exact format. Do not invent new fields unless `promote-memory`'s gate logic is
updated to read them.
