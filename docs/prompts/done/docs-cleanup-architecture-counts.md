# Docs Cleanup — Refresh Agent/Skill Counts in ARCHITECTURE.md

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/ARCHITECTURE.md` line 16 claims "25 agents" and line 17 claims "56 skills". The current
live counts are 39 agents and 69 skills. The file also describes the platform count as six, but
`shared/platform-registry.json` may reflect a higher number. These stale values cause the
`scripts/check-inventory-drift.sh` fitness function to report drift.

## Scope

**Op 1 — Establish live counts:**

Run the following to get authoritative numbers:
```bash
ls shared/agents/*.md | wc -l          # agent count
ls shared/skills/*/SKILL.md | wc -l    # skill count
cat shared/platform-registry.json | grep '"name"' | wc -l   # platform count
```

**Op 2 — Update ARCHITECTURE.md:**

- Fix the numeric claims on lines 16–17 to match the live counts.
- Scan the rest of the file for any other embedded numeric counts referencing agents, skills,
  or platforms and update them. Do not patch numbers in isolation — read each paragraph that
  embeds a count and verify the surrounding context is still accurate.
- If tier/generation details (e.g. "Generation X" labels, "Tier 1/2/3" counts) appear stale,
  note them in the report but do not silently patch what you cannot verify. Flag for human review.

**Op 3 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Refresh `docs/ARCHITECTURE.md`" item as `[x]`.

## Guardrails

- Commit per logical op is fine; a single commit covering all three ops is also acceptable.
- Conventional commit: `docs(architecture): refresh agent/skill/platform counts`
- Stage only modified files — do not `git add -A`.
- Do not rewrite prose sections; only correct numeric claims and their immediate context.

## Escalation

Stop and report if:
- The live counts differ substantially from both the old values and the expected 39/69 — the
  inventory may be in flux.
- Tier or generation details are ambiguous and cannot be verified against current files.

## Report

On completion, confirm:
- Old counts → new counts for agents, skills, platforms
- Any stale tier/generation details flagged for human review
- Commit hash(es)
