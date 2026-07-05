# Contract: performance-report.md

**Produced by**: performance-engineer
**Consumed by**: developer (mandates must be implemented), data-engineer (if invoked, for query-pattern alignment)

## Required Sections (exact heading text and level)
- `## 1. Idempotency Guarantees`
- `## 2. Timeout & Circuit Breaker Mandates`
- `## 3. N+1 Query Prevention`
- `## 4. Hot Path Caching`
- `## Notes for Developer`

## Validation Rule
Each of the four numbered sections must contain a `**Status**` line — a section missing its status line means
the agent skipped judging that risk category rather than explicitly clearing it.

This is a structural check only. It does not verify the mandates are technically correct or sufficient —
that's the human PAUSE checkpoint and the developer's/code-reviewer's job.
