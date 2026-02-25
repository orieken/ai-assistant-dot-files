---
name: deliver-feature
description: Orchestrates the full feature delivery pipeline. Runs analyst → developer → qa-engineer → tech-writer → devops-engineer in sequence. Invoke with /deliver-feature path/to/feature.md
---

# Feature Delivery Pipeline

You are orchestrating a full feature delivery. The feature spec is: $ARGUMENTS

Follow this pipeline exactly:

## Step 1: Setup
Create `.claude/feature-workspace/` directory if it doesn't exist.
If `$ARGUMENTS` is not a file path, save the inline description to `features/[kebab-case-name].md` first.

## Step 2: Analysis (PAUSE FOR APPROVAL)
Use the `analyst` subagent to analyze the feature at `$ARGUMENTS`.
Wait for it to produce `.claude/feature-workspace/analysis.md`.

**STOP HERE.** Show the user a summary of the analysis and ask:
> "The analyst has completed the technical breakdown. Here's what's planned:
> [summarize key points from analysis.md]
> 
> Does this look correct? Should I proceed with implementation?"

Do not continue until the user confirms.

## Step 3: Implementation
Use the `developer` subagent to implement the feature.
It should read `analysis.md` and produce code + `implementation-notes.md`.

## Step 4: Testing
Use the `qa-engineer` subagent to write and run tests.
It should read `analysis.md` and `implementation-notes.md` and produce tests + `qa-report.md`.

If QA reports test failures that couldn't be resolved, pause and report to the user before continuing.

## Step 5: Documentation
Use the `tech-writer` subagent to update all documentation.
It should read all previous outputs and produce doc updates + `docs-report.md`.

## Step 6: CI/CD
Use the `devops-engineer` subagent to handle CI/CD and infrastructure.
It should read all previous outputs and produce CI changes + `devops-report.md`.

## Step 7: Delivery Summary
Write `.claude/feature-workspace/delivery-summary.md` with this format:

```markdown
# Delivery Summary: [Feature Name]

**Status**: ✅ Complete / ⚠️ Complete with notes / ❌ Blocked

## What Was Built
[2-3 sentence plain English summary]

## Files Changed
[Count of files created/modified by category: source, tests, docs, CI]

## Acceptance Criteria
- [x] Criterion 1 — verified by [test name]
- [x] Criterion 2 — verified by [test name]

## Manual Steps Required
[Any steps the human must take: secrets to set, migrations to run, etc.]

## Known Issues / Follow-ups
[Anything that was deferred or needs attention]
```

Then present this summary to the user.
