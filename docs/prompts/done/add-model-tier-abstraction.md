# Add Portable `model_tier` Abstraction to Agent Frontmatter

Introduce a portable, platform-agnostic way for agents, subagents, and workflows to declare *what class of model they need* (`light`, `default`, `heavy`) without pinning a concrete Anthropic model ID that would break installs on Cursor / Copilot / Gemini. The concrete model is resolved at install time by looking up `shared/model-defaults.yaml` for the target platform. Where the platform can't honor per-agent model selection, the installer strips the field and emits a clear warning — no silent degradation.

## Why this matters

Right now every agent declares `model: inherit` (Claude-Code-shaped). Two problems:

1. **Cost tuning is blocked.** Read-only auditors (`prompt-evaluator`, `context-auditor`, `documentation-auditor`, `pattern-reviewer`) don't need Opus — they're pattern-matching against a rubric. But we can't declare `model: claude-haiku-4-5-*` in shared agents because that would break Cursor/Copilot installs.
2. **No portable declaration of intent.** Even where a platform *does* honor a per-agent model choice, we have no cross-platform vocabulary for "this agent is light-weight" that survives translation to each platform's config.

The fix is a two-layer split: agents declare **abstract tier**; a central mapping file declares **concrete model per platform**; the installer resolves + warns.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prerequisites

- Working tree clean
- User confirmation on the tier taxonomy (see Phase A) before Phase B starts
- User confirmation on the per-agent tier assignments (see Phase C) before Phase C commits land

## Scope

### Phase A — Design doc (one commit)

Draft `docs/aos/model-tier-abstraction-plan.md` with:

- **Tier taxonomy** — proposed values and their operational definitions:
  - `light` — pattern-matching, rubric-scoring, structural validation. Read-only auditors. Small, cheap model.
  - `default` — most producer agents (analyst, developer, code-reviewer, etc.). Session default; on Claude Code this maps to `inherit`.
  - `heavy` — architectural / security reasoning where cost-to-quality tradeoff favors the top-tier model.
- **Frontmatter contract change** — `model_tier` is REQUIRED going forward; `model` becomes optional. If both are present, `model` wins for Claude Code (explicit override); other platforms use `model_tier` since they ignore `model` anyway.
- **Per-platform resolution behavior** — enumerate what happens on Claude Code / Cursor `.cursor/agents/` / Cursor `.cursor/rules/` / Copilot / Gemini Antigravity / Roo Code / Cline. Where the platform doesn't honor per-agent model, resolution is a no-op and the installer emits a WARNING with a link to how the user sets a platform-wide default.
- **User override precedence** — `.claude/model-overrides.yaml` (and the parallel `.cursor/model-overrides.yaml` etc.) is checked *before* `shared/model-defaults.yaml`. Document that user override always wins.
- **Rollout strategy** — how to backfill existing agents without a big-bang. Recommendation: Phase C ships tier assignments per agent in per-file commits with a short justification for the tier picked; agents without an explicit tier default to `default` for one release cycle before the health check upgrades from WARN to FAIL.

Commit: `docs(aos): draft model-tier abstraction plan (Op A)`.

**Pause for user approval on the taxonomy and rollout strategy.**

### Phase B — Schema + registry (multiple commits)

#### B1 — Extend the frontmatter JSON schema

Update `shared/schemas/agent-frontmatter.schema.json` to:
- Add `model_tier` as a required property with `enum: ["light", "default", "heavy"]`
- Keep `model` as optional (was already optional; explicitly document the "if both present, `model` wins on Claude Code" precedence in the description)

Update `shared/contracts/agent-frontmatter-contract.md` with the same information in prose so `validate-artifact` and human reviewers agree.

Commit: `feat(schemas): add model_tier to agent frontmatter schema (Op B1)`.

#### B2 — Create `shared/model-defaults.yaml`

New file structured as:

```yaml
# Central mapping from portable model_tier → concrete model per platform.
# Agents declare `model_tier: light|default|heavy` in their frontmatter.
# The installer resolves the tier against the target platform's row below.
#
# `null` means: this platform does not honor per-agent model selection.
#   The installer strips the field from the exported agent and emits a
#   WARNING pointing at the platform's own model-selection docs.

claude_code:
  light: claude-haiku-4-5-20251001
  default: inherit
  heavy: claude-opus-5

cursor:
  # Cursor's `.cursor/agents/` subagents accept a model field but the
  # per-agent honoring is inconsistent across Cursor versions — treat as
  # advisory only. Users can pin via Cursor Settings → Model.
  light: null
  default: null
  heavy: null

copilot:
  # GitHub Copilot has no per-agent model API. Global model is picked in
  # VS Code Copilot settings.
  light: null
  default: null
  heavy: null

gemini_antigravity:
  # Gemini/Antigravity has no per-agent model selection.
  light: null
  default: null
  heavy: null

roo_code:
  # TBD — depends on Epic 42 landing. Placeholder null for now.
  light: null
  default: null
  heavy: null

cline:
  # TBD — depends on Epic 42 landing. Placeholder null for now.
  light: null
  default: null
  heavy: null
```

Include an inline pointer to `shared/platform-registry.json` explaining that per-platform capabilities are the source of truth for what each platform actually supports.

Commit: `feat(model-defaults): add tier→model mapping registry (Op B2)`.

#### B3 — Update `shared/templates/agent.template.md`

Add `model_tier: default` to the template frontmatter block (immediately below `model: inherit`), and add a row to the "Frontmatter reference" table explaining the field, values, and precedence with `model:`. Reference `shared/model-defaults.yaml` for the mapping.

Commit: `docs(templates): document model_tier field in agent template (Op B3)`.

### Phase C — Backfill existing agents (one commit per agent group)

For each agent in `shared/agents/`, add a `model_tier:` line to its frontmatter matching its operational profile. Grouping suggestion (commit per group, not per agent — 4-5 commits total):

- **`light` tier (read-only auditors + evaluators)**: `context-auditor`, `documentation-auditor`, `knowledge-auditor`, `memory-auditor`, `pattern-reviewer`, `privacy-auditor`, `prompt-evaluator`, `retrieval-evaluator`, `rule-auditor`, `tool-validator`, `agent-evaluator`
  - Commit: `feat(agents): tag read-only auditors as model_tier: light (Op C1)`
- **`default` tier (most producers)**: `analyst`, `developer`, `code-reviewer`, `qa-engineer`, `tech-writer`, `devops-engineer`, `accessibility-engineer`, `data-engineer`, `performance-engineer`, `sre-engineer`, `dependency-auditor`, `release-manager`, `unit-tester`, `test-driven-developer`, `documentation-manager`, `spec-writer`, `context-engineer`, `product-owner`, `dx-engineer`, `finops-engineer`, `chaos-engineer`, `api-test-generator`, `modernization-supervisor`, `claude`
  - Commit: `feat(agents): tag producer agents as model_tier: default (Op C2)`
- **`heavy` tier (deep reasoning)**: `architect`, `security-reviewer`
  - Commit: `feat(agents): tag architect + security-reviewer as model_tier: heavy (Op C3)`
- **Any custom / plugin agents outside the above list**: assign per operational profile and commit as `Op C4`.

For each agent, add a short comment ABOVE the `model_tier:` line explaining WHY that tier was picked (e.g., `# Rubric-driven read-only auditor — no code generation`). Not a docstring, a single-line YAML `#` comment.

**Pause between C1, C2, C3 for user approval of the tier picks.** Cost/quality tradeoffs are opinions, not facts — the user must sign off.

### Phase D — Installer wire-up (one commit)

Update `install.sh` (and any per-platform sub-installers) to:

1. Read each shared agent's `model_tier` at install time.
2. Look up the resolved model in `shared/model-defaults.yaml` for the target platform.
3. If resolved value is a concrete model ID: write `model: <resolved-id>` into the exported agent copy.
4. If resolved value is `null` for that platform: strip the `model:` line from the exported agent AND emit a WARNING to stdout:
   ```
   WARN: agent <name> requested model_tier: <tier>
         platform <platform> does not honor per-agent model selection.
         Set your global model preference via <link to platform docs>.
   ```
5. Check `.claude/model-overrides.yaml` (or platform-equivalent path) FIRST — user overrides win over `shared/model-defaults.yaml`.

Commit: `feat(install): resolve model_tier per target platform (Op D)`.

### Phase E — Health check + counter-agent (one commit each)

#### E1 — Extend `scripts/health-check.sh`

Add checks:
- Every agent under `shared/agents/*.md` MUST have `model_tier:` set.
- `model_tier:` value MUST be one of `light`, `default`, `heavy`.
- If `model:` is also present, warn that Claude Code will use `model:` but other platforms will use `model_tier:`.
- WARN (not FAIL) for one release cycle; upgrade to FAIL in the next release per the Phase A rollout strategy.

Commit: `feat(health-check): validate model_tier frontmatter (Op E1)`.

#### E2 — Add `model-tier-auditor` counter-agent

New read-only agent at `shared/agents/model-tier-auditor.md` following the counter-agent pattern (`Read, Glob, Grep` only). Scans:
- Missing `model_tier:` on any shared agent
- Tier value not in the enum
- Tier assignment that looks wrong for the agent's actual behavior (e.g., an agent that produces code tagged `light`; an auditor tagged `heavy`) — heuristic-driven, produces findings for human review, never mutates.

Register in `docs/AGENT_REFERENCE.md` alongside the other counter-agents.

Commit: `feat(agents): add model-tier-auditor counter-agent (Op E2)`.

## Discipline (non-negotiable)

- One commit per op (per phase step, or per agent group in Phase C).
- Conventional Commits.
- **NEVER `git add -A`.**
- Do NOT push.
- After Phase A: pause for user approval on taxonomy + rollout.
- Between Phase C1/C2/C3: pause for user approval on tier assignments.

## Escalation criteria

Stop and report if:
- The Phase A taxonomy discussion surfaces a fourth tier the framework actually needs (e.g., `vision` for screenshot-reading agents) — halt, describe, get approval before continuing.
- A per-agent tier assignment in Phase C is genuinely ambiguous (an agent that behaves like `default` on some features and `heavy` on others) — halt, describe, ask the user to disambiguate rather than guess.
- The installer resolution logic in Phase D would need to change more than the shared agents themselves (e.g., generated Copilot instruction files also need tier annotation somehow) — halt, describe, may need a scope expansion.
- `shared/platform-registry.json`'s existing platform list disagrees with the platform list in `shared/model-defaults.yaml` — halt, reconcile.

## Report format (per phase, under 300 words)

### Phase A report
```
Plan commit: <sha>
Tier taxonomy confirmed: light | default | heavy | <any additions>
Rollout strategy: <WARN-then-FAIL cadence>
Platform resolution behavior documented: <count of platforms covered>
Open questions raised for user: <list>
```

### Phase B report
```
B1 (schema): <sha>
B2 (model-defaults.yaml): <sha> — <count> platforms mapped, <count> null
B3 (template): <sha>
```

### Phase C report (per commit)
```
Op C<n>: <sha>
Agents tagged: <count>
Tier picked: <light | default | heavy>
Justification comments added: verified
User-approved before commit: yes
```

### Phase D report
```
Installer commit: <sha>
Platforms with resolution logic: <list>
Warning emission verified against a null-mapped platform: yes/no
Override file precedence verified: yes/no
```

### Phase E report
```
E1 (health-check): <sha> — new checks added: <list>
E2 (counter-agent): <sha> — model-tier-auditor registered in AGENT_REFERENCE.md: yes
```

Go.
