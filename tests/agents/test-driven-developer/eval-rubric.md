# Eval Rubric: test-driven-developer / input-acceptance-criteria.md

- **Tests are written before implementation**: the output shows a failing test (Red phase) before any implementation code — not "here's the implementation + here are the tests for it."
- **Domain-layer check is tested independently of the DB constraint**: a unit test exercises the domain rule through a mock repository — the test does NOT depend on a real database throwing a unique-constraint error.
- **Anti-enumeration is addressed in the test**: there is a test case that verifies the 409 response body does NOT differ based on whether the email exists (or the timing is controlled) — not just a test that checks the 409 is returned.
- **409 response body is exact**: the test asserts the exact JSON `{ "error": "EMAIL_ALREADY_REGISTERED" }` — not just that the status is 409.
- **Refactor step is explicit**: after the Green phase, the output includes a Refactor step that either does a named refactoring (Extract Function, Rename Variable, etc.) or explicitly states "no refactor needed and why."

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
