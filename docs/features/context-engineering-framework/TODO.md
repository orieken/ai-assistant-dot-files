# Context Engineering Framework — Implementation Checklist

> **Goal**: One repo, one clone, one `./install.sh` — any machine, any AI tool, full framework.
>
> **Platforms**: Claude Code, Cursor, Gemini (Antigravity), GitHub Copilot, Windsurf, OpenAI Codex
>
> See [analysis.md](./analysis.md) for detailed acceptance criteria and [architecture-notes.md](./architecture-notes.md) for structural decisions.

---

## Phase 1: Foundation

### Epic 1 — Canonical shared layer
- [x] Create `shared/rules/` and move `architecture-guardrails.md`, `design-principles.md`, `approval-gates.md` from `.claude/rules/`
- [x] Create `shared/agents/` and move all 20+ agent `.md` files from `.claude/agents/` — 24 agents present
- [x] Create `shared/skills/` and move all 35+ skill directories from `.claude/skills/` — 37 skill dirs present
- [x] Move `ARCHITECTURE_RULES.md` and `DOMAIN_DICTIONARY.md` into `shared/`
- [x] Make `.claude/rules/`, `.claude/agents/`, `.claude/skills/` symlink to `shared/` equivalents
- [x] Create `shared/platform-registry.json` with platform name, config path, format, capability tier
- [x] Verify existing Claude Code functionality is unbroken after restructure — `scripts/check-parity.sh` passes clean across all 6 platforms (24 agents synced everywhere); found and removed two stale broken symlinks left over from before the tier system (`.openai/agents`, `.openai/skills` — pointed at a relative path missing `../`, and weren't part of the current generator or parity check anyway)

### Epic 9 — Persona vs. agent formalization
- [x] Add **Persona** definition to `DOMAIN_DICTIONARY.md` (context frame, no tools, no autonomy) — already present under Entities
- [x] Add **Agent** definition clarification (persona + tools + process + pipeline participation) — already present under Entities
- [x] Add **Capability tier** to domain dictionary (Full / Personas+Rules / System Prompt) — already present under Entities
- [x] Update platform configs to use correct term per platform capability — `collect_agent_roster()` in `scripts/generate-configs.sh` was mixing "Agent / Persona Roster" wording even though it's only ever used for Tier 2/3 output; renamed to "Persona Roster" consistently and added a one-line note pointing to `DOMAIN_DICTIONARY.md` and clarifying full agent orchestration is Tier 1 (Claude Code) only. Regenerated all configs and re-ran `check-parity.sh` — clean

### Epic 3 — Universal install/uninstall
> Audit note: `install.sh`/`uninstall.sh` already existed before this epic was picked up (built in an
> earlier, undocumented pass — same situation as Epic 1) and were substantially further along than this
> checklist suggested. Verified each item against the actual script rather than assuming either "already
> done" or "needs building from scratch."
- [x] Create `install.sh` with `--global` (symlinks to `~/`) and `--project <path>` (copies to target) modes — already existed, confirmed working via `--dry-run` and a real `--project` run into a scratch directory
- [x] Add platform auto-detection (check for `.cursor/`, `gh copilot`, `gemini` CLI, etc.) — already existed (`detect_platforms()`), matches the checklist's exact examples
- [x] Add `--dry-run` flag — already existed, confirmed it suppresses all real filesystem changes
- [x] Make idempotent (backup existing, skip if identical) — symlink mode already did this correctly; **found and fixed a real gap** in `--copy` mode, which always backed up + re-copied even when content was byte-identical. Added a `diff -rq` check; verified a second `--copy` run now reports `[skip] (already copied, content identical)` for everything instead of creating a fresh `.bak.<timestamp>` every time
- [x] Create `uninstall.sh` (remove symlinks, restore backups) — already existed and correctly restores the latest `.bak.<timestamp>` after removing a framework-owned symlink or copy
- [x] Update `scaffold-team.sh` to delegate to `install.sh --project --platform claude` — already did this, using the registry's actual platform name (`claude-code`, not `claude` as this checklist item literally says)
- [x] Print verification summary at end (agent count, skill count, platform count) — agent/skill/rule counts already existed; **added the missing platform count line** (`Platforms: 6 (claude-code cursor windsurf github-copilot gemini openai-codex)`)
- [x] Support macOS, Linux, and WSL — no platform-specific dependencies — audited for bash-4-only syntax (associative arrays, `${var,,}`, `mapfile`) and macOS-only flags (`sed -i ''`, `readlink -f`) — found none; already portable
- [x] Add `--copy` mode fallback for Windows (non-WSL) where symlinks don't work — already existed, now also properly idempotent (see above)

---

## Phase 2: Cross-platform parity

### Epic 2 — Config generation
- [x] Create `scripts/generate-configs.sh` that reads `shared/` + `platform-registry.json`
- [x] Generate `.claude/` config (Tier 1 — symlinks to shared agents/skills/rules)
- [x] Generate `.cursor/rules/` as multiple focused `.mdc` files (Cursor requires all content inline, no file refs):
  - [x] `architecture.mdc` — guardrails + clean architecture rules (alwaysApply: true)
  - [x] `design-principles.mdc` — Kent Beck, Fowler, Sandi Metz rules (alwaysApply: true)
  - [x] `approval-gates.mdc` — human checkpoint rules (alwaysApply: true)
  - [x] `agent-roster.mdc` — full agent/persona roster with descriptions (alwaysApply: true)
  - [x] `testing.mdc` — Saturday/Sunday framework testing rules (globs: `["**/*.spec.*", "**/*.test.*"]`)
  - [x] `go-backend.mdc` — Go-specific conventions (globs: `["**/*.go"]`)
  - [x] `vue-frontend.mdc` — Vue 3 + Tailwind rules (globs: `["**/*.vue", "**/*.tsx"]`)
- [x] Generate `.cursorrules` (Tier 2 — flat concatenation of all rules for legacy Cursor support)
- [x] Generate `.windsurfrules` (Tier 2 — same format as `.cursorrules`)
- [x] Generate `.github/copilot-instructions.md` (Tier 3 — rules + agent awareness)
- [x] Generate `.gemini/antigravity/instructions.md` (Tier 3 — rules + agent awareness)
- [x] Generate `.openai.md` (Tier 3 — rules + agent awareness)
- [x] Extend `scripts/check-parity.sh` to diff generated configs vs. shared rules, fail on drift
- [x] Add CI fitness function: `check-parity.sh` runs on every PR — satisfied as a side effect of Epic 20's `.github/workflows/framework-ci.yml` (`check-parity` job runs on every push/PR to main), not built separately

### Epic 11 — Cross-platform agent/persona translation
- [x] For Cursor: generate `.cursor/rules/<agent-name>.mdc` persona files (all content inlined, short and directive, use ALWAYS/NEVER/CRITICAL keywords, valid YAML frontmatter) — **this was checked off but didn't actually exist**: `.cursor/rules/` only ever had the 7 always-apply rule files, never one `.mdc` per agent. Built `generate_cursor_personas()` in `scripts/generate-configs.sh`: strips each agent's `.claude/rules/*.md` preamble (a file reference Cursor can't follow) and frontmatter, inlines the body verbatim (matching the existing mechanical pattern the other `.mdc` generators already use — no attempt at LLM-style condensing from a deterministic bash script), `alwaysApply: false` since a persona should be invoked deliberately, not always loaded. Found a real bug while building this: 5 agent descriptions contain embedded double quotes (e.g. dependency-auditor's `"audit dependencies"`), which broke the YAML frontmatter until escaped. Validated all 24 generated files parse as valid YAML via `python3 -c "import yaml..."`, and added a parity check to `check-parity.sh` — verified both its pass path (24/24 present) and fail path (removed one file, confirmed DRIFT reported, restored it)
- [x] For Gemini: generate persona blocks in `.gemini/antigravity/instructions.md` — verified genuinely present (not a repeat of the Cursor false-positive), a full "Persona Roster" section already existed
- [x] For Copilot: generate persona reference section in `copilot-instructions.md` — same verification, genuinely present
- [x] Include agent roster summary in all Tier 2/3 configs ("these are the specialists available — invoke by name")
- [ ] Test: verify each platform's AI tool acknowledges the persona/agent roster when prompted — **Cursor confirmed 2026-07-02** (`tests/platform-verification/results/cursor-2026-07-02.md`): always-apply rules, glob-triggered Auto Attach, and manual `@developer.mdc` persona invocation all behave correctly. **Gemini Antigravity confirmed 2026-07-02** (`tests/platform-verification/results/antigravity-2026-07-02.md`, v2.1.1): reads `AGENTS.md` for rules, genuinely invokes skills (not just describes them) from the global skills root — this also caught and fixed a real bug (the previously-generated `.gemini/antigravity/instructions.md` was confirmed unread and removed; `install.sh`'s global skills path was wrong and corrected to `~/.gemini/config/skills/`). One sub-item still open: project-level `.agents/skills/` specifically (vs. the global root actually exercised) — see Test 5 in `antigravity.md`. **GitHub Copilot confirmed 2026-07-02** (`tests/platform-verification/results/copilot-2026-07-02.md`, VS Code Copilot Chat): repo-wide instructions and `.go`/`.vue` path-scoped instructions both verified against their actual file content, not just plausible-sounding output. One sub-item inconclusive: the `.test.ts`/`testing.instructions.md` result showed zero overlap with that file's actual Saturday/Sunday Framework content when checked — likely a false positive from the fixture's own violations being generic enough to flag without any custom instructions loaded. Sharper disambiguation prompt written into `copilot.md`'s Test 3, not yet re-run.

