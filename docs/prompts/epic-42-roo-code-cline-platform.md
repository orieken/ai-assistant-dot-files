# Epic 42 — Roo Code / Cline Multi-Mode Platform Integration

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 1 (2nd priority per audit).

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

`shared/platform-registry.json` documents 6 platform targets (Claude Code, Cursor, Windsurf, GitHub Copilot, Gemini/Antigravity, OpenAI Codex). Missing: Roo Code (fork of Cline with custom modes) + Cline (VS Code MCP-native agent). Both have growing user bases and their own config file formats.

Investigation needed first — Roo Code and Cline configs may have shifted since the audit was written. Read their current docs before assuming the format.

## Scope

**Phase A — Investigate + design** (one commit): read Roo Code and Cline current docs to confirm the config file shapes. Draft:
- `shared/platform-registry.json` entry proposal for each platform (Tier assignment, capabilities, config paths)
- Format spec for `.roomodes` and `.clinerules` files based on what those tools currently expect
- Which of the framework's agents/skills/rules translate to their format vs. get dropped

Commit as: `docs(platform-registry): investigate Roo Code + Cline integration (Epic 42 Phase A)`.

**Pause for user approval before Phase B.**

**Phase B — Implementation** (multiple commits):

1. Add Roo Code + Cline entries to `shared/platform-registry.json`
2. Extend `scripts/generate-configs.sh` to produce `.roomodes` and `.clinerules` from the shared sources
3. Extend `scripts/check-parity.sh` to include the two new platforms in drift detection
4. Update `install.sh` to accept `--platform roo-code` and `--platform cline` flags
5. Update `README.md` to list the two new platforms in the supported-tools table
6. Verify: run generation, confirm output files match the expected format, install into a scratch project + validate

Commit per file: `feat(platform-registry): add Roo Code entry (Epic 42 Phase B Op N)`, etc.

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- If Roo Code or Cline's config format is significantly more complex than existing platform formats — halt, describe the delta. May need a new pattern in `scripts/generate-configs.sh`.
- If the two platforms have overlapping/conflicting conventions (e.g., both claim `.rules` as their config) — halt, describe.
- If Tier assignment is unclear (Tier 1 = full agents+skills+rules, Tier 2 = some subset) — halt, propose a defensible tier and rationale.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A findings:
  - Roo Code config format: <describe>
  - Cline config format: <describe>
  - Tier assignments proposed: <Roo=X, Cline=Y>
  - Config-path collision: yes | no

Phase B commits (if approved):
  <sha> <message>
  ...

Post-B verification: generation produces valid configs, install works on scratch project, parity check passes.
```

Go.
