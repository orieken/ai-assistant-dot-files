# Pipeline Retrospectives

Cross-delivery timing and iteration-count trend reports, produced by the `pipeline-retrospective` skill
(`shared/skills/pipeline-retrospective/SKILL.md`).

## Convention
```
docs/pipeline-retrospectives/retrospective-<YYYY-MM-DD>.md
```

One file per analysis run (not a fixed cadence — run it whenever you want to check whether an agent is
getting slower, more retried, or trending in the right direction). Each file analyzes the last N
`docs/features/*/pipeline-trace.json` files (default 10) and reports, per agent: average duration, average
iterations, and whether each is `IMPROVING`/`STABLE`/`DEGRADING` — plus whether that trend's boundary lines
up with an `agentVersion` change, which is the difference between "this prompt edit caused it" and "the mix
of features analyzed just got harder."

## This is distinct from
- `docs/features/<name>/retrospective.md` — one delivery's qualitative narrative.
- `docs/agent-metrics/` — cross-delivery *quality* scores (was the output good), not timing.
- `docs/lessons-learned/` — recurring patterns worth promoting to a rule, prompt change, or KI.

`pipeline-retrospective` is the only one of the four that's purely about timing and retry counts — it never
judges whether an agent's output was actually correct.