---

## Phase 3: Pipeline hardening

### Epic 4 — Wire context-engineer into pipeline
- [x] Update `deliver-feature/SKILL.md` Phase 0 → add step: invoke context-engineer after setup
- [x] Context-engineer produces `context-manifest.md` in `.claude/feature-workspace/`
- [x] Downstream agents read manifest for pinpointed file list (not full directory scans) — wired into `analyst.md` and `developer.md`
- [x] Add token budget estimation to manifest (flag if > 80% of context window) — tiered thresholds (Analyst/Architect 60%, Developer 80%, Reviewer 40%) per Epic 16

### Epic 5 — Inter-agent contracts
- [x] Create `shared/contracts/analysis-contract.md` (required sections for analyst output)
- [x] Create `shared/contracts/architecture-contract.md` (required sections for architect output)
- [x] Create `shared/contracts/implementation-contract.md` (required sections for developer output)
- [x] Create `shared/contracts/review-contract.md` (required sections for code-reviewer output)
- [x] Create `shared/contracts/security-contract.md` (required sections for security-reviewer output)
- [x] Create `shared/contracts/qa-contract.md` (required sections for qa-engineer output)
- [x] Create `shared/contracts/observability-contract.md` (required sections for sre-engineer output)
- [x] Create `shared/skills/validate-artifact/SKILL.md` (reads contract + artifact, fails if sections missing) — also enforces a few contract-specific content rules (e.g., Overall Status literal match, Failed: 0, PII status) beyond just heading presence
- [x] Wire `validate-artifact` into `deliver-feature` between each agent handoff — added as its own numbered step after each of analyst/architect/developer/code-reviewer/security-reviewer/qa-engineer/sre-engineer; not wired for performance-engineer, data-engineer, accessibility-engineer, tech-writer, devops-engineer since they have no contract yet

