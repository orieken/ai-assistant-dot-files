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
- [x] Test: verify each platform's AI tool acknowledges the persona/agent roster when prompted — **fully confirmed across all 3 platforms, 2026-07-02** (`tests/platform-verification/results/`). **Cursor**: always-apply rules, glob-triggered Auto Attach, and manual `@developer.mdc` persona invocation all behave correctly. **Gemini Antigravity** (v2.1.1): reads `AGENTS.md` for rules, genuinely invokes skills from both the global root (`~/.gemini/config/skills/`) and project-level `.agents/skills/` — caught and fixed two real bugs along the way (`.gemini/antigravity/instructions.md` confirmed unread and removed; `install.sh`'s global skills path corrected). **GitHub Copilot** (VS Code Copilot Chat): repo-wide instructions and all three path-scoped instruction files (`go-backend`, `vue-frontend`, `testing`) confirmed — the `testing.instructions.md` result initially came back inconclusive (the first test question couldn't distinguish "rule loaded" from "generic good advice"), but a sharper disambiguating question definitively confirmed it (named the Saturday/Sunday Framework and specific class names no generic response would invent). All findings verified against actual file content before being accepted, not taken at face value.

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
- [x] Create `shared/contracts/context-manifest-contract.md` (required sections for context-engineer output) — added 2026-07-04 after this exact gap was independently flagged by all 3 of Antigravity/Codex/Copilot's audits
- [x] Create `shared/contracts/{performance,data-engineering,accessibility,docs,devops}-contract.md` and close the last 5 agents' contract gap (performance-engineer, data-engineer, accessibility-engineer, tech-writer, devops-engineer) — added 2026-07-04, same rationale
- [x] Wire `validate-artifact` into `deliver-feature` between every agent handoff, including context-engineer and the 5 above — every pipeline agent now has a contract; nothing left uncovered

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

## Phase 9: Memory engineering and governance polish (post-v2.0.0)

> Work done after the Epic 21 `v2.0.0` tag, driven by three independent external audits (Antigravity,
> Codex, GitHub Copilot) run against `docs/runbooks/self-audit-prompt.md`, plus a scoping decision to bring
> memory-engineering prompts into this same v2 space rather than deferring them to the separate `docs/aos/`
> v3 prototyping effort. Backfilled here because this checklist is supposed to be this feature's single
> source of truth for scope, and it had drifted out of sync with real repo state — the exact class of doc
> drift `docs/runbooks/self-audit-prompt.md` dimension 5 exists to catch, just pointed at itself this time.

### Epic 22 — Memory Engineering (v2 scope, split from AOS/v3)
- [x] Create `shared/memory-registry.json` — catalog of every memory source (KIs, ADRs, feature archive,
      lessons-learned, `DOMAIN_DICTIONARY.md`, `TEAM_TOPOLOGY.md`) with owner/portability/retrieval-backend
      metadata, plus an `optionalPaths` list for conditionally-existing paths like `.claude/knowledge/`
- [x] Write `docs/runbooks/memory-engineering.md` — the Capture → Candidate → Audit → Approve → Index →
      Retrieve → Expire lifecycle, plus the Memory Contract (Candidate/Audit required fields)
- [x] Create `shared/skills/memory-engineer/SKILL.md` — periodic KI-corpus duplicate/expiration sweep
- [x] Create `shared/skills/promote-memory/SKILL.md` — evaluates one retrospective immediately (distinct
      from `extract-lessons`' cross-delivery pattern threshold)
- [x] Create `shared/skills/query-memory/SKILL.md` — registry-aware search, delegates KI/ADR lookups to
      `search-ki` instead of duplicating that logic
- [x] Update `search-ki` and `context-engineer` (agent + skill twin, bumped 2.0.0 -> 2.1.0) with memory
      registry awareness
- [x] Extend `scripts/health-check.sh` with a Memory Registry check (valid JSON, every path exists or is
      marked optional, no duplicate KI frontmatter names)
- [x] Write `docs/runbooks/lightrag-integration.md` as a deliberately deferred "when and how" runbook — no
      LightRAG code or dependency added; YAGNI until the KI corpus actually outgrows lexical search
- [x] Cross-reference from `context-engineering.md`, `DOMAIN_DICTIONARY.md`, README.md, `runbooks/README.md`

### Epic 23 — Contract completeness closure (Epic 5 follow-through)
- [x] Add `shared/contracts/context-manifest-contract.md` — closed a gap independently flagged by all 3
      external audits (context-engineer's own output had no contract)
- [x] Add the last 5 pipeline agents' contracts: `performance-contract.md`, `data-engineering-contract.md`,
      `accessibility-contract.md`, `docs-contract.md`, `devops-contract.md` — read each producing agent's
      actual current Output Format headings fresh rather than assuming
- [x] Wire all 5 into `validate-artifact`'s Contract Mapping table and `deliver-feature` (renumbered the
      pipeline from 34 to 39 steps, fixed every cross-reference including two stale trace-JSON examples)
- [x] Rewrite `scripts/health-check.sh`'s "Inter-Agent Contracts" check to parse `validate-artifact`'s own
      table directly instead of maintaining a second, independently-drifting hardcoded agent/contract list

### Epic 24 — Agent Reference documentation
- [x] Create `docs/AGENT_REFERENCE.md` — Role, Counterbalance (Structural contract / Downstream agent
      review / Human approval gate / Aggregate-delayed metric), and an explicit Gap for all 24 agents,
      built by reading every `shared/agents/*.md` file fresh and cross-referencing `deliver-feature`,
      `validate-artifact`, and `shared/rules/approval-gates.md`
- [x] Link it from README.md's Agent Roster section and "For deeper detail" list

### Epic 25 — Team Topologies real ownership
- [x] Fill in `shared/TEAM_TOPOLOGY.md`'s placeholder team names — all 5 bounded-context rows set to the
      actual (solo) owner, with an explicit "current state: solo, pre-multi-team" framing and a stated
      trigger for when Team Type/Interaction Mode discipline starts mattering for real

### Epic 26 — documentation-manager narrowed to ad-hoc-session counterpart of promote-memory
- [x] Redesign `documentation-manager` (1.0.0 -> 2.0.0, **major**): previously wrote directly to
      `ARCHITECTURE.md`/`RUNBOOKS.md`/`GOTCHAS.md`/`ONBOARDING.md` with no review step — an undocumented
      overlap with the Epic 22 memory-engineering skills flagged in `docs/AGENT_REFERENCE.md`. Now produces
      Candidate Records via the same Memory Contract as `promote-memory`, requires explicit human approval
      before any KI/ADR/rule/living-doc edit, and retires `GOTCHAS.md` as a target (gotchas are Knowledge
      Items now, via `create-ki`)
- [x] Update `docs/runbooks/memory-engineering.md`, `shared/DOMAIN_DICTIONARY.md`'s Candidate Record entry,
      and `docs/AGENT_REFERENCE.md`'s entry to reflect `documentation-manager` as a third Candidate producer

### Epic 27 — CI-parity verification tooling
- [x] Create `scripts/ci-check.sh` — runs `check-parity.sh`, `test-agents.sh`, and `health-check.sh
      --verbose` inside a `docker run ubuntu:24.04` container matching the actual GitHub Actions runner,
      instead of trusting a local run that might be on a different bash version. Built after a real
      CI-breaking bug (`((var++))` evaluating to a "false" exit status under `set -e` when the pre-increment
      value is 0 — invisible on macOS's bundled bash 3.2, real on Ubuntu's bash 5.2.21) shipped past a local
      check. BATS was considered and explicitly rejected: CI was already correctly catching the bug on every
      push, so the real gap was not checking CI status, not a testing-coverage gap
- [x] Document it in `docs/CONTRIBUTING.md`'s new "Before you push" section

### Epic 28 — Framework diagrams
- [x] Add a "The Framework at a Glance" Mermaid diagram to README.md (shared/ → 6 tools → deliver-feature
      pipeline → Context/Memory/Learning loop feeding back into shared/)
- [x] Add a Context/Memory/Learning cycle Mermaid diagram to `docs/runbooks/context-engineering.md`
- [x] Validate both via a real render pass (`@mermaid-js/mermaid-cli` through a headless browser), not just
      syntax linting

### Epic 29 — External audit fixes (2026-07-05)
- [x] Remove hardcoded personal machine paths from `scripts/api-generator/index.ts` — CLI args or
      `API_GENERATOR_GO_DIR`/`API_GENERATOR_TS_DIR` env vars instead
- [x] Mark `scripts/api-generator` explicitly experimental/unsupported (new README.md + package.json
      description) rather than leaving its lack of tests/CI integration unstated
- [x] Make `context-engineer`'s Prior-Deliveries lookup case-insensitive (2.1.0 -> 2.1.1) — the archived
      `context-engineering-framework/analysis.md` uses `**Owning context**` (lowercase c), which the old
      exact-match grep silently missed. Left the archived doc itself untouched — the feature archive is
      treated as an immutable historical record elsewhere in this framework
- [x] Delete `docs/files.zip`, `docs/spec-writer.zip`, `docs/dotfiles-additions.zip` — unreferenced anywhere
      in the repo, contents fully superseded by canonical `shared/agents/` files. Left
      `docs/clean-code-guidelines.docx` alone since it's deliberately referenced from `CLAUDE.md`
- [x] Add a "Known Judgment Calls" section to `docs/runbooks/self-audit-prompt.md` recording two findings
      (uneven standalone-agent governance depth, some platform behaviors marked unconfirmed) that were
      reviewed and confirmed as already-deliberate tradeoffs, so a future audit doesn't re-flag them as new

### Epic 30 — Cursor native skills/agents parity (2026-07-06/07)
> Cursor shipped native Agent Skills (`.cursor/skills/*/SKILL.md`) and subagent (`.cursor/agents/*.md`)
> support, using the same open standard `shared/agents/`/`shared/skills/` already follow (confirmed
> against `cursor.com/docs/subagents`, `cursor.com/docs/skills`). This closed most of the gap between
> Cursor's Tier 2 classification and Claude Code's Tier 1 — see the updated Tier system section in
> `docs/ARCHITECTURE.md` and `shared/platform-registry.json`'s Cursor entry.
- [x] **Phase 1 — prerequisite frontmatter fixes (all 24 agents)**: `model: sonnet` -> `model: inherit`
      (both Claude Code and Cursor default to `inherit` when the field is omitted, and both accept the
      literal keyword) and relocated a "read these rule files" preamble that sat *before* the opening
      `---` on 23 of 24 agents into the body, using canonical `shared/rules/` paths instead of the
      Claude-Code-only `.claude/rules/` prefix. Invisible to Claude Code's tolerant loader and
      `health-check.sh`'s lenient grep-anywhere frontmatter check, but would have broken Cursor's
      stricter parser once agents were symlinked directly in Phase 3. All 24 patch-bumped, one
      CHANGELOG entry.
- [x] **Phase 2 — symlink feasibility check**: discovered `.cursor/agents`/`.cursor/skills` already
      existed in this repo as symlinks, committed 2026-04-09 (`d0b54d3`, "expanded to work for all
      platforms") — an earlier, forgotten attempt at this exact idea that predates the whole Tier
      system and was never wired into `check-parity.sh`/`platform-registry.json`/docs. User confirmed
      live in Cursor that the `analyst` subagent and `search-ki` skill both load and behave correctly
      (not just generically) via this mechanism.
- [x] **Phase 3 — retire the old workaround, wire up the new one**: removed
      `generate_cursor_personas()` and the now-dead `extract_agent_body()` helper from
      `scripts/generate-configs.sh` (Cursor's generate step now only produces the 7 always-apply/
      glob-triggered rule `.mdc` files — rules still require inlining, agents/skills don't). Added
      `install_cursor()` to `install.sh` (symlinks `.cursor/agents`/`.cursor/skills` -> `shared/`,
      mirroring `install_claude_code()`) and matching removal logic to `uninstall.sh` — both verified
      with real installs to a scratch directory, not just `--dry-run`. Replaced
      `check-parity.sh`'s 24-persona-file existence check with a symlink-resolution check. Fixed this
      repo's own pre-existing symlinks to point directly at `../shared/{agents,skills}` instead of
      double-hopping through `../.claude/`. Deleted the 24 now-orphaned persona `.mdc` files.
- [x] **Phase 4 — registry and doc sync**: updated `shared/platform-registry.json`'s Cursor
      capabilities (`agents`/`skills`/`subAgentOrchestration` -> `true`, `hooks` stays `false`),
      `docs/ARCHITECTURE.md`'s Tier system section and generation-strategy bullet,
      `README.md`'s Platform Capability Matrix, `shared/DOMAIN_DICTIONARY.md`'s Capability Tier entry
      (now documents Cursor as a mixed profile, not a clean single tier), and
      `docs/CONTRIBUTING.md`'s agent-creation frontmatter example (`model: inherit`, explicit
      "must start with `---` on line 1" requirement). Documented the one real permanent gap: Cursor
      subagents have no `tools:` allowlist — they inherit all parent tools with only a coarse
      `readonly: true/false`.
- [x] Verified throughout: `check-parity.sh`, `health-check.sh --verbose` (159 passed/0 failed),
      `test-agents.sh`, real `install.sh`/`uninstall.sh` cycles to scratch directories — all green at
      every phase, committed and pushed separately per phase.

### Epic 31 — Language/framework convention files (2026-07-06/07)
> Scoped while discussing whether the framework should document preferred packages and project
> structure per language: `go_backend_rules_body()`, `vue_frontend_rules_body()`, and
> `testing_rules_body()` in `scripts/generate-configs.sh` already gave architectural *patterns* (Clean
> Architecture layers, Composition API, Site-Centric pattern) but zero concrete *package* or
> *directory-structure* guidance -- and, more importantly, none of that content was actually sourced
> from `shared/rules/` (which only had 3 files). It was hardcoded as literal bash strings inside the
> generator script, reaching Cursor and Copilot but **not Claude Code** (no `shared/rules/go-
> conventions.md` for `.claude/rules/` to symlink) -- a single-source-of-truth gap for this one content
> type specifically. Scope expanded from the original Go/Vue/testing-only plan to 5 languages (Go,
> TypeScript, Python, C#, Java) with an explicit testing-tooling checklist per language: fake-data
> (faker-equivalent), factories/fixtures (fishery-equivalent), Playwright bindings, k6, unit test
> framework, and reporting tools.
- [x] Fixed the sourcing gap: `testing_rules_body()` and `go_backend_rules_body()` now read from real
      `shared/rules/testing-conventions.md` / `shared/rules/go-conventions.md` files via
      `extract_rule_content` instead of embedding literal strings -- content unchanged for
      `testing_rules_body()`, expanded with a Testing & QA Tooling section for `go-conventions.md`
- [x] Created 4 new files: `shared/rules/typescript-conventions.md`, `python-conventions.md`,
      `csharp-conventions.md`, `java-conventions.md`. Python and C# grounded directly against
      `saturday-monorepo-python`'s and `saturday-monorepo-csharp`'s own READMEs (uv/ruff/pytest/
      pytest-bdd/Faker/polyfactory for Python; .NET 8/NuGet CPM/Reqnroll+xUnit/Bogus/AutoFixture/
      `Saturday.K6Exporter`/`Saturday.Reporting` for C#) -- not assumed. Java has no internal reference
      repo yet, so it's explicitly labeled as industry-standard picks (DataFaker, Instancio, JUnit5+
      Mockito, Allure) rather than a confirmed internal decision. TypeScript fills the one gap the
      existing Saturday/Sunday docs hadn't named yet: `@faker-js/faker` + `fishery` specifically.
- [x] Wired 4 new `generate_mdc()` calls into `generate_cursor()` (glob-scoped per language: `**/*.ts`,
      `**/*.py`, `**/*.cs`, `**/*.java` -- `.ts` deliberately excludes `.tsx`/`.vue`, which stay under
      the existing `vue-frontend.mdc`) and 4 new `generate_instructions_md()` calls into
      `generate_copilot_scoped_instructions()`. Cursor's rule-file count: 7 -> 11. Copilot's scoped
      instructions: 3 -> 7.
- [x] Claude Code confirmed receiving all of this content now via the existing `.claude/rules` symlink
      (verified: `ls .claude/rules/` shows all 9 files, zero extra wiring needed -- the whole point of
      fixing the sourcing gap).
- [x] Updated `check-parity.sh` (4 new `.mdc` existence + content checks, 4 new scoped-instructions
      checks), `docs/ARCHITECTURE.md`'s generation-strategy section, and `README.md`'s Platform
      Capability Matrix (Cursor's file count, Copilot's scoped-instructions list) to match.
- [x] Verified: `check-parity.sh`, `health-check.sh --verbose` (159 passed/0 failed), `test-agents.sh`,
      and a direct YAML-frontmatter parse check on all 8 new generated files (4 `.mdc` + 4
      `.instructions.md`) -- all valid.

### Epic 32 — Install verification matrix (2026-07-08)
> Prompted by an external audit (`docs/audits/perplex-audit.md`, Perplexity) flagging that the repo's
> multi-platform support claims were stronger on metadata/static checks than on real
> installation/runtime verification. Several of that audit's other suggestions turned out to already
> exist under different names (`health-check.sh --fix` is its "doctor command"; `agent-scorecard` +
> `agent-eval` + `tests/agents/` collectively are its "prompt quality benchmark suite";
> `docs/AGENT_REFERENCE.md` is most of its "agent coverage report") -- this was the one genuinely new,
> well-motivated gap: `check-parity.sh` only validates this repo's own already-generated output, never
> a fresh install into a real target project.
- [x] Created `scripts/test-install.sh` -- runs a real `install.sh --project <scratch-dir> --platform
      <name>` for all 6 platforms and asserts the expected symlinks/files exist and resolve correctly
      (24-agent/53-skill/9-rule symlinks for Claude Code/Cursor/Gemini, 11 `.mdc` files for Cursor, 7
      `.instructions.md` files for Copilot, flat files for Windsurf/OpenAI). Verified both the pass path
      (32/32 checks green on a real install) and the fail path (deliberately broke a symlink and a file,
      confirmed both were caught) before trusting it.
- [x] Wired into `scripts/ci-check.sh` as a 4th check, running inside the same Docker container as the
      other three -- writes to the container's own `/tmp`, not the read-only `/repo` mount, so it's safe
      there without changing the mount's permissions.
- [x] Initially held back from `.github/workflows/framework-ci.yml` since that's a real CI/CD pipeline
      change, gated by `shared/rules/approval-gates.md` #7 ("Wiring a New Fitness Function") -- flagged
      to the user as a separate, explicit decision rather than bundled into the original commit.
      Approved 2026-07-08: added as a 5th `test-install` job in `framework-ci.yml`, running alongside
      `check-parity`/`test-agents`/`health-check` on every push and PR.
- [x] Updated `docs/CONTRIBUTING.md`'s "Before you push" section to describe the new check.

### Epic 33 — docs/ cleanup (2026-07-08)
> A user-directed cleanup pass, not audit-driven. `docs/dotfiles-additions/` and `docs/spec-writer/`
> turned out to be extracted, still-tracked copies of the exact zips already deleted in Epic 29
> (`dotfiles-additions.zip`, `spec-writer.zip`) -- Epic 29 only deleted the archives, missing that the
> unzipped directories were sitting right next to them, equally stale. Investigated and removed
> alongside them: `build-out-prompts.md`/`master-build-out-prompts.md`/`dotfiles-remediation.md`
> (2,531 lines of March-era "hand this prompt to Claude Code" bootstrap scratch, fully superseded by
> `CONTRIBUTING.md`'s actual current process), `thoughtworks-specialist.md` (mostly redundant with
> `shared/rules/design-principles.md`), `docs/mcp/framework-tools-prompts.md` + `rag-example.md`
> (build-out prompts for an MCP server at `mcp/` that was never created in this repo -- confirmed via
> `git log --all`, no history for that path -- superseded by the separate `aakg-mcp` integration
> instead), and `docs/patterns/README.md` (a stub describing a "Reusable Patterns" directory that was
> never actually populated -- `docs/README.md` falsely claimed agents read it).
- [x] Deleted all of the above (verified each was unreferenced by any current operational doc before
      removing -- only the already-stale `docs/README.md` pointed at any of it)
- [x] Rewrote `docs/README.md` to reflect the actual current `docs/` structure (it hadn't been updated
      since a very early snapshot -- missing `AGENT_REFERENCE.md`, `MIGRATION.md`, `agent-metrics/`,
      `lessons-learned/`, `pipeline-retrospectives/`, `runbooks/`, `blog-posts/`, `audits/`, and
      pointing at everything just deleted)
- [x] Populated `docs/patterns/` for real (2026-07-08): `gang-of-four-patterns.md` (adds
      structure/example/trade-off detail to `CLAUDE.md`'s existing one-line decision table, doesn't
      duplicate it), `saturday-framework-patterns.md` (`BaseSite`/`BasePage`/`BaseElement`/`BaseFlow`/
      `Filters`/`SiteManager`/`TabManager`), `sunday-framework-patterns.md` (`BaseApiClient`/
      `IHttpAdapter`/Fluent Matchers/Schema Validation/Resilience Primitives), and
      `clean-architecture-layers.md` (expands `architecture-guardrails.md`'s dependency-direction rule
      per layer). Grouped by category rather than one file per pattern -- most individual patterns are a
      few paragraphs, splitting each into its own file would've added navigation cost with no real
      benefit. `docs/README.md`'s directory listing and "How Agents Use These Docs" section updated to
      reference `patterns/` again now that it has real content.

### Epic 34 — Expanded pattern catalog: DDD, EIP, Stability Patterns, 12-Factor (2026-07-08)
> Follow-on to the Epic 33 patterns work: does the catalog include Domain-Driven Design? Should it just
> list all 23 GoF patterns? Are there newer concepts (12-Factor) worth adding? User's call: add them all.
> Prioritized additions that are already-practiced-but-undocumented (same gap class Epic 33 fixed)
> over speculative new theory -- `analyst.md`/`architect.md` already explicitly channel Eric Evans,
> Gregor Hohpe, and Michael Nygard by name; none of those three had a pattern doc until now.
- [x] `domain-driven-design.md` -- Entity, Value Object, Aggregate/Aggregate Root, Repository, Domain
      Service, Domain Event, Bounded Context, Context Map, Anti-Corruption Layer, Ubiquitous Language.
      Grounded directly in existing repo mechanics rather than generic DDD theory: Domain Events
      cross-reference `DOMAIN_DICTIONARY.md`'s own event table, Bounded Context references the same
      file's 5 core domains, Anemic Domain Model warning cross-references
      `design-principles.md`'s Anti-Pattern Radar.
- [x] `enterprise-integration-patterns.md` -- Message Channel, Content-Based Router, Message Translator,
      Publish-Subscribe Channel, Dead Letter Channel, Correlation Identifier, Saga. This repo doesn't run
      literal message-queue infrastructure, so each pattern is grounded in a structurally-equivalent
      mechanism that already exists: `pipeline-trace.json` as Correlation Identifier,
      `deliver-feature`'s checkpointed pipeline + `resume-pipeline` rollback as Saga, `.history/` backups
      as Dead Letter Channel.
- [x] `stability-patterns.md` -- Circuit Breaker (general form; Sunday's API-specific version already
      existed), Bulkhead, Timeout, Fail Fast, Steady State. Timeout and Circuit Breaker both grounded in
      already-hard guardrails (`architecture-guardrails.md` #5), not just Nygard's book.
- [x] `twelve-factor-app.md` -- all 12 factors. The one axis nothing else in `docs/patterns/` covers:
      production runtime characteristics, not code structure. Cross-referenced against existing partial
      coverage per factor (Config -> no-hardcoded-secrets guardrail, Dependencies -> each language's
      package-manager convention, Admin Processes -> `db-migration`'s expand/contract skill) rather than
      presented as unrelated new theory.
- [x] Added Template Method, Composite, and State to the existing `gang-of-four-patterns.md` instead of
      blanket-adding all 23 canonical GoF patterns -- explicitly labeled "not in CLAUDE.md's table," named
      only because this codebase's own mechanisms (`BasePage`/`BaseFlow`, `BaseElement` composition,
      `CircuitBreaker`/`pipeline-state.json` phase tracking) already use them unnamed. The other ~13
      unused GoF patterns (Singleton, Flyweight, Interpreter, Memento, etc.) were deliberately not added
      -- no evidence anything in this stack reaches for them.
- [x] Found and fixed one citation error while verifying: `twelve-factor-app.md`'s Dev/Prod Parity
      section originally cited a `docker-compose.e2e.yml` detail that lives in the real
      `saturday-monorepo-csharp` README, not in this repo's own condensed `csharp-conventions.md` --
      corrected to cite `scripts/ci-check.sh` instead, which is a real, already-verified fact from this
      session (the bash-version CI bug it exists to catch).
- [x] Verified: `health-check.sh --verbose` (167 passed/0 failed, DOMAIN_DICTIONARY.md orphaned-terms
      warnings dropped from 15 to 7 since several previously-orphaned terms are now referenced from the
      new pattern docs), `check-parity.sh`, and a cross-reference resolution check confirming every
      backtick-quoted filename across all `docs/patterns/*.md` resolves to a real file in the repo.

### Epic 35 — Second pattern catalog round: Security, Observability, Expand/Contract, API Design (2026-07-08/09)
> Same "already-practiced-but-undocumented" prioritization as Epic 34, applied to a second round of
> candidates. Two more agents turned out to have fully-worked frameworks locked inside their own prompts:
> `security-reviewer.md`'s complete STRIDE table with Saturday/Sunday-specific examples, and
> `sre-engineer.md`'s four governing observability principles.
- [x] `security-patterns.md` -- STRIDE (full 6-category table with concrete example gaps, lifted
      directly from `security-reviewer.md`'s own contract), Secure by Default, Defense in Depth, Least
      Privilege, Paved Road / Golden Path. Least Privilege cross-referenced against Cursor's `readonly`
      subagent flag and `shared/agents/*.md`'s own `tools:` allowlist as real, already-shipped instances.
- [x] `observability-patterns.md` -- SLI (already-required, lifted from `sre-engineer.md`'s own
      contract), plus SLO and Error Budget -- the one layer above SLI this repo uses the foundation of
      without ever naming. Low-Cardinality Logging, Structured Tracing, and No PII in Telemetry lifted
      directly from `sre-engineer.md`'s governing principles, including its own bad/good logging example.
- [x] `expand-contract-migrations.md` -- upgraded from "a rule that's referenced" to a full pattern doc,
      same treatment Circuit Breaker got in Epic 34's `stability-patterns.md`. Uses `db-migration/
      SKILL.md`'s own worked example (`lockout_expires_at` -> `locked_until`) verbatim, and documents its
      own named approval gate (`approval-gates.md` #4, the Contract-phase-specific one) as part of the
      pattern, not just the migration mechanics.
- [x] `api-design-patterns.md` -- Resource-Oriented Design, Idempotency Keys, Status Code Discipline,
      Pagination by Default, User Enumeration Prevention, Schema-First Contract -- every one of these is
      already an enforced guardrail inside `openapi/SKILL.md`, just never extracted as a standalone
      reference. Cross-references `security-patterns.md` (User Enumeration Prevention is literally
      STRIDE's Information Disclosure category applied to API design) and `sunday-framework-patterns.md`
      (the skill's own "Sunday Framework Mapping" output section connects API design directly to
      `BaseApiClient`/Zod schema conventions).
- [x] Verified: `health-check.sh --verbose`, `check-parity.sh`, cross-reference resolution across all 9
      `docs/patterns/*.md` files, and a direct quote-check against `approval-gates.md` #4's exact gate
      text before citing it.

---

### Epic 36 — New agent: unit-tester, gets memory awareness alongside test-driven-developer (2026-07-13)
> Surfaced by a direct question: `test-driven-developer` writes tests and changes code to satisfy them, but
> there was no standalone agent for the opposite case -- adding tests to code that must *not* change,
> whether to raise coverage on trusted code or to build a Michael Feathers-style characterization-test
> safety net before a legacy refactor/migration. `qa-engineer` already had the right guidance for this (its
> own "Testing Legacy Code" section) but only reachable through `deliver-feature`'s workspace artifacts.
> While closing that gap, also closed a smaller one found first: `test-driven-developer` itself started
> every run cold, with no KI lookup and no path for its learnings to reach the memory system afterward.
- [x] `test-driven-developer` 1.0.1 -> 1.1.0 (Minor): new process step invokes `search-ki` before test
      design (read-only, non-blocking). New "Knowledge Consulted" output section. New rule recommending
      `documentation-manager` after a substantial session, without auto-invoking it.
- [x] New agent `unit-tester` (1.0.0): standalone, writes/backfills unit tests for existing code without
      modifying it. Stricter than `qa-engineer` -- never modifies source, not even to fix a discovered bug;
      a required seam is treated as `approval-gates.md` gate #6 (Writing Files out of Boundary) and held
      for explicit approval rather than performed automatically. Same `search-ki` + non-auto-invoked
      `documentation-manager` pattern as the `test-driven-developer` update above.
- [x] New skill `backfill-unit-tests`: coordinates `unit-tester` then automatically runs `code-reviewer`
      against just the new test files (never the untouched source) -- unlike the soft
      `documentation-manager` recommendation, this counterbalance is auto-chained rather than merely
      suggested, since a code-quality pass on new tests is cheap and useful every time, with none of
      `documentation-manager`'s "most sessions produce nothing worth promoting" waste risk. Mirrors
      `review-pr`, this repo's existing precedent for a thin orchestration skill.
- [x] `docs/AGENT_REFERENCE.md`: updated agent-count references (24 -> 25), added entry #25 for
      `unit-tester`, updated the survey summary's standalone-agent count (9 -> 10).
- [x] `README.md` / `docs/ARCHITECTURE.md`: agent/skill counts updated (24 -> 25 agents, 53 -> 54 skills).
- [x] Verified: `generate-configs.sh`, `check-parity.sh`, `health-check.sh --verbose`, `test-agents.sh` all
      green.

---

### Epic 37 — deliver-atdd: config-driven, trust-progression ATDD workflow (2026-07-15)
> Surfaced by a direct question: what if a team wants an ATDD/BDD-shaped delivery (pair on spec ->
> qa-engineer writes Gherkin -> human review -> qa-engineer writes step defs -> human review ->
> test-driven-developer implements to green -> qa-engineer runs -> ship review) with the ability to
> phase out mechanical human-review gates as trust is earned, rather than either running
> `deliver-feature`'s full 14-agent pipeline or manually orchestrating three agents by hand every time?
> The two existing orchestration skills (`review-pr`, `backfill-unit-tests`) are stateless -- every
> invocation runs the same fixed chain. This one's whole point is that its gate configuration is a
> repo-persisted property that changes over time, so a static skill can't capture what makes it
> interesting.
- [x] New skill `deliver-atdd`: coordinates `qa-engineer` (scenario writing, step definitions, and
      final acceptance run) with `test-driven-developer` (autonomous inner red-green loop, per its
      already-established v1.1.0 contract -- no gate around it by design).
- [x] Config-driven trust progression via `.claude/atdd-config.json` (project-root, checked into git so
      the trust curve is auditable in history, not tribal knowledge). Two configurable gates:
      `scenario-review` and `test-code-review`, each `active` (pause for human) or `phased-out`
      (skip). Two other gates are non-configurable and always run: the initial `spec-writer`/`analyst`
      pairing (pre-pipeline, matching `deliver-feature`'s own pre-pipeline stance) and the final
      ship-readiness gate (matching `approval-gates.md`'s principle that irreversible actions never
      delegate the final "yes").
- [x] "Suggest, don't act" pattern for gate progression: at end of a successful run, if a gate has
      been active for 5+ consecutive runs with no human-requested edits, surface a suggestion to the
      user to consider phasing it out. Never mutates the config file itself -- matches
      `memory-engineer`'s existing "surface, don't act" pattern for KI expiration.
- [x] Reuses `deliver-feature`'s proven infrastructure verbatim -- `feature-workspace/`,
      `pipeline-state.json` (with new `"pipeline": "deliver-atdd"` field so resumer can distinguish),
      `pipeline-trace.json`, `docs/features/<name>/` archive shape. Divergences from `deliver-feature`
      called out explicitly in the skill's own "When To Use": no architect/perf/data/a11y/security
      reviewers (use `deliver-feature` if any of those are needed -- this is a scope narrowing, not a
      superset), no Friday POST at the end (not blocked, just not built -- flagged as a worthwhile
      future extension for teams that want it).
- [x] `README.md` / `docs/ARCHITECTURE.md`: skill count updated (54 -> 55).
- [x] Verified: `generate-configs.sh`, `check-parity.sh`, `health-check.sh --verbose`, `test-agents.sh`
      all green.

---

### Epic 38 — Testing taxonomy explicit + annotation convention + honest TDD scoping (2026-07-15)
> Surfaced by two related questions: does the framework distinguish clearly between unit / integration
> / acceptance / E2E tests, and should tests annotate their originating issue and specific AC for
> traceability? Both real gaps -- the framework had structural separation (Saturday, Sunday,
> test-driven-developer, unit-tester) but no doc naming the pyramid or its principles per level, and
> traceability was only at qa-report.md time (evaporates when files get renamed). A third gap surfaced
> during the conversation itself: single-agent TDD doesn't produce XP TDD's design benefit the way pair
> programming does -- role separation is what makes the discipline work. Rather than cargo-cult TDD
> across the framework, this update states the honest scope.
- [x] New pattern doc `docs/patterns/testing-pyramid.md`: five test levels (Unit, Integration, API
      Contract, Acceptance, E2E/UI), each with context/principles/writing-agent, and cross-references
      to `saturday-framework-patterns.md` / `sunday-framework-patterns.md` for the top two levels
      rather than restating them. FIRST principles (Uncle Bob) as unit-test properties. Three Laws of
      TDD stated once, canonically, with an honest "When the discipline actually applies to agent-
      written code" section: role-separated (via `deliver-atdd`) preserves XP TDD's design pressure;
      standalone `test-driven-developer` doesn't, and that's fine -- design pressure for solo agents
      comes from complexity thresholds, SOLID, `code-reviewer`, and refactoring passes, not from the
      test-first ritual alone. Explicitly notes `unit-tester` never follows the Three Laws (impossible
      when the code came first -- its discipline is Feathers' characterization).
- [x] Two new sections in `shared/rules/testing-conventions.md`: Test Categories (a level -> agent ->
      framework -> speed-budget -> principles table, enforcement side of the pattern doc's philosophy)
      and Test Annotation Convention (issue-ref + AC-ref per test, with per-language mechanisms --
      JSDoc for TS, docstring for pytest, `@Tag`/`@DisplayName` for JUnit, `[Trait]` for xUnit, comment
      for Go, `@issue:...` tag for Gherkin). Documented-only, no CI fitness function -- matches how
      Sandi Metz's rules are handled today.
- [x] `test-driven-developer` 1.1.0 -> 1.2.0: cite Three Laws + FIRST explicitly in the preamble;
      state the standalone-vs-role-separated scoping honestly; add annotation-convention step; read
      `testing-conventions.md`.
- [x] `unit-tester` 1.1.0 -> 1.2.0: state explicitly that it does NOT follow the Three Laws (its
      discipline is characterization, not TDD); its tests satisfy FIRST as properties; add
      annotation-convention step with characterization-mode guidance; read `testing-conventions.md`.
- [x] `qa-engineer` 1.1.1 -> 1.2.0: multi-level scope stated explicitly (this agent legitimately owns
      integration / API contract / acceptance / E2E); annotation step with Gherkin-specific guidance;
      read `testing-conventions.md`.
- [x] `deliver-atdd` gets a "Why this workflow specifically" paragraph in its When To Use section
      pointing at the testing-pyramid doc's role-separation framing -- this is the framework's
      strongest agent-TDD shape and worth calling out as such.
- [x] `docs/patterns/README.md`: new entry for testing-pyramid.md at the end of the pattern list.
- [x] Verified: `generate-configs.sh`, `check-parity.sh`, `health-check.sh --verbose`, `test-agents.sh`
      all green.

---

### Epic 39 — Saturday concept expansion: Partials, Models, Factories, Moist tests (2026-07-17)
> Surfaced by the user sharing a personal reformulation of Saturday's concepts (Site / Pages [with
> Elements, Filters, Partials] / Flows) and asking whether it made sense. Three genuine additions came
> out of that: `BasePartial` was a real gap (`BaseElement` is documented as page-scoped, leaving
> cross-page shared UI like headers/footers/global-nav with no clean home), `Model`/`Factory` were
> framework-wide `CLAUDE.md` conventions that were never explicitly claimed as first-class Saturday
> concepts even though every Saturday test uses them, and the "Moist tests" principle (DRY the setup
> noise, keep the critical assertion path visible — related to DAMP) sharpens the existing "always
> practice TDD/BDD" rule with a nuance the framework practices implicitly but never stated.
- [x] `docs/patterns/saturday-framework-patterns.md`: added "Mental Model" section at the top laying
      out the top-down hierarchy (Site → Pages → Flows, plus Models/Factories as data, plus
      Coordinators). Expanded `BasePage`'s Structure and Related to explicitly name its three
      sub-parts (`BaseElement`, `Filters`, `BasePartial`) rather than treating them as unrelated
      sibling patterns. Expanded Site-Centric Pattern's "Why not POM" to include the "shared header
      across every page" failure mode that `BasePartial` specifically addresses.
- [x] New pattern entries in the same doc: `BasePartial` (Context/Structure/Example/Trade-offs/Related
      with the naming-convention note that partials deliberately do NOT follow `FooPage`), `Model`
      (test-context application of `CLAUDE.md`'s general Models convention), `Factory` (test-context
      application of `CLAUDE.md`'s general Factories convention, referencing the per-language factory
      libraries already documented in `shared/rules/<language>-conventions.md`).
- [x] `BaseFlow` Trade-offs section rewritten to fold in the Moist principle explicitly — DRY the setup
      noise, keep the critical assertion path visible, with reference to DAMP as the related
      established idea.
- [x] `shared/DOMAIN_DICTIONARY.md`: added `BasePartial` row to the Saturday Framework section with
      synonyms to avoid (`SharedComponent`, `LayoutFragment`, `PartialPage`). Clarified `BaseElement`'s
      description to say "within a single BasePage" to make the page-scoped-vs-cross-page distinction
      between `BaseElement` and `BasePartial` unambiguous.
- [x] `shared/rules/testing-conventions.md`: added the Moist principle to Test Quality as an ALWAYS
      rule with a concrete example (search-results test — the search action stays visible, the click
      mechanics don't). Documented-only, no CI check, matching how the other Test Quality items are
      handled.
- [x] `docs/patterns/README.md`: updated the Saturday entry to reflect the top-down hierarchy and
      mention all the new concepts.
- [x] Verified: `generate-configs.sh`, `check-parity.sh`, `health-check.sh --verbose`, `test-agents.sh`
      all green. No agent version bumps this round — pure docs and dictionary additions; agent
      behavior unchanged.
- [x] **Follow-up flagged, not done here**: the Saturday implementation repos themselves
      (`saturday-monorepo` TS, `saturday-monorepo-csharp`, `saturday-monorepo-python`,
      `saturday-monorepo-java`) need `BasePartial` added as a real class/interface to match the
      documented concept. Explicitly out of scope for this epic — a separate feature spec per repo,
      matching the pattern already used for the earlier factory-tooling specs in each of those repos.

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
| No durable memory lifecycle (capture/audit/expire) beyond raw KIs | Epic 22 | 9 |
| Contract coverage incomplete (context-engineer + 5 pipeline agents) | Epic 23 | 9 |
| No single reference for agent Role/Counterbalance/Gap | Epic 24 | 9 |
| Team Topology registry still all placeholders | Epic 25 | 9 |
| documentation-manager overlapped memory-engineering with no review gate | Epic 26 | 9 |
| Local script runs didn't match CI's actual OS/bash version | Epic 27 | 9 |
| No visual diagram of how the framework fits together | Epic 28 | 9 |
| Machine-specific paths, casing drift, stale binary artifacts | Epic 29 | 9 |
| Cursor's Tier 2 classification was stale -- native agents/skills unrecognized | Epic 30 | 9 |
| No preferred-package/structure guidance per language, and Claude Code never received what little existed | Epic 31 | 9 |
| No verification that a fresh install actually produces correct output, only that this repo's own output is in sync | Epic 32 | 9 |
| docs/ accumulated fully-superseded early bootstrap material, and docs/README.md itself was stale | Epic 33 | 9 |
| Pattern catalog covered GoF/Saturday/Sunday/Clean Architecture but not DDD, EIP, Stability Patterns, or 12-Factor -- three of which were already named influences in architect.md with no documentation | Epic 34 | 9 |
| STRIDE and SLI/observability principles were fully worked out inside security-reviewer.md/sre-engineer.md but never extracted as standalone references; Expand/Contract and API design guardrails existed only as rules, not documented patterns | Epic 35 | 9 |
| test-driven-developer started every run cold (no KI lookup, no memory feedback loop); no standalone agent existed for adding tests to code that must not change (coverage backfill / legacy characterization) | Epic 36 | 9 |
| No ATDD-shaped delivery workflow that supports trust-progressive phase-out of mechanical review gates -- teams either ran full deliver-feature (heavier than needed) or hand-orchestrated qa-engineer/test-driven-developer every time (no persistence, no trust curve visible in repo) | Epic 37 | 9 |
| Testing taxonomy structurally implicit but never named as a pyramid; no in-test annotation convention for issue/AC traceability (report-time-only, evaporates on rename); framework silently equated single-agent test-first with XP TDD's role-separated design discipline | Epic 38 | 9 |
| Saturday had no dedicated concept for cross-page shared UI (headers/footers/global-nav awkwardly stuffed into BaseElement or base classes); Model/Factory were framework-wide but never claimed as first-class Saturday concepts; DRY-vs-DAMP tension in test authoring was practiced implicitly but never stated as a rule | Epic 39 | 9 |

---

## Decisions (resolved)

1. **Global symlinks vs. copies**: Generated configs only — avoids accidental edits to canonical source.
2. **Cursor file references**: Cursor cannot follow "read this file" instructions in `.mdc` rules. LLMs don't dynamically traverse the file tree. All rules must be **inlined** in the `.mdc` body. Best practices: keep rules short and direct, use ALWAYS/NEVER/CRITICAL keywords, break into small focused `.mdc` files per concern, validate YAML frontmatter or Cursor silently ignores them. This means `generate-configs.sh` must **concatenate and inline** shared rules into each `.mdc` file — no file references.
3. **Agent test strategy**: Structural checks (grep for required sections/findings) for CI. LLM-backed tests for manual verification only.
4. **KI storage**: Both — `shared/knowledge/` for universal patterns (portable across machines), `.claude/knowledge/` for project-specific context.
5. **Auto-promotion threshold**: 3 occurrences across different features triggers auto-promotion of a finding to a rule.
