# Contract: code-review-report.md

**Produced by**: code-reviewer
**Consumed by**: security-reviewer, qa-engineer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## Overall Status`
- `## Design Narrative`
- `## Design Score`
- `## Security Surface`
- `## Performance Surface`
- `## Test Design Review`
- `## Verification of Developer Self-Review`
- `## Feedback for the Developer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## Overall Status` must contain exactly one of `APPROVED` or `CHANGES REQUESTED` (bolded, per the agent's own template) — anything else is a FAIL, since the orchestrator's CHANGES REQUESTED loop parses this literal string.
- `## Design Score` must contain all four dimensions (Clarity, Cohesion, Coupling, Craft) with a numeric 1-5 rating each — a missing dimension is a FAIL.

This is a structural check only. It does not re-judge whether APPROVED was the right call — that's the security-reviewer and qa-engineer's job to catch downstream, and the human's job at the CHANGES REQUESTED checkpoint.
