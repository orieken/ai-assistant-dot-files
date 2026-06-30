---
name: spec-writer
description: Creates or reviews feature specifications through structured interviews, then critiques readiness against every downstream agent's needs.
triggers:
  keywords: ["spec", "specification", "write spec", "review spec", "spec-writer"]
  intentPatterns: ["Write a spec for *", "Review this spec *", "Is this spec ready?", "/spec-writer *"]
standalone: true
---

## When To Use
When the user wants to create a new feature specification from scratch (Write Mode) or critique an existing spec for pipeline readiness (Review Mode).

Do NOT use when the user wants the full delivery pipeline — use `/deliver-feature` instead.
Do NOT use when the user wants a guided feature creation experience with directory setup — use `/new-feature` instead.

## Context To Load First
1. `features/TEMPLATE.md` — the spec structure to follow
2. `DOMAIN_DICTIONARY.md` — ubiquitous language reference
3. `CLAUDE.md` — project constraints

## Process

### Mode Detection
- **No file path given** or user says "write a spec for X" -> **Write Mode**
- **File path given** or user says "review this spec" -> **Review Mode**

### Write Mode
1. Invoke the `spec-writer` agent. It interviews the user one question at a time.
2. After all questions, the agent writes the spec to `features/<kebab-title>.md`.
3. The agent immediately enters Review Mode on the draft.

### Review Mode
1. Read the target spec file.
2. Run the readiness critique against every downstream agent's needs (analyst, architect, developer, QA, security, devops).
3. Produce a structured critique with verdict: READY or NEEDS WORK.
4. If NEEDS WORK: offer to fix the gaps and re-run the critique.
5. If READY: suggest running `/deliver-feature features/<name>.md` or `/new-feature`.

## Output Format

### Write Mode Output
- `features/<kebab-title>.md` — the spec file following the template structure

### Review Mode Output
Inline critique displayed to the user:

```markdown
## Spec Readiness Critique: [Title]

### Verdict
READY | NEEDS WORK

### Agent Readiness
| Agent | Status | Gap (if any) |
|---|---|---|
| Analyst | PASS / FAIL | [specific gap or "—"] |
| Architect | PASS / FAIL | [specific gap or "—"] |
| Developer | PASS / FAIL | [specific gap or "—"] |
| QA | PASS / FAIL | [specific gap or "—"] |
| Security | PASS / FAIL | [specific gap or "—"] |
| DevOps | PASS / FAIL | [specific gap or "—"] |

### Required Changes
[Only if verdict is NEEDS WORK — specific, actionable fixes]
```

## Guardrails
- One question at a time during the interview — never batch questions.
- Push back on implementation language in acceptance criteria — every time, without exception.
- Never declare READY with an open question that would block a downstream agent.
- The critique evaluates the spec, not the person — be direct about gaps.
- Never start the delivery pipeline without explicit user confirmation.

## Standalone Mode
Works entirely offline. No external services required.
