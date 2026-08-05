---
name: five-whys
description: Structured root cause analysis for bugs, outages, and recurring problems.
triggers:
  keywords: ["whys", "root cause", "post-mortem", "incident"]
  intentPatterns: ["Five whys on *", "Root cause *", "Why did * happen?", "Post-mortem on *"]
standalone: true
---

## When To Use
Triggered to investigate the root cause of an incident, outage, test failure, or recurring problem.

## Context To Load First
Any relevant files mentioned (error logs, test failures, workspace artifacts, git log for affected files).

## Process
1. State the symptom clearly: "The symptom is: [X]"
2. Ask "Why did [X] happen?" — wait, never suggest
3. Reflect: "So [X] happened because [their answer]. Why did [their answer] happen?"
4. Repeat until stop condition
5. Produce Root Cause Report
6. **Persist incident record** (production incidents and outages only — skip for non-incident analyses
   such as failing tests or design post-mortems): write `docs/incidents/<YYYY-MM-DD>-<kebab-slug>.md`
   using the Incident Record Schema below. Populate the **Candidate Records** section with any
   recurrence-prevention lessons (rule change, fitness function, prompt improvement) that emerged from
   the why chain — use `promote-memory`'s exact Candidate Record format so the same gated promotion
   machinery can consume them unchanged. Create `docs/incidents/` if it does not yet exist.

## Output Format

**Workspace output** (always): `.claude/feature-workspace/five-whys-[kebab-topic].md`

```markdown
# Five Whys: [Problem Statement]

Date: [today]

## The Symptom
[What was observed — specific, observable]

## The Why Chain
1. Why did [symptom] happen? → [answer]
2. Why did [answer] happen? → [answer]
3. ...

## Root Cause
[The actionable thing at the bottom of the chain]

## Investigation Gaps
[Any "I don't know" answers needing further investigation]

## Recommended Action
- Immediate: [what to do right now]
- Preventive: [fitness function, test, or process change that prevents recurrence]
- Owner: [who should act]
```

**Incident record** (production incidents only): `docs/incidents/<YYYY-MM-DD>-<kebab-slug>.md`
See Incident Record Schema below.

## Incident Record Schema

```markdown
# Incident: [kebab-slug]

**Date**: YYYY-MM-DD
**Severity**: P0-Critical | P1-High | P2-Medium | P3-Low
**Status**: Resolved | Ongoing
**Affected Feature**: [link to docs/features/<name>/ if the incident traces to a delivered feature — or "No matching feature delivery"]

## Summary
[One sentence: what broke and what the user-visible impact was]

## Timeline
| Time | Event |
|---|---|
| HH:MM | Alert fired / symptom reported |
| HH:MM | Root cause identified |
| HH:MM | Fix applied |
| HH:MM | Incident resolved |

## Five Whys Chain
1. Why did [symptom] happen? → [answer]
2. Why did [answer] happen? → [answer]
3. ...

## Root Cause
[The actionable thing at the bottom of the chain]

## Fix Applied
[Rollback to commit X / Hotfix: description / Pending]

## Candidate Records

### Candidate: [short title]
- **Source**: docs/incidents/YYYY-MM-DD-[slug].md, "[section name]"
- **Type**: KI | ADR-worthy | Rule-change-worthy | Lesson
- **Evidence**: [exact quote or close paraphrase supporting this candidate]
- **Tags**: [proposed frontmatter tags, if Type is KI]
- **Expiration condition**: [what would make this stop being true]
- **Existing overlap checked**: [KI or rule checked, result]

— or "None — no promotable lessons identified in this incident"
```

## Guardrails
- Never answer your own questions. Hypotheses stay private.
- Reflect back exactly what the user said before asking the next why.
- Ask one question at a time. Never stack questions.
- If the answer is vague ("because it broke"), ask for specificity.
- Stop when: answer is actionable / user says "I don't know" twice / five iterations complete.
- Never suggest an answer during facilitation. Never skip the reflection step.
- Never produce the report until the why chain is complete or a stop condition is hit.

## Standalone Mode
Fully conversational. No external tools needed.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
