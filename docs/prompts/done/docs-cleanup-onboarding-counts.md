# Docs Cleanup — Refresh Stale Counts in ONBOARDING.md

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/TODO.md` flags that `docs/ONBOARDING.md` claims "24 agents and 53 skills" while the live
counts are 39 and 69 respectively. The preferred fix is to replace embedded counts with pointers
to generated/current inventory so the numbers cannot drift again.

## Scope

**Op 1 — Verify current state:**

Read `docs/ONBOARDING.md` in full and grep for any numeric agent/skill count claims:
```bash
grep -n "agents\|skills" docs/ONBOARDING.md | grep -E "[0-9]+"
```

If no stale numeric claims exist (they may have already been removed), document that finding
and close the TODO item without further changes.

**Op 2 — Fix stale claims (if present):**

For each hardcoded count found:
- If it is stale, replace it with a pointer to `docs/AGENT_REFERENCE.md` (for agent count) or
  `docs/ARCHITECTURE.md` (for skill count) rather than re-embedding the current number.
  Example: replace "39 agents" with "the full agent list in `docs/AGENT_REFERENCE.md`".
- If the phrasing genuinely needs a number (e.g. in a specific comparison), update to the live
  count from:
  ```bash
  ls shared/agents/*.md | wc -l
  ls shared/skills/*/SKILL.md | wc -l
  ```

**Op 3 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Refresh `docs/ONBOARDING.md`" item as `[x]`.

## Guardrails

- Conventional commit: `docs(onboarding): replace hardcoded counts with inventory pointers`
  (or `docs(onboarding): verify counts are current — no changes needed` if Op 1 finds nothing stale)
- Stage only `docs/ONBOARDING.md` and `docs/TODO.md`.
- Do not rewrite onboarding prose beyond the count fix.

## Escalation

Stop and report if the onboarding document structure seems significantly outdated beyond just
counts — broader refresh is out of scope for this prompt.

## Report

On completion, confirm:
- Whether stale counts were found (and what they were) or not
- What changes were made (or why none were needed)
- Commit hash
