# Epic 55 — Agent Golden-File Eval Expansion

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 5.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

`tests/agents/` contains golden-file fixtures for a subset of agents (not all 36). The `agent-eval` skill and (now) the `agent-evaluator` agent (post-Phase-2) can regression-test agents against their golden output, but coverage is thin. Adding a new agent today doesn't force a golden-file test to exist.

## Scope

**Phase A — Audit current coverage** (one commit): enumerate which of the 36 agents have golden-file fixtures under `tests/agents/` and which don't. Draft a coverage table to include in the eventual expansion.

Commit: `docs(tests): audit current agent-eval coverage (Epic 55 Phase A)`.

**Phase B — Fill in gaps** (multiple commits, one per agent OR grouped by category — your judgment):

For each un-covered agent:
1. Read the agent's frontmatter + Process to understand what a canonical invocation looks like
2. Craft a minimal input fixture (sample feature spec if the agent is a pipeline agent; sample corpus if the agent is a counter agent)
3. Run the agent, capture the output as the golden file
4. Add the golden to `tests/agents/<agent-name>/expected-output.md`
5. Wire into `tests/agents/README.md` if that's how the current runner discovers goldens

Follow the shape of an existing well-covered agent's fixture directory (pick one from Phase A findings as the template).

Commit per agent (or grouped by category — e.g., "counter agents batch"): `test(agents): add golden fixtures for <name> (Epic 55 Phase B)`.

**Phase C — Make it enforced going forward** (one commit):

Update `scripts/ci-check.sh` (if it exists) or `scripts/health-check.sh` to FAIL when a `shared/agents/<name>.md` exists without a corresponding `tests/agents/<name>/expected-output.md`. This is the fitness function that prevents future coverage drift.

Commit: `feat(ci): enforce agent-eval golden-file coverage (Epic 55 Phase C)`.

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- If capturing a golden requires expensive/non-deterministic operations (network calls, real LLM invocations, large corpus scans) — halt for that specific agent, propose either a mock-based fixture or skip with justification.
- If the current golden-file format has drifted such that some fixtures fail against agents that haven't changed behavior — halt, describe the format issue.
- If Phase B would push coverage above ~25 fixtures created in one prompt run — split into multiple prompts. Overloading one execution session risks the same session-limit issues we hit during mcp-expand M1.

## Report (under 250 words)

```
Phase A commit: <sha>
Pre-expansion coverage: N/36 agents covered
Phase B commits: <count>
Post-expansion coverage: N/36 agents covered
Phase C commit: <sha>
CI enforcement: FAIL when new agent lacks fixture — verified

Any fixtures that were intentionally skipped: <list with justification>
```

Go.
