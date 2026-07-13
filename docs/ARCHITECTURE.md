# Architecture

This document explains how `shared/` becomes six different platform configs, and how context flows through
the agent pipeline once a feature is being delivered. For "how do I add X," see
[CONTRIBUTING.md](CONTRIBUTING.md). For a narrative walkthrough, see the [README](../README.md). For what
each individual agent does and what checks its work, see [AGENT_REFERENCE.md](AGENT_REFERENCE.md).

---

## 1. The `shared/` Canonical Layer

Every agent, skill, and rule is authored exactly once, in `shared/`:

```
shared/
├── agents/          25 agents — .md, YAML frontmatter (name, description, tools, model, version)
├── skills/          54 skills — .md, YAML frontmatter (name, description, triggers)
├── rules/           architecture-guardrails.md, design-principles.md, approval-gates.md
├── contracts/       required-section contracts for pipeline agent handoffs (Epic 5)
├── knowledge/       portable Knowledge Items (KIs) — searchable via search-ki / query-memory
├── ARCHITECTURE_RULES.md
├── DOMAIN_DICTIONARY.md
├── TEAM_TOPOLOGY.md — Bounded Context -> team/type/interaction-mode registry (Skelton & Pais), checked by architect and team-topology-check
├── memory-registry.json — catalog of every durable memory source + retrieval backend, checked by search-ki/query-memory/memory-engineer
└── platform-registry.json
```

