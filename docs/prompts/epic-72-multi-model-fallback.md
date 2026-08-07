# Epic 72 — Dynamic Provider Fallback & Multi-Model Orchestration

Source: `docs/audits/framework-audit-2026-08-07.md` §3 item 4.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Prior context

`shared/model-defaults.yaml` maps three tiers to concrete model IDs:

```yaml
claude_code:
  light:   claude-haiku-4-5-20251001
  default: inherit
  heavy:   claude-opus-5-...
```

`scripts/resolve-model-tier.py` reads this file and returns a single model ID for a given
tier. Agents declare `model_tier: light|default|heavy` in frontmatter.

No fallback exists. If the resolved model hits a rate limit or provider outage during a long
pipeline run (`deliver-feature` can span 20+ agent phases), the run stalls or errors with
no alternative path.

Related files:
- `shared/schemas/model-defaults.schema.json` — JSON schema for `model-defaults.yaml`
- `shared/agents/*.md` — all declare `model_tier`
- `scripts/run-agent-evals.sh` — uses `resolve-model-tier.py` at eval time

## Scope

**Phase A — Schema Design (one commit, then PAUSE for user approval):**

Draft and commit as `docs(orchestration): provider fallback schema design (Epic 72 Phase A)`:

Produce `docs/patterns/multi-model-fallback-design.md` answering these rulings:

1. **Schema extension**: how does `model-defaults.yaml` express fallback chains? Options:
   - A `fallback:` list under each tier (ordered, first available wins)
   - A top-level `providers:` section with per-provider config + per-tier ordering
   - Propose one option with the full YAML shape.

2. **Resolver behavior**: does `resolve-model-tier.py` return the full chain (for the
   caller to iterate), or does it probe provider health and return the best available model?
   Decision: the script must not make network calls (that violates KISS and adds latency);
   it returns the full chain and the caller decides retry policy.

3. **Invocation surface**: where in the delivery pipeline is the fallback actually used?
   Identify the 1–3 places where a model call is initiated that could benefit from a
   fallback chain. The framework is skill-prompts executed by the LLM, not direct API
   calls — confirm whether fallback applies at the skill-routing layer, the `run-agent-evals.sh`
   layer, or elsewhere.

4. **Multi-provider model IDs**: confirm the canonical model ID strings for the proposed
   fallback providers (Gemini 1.5 Pro, GPT-4o, etc.) that Claude Code or the framework
   runner can actually invoke. If no cross-provider invocation mechanism exists in the
   current runner, note this as a hard constraint.

**Phase B — Implementation (after approval; one commit per op):**

Op 1 — `feat(orchestration): extend model-defaults.yaml with fallback chains (Epic 72 Op 1)`:
- Update `shared/model-defaults.yaml` with the schema chosen in Phase A.
- Add fallback chains for all three tiers.
- Update `shared/schemas/model-defaults.schema.json` to validate the new shape.

Op 2 — `feat(scripts): resolve-model-tier.py fallback chain output (Epic 72 Op 2)`:
- Update `scripts/resolve-model-tier.py` to accept `--chain` flag: returns JSON array
  of model IDs in priority order (`["claude-sonnet-5", "gemini-1.5-pro", "gpt-4o"]`).
- Default behavior (no `--chain`) unchanged: returns the primary model ID string.
- Unit test the chain output with a fixture `model-defaults.yaml`.

Op 3 — `feat(evaluation): run-agent-evals.sh fallback chain support (Epic 72 Op 3)`:
- Update `scripts/run-agent-evals.sh` to use `--chain` output: on `claude` invocation
  failure (non-zero exit), retry with the next provider in the chain (up to 2 retries).
- Log each fallback: `[FALLBACK] <agent>: primary failed, retrying with <model>`.
- If all providers fail, emit SKIP (not FAIL) and log the full chain tried.

Op 4 — `feat(health-check): validate fallback chain schema (Epic 72 Op 4)`:
- Add a FAIL-level check in `scripts/health-check.sh`: `model-defaults.yaml` must be
  schema-valid against the updated `model-defaults.schema.json`.
- Add a WARN-level check: each tier should have at least one fallback defined.
- `bash scripts/health-check.sh` green.

Op 5 — `docs(orchestration): update model-defaults docs (Epic 72 Op 5)`:
- Update any prose in `docs/ARCHITECTURE.md` or `docs/AGENT_REFERENCE.md` that describes
  the model tier system to mention fallback chains.
- Add a note to `shared/model-defaults.yaml` header comment explaining the fallback
  contract.

After every commit: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- If Phase A ruling 4 concludes that no cross-provider invocation mechanism exists (Claude
  Code only calls Anthropic; no Gemini/OpenAI path), halt Phase A with that finding. The
  fallback chain can still be stored in `model-defaults.yaml` as forward-compatible
  metadata, but Op 3 (run-agent-evals retry) would be limited to Anthropic-model fallbacks
  only. Document this constraint and ship Op 1–2 only.
- If updating `resolve-model-tier.py` to support `--chain` would push the script's
  cyclomatic complexity above 6 (the framework-wide cap), Extract the chain-resolution
  logic into a separate helper function before extending it.
- If `model-defaults.schema.json` does not exist today, create it in Op 1 rather than Op 4
  so subsequent ops can validate against it.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A rulings:
  - Schema shape: <tier.fallback list | providers section | other>
  - Resolver output: <chain array | primary only (constraint found)>
  - Invocation surface: <run-agent-evals only | pipeline orchestrator | other>
  - Cross-provider constraint: <none found | hard constraint: Anthropic-only runner>

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, resolve-model-tier --chain <correct output>,
run-agent-evals fallback retry <exercised in test>.
```

Go.
