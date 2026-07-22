<!--
Template for qa-report.md. Consumed by the qa-engineer agent.
Structure defined here; contract in shared/contracts/qa-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# QA Report: [Feature Name]

## Test Files Created
- `tests/test_feature.py` — [what it tests]

## Test Files Modified
- `tests/test_existing.py` — [what was added]

## Coverage Summary
- Acceptance criteria covered: X/Y
- Total new tests: N
- Total test assertions: N

## Test Results
- Passed: N
- Failed: 0 (all failures resolved)
- Skipped: N (with reason)

## Accessibility Check
- [Violations flagged as `[A11Y]` or "None identified"]

## Bugs Found
- [Bug description]: [How it was fixed] — or "None"

## Known Gaps
- [Any acceptance criteria that couldn't be tested and why]

## Notes for Tech Writer
- [Any behavior that was surprising or non-obvious that docs should clarify]
