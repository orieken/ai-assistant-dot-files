<!--
Template for refactoring-notes.md. Consumed by the refactor-engineer agent.
Contract in shared/contracts/refactoring-contract.md validates required headings.
Preserve every heading exactly. Write "None" for inapplicable sections — never delete headings.
-->

# Refactoring Notes: [Target / Feature Name]

## Summary

[One paragraph: what was refactored and why (complexity violation / Boy Scout Rule / migration
sprint), the Fowler operations applied, and the net outcome (complexity delta, files changed).]

## Target Files

| File | Function / Method | Pre-Refactor Complexity | Post-Refactor Complexity |
|---|---|---|---|
| [path/to/file.ts] | [functionName] | [N] | [N] |

## Pre-Refactor State

[Paste or quote the complexity tool output for each target function before any changes.
Include the test-suite baseline result (pass/fail counts) before the first edit.]

**Complexity baseline** (from `[tool name, e.g. eslint / gocyclo / radon]`):
```
[tool output]
```

**Test baseline**:
- Command: `[test command]`
- Result: [N passed, N failed, N skipped]

## Operations Applied

[Ordered list. Each entry: Fowler operation name (exact, from shared/rules/design-principles.md §2),
the target function, and one sentence describing the transformation. If refactor-to-pattern was
called as a sub-step, name the GoF pattern and the plan file it produced.]

1. **[Extract Function]** — `[functionName]` → extracted `[newFunctionName]` ([N lines → N lines each)
2. **[Replace Conditional with Polymorphism]** — `[functionName]` → introduced `[StrategyName]` via `refactor-to-pattern`
3. **[Introduce Parameter Object]** — `[functionName]` → grouped `[param1, param2, param3]` into `[ObjectName]`

## Post-Refactor State

[Paste or quote the complexity tool output after all operations are complete.
Every originally violating function must now show < 7, or carry a documented exception below.]

**Complexity after refactoring**:
```
[tool output]
```

**Exceptions** (functions that could not reach < 7 without adding behavior):
- [functionName]: complexity [N] — [one-sentence rationale why structural reduction alone cannot
  reach < 7, e.g. "inherently complex algorithm; decomposition would require behavioral changes"]

## Behavior Preservation Evidence

[Test suite result after all refactoring operations are complete. Must state pass/fail counts.]

- Command: `[test command]`
- Result: [N passed, 0 failed, N skipped] — **all tests green**
- Tests added by this run: [N characterization tests added by unit-tester, or "None — existing
  coverage was sufficient"]

## Bugs Observed (Not Fixed)

[Bugs spotted during refactoring that were intentionally left untouched. Each entry: file, line,
description, and a note that it must be fixed in a separate commit. Write "None" if none found.]

- `[path/to/file.ts]` line [N]: [description] — deferred to separate fix commit.

## No-Behavior-Added Attestation

No new behavior was added during this refactoring run. Every changed line is a structural
transformation only — no new conditional branches, no new business logic, no new observable
outputs. Any bugs observed above were noted and left for a dedicated fix commit.
