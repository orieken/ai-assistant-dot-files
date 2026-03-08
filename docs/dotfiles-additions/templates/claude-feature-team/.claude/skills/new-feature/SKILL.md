---
name: new-feature
description: Scaffolds a new feature spec file from the project template, creates the feature workspace, and guides you through the key fields before kicking off delivery. Invoke with /new-feature [feature-name] or /new-feature "short description of the feature"
---

# New Feature Setup

You are setting up a new feature for delivery. The input is: $ARGUMENTS

## Step 1 — Derive the feature name

If `$ARGUMENTS` is already a kebab-case slug (e.g. `user-avatar-upload`), use it directly.
If `$ARGUMENTS` is a natural language description (e.g. "add rate limiting to the login endpoint"),
convert it to kebab-case: `add-rate-limiting-to-login-endpoint`.

Set `FEATURE_NAME` = the derived kebab-case slug.
Set `FEATURE_FILE` = `features/${FEATURE_NAME}.md`

## Step 2 — Preflight checks

1. Verify `features/TEMPLATE.md` exists. If it doesn't, tell the user:
   > "No `features/TEMPLATE.md` found. Run `scaffold-team.sh` to install the feature team
   > template, or create `features/TEMPLATE.md` manually."
   Then stop.

2. Verify `${FEATURE_FILE}` does NOT already exist. If it does, tell the user:
   > "`${FEATURE_FILE}` already exists. Did you mean to continue work on an existing feature?
   > Run `/deliver-feature ${FEATURE_FILE}` to kick off the pipeline, or choose a different name."
   Then stop.

## Step 3 — Create the feature file

Copy `features/TEMPLATE.md` to `${FEATURE_FILE}`.

Replace the template title placeholder with the feature name in title case.

## Step 4 — Interactive field population

Ask the user the following questions **one at a time**. After each answer, fill in the
corresponding section of `${FEATURE_FILE}` before asking the next question.

Do not ask all questions at once. One question, wait for answer, write it, next question.

**Question 1 — Summary:**
> "In one or two sentences, what does this feature do and why does it exist?
> (This becomes the Summary section — write for a non-technical stakeholder.)"

**Question 2 — Acceptance Criteria:**
> "List your acceptance criteria. Use plain language — I'll format them as Given/When/Then.
> One criterion per line. Type 'done' when finished."

Keep asking for criteria until the user types "done". Format each one as:
`- [ ] Given [context], when [action], then [observable outcome]`

**Question 3 — Out of Scope:**
> "What is explicitly NOT included in this feature? (Helps the analyst stay focused.
> Type 'none' if nothing comes to mind.)"

**Question 4 — Test Approach:**
> "What test framework should the QA agent use?
> - Saturday/Cucumber — for UI flows with clear BDD scenarios (default for most features)
> - Saturday/Playwright — for component/visual tests without Gherkin
> - Sunday/API — for pure API/contract testing with no UI
> - Leave blank — auto-detect from feature content
> Type the option or press enter to auto-detect."

**Question 5 — Open Questions:**
> "Any ambiguities or open questions the analyst should flag before implementation starts?
> Type 'none' if not."

## Step 5 — Create the feature workspace

```bash
mkdir -p .claude/feature-workspace
```

## Step 6 — Confirm and offer to launch

Show the user a summary:

```
✅ Feature spec created: ${FEATURE_FILE}
✅ Workspace ready:      .claude/feature-workspace/

Summary:     [their summary]
Criteria:    [N acceptance criteria]
Test approach: [their choice or "auto-detect"]
```

Then ask:
> "Ready to launch the delivery pipeline? I'll hand this to the analyst first and pause
> for your approval before any code is written.
>
> Type 'yes' to start, or 'no' to review `${FEATURE_FILE}` first."

If yes: immediately invoke the `analyst` subagent with `${FEATURE_FILE}` and follow the
full delivery pipeline as defined in `CLAUDE.md`.

If no: tell the user:
> "Your feature spec is ready at `${FEATURE_FILE}`. Edit it anytime, then run:
> `/deliver-feature ${FEATURE_FILE}`"

## Rules

- Never skip the one-at-a-time question flow — batching questions produces worse specs
- Never start the delivery pipeline without explicit user confirmation
- Never overwrite an existing feature file
- If the user provides a very detailed description as `$ARGUMENTS`, pre-populate what you
  can from it before asking questions (skip questions already answered by the description)
