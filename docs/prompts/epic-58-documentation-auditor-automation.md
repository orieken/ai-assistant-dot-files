# Epic 58 — Documentation Auditor Automation

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 5; re-confirmed open (and the
smallest remaining standalone epic) by `docs/audits/framework-gap-audit-2026-07-31.md` § F5.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

The `documentation-auditor` agent (read-only counter agent; audits README.md,
`docs/ARCHITECTURE.md`, `docs/AGENT_REFERENCE.md`, and prose docs for staleness against current
agent/skill inventories) exists with a complete golden-file fixture in
`tests/agents/documentation-auditor/`. What's missing is any *repeatable* invocation path — today
it only runs when a human remembers to ask. Adjacent machinery that already exists:

- `scheduler` skill (shared/skills/scheduler/) — orchestrates scheduled or hook-driven pipeline
  runs (cron triggers, periodic health checks).
- `shared/hooks/` layer (README, hooks-schema.md, examples/) — event-driven hook definitions,
  all opt-in.
- `scripts/check-inventory-drift.sh` — already catches *count* drift deterministically; the
  documentation-auditor's value-add is *prose* staleness that grep can't catch (renamed concepts,
  obsolete workflow descriptions, dead references).

## Scope

**Do NOT wire an LLM agent into `health-check.sh` itself** — that script is deterministic bash and
must stay runnable in CI without an LLM. The integration is config + docs, not code:

1. `shared/hooks/examples/on-inventory-change-doc-audit.<ext>` — a hook example (following the
   existing examples' schema) that triggers a `documentation-auditor` run when files under
   `shared/agents/` or `shared/skills/` change. Opt-in, like every hook.
2. A scheduler entry example: extend the `scheduler` skill's documented examples with a periodic
   (e.g. weekly) documentation-auditor run, output landing in a dated findings file (follow
   whatever output-location convention the agent's own prompt already declares — read it first).
3. `scripts/health-check.sh` — add a WARN-level (never FAIL) freshness pointer: if the newest
   documentation-auditor findings file is older than N days (or absent), print "consider running
   documentation-auditor". Pure file-mtime check, no LLM. Skip silently if the findings directory
   convention doesn't exist yet (opt-in guarantee).
4. Update `docs/AGENT_REFERENCE.md` documentation-auditor entry to describe the automation paths.

Commit sequence (one per op above):

1. `feat(hooks): add on-inventory-change doc-audit hook example (Epic 58)`
2. `docs(scheduler): add periodic documentation-auditor example (Epic 58)`
3. `feat(health-check): WARN when doc-audit findings are stale (Epic 58)`
4. `docs(agents): document documentation-auditor automation paths (Epic 58)`

After commits: `bash scripts/health-check.sh` green on a repo with NO findings directory (proves
the opt-in guarantee holds).

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If the documentation-auditor agent prompt declares no findings-output convention, halt and
  propose one (location + filename pattern) before building the mtime check against it.
- If the hooks schema can't express "on file change under a path", halt — extending the hooks
  schema is Phase-3-adjacent scope, not this epic.

## Report (under 100 words)

```
Commits:
  <sha> <message>
  ...
Findings-output convention used: <path pattern, and whether it pre-existed or was proposed>
health-check on pristine repo: <pass — WARN absent/silent?>
```

Go.
