# Epic 51 — Enterprise Memory Sync (`install.sh --sync-memory`)

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 4.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

Knowledge Items reside in two local paths per `shared/memory-registry.json`:
- `shared/knowledge/` — portable KIs (shipped with the framework itself)
- `.claude/knowledge/` — project-specific KIs (local to a single codebase, not portable)

No mechanism today for an ORGANIZATION to sync KIs across multiple installed projects — every team learns in isolation. Enterprise use case: KIs discovered in project A should be discoverable by team B without manual copy-paste.

## Scope

**Phase A — Design** (one commit): draft the sync protocol.

Answer:
- Where does the org's shared KI corpus live? Options: separate git repo (`<org>/knowledge-hub`), git submodule, org-scoped npm-like registry, or cloud sync via S3
- Recommend: separate git repo model — simplest, familiar, git-native conflict resolution, aligns with how the framework itself is distributed
- Auth: SSH-based (GitHub org membership) as the M1 path; OAuth/token later
- Conflict resolution: local wins on divergence for `.claude/knowledge/` (project-specific); org-repo wins for `shared/knowledge/` promotions
- Cadence: manual `install.sh --sync-memory` pull; push flow separate (opt-in)

Draft the design doc at `docs/features/enterprise-memory-sync/design.md` OR as an ADR at `docs/adrs/ADR-003-enterprise-memory-sync.md` (recommend ADR — this is a real architectural decision).

Commit: `docs(adr): ADR-003 enterprise memory sync design (Epic 51 Phase A)`.

**Pause for user approval before Phase B.**

**Phase B — Implementation** (multiple commits):

1. Add `--sync-memory` flag to `install.sh` with subcommands `pull` (default) and `push`
2. `scripts/sync-memory.sh` — the actual sync logic: git clone/pull the org KI repo to a cache dir, merge KIs into `shared/knowledge/` respecting conflict-resolution rules, verify frontmatter still validates against `shared/schemas/ki-frontmatter.schema.json`
3. Push flow (`install.sh --sync-memory push`): stage local KI changes → PR them to the org repo (never direct commit)
4. Update `shared/memory-registry.json` with a new field indicating sync-source per KI (locally-authored vs. sync-pulled)
5. Update `README.md` with the sync workflow
6. `scripts/health-check.sh` reports sync health (last-sync date, any conflicts)

Commit per op.

## Discipline

- Every sync operation must be a diff-and-approve step, never an auto-mutate. The pull command should show what it will change; require `--confirm` flag to apply.
- Standard git hygiene — never `git add -A`.
- No push during Phase B — human step.

## Escalation

- If the design decision on "where org KIs live" is genuinely unclear — halt, propose 2-3 options with tradeoffs, get user pick.
- If auth strategy needs to support enterprises where SSH is disabled — halt, propose HTTPS + token alternative.
- If a project already has an unrelated `.claude/knowledge/` KI with the same name as an org KI — halt at Phase B Op 2, describe the collision behavior.

## Report (under 250 words)

```
Phase A commit: <sha> (ADR-003)
Chosen architecture: <separate git repo | submodule | ...>
Auth strategy: <SSH | HTTPS+token | ...>
Conflict resolution: <local wins | org wins | ...>

Phase B commits (if approved):
  <sha> <message>
  ...

Verification: sync against a scratch org-repo works both directions, conflict detection works, health-check reports last-sync date.
```

Go.
