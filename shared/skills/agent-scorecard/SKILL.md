---
name: agent-scorecard
description: Scores each pipeline agent against a defined quality metric using artifacts from recent deliveries in docs/features/, compares against the previous month's scorecard to flag improving/degrading agents, and persists the result to docs/agent-metrics/scorecard-YYYY-MM.md.
triggers:
  keywords: ["agent-scorecard", "score agents", "agent metrics", "agent performance"]
  intentPatterns: ["Score the agents", "How are the agents performing", "Generate this month's scorecard", "/agent-scorecard *"]
standalone: true
---

## When To Use
Monthly cadence (manually triggered, or wire up via the host tool's own scheduling capability — e.g. Claude
Code's `schedule` skill — if the user wants it automatic; this is a platform capability, not one of this
repo's `shared/skills/`, so it isn't auto-invoked by `deliver-feature` itself the way `/retrospective` is).
Also on-demand whenever the user asks how a specific agent, or the pipeline generally, is performing.

Do NOT use for a single delivery's story (what happened on this one feature) — use `/retrospective`. Do NOT
use for cross-delivery timing/iteration trends without a quality judgment attached — use
`pipeline-retrospective` (this skill is complementary to it, not a replacement: `pipeline-retrospective` says
*how long/how many retries*, this skill says *was the output actually good*).

## Metric Definitions
These are the only four scored metrics. Don't invent new ones without updating this list first — an
undocumented metric can't be tracked for a trend.

| Agent | Metric | Computed From | Underperforming Floor |
|---|---|---|---|
| `security-reviewer` | True positive rate (proxy) | `security-report.md`: fraction of CRITICAL/HIGH findings with a non-"Recommendation only" `Fix applied` line, adjusted down for any finding a later `retrospective.md` explicitly disputes as a false alarm | < 80% |
| `code-reviewer` | First-pass acceptance rate | `pipeline-trace.json`: fraction of deliveries where `changesRequestedCount == 0` | < 50% |
| `analyst` | Completeness score | `analysis.md` vs. `shared/contracts/analysis-contract.md`: fraction of required sections present **and** containing real content (not leftover `[...]` template placeholders) | < 90% |
| `architect` | Fitness function coverage | `architecture-notes.md`: fraction of `## Structural Decisions` entries with a concrete `**Fitness Function**` + `**Enforcement**` line, vs. entries explicitly flagged `judgment-only` | < 70% |

**Known limitation**: security-reviewer's "true positive rate" is a proxy, not a real confirmed/false-positive
rate — there's no mechanism yet for a human or downstream agent to formally dispute a finding after the fact.
Treat this metric as directional, not exact, until that dispute-tracking mechanism exists (tracked as an open
item — see `docs/features/context-engineering-framework/TODO.md`, Epic 15).

## Context To Load First
1. `docs/features/*/delivery-summary.md` — determine which features fall in the scoring period
2. `docs/features/*/pipeline-trace.json` — for code-reviewer's `changesRequestedCount`
3. `docs/features/*/security-report.md`, `analysis.md`, `architecture-notes.md` — per-agent artifacts
4. `docs/features/*/retrospective.md` (if present) — for disputed-finding signals
5. `shared/contracts/analysis-contract.md` — required section list for the analyst completeness score
6. The most recent `docs/agent-metrics/scorecard-*.md` (if one exists) — for trend comparison

## Process
1. **Determine scope**: all features in `docs/features/` whose `delivery-summary.md` falls in the current
   calendar month. If fewer than 3 features exist this month, widen to the last 3 delivered features instead
   (say so explicitly in the output — don't silently score on 1 data point).
2. **Compute each metric** per the table above, across every feature in scope.
3. **Compare to the previous scorecard** (`docs/agent-metrics/scorecard-<prior-YYYY-MM>.md`), if one exists.
   Trend = `IMPROVING` / `STABLE` / `DEGRADING` (>10 percentage-point swing = not stable; smaller threshold
   than `pipeline-retrospective`'s 15% because these are already rate/percentage metrics, not raw durations).
4. **Flag underperforming agents** — any metric below its floor in the table above.
5. **Write the scorecard** to `docs/agent-metrics/scorecard-[YYYY-MM].md`.

## Output Format
```markdown
# Agent Scorecard: [YYYY-MM]

## Scope
- Features scored: [N] — [list feature names]
- Period: [start date] to [end date]
- Note: [if scope was widened due to <3 features this month]

## Metrics

| Agent | Metric | Score | Prior Month | Trend | Status |
|---|---|---|---|---|---|
| security-reviewer | True positive rate (proxy) | [X%] | [Y%] / n/a | IMPROVING/STABLE/DEGRADING/n/a | OK/UNDERPERFORMING |
| code-reviewer | First-pass acceptance rate | [X%] | [Y%] / n/a | ... | ... |
| analyst | Completeness score | [X%] | [Y%] / n/a | ... | ... |
| architect | Fitness function coverage | [X%] | [Y%] / n/a | ... | ... |

## Underperforming Agents
- [agent]: [metric] at [X%], below the [floor]% floor. [Brief diagnosis grounded in the actual artifacts —
  e.g. "3 of 5 features this month had analysis.md missing a non-empty Edge Cases and Risks section"]
— or "None this period"

## Methodology & Known Limitations
- [Restate the security-reviewer proxy caveat if it's scored this period]
- [Any other scope caveats, e.g. widened window]

## Recommendations
- [Specific, tied to a specific metric/agent — e.g. "Review analyst.md's Output Format template; the Edge
  Cases section may need a stronger prompt nudge since it's the most frequently thin section"]
```

## Guardrails
- **Never** score an agent that wasn't invoked in any feature this period — report `n/a`, not `0%`
  (architect and security-reviewer are conditional; a period with no architecturally-flagged features
  should not show architect as "failing").
- **Never** fabricate a prior-month comparison if no prior scorecard exists — say "n/a, first scorecard".
- **Never** silently change a metric's definition or floor between runs — if you believe a metric needs to
  change, say so explicitly and note it breaks trend continuity with prior scorecards.
- This is a read-only analysis — it does not modify any agent prompt, contract, or artifact.

## Standalone Mode
Pure local file reads and aggregation. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
