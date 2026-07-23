# Framework Improvement Prompts

Self-contained agent handoff prompts for framework-level improvements to `ai-assistant-dot-files`. Each file is a fully-briefed prompt a fresh Claude Code chat or subagent can execute standalone.

Separate from `docs/aos/prompts/` (AOS-migration-specific) — these are general framework improvements not tied to the AOS phased migration.

## How to use

1. Open a fresh Claude Code chat or spawn a subagent.
2. Copy the entire target prompt file into the message.
3. The agent has everything needed: repo path, prior state, scope, guardrails, commit discipline, escalation criteria, report format.

## Prompts

| File | Scope | Estimated size |
|---|---|---|
| [add-frontmatter-contracts.md](add-frontmatter-contracts.md) | Add `shared/contracts/agent-frontmatter-contract.md` + `skill-frontmatter-contract.md` + `ki-frontmatter-contract.md` so `validate-artifact` can grep-check these the same way it checks pipeline artifacts | Small — 3 contract files + validate-artifact skill update, 1-2 commits |
| [add-frontmatter-json-schemas.md](add-frontmatter-json-schemas.md) | Add `shared/schemas/agent-frontmatter.schema.json` etc. Enables IDE autocomplete + enum-value validation. Wire into VS Code / Cursor settings templates | Medium — 3 schema files + IDE settings + optional health-check integration, 2-4 commits |

## Convention

Every prompt in this directory:
- Names the target repo path (`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`)
- Restates that this repo IS the git repo — commits land here directly, not in a parent
- Lists the exact ops in scope with clear file paths
- Enumerates guardrails (commit per op, conventional commits, `git add` explicit paths only)
- Names escalation criteria (when to stop and ask)
- Requests a specific report format

When a prompt is executed and shipped, either delete the file or move it to `docs/prompts/done/` (convention TBD when the first one ships).
