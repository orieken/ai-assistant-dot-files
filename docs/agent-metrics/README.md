# Agent Metrics

Monthly quality scorecards for the pipeline agents, produced by the `agent-scorecard` skill
(`shared/skills/agent-scorecard/SKILL.md`).

## Convention
```
docs/agent-metrics/scorecard-<YYYY-MM>.md
```

One file per calendar month it was generated for. Each scorecard compares against the previous month's file
to flag agents that are improving, stable, or degrading — see `agent-scorecard`'s "Metric Definitions"
section for exactly what's measured and why.

This is distinct from:
- `docs/features/<name>/retrospective.md` — a single delivery's qualitative narrative.
- `docs/pipeline-retrospectives/` — cross-delivery timing/iteration trends (not quality judgments).
- `docs/agent-metrics/evals/` (`agent-eval`) — a fixed fixture + rubric regression check for one agent right
  after a prompt edit, not a monthly trend across real deliveries. See `evals/README.md`.

`agent-scorecard` is the only one of these that says whether an agent's *output* was actually good across
real deliveries, not just how long it took, how the process felt for one feature, or how it handles one
fixed test case.
