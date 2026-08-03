# Framework Improvement Prompts

Self-contained agent handoff prompts for framework-level improvements to `ai-assistant-dot-files`. Each file is a fully-briefed prompt a fresh Claude Code chat or subagent can execute standalone.

Separate from `docs/aos/prompts/` (AOS-migration-specific) — these are general framework improvements not tied to the AOS phased migration.

## How to use

1. Open a fresh Claude Code chat or spawn a subagent.
2. Copy the entire target prompt file into the message.
3. The agent has everything needed: repo path, prior state, scope, guardrails, commit discipline, escalation criteria, report format.

For tasks that require a human (not a fireable agent prompt), see [../human-tasks.md](../human-tasks.md).

## Active Prompts

### Cross-project setup + maintenance patterns

| File | Scope | Estimated size |
|---|---|---|
| [install-framework-with-mcp-bridge.md](install-framework-with-mcp-bridge.md) | Install framework into a target project **AND** bridge selected framework tools (analyze_complexity, check_accessibility, search_ki, search_docs, etc.) into that project's existing MCP server. Follows Option B (copy source, keep MCP sovereign, adjust paths + monorepo `packagePath` filter for the target's layout). Reusable pattern for projects with Saturday-like structure but not using Saturday itself | Medium-large — Phase A investigate + Phase B plan + Phase C multi-op execution; ~8-10 commits in the target repo |
| [update-installed-framework.md](update-installed-framework.md) | Propagate framework updates into an already-installed target project. Auto-detects which of 4 install patterns (traditional symlink / bridge-copy / saturday-mcp-adopted / generated platform configs) apply — often multiple — and runs the correct update sequence for each. Includes contract-drift check for previously-delivered artifacts | Medium — Phase A detect + Phase B update per pattern + Phase C contract-drift check; commit count depends on pattern mix |
| [add-model-tier-abstraction.md](add-model-tier-abstraction.md) | Add a portable `model_tier:` field to agent frontmatter (`light` / `default` / `heavy`) + a central `shared/model-defaults.yaml` mapping tier → concrete model per platform. Install script resolves the tier for the target platform (or strips + warns when the platform doesn't honor per-agent model). Unlocks cost-tier optimization without breaking Cursor/Copilot/Gemini installs | Medium-large — Phase A design + Phase B schema/registry + Phase C agent backfill + Phase D installer wire-up + Phase E counter-agent |

### Phase 10 Roadmap Epics (from `docs/audits/framework-gap-audit-2026-07-25.md`)

Ten Epic-numbered handoffs drafted 2026-07-26 from that audit. Recommended execution order below matches the audit's priority ranking.

| Epic | File | Estimated size |
|---|---|---|
| ~~**53**~~ | ~~[epic-53-inventory-drift-check.md](done/epic-53-inventory-drift-check.md)~~ — shipped `41d9c27` | ✓ Done |
| ~~**42**~~ | ~~[epic-42-roo-code-cline-platform.md](done/epic-42-roo-code-cline-platform.md)~~ — shipped `0e47b2f`→`523d4e1` | ✓ Done |
| ~~**44**~~ | ~~[epic-44-jetbrains-exporter.md](done/epic-44-jetbrains-exporter.md)~~ — shipped `d047e69`→`3e953df` | ✓ Done |
| ~~**46**~~ | ~~[epic-46-visual-qa-engineer.md](done/epic-46-visual-qa-engineer.md)~~ — shipped `472c961`→`5d28b74` | ✓ Done |
| ~~**48**~~ | ~~[epic-48-iac-conventions.md](done/epic-48-iac-conventions.md)~~ — shipped `aa954e2` | ✓ Done |
| ~~**49**~~ | ~~[epic-49-mobile-conventions.md](done/epic-49-mobile-conventions.md)~~ — shipped `aea8172`→`0bc5107` | ✓ Done |
| ~~**50**~~ | ~~[epic-50-rust-conventions.md](done/epic-50-rust-conventions.md)~~ — shipped `e67653a` | ✓ Done |
| ~~**51**~~ | ~~[epic-51-enterprise-memory-sync.md](done/epic-51-enterprise-memory-sync.md)~~ — shipped `146e3bb`→`69d6c3d` | ✓ Done |
| ~~**55**~~ | ~~[epic-55-agent-eval-expansion.md](done/epic-55-agent-eval-expansion.md)~~ — shipped batch 1 + batch 2 + Phase C | ✓ Done |
| ~~**57**~~ | ~~[epic-57-context-manifest-fixtures.md](done/epic-57-context-manifest-fixtures.md)~~ — shipped | ✓ Done |

**Not drafted** — per-epic dispositions live in `docs/audits/framework-gap-audit-2026-07-31.md` §F5 (the 07-25 audit itself carries no inline reasons): Epics 52, 54, and 56 are subsumed by AOS Phase 3 (`docs/aos/prompts/phase-3-runtime.md`) and must not be drafted separately. Epics 43, 45, 47, and 58 were drafted 2026-07-31 — see the table below.

### Remaining standalone epics (drafted 2026-07-31, from `framework-gap-audit-2026-07-31.md` §F5)

Recommended execution order matches the 07-31 audit's priority ranking (47 first, 43 last).

| Epic | File | Estimated size |
|---|---|---|
| ~~**47**~~ | ~~[epic-47-ship-feature.md](done/epic-47-ship-feature.md)~~ — shipped `58a429c`→`f94e5c0` (incl. tool-validator fix) | ✓ Done |
| ~~**58**~~ | ~~[epic-58-documentation-auditor-automation.md](done/epic-58-documentation-auditor-automation.md)~~ — shipped `98159fe`→`ba34794` | ✓ Done |
| **45** | [epic-45-refactor-engineer.md](epic-45-refactor-engineer.md) — `refactor-engineer` agent + refactoring contract + fixture + full 38→39 roster ripple | Large — Phase A design pause, then ~7 commits |
| **43** | [epic-43-mcp-exporter.md](epic-43-mcp-exporter.md) — `shared/mcp/` exporter; Phase A must first rule generator-vs-manifest-vs-superseded (audit flagged scope-clarification needed) | Unknown until Phase A ruling — possibly closes as superseded |

### Project-as-RAG optimization pair (drafted 2026-07-31, from retrieval-optimization discussion)

Exploits the fact that the pipeline *generates* the installed project's corpus — enrich at write
time so every retrieval tier (today's lexical, ADR-002's future BM25/vector) inherits a better
corpus. 59 is independent and highest-leverage; 60 is best run after it. Neither depends on AOS
Phase 3.

| Epic | File | Estimated size |
|---|---|---|
| **59** | [epic-59-retrieval-write-time-enrichment.md](epic-59-retrieval-write-time-enrichment.md) — artifact retrieval frontmatter (WARN-first), domain-dictionary query expansion in search-ki/query-memory, persisted summary surrogates, citation-link convention | Medium — 4 commits, no design pause |
| **60** | [epic-60-retrieval-index-freshness-eval.md](epic-60-retrieval-index-freshness-eval.md) — index-freshness hook examples (against saturday-mcp's reindex), registry retrievalBackends fitness function (ADR-002's deferred check), CODEMAP generator for the source tier, telemetry-sourced retrieval regression set | Medium — 4 commits, per-op halt conditions |

### Structural-gap epics (drafted 2026-07-31, from `framework-gap-audit-2026-07-31.md` §3b)

Eight gaps from the same-day structural review — not backlog items but blind spots ("what is the
framework missing"). Table order IS the priority order.

| Epic | File | Estimated size |
|---|---|---|
| **61** | [epic-61-prompt-regression-harness.md](epic-61-prompt-regression-harness.md) — headless eval runner over the golden-file fixtures; catches model-change behavior drift across all 38 agents. Assembly, not invention | Medium-large — Phase A design pause (runner/judge/cost rulings), then 4 commits |
| ~~**62**~~ | ~~[epic-62-gate-decision-cost-telemetry.md](done/epic-62-gate-decision-cost-telemetry.md)~~ — shipped `2b7ed81`→`e5218c8` | ✓ Done |
| **63** | [epic-63-parallel-delivery-isolation.md](epic-63-parallel-delivery-isolation.md) — break the feature-workspace singleton; workspace-per-feature vs git-worktree ruling; legacy resume must survive | Large — Phase A design pause, blast-radius inventory, then per-subsystem commits |
| ~~**64**~~ | ~~[epic-64-shipped-linter-configs.md](done/epic-64-shipped-linter-configs.md)~~ — shipped `2eb31cc`→`c245bcb` | ✓ Done |
| ~~**65**~~ | ~~[epic-65-framework-threat-model.md](done/epic-65-framework-threat-model.md)~~ — shipped `839b610`→`2699b03` | ✓ Done |
| ~~**66**~~ | ~~[epic-66-capability-inventory-lifecycle.md](done/epic-66-capability-inventory-lifecycle.md)~~ — shipped `b95df04`→`2168505` | ✓ Done |
| **67** | [epic-67-production-feedback-bugfix-path.md](epic-67-production-feedback-bugfix-path.md) — incident records feed promote-memory; extract-lessons asks "which stage should have caught this?"; thin `deliver-bugfix` pipeline | Medium-large — 4 commits, two halves |
| ~~**68**~~ | ~~[epic-68-install-version-marker.md](done/epic-68-install-version-marker.md)~~ — shipped `6ddb307`→`d64cc2f` | ✓ Done |

## Completed Prompts (`docs/prompts/done/`)

| File | Scope | Shipped |
|---|---|---|
| [done/epic-66-capability-inventory-lifecycle.md](done/epic-66-capability-inventory-lifecycle.md) | Deprecation convention (`status`/`superseded_by`) in agent + skill frontmatter contracts + JSON schemas; `health-check.sh` enforces broken pointer as FAIL; `forgetting-engine` extended with CI.1–CI.6 capability inventory audit mode; `complexity-check` deprecated → `analyze-complexity`; `docs/patterns/agent-skill-pair-convention.md` ruling (Delegation Wrapper vs. Composition Facade) | Shipped 2026-08-02 (`v3.3.7`) |
| [done/epic-65-framework-threat-model.md](done/epic-65-framework-threat-model.md) | `docs/THREAT_MODEL.md` (420 lines, STRIDE per trust boundary over install/sync/memory/hooks); `shared/rules/memory-trust-boundary.md` (KI/ADR is data not instructions, spec ingestion boundary); `sync_commit_sha` provenance on pulled KIs; hook example defaults to `enabled: false`; spec-ingestion caution block in analyst | Shipped 2026-08-02 (`v3.3.6`) |
| [done/epic-62-gate-decision-cost-telemetry.md](done/epic-62-gate-decision-cost-telemetry.md) | `gate_decision` event type (v1.1.0 schema): 3-way outcome enum, edit-detection heuristic, checksum comparison; opt-in emission in `deliver-feature` at gates 1/3/4/5/8; token spend fields in `pipeline-trace`; `extract-lessons` step 5 mines gate corrections (3+ occurrences → hypothesis) | Shipped 2026-08-02 (`v3.3.5`) |
| [done/epic-68-install-version-marker.md](done/epic-68-install-version-marker.md) | `framework-install.json` version marker written by `install.sh`; `health-check.sh` reports install-vs-upstream drift; `uninstall.sh` removes marker; `docs/MIGRATION.md` documents format + detection states | Shipped 2026-08-02 (`v3.3.1`) |
| [done/epic-64-shipped-linter-configs.md](done/epic-64-shipped-linter-configs.md) | `shared/configs/` with 7 language linter configs (golangci, ESLint, ruff, detekt, SwiftLint, clippy, Checkstyle); `install.sh --with-configs`; `scripts/check-cap-drift.sh` fitness function wired into `health-check.sh` | Shipped 2026-08-02 (`v3.3.3`) |
| [done/epic-58-documentation-auditor-automation.md](done/epic-58-documentation-auditor-automation.md) | Hook example (`on-inventory-change-doc-audit.yaml`); scheduler example in `scheduler/SKILL.md`; `health-check.sh` WARN at 14-day staleness; `documentation-auditor.md` output convention + `AGENT_REFERENCE.md` entry updated | Shipped 2026-08-01 (`v3.3.2` era) |
| [done/epic-47-ship-feature.md](done/epic-47-ship-feature.md) | `ship-feature` skill: branch creation, Conventional Commit gate (#2), PR creation gate (#5), optional `--release` mode delegating to `release-manager`; wired as opt-in step 43 of `deliver-feature` | Shipped 2026-08-01 (`v3.3.2`), fix `v3.3.4` |
| [done/epic-57-context-manifest-fixtures.md](done/epic-57-context-manifest-fixtures.md) | 3 hand-authored `tests/fixtures/context-manifests/` fixtures (passing, warning-no-cuts, missing-status); `check-context-budget.sh` rewritten to scan both real manifests and fixtures, make missing-Status a FAIL instead of SKIP, and document expected fixture outcomes | Shipped 2026-07-31 |
| [done/epic-55-agent-eval-expansion.md](done/epic-55-agent-eval-expansion.md) | Golden-file fixture expansion from 5/38 to 32/38 agents (84%). Batch 1: 14 pipeline agents (developer, sre-engineer, tech-writer, devops-engineer, performance-engineer, data-engineer, accessibility-engineer, visual-qa-engineer, context-engineer, spec-writer, product-owner, release-manager, test-driven-developer, unit-tester). Batch 2: 13 counter/auditor agents (agent-evaluator, context-auditor, documentation-auditor, documentation-manager, knowledge-auditor, memory-auditor, model-tier-auditor, pattern-reviewer, privacy-auditor, prompt-evaluator, retrieval-evaluator, rule-auditor, tool-validator). Phase C: `test-agents.sh` contract_for_agent() expanded to all 33 covered agents; `health-check.sh` fixture-coverage enforcement section added (FAIL when non-deferred agent lacks fixture directory) | Shipped 2026-07-31 |
| [done/epic-51-enterprise-memory-sync.md](done/epic-51-enterprise-memory-sync.md) | ADR-003 (separate git repo model) + `scripts/sync-memory.sh` (pull/push with --confirm gate) + `install.sh --sync-memory` flag + `memory-registry.json` enterpriseSync stanza + README workflow docs + `health-check.sh` Memory Sync section | Shipped 2026-07-31 in commits `146e3bb`→`69d6c3d` |
| [done/epic-50-rust-conventions.md](done/epic-50-rust-conventions.md) | Add `shared/rules/rust-conventions.md` covering Clean Architecture layer separation, cargo + clippy as CI fitness function, clippy::cognitive_complexity capped at 6, thiserror/anyhow error strategy, tokio async runtime, #![forbid(unsafe_code)], proptest + rstest + mockall + criterion | Shipped 2026-07-30 in commit `e67653a` |
| [done/epic-49-mobile-conventions.md](done/epic-49-mobile-conventions.md) | Add `shared/rules/swift-conventions.md` (SwiftUI + Swift Concurrency, SwiftLint `< 7`, XCTest + Nimble + swift-snapshot-testing) and `shared/rules/kotlin-conventions.md` (Jetpack Compose + Coroutines, detekt `< 7`, JUnit 5 + MockK + Paparazzi). No agent wiring today — mobile features will pick up on first use | Shipped 2026-07-30 in commits `aea8172`→`0bc5107` |
| [done/epic-48-iac-conventions.md](done/epic-48-iac-conventions.md) | Add `shared/rules/iac-conventions.md` covering Terraform/OpenTofu, Dockerfile, Kubernetes/Helm, and GitHub Actions guardrails. Wired into `devops-engineer.md` pre-read block. Propagated to all platform generated configs | Shipped 2026-07-30 in commit `aa954e2` |
| [done/epic-46-visual-qa-engineer.md](done/epic-46-visual-qa-engineer.md) | Add `visual-qa-engineer` pipeline agent. Wraps `@orieken/saturday-ml-analyzer` (interaction heatmap cold-spot analysis) and Playwright `toHaveScreenshot()` (visual regression). Conditional: UI features with `heatmap-data/` or snapshot baselines. Adds contract, template, pipeline slot after qa-engineer (steps 28-29), workflow diagram, validate-artifact mapping | Shipped 2026-07-30 in commits `472c961`→`5d28b74` |
| [done/epic-44-jetbrains-exporter.md](done/epic-44-jetbrains-exporter.md) | Add JetBrains AI Assistant + Junie as a single platform entry. Generates `.aiassistant/rules/` (10 files, IDE mode hints) and `.junie/guidelines.md`. Confirmed Junie already reads root `AGENTS.md` so partial support existed — this adds the AI Assistant project-rules layer and explicit Junie path | Shipped 2026-07-30 in commits `d047e69`→`3e953df` |
| [done/epic-42-roo-code-cline-platform.md](done/epic-42-roo-code-cline-platform.md) | Add Roo Code (custom modes via `.roomodes` YAML, 37 agents mapped) and Cline (`.clinerules/` directory, 10 path-scoped files) as platform targets. Registry entries, config generators, parity checks, install detection, README table — full Phase A+B implementation | Shipped 2026-07-30 in commits `0e47b2f`→`523d4e1` |
| [done/epic-53-inventory-drift-check.md](done/epic-53-inventory-drift-check.md) | Add `scripts/check-inventory-drift.sh` — counts actual `shared/` files and greps authoritative prose docs for stated counts; wired into `health-check.sh` as WARN-level. Corrected 5 pre-existing drifts in README and AGENT_REFERENCE | Shipped 2026-07-30 in commit `41d9c27` |
| [done/add-mcp-patterns-directory.md](done/add-mcp-patterns-directory.md) | Extract framework-generic MCP tool patterns (retrievers, analyzers, walkutil, 6 M1 tools, response structs, registration pattern) from `saturday-mcp` into `shared/mcp-patterns/`, then re-point the bridge/update prompts at that directory. Broke the accidental coupling that forced downstream users to also clone saturday-monorepo | Shipped 2026-07-27 in commits `39747a5` → `349281c` (Ops A → D3 + Op C follow-up) |
| [done/write-blog-posts.md](done/write-blog-posts.md) | Draft dev.to + LinkedIn blog posts covering recent framework + agent developments (Posts 05–10) | Shipped 2026-07-25 in commits `1abd94f` → `47dccce` |
| [done/implement-automated-delivery-tier-a.md](done/implement-automated-delivery-tier-a.md) | IMPLEMENT interim Tier A auto-proceed & Tier B contract retries in `/deliver-feature` skill based on `.claude/delivery-policy.yaml` | Shipped 2026-07-25 in commit `25b1a50` |
| [done/automate-deliver-feature.md](done/automate-deliver-feature.md) | DESIGN automated `/deliver-feature` workflow (Tier A/B/C classification, policy schema, interim automation, design doc + implementation handoff prompt) | Shipped 2026-07-25 in commits `04ea74d` → `4d13163` |
| [done/capture-session-history.md](done/capture-session-history.md) | Capture AOS foundations + frontmatter + MCP retrofit session history via `/retrospective` and `/extract-lessons` | Shipped 2026-07-25 in commits `e0a48e8` → `61a17e7` |
| [done/framework-hygiene-sweep.md](done/framework-hygiene-sweep.md) | Eight small pending items: commit `.gitignore` mod, track blog drafts, delete redundant AOS zip, add KI template, adopt `done/` convention, investigate audits doc, fix 2 pre-existing health-check WARNs | Shipped 2026-07-24 in commits `5a4a35d` → `9717fca` |
| [done/add-frontmatter-contracts.md](done/add-frontmatter-contracts.md) | Add `shared/contracts/agent-frontmatter-contract.md` + `skill-frontmatter-contract.md` + `ki-frontmatter-contract.md` so `validate-artifact` can grep-check these | Shipped 2026-07-22 in commits `ae1e440` → `38b14a9` |
| [done/add-frontmatter-json-schemas.md](done/add-frontmatter-json-schemas.md) | Add `shared/schemas/agent-frontmatter.schema.json` etc. Enables IDE autocomplete + enum-value validation. Wire into VS Code / Cursor settings templates | Shipped 2026-07-22 in commits `fdfdd9b` → `a522e14` |

## Convention

Every prompt in this directory:
- Names the target repo path (`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`)
- Restates that this repo IS the git repo — commits land here directly, not in a parent
- Lists the exact ops in scope with clear file paths
- Enumerates guardrails (commit per op, conventional commits, `git add` explicit paths only)
- Names escalation criteria (when to stop and ask)
- Requests a specific report format

When a prompt is executed and shipped, move it to `docs/prompts/done/` and update the table above.

