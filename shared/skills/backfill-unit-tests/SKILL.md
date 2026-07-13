---
name: backfill-unit-tests
description: Writes unit tests for existing code without modifying it, then automatically runs code-reviewer against just the new test files -- for raising coverage on working code or building a characterization-test safety net around legacy code before a refactor/migration. Coordinates the unit-tester and code-reviewer agents.
triggers:
  keywords: ["backfill tests", "add unit tests to", "raise coverage", "characterization tests", "wrap tests around legacy", "/backfill-unit-tests"]
  intentPatterns: ["Add unit tests to *", "Backfill coverage for *", "Write characterization tests for *", "/backfill-unit-tests *"]
standalone: true
---

## When To Use
When the user wants tests added to code that already exists and should not change — raising coverage on
trusted code, or building a characterization-test safety net before refactoring or migrating legacy code.
Accepts a file, directory, or module.

Do NOT use when the code doesn't exist yet or is expected to change to satisfy the tests — use
`test-driven-developer` instead (tests-first, code conforms to the tests). Do NOT use inside a
`deliver-feature` run — `qa-engineer` already covers testing for a feature just implemented, including its
own Legacy Code section for touching pre-existing files in that context; this skill is for standalone
coverage/characterization work with no accompanying feature delivery.

## Context To Load First
1. `CLAUDE.md` — project constraints and clean code rules
2. `ARCHITECTURE_RULES.md` — architectural guardrails, particularly the 85% coverage threshold
3. `shared/rules/approval-gates.md` — gate #6 (Writing Files out of Boundary) governs any seam `unit-tester`
   flags as a blocker

## Process
1. **Resolve the scope** — the file, directory, or module the user pointed at.
2. **Capture baseline coverage** — invoke `run-tests` for the target scope before any new tests exist, so
   the final report shows a real before/after delta. Skip only if the target has no existing test
   infrastructure to run at all (typical for genuinely untested legacy code) — note "no baseline" rather
   than fabricating a 0%.
3. **Run `unit-tester`** against the resolved scope. It determines coverage-backfill vs. characterization
   mode itself and produces `.claude/feature-workspace/unit-test-report.md`.
4. **Run `code-reviewer`** against only the new/modified test files from step 3 — never against the
   untouched source. This is the counterbalance `unit-tester` doesn't have on its own: nothing else checks
   whether the tests it wrote are well-structured, correctly scoped, and free of the complexity/SOLID issues
   this framework flags everywhere else.
5. **Capture final coverage** — invoke `run-tests` again for the same scope, compute the delta against step
   2's baseline.
6. **Produce the combined report** — display to the user.

## Output Format

```markdown
# Test Backfill: [target]

## Mode
Coverage backfill | Characterization (legacy/migration)

## Coverage Delta
- **Before**: N% (or "No baseline — untested")
- **After**: N%

## Tests Written
- `tests/test_x.py` — [what it covers]

## Code Review (tests only)
| Check | Result | Details |
|---|---|---|
| Cyclomatic complexity | PASS / FAIL | [specifics] |
| Test independence / no brittle coupling | PASS / FAIL | [specifics] |
| Naming conventions | PASS / FAIL | [specifics] |

## Behavior Notes (characterization mode only)
- [Bug-like behavior captured as-is, not fixed] / "N/A"

## Blocked by Structure
- [Code that needs a seam to be testable — proposed seam, awaiting explicit "approve file write"] / "None"

## Summary
[2-3 sentences: what's covered now, what's still not, whether anything needs the user's attention]
```

## Guardrails
- Never let `unit-tester` or `code-reviewer` modify the source under test — both operate read-only against
  it. A seam is the one structural exception, and it requires the explicit approval described in
  `shared/rules/approval-gates.md` gate #6 — never performed automatically by this skill.
- Never mark the backfill complete if `code-reviewer` finds a Critical-severity issue in the new tests —
  fix the tests (not the source) before reporting done.
- In characterization mode, never "correct" behavior the tests capture — that's what the next
  refactor/migration is for, not this skill.

## Standalone Mode
Pure local analysis + test execution — no external calls required beyond whatever the project's own test
runner needs.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
