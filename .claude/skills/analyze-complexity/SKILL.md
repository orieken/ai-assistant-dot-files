---
name: analyze-complexity
description: Analyzes a file or directory for cyclomatic complexity and function length. Fails if complexity > 7 or function length > 30 LOC according to clean code rules.
triggers:
  keywords: ["complexity", "refactor", "too long", "cyclomatic", "mccabe", "clean up"]
  intentPatterns: ["check code complexity", "enforce clean code rules", "is this function too long"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill when reviewing code, checking pull request changes, or proactively cleaning up technical debt.
Specifically apply it to any file modified during feature implementation.
Do NOT use when checking data files, documentation, or configuration files (JSON, YAML, MD).

## Context To Load First
1. `ARCHITECTURE_RULES.md`
2. `.claude/feature-workspace/implementation-notes.md` (if available)

## Process
1. Determine the source files or directories you want to check.
   - What to do: Identify modified or relevant code files (e.g., `.ts`, `.js`, `.py`, `.go`).
   - What to produce: A list of targets.
   - What to check: Verify the files exist.
2. Run the complexity check on the targets.
   - What to do: Run `./.claude/skills/analyze-complexity/check.sh [target]`.
   - What to produce: The console output of the script.
   - What to check: Check the exit code. If non-zero, complexity rules are violated.
3. Handle violations.
   - What to do: If the script fails, instruct the developer to refactor the code (Extract Method, Inline, etc.) to meet the `< 7` complexity and `< 30` LOC standards.
   - When to pause: If complex refactoring alters the architecture significantly, ask for approval before applying.

## Output Format
If errors are found, output a markdown list of violations:
```markdown
# Complexity Report

## Violations
- `path/to/file.ts`: Function `doComplexThing` complexity is 8 (Limit: 7)
- `path/to/file.ts`: Function `longMethod` length is 45 lines (Limit: 30)

## Recommended Refactorings
- [Suggested Extract Method or Strategy Pattern application]
```

## Guardrails
- You MUST NOT approve code that violates the complexity limit without explicit human approval.
- Do NOT silently ignore failures from `check.sh`.

## Standalone Mode
If `check.sh` is unavailable or fails to execute, manually review the code using static analysis: Count the `if`, `else`, `for`, `while`, `switch` statements to estimate cyclomatic complexity. Count lines of code manually. 
