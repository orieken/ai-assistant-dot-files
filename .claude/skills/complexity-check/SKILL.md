---
name: complexity-check
description: On-demand complexity analysis of any file or directory.
triggers:
  keywords: ["complexity", "technical debt", "refactoring"]
  intentPatterns: ["Check complexity of *", "Is * too complex?", "Complexity report on *", "What's the most complex file in *?"]
standalone: true
---

## When To Use
Triggered when the user asks for a complexity analysis of a file, directory, or package, or mentions "technical debt" without a specific task.

## Context To Load First
1. `ARCHITECTURE_RULES.md`
2. The target file(s)

## Process
1. Determine the language(s) of the target
2. Run the appropriate complexity tool
   - TypeScript: `npx eslint --rule '{"complexity": ["error", 6]}' [file]`
   - Python: `radon cc [file] -s` or `flake8 --max-complexity=6 [file]`
   - Go: `gocyclo -over 6 [file]`
   - Java: manual heuristics only — recommend Checkstyle/SonarQube to user
3. For each function exceeding complexity 6: name the smell, name the Fowler operation, estimate effort (trivial / one session / needs design discussion)
4. Rank findings by impact: highest complexity first, then by call graph depth
5. Produce Complexity Report

## Output Format
`.claude/feature-workspace/complexity-report.md`

```markdown
# Complexity Report: [target]

Threshold: Cyclomatic complexity > 6 (ARCHITECTURE_RULES.md)

## Summary
- Files scanned: N
- Functions exceeding threshold: N
- Highest complexity: [function name] at [N]

## Findings (ranked by impact)

### [FunctionName] — complexity [N]
**File**: `path/to/file.ts` line X
**Smell**: [specific description]
**Refactoring**: [named Fowler operation]
**Effort**: trivial / one session / needs design discussion

## Refactoring Roadmap
1. [First because: quick win / unblocks other refactors / highest risk]

## What's Clean
[Functions within threshold — always acknowledge good work]
```

## Guardrails
- Read-only. Never modify files. Never run tests. 
- Fall back to heuristic analysis if CLI tools unavailable — still produce a report, note it's heuristic.

## Standalone Mode
Runs external tools via the terminal or uses heuristic analysis via string processing / code reading.
