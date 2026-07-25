---
name: memory-compression
description: Deduplicates, consolidates, and summarizes stale Knowledge Items (KIs) in shared/knowledge/. Opposing-force pair with memory-expansion.
triggers:
  keywords: [memory compression, compress memory, deduplicate ki, ki summary]
  intentPatterns: ["deduplicate knowledge items", "compress memory corpus"]
standalone: true
---

## When To Use

Use when:
- Resolving duplicate or overlapping Knowledge Items flagged by `memory-auditor`.
- Consolidating obsolete KIs into summary references.

Do NOT use when:
- Creating new KIs from retrospectives (use `memory-expansion`).
- Editing core framework rules (use human review).

## Context To Load First

1. `shared/knowledge/README.md`
2. `shared/memory-registry.json`
3. Audit reports from `.claude/audits/`

## Process

1. **Read Audit Findings**: Read latest `memory-auditor` report.
2. **Identify Merge Candidates**: Evaluate candidate duplicate KI pairs.
3. **Draft Consolidated KI**: Create proposed merged KI in `.claude/feature-workspace/compressed-ki.md`.
4. **Pause for Human Confirmation**: Prompt human operator with proposed merge diff.
5. **Execute Consolidation**: Replace duplicate files with consolidated KI and update `shared/memory-registry.json`.

## Output Format

```markdown
# Memory Compression Proposal

## Target KIs to Merge
- `shared/knowledge/ki-old-1.md`
- `shared/knowledge/ki-old-2.md`

## Consolidated Outcome
- `shared/knowledge/ki-consolidated.md`

## Proposed Diff
[Merged content]
```

## Guardrails

- Never delete or overwrite a Knowledge Item without explicit human approval ("approve merge" or "confirm").

## Standalone Mode

Operates purely using local file editing tools.
