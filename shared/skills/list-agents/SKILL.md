---
name: list-agents
description: Search for custom agent configurations in .claude/ directory and present them in a formatted list.
triggers:
  keywords: [list agents, custom agents, show agents, what agents]
  intentPatterns: ["list all custom agents", "what agents do we have"]
standalone: true
---

## When To Use
Use this skill when the user asks to see what custom agents or roles are available in the repository.

## Context To Load First
1. Verify that the `.claude/agents/` directory exists.

## Process
1. Use a tool to list the contents of the `.claude/agents/` directory.
2. Filter the output to show only `.md` files.
3. Present the list of agents back to the user in a well-formatted markdown list.
4. If the directory is missing or empty, inform the user that no custom agents were found.

## Output Format
A clear markdown list where each item is the name of the agent (the file basename without the `.md` extension).

## Guardrails
- Do not run any commands outside of listing the directory.

## Standalone Mode
Fully supported without external MCP servers. Just uses directory listing tools.
