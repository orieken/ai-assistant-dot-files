# Add JSON schemas for agent / skill / KI frontmatter

Create `shared/schemas/*.schema.json` for each frontmatter shape. Unlocks IDE autocomplete (VS Code, Cursor) via `yaml.schemas` settings, and enables enum-value validation that `scripts/health-check.sh` currently can't do (today it checks field presence only, not values).

Depends on: `add-frontmatter-contracts.md` — the schemas should match the contracts. Do the contracts first for clean sourcing.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` — this IS the git repo.

## Prior context

The reference for what fields exist and what values are valid is `docs/patterns/frontmatter-conventions.md`. Read it before writing schemas.

Today's enforcement (`scripts/health-check.sh`) checks field PRESENCE only — so `tools: WhatEverRandomName` would pass. Real IDE-level validation with autocomplete needs JSON Schema.

Only existing schema in the repo: `shared/telemetry/event-schema.md` (added in AOS Phase 1). That's a Markdown-documented schema, not a JSON-Schema file. Different pattern.

## Scope

### Op 1: create `shared/schemas/agent-frontmatter.schema.json`

JSON Schema (draft-07 or newer) for the agent frontmatter shape. Constraints:
- `name` — string, kebab-case pattern (`^[a-z][a-z0-9-]*$`)
- `description` — string, min length ~30 chars (avoids empty descriptions)
- `tools` — string, comma-separated enum of valid Claude Code tools (Read, Write, Edit, MultiEdit, Bash, Glob, Grep, Task, Skill, WebFetch, WebSearch — see the current tool set)
- `model` — enum: `inherit`, `claude-opus-4-8`, `claude-sonnet-5`, `claude-haiku-4-5-20251001`, `claude-opus-4-7` (verify against `docs/CLAUDE.md`'s "most recent Claude models" line before hardcoding — model list may need updating)
- `version` — string, semver pattern (`^\d+\.\d+\.\d+$`)
- `isolation` — optional enum: `worktree`

### Op 2: create `shared/schemas/skill-frontmatter.schema.json`

Constraints:
- `name` — string, kebab-case
- `description` — string, min length
- `triggers` — object with required `keywords` (array of strings) and `intentPatterns` (array of strings)
- `standalone` — boolean, must be `true` (enum with single value — enforces the "everything works offline" discipline)

### Op 3: create `shared/schemas/ki-frontmatter.schema.json`

Constraints:
- `name` — string, kebab-case
- `tags` — array of strings, each kebab-case
- `domain` — string, kebab-case
- `created` — string, ISO date format (`^\d{4}-\d{2}-\d{2}$`)

### Op 4: wire schemas into VS Code / Cursor settings templates

Create or update `.vscode/settings.json.example` and `.cursor/settings.json.example` with:

```json
{
  "yaml.schemas": {
    "./shared/schemas/agent-frontmatter.schema.json": ["shared/agents/*.md"],
    "./shared/schemas/skill-frontmatter.schema.json": ["shared/skills/*/SKILL.md"],
    "./shared/schemas/ki-frontmatter.schema.json": ["shared/knowledge/*.md", ".claude/knowledge/*.md"]
  }
}
```

Add a section to `README.md` (or `docs/runbooks/`) explaining how contributors opt in.

### Op 5: (optional) extend `scripts/health-check.sh` to validate against schemas

If a shell-friendly JSON Schema validator is available (e.g., `check-jsonschema`, `ajv-cli`, `python -m jsonschema`), add a validation step. If none is portable/installable-free enough for the repo's zero-dependency goal, skip and note as a follow-up.

### Op 6: update `docs/patterns/frontmatter-conventions.md`

Mark the "JSON schemas" bullet in the Gaps section as done. Cross-reference the new schema files.

## Discipline

- One commit per op ideally, 5-6 total.
- Conventional Commits: `feat(schemas): add agent-frontmatter JSON schema`, etc.
- **NEVER `git add -A`** — explicit paths only.
- Verify one existing agent/skill/KI validates against its schema as a sanity check before committing each schema.

## Escalation criteria

Stop and report if:
- An existing frontmatter file fails validation under the schema — either loosen the schema or fix the offender first (real drift the schema exposes)
- The model enum in Op 1 can't be pinned reliably (models release frequently) — recommend removing the enum in favor of a `pattern` matching `inherit` OR any `claude-*` id, with a comment
- No portable shell JSON Schema validator exists for Op 5 — mark it as a follow-up rather than adding a heavy install dependency

## Report format (under 200 words)

```
STATUS: complete | stopped-at-op-N
Commits: <sha> <message>, ...

Schemas created:
- shared/schemas/agent-frontmatter.schema.json
- shared/schemas/skill-frontmatter.schema.json
- shared/schemas/ki-frontmatter.schema.json

IDE templates:
- .vscode/settings.json.example (created / not created)
- .cursor/settings.json.example (created / not created)

Health-check schema validation:
- integrated (via <tool>) / skipped (reason: <why>)

Any existing frontmatter that fails the schema (should be zero if contracts were done first):
- <list, or "None">
```

Go.
