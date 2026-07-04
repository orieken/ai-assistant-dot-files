# Architecture Notes: Context Engineering Framework

## Structural Decisions

### Decision 1: Canonical `shared/` layer as single source of truth
**Decision**: All rules, agents, skills, and the domain dictionary live in `shared/`. Platform-specific directories (`.claude/`, `.cursor/`, `.github/`, `.gemini/`, `.openai/`, `.windsurf/`) are generated outputs, not sources.
**Reversibility**: Moderate — requires moving files and updating all references, but no data loss.
**Fitness function**: `scripts/check-parity.sh` — diffs generated platform configs against `shared/` canonical content. Fails on any drift.
**Enforcement**: CI runs `check-parity.sh` on every PR.

### Decision 2: Platform capability tiers
**Decision**: Platforms are classified into three capability tiers:
- **Tier 1 (Full)**: Claude Code — supports agents, skills, rules, hooks, sub-agent orchestration
- **Tier 2 (Personas + Rules)**: Cursor, Windsurf — supports persona-level context shaping and rule files, but no native agent orchestration
- **Tier 3 (System Prompt)**: Copilot, Gemini, OpenAI — single instruction file only

Generated configs match the platform's tier. Tier 2/3 platforms inline the agent roster for awareness but cannot orchestrate agents natively.
**Reversibility**: Cheap — tiers can be upgraded as platforms add features.
**Fitness function**: `platform-registry.json` declares each platform's tier; `generate-configs.sh` respects it.
**Enforcement**: Code review + registry schema validation.

### Decision 6: Cursor requires fully inlined rules (no file references)
**Decision**: Cursor's LLM cannot dynamically traverse the file tree from `.mdc` rules. All instructions must be written directly in the `.mdc` markdown body. `generate-configs.sh` must concatenate and inline shared rules into each `.mdc` file. Best practices: keep each `.mdc` short and focused (one per concern), use ALWAYS/NEVER/CRITICAL for compliance, break the agent roster into separate `.mdc` files, and ensure valid YAML frontmatter (Cursor silently ignores malformed files).
**Reversibility**: Cheap — if Cursor adds file-reference support later, switch to references.
**Fitness function**: `check-parity.sh` verifies each `.mdc` file's content matches the corresponding `shared/` source.
**Enforcement**: CI + config generation script.

### Decision 3: Symlinks for global install, copies for project install
**Decision**: `install.sh --global` symlinks to `~/.claude/`, `~/.cursor/`, etc. `install.sh --project <path>` copies files to the target project. Symlinks keep global configs always current; copies let projects pin a version.
**Reversibility**: Cheap — symlinks are trivially removed, copies are just files.
**Fitness function**: `health-check` skill verifies symlinks resolve and copies match `shared/` checksums.
**Enforcement**: `install.sh` post-install verification.

### Decision 4: File-based agent tracing (no external dependencies)
**Decision**: Pipeline observability uses a JSON trace file (`pipeline-trace.json`) persisted alongside delivery artifacts. No OpenTelemetry collector, no external service. Structured enough to query with `jq`, simple enough to work offline.
**Reversibility**: Cheap — trace format can be upgraded to OTel-compatible later.
**Fitness function**: `deliver-feature` pipeline must produce a `pipeline-trace.json` for every run.
**Enforcement**: Pipeline guardrail — delivery summary step fails if trace file is missing.

### Decision 5: Golden-file agent tests use fuzzy matching
**Decision**: Agent evaluation tests check for the presence of key findings (phrases, section headings, severity levels) rather than exact text output. This makes tests resilient to prompt wording changes while still catching functional regressions.
**Reversibility**: Cheap — test assertions are just grep patterns.
**Fitness function**: `scripts/test-agents.sh` must pass before any agent file is merged.
**Enforcement**: CI + pre-merge check.

---

## Component Placement

| Component | Location | Layer | Purpose |
|---|---|---|---|
| Canonical rules | `shared/rules/` | Governance | Single source of truth for all constraints |
| Canonical agents | `shared/agents/` | Orchestration | Platform-agnostic agent definitions |
| Canonical skills | `shared/skills/` | Orchestration | Platform-agnostic skill definitions |
| Platform registry | `shared/platform-registry.json` | Infrastructure | Declares platforms, paths, tiers |
| Config generator | `scripts/generate-configs.sh` | Infrastructure | Produces platform-specific configs |
| Parity checker | `scripts/check-parity.sh` | Infrastructure | Fitness function for config drift |
| Agent tests | `tests/agents/<name>/` | Verification | Golden-file regression tests |
| Agent test runner | `scripts/test-agents.sh` | Infrastructure | Runs all agent golden-file tests |
| Inter-agent contracts | `shared/contracts/` | Orchestration | Required artifact schemas |
| Artifact validator | `shared/skills/validate-artifact/` | Orchestration | Checks artifacts against contracts |
| Pipeline tracer | `shared/skills/pipeline-trace/` | Observability | Logs agent invocation metrics |
| Pipeline retro | `shared/skills/pipeline-retrospective/` | Feedback | Analyzes past pipeline traces |
| Install script | `install.sh` | Infrastructure | One-command setup |
| Uninstall script | `uninstall.sh` | Infrastructure | Clean teardown |
| Agent changelog | `shared/agents/CHANGELOG.md` | Documentation | Tracks agent prompt evolution |

---

## Proposed Directory Structure

