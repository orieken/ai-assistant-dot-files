---
name: ast-grep-structural-search
tags: [context-engineering, structural-search, refactoring, cli-tools, optional-dependency]
domain: context-engineering
created: 2026-08-21
---

ast-grep (`github.com/ast-grep/ast-grep`) is an AST-level structural search and rewrite tool.
It finds code by shape, not text — `$A.findById($B)` matches all callers of `findById`
regardless of variable names. This eliminates false positives (text matches on unrelated
strings) and false negatives (text misses because of aliasing or formatting) that plain
`grep` produces.

Supports TypeScript, JavaScript, Go, Python, Rust, Java, C#, and others natively.

## Integration point — context-engineer Step 1 (Identify Target Component Scope)

When scoping a bounded context, use ast-grep to find all files that implement or call a
specific interface, rather than text-searching for the interface name:

```bash
# Find all implementations of an interface (TypeScript)
sg run -p 'implements $IFACE' --lang ts

# Find all call sites of a method (Go)
sg run -p '$RECV.MethodName($$$ARGS)' --lang go

# Find all files that import a specific package
sg run -p 'import "$PKG"' --lang ts -r
```

Add discovered files to the candidate file list before applying churn and budget filters.

If ast-grep is not installed, fall back to `grep -r` and note:
```
context-debt: ast-grep not installed — file discovery used text search (may have false positives)
```

## Integration point — refactor-engineer Phase 0

Before a structural refactoring campaign, use ast-grep to produce an exhaustive call-site
inventory — every location that will need updating:

```bash
# All usages of the old pattern
sg run -p '<old-pattern>' --lang <lang> --json | jq '.[] | .file'
```

This gives refactor-engineer an accurate scope baseline rather than a grep-approximated one.

## Rewrite mode (refactor-engineer Phase 3+)

ast-grep can also rewrite matches in place — useful for mechanical rename/extract operations:

```bash
sg run -p '$OBJ.oldMethod($$$ARGS)' \
       -r '$OBJ.newMethod($$$ARGS)' \
       --lang ts --update-all
```

Always run rewrites on a clean git branch; verify with the test suite before committing.

## Installation

```bash
brew install ast-grep
# or
cargo install ast-grep --locked
```

## Guardrails

- Patterns are language-specific. A Go pattern will not match TypeScript even if the syntax
  looks identical. Always pass `--lang`.
- `$$$` matches zero or more nodes (variadic). `$$` matches one or more. `$` matches exactly
  one. Mixing these up produces silent mismatches.
- ast-grep does not resolve imports — it finds structural matches in the files it scans, not
  cross-file symbol resolution. For "who implements interface X across the monorepo," you still
  need to provide the right directory scope.
- Do not use ast-grep rewrites (`--update-all`) without characterization tests in place first.
  Structural rewrite ≠ semantic equivalence — the test suite is the verification.

## See also

- `shared/skills/context-engineer/SKILL.md` — Step 1 (Identify Target Component Scope)
- `shared/agents/refactor-engineer.md` — Phase 0 scope baseline and Phase 3+ rewrite operations
- `shared/knowledge/git-churn-risk-signal.md` — pairs for Step 4b: churn filters the
  structural candidate list by risk priority
