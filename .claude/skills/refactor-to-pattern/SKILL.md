---
name: refactor-to-pattern
description: A surgical tool to transition messy procedural code into a specific GoF or Enterprise Integration pattern.
triggers:
  keywords: ["refactor", "pattern", "GoF", "clean code", "strategy", "factory", "decorator"]
  intentPatterns: ["Refactor this to use the * pattern", "Rewrite * to be a * pattern", "Clean up this messy code into a * pattern"]
standalone: true
---

## When To Use
Triggered when a developer specifically requests a codebase/file to be transformed into a documented design pattern, or when code-reviewer specifically recommends a pattern refactoring.

## Context To Load First
1. The target file(s).
2. `ARCHITECTURE_RULES.md`
3. Any unit tests covering the target file.

## Process
1. Analyze the original code: identify the "smell" (e.g., giant switch statement, duplicated behavior, hardcoded dependencies).
2. Validate the requested pattern: Does the requested pattern actually solve the smell? If no, recommend a better pattern. If yes, proceed.
3. Outline the transformation steps (e.g., 1. Extract interface, 2. Create concrete classes, 3. Inject context).
4. Present the refactoring plan to the user for approval.
5. Once approved, rewrite the code specifically adhering to the pattern constraints, preserving exact business logic.
6. Verify test safety: What tests need to be updated to support the new structure?

## Output Format
`.claude/feature-workspace/refactor-[pattern]-plan.md`

```markdown
# Refactoring to [Pattern Name]

## The Smell Addressed
[E.g. "Giant switch statement violating Open-Closed Principle"]

## Transformation Plan
1. [Extract new interface X]
2. [Create concrete implementations Y, Z]
3. [Implement the client using Dependency Injection]

## Test Impact
- [What tests will break due to structure changes, though behavior remains identical]
```

## Guardrails
- **Never** change business logic during a structural refactoring. The output must be functionally identical.
- Ensure the refactored code passes Sandi Metz's complexity limits.
