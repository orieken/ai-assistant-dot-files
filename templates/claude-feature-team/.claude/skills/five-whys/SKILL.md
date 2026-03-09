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

## Output Format
`.claude/feature-workspace/five-whys-[kebab-topic].md`

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
