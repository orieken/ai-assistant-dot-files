# Epic 44 — JetBrains / Junie System Prompt Exporter

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 1.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

JetBrains AI Assistant (with Junie as the newer agentic mode) is a real platform with growing adoption. Framework doesn't export to it today. Config location and format need investigation — JetBrains stores AI settings in per-IDE `<Config Dir>/AIAssistant/` typically, and Junie has its own `.junie/guidelines.md` convention.

## Scope

**Phase A — Investigate** (one commit): read current JetBrains AI Assistant + Junie docs. Confirm the file paths, format shape, tier capabilities (rules only? agents? skills?).

Draft `shared/platform-registry.json` entry with Tier assignment + config paths. Draft the generator format spec.

Commit as: `docs(platform-registry): investigate JetBrains + Junie integration (Epic 44 Phase A)`.

**Pause for user approval before Phase B.**

**Phase B — Implementation** (multiple commits):

1. Add JetBrains + Junie entries to `shared/platform-registry.json`
2. Extend `scripts/generate-configs.sh` to produce the JetBrains + Junie config files
3. Extend `scripts/check-parity.sh` for drift detection on the new files
4. Update `install.sh` for `--platform jetbrains` and `--platform junie`
5. Update README.md's supported-tools table
6. Verify: generate configs, install into a scratch project, validate JetBrains actually reads them

Commit per file, matching Epic 42's shape.

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- JetBrains AI Assistant and Junie may be a single tier or two separate ones — halt if unclear, propose.
- If JetBrains config is fundamentally per-IDE (IntelliJ ≠ WebStorm ≠ Rider) — halt, describe whether the framework generates per-IDE variants or a shared file.
- If Junie's `.junie/guidelines.md` and JetBrains AI Assistant have completely different capability tiers — halt, describe.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A findings:
  - JetBrains AI Assistant config path + format: <describe>
  - Junie config path + format: <describe>
  - Tier assignments proposed: <values>
  - Per-IDE differentiation needed: yes | no

Phase B commits (if approved):
  <sha> <message>
  ...

Post-B verification: generation valid, install works, JetBrains actually reads the generated file.
```

Go.
