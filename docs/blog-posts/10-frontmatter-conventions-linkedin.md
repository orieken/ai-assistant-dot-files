YAML frontmatter is the standard header for AI agent, skill, and rule files.

The problem: Cursor, Claude Code, Windsurf, and Copilot all expect different frontmatter keys—or ignore frontmatter entirely.

In `ai-assistant-dot-files`, we solved multi-platform config drift with a **Canonical Source + Automated Projection** model:

1. Maintain canonical agents, skills, and rules in `shared/`.
2. Validate frontmatter against JSON Schemas (`shared/schemas/*.schema.json`).
3. Project configurations out to 6 AI platforms via `scripts/generate-configs.sh`.
4. Enforce zero drift in CI with `scripts/check-parity.sh`.

Never let platform-specific config files become your source of truth.

Full technical breakdown: TODO_DEVTO_URL
