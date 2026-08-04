# Contract: refactoring-notes.md

**Produced by**: refactor-engineer
**Consumed by**: code-reviewer (light-touch post-refactor review), modernization-supervisor
(when invoked as part of a swarm), orchestrator (deliver-feature if wired as an optional step)

## Required Sections (exact heading text and level)

- `## Summary`
- `## Target Files`
- `## Pre-Refactor State`
- `## Operations Applied`
- `## Post-Refactor State`
- `## Behavior Preservation Evidence`
- `## Bugs Observed (Not Fixed)`
- `## No-Behavior-Added Attestation`

## Validation Rules

`validate-artifact` checks presence of every heading above, plus:

- `## Operations Applied` must contain at least one named Fowler operation from this set:
  `Extract Function`, `Inline Function`, `Extract Variable`, `Rename Variable`, `Move Method`,
  `Move Field`, `Replace Conditional with Polymorphism`, `Introduce Parameter Object`,
  `Remove Dead Code`, `Separate Query from Modifier`, `Preserve Whole Object`. Free-text
  operation names are a FAIL — the operation must be named exactly.

- `## Pre-Refactor State` must contain at least one complexity measurement (a number ≥ 7, or
  the word "complexity"). A state section that only describes the code qualitatively without
  measuring is a FAIL.

- `## Behavior Preservation Evidence` must contain the word `pass` or `green` (case-insensitive)
  to indicate the test suite result after refactoring. A section that only says "tests were run"
  without a pass confirmation is a FAIL.

- `## No-Behavior-Added Attestation` must contain either the exact phrase `no behavior added`
  or `no new behavior` (case-insensitive). A missing or vague attestation is a FAIL.

## Section Semantics

### `## Summary`
One paragraph: what was refactored, why (complexity violations / Boy Scout Rule / migration),
and the outcome (complexity delta, file count changed).

### `## Target Files`
List of files touched with their pre- and post-refactor complexity for each target function.

### `## Pre-Refactor State`
Measured cyclomatic complexity per function, Sandi Metz violations, and the baseline test
result (pass count / fail count before any changes). Quote the complexity tool output directly.

### `## Operations Applied`
Ordered list of operations performed: Fowler operation name, target function, and one sentence
describing what was moved/extracted/renamed. If `refactor-to-pattern` was called as a sub-step,
name the GoF pattern and the plan file it produced.

### `## Post-Refactor State`
Re-measured complexity per target function. Every originally violating function must now be
< 7, or must carry a documented "judgment-only" exception.

### `## Behavior Preservation Evidence`
Test suite result after refactoring: pass/fail counts, any new test failures (which are a
signal of regression, not expected). Must state the test command used.

### `## Bugs Observed (Not Fixed)`
Bugs spotted during refactoring that were intentionally left alone (to be fixed in a separate
commit). "None" is valid.

### `## No-Behavior-Added Attestation`
Explicit statement that no new logic, conditional branches, or observable behavior was
introduced. Required even when the answer is obvious.

## Retrieval Frontmatter (WARN)

Pipeline artifacts should include a YAML frontmatter block at the very top of the file. Missing or incomplete retrieval frontmatter triggers a **WARN** from `validate-artifact` — not a FAIL. Existing archived artifacts without frontmatter are unaffected.

```yaml
---
feature: "<feature-name>"             # kebab-case slug derived from the feature file name
bounded_context: "<context>"          # owning bounded context (from DOMAIN_DICTIONARY.md domain list)
domain_terms: []                      # canonical terms from DOMAIN_DICTIONARY.md used in this feature
files_touched: []                     # repo-relative paths of files created or modified
issue_refs: []                        # ticket/issue references (e.g., PROJ-123, #456)
linked_adrs: []                       # repo-root-relative paths to referenced ADRs
linked_kis: []                        # repo-root-relative paths to referenced Knowledge Items
---
```

Once frontmatter adoption is visible across a project's feature archive, this check will be promoted to FAIL in a future release.
