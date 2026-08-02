# Contract: Knowledge Item frontmatter (`shared/knowledge/*.md`, `.claude/knowledge/*.md`)

**Produced by**: humans authoring Knowledge Items (typically via the `create-ki` skill) under
`shared/knowledge/` (portable, cross-project) or `.claude/knowledge/` (project-specific)
**Consumed by**: `search-ki` (tag-based and domain-based filtering), `memory-auditor` (coverage
analysis, duplicate detection), `memory-engineer` (periodic curation and deduplication), the
`shared/memory-registry.json` index, and other KIs linking via `[[link]]` syntax by `name`.

This contract governs the YAML frontmatter block at the top of every Knowledge Item file — not the body.
The body content (the actual knowledge captured) is judgment-only; only the frontmatter is
contract-bound because it is what the memory tooling grep-parses.

`README.md` files inside knowledge directories are excluded — they are documentation of the directory,
not Knowledge Items themselves.

## Required Fields

| Field | Type | Notes |
|---|---|---|
| `name` | string | kebab-case; must match filename base (`workflow-tool-wraps-domain-workflow-for-mcp.md` → `name: workflow-tool-wraps-domain-workflow-for-mcp`). Referenced by `[[link]]` syntax from other KIs and by the memory-registry index. |
| `tags` | list | Free-form tag list. Convention: use lowercase kebab-case tags (`tag-name` not `TagName`); reuse existing tags in the corpus when possible before inventing new ones. Consumed by `search-ki` for tag-based filtering and by `memory-auditor` for coverage analysis. |
| `domain` | string | Which bounded context this KI applies to (e.g., `testing`, `architecture`, `deployment`, `retrieval`, `mcp-server-pattern`). Consumed by `context-engineer` when building context manifests for a feature in a matching domain. |
| `created` | ISO date (`YYYY-MM-DD`) | Immutable — set once at authoring time. Consumed by the forgetting engine (future) for staleness calculations. |

## Validation Rule

`validate-artifact` checks:

1. **Field presence** — every required field above (`name`, `tags`, `domain`, `created`) must appear
   as a top-level YAML key inside the opening `---` / closing `---` frontmatter block. Missing any one
   is a FAIL. This matches the field-presence check in `scripts/health-check.sh` step 8 (the two
   enforcement paths agree on shape by design — this contract is the referenceable version of what the
   health-check script already enforces).
2. **`created` format** — the `created` value must match ISO 8601 date format
   `^[0-9]{4}-[0-9]{2}-[0-9]{2}$` (e.g., `2026-07-19`). Timestamps, slash-separated dates, and
   textual dates ("July 19, 2026") are FAILs. The staleness math downstream expects this exact shape.
3. **`tags` shape** — must be a YAML list (either flow `[a, b, c]` or block `- a\n- b\n- c` form).
   A bare string is a FAIL — `search-ki` expects to iterate.
4. **`name` shape** — must be lowercase kebab-case matching `^[a-z][a-z0-9-]*$` and must equal the
   filename base.
5. **No duplicate `name` across the KI corpus** — `scripts/health-check.sh` already flags this at the
   Memory Registry step; validate-artifact reproduces the check for standalone use. Two KIs sharing an
   exact frontmatter `name:` is a FAIL — `memory-engineer` should audit for a merge before either can
   ship.

This is a structural check only. It does not verify the KI's body content is accurate, non-duplicative
in *meaning* (only in `name`), or worth keeping — those judgments belong to `memory-engineer` on the
periodic curation pass and to `promote-memory` when the KI is first authored from a retrospective.

## Sync Provenance Fields (optional — set by `sync-memory.sh pull`)

KIs pulled from an org knowledge-hub via ADR-003 sync carry three additional frontmatter fields
stamped automatically by `sync-memory.sh pull`. They are optional (locally-authored KIs never have
them), but when present must conform to the shapes below. Declared in `ki-frontmatter.schema.json`.

| Field | Type | Notes |
|---|---|---|
| `sync_source` | string | Remote URL of the org knowledge-hub (e.g., `git@github.com/org/knowledge-hub`). Presence signals externally-sourced content. Agents must apply `shared/rules/memory-trust-boundary.md` when this field is set. |
| `sync_pulled` | ISO date (`YYYY-MM-DD`) | Date `sync-memory.sh pull` applied this KI. |
| `sync_commit_sha` | string (7–40 hex chars) | Git commit SHA of the org repo at pull time. Enables auditors to trace the exact org-repo state that introduced this KI. |

These fields enable traceability for security audits (THREAT_MODEL.md F-04, F-08) and allow
`memory-auditor` to distinguish locally-authored KIs from externally-sourced ones when assessing
trust level.
