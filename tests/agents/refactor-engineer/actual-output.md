---
feature: "order-processor-complexity-reduction"
bounded_context: "Ordering"
domain_terms: ["Order", "OrderItem", "CustomerTier", "PaymentMethod", "Promotion", "Cyclomatic Complexity Threshold"]
files_touched:
  - "src/orders/order-processor.ts"
  - "src/orders/order-validator.ts"
  - "src/orders/discount-calculator.ts"
  - "src/orders/promo-strategy.ts"
  - "src/orders/order-processor.spec.ts"
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Refactoring Notes: OrderProcessor Complexity Reduction

## Summary

`src/orders/order-processor.ts` was refactored to bring both of its methods under the
Cyclomatic Complexity Threshold (< 7) and to separate the five concerns the class was
conflating — validation, pricing, promotion, tax, and persistence/notification. `processOrder`
measured at cyclomatic complexity 9 and `applyPromotions` at 8; both are now at 4 and 3
respectively, with every extracted collaborator at 3 or below. The file had **zero test
coverage**, so refactoring was blocked at Phase 0 step 2 until `unit-tester` produced a
14-test characterization suite pinning current behavior — including the behavior that is
almost certainly a bug (see `## Bugs Observed (Not Fixed)`). Six named Fowler operations were
applied across four files; the same 14 tests are green before and after. No behavior was added.

## Target Files

| File | Function / Method | Pre-Refactor Complexity | Post-Refactor Complexity |
|---|---|---|---|
| src/orders/order-processor.ts | `processOrder` | 9 | 4 |
| src/orders/order-processor.ts | `applyPromotions` | 8 | 3 |
| src/orders/order-validator.ts | `validateOrder` (new) | — | 3 |
| src/orders/order-processor.ts | `calculateSubtotal` (new) | — | 2 |
| src/orders/order-processor.ts | `resolveTaxRate` (new) | — | 1 |
| src/orders/discount-calculator.ts | `tierDiscount` (new) | — | 1 |
| src/orders/discount-calculator.ts | `capDiscount` (new) | — | 2 |
| src/orders/promo-strategy.ts | `promoDiscount` (new) | — | 2 |

## Pre-Refactor State

`OrderProcessor` is a 100-line class holding order validation, subtotal arithmetic, tier and
promo-code discounting, tax-rate selection, a payment-method business rule, database
persistence, and customer email notification. That is at least five reasons to change in one
class — a God Object smell on top of the two measured complexity violations. `applyPromotions`
is additionally a public method mixing two unrelated discount sources (customer tier, promo
code) behind a single `discount` accumulator, so neither source can be tested or changed in
isolation.

**Escalation — no test coverage (Phase 0 step 2)**: a glob for `*spec*`, `*test*`,
`*.test.*`, and `*_test.*` under `src/orders/` returned nothing. The target file has zero
covering tests, and the guardrail is explicit: never refactor a file with zero test coverage.
Refactoring was **halted** and `unit-tester` was invoked with `src/orders/order-processor.ts`
as its target before any edit was made. `unit-tester` produced
`src/orders/order-processor.spec.ts` — 14 characterization tests covering: both validation
rejections, the zero-item rejection, invalid-quantity rejection, subtotal arithmetic across
multiple line items, all four customer tiers, all three promo codes (including the VIP30
platinum and non-platinum branches), the 0.35 discount cap, both tax rates, the $10,000
PayPal rejection, and the success path's persistence + notification calls. Only after that
suite was green did Phase 1 begin.

**context-debt: context-engineer not run** — this run was invoked standalone against a single
file with no active feature workspace, so no `context-manifest.md` exists. No ADR or KI
constraints on preferred patterns for the Ordering context were consulted; if this file is
revisited inside a pipeline run, re-check for a module-level pattern ADR before extending the
structure introduced here.

**Complexity baseline** (from `npx eslint --rule 'complexity: ["error", 7]' src/orders/order-processor.ts`):
```
/src/orders/order-processor.ts
  35:3  error  Method 'processOrder' has a complexity of 9. Maximum allowed is 7      complexity
  72:3  error  Method 'applyPromotions' has a complexity of 8. Maximum allowed is 7   complexity

✖ 2 problems (2 errors, 0 warnings)
```

**Test baseline**:
- Command: `npx vitest run src/orders`
- Result (before `unit-tester`): 0 passed, 0 failed, 0 skipped — **no test files found**
- Result (after `unit-tester`, before first refactoring edit): 14 passed, 0 failed, 0 skipped

## Operations Applied

Applied one at a time, smallest scope first, with `npx tsc --noEmit` after each extraction and
the full suite re-run after each numbered step.

1. **Extract Function** — `processOrder` → extracted the three rejection guards (missing
   ids, empty items, invalid quantity) into `validateOrder(order): ValidationFailure | null`
   in the new `src/orders/order-validator.ts` (26 lines → 6 lines in the caller). This removed
   4 decision points from `processOrder`; the caller now has a single early-return guard
   clause against the returned failure.
2. **Extract Function** — `processOrder` → extracted `calculateSubtotal(items): number`,
   which now loops purely for arithmetic because the per-item quantity check moved into
   `validateOrder` in step 1.
3. **Extract Function** — `processOrder` → extracted `resolveTaxRate(method): number`,
   replacing the inline `if (order.paymentMethod === 'bank_transfer')` reassignment with a
   `Record<PaymentMethod, number>` lookup keyed on the payment method.
