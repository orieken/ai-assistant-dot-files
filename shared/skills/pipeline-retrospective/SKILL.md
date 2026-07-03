---
name: pipeline-retrospective
description: Analyzes pipeline-trace.json across the last N feature deliveries in docs/features/ to find cross-delivery trends — which agent is the slowest or most-retried, whether trends are improving or degrading, and persistent bottlenecks. Complements the single-delivery retrospective skill, which only looks at one feature at a time.
triggers:
  keywords: ["pipeline-retrospective", "trend analysis", "cross-delivery", "agent trends"]
  intentPatterns: ["Run a pipeline retrospective", "How are our agents trending?", "Which agent is the bottleneck?", "/pipeline-retrospective *"]
standalone: true
---

## When To Use
When the user wants to know how the pipeline is performing *across* deliveries, not within one — "is
code-reviewer getting slower", "which agent causes the most rework", "are we improving". Manually invoked;
not on a fixed cadence (unlike the single-delivery `/retrospective`, which `deliver-feature` auto-invokes
every 5th delivery).

Do NOT use for a single delivery's narrative (what went well/poorly on this one feature) — use
`/retrospective` instead. Do NOT use for a single run's raw numbers — use `pipeline-trace` instead.

## Context To Load First
1. `docs/features/*/pipeline-trace.json` — collect the most recent N (default 10), ordered by each
   feature directory's `pipeline-trace.json` `completedAt` timestamp
2. The most recent `docs/agent-metrics/scorecard-*.md` if one exists (see `agent-scorecard`), for
   qualitative cross-reference alongside the timing/iteration trend

## Process
1. Collect the last N `pipeline-trace.json` files (default N=10; use fewer if fewer exist, but report
   "insufficient data" instead of a trend claim if fewer than 3 exist).
2. For each agent that appears across the collected traces, compute:
   - Average `durationSeconds` and average `iterations`.
   - Trend: split the N traces into an older half and a newer half (by `completedAt`); compare each
     metric's average between halves. Report `IMPROVING` / `STABLE` / `DEGRADING` (>15% change = not stable).
   - If `agentVersion` changes partway through the collected traces, check whether the trend boundary lines
     up with that version change — a DEGRADING (or IMPROVING) trend that starts exactly at a version bump
     is strong evidence the prompt edit caused it, not noise. A trend with no corresponding version change
     is weaker evidence (could be caused by the mix of features analyzed getting harder/easier).
3. Identify the single biggest bottleneck: the agent with the highest average `durationSeconds` across all
   traces.
4. Identify the single most-retried agent: the agent with the highest average `iterations` (this usually
   means either `code-reviewer` (CHANGES REQUESTED loops) or a contract-bound agent failing
   `validate-artifact` repeatedly).
5. If a recent `agent-scorecard` exists, cross-reference: does a DEGRADING trend line up with a low score
   for that same agent? Call this out explicitly — it's the strongest evidence of a real regression versus
   noise.
6. Write the report to `docs/pipeline-retrospectives/retrospective-[YYYY-MM-DD].md`.

## Output Format
```markdown
# Pipeline Retrospective: [YYYY-MM-DD]

## Sample
- Deliveries analyzed: [N] (out of [total] available)
- Date range: [oldest completedAt] to [newest completedAt]

## Per-Agent Trends
| Agent | Avg Duration | Avg Iterations | Duration Trend | Iteration Trend | Version Change Aligns? |
|---|---|---|---|---|---|
| analyst | [Ns] | [N] | IMPROVING/STABLE/DEGRADING | IMPROVING/STABLE/DEGRADING | Yes (v[x.y.z] -> v[x.y.z]) / No version change in window |
| ... | ... | ... | ... | ... | ... |

## Biggest Bottleneck
[agent] — average [Ns] per run, [context on why, e.g. "consistently the longest single step"]

## Most-Retried Agent
[agent] — average [N] iterations per run, [context, e.g. "usually from code-reviewer CHANGES REQUESTED, not contract failures"]

## Cross-Reference with Agent Scorecard
[If a recent scorecard exists: does the trend match the score? If no scorecard exists: "No scorecard available — see agent-scorecard skill."]

## Recommendations
- [Specific, actionable: which agent prompt to review, which contract may be too strict, etc.]
```

## Guardrails
- **Never** report a trend from fewer than 3 traces — say "insufficient data for a trend, only N deliveries
  recorded" instead of guessing.
- **Never** blame an agent's prompt without also checking whether the underlying features analyzed were
  just harder/riskier in that period (e.g. more features touching auth/payments would legitimately raise
  security-reviewer's average duration without the prompt regressing).
- This is a read-only analysis — it does not modify any agent prompt, contract, or trace file.

## Standalone Mode
Pure local file reads and aggregation. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
