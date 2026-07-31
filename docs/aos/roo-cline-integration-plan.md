# Roo Code + Cline Integration Plan (Epic 42)

Drafted 2026-07-30 from live doc research. Confirms config formats and tier assignments before Phase B implementation.

## Investigation Findings

### Roo Code

- **Config file**: `.roomodes` at workspace root (YAML or JSON; YAML preferred)
- **Global config**: `~/.roo/settings/custom_modes.yaml`
- **Per-mode rules**: `.roo/rules-{slug}/` directory (project-level)
- **Global rules**: `.roo/rules/` directory (all modes)
- **Schema source**: `https://docs.roocode.com/schemas/custom-modes.schema.json`

**Mode structure:**
```yaml
customModes:
  - slug: analyst                  # kebab-case, [a-zA-Z0-9-]+ only
    name: Analyst                  # display name
    description: "Short UI summary"
    roleDefinition: |              # injected as system prompt for this mode
      Full prose body of the agent
    whenToUse: "Automated decision guidance"
    customInstructions: ""         # optional extra instructions (appended at end)
    groups:
      - read                       # string = unrestricted
      - [edit, {fileRegex: ".*"}]  # tuple = restricted by regex
      - command
      - mcp
      - browser
```

**Tool group mapping from shared/agents `tools:` field:**

| `tools:` values | Roo group |
|---|---|
| `Read`, `Glob`, `Grep`, `WebFetch`, `WebSearch` | `read` |
| `Write`, `Edit`, `MultiEdit`, `NotebookEdit` | `edit` |
| `Bash` | `command` |
| `Agent`, `TaskCreate`, `TaskUpdate`, `Artifact`, `mcp__*` | `mcp` |
| `*` (all tools) | all five groups |

Default for any agent without a tools field: `read` only.

### Cline

- **Config file**: `.clinerules/` directory containing `.md` and `.txt` files
- **Also accepted**: single `.clinerules` file at root
- **Global config**: `~/Documents/Cline/Rules/` (macOS/Linux)
- **Cross-tool**: also reads `~/.agents/AGENTS.md` (same convention as agents.md standard)
- **Optional frontmatter**: `paths:` list of glob patterns for conditional activation

**File structure:**
```
.clinerules/
├── 00-approval-gates.md          # always-active
├── 01-design-principles.md       # always-active
├── 02-architecture-guardrails.md # always-active
├── 03-agent-roster.md            # always-active
├── 04-testing-conventions.md     # paths: [**/*.spec.*, ...]
├── 05-go-conventions.md          # paths: [**/*.go]
├── 06-typescript-conventions.md  # paths: [**/*.ts]
├── 07-python-conventions.md      # paths: [**/*.py]
├── 08-csharp-conventions.md      # paths: [**/*.cs]
└── 09-java-conventions.md        # paths: [**/*.java]
```

No custom modes, no agent-equivalent capability.

## Tier Assignments

| Platform | Tier | Label | Rationale |
|---|---|---|---|
| Roo Code | 2 | `Personas + Modes` | Custom modes give per-mode tool access scoping (agent-like) but no sub-agent spawning, no skill invocation, no hooks. Closer to Tier 1 than Windsurf/Copilot but not full orchestration. `capabilities.agents: true` via modes. |
| Cline | 2 | `Personas + Rules` | Plain markdown rules, optional path scoping. No modes, no agent concept. Identical capability level to Windsurf. |

## Config Path Collision

None. Paths are fully distinct:
- Roo Code: `.roomodes`, `.roo/rules/`
- Cline: `.clinerules/`

## Translation Strategy

### Roo Code — agents → custom modes

Every `shared/agents/*.md` file maps to one Roo Code custom mode. The agent body (after frontmatter) becomes `roleDefinition`. The `description` frontmatter field maps to both `description` and `whenToUse`. Tool access is inferred from the `tools:` frontmatter field using the mapping table above.

Global framework rules (`design-principles.md`, `architecture-guardrails.md`, `approval-gates.md`) are placed in `.roo/rules/` so they apply to ALL modes — agents already pre-read these in their process instructions, so this avoids duplication.

CHANGELOG.md is excluded (not an agent definition).

### Cline — rules directory

Each `shared/rules/*.md` file becomes a numbered file in `.clinerules/` with an appropriate `paths:` frontmatter where the rule is language/test-specific. Always-active rules (approval-gates, design-principles, architecture-guardrails) get no paths restriction. The agent roster is appended as a separate file.

Skills are not translated — Cline has no skill invocation mechanism.
Hooks are not translated — Cline has no hooks.

## What Gets Dropped

| Framework concept | Roo Code | Cline |
|---|---|---|
| Skills (slash commands) | Not supported | Not supported |
| Hooks | Not supported | Not supported |
| Inter-agent contracts | Not applicable | Not applicable |
| Pipeline orchestration | Not applicable | Not applicable |
| `tools:` fine-grained allowlist | Approximated via groups | Not applicable |

## Phase B Plan

- B1: Add Roo Code + Cline to `shared/platform-registry.json`
- B2: Add `scripts/generate-roomodes.py` + `generate_roo_code()` in `generate-configs.sh`
- B3: Add `generate_cline()` in `generate-configs.sh`
- B4: Add Roo Code + Cline checks to `scripts/check-parity.sh`
- B5: Add `install_roo_code()` + `install_cline()` + detection to `install.sh`
- B6: Update `README.md` supported-tools table
