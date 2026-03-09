---
name: event-storm
description: Collaborative domain modeling using Event Storming Lite before any feature starts.
triggers:
  keywords: ["event storm", "model domain", "domain events"]
  intentPatterns: ["Event storm *", "Model the domain for *", "What are the domain events for *?", "Let's do an event storm on *"]
standalone: true
---

## When To Use
Use for any feature that involves complex business logic, state machines, or multiple aggregates.
Do NOT use for simple CRUD with no business logic, UI-only changes with no domain events, or infrastructure/configuration changes.

## Context To Load First
1. The feature file
2. `DOMAIN_DICTIONARY.md`

## Process
1. State the feature scope: "We're modeling [feature]. Let's start from what can happen."
2. Ask: "What are the most important things that can *happen* in this feature? Give them as past-tense events — `[Noun][PastVerb]`."
   Wait. List them. Don't suggest yet.
3. For each event: "What *command* triggers [Event]? Who or what sends that command?"
4. "Which *aggregate* (the thing that owns the state) receives each command and emits each event?"
5. "Are there any *policies* — 'when [Event] happens, automatically trigger [Command]'?"
6. "What does the UI need to *read*? What projections or read models does this require?"
7. "Does this feature cross a *bounded context boundary*? If so, how do the contexts communicate?"
8. Produce the Event Storm Map

## Output Format
`.claude/feature-workspace/event-storm.md`

```markdown
# Event Storm: [Feature Name]

## Domain Events
Past-tense — things that happened:
- `UserLockedOut` — triggered after 5 consecutive failed login attempts

## Commands
What triggers each event:
| Command | Triggered by | Produces Event |
|---|---|---|
| `LockUserAccount` | Auth service (policy) | `UserLockedOut` |

## Aggregates
Who owns what:
| Aggregate | Owns Events | Key State |
|---|---|---|
| `UserAccount` | `UserLockedOut` | `failedAttempts`, `lockedUntil` |

## Policies
When → Then rules:
- When `LoginAttemptFailed` fires for the 5th time → trigger `LockUserAccount`

## Read Models
What the UI needs:
- `LoginPageState`: `isLocked: boolean`

## Bounded Context
- Owning context: Identity & Access
- Context crossings: None / [describe crossing + integration pattern]

## Hotspots
- [Question]: [Why it matters]

## DOMAIN_DICTIONARY.md Additions
- `UserLockedOut`: [definition]
```

## Guardrails
- Never suggest domain events before the user has given their first answer
- Never rename or combine the user's domain events without asking
- Hotspots are questions, never answers
- The output feeds directly into the analyst's `## Domain Events` section

## Standalone Mode
Fully conversational. No external tools needed.
