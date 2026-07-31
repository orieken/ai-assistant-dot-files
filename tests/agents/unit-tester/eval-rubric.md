# Eval Rubric: unit-tester / input-source.ts

- **Source is not modified**: the unit-tester's output contains only test code — no edits to `discount-calculator.ts` itself, even if the agent notices a bug.
- **The 40% cap boundary is tested**: there is a test for the exact cap — e.g., `enterprise + SAVE20 = 45%, capped to 40%` — not just a test that discount is "applied."
- **Free tier (zero discount) is an explicit test case**: there is a test confirming `free` tier with no coupon returns `basePrice * 1.0` — not assumed to be covered by other cases.
- **Each coupon code is tested independently**: SAVE10 and SAVE20 each have at least one test case verifying their specific discount amount, not just that "a coupon was applied."
- **Tests follow Arrange/Act/Assert**: each test body has a clear setup, a single call to `calculateDiscount`, and an assertion — no multi-step logic or multiple calls per test.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
