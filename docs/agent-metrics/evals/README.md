# Agent Evals

Per-agent qualitative eval results produced by the `agent-eval` skill (`shared/skills/agent-eval/SKILL.md`),
run against the fixtures and rubrics in `tests/agents/<agent-name>/`. Each fixture directory (input + fuzzy
patterns + rubric) is an **Agent Eval Case** — see `DOMAIN_DICTIONARY.md`.

## Convention
```
docs/agent-metrics/evals/<agent-name>-eval-<YYYY-MM-DD>.md
```

One file per run. Each run compares its per-criterion results against the most recent prior file for that
agent (if one exists) to flag regressions — a criterion that passed before and fails now, likely caused by
the prompt edit made since.

This is distinct from:
- `docs/agent-metrics/scorecard-<YYYY-MM>.md` (`agent-scorecard`) — monthly quality trend across real
  pipeline deliveries, not a fixed fixture.
- `tests/agents/<agent-name>/actual-output.md` — the raw scratch output for the most recent run, not
  checked in and not dated.

`agent-eval` is the only one of these that specifically answers "did editing this agent's prompt just make
it worse at the exact case it used to handle correctly."
