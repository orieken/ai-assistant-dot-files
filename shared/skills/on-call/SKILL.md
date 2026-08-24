---
name: on-call
description: Active incident response — for when something is broken right now.
triggers:
  keywords: ["incident", "outage", "broken", "alert", "rollback"]
  intentPatterns: ["Something is broken in production", "The site is down", "Alert firing for *", "Incident on *", "Rollback *", "* is failing in CI"]
standalone: true
---

## When To Use
Triggered when an active incident or urgency occurs. Has two modes:
- **Active Incident**: system is broken RIGHT NOW — speed matters, minimize cognitive load
- **Post-Incident**: system is stable — shift to Five Whys for root cause, this skill for timeline reconstruction and communication

## Context To Load First
Incident logs, alerts, deployment info.

## Process
**Active Incident Process (minimize steps — time matters):**
1. Establish the blast radius immediately:
   "What is the observable symptom? How many users/requests are affected?"
2. Check for a recent deploy: `git log --oneline -10`
   If yes → rollback decision: "Is reverting the last deploy safe and faster than a hotfix?"
3. Rollback or hotfix decision tree:
   - Rollback: `git revert HEAD` (a commit — requires the user to say "commit" or "approve commit" per
     `shared/rules/approval-gates.md` Gate #2) or redeploy previous image (a deploy — requires "deploy" per
     Gate #8) — safe if no DB migration in deploy. Incident urgency does not waive these gates; state the
     proposed action and get the explicit confirmation word before executing it, even under time pressure.
   - Hotfix: only if rollback is not safe (migration already ran) or fix is trivially small
4. Communicate: draft the incident message (who, what, ETA)
5. Once stable: capture timeline for the Five Whys session
6. **Persist incident record**: after resolution (or after the five-whys session completes), write
   `docs/incidents/<YYYY-MM-DD>-<slug>.md`. Populate with the timeline, blast radius, and fix applied
   from this report, plus the five-whys chain and root cause if `/five-whys` has already run. If the
   five-whys session runs later, link back to this file from it. Add a **Candidate Records** section
   for any recurrence-prevention lessons immediately apparent — use `promote-memory`'s exact format so
   they flow through the same gated promotion machinery. Check `docs/features/` for a matching feature
   delivery to link under "Affected Feature." Create `docs/incidents/` if it does not yet exist.

## Output Format
`.claude/feature-workspace/incident-[description].md`

```markdown
# Incident Report: [Description]

## Timeline
| Time | Event |
|---|---|
| HH:MM | Alert fired / symptom reported |
| HH:MM | Blast radius established |
| HH:MM | Rollback/hotfix decision made |
| HH:MM | Fix deployed |
| HH:MM | Incident resolved |

## Blast Radius
- Users affected: [N / all / subset — be specific]
- Services affected: [list]
- Data integrity: [affected / not affected]

## Root Cause (preliminary)
[One sentence — full RCA via Five Whys to follow]

## Fix Applied
[Rollback to commit X / Hotfix: description]

## Rollback Safety
Safe ✅ / Unsafe ⚠️ [reason — DB migration state]

## Communication Sent
- [Channel]: [timestamp] — [which template used]

## Follow-Up Required
- [ ] Five Whys session scheduled: [date]
- [ ] Fitness function added to prevent recurrence: [yes / not yet identified]
- [ ] Incident record persisted: docs/incidents/[YYYY-MM-DD]-[slug].md [pending / done]

## Candidate Records

### Candidate: [short title]
- **Source**: docs/incidents/[slug].md, "[section]"
- **Type**: KI | ADR-worthy | Rule-change-worthy | Lesson
- **Evidence**: [quote or paraphrase]
- **Tags**: [proposed tags if KI]
- **Expiration condition**: [what would make this stop being true]
- **Existing overlap checked**: [result]

— or "None — no promotable lessons immediately apparent"
```

## Guardrails
- During active incident: never ask more than one question at a time
- Never recommend a complex fix when a rollback is safe
- Never lower a fitness function threshold to make a failing check green
- Rollback + Five Whys is almost always better than hotfix under pressure
- **Never** execute the rollback commit or redeploy without the explicit "commit"/"deploy" confirmation
  `shared/rules/approval-gates.md` requires — incident urgency shortens the number of questions asked, it
  does not remove the approval gate itself

## Standalone Mode
Pure reasoning and templates. No external tools required during active incident. All commands are standard git/pnpm/kubectl operations.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
