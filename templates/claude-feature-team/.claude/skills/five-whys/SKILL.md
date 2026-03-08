---
name: five-whys
description: A standalone conversational skill for structured root cause analysis of bugs, outages, and recurring problems.
triggers:
  keywords: ["five whys", "root cause", "post-mortem", "recurring bugs", "keeps happening"]
  intentPatterns: ["five whys on {problem}", "root cause {issue}", "why did {x} happen?", "post-mortem on {incident}"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill when the user asks to perform a root cause analysis, "five whys", post-mortem, or when they express frustration about recurring bugs/incidents ("this keeps happening").
Do NOT use this to just give them an answer to a coding problem. It is a facilitation tool.

## Context To Load First
1. Any relevant files mentioned by the user (error logs, test failures, feature workspace artifacts, git log for affected files).

## Process
1. **State the problem clearly**: Ask the user: "The symptom is: [X]" (derive X from their prompt).
2. **Ask "Why did [X] happen?"**: Wait for the user's answer. **NEVER suggest an answer.**
3. **Reflect back**: When the user answers, reply with: "So [X] happened because [their answer]. Let me ask..."
4. **Ask the next why**: "Why did [their answer] happen?" — wait for the user's answer. **NEVER suggest one.**
5. **Iterate**: Repeat steps 3 and 4 until ONE of these stop conditions is met:
   - The answer exposes something actionable (a process gap, a missing test, a design smell, a missing fitness function).
   - The user says "I don't know" twice in a row (this is an investigation gap).
   - Five full iterations are complete.
6. **Produce Root Cause Report**: Once a stop condition is met, generate the final report using the format below.

### Facilitation Rules
- Never answer your own questions. If you have a hypothesis, keep it to yourself.
- Reflect back exactly what the user said before asking the next why.
- Ask ONE question at a time. Never stack questions.
- If the user's answer is vague ("because it broke"), ask for specificity: "Can you be more specific? What exactly broke?"
- End the session with questions and a recommended action, never a lecture.

## Output Format
Create `.claude/feature-workspace/five-whys-[kebab-topic].md`:

```markdown
# Five Whys: [Problem Statement]

Date: [today]
Facilitator: Claude
Participants: [from git config or ask user]

## The Symptom
[What was observed — specific, observable]

## The Why Chain
1. Why did [symptom] happen? → [answer]
2. Why did [answer] happen? → [answer]
3. Why did [answer] happen? → [answer]
4. Why did [answer] happen? → [answer]
5. Why did [answer] happen? → [root cause]

## Root Cause
[The actionable thing at the bottom of the chain]

## Investigation Gaps
[Any "I don't know" answers that need further investigation]

## Recommended Action
- Immediate: [what to do right now]
- Preventive: [what fitness function, test, or process change prevents recurrence]
- Owner: [who should act on this]
```

## Guardrails
- **No Suggestions**: You MUST NEVER suggest an answer during the facilitation process.
- **Always Reflect**: You MUST NEVER skip the reflection step before asking the next why.
- **Hold the Report**: You MUST NEVER produce the final report until the why chain is complete or a stop condition is explicitly hit.

## Standalone Mode
Fully conversational. No external tools or MCP servers are required. Relies purely on chat interaction with the user.