### Epic 12 — Pipeline rollback & recovery
- [x] Add checkpoint system to `deliver-feature` — persist pipeline state after each phase — checkpoints marked at every artifact-producing step (finer-grained than per-phase), plus Phase 0 now checks for an existing state file before touching the workspace
- [x] Create `shared/skills/resume-pipeline/SKILL.md` — reads checkpoint state, resumes from last successful phase — also handles `--from-phase N` jump and artifact rollback (Modes 1-3)
- [x] If an agent produces bad output, allow rollback to previous agent's artifact and re-run — via `.claude/feature-workspace/.history/` backups + `resume-pipeline` Mode 3, documented in `deliver-feature`'s new "Rollback" section
- [x] Pipeline state file: `.claude/feature-workspace/pipeline-state.json` (current phase, completed agents, artifact checksums) — schema documented in `deliver-feature`'s new "Checkpointing & Pipeline State" section
- [x] Add `--from-phase N` flag to `deliver-feature` for manual resume — implemented as `resume-pipeline` Mode 2 (natural-language equivalent: "resume delivery on <feature> from phase N")

---

## Phase 4: Quality & observability

### Epic 6 — Agent golden-file tests
- [x] Create `tests/agents/` directory structure
- [x] Create `tests/agents/security-reviewer/` — vulnerable code input + expected findings — smoke-tested end-to-end: ran the real security-reviewer prompt against the fixture, 20/20 checks passed
- [x] Create `tests/agents/code-reviewer/` — smelly code input + expected flags
- [x] Create `tests/agents/analyst/` — feature spec input + expected analysis sections
- [x] Create `tests/agents/architect/` — analysis input with architectural flags + expected structural decisions
- [x] Create `tests/agents/qa-engineer/` — implementation input + expected test plan sections
- [x] Create `scripts/test-agents.sh` — runs golden-file tests with fuzzy matching (grep patterns) — bash 3.2 compatible (macOS default has no associative arrays; used a case-based lookup function instead)
- [x] Add structural checks: verify agent output contains all contract-required sections — reuses `shared/contracts/*.md` from Epic 5 rather than duplicating section lists
- [x] Document: "run `./scripts/test-agents.sh` after editing any agent prompt" — in `tests/agents/README.md`

