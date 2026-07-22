# AOS Migration Handoff Prompts

Self-contained agent handoff prompts for executing the AOS migration defined in `docs/aos/migration-plan.md`. Each prompt drives one phase of the migration in a fresh chat — no dependency on prior conversation memory.

## How to use

1. Read `docs/aos/migration-plan.md` for the full plan context.
2. Open a fresh Claude Code chat (or spawn a subagent).
3. Copy the entire target phase prompt into the message.
4. The agent has everything needed: repo path, plan reference, phase scope, backward-compat guarantee, commit discipline, escalation criteria, report format.

## Prompts

| File | Scope | Estimated size |
|---|---|---|
| [phase-1-foundations.md](phase-1-foundations.md) | Ops 1.1–1.7 — telemetry + evaluation directories, first counter agent (`memory-auditor`), health-check update, migration-guide stub. Ships as v3.0.0. | 7 ops, ~1 subagent session |
| _phase-2-governance.md_ | Not yet drafted. Ops 2.1–2.8 — 10 remaining audit-pair counter agents, 4 opposing-force skill pairs, hooks/ layer. Ships as v3.1.0. | ~15 ops |
| _phase-3-runtime.md_ | Not yet drafted. Ops 3.1–3.13 — orchestration/ + rag/ + Learning/Forgetting + trinity-native workflow refactor. Ships as v3.2.0. | ~13 ops (heaviest phase) |
| _phase-4-policy.md_ | Not yet drafted. Ops 4.1–4.8 — policy layer + evaluator + sample policies. Ships as v3.3.0. | ~8 ops |

## Convention

Every prompt in this directory:
- Names the target repo path explicitly
- Points at `docs/aos/migration-plan.md` as source of truth
- Restates the backward-compat guarantee (a team on prior version upgrading MUST see zero change from what they use today)
- Lists the exact ops in scope
- Enumerates guardrails (commit per op, conventional commits, never `git add -A`)
- Names escalation criteria (when to stop and ask instead of pushing through)
- Requests a specific report format

Draft new phase prompts as `phase-N-<slug>.md` when the prior phase lands cleanly. Do NOT draft all four upfront — each phase's prompt should be informed by what the prior phase's execution actually taught us (real learnings amend `docs/aos/migration-plan.md` as `## Plan Amendments`, same pattern as `mcp-add-plan.md` used).
