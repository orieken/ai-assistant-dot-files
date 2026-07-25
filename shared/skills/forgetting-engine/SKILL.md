---
name: forgetting-engine
description: Identifies and flags obsolete or superseded lessons and Knowledge Items for expiration. Opposing-force pair with learning-engine.
triggers:
  keywords: [forgetting engine, expire lesson, archive ki, decay memory]
  intentPatterns: ["identify obsolete lessons", "run forgetting engine"]
standalone: true
---

## When To Use

Use when:
- Sweeping `docs/lessons-learned/` or `shared/knowledge/` for obsolete rules superseded by newer architecture.
- Flagging items for draft expiration.

Do NOT use when:
- Learning new lessons from retrospectives (use `learning-engine`).

## Context To Load First

1. `shared/knowledge/README.md`
2. `docs/lessons-learned/`

## Process

1. **Scan Corpus**: Glob `shared/knowledge/*.md` and `docs/lessons-learned/*.md`.
2. **Identify Obsolete Patterns**: Find items superseded by explicit ADR decisions in `docs/adrs/`.
3. **Draft Expiration Proposal**: Produce `.claude/feature-workspace/proposed-expirations.md`.
4. **Pause for Human Confirmation**: Prompt human operator to approve expiration list.
5. **Execute Expiration**: Archive approved obsolete files to `docs/lessons-learned/archive/`.

## Output Format

```markdown
# Proposed Expirations

## Obsolete Items Flagged
- `shared/knowledge/ki-obsolete.md` — Superseded by ADR-002
- `docs/lessons-learned/lessons-2025-01-01.md` — Obsolete framework setup
```

## Guardrails

- Never delete or move files without explicit human confirmation ("approve expiration" or "confirm").

## Standalone Mode

Operates purely using local file editing tools.