### Epic 7 — Agent observability & feedback loop
- [x] Create `shared/skills/pipeline-trace/SKILL.md` — logs agent name, duration, status, iteration count — deliver-feature writes the file directly at each Checkpoint (same pattern as pipeline-state.json); this skill owns the schema and answers ad-hoc single-run questions
- [x] Persist `pipeline-trace.json` to `docs/features/<name>/` — wired into deliver-feature Phase 4 persistence and the artifact tree
- [x] Create `shared/skills/pipeline-retrospective/SKILL.md` — analyzes past N traces for patterns — cross-delivery trend analysis, distinct from the existing single-delivery `retrospective` skill; writes to `docs/pipeline-retrospectives/`
- [x] Update analyst agent to read 3 most recent delivery summaries (feedback loop) — also cross-checks `retrospective.md` if present for the same 3 features, since that's where "What To Improve" actually lives
- [x] Update `deliver-feature` to auto-invoke `/retrospective` after every 5th delivery — counts `docs/features/*/delivery-summary.md`, triggers on exact multiples of 5

### Epic 13 — Agent performance metrics & scoring
- [x] Define agent quality metrics — documented in `shared/skills/agent-scorecard/SKILL.md`'s "Metric Definitions" table:
  - Security reviewer: true positive rate (proxy: Critical/High fix-applied rate, adjusted for disputes noted in retrospectives — a real confirmed/false-positive rate needs dispute tracking, see Epic 15)
  - Code reviewer: first-pass acceptance rate (from pipeline-trace.json's changesRequestedCount, added in Epic 7)
  - Analyst: completeness score (contract sections present AND non-placeholder, against analysis-contract.md from Epic 5)
  - Architect: fitness function coverage (decisions with concrete Enforcement vs. judgment-only)
- [x] Create `shared/skills/agent-scorecard/SKILL.md` — reads past N delivery artifacts and scores each agent
- [x] Persist scorecard to `docs/agent-metrics/scorecard-YYYY-MM.md`
- [x] Surface underperforming agents in retrospective output — added an "Agent Scorecard Cross-Reference" section to the single-delivery `retrospective` skill, and `pipeline-retrospective` (Epic 7) already cross-references the latest scorecard for cross-delivery trends
- [x] Track metrics over time: is agent quality improving or degrading after prompt edits? — each scorecard compares against the previous month's file, IMPROVING/STABLE/DEGRADING per metric

### Epic 8 — Agent versioning & changelog
- [x] Add `version: 1.0.0` to every agent's frontmatter — all 24 agents in `shared/agents/`
- [x] Create `shared/agents/CHANGELOG.md` with initial entries — this file has no agent frontmatter, which exposed a real `set -e`/`pipefail` bug in `scripts/check-parity.sh` and `scripts/generate-configs.sh` (the same class of bug as the pre-commit hook's, below): a bare `grep` finding no `name:` line aborted the whole script instead of hitting the existing "skip non-agent files" check further down. Fixed both scripts and confirmed `check-parity.sh` passes clean (still 24 agents, CHANGELOG.md correctly excluded from the roster)
- [x] Create pre-commit hook: agent file change requires version bump + changelog entry — `scripts/hooks/pre-commit`, bash 3.2 compatible, opt-in (not wired to `.git/hooks/` automatically — git doesn't track its own hooks dir; requires `git config core.hooksPath scripts/hooks`, documented in the hook's header). Tested both the FAIL (no bump / no changelog entry) and PASS paths for real; fixed a real `set -e`/`pipefail` bug found while testing (grep finding no version in HEAD for a brand-new agent was aborting the whole script silently)
- [x] Include agent versions in delivery summary output — added a Version column to deliver-feature's Pipeline Run table
- [x] Include agent versions in pipeline-trace.json (correlate version to performance) — added `agentVersion` field to the schema, and wired `pipeline-retrospective` to check whether a duration/iteration trend's boundary lines up with a version change

---

## Phase 5: Knowledge & memory

### Epic 14 — Knowledge Items (KI) infrastructure
- [x] Create `shared/knowledge/` directory for reusable Knowledge Items
- [x] Define KI format: markdown file with frontmatter (tags, domain, created date) — see `shared/knowledge/README.md`
- [x] Create `shared/skills/create-ki/SKILL.md` — captures a pattern, bug fix, or decision as a searchable KI
- [x] Create `shared/skills/search-ki/SKILL.md` — searches KIs by tag/domain before agents start work
- [x] Wire context-engineer to search KIs during manifest creation (Proactive RAG) — refactored to invoke `search-ki` instead of ad-hoc grepping `shared/knowledge/`/`.claude/knowledge/`/`docs/adrs/` directly (context-engineer bumped to v1.1.0, CHANGELOG updated)
- [x] Seed initial KIs from existing ADRs and runbooks — 2 more added: `docs-directory-follows-rag-friendly-structure.md` (from ADR-001) and `subagent-isolation-is-a-hard-boundary.md` (from both runbooks); not an exhaustive migration of every ADR/runbook, but the genuinely reusable nuggets are captured

### Epic 15 — Cross-delivery learning
- [x] Create `shared/skills/extract-lessons/SKILL.md` — after delivery, extracts reusable patterns:
  - Security findings that should become rules
  - Code review rejections that indicate missing guardrails
  - Architecture decisions that should become KIs
- [x] Auto-promote recurring security findings (3+ occurrences) to `shared/rules/` as new guardrails — implemented as **draft + require explicit approval**, not silent auto-write. `.claude/rules/approval-gates.md` Gate #7 ("Wiring a New Fitness Function") already requires human sign-off before any `shared/rules/` change, and that gate isn't something this skill can or should route around just because a pattern is well-evidenced.
- [x] Auto-promote recurring code review patterns to agent prompt improvements — same gating: drafts the prompt edit + version bump + changelog line, requires explicit confirmation before touching `shared/agents/`, consistent with Epic 8's versioning requirement.
- [x] Create `docs/lessons-learned/` directory for persisted lessons
- [x] Track: which KIs are actually referenced by agents (usage analytics) — tallies KI references across all `context-manifest.md` files, flags zero-reference KIs

---

## Phase 6: Context budget & optimization

### Epic 16 — Dynamic context budget management
- [x] Define token budget per agent tier:
  - Analyst/Architect: up to 60% of context window (need broad codebase awareness)
  - Developer: up to 80% (needs code + analysis + architecture)
  - Reviewer agents: up to 40% (focused on specific output)
- [x] Context-engineer estimates token count per file in manifest
- [x] Add `budget_utilization` field to pipeline-trace.json per agent — unblocked now that Epic 7 exists; added as `budgetUtilization` (fraction of the agent's tier ceiling, `null` for agents that don't consume context-manifest.md directly), copied from context-engineer's upfront estimate rather than independently re-measured (no live hook into actual model context usage exists)
- [x] Create `shared/skills/context-audit/SKILL.md` — analyzes a conversation or pipeline run for context waste:
  - Files loaded but never referenced in output
  - Duplicate information across loaded files
  - Large files loaded when a line-range read would suffice
- [x] Create fitness function: "no agent exceeds its token budget tier" — `scripts/check-context-budget.sh`, operationalized as "no persisted context-manifest.md reports WARNING without actionable cut recommendations" (a live token-usage check isn't possible from a shell script; this enforces context-engineer's own guardrail as a real, testable check instead of trusting prose). Tested all 3 cases (OK / WARNING-with-recommendations / WARNING-without) against a fixture tree — correct pass/fail on each

### Epic 17 — Context pruning automation
- [x] Update context-engineer to auto-prune based on bounded context mapping:
  - If task is in `billing` domain, exclude `auth` domain files unless explicitly crossing
  - If task is UI-only, exclude infrastructure/migration files
  (context-engineer bumped 1.1.0 -> 1.2.0, standalone skill twin updated to match, CHANGELOG entry added)
- [x] Add "context decay" — summarize artifacts older than 2 phases in the pipeline instead of passing full text — documented as a general rule in `deliver-feature/SKILL.md`'s new "Context Decay" section, and concretely applied to `qa-engineer` and `tech-writer` (both read `analysis.md` via `summarize-artifact` now instead of in full, since it's 2 phases old by Phase 3) rather than left as prose with no real agent actually doing it
- [x] Create `shared/skills/summarize-artifact/SKILL.md` — produces a 200-word summary of any agent artifact for downstream context compression

---

## Phase 7: Documentation & onboarding

### Epic 18 — Framework documentation
- [x] Update README.md with:
  - Architecture diagram (shared layer → platform configs → project install)
  - Quick start guide (clone → install → verify) — verified `install.sh --global --dry-run` actually runs clean (24 agents, 47 skills, 3 rules) before documenting it; discovered `install.sh`/`uninstall.sh` already exist and work (Epic 3 is further along than the TODO suggested — left as-is per explicit instruction to save Epic 3 for last, just documented what's demonstrably true now)
  - Platform capability matrix (what works where)
  - Agent roster with one-line descriptions (24)
  - Skill catalog with trigger keywords (47, grouped by purpose)
- [x] Create `docs/CONTRIBUTING.md` — how to add a new agent, skill, rule, or platform
- [x] Create `docs/ARCHITECTURE.md` — the canonical `shared/` layer design, tier system, and context flow
- [x] Create `docs/runbooks/adding-a-new-platform.md` — step-by-step guide
- [x] Create `docs/runbooks/editing-agent-prompts.md` — versioning, testing, and changelog requirements — also created `docs/pipeline-retrospectives/README.md` (existed only as a path reference from Epic 7, never actually created) for consistency with `docs/agent-metrics/` and `docs/lessons-learned/`, and linked all runbooks from `docs/runbooks/README.md` (its own stated convention wasn't being followed even for pre-existing runbooks)

### Epic 19 — Onboarding experience
- [x] Create `shared/skills/onboard/SKILL.md` — interactive tour for new users:
  - Explains the three context layers (rules, agents, skills)
  - Shows how to invoke an agent
  - Shows how to trigger a skill
  - Shows how to run a pipeline
  - Lists available approval gates
- [x] Add `install.sh --tour` flag that runs the onboarding skill after setup — the flag already existed as a parsed-but-unused stub (`SHOW_TOUR` was never checked anywhere); wired it to actually print guidance. A shell script can't invoke an AI skill directly, so it prints the exact `/onboard` invocation and what it covers rather than silently doing nothing. Tested both paths: `--tour --dry-run` correctly suppresses it, a real `--project` install with `--tour` correctly shows it
- [x] Create `shared/templates/my-first-feature.md` — a tutorial feature spec that walks through the full pipeline — deliberately touches auth + a new data table + a new env var so `architect`, `data-engineer`, `security-reviewer`, and `devops-engineer` all activate alongside the always-on agents, not just the minimal path

---

## Phase 8: Polish & hardening

### Epic 10 — Health check & self-test
- [x] Implement `health-check` skill to verify (a `health-check` skill already existed but was stale — referenced the old `./install` script and didn't cover version/contract/changelog/KI checks at all; rewrote it to wrap a new `scripts/health-check.sh`):
  - [x] All symlinks resolve
  - [x] All agents have valid frontmatter (name, description, tools, model, version)
  - [x] All skills have valid SKILL.md with triggers
  - [x] All platform configs generated from current `shared/` (no drift) — delegates to `check-parity.sh` rather than duplicating it
  - [x] Domain dictionary has no orphaned terms — best-effort grep-based check across `shared/`+`docs/`; correctly flagged 13 terms with zero in-repo references (mostly Saturday/Sunday framework classes and domain events that are only meaningful in a generated project, not this repo itself)
  - [x] All inter-agent contracts exist for pipeline agents
  - [x] Agent changelog is up to date (no version mismatches)
  - [x] Knowledge Items have valid frontmatter and tags
- [x] `install.sh` runs health-check automatically at end — skipped on `--dry-run`; tested with a real `--project` install
- [x] Add `--verbose` flag for detailed diagnostics
- [x] Add `--fix` flag to auto-repair common issues (regenerate configs, fix symlinks) — found and fixed a real `set -e`/`pipefail` bug while testing (same class as twice before: a `grep` finding zero matches for a domain term killed the script silently)

### Epic 20 — CI/CD integration
- [x] Create `.github/workflows/framework-ci.yml`:
  - Runs `check-parity.sh` (config drift detection)
  - Runs `test-agents.sh` (agent regression tests)
  - Runs `health-check` (structural validation) — `scripts/health-check.sh --verbose`
  - Validates agent version bumps on agent file changes — new `scripts/check-agent-versions-ci.sh` (PR-only job; the existing pre-commit hook only handles staged-vs-HEAD, not base-branch-vs-PR-head, so this is a separate script rather than a forced awkward reuse). Verified against real history (`d5b2bd1..9fc90b0`, Epic 14's version bump) before trusting it, not just written and assumed correct
- [x] Create `Makefile` with targets: `install`, `uninstall`, `generate`, `check`, `test-agents`, `health` — ran `make check`/`generate`/`test-agents`/`health` for real; deliberately did not run `make install`/`uninstall` since those touch the real home directory (install.sh/uninstall.sh already independently verified via `--project` into scratch dirs earlier)
- [x] Add badge to README showing CI status

### Epic 21 — Rollout & migration
- [x] Create `scripts/migrate-v1-to-v2.sh` — moves existing `.claude/agents/` to `shared/agents/`, creates symlinks (also handles skills/, rules/, and the two root reference files) — verified against a synthetic v1-style scratch repo: dry-run, real migration with content preservation, and idempotent re-run all confirmed working, not just written and assumed correct
- [x] Document breaking changes from current structure to `shared/` canonical structure — `docs/MIGRATION.md`
- [x] Tag current state as `v1.0.0` before restructure — applied retroactively to `e7d5557`, the last commit before `cc841a8` began the restructure (the restructure itself predates this TODO-driven session)
- [x] Tag completed framework as `v2.0.0` — applied to HEAD after this epic's commit

---

## Summary: Gap Coverage Matrix

| Gap Identified | Epic(s) | Phase |
|---|---|---|
| No feedback loop / learning system | Epic 7, 14, 15 | 4, 5 |
| No agent evaluation / quality metrics | Epic 6, 13 | 4 |
| Context-engineer not in pipeline | Epic 4 | 3 |
| No inter-agent contract schema | Epic 5 | 3 |
| No agent testing | Epic 6 | 4 |
| No agent observability | Epic 7, 13 | 4 |
| No dynamic context budget | Epic 16, 17 | 6 |
| No RAG retrieval infrastructure | Epic 14 | 5 |
| No cross-platform parity | Epic 2, 11 | 2 |
| No agent versioning / changelog | Epic 8 | 4 |
| No pipeline rollback / recovery | Epic 12 | 3 |
| Persona vs. agent undefined | Epic 9 | 1 |
| No onboarding experience | Epic 18, 19 | 7 |
| No CI for the framework itself | Epic 20 | 8 |
| No migration path from current structure | Epic 21 | 8 |

---

## Decisions (resolved)

1. **Global symlinks vs. copies**: Generated configs only — avoids accidental edits to canonical source.
2. **Cursor file references**: Cursor cannot follow "read this file" instructions in `.mdc` rules. LLMs don't dynamically traverse the file tree. All rules must be **inlined** in the `.mdc` body. Best practices: keep rules short and direct, use ALWAYS/NEVER/CRITICAL keywords, break into small focused `.mdc` files per concern, validate YAML frontmatter or Cursor silently ignores them. This means `generate-configs.sh` must **concatenate and inline** shared rules into each `.mdc` file — no file references.
3. **Agent test strategy**: Structural checks (grep for required sections/findings) for CI. LLM-backed tests for manual verification only.
4. **KI storage**: Both — `shared/knowledge/` for universal patterns (portable across machines), `.claude/knowledge/` for project-specific context.
5. **Auto-promotion threshold**: 3 occurrences across different features triggers auto-promotion of a finding to a rule.
