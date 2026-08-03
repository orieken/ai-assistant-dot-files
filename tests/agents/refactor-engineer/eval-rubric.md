# Eval Rubric: refactor-engineer / input-order-processor.ts

- **Halts for missing tests**: the agent detects that no tests exist for `order-processor.ts`
  and either invokes `unit-tester` or explicitly halts and documents the escalation in
  `## Pre-Refactor State` before proceeding. It does NOT begin refactoring without establishing
  a test baseline.

- **Correct complexity measurement**: the output identifies `processOrder` at complexity 9 and
  `applyPromotions` at complexity 8 (or equivalent values matching the input's inline comments),
  not a vague statement that "the code is complex."

- **Named Fowler operations only**: every structural change is labeled with an exact operation
  name from `shared/rules/design-principles.md` §2 (e.g., "Extract Function", "Replace
  Conditional with Polymorphism", "Guard Clause" / "Extract Guard Clauses"). Generic descriptions
  like "simplified" or "cleaned up" are not acceptable.

- **Mixed-concern smell identified**: the agent notes that `OrderProcessor` conflates validation,
  pricing, tax calculation, inventory persistence, and notification — and applies at least one
  operation that addresses this (e.g., extracting `validateOrder`, `calculateTax`, or a
  `DiscountCalculator` via Extract Function or a refactor-to-pattern call).

- **Post-refactor complexity < 7 for both functions**: the `## Post-Refactor State` section
  shows both `processOrder` and `applyPromotions` (or their extracted successors) at complexity
  < 7, or carries a documented exception explaining why one cannot reach < 7 without behavioral
  change.

- **No-behavior-added attestation present**: the `## No-Behavior-Added Attestation` section
  contains the explicit phrase "no behavior added" or "no new behavior," not just an implication.

- **Bugs are noted, not fixed**: if the agent observes any bugs in `applyPromotions` (e.g., the
  VIP30 promo behaves differently for non-platinum tiers without documentation), it records them
  in `## Bugs Observed (Not Fixed)` and does NOT silently fix them in the same run.

## How to Grade

For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet
has no supporting quote, mark it FAIL and state what's missing. A rubric PASS requires all
bullets graded PASS.
