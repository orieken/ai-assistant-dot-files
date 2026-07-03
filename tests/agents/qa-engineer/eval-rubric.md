# Eval Rubric: qa-engineer / input-implementation-notes.md

- **Three distinct tests, matching the three "Notes for QA" bullets**: separate test cases for (1) identical
  response shape/status for registered vs. unregistered email, (2) token expiry after 30 minutes, and
  (3) single-use enforcement (second confirm attempt fails) — not one combined or partial test.
- **Enumeration test framed as a negative/equivalence check**: verifies the two responses are the *same*
  (status code and payload shape), not just that each individually "works" — the whole point of the test is
  the absence of a difference.
- **Tests actually reported as run and green**: the report states tests were executed and passing (e.g.,
  "Failed: 0"), consistent with the qa-engineer contract's rule that tests must be green before reporting
  done — not just a list of written test names with no run confirmation.
- **Edge cases tied back to the implementation's actual interface**: references the real `ResetToken` shape
  (`expiresAt`, `usedAt`) or the actual endpoint paths from "Files Created," showing the tests target this
  implementation specifically rather than being generic boilerplate.

## How to Grade
For each bullet above, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no
supporting quote, mark it FAIL and say what's missing.
