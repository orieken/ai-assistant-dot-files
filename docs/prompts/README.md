# Framework Improvement Prompts

Self-contained agent handoff prompts for framework-level improvements to `ai-assistant-dot-files`. Each file is a fully-briefed prompt a fresh Claude Code chat or subagent can execute standalone.

Separate from `docs/aos/prompts/` (AOS-migration-specific) — these are general framework improvements not tied to the AOS phased migration.

## How to use

1. Open a fresh Claude Code chat or spawn a subagent.
2. Copy the entire target prompt file into the message.
3. The agent has everything needed: repo path, prior state, scope, guardrails, commit discipline, escalation criteria, report format.

For tasks that require a human (not a fireable agent prompt), see [../human-tasks.md](../human-tasks.md).

## Active Prompts

| File | Scope | Estimated size |
|---|---|---|
| [framework-hygiene-sweep.md](framework-hygiene-sweep.md) | Eight small pending items: commit `.gitignore` mod, track blog drafts, delete redundant AOS zip, add KI template, adopt `done/` convention, investigate audits doc, fix 2 pre-existing health-check WARNs | Small — 7-8 commits, all mechanical |
| [capture-session-history.md](capture-session-history.md) | Investigate the framework's existing archive patterns (`docs/features/`, `docs/lessons-learned/`, `/retrospective` skill), propose one of four options for capturing the 2026-07 AOS foundations + frontmatter + retrofit sessions as history, execute on user approval | Small — investigation + user decision + 0-3 commits |
| [automate-deliver-feature.md](automate-deliver-feature.md) | DESIGN an automated `/deliver-feature` workflow — policy-driven graduated automation preserving the framework's stage-3 human-in-the-loop stance. Produces `docs/aos/automated-delivery-design.md` + a follow-up Tier A implementation handoff. No code changes to the skill in this pass | Medium — design doc + handoff prompt + index update, 3 commits |
| [write-blog-posts.md](write-blog-posts.md) | Draft dev.to + LinkedIn blog posts covering recent framework + agent developments. Six candidate topics ranging from the mcp-add retrofit to the AOS migration to corpus-aware RAG | Small-Medium — one topic = ~2 files + 1 commit; menu-driven, do 1-3 topics per session |

## Completed Prompts (`docs/prompts/done/`)

| File | Scope | Shipped |
|---|---|---|
| [done/add-frontmatter-contracts.md](done/add-frontmatter-contracts.md) | Add `shared/contracts/agent-frontmatter-contract.md` + `skill-frontmatter-contract.md` + `ki-frontmatter-contract.md` so `validate-artifact` can grep-check these | Shipped 2026-07-22 in commits `ae1e440` → `38b14a9` |
| [done/add-frontmatter-json-schemas.md](done/add-frontmatter-json-schemas.md) | Add `shared/schemas/agent-frontmatter.schema.json` etc. Enables IDE autocomplete + enum-value validation. Wire into VS Code / Cursor settings templates | Shipped 2026-07-22 in commits `fdfdd9b` → `a522e14` |

## Convention

Every prompt in this directory:
- Names the target repo path (`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`)
- Restates that this repo IS the git repo — commits land here directly, not in a parent
- Lists the exact ops in scope with clear file paths
- Enumerates guardrails (commit per op, conventional commits, `git add` explicit paths only)
- Names escalation criteria (when to stop and ask)
- Requests a specific report format

When a prompt is executed and shipped, move it to `docs/prompts/done/` and update the table above.

