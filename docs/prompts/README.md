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
| [automate-deliver-feature.md](automate-deliver-feature.md) | DESIGN an automated `/deliver-feature` workflow — policy-driven graduated automation preserving the framework's stage-3 human-in-the-loop stance. Produces `docs/aos/automated-delivery-design.md` + a follow-up Tier A implementation handoff. No code changes to the skill in this pass | Medium — design doc + handoff prompt + index update, 3 commits |
| [write-blog-posts.md](write-blog-posts.md) | Draft dev.to + LinkedIn blog posts covering recent framework + agent developments. Six candidate topics ranging from the mcp-add retrofit to the AOS migration to corpus-aware RAG | Small-Medium — one topic = ~2 files + 1 commit; menu-driven, do 1-3 topics per session |
| [add-frontmatter-contracts.md](add-frontmatter-contracts.md) — DONE | Add `shared/contracts/agent-frontmatter-contract.md` + `skill-frontmatter-contract.md` + `ki-frontmatter-contract.md` so `validate-artifact` can grep-check these the same way it checks pipeline artifacts. Shipped 2026-07-22 in commits `ae1e440` → `38b14a9` | Small — 3 contract files + validate-artifact skill update, 1-2 commits |
| [add-frontmatter-json-schemas.md](add-frontmatter-json-schemas.md) — DONE | Add `shared/schemas/agent-frontmatter.schema.json` etc. Enables IDE autocomplete + enum-value validation. Wire into VS Code / Cursor settings templates. Shipped 2026-07-22 in commits `fdfdd9b` → `a522e14` | Medium — 3 schema files + IDE settings + optional health-check integration, 2-4 commits |

## Convention

Every prompt in this directory:
- Names the target repo path (`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`)
- Restates that this repo IS the git repo — commits land here directly, not in a parent
- Lists the exact ops in scope with clear file paths
- Enumerates guardrails (commit per op, conventional commits, `git add` explicit paths only)
- Names escalation criteria (when to stop and ask)
- Requests a specific report format

When a prompt is executed and shipped, either delete the file or move it to `docs/prompts/done/` (convention TBD when the first one ships).
