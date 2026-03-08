---
name: complexity-check
description: A standalone skill for on-demand complexity analysis of any file or directory, identifying technical debt and refactoring targets.
triggers:
  keywords: ["complexity", "technical debt", "most complex file", "cyclomatic"]
  intentPatterns: ["check complexity of {target}", "is {target} too complex?", "complexity report on {target}", "what's the most complex file in {package}?"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill when the user asks to evaluate the complexity of a specific file, module, or package. Use it when they ask about technical debt without a specific feature task attached.
Do NOT use this during standard feature development unless explicitly asked (that is the `code-reviewer` agent's job).

## Context To Load First
1. The target file(s) or directory specified by the user.
2. `ARCHITECTURE_RULES.md` Section IV (to confirm the exact threshold: McCabe complexity ceiling < 7).

## Process
1. **Determine the language(s)**: Check the file extensions of the target.
2. **Run the appropriate complexity tool**:
   - What to do: Execute one of the following commands based on language:
     - TypeScript/JS: `npx eslint --rule '{"complexity": ["error", 6]}' [file]`
     - Python: `radon cc [file] -s` or `flake8 --max-complexity=6 [file]`
     - Go: `gocyclo -over 6 [file]`
     - Java/Other: If no CLI tool is available, explicitly state that Checkstyle/SonarQube should be used in CI, and fall back to manual heuristic analysis.
   - What to produce: Tool output or heuristic findings.
3. **Analyze findings**:
   - What to do: For each function exceeding complexity 6, identify:
     - The specific smell (e.g., nested conditionals, long switch, multiple responsibilities).
     - The exact Fowler refactoring operation to apply (e.g., Extract Function, Replace Conditionals with Polymorphism).
     - Effort estimate: (trivial / one session / needs design discussion).
4. **Rank findings**: Rank them by highest complexity score first, then by depth in the call graph (how many other functions depend on it).
5. **Produce Complexity Report**: Output the report.

## Output Format
Create `.claude/feature-workspace/complexity-report.md` (and show it to the user):

```markdown
# Complexity Report: [target]

Date: [today]
Threshold: Cyclomatic complexity > 6 (ARCHITECTURE_RULES.md Section IV)

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
Suggested order to attack the debt:
1. [First because: quick win / unblocks other refactors / highest risk]
2. ...

## What's Clean
[Functions that are well within threshold — always acknowledge good work]
```

## Guardrails
- **Read-only**: You MUST NEVER modify source files while running this skill.
- **No Test Execution**: You MUST NEVER run the test suite as part of this check.
- **Clear Caveats**: If a complexity tool is not installed or available, you MUST explicitly state that the scan is using manual heuristics and is not a strictly measured score.

## Standalone Mode
If the specified CLI tools (e.g., eslint, radon, gocyclo) are unavailable, this skill degrades gracefully. You will manually read the source files, estimate cyclomatic complexity by counting branching statements (`if`, `else`, `switch`, `for`, `while`, `catch`), and produce the report with a clear "Heuristic Mode" disclaimer.
