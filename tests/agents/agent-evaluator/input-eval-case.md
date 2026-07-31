# Agent Eval Case: analyst

## Input Fixture
`tests/agents/analyst/input-feature-spec.md` — Password Reset via Email feature.

## Expected Patterns
See `tests/agents/analyst/expected-patterns.txt`:
- `expir` — token expiry addressed
- `enumeration` — non-enumeration constraint covered
- `notification` — notification service dependency identified
- `token` — reset token mechanics described

## Eval Rubric
See `tests/agents/analyst/eval-rubric.md` — 5 qualitative criteria including:
- Three distinct edge cases (non-enumeration, token expiry, single-use) are separate items
- Bounded context correctly identified with notification service as a Context Crossing
- New dependency flagged with resilience angle
- Feature Flag strategy explicitly named
- Acceptance criteria are behavior-focused (given/when/then)

## Actual Output Location
`tests/agents/analyst/actual-output.md` (not checked in — generate locally)
