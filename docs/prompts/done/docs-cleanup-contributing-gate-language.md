# Docs Cleanup — Fix Gate #7 Misapplication in CONTRIBUTING.md

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/CONTRIBUTING.md`, section "Adding a new rule" (around line 62–69), states:

> This requires human sign-off — `.claude/rules/approval-gates.md` Gate #7 ("Wiring a New
> Fitness Function") applies to any change here, since every agent treats these as session-long
> law.

This is incorrect. Gate #7 specifically governs **wiring a fitness function into CI/CD**. It
does not apply to all `shared/rules/` edits. Editing an existing rule (e.g. correcting a
guideline) or adding a new rule file does not inherently trigger Gate #7 — only if that rule
change is simultaneously wired as a CI check does the gate apply.

The gate wording in `shared/rules/approval-gates.md` (Gate #7):
> **Action**: Modifying CI/CD pipelines to enforce a new architectural property.

## Scope

**Op 1 — Read the affected section:**

Read `docs/CONTRIBUTING.md` lines 62–75. Confirm the current gate language.

Also read `shared/rules/approval-gates.md` Gate #7 to confirm its exact scope.

**Op 2 — Rewrite step 2 of the "Adding a new rule" section:**

Replace the current step 2 with language that distinguishes:

**(a) Editing or adding a rule that does NOT wire a CI check:**
Human review is recommended (because agents load these as law), but Gate #7 does not apply.
A PR review or the user explicitly saying "approved" is sufficient.

**(b) Adding a rule that IS simultaneously wired as a CI fitness function:**
Gate #7 applies. The user must say "approve fitness function" or "add to CI" before the
CI wiring step is committed.

Keep the step concise — two sub-bullets or a short conditional is fine. Do not over-explain.

**Op 3 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Correct or clarify the approval language in `docs/CONTRIBUTING.md`"
item as `[x]`.

## Guardrails

- Conventional commit: `docs(contributing): correct Gate #7 scope for shared/rules changes`
- Stage only `docs/CONTRIBUTING.md` and `docs/TODO.md`.
- Do not rewrite other sections of CONTRIBUTING.md.
- The fix must be accurate — do not introduce a new inaccuracy by being too permissive or
  too restrictive about when human approval is needed.

## Escalation

Stop and report if the Gate #7 wording in `shared/rules/approval-gates.md` is itself ambiguous
or contradicts this prompt's reading of its scope — flag for human review before changing
CONTRIBUTING.md.

## Report

On completion, confirm:
- The old step 2 text (quote it)
- The new step 2 text (quote it)
- Commit hash
