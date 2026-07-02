# Architecture

This document explains how `shared/` becomes six different platform configs, and how context flows through
the agent pipeline once a feature is being delivered. For "how do I add X," see
[CONTRIBUTING.md](CONTRIBUTING.md). For a narrative walkthrough, see the [README](../README.md).

---

## 1. The `shared/` Canonical Layer

Every agent, skill, and rule is authored exactly once, in `shared/`:

```
shared/
├── agents/          24 agents — .md, YAML frontmatter (name, description, tools, model, version)
├── skills/          47 skills — .md, YAML frontmatter (name, description, triggers)
├── rules/           architecture-guardrails.md, design-principles.md, approval-gates.md
├── contracts/       required-section contracts for pipeline agent handoffs (Epic 5)
├── knowledge/       portable Knowledge Items (KIs) — searchable via search-ki
├── ARCHITECTURE_RULES.md
├── DOMAIN_DICTIONARY.md
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
| **2** | Personas + Rules | Persona-level context shaping + rule files, no native orchestration | Cursor, Windsurf | **Persona** |
| **3** | System Prompt | Single instruction file, everything inlined | GitHub Copilot, Gemini, OpenAI/Codex | **Persona** |

The distinction matters because it's enforced in the generated output, not just documented: `Persona` is a
context frame with no tool access and no autonomous workflow (see `DOMAIN_DICTIONARY.md`'s Entity table) —
Tier 2/3 configs consistently say "persona" in generated roster text (`collect_agent_roster()` in
`scripts/generate-configs.sh`), never "agent." Only Tier 1 gets the word "agent," because only Tier 1
actually runs multi-step orchestration with tool access.

### Generation strategy per tier
- **Tier 1 (Claude Code)**: symlink. `install.sh` creates `.claude/{agents,rules,skills}` -> `shared/`
  equivalents. Always current after a `git pull`; no generation step needed.
- **Tier 2 (Cursor, Windsurf)**: generate-inline. Cursor gets one `.mdc` file per concern
  (`architecture.mdc`, `design-principles.mdc`, `agent-roster.mdc`, `testing.mdc`, etc.) with YAML
  frontmatter (`alwaysApply`, `globs`) and content fully inlined — Cursor silently ignores a `.mdc` file with
  invalid frontmatter, so `generate_mdc()` in `generate-configs.sh` is careful about exact YAML shape.
  Windsurf gets one flat `.windsurfrules` (no per-file globs support).
- **Tier 3 (Copilot, Gemini, OpenAI)**: generate-inline, single file. `generate_tier3()` concatenates rules +
  craftsmanship section + persona roster into one instruction file per platform.

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
  hooks/pre-commit                opt-in: agent version bump + changelog gate
install.sh / uninstall.sh        --global | --project <path>, --copy, --platform, --dry-run
tests/agents/                    fixtures + expected patterns for golden-file tests
docs/
  features/<name>/                every delivered feature's full pipeline artifact set
  adrs/                           Architecture Decision Records
  agent-metrics/                  monthly agent quality scorecards
  pipeline-retrospectives/        cross-delivery timing/iteration trend reports
  lessons-learned/                cross-delivery pattern extraction
  runbooks/                       operational guides
```
