# Contract: accessibility-report.md

**Produced by**: accessibility-engineer
**Consumed by**: qa-engineer (keyboard/screen-reader test points), security-reviewer (if invoked next)

## Required Sections (exact heading text and level)
- `## Evaluation Summary`
- `## Findings & Fixes`
- `## Notes for QA`

## Validation Rule
`## Evaluation Summary` must contain all four of: `**Semantic HTML**`, `**Interactive Elements**`,
`**ARIA & Labels**`, `**Keyboard Navigation**` — each with a Pass/Fail/Notes value. A missing category means
that risk area was never actually evaluated, not that it passed silently.

This is a structural check only. It does not verify a fix is actually accessible — that requires the human
PAUSE checkpoint or real assistive-technology testing, neither of which this contract can substitute for.