Nothing outside `shared/` is a source of truth. `.claude/agents/`, `.claude/skills/`, and `.claude/rules/`
are symlinks to their `shared/` equivalents (`ln -s ../shared/agents .claude/agents`, etc.) — editing through
the symlink and editing the canonical file are the same operation for Claude Code. Every other platform's
config is *generated*, not symlinked, because none of them can follow file references the way Claude Code
can (Cursor's own docs are explicit about this: `.mdc` rules must be fully self-contained).

### Why a single canonical layer
Before this structure existed, the same instructions were hand-copied into `.cursorrules`,
`copilot-instructions.md`, and `CLAUDE.md` independently, and they drifted every time one was edited without
the others. `scripts/check-parity.sh` exists specifically to catch that drift — it diffs every generated
config against `shared/` and fails if any platform is missing an agent, a rule concept, or has fallen out of
sync.

---

## 2. The Capability Tier System

Not every AI tool can do the same things. `shared/platform-registry.json` classifies each platform into one
of three tiers, and `DOMAIN_DICTIONARY.md` defines exactly what each tier means:

| Tier | Label | Capabilities | Platforms | Terminology |
|---|---|---|---|---|
| **1** | Full | Agents with tool access, autonomous multi-step process, pipeline participation, hooks | Claude Code | **Agent** |
| **2** | Personas + Rules | Persona-level context shaping + rule files, no native orchestration | Windsurf, GitHub Copilot | **Persona** |
| **2*** | Personas + Rules, agents/skills now Tier-1-equivalent | Real subagent + skill loading (confirmed 2026-07-06) via `.cursor/agents/`/`.cursor/skills/`, plus rule files (still inlined, no orchestration at the rules layer) | Cursor | **Agent** for `.cursor/agents/`, **Persona** for `.cursor/rules/` |
| **3** | System Prompt (rules); real skill invocation confirmed on top | Single rules file (`AGENTS.md`), plus genuine skill execution — not just description | Gemini/Antigravity, OpenAI/Codex | **Persona** |

Cursor is the one platform whose tier number alone no longer tells the whole story: it shipped native
Agent Skills (`.cursor/skills/*/SKILL.md`) and subagents (`.cursor/agents/*.md`) using the same open
standard `shared/agents/`/`shared/skills/` already follow (confirmed 2026-07-06 against
`cursor.com/docs/subagents`, `cursor.com/docs/skills`, and a live check in this repo — the analyst subagent
and search-ki skill both loaded and behaved correctly, not just generically). `install.sh` now symlinks
`.cursor/agents/`/`.cursor/skills/` directly, the same zero-drift mechanism Claude Code has always used —
see `shared/platform-registry.json`'s `capabilities` object, which is what's actually authoritative per
capability now, not the single `tier` number. Cursor Rules (`.mdc`) are unaffected by this — still fully
inlined, still no orchestration at that layer, still genuinely Tier 2. One real, permanent capability gap
versus Claude Code: Cursor subagents have no `tools:` allowlist field — they inherit *all* of the parent's
tools (MCP tools included) with only a coarse `readonly: true/false` available, so `shared/agents/*.md`'s
`tools:` frontmatter is simply ignored by Cursor's parser.

Copilot moved from Tier 3 to Tier 2 in 2026-07 after confirming (via GitHub's own docs) that it supports
path-scoped `.github/instructions/*.instructions.md` files alongside the repo-wide instructions file — the
same "multiple rule files, no orchestration" shape Windsurf already had.

Gemini/Antigravity was live-tested 2026-07-02 (see `tests/platform-verification/antigravity.md` and its
results file) rather than left on secondary-source guesswork: it reads `AGENTS.md` for rules (confirmed —
asking it to list approval gates returned an exact match against `shared/rules/approval-gates.md`), and it
genuinely *invokes* skills rather than just describing them (asking it to run `complexity-check` against a
fixture correctly applied the real thresholds). The framework's previous best guess —
`.gemini/antigravity/instructions.md` — was confirmed **not** read at all and has been removed. Skills
loaded from `~/.gemini/config/skills/` (the global root) in this test since project-level `.agents/skills/`
didn't exist yet at session start; that project-level path itself remains unconfirmed, not contradicted.

The distinction matters because it's enforced in the generated output, not just documented: `Persona` is a
context frame with no tool access and no autonomous workflow (see `DOMAIN_DICTIONARY.md`'s Entity table) —
Tier 2/3 configs consistently say "persona" in generated roster text (`collect_agent_roster()` in
`scripts/generate-configs.sh`), never "agent." Only Tier 1 gets the word "agent," because only Tier 1
actually runs multi-step orchestration with tool access.

### Generation strategy per tier
- **Tier 1 (Claude Code)**: symlink. `install.sh` creates `.claude/{agents,rules,skills}` -> `shared/`
  equivalents. Always current after a `git pull`; no generation step needed.
- **Cursor (mixed strategy)**: symlink for agents/skills, generate-inline for rules.
  - `install.sh`'s `install_cursor()` symlinks `.cursor/agents` -> `shared/agents` and `.cursor/skills` ->
    `shared/skills` directly — same zero-drift mechanism as Tier 1, confirmed working 2026-07-06. This
    retired the earlier `generate_cursor_personas()` workaround (built in Epic 11, before Cursor could do
    real skill/agent loading), which used to flatten each agent into a standalone `.cursor/rules/<name>.mdc`
    persona file.
  - Rules still can't follow file references (no evidence Cursor Rules support them, unlike agents/skills),
    so `generate-configs.sh` still generates 11 `.mdc` files for rules: `architecture.mdc`,
    `design-principles.mdc`, `agent-roster.mdc`, `approval-gates.mdc`, `testing.mdc`, `go-backend.mdc`,
    `vue-frontend.mdc`, plus `typescript-conventions.mdc`, `python-conventions.mdc`,
    `csharp-conventions.mdc`, `java-conventions.mdc` (Epic 31, 2026-07-07 — preferred packages/structure
    per language, sourced from `shared/rules/<language>-conventions.md`). `approval-gates.mdc` and
    `agent-roster.mdc` are the only `alwaysApply: true` files — the rest Auto Attach on a language-specific
    or broad source-file glob instead, since combined they'd otherwise blow well past Cursor's own
    recommended ~2,000-token always-apply budget. Cursor silently ignores a `.mdc` file with invalid
    frontmatter, so `generate_mdc()` is careful about exact YAML shape.
- **Tier 2 (Windsurf, GitHub Copilot)**: generate-inline, multi-file.
  - Windsurf gets one flat `.windsurfrules` (no per-file globs support in the legacy format).
  - Copilot gets the Tier-3-style `copilot-instructions.md` (roster + rules inlined) **plus** the same
    7 scoped `.github/instructions/*.instructions.md` files as Cursor's non-agent-roster `.mdc` set
    (`testing`, `go-backend`, `vue-frontend`, `typescript-conventions`, `python-conventions`,
    `csharp-conventions`, `java-conventions`), each with an `applyTo` frontmatter field (comma-separated
    glob string, not an array like Cursor's `globs`) — both coexist and combine per GitHub's docs.
- **Tier 3 (OpenAI)**: generate-inline, single file. `generate_tier3()` concatenates rules + craftsmanship
  section + persona roster into one instruction file.
- **Gemini/Antigravity**: generates root `AGENTS.md` only (the [agents.md](https://agents.md) cross-tool
  convention — confirmed read, see above). `install.sh` symlinks `shared/skills/` to
  `~/.gemini/config/skills/` on a `--global` install (confirmed global skills root) or to
  `.agents/skills/`/`shared/rules/` to `.agents/rules/` on a `--project` install (documented project-scope
  convention, not yet directly exercised by testing).

---

## 3. Context Flow (Six-Layer Taxonomy)

From `docs/runbooks/context-engineering.md`, the taxonomy every agent operates within, ordered from
permanent/static to ephemeral/dynamic:

| Layer | Name | Source | Lifetime |
|---|---|---|---|
| 1 | System Context | `CLAUDE.md`, `.cursorrules` | Session-long |
| 2 | Rule Context | `ARCHITECTURE_RULES.md`, `shared/rules/` | Session-long |
| 3 | Knowledge Context | `shared/knowledge/` KIs, `docs/adrs/` | Demand-driven |
| 4 | Task/Goal Context | Feature spec, `analysis.md` | Task-long |
| 5 | Historical Context | Thread history, `docs/features/` | Ephemeral |
| 6 | Runtime Context | Open files, tool outputs | Real-time |

`context-engineer` is the agent responsible for keeping Layers 3-6 high-signal before the rest of the
pipeline starts: it maps the task to a Bounded Context, auto-prunes files from unrelated contexts (Epic 17),
searches Layer 3 via `search-ki` before letting `analyst` reason independently (Proactive RAG), and estimates
a token budget per pipeline-agent tier (Analyst/Architect ≤60%, Developer ≤80%, Reviewers ≤40% of a
200k-token window).

### Context decay
An artifact 2+ pipeline phases old gets read as a `summarize-artifact` summary, not the full file — e.g.
`qa-engineer` and `tech-writer` (Phase 3) get `analysis.md`'s (Phase 1) gist, not its full text, since
`implementation-notes.md` already restates what matters for their job. This never applies to the artifact an
agent is *immediately* reviewing.

### Subagent isolation
Spawning a subagent is a clean-slate context — it sees only its own definition plus the specific
artifact/task handed to it, never the orchestrator's full conversation history or other subagents' internal
reasoning. The orchestrator (`deliver-feature`) only ever consumes a subagent's final structured report. This
is why `shared/contracts/` exists: the report *is* the entire interface between agents, so its shape has to
be both complete and predictable.

---

## 4. The Pipeline's Own Observability

The pipeline instruments itself the same way you'd instrument a production system:

- **`pipeline-state.json`** (per feature, in `.claude/feature-workspace/`, persisted to `docs/features/<name>/`)
  — resumability: current phase, completed agents, artifact checksums. `resume-pipeline` reads this.
- **`pipeline-trace.json`** (same location) — timing, status, iteration counts, and `budgetUtilization` per
  agent. `pipeline-trace` (single run) and `pipeline-retrospective` (cross-delivery trends) read this.
- **`docs/agent-metrics/scorecard-YYYY-MM.md`** — monthly quality scores per agent (security TPR proxy,
  code-reviewer first-pass acceptance, analyst completeness, architect fitness-function coverage),
  trend-compared month over month.
- **`docs/lessons-learned/`** — cross-delivery pattern extraction; recurring findings get *drafted* as rule
  or prompt changes, never auto-applied (see `.claude/rules/approval-gates.md` Gate #7 — a rule change
  always requires explicit human sign-off).

---

## 5. Versioning and Testing the Agents Themselves

Agents are prompts, but prompts are code here: every agent has a `version:` field
(`shared/agents/CHANGELOG.md` tracks history), a pre-commit hook (`scripts/hooks/pre-commit`, opt-in) that
requires a version bump + changelog entry for any behavior change, and golden-file structural tests
(`tests/agents/`, run via `scripts/test-agents.sh`) for the five agents most likely to regress silently.
See [docs/runbooks/editing-agent-prompts.md](runbooks/editing-agent-prompts.md) for the full workflow.

---

## Directory Reference

```
shared/                          canonical source — see section 1
scripts/
  generate-configs.sh            shared/ -> six platform configs
  check-parity.sh                fitness function: configs match shared/
  test-agents.sh                 golden-file structural tests for agent prompts
  check-context-budget.sh        fitness function: no WARNING manifest without cut recommendations
  health-check.sh                symlinks, frontmatter, drift, contracts, changelog, KI validity
  check-agent-versions-ci.sh     CI equivalent of hooks/pre-commit (base-branch vs. PR-head, not staged vs. HEAD)
  hooks/pre-commit                opt-in: agent version bump + changelog gate
install.sh / uninstall.sh        --global | --project <path>, --copy, --platform, --dry-run
Makefile                          install, uninstall, generate, check, test-agents, health
.github/workflows/framework-ci.yml  check-parity, test-agents, health-check, agent-versions (PRs only)
tests/agents/                    fixtures + expected patterns for golden-file tests
docs/
  features/<name>/                every delivered feature's full pipeline artifact set
  adrs/                           Architecture Decision Records
  agent-metrics/                  monthly agent quality scorecards
  pipeline-retrospectives/        cross-delivery timing/iteration trend reports
  lessons-learned/                cross-delivery pattern extraction
  runbooks/                       operational guides
```
