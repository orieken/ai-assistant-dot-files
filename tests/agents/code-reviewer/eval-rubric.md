# Eval Rubric: code-reviewer / input-smelly-code.ts

- **Correct diagnosis of the infrastructure leak**: identifies that `Invoice` (a domain class) directly
  imports and queries `pg`'s `Pool` as a Clean Architecture / dependency-inversion violation — names *why*
  it's wrong (domain depending on infrastructure), not just "this uses pg directly."
  - Fix references [architecture-guardrails.md](../../../shared/rules/architecture-guardrails.md) rule #1
    (dependency direction) in spirit, even if not by filename.
- **Single Responsibility named explicitly**: calls out that `processAndPrint` both calculates a total AND
  logs/prints it as two separate reasons to change — not just "this method does a lot" without naming the
  two responsibilities.
- **Named refactoring operations used**: uses Fowler's actual vocabulary (e.g., "Extract Function",
  "Replace Conditional with Polymorphism" or "Replace Conditional with Lookup Table" for the if/else-if
  pricing chain) rather than vague instructions like "clean this up" or "make this better."
- **Magic numbers named specifically**: flags `1.08` and `0.9` by value, not just "avoid magic numbers" in
  the abstract.
- **Verdict matches severity**: overall verdict is CHANGES REQUESTED (not APPROVED or APPROVED WITH
  COMMENTS) given the infrastructure-leak violation is a real architectural issue, not a style nit.

## How to Grade
For each bullet above, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no
supporting quote, mark it FAIL and say what's missing.