```
ai-assistant-dot-files/
├── shared/                          ← CANONICAL SOURCE OF TRUTH
│   ├── rules/
│   │   ├── architecture-guardrails.md
│   │   ├── design-principles.md
│   │   └── approval-gates.md
│   ├── agents/
│   │   ├── CHANGELOG.md
│   │   ├── analyst.md
│   │   ├── architect.md
│   │   ├── developer.md
│   │   ├── code-reviewer.md
│   │   ├── security-reviewer.md
│   │   └── ... (all 20+ agents)
│   ├── skills/
│   │   ├── deliver-feature/SKILL.md
│   │   ├── validate-artifact/SKILL.md  ← NEW
│   │   ├── pipeline-trace/SKILL.md     ← NEW
│   │   ├── pipeline-retrospective/SKILL.md ← NEW
│   │   └── ... (all 35+ skills)
│   ├── contracts/                      ← NEW
│   │   ├── analysis-contract.md
│   │   ├── architecture-contract.md
│   │   ├── implementation-contract.md
│   │   └── review-contract.md
│   ├── platform-registry.json          ← NEW
│   ├── domain-dictionary.md
│   └── architecture-rules.md
│
├── .claude/                         ← GENERATED (Tier 1: full)
│   ├── agents/ → symlink to shared/agents/
│   ├── skills/ → symlink to shared/skills/
│   ├── rules/ → symlink to shared/rules/
│   └── settings.json
│
├── .cursor/                         ← GENERATED (Tier 2: personas + rules)
│   └── rules/
│       └── global.mdc
│
├── .github/                         ← GENERATED (Tier 3: system prompt)
│   └── copilot-instructions.md
│
├── .gemini/                         ← GENERATED (Tier 3: system prompt)
│   └── antigravity/
│       └── instructions.md
│
├── .openai.md                       ← GENERATED (Tier 3: system prompt)
├── .cursorrules                     ← GENERATED (Tier 2: alias)
├── .windsurfrules                   ← GENERATED (Tier 2: alias)
│
├── tests/                           ← NEW
│   └── agents/
│       ├── security-reviewer/
│       │   ├── input.md             (known-vulnerable code)
│       │   └── expected.txt         (grep patterns for expected findings)
│       ├── code-reviewer/
│       └── analyst/
│
├── scripts/
│   ├── generate-configs.sh          ← NEW
│   ├── check-parity.sh             ← EXISTS (extend)
│   └── test-agents.sh              ← NEW
│
├── templates/
│   └── claude-feature-team/         ← EXISTS (extend for multi-platform)
│
├── install.sh                       ← NEW
├── uninstall.sh                     ← NEW
├── scaffold-team.sh                 ← EXISTS (wrap install.sh)
│
├── CLAUDE.md
├── ARCHITECTURE_RULES.md
├── DOMAIN_DICTIONARY.md
├── README.md
└── docs/
    ├── features/
    ├── adrs/
    └── runbooks/
```

---

## Implementation Priority

### Phase 1: Foundation (do first — everything else depends on this)
1. Epic 1 — Canonical `shared/` layer
2. Epic 9 — Domain dictionary update (persona vs. agent)
3. Epic 3 — Install/uninstall scripts

### Phase 2: Cross-platform (the main deliverable)
4. Epic 2 — Config generation + parity tests

### Phase 3: Pipeline hardening
5. Epic 4 — Wire context-engineer into pipeline
6. Epic 5 — Inter-agent contracts

### Phase 4: Quality & observability
7. Epic 6 — Agent golden-file tests
8. Epic 7 — Agent observability & feedback loop
9. Epic 8 — Agent versioning & changelog

### Phase 5: Polish
10. Epic 10 — Health check & self-test

---

## Anti-Pattern Check

- [x] **Checked for distributed monolith**: Not applicable — single repo, file-based orchestration
- [x] **Checked for anemic domain model**: Agents have rich process definitions, not just data
- [x] **Checked for God object**: `deliver-feature` SKILL.md is complex (153 lines, 22 steps) but appropriately so — it's the orchestrator. No extraction needed.
- [x] **Checked for shotgun surgery**: Adding a new platform currently requires editing 6+ files. The `shared/` + `generate-configs.sh` pattern reduces this to 1 registry entry + 1 template.
- [x] **Checked for leaky abstraction**: Platform configs currently expose Claude-specific paths (`.claude/agents/`). Generated configs should reference concepts ("the architect agent") not paths.
- [x] **Checked for premature generalization**: Platform registry is the right level of abstraction — it captures real variation (3 tiers) without over-engineering.

---

## Developer Handoff Notes

- `scaffold-team.sh` is the existing entry point for Claude-only setup. It should continue to work unchanged but internally delegate to `install.sh --project --platform claude`.
- The `archive/` directory contains old versions of platform configs. These are useful for understanding format evolution but should not be referenced by any new code.
- `check-parity.sh` already exists at the repo root. Extend it rather than creating a new script.
- The `.claude/settings.json` hook (`afterEdit: npm test --changed`) is project-specific and should NOT be included in the shared layer — it belongs in the project template.

---

## Open Architectural Questions

- **Q1**: Should `install.sh --global` symlink the entire `shared/` directory, or only the platform-specific generated config? Symlink-all is simpler but may cause confusion if the user edits a file thinking it's local.
- **Q2**: For Cursor's `.mdc` format — should the generated file include the full agent roster inline, or reference `shared/agents/` with a "read these files" instruction? Cursor may or may not follow file references.
- **Q3**: Should agent golden-file tests run against the actual LLM (expensive, non-deterministic) or against a mocked response with structural checks only (cheap, deterministic)? Recommendation: structural checks for CI, LLM tests as a manual verification step.
