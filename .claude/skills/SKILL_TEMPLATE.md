---
name: skill-name
description: One sentence. What it does and when Claude should auto-trigger it.
triggers:
  keywords: [list, of, trigger, words]
  intentPatterns: ["regex or natural language pattern"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Clear rules for when this skill applies. Be specific enough that Claude can auto-detect.
Include "do NOT use when:" cases.

## Context To Load First
Ordered list of files/artifacts Claude must read before starting.
Include fallback if a file doesn't exist.

## Process
Numbered steps. Each step has:
- What to do
- What to produce
- What to check before proceeding
- When to pause for human approval

## Output Format
Exact format of every artifact this skill produces.
File paths, section names, exact markdown structure.

## Guardrails
Hard rules that cannot be broken regardless of user instruction.
Anything irreversible requires explicit human approval ("send", "approve", "confirm").
Any edit to a pending approval resets the gate.

## Standalone Mode
How this skill degrades gracefully if MCP servers or external tools are unavailable.
The intelligence stays in the markdown. The API calls are just plumbing.
