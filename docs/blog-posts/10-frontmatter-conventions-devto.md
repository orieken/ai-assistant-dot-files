---
title: "Frontmatter Conventions Nobody Agrees On: Managing Multi-Platform AI Metadata"
published: false
description: "How we unified divergent agent, skill, and rule metadata across 6 AI coding platforms using canonical JSON schemas."
tags: ai, architecture, devtools, markdown, software-engineering
canonical_url:
cover_image:
---

YAML frontmatter has become the universal metadata header for markdown-based AI agents, skills, and rule files.

There is just one problem: **no two AI coding tools agree on the frontmatter schema**.

- **Claude Code** expects `name`, `description`, `tools`, `model`, and `version` inside `.claude/agents/*.md`.
- **Cursor** expects `description`, `globs`, and `alwaysApply` inside `.cursor/rules/*.mdc`.
- **Windsurf** and **GitHub Copilot** ignore YAML frontmatter entirely and expect flat system instructions.
- **Custom Agent Frameworks** introduce custom keys like `triggers`, `standalone`, or `domain`.

If you manage a team across multiple IDEs and CLI tools, attempting to maintain separate prompt files for each platform leads to massive configuration drift.

Here is how we solved multi-platform frontmatter fragmentation in `ai-assistant-dot-files` using **Canonical Schemas and Automated Platform Projection**.

---

## The Core Strategy: One Canonical Source, Six Projections

Instead of authoring tool-specific configurations by hand, our repository maintains one canonical definition layer in `shared/`:

- `shared/agents/*.md`: Canonical agent definition files.
- `shared/skills/*/SKILL.md`: Canonical skill instructions with YAML frontmatter.
- `shared/rules/*.md`: Canonical architectural rules.

We then use a generator script (`scripts/generate-configs.sh`) driven by `shared/platform-registry.json` to automatically project those canonical files into six target AI tool environments:

1. **Claude Code** (`.claude/agents/`, `.claude/skills/`)
2. **Cursor** (`.cursor/rules/*.mdc`, `.cursor/agents/`, `.cursor/skills/`)
3. **Windsurf** (`.windsurfrules`)
4. **GitHub Copilot** (`.github/copilot-instructions.md`, `.github/instructions/`)
5. **Gemini / Antigravity** (`AGENTS.md`)
6. **OpenAI Codex** (`.openai.md`)

---

## Schema-First Validation: JSON Schemas for Markdown Frontmatter

To prevent invalid frontmatter keys from breaking platform generators, we authored explicit JSON Schemas in `shared/schemas/`:

- `agent-frontmatter.schema.json`
- `skill-frontmatter.schema.json`
- `ki-frontmatter.schema.json`

During our automated health check (`scripts/health-check.sh`), PyYAML parses the frontmatter blocks and validates them against the JSON Schemas:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["name", "description", "tools", "version"],
  "properties": {
    "name": { "type": "string", "pattern": "^[a-z0-9-]+$" },
    "description": { "type": "string", "minLength": 10 },
    "tools": { "type": "string" },
    "version": { "type": "string", "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$" }
  }
}
```

If an author misspells `description` as `desc` or adds an invalid tool string, the health check fails immediately before platform configs are built.

---

## Zero-Drift Verification: `check-parity.sh`

To guarantee that no platform config drifts out of sync with `shared/`, our CI pipeline runs `scripts/check-parity.sh`.

It cross-references every projected file (`.cursorrules`, `.windsurfrules`, `.github/copilot-instructions.md`, `AGENTS.md`) against the canonical `shared/` definitions. If someone edits `.cursorrules` directly without modifying `shared/`, `check-parity.sh` flags the drift and blocks the build.

---

## Takeaways for AI System Architects

1. **Never make platform-specific rule files your source of truth.**
2. **Author canonical markdown artifacts in a central `shared/` directory.**
3. **Enforce YAML frontmatter schemas using JSON Schema validation.**
4. **Automate projection to Cursor, Copilot, Windsurf, and Claude Code.**

---

## Image Prompt

> Hero image prompt: An isometric technical diagram showing a single central golden blueprint document branching into 6 stylized glowing streams, projecting into 6 distinct platform icon displays (VS Code, Cursor, Terminal, Copilot). Clean dark mode background, vibrant emerald, amber, and cyan neon accents.
