---
name: review-pr
description: Runs the code-reviewer and security-reviewer agents against a git diff or branch — useful for hotfixes, dependency updates, and changes that didn't go through the full delivery pipeline. Invoke with /review-pr or /review-pr [branch-name] or /review-pr [base-branch] [head-branch]
---

# PR / Diff Review

You are running a focused code and security review outside the full delivery pipeline.
Input: $ARGUMENTS

## Step 1 — Determine the diff scope

Parse `$ARGUMENTS`:

- **No arguments**: diff staged changes (`git diff --cached`) + unstaged (`git diff`)
- **One argument** (branch name, e.g. `feature/my-branch`): diff against `main`
  → `git diff main...feature/my-branch`
- **Two arguments** (e.g. `main feature/my-branch`): diff between those two refs
  → `git diff main...feature/my-branch`
- **`--staged`**: only staged changes → `git diff --cached`
- **`--last`**: last commit → `git diff HEAD~1 HEAD`

Run the appropriate git command and capture the diff output.

If the diff is empty, tell the user:
> "No changes found to review. Make sure your changes are staged, or specify a branch:
> `/review-pr feature/my-branch`"
Then stop.

## Step 2 — Diff summary

Before running agents, produce a quick summary for the user:

```
📋 Review scope: [what you're diffing]
📁 Files changed: [list of files from the diff]
📏 Lines changed: +[N] -[N]
```

Ask:
> "Does this look like the right diff? Type 'yes' to proceed or specify a different scope."

## Step 3 — Context gathering

Read the following if they exist — agents need this context:
- `ARCHITECTURE_RULES.md` — the hard constraints
- `DOMAIN_DICTIONARY.md` — ubiquitous language to check against
- `.claude/feature-workspace/analysis.md` — if this diff is part of an in-progress feature

If none of the above exist, note it — the review will be based on general craftsmanship
rules only.

## Step 4 — Code Review

Invoke the `code-reviewer` subagent.

Pass it:
- The full diff output
- The list of changed files (so it can read them in full context, not just the diff)
- Any feature workspace context from Step 3

The code-reviewer should evaluate against all criteria in its agent definition:
- Architecture & layer boundary violations
- Sandi Metz rules (complexity, LOC, params)
- Frontend craftsmanship (a11y, semantic HTML)
- Shift-left reliability (timeouts, idempotency, N+1)
- Fowler smells (Feature Envy, Magic Numbers, DRY)
- Ubiquitous Language violations

Output: `.claude/feature-workspace/code-review-report.md`

If the code-reviewer marks **CHANGES REQUESTED**, pause here and show the findings to the
user before proceeding to security review. Ask:
> "The code review found issues. Should I continue to security review anyway, or stop here
> so you can address the code quality issues first?"

## Step 5 — Security Review (conditional)

Invoke the `security-reviewer` subagent if the diff touches **any** of:
- Auth, session, or token handling
- API endpoints (new or modified routes)
- User-supplied input (forms, query params, file uploads)
- Secret or credential handling
- Rate limiting or account lockout
- Data crossing a trust boundary (external API calls, file reads, DB queries on user data)

Skip security review if the diff is purely:
- Documentation changes
- CI/CD config with no secrets
- Refactoring with no behavior change
- Test file additions with no production code changes

If skipping, tell the user why: "Skipping security review — no security surface in this diff."

Output: `.claude/feature-workspace/security-report.md`

## Step 6 — Review Summary

Write `.claude/feature-workspace/pr-review-summary.md`:

```markdown
# PR Review Summary

**Diff**: [scope reviewed]
**Date**: [today]

## Code Review: [APPROVED ✅ | CHANGES REQUESTED ⚠️]

### Findings
[Summary of code-review-report.md findings, or "None — approved as written"]

## Security Review: [APPROVED ✅ | FINDINGS ⚠️ | SKIPPED ⏭️]

### Findings
[Summary of security-report.md findings by severity, or "None identified" or "Not applicable"]

## Action Items
- [ ] [Specific change required — file, line, named refactoring operation]
— or "None — ready to merge"

## What Looks Good
[2-3 things done well — always include this, reviewers who only criticize are less effective]
```

Present this summary to the user.

## Rules

- Never modify source files — this is read-only review, not implementation
- Always include "What Looks Good" in the summary — pure criticism is less useful
- If the diff is very large (>500 lines), warn the user:
  > "This is a large diff ([N] lines). Consider breaking it into smaller PRs for more
  > focused review. I'll proceed but coverage may be less thorough."
- The code-reviewer iterative loop does NOT apply here — this is a review pass, not a
  delivery pipeline. Surface findings, don't fix them.
