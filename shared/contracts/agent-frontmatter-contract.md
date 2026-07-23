# Contract: agent frontmatter (`shared/agents/*.md`)

**Produced by**: humans authoring agent files under `shared/agents/`
**Consumed by**: `install.sh` (symlinks `shared/agents/` into `.claude/agents/`), Claude Code / Cursor
loaders (registering agents by `name`, routing on `description`), and orchestrator agents citing this
agent by role (e.g., `deliver-feature` invoking `analyst`, `developer`, `qa-engineer` by their
frontmatter `name`).

This contract governs the YAML frontmatter block at the top of every agent file — not the body. The body
prompt is judgment-only; only the frontmatter is contract-bound because it is what downstream tooling
grep-parses.

## Required Fields

| Field | Type | Notes |
|---|---|---|
| `name` | string | kebab-case; must match filename base (`analyst.md` → `name: analyst`). Referenced by the `Task` tool's `subagent_type` parameter and by other agents citing this one by role. |
| `description` | string | One sentence. Include `PROACTIVELY` when the pipeline should invoke unconditionally; include `MUST` for hard-required ordering (e.g., "MUST be invoked before developer"). Consumed by the Claude Code / Cursor UI for agent picking and by the orchestrator when deciding pipeline sequencing. |
| `tools` | comma-separated string | Claude Code tool names (`Read`, `Write`, `Edit`, `MultiEdit`, `Bash`, `Glob`, `Grep`). Use least-privilege — read-only agents (auditors, reviewers) should NOT have `Write`. Enforces capability boundaries. |
| `model` | string | Usually `inherit`. Can pin (`claude-opus-4-8`, `claude-sonnet-5`, `claude-haiku-4-5`) for cost/quality tradeoffs on specific agents. `inherit` is the norm because the parent session's model choice usually applies; pin only when there's a real reason. |
| `version` | semver string (`X.Y.Z`) | Bump minor on behavior-relevant change (output-format refactor, tool-list change); patch on prose-only edits; major only if the agent's contract with callers breaks. Downstream installs pull versioned agents; `shared/agents/CHANGELOG.md` tracks bumps. |

## Optional Fields

| Field | Type | Notes |
|---|---|---|
| `isolation` | string | `worktree` — agent runs in a temporary git worktree isolated from the main working copy. See `shared/agents/developer.md` for the reference use. |

## Validation Rule

`validate-artifact` checks:

1. **Field presence** — every required field above (`name`, `description`, `tools`, `model`, `version`)
   must appear as a top-level YAML key inside the opening `---` / closing `---` frontmatter block.
   Missing any one is a FAIL. This matches the field-presence check in `scripts/health-check.sh` step 2
   (the two enforcement paths agree on shape by design — this contract is the referenceable version of
   what the health-check script already enforces).
2. **`version` format** — the `version` value must be a valid semver string matching
   `^[0-9]+\.[0-9]+\.[0-9]+$`. `1.2.0` passes; `1.2` and `v1.2.0` do not.
3. **`name` shape** — must be lowercase kebab-case matching `^[a-z][a-z0-9-]*$` and must equal the
   filename base (`analyst.md` → `name: analyst`).

This is a structural check only. It does not verify that `description` accurately describes the agent,
that `tools` matches the tool set the body actually uses, or that `version` was bumped when it should
have been — those judgments belong to the human reviewer and to `memory-engineer` / `agent-scorecard`
over time.

Agent files that don't declare frontmatter at all (e.g., `shared/agents/CHANGELOG.md`, which is the
changelog itself, not an agent) are excluded — the validator skips any file listed in
`scripts/health-check.sh`'s exclusion list (`CHANGELOG.md`).
