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

`agent-scorecard` is the only one of the three that says whether an agent's *output* was actually good, not
just how long it took or how the process felt for one feature.
