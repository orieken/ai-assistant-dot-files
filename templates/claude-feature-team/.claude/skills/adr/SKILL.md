---
name: adr
description: A standalone skill to interactively draft and record Architecture Decision Records (ADRs).
triggers:
  keywords: ["adr", "architecture decision record", "document decision"]
  intentPatterns: ["write an adr for {decision}", "document the decision to {decision}", "we decided to {decision}, record it", "create an architecture decision record"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill when the user asks to document, record, or write an architectural decision.
Do NOT use when the user is asking the Architect to *make* a decision. This skill is for *recording* a decision that has already been made or proposed.

## Context To Load First
1. All existing ADRs in `docs/adr/` (to determine the next number and match tone).
2. `ARCHITECTURE_RULES.md`
3. `DOMAIN_DICTIONARY.md`
4. `.claude/feature-workspace/architecture-notes.md` (if it exists).

## Process
1. **Determine the next ADR number**: Scan the `docs/adr/` directory for existing files to find the highest `ADR-[NNN]` number. The new ADR will be `NNN + 1`.
2. **Interactive Q&A**: Ask the user the following questions one by one (or extract the answers if already provided in the prompt):
   - "What was decided?" (must be one sentence, active voice)
   - "What was the context that made this decision necessary?"
   - "What alternatives were considered and why were they rejected?"
   - "What are the consequences — what becomes easier, harder, or different?"
   - "Does this decision produce a fitness function? If so, how is it enforced?"
3. **Draft the ADR**: Create the markdown content using the Output Format below.
   - What to check: Verify the title uses active voice (e.g., "Use X", not "Decision to use X").
   - When to pause: **STOP HERE explicitly for user approval.** Show the drafted ADR to the user and ask: "Does this look correct? Reply 'approve' to save."
4. **Save the ADR**: Once approved, write the file to `docs/adr/ADR-[NNN]-[kebab-title].md`.
5. **Update the Index**: Update the `docs/adr/README.md` file to include the new ADR. If `docs/adr/README.md` does not exist, create it with a simple markdown list of all ADRs.

## Output Format
Create `docs/adr/ADR-[NNN]-[kebab-title].md`:

```markdown
# ADR-[NNN]: [Title]

Date: YYYY-MM-DD
Status: Accepted
Deciders: [from git config user.name or ask user]
Technical Story: [feature or ticket that prompted this]

## Context
[What situation, constraint, or force prompted this decision]

## Decision
[What was decided — specific, concrete, active voice]

## Alternatives Considered
| Option | Why Rejected |
|---|---|

## Consequences
- **Easier**: [what this unlocks]
- **Harder**: [what this constrains]
- **Changed**: [what is different going forward]

## Fitness Function
[How this decision is enforced in CI — or "Judgment-only: [reason]"]
```

## Guardrails
- **Mandatory Approval**: You MUST NEVER write the ADR file to disk without explicit user approval of the draft.
- **Sequence Integrity**: You MUST NEVER number an ADR without checking the existing sequence in `docs/adr/`.
- **Active Voice**: ADR titles MUST use active voice (e.g., "Use React", not "Decision to use React").

## Standalone Mode
Works fully without MCP. Uses standard local file reading to scan `docs/adr/` and writes local files once approved by the user through chat.
