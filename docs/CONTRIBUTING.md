# Contributing

How to add a new agent, skill, rule, or platform to this framework. For the overall design, see
[ARCHITECTURE.md](ARCHITECTURE.md).

Everything you add lives in `shared/` and only `shared/` — never edit a generated platform config directly
(anything in `.cursor/`, `.github/`, `.gemini/`, `.openai.md`, `.windsurfrules`), it'll just get overwritten
by the next `scripts/generate-configs.sh` run.

---

## Adding a new agent

1. Create `shared/agents/<name>.md` following an existing agent as a template. Required frontmatter:
   ```yaml
   ---
   name: your-agent-name
   description: When to use it — be specific enough that an orchestrator or human can decide without guessing.
   tools: Read, Write, Edit, Bash, Glob, Grep   # only what it actually needs
   model: sonnet
   version: 1.0.0
   ---
   ```
2. Write the persona, process, output format, and rules sections — see any existing agent for the expected
   shape. Keep the process numbered and concrete; vague steps produce vague output.
3. If this agent participates in `deliver-feature`, decide where in the pipeline it runs and update
   `shared/skills/deliver-feature/SKILL.md`'s numbered steps accordingly.
4. If its output is consumed by other agents, consider adding a contract in `shared/contracts/` (see
   "Adding a contract" below) and wiring `validate-artifact` into `deliver-feature` after its step.
5. Add a `## 2026-MM-DD — <reason>` entry to `shared/agents/CHANGELOG.md` in the same commit — the
   pre-commit hook (`scripts/hooks/pre-commit`, opt-in) checks for this if it's enabled.
6. Consider a golden-file fixture in `tests/agents/<name>/` plus an `eval-rubric.md` (see
   [editing-agent-prompts.md](runbooks/editing-agent-prompts.md) and the `agent-eval` skill) if this agent is
   likely to regress silently on future edits.
7. Run `scripts/check-parity.sh` and `scripts/generate-configs.sh` — the new agent needs to show up in
   every platform's persona roster.

## Adding a new skill

1. Copy `shared/skills/SKILL_TEMPLATE.md` into `shared/skills/<name>/SKILL.md`.
2. Fill in `triggers.keywords` and `triggers.intentPatterns` — be specific; overly broad keywords cause a
   skill to fire when the user meant something else.
3. Write "When To Use" with explicit "Do NOT use when" cases — this is what lets Claude (or a human)
   disambiguate between similar skills (e.g. `retrospective` vs. `pipeline-retrospective` vs.
   `agent-scorecard` all sound similar but analyze different things).
4. If the skill is meant to run inside `deliver-feature`, wire it into the numbered pipeline steps. If it's
   standalone, it needs no wiring — it's just discoverable by its trigger patterns.
5. Run `scripts/check-parity.sh` — skills aren't part of the persona roster check, but confirm nothing else
   broke.

## Adding a new rule

1. Add or edit a file in `shared/rules/` (`architecture-guardrails.md`, `design-principles.md`,
   `approval-gates.md`, or a new one).
2. **This requires human sign-off** — `.claude/rules/approval-gates.md` Gate #7 ("Wiring a New Fitness
   Function") applies to any change here, since every agent treats these as session-long law. Don't land a
   rule change without the explicit "approve fitness function" or "add to CI" confirmation that gate
   describes.
3. If you added a new rule *file* (not just edited an existing one), update `generate-configs.sh`'s
   `collect_rules()` and, for Cursor, add a new `generate_mdc()` call in `generate_cursor()` so it gets its
   own `.mdc` file (Cursor needs focused, small rule files, not one giant one).
4. Run `scripts/generate-configs.sh` then `scripts/check-parity.sh`.

## Adding a contract (`shared/contracts/`)

Contracts define the required section headings for a pipeline agent's output, enforced structurally by
`validate-artifact` between handoffs (see Epic 5 in
`docs/features/context-engineering-framework/TODO.md` for the original design rationale).

1. Create `shared/contracts/<agent>-contract.md` — list every required heading **exactly** as the producing
   agent's own Output Format template writes it, including heading level (`##` vs `###`).
2. Add any content-level rules beyond "heading exists" (e.g. `qa-contract.md` requires `Failed: 0`,
   `review-contract.md` requires the literal `APPROVED`/`CHANGES REQUESTED` string).
3. Add the agent to `validate-artifact`'s Contract Mapping table and to `scripts/test-agents.sh`'s
   `contract_for_agent()` function.
4. Wire a `validate-artifact` step into `deliver-feature` immediately after the producing agent's step.

## Adding a new platform

See [docs/runbooks/adding-a-new-platform.md](runbooks/adding-a-new-platform.md) — it's involved enough to
warrant its own step-by-step guide (registry entry, generator function, parity check, tier/terminology
decision).

## Editing an existing agent prompt

See [docs/runbooks/editing-agent-prompts.md](runbooks/editing-agent-prompts.md) — versioning, changelog, and
testing requirements.

## Before you commit

- `scripts/check-parity.sh` — no drift between `shared/` and generated configs
- `scripts/test-agents.sh` — if you touched an agent with a `tests/agents/` fixture
- the `agent-eval` skill — if that fixture also has an `eval-rubric.md`, for the qualitative regression check
- `scripts/check-context-budget.sh` — if you touched `context-engineer` or its output format
- A `shared/agents/CHANGELOG.md` entry — if you touched any `shared/agents/*.md` behavior
