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
| **44** | [epic-44-jetbrains-exporter.md](epic-44-jetbrains-exporter.md) — JetBrains + Junie system-prompt exporter | Small-medium (same Phase A/B pattern) |
| **46** | [epic-46-visual-qa-engineer.md](epic-46-visual-qa-engineer.md) — visual-qa-engineer agent + Saturday.ML integration | Medium (Phase A design + Phase B agent/contract/template/pipeline wiring) |
| **48** | [epic-48-iac-conventions.md](epic-48-iac-conventions.md) — Terraform / Docker / K8s / GHA rule doc | Small-medium (1 file + agent wire-up) |
| **49** | [epic-49-mobile-conventions.md](epic-49-mobile-conventions.md) — Swift + Kotlin conventions | Medium (2 files, one commit each) |
| **50** | [epic-50-rust-conventions.md](epic-50-rust-conventions.md) — Rust conventions | Small-medium (1 file) |
| **51** | [epic-51-enterprise-memory-sync.md](epic-51-enterprise-memory-sync.md) — `install.sh --sync-memory` + ADR-003 | Medium-large (design + protocol + tooling) |
| **55** | [epic-55-agent-eval-expansion.md](epic-55-agent-eval-expansion.md) — golden-file coverage for all 36 agents | Medium (audit + fill + enforce) |
| **57** | [epic-57-context-manifest-fixtures.md](epic-57-context-manifest-fixtures.md) — fixtures for `check-context-budget.sh` | Small (1 commit) |

**Not drafted — reasons documented inline in the audit** (superseded, resolved, or scope-clarification needed): Epics 43, 45, 47, 52, 54, 56, 58.

## Completed Prompts (`docs/prompts/done/`)

| File | Scope | Shipped |
|---|---|---|
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

