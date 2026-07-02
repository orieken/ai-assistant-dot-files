# Contract: qa-report.md

**Produced by**: qa-engineer
**Consumed by**: sre-engineer, tech-writer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## Test Files Created`
- `## Test Files Modified`
- `## Coverage Summary`
- `## Test Results`
- `## Accessibility Check`
- `## Bugs Found`
- `## Known Gaps`
- `## Notes for Tech Writer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## Test Results` must show `Failed: 0` — per the agent's own rule, tests must be green before the pipeline proceeds. A non-zero failure count is a FAIL, not a warning.

This is a structural check only. It does not re-run the test suite — that's already qa-engineer's job before writing the report.
