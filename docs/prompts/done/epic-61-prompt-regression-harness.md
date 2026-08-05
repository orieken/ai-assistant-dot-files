# Epic 61 — Automated Prompt-Regression Eval Harness

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #1). The gap: a model change
silently alters behavior across all 38 agent prompts and nothing detects it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

Every piece exists except the runner:

- `tests/agents/<agent>/` — 32/38 agents have complete fixtures: `input-*` file,
  `expected-patterns.txt` (structural grep checks), `eval-rubric.md` (qualitative criteria).
  6 specialists deferred, documented in `tests/agents/README.md`.
- `actual-output.md` is gitignored (locally generated, non-deterministic); CI's
  `scripts/test-agents.sh` only validates outputs someone generated manually in a live session.
- `agent-eval` skill — already defines the grading flow: act as the agent against its fixture,
  grade against patterns + rubric, flag regressions vs. the last recorded eval.
- `agent-evaluator` agent — the persona form of the same logic. `shared/evaluation/` — where
  regression metrics belong.
- Headless execution options: `claude -p` (Claude Code non-interactive mode) or the Claude Agent
  SDK. Both cost real API money per run.

## Scope

**Phase A — Design (one commit, then PAUSE for user approval):**

Commit `docs(evaluation): design prompt-regression harness (Epic 61 Phase A)` answering:

1. **Runner tech**: `claude -p` subshell per agent vs. an Agent SDK script. Weigh: dependency
   footprint, per-agent model_tier resolution (`shared/model-defaults.yaml` +
   `scripts/resolve-model-tier.py` exist), output capture fidelity.
2. **Grading tech**: pattern checks are deterministic (bash, free); rubric grading needs an LLM
   judge — same run or second cheaper pass? Propose the judge's model tier.
3. **Cost + cadence policy**: estimated tokens for a full 32-agent sweep; propose default cadence
   (on-demand + pre-release, NOT nightly by default) and a `--agents <subset>` flag.
4. **Regression record format**: where results land in `shared/evaluation/` (schema: agent,
   model id, date, pattern pass/fail, rubric grade, delta vs. previous) — the "last recorded
   eval" the agent-eval skill already refers to, made concrete.

**Phase B — Implementation (after approval):**

1. `scripts/run-agent-evals.sh` (+ helper per Phase A ruling) — runs fixtures headless, captures
   `actual-output.md`, runs pattern checks, invokes rubric grading, writes the regression record.
   Commit: `feat(evaluation): headless agent-eval runner (Epic 61 Op 1)`
2. Regression-record schema doc in `shared/evaluation/` + comparison logic (fail the run when an
   agent's pattern-pass or rubric grade drops vs. last record).
   Commit: `feat(evaluation): regression record schema + diffing (Epic 61 Op 2)`
3. Wire an OPTIONAL CI job into `.github/workflows/framework-ci.yml` behind a secret check
   (`if: secrets.ANTHROPIC_API_KEY` present) so forks/PRs without keys skip cleanly.
   Commit: `feat(ci): opt-in agent-eval regression job (Epic 61 Op 3)`
4. Update `tests/agents/README.md` — the "manual-only" caveat becomes "manual OR harness".
   Commit: `docs(tests): document eval harness (Epic 61 Op 4)`

`bash scripts/health-check.sh` green after every commit; the harness itself must SKIP (not FAIL)
when no API key is configured — same lesson as commit `6c422cb`.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push. Never commit API keys or recorded outputs containing them.

## Escalation

- Full-sweep cost estimate exceeds what a reasonable pre-release check should cost — halt with
  numbers and subset options.
- `claude -p` can't faithfully reproduce agent invocation (tool access, system prompt assembly) —
  halt; a harness that tests something other than real behavior is worse than none.
- Rubric grading proves too noisy to diff meaningfully (grade flapping run-to-run) — halt,
  propose pattern-only regression + rubric as advisory.

## Report (under 150 words)

```
Phase A commit: <sha>  Rulings: runner=<>, judge=<>, cadence=<>, record location=<>
Phase B commits: <sha> x4
Sweep cost estimate: ~<tokens/$> for 32 agents
Demo run: <n> agents evaluated, <n> pattern pass, regressions vs. baseline: <n | first-run baseline established>
No-API-key behavior: SKIP confirmed
```

Go.