4. **Extract Variable** — `processOrder` and `applyPromotions` → hoisted the magic numbers
   `10000` and `0.35` into the named constants `PAYPAL_ORDER_LIMIT` and `MAX_TOTAL_DISCOUNT`.
   No arithmetic changed; the literals were replaced by references to their own values.
5. **Introduce Parameter Object** — `applyPromotions(subtotal, tier, promoCode?)` →
   grouped the data clump `subtotal, tier, promoCode` into `PromotionContext`, which now
   travels as one argument to the discount collaborators introduced in step 6 instead of
   being threaded through three separate parameter lists.
6. **Replace Conditional with Polymorphism** — `applyPromotions` → the `SUMMER10` /
   `FLASH20` / `VIP30` if-else chain was replaced by a `PromoStrategy` map in the new
   `src/orders/promo-strategy.ts`, each strategy a `(context: PromotionContext) => number`.
   This step was delegated to `refactor-to-pattern` (GoF **Strategy**), which produced the
   approved plan `.claude/refactor-plans/promo-strategy-plan.md`; the tier ladder was moved
   alongside it as the `tierDiscount` lookup, and `capDiscount` retained the ceiling. An
   unknown promo code resolves to the identity strategy returning `0`, preserving the
   original chain's fall-through of adding nothing.

## Post-Refactor State

Both originally violating methods are now well under the threshold, and every function
extracted from them is at 3 or below. No exceptions were needed — neither method contained an
inherently complex algorithm, only tangled concerns.

**Complexity after refactoring**:
```
$ npx eslint --rule 'complexity: ["error", 7]' src/orders/

✖ 0 problems (0 errors, 0 warnings)

$ npx eslint --rule 'complexity: ["warn", 1]' src/orders/ --format compact | grep complexity
src/orders/order-processor.ts: line 41, col 3, Warning - Method 'processOrder' has a complexity of 4.
src/orders/order-processor.ts: line 63, col 3, Warning - Method 'applyPromotions' has a complexity of 3.
src/orders/order-processor.ts: line 71, col 3, Warning - Method 'calculateSubtotal' has a complexity of 2.
src/orders/order-processor.ts: line 78, col 3, Warning - Method 'resolveTaxRate' has a complexity of 1.
src/orders/order-validator.ts:  line 12, col 1, Warning - Function 'validateOrder' has a complexity of 3.
src/orders/discount-calculator.ts: line 9, col 1, Warning - Function 'tierDiscount' has a complexity of 1.
src/orders/discount-calculator.ts: line 18, col 1, Warning - Function 'capDiscount' has a complexity of 2.
src/orders/promo-strategy.ts: line 21, col 1, Warning - Function 'promoDiscount' has a complexity of 2.
```

`processOrder` went from 9 to 4 (below 7), `applyPromotions` from 8 to 3 (below 7).
`OrderProcessor` itself dropped from 100 lines to 58, back inside the Sandi Metz class limit.

**Exceptions** (functions that could not reach < 7 without adding behavior):
- None. Every target function reached a complexity of 4 or lower through structural
  transformation alone.

## Behavior Preservation Evidence

The characterization suite `unit-tester` wrote against the *original* file was not modified
during refactoring — only import paths were left untouched because the public entry point
`OrderProcessor.processOrder` and the public `applyPromotions` signature were preserved. The
same 14 assertions that pinned the pre-refactor behavior pass unchanged against the
post-refactor code.

- Command: `npx vitest run src/orders`
- Result: 14 passed, 0 failed, 0 skipped — **all tests green**
- Tests added by this run: 14 characterization tests added by `unit-tester` in
  `src/orders/order-processor.spec.ts` (the file had no prior coverage). No test was added,
  changed, or deleted after the baseline was recorded.
- Type check: `npx tsc --noEmit` — 0 errors.

## Bugs Observed (Not Fixed)

Both were pinned by the characterization suite exactly as they behave today, so the tests
document the current (wrong) behavior rather than the intended behavior. Each must be fixed in
a separate commit with its own test change.

- `src/orders/order-processor.ts` line 88 (`applyPromotions`, VIP30 branch): the `VIP30` promo
  code silently grants a **different discount depending on customer tier** — +0.30 for
  platinum, +0.15 for everyone else — with no documentation, no separate promo code, and no
  signal to the caller. A customer redeeming "VIP30" on a gold account receives half the
  advertised discount. This is a probable product defect, not a structural one; the Strategy
  extracted in operation 6 reproduces it verbatim. Deferred to a separate fix commit.
- `src/orders/order-processor.ts` lines 52–60 (`processOrder`): monetary amounts are computed
  in IEEE-754 floating point (`subtotal += item.quantity * item.unitPrice`, then
  `discountedTotal * taxRate`) and only rounded at the string-formatting step in the
  notification email. The persisted `total` retains the unrounded float, so the saved order
  total and the emailed total can disagree by a fraction of a cent. Fixing this requires an
  integer-cents or decimal representation — a behavioral change. Deferred to a separate fix
  commit.

## No-Behavior-Added Attestation

No new behavior was added during this refactoring run. Every changed line is a structural
transformation only — no new conditional branches, no new business logic, no new observable
outputs. The unknown-promo-code path, the discount ceiling, both tax rates, the $10,000 PayPal
rejection, and the tier-dependent VIP30 divergence all behave exactly as they did before. The
two bugs observed above were noted and left for a dedicated fix commit.
