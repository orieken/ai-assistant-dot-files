---
name: design-review
description: Standalone design critique of any file, not tied to feature delivery.
triggers:
  keywords: ["review", "design", "smells", "fowler", "critique"]
  intentPatterns: ["Review the design of *", "Is this well-designed?", "Check * for design smells", "Would Fowler approve of this?"]
standalone: true
---

## When To Use
Triggered when the user asks for a design review, code critique, or smell checking of a specific file or directory outside the context of the standard feature delivery pipeline. Do NOT use if the user just asks to "review my code" during feature implementation; that's the code-reviewer's job.

## Context To Load First
1. `ARCHITECTURE_RULES.md`
2. `DOMAIN_DICTIONARY.md`
3. The target file(s)
4. Any files the target imports from or that import it (one level of dependency graph).

## Process
1. Write a Design Narrative (what is this thing trying to do?)
2. Map the dependency graph for the target (what does it depend on? what depends on it?)
3. Check all Fowler smells from ARCHITECTURE_RULES.md
4. Check Clean Architecture layer placement
5. Check Simple Design rules (Beck's 4 rules) in priority order
6. Check Sandi Metz constraints (class/method size, param count)
7. Identify the single most impactful refactoring to apply first
8. Produce a Design Report

## Output Format
`.claude/feature-workspace/design-report.md`

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
Read-only. Never modify source files. Never run tests. Surface findings only.

## Standalone Mode
Fully conversational. No external tools needed beyond file reading.
