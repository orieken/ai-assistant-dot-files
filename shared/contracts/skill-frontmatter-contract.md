# Contract: skill frontmatter (`shared/skills/*/SKILL.md`)

**Produced by**: humans authoring skill files under `shared/skills/<skill-name>/SKILL.md`
**Consumed by**: `install.sh` (symlinks `shared/skills/` into `.claude/skills/`), Claude Code / Cursor
loaders (registering skills by `name`, routing slash commands like `/adr`), and the auto-trigger routing
table that fires a skill when a user's input matches its `triggers.keywords` or `triggers.intentPatterns`.

This contract governs the YAML frontmatter block at the top of every `SKILL.md` file — not the body. The
body content (Process, Guardrails, Output Format, etc.) is judgment-only; only the frontmatter is
contract-bound because it is what the Claude Code loader and the auto-trigger router grep-parse.

## Required Fields

| Field | Type | Notes |
|---|---|---|
| `name` | string | kebab-case; must match the parent directory name (`shared/skills/adr/SKILL.md` → `name: adr`). Referenced by the `Skill` tool when invoking, and by slash-command routing (`/adr`). |
| `description` | string | One sentence. What the skill does and when Claude should auto-trigger it. Consumed by the Claude Code UI for slash-command discovery and auto-trigger routing. |
| `triggers` | object | See sub-structure below. Has two required sub-keys (`keywords`, `intentPatterns`). The auto-trigger routing table. |
| `standalone` | boolean | Must be `true` in this framework. Skills that require MCP servers or external network calls to function are not portable across installs. Enforces the "everything must work offline" discipline. |

### `triggers` sub-structure

```yaml
triggers:
  keywords: [list, of, single, words, or, short, phrases]
  intentPatterns: ["Natural language sentence pattern *", "Another pattern *"]
```

- `triggers.keywords` — literal string matches; routing fires when any keyword appears in user input.
- `triggers.intentPatterns` — fuzzy sentence patterns with `*` wildcards for user-supplied nouns.

Both sub-keys must be present. An empty list is acceptable if the skill is manual-invocation-only (e.g.,
only fires when the user explicitly types the slash-command), but the key itself must exist so parsers
don't have to distinguish "missing" from "explicitly empty".

## Validation Rule

`validate-artifact` checks:

1. **Field presence** — every required top-level field above (`name`, `description`, `triggers`,
   `standalone`) must appear inside the opening `---` / closing `---` frontmatter block. Missing any
   one is a FAIL. This matches the field-presence check in `scripts/health-check.sh` step 3 (the two
   enforcement paths agree on shape by design — this contract is the referenceable version of what the
   health-check script already enforces).
2. **`triggers` sub-key presence** — the `triggers` object must contain both `keywords` and
   `intentPatterns` as sub-keys. Missing either is a FAIL.
3. **`standalone` value** — must parse to boolean `true` (a trailing YAML comment on the same line is
   fine, per YAML spec — `standalone: true   # must work without MCP/external systems` is valid).
   `false`, `yes`, or missing is a FAIL. This is the "everything must work offline" enforcement — a
   skill that can't declare `standalone: true` doesn't belong in this framework.
4. **`name` shape** — must be lowercase kebab-case matching `^[a-z][a-z0-9-]*$` and must equal the
   parent directory name (`shared/skills/adr/SKILL.md` → `name: adr`).

This is a structural check only. It does not verify that `description` accurately describes the skill,
that the `triggers` actually fire on the intended inputs, or that the skill body is well-formed — those
judgments belong to the human reviewer and to `agent-scorecard` / `pipeline-retrospective` over time.
