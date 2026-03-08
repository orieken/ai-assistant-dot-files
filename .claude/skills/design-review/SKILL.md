---
name: design-review
description: A standalone skill for reviewing the design of any code, completely detached from feature delivery.
triggers:
  keywords: ["design review", "well-designed", "design smells", "fowler"]
  intentPatterns: ["review the design of {target}", "is this well-designed", "check {target} for design smells", "would fowler approve of this"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill when asked to evaluate the design, architecture, or craftsmanship of a specific file, class, or module outside of standard feature implementation. 
Do NOT use when asked to run tests, write implementations, or deploy code.

## Context To Load First
1. `ARCHITECTURE_RULES.md`
2. `DOMAIN_DICTIONARY.md`
3. The target file(s) specified by the user.
4. Any files the target imports from or that import it (explore one level of the dependency graph).

## Process
1. **Write a Design Narrative**: Synthesize a 2-3 sentence plain-English description of what the target is trying to do. Assess if it does only one thing.
2. **Map the dependency graph**: Identify what the target depends on and what depends on it (just one level deep).
3. **Check all 10 Fowler smells**: Explicitly check the code against the 10 Refactoring Catalog operations in `ARCHITECTURE_RULES.md` Section VII (e.g., Extract Function, Replace Conditionals).
4. **Check Clean Architecture layer placement**: Verify if the component belongs in Entities, Use Cases, Adapters, or Frameworks, and check if its dependencies point inward appropriately (Section I).
5. **Check Simple Design rules**: Apply Kent Beck's 4 rules (Section VI) strictly in priority order: Passes tests, Reveals intention, No duplication, Fewest elements.
6. **Check Sandi Metz constraints**: Evaluate class size ($\le$ 100 lines), method size ($\le$ 5 lines), and method parameters ($\le$ 4).
7. **Identify the single most impactful refactoring**: Decide which named refactoring operation would provide the most value if applied first.
8. **Produce a Design Report**: Output the findings using the exact structure below.

## Output Format
Create `.claude/feature-workspace/design-report.md` (and show it to the user):

```markdown
# Design Review: [ClassName or filename]

## Design Narrative
[2-3 sentences: what is this thing, what is its job, does it do one thing?]

## Dependency Graph
- Depends on: [list]
- Depended on by: [list]
- Layer: [which Clean Architecture layer]
- Direction correct: yes/no

## Design Smells Found
| Smell | Location | Severity | Fowler Operation |
|---|---|---|---|

## Simple Design Score
1. Passes tests: ✅/❌
2. Reveals intention: ✅/⚠️ [what to rename]
3. No duplication: ✅/⚠️ [what to extract]
4. Fewest elements: ✅/⚠️ [what to remove]

## Recommended Refactoring (Priority Order)
1. [Most impactful named operation] — [why, what it unlocks]
2. ...

## What's Working Well
[Always include — good design deserves recognition]
```

## Guardrails
- **Read-only**: You MUST NEVER modify source files while running this skill.
- **No execution**: Never run tests or execute code.
- **Surface findings only**: You are an observer and a critic in this mode. Do not attempt to fix the problems you find unless the user explicitly asks you to in a completely separate prompt.

## Standalone Mode
Works entirely without external systems or MCP tools. Uses standard local file reading to parse the source code and static dependency analysis.
