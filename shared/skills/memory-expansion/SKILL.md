---
name: memory-expansion
description: Promotes recurring lessons and delivery retrospectives into portable Knowledge Items (KIs) in shared/knowledge/. Opposing-force pair with memory-compression.
triggers:
  keywords: [memory expansion, promote lesson, expand memory, ki promotion]
  intentPatterns: ["promote lessons to knowledge items", "expand memory corpus"]
standalone: true
---

## When To Use

Use when:
- Synthesizing recurring patterns from `docs/lessons-learned/` or `docs/features/*/retrospective.md`.
- Proposing new Knowledge Items (KIs) for `shared/knowledge/`.

Do NOT use when:
- Deduplicating or removing KIs (use `memory-compression`).
- Auditing KI schema compliance (use `memory-auditor`).

## Context To Load First

1. `shared/knowledge/README.md`
2. `shared/templates/ki.template.md`
3. `docs/lessons-learned/` files

## Process

1. **Scan Retrospectives & Lessons**: Glob `docs/lessons-learned/*.md` and recent `retrospective.md` files.
2. **Identify Candidate Patterns**: Extract recurring architectural fixes or process rules that appeared across 2+ deliveries.
3. **Draft Knowledge Item**: Create draft KI in `.claude/feature-workspace/proposed-ki.md` conforming to `shared/schemas/ki-frontmatter.schema.json`.
4. **Pause for Human Confirmation**: Prompt human operator to approve proposed KI before writing to `shared/knowledge/`.
5. **Persist Approved KI**: Copy approved KI to `shared/knowledge/<slug>.md` and update `shared/memory-registry.json`.

## Output Format

Draft artifact at `.claude/feature-workspace/proposed-ki.md`:

```markdown
---
name: kebab-case-slug
tags: [tag1, tag2]
domain: framework
created: YYYY-MM-DD
---

# Title

## Context
[Problem description]

## Pattern
[Solution details]
```

## Guardrails

- Never write to `shared/knowledge/` without explicit human confirmation ("approve" or "yes").
- Must conform strictly to `shared/templates/ki.template.md`.

## Standalone Mode

Operates purely using local file editing tools (`Read`, `Write`, `Glob`, `Grep`).
