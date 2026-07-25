---
name: learning-engine
description: Analyzes past pipeline delivery retrospectives and extracts candidate lessons for docs/lessons-learned/. Opposing-force pair with forgetting-engine.
triggers:
  keywords: [learning engine, extract learning, process feedback, retrospective analysis]
  intentPatterns: ["extract lessons from retrospectives", "run learning engine"]
standalone: true
---

## When To Use

Use when:
- Running a feedback sweep across completed feature deliveries in `docs/features/`.
- Proposing new entries for `docs/lessons-learned/`.

Do NOT use when:
- Expiring or archiving old lessons (use `forgetting-engine`).

## Context To Load First

1. `docs/lessons-learned/`
2. `docs/features/*/retrospective.md`

## Process

1. **Scan Feature Retrospectives**: Glob `docs/features/*/retrospective.md`.
2. **Extract Recurrent Process Failures**: Identify repeated test breakages or contract retry patterns.
3. **Draft Lessons Learned Proposal**: Produce `.claude/feature-workspace/proposed-lessons.md`.
4. **Pause for Human Confirmation**: Prompt human operator to confirm adding lesson entry.
5. **Persist Approved Lesson**: Copy approved lesson to `docs/lessons-learned/lessons-YYYY-MM-DD.md`.

## Output Format

```markdown
# Proposed Lessons: YYYY-MM-DD

## Executive Summary
[Brief description of pattern]

## Findings & Recommendations
- [ ] Recommendation 1
```

## Guardrails

- Never write to `docs/lessons-learned/` without explicit human confirmation ("approve" or "yes").

## Standalone Mode

Operates purely using local file editing tools.
