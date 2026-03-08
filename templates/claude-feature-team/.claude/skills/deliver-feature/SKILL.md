---
name: deliver-feature
description: Orchestrates the full feature delivery pipeline. Runs analyst → developer → qa-engineer → tech-writer → devops-engineer in sequence. Invoke with /deliver-feature path/to/feature.md
triggers:
  keywords: ["/deliver-feature", "deliver feature", "build this feature", "implement feature in"]
  intentPatterns: ["run the feature pipeline", "deliver feature {path}"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill to initiate the fully automated feature delivery pipeline. It starts the analyst, goes through development, QA, doc updates, and CI infrastructure.
Do NOT use if the user just wants a specific subtask executed unless they explicitly ask to deliver the whole feature.

## Context To Load First
1. Provided `FEATURE.md` argument or description.
2. `README.md` (to understand agent orchestration pipeline).
3. `ARCHITECTURE_RULES.md`

## Process
1. Setup Workspace.
   - What to do: Create `.claude/feature-workspace/` directory if it doesn't exist. If the argument is an inline description, save it to `features/[kebab-case-name].md`.
   - What to produce: The `.claude/feature-workspace/` directory.
   - What to check: Verify the directory exists.
2. Analysis Phase
   - What to do: Use the `analyst` subagent to analyze the feature at `$ARGUMENTS`. Wait for it to produce `.claude/feature-workspace/analysis.md`.
   - What to produce: `.claude/feature-workspace/analysis.md`
   - When to pause: **STOP HERE for approval.** Show the user a summary of the analysis and ask: 
     > "The analyst has completed the technical breakdown. Here's what's planned: [summary]. Does this look correct? Should I proceed with implementation?"
3. Implementation Phase
   - What to do: Use the `developer` subagent to implement the feature by reading `analysis.md`.
   - What to produce: Code updates and `implementation-notes.md`.
   - What to check: Wait until implementation is complete before moving to QA.
4. Testing Phase
   - What to do: Use the `qa-engineer` subagent to write and run tests based on `analysis.md` and `implementation-notes.md`.
   - What to produce: Test files and `qa-report.md`.
   - When to pause: If QA reports test failures that couldn't be resolved, pause and report to the user before continuing.
5. Documentation Phase
   - What to do: Use the `tech-writer` subagent to update all documentation.
   - What to produce: Doc updates and `docs-report.md`.
6. CI/CD Phase
   - What to do: Use the `devops-engineer` subagent to handle CI/CD.
   - What to produce: CI pipeline changes and `devops-report.md`.

## Output Format
Produce a Delivery Summary in `.claude/feature-workspace/delivery-summary.md` and display it to the user:

```markdown
# Delivery Summary: [Feature Name]

**Status**: ✅ Complete / ⚠️ Complete with notes / ❌ Blocked

## What Was Built
[2-3 sentence plain English summary]

## Files Changed
[Count of files created/modified by category: source, tests, docs, CI]

## Acceptance Criteria
- [x] Criterion 1 — verified by [test name]

## Manual Steps Required
[Any steps the human must take: secrets to set, migrations to run, etc.]

## Known Issues / Follow-ups
[Anything that was deferred or needs attention]
```

## Guardrails
- You MUST wait for human approval after Step 2 (Analysis Phase). Do not proceed to the Developer until approved.
- Any unresolved test failures in Step 4 must pause the pipeline for human insight.

## Standalone Mode
If subagents are unavailable, degrade gracefully: prompt the user for manual steps in sequence (Analysis -> Implementation -> QA), outputting the necessary artifacts manually.
