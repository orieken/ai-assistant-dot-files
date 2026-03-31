---
name: debug-environment
description: Debug environment or configuration issues systematically
triggers:
  keywords: [fix PATH, debug config, environment issue, fix environment]
  intentPatterns: ["fix my path", "debug my configuration"]
standalone: true
---

## When To Use
When the user asks you to debug environment issues, such as missing commands, broken PATHs, or invalid configuration.

## Context To Load First
1. Check relevant config files based on context (e.g., `.zshrc`, `.bashrc`, `.bash_profile`, `package.json`, `.env`).

## Process
1. Run diagnostic commands to verify the current state of the environment or variables.
2. Formulate specific fixes.
3. Propose these fixes to the user with before/after verification commands.
4. Execute verify commands if the fix is applied.

## Output Format
Markdown explanation of the root cause, followed by a proposed fix snippet.

## Guardrails
- Do not apply changes to global bash/zsh profiles without explicitly stating the implication.
- Ensure that you use safe diagnostics commands before blindly exporting.

## Standalone Mode
Fully supported without external MCP servers.
