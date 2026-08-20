# Docs Cleanup — Extend check-inventory-drift.sh to Catch Stale Doc Counts

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## ⚠️ Approval Gate Required

**This prompt includes a CI fitness function wiring step. Approval Gate #7 ("Wiring a New
Fitness Function") applies. Do NOT commit the CI wiring until the user explicitly says
"approve fitness function" or "add to CI".**

Implement the script changes and present them for review. Only commit and wire into CI after
explicit approval. The script improvement itself (Ops 1–2) may be committed independently
without the gate.

## Context

`scripts/check-inventory-drift.sh` was added in Epic 53 to count agents and skills and compare
them against authoritative prose. However, `docs/ARCHITECTURE.md` and `docs/ONBOARDING.md`
both contained stale counts that the check did not catch — the check reported zero drift while
both files had wrong numbers. The fitness function needs to be extended to inspect these two
files directly.

## Scope

**Op 1 — Read the current script:**

Read `scripts/check-inventory-drift.sh` in full to understand its current structure and output
conventions (WARN/FAIL format, how counts are computed, how it reports).

**Op 2 — Add doc-count checks:**

Add a new section to the script that:

1. Gets live agent and skill counts:
   ```bash
   LIVE_AGENTS=$(ls shared/agents/*.md 2>/dev/null | wc -l | tr -d ' ')
   LIVE_SKILLS=$(ls shared/skills/*/SKILL.md 2>/dev/null | wc -l | tr -d ' ')
   ```

2. Greps `docs/ARCHITECTURE.md` for numeric agent/skill claims using patterns like
   `[0-9]+ agents` and `[0-9]+ skills`. Extract the number and compare against live counts.
   If any extracted number does not match the live count, report WARN or FAIL (use WARN to
   match the existing check's severity — do not introduce a new FAIL level without the gate).

3. Does the same for `docs/ONBOARDING.md`.

4. If a file has no matching pattern (counts have been removed and replaced with links), skip
   it and report "no hardcoded counts found — OK".

Use the existing script's output pattern for consistency. Do not alter existing checks.

**Op 3 — Present for Gate #7 approval:**

After implementing Ops 1–2, output the diff and the full updated script for human review.
State clearly:

> Gate #7 check: the above change adds a new fitness function. It should be wired into
> `health-check.sh` only after you say "approve fitness function" or "add to CI".

**Op 4 — Wire into health-check.sh (AFTER GATE APPROVAL ONLY):**

Once the user approves, add a call to the new check section in `scripts/health-check.sh`,
following the pattern used for the existing drift check. Confirm the gate was cleared before
making this commit.

**Op 5 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Extend `scripts/check-inventory-drift.sh`" item as `[x]`
after both the script change and the CI wiring are complete.

## Guardrails

- Conventional commit for script change: `chore(scripts): extend drift check to catch stale doc counts`
- Conventional commit for CI wiring: `chore(ci): wire doc-count drift check into health-check`
- Two separate commits — one for the script, one for the CI wiring (post-approval only).
- Stage only modified files — do not `git add -A`.

## Escalation

Stop and report if:
- The script structure makes it unclear where to add the new checks.
- Grepping for numeric patterns produces too many false positives (e.g. version numbers, dates).
- The user does not approve the CI wiring — commit only the script change and close.

## Report

On completion (after full approval and wiring), confirm:
- What patterns the new checks use to extract counts from the two doc files
- Whether any edge cases (no hardcoded counts) are handled
- Whether Gate #7 was cleared and by what confirmation
- Commit hash(es) for script change and CI wiring separately
