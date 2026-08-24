# ADR-004: Bugfix Artifacts Archive in docs/features/, Not a Separate Namespace

## Status

Accepted

## Date

2026-08-05

## Context

When `deliver-bugfix` (the lightweight bugfix pipeline introduced in Epic 67) produces artifacts
— characterization test, implementation-notes.md, code-review-report.md, qa-report.md — they
need a permanent home in the repository. Two options were considered.

### Considered options

| Option | Pros | Cons |
|---|---|---|
| **Shared `docs/features/<bug-slug>/`** (same as features) | `extract-lessons`, `retrospective`, `analyst`, and `context-engineer` all already scan `docs/features/` — bugs automatically join the learning loop; no new scan target needed | Bug artifacts mixed with feature artifacts; slug must be distinct (e.g., `fix-login-token-expiry` vs `login-auth`) |
| Separate `docs/bugs/<bug-slug>/` | Clear namespace separation; bugs never confused with features | Every downstream consumer (`extract-lessons` Step 1-2, `retrospective`, `analyst` prior-delivery scan, `context-engineer` bounded-context grep) would need a new scan path; risk of forgetting one and creating a blind spot in the learning loop |

## Decision

Bugfix artifacts are stored at `docs/features/<bug-slug>/` — the same namespace as full feature
deliveries — using a slug that describes the fix (e.g., `fix-login-token-expiry`), not a generic
`bug-<id>` pattern.

The incident record at `docs/incidents/<date>-<slug>.md` separately preserves the production
signal (severity, blast radius, timeline, five-whys chain) that feature artifacts don't carry.
The two records are linked: `deliver-bugfix` Phase 5 writes a "Fixed by" pointer into the
incident record pointing at `docs/features/<bug-slug>/`.

## Consequences

### What becomes easier

- `extract-lessons` Steps 1-2 (code-review pattern mining) and Step 6 (incident-feature pair
  mining) automatically see bugfix artifacts alongside feature artifacts with no config change.
- `retrospective` sees bugfixes as first-class entries ("did the fix close the incident
  fully?" is a question a retrospective should ask).
- `analyst` and `context-engineer` surface related bugfixes when a new feature touches the same
  bounded context — without any new scan wiring.
- One archive to search, one index to maintain.

### What becomes harder

- A slug collision between a feature and a bugfix targeting the same subsystem is possible and
  must be caught manually. Convention: use `fix-` prefix for bug slugs to distinguish
  (`fix-user-session-expiry` vs `user-auth`).

### Fitness function

`health-check.sh` already enforces `docs/features/` writability. No new check required.
Slug collision is a human guardrail, not a CI check (low-frequency, easily caught during naming).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
