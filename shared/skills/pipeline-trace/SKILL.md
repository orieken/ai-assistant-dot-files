---
name: pipeline-trace
description: Owns the pipeline-trace.json schema (per-agent duration, status, iteration count) that deliver-feature writes during a run, and answers ad-hoc questions about a specific run's trace ("how long did the code-reviewer loop take on user-auth?", "show me the trace for the current run so far").
triggers:
  keywords: ["pipeline-trace", "trace", "how long did", "show the trace"]
  intentPatterns: ["Show the pipeline trace for *", "How long did * take on *", "/pipeline-trace *"]
standalone: true
---

## When To Use
- `deliver-feature` writes to `.claude/feature-workspace/pipeline-trace.json` directly as part of its own
  Checkpoint bookkeeping (see `deliver-feature/SKILL.md`, "Pipeline Tracing") — it does not invoke this
  skill on every step, the same way it doesn't invoke a separate skill to write `pipeline-state.json`.
- Invoke this skill when a human asks an ad-hoc question about timing/iterations for a specific feature's
  run, current or past: "why did this delivery take so long", "which agent looped the most on user-auth".
- For cross-delivery trend analysis (is code-reviewer getting slower over the last 10 features?), use
  `pipeline-retrospective` instead — this skill only looks at one run at a time.

## Context To Load First
1. `.claude/feature-workspace/pipeline-trace.json` (in-progress run) or `docs/features/<feature-name>/pipeline-trace.json` (completed run)

## Schema
```json
{
  "featureName": "user-auth",
  "startedAt": "2026-07-02T10:00:00Z",
  "completedAt": "2026-07-02T11:30:00Z",
  "totalDurationSeconds": 5400,
  "agents": [
    {
      "agent": "analyst",
      "agentVersion": "1.0.0",
      "step": 8,
      "startedAt": "2026-07-02T10:05:00Z",
      "completedAt": "2026-07-02T10:15:00Z",
      "durationSeconds": 600,
      "status": "PASS",
      "iterations": 1,
      "contractRetries": 0,
      "budgetUtilization": 0.42
    },
    {
      "agent": "code-reviewer",
      "agentVersion": "1.0.0",
      "step": 19,
      "startedAt": "2026-07-02T10:40:00Z",
      "completedAt": "2026-07-02T11:10:00Z",
      "durationSeconds": 1800,
      "status": "APPROVED",
      "iterations": 3,
      "changesRequestedCount": 2,
      "budgetUtilization": null
    }
  ]
}
```

Field notes:
- `budgetUtilization`: copied from `context-manifest.md`'s Token Budget section for this agent's tier at the
  time it ran — estimated tokens for pinned files ÷ that tier's ceiling (Analyst/Architect 60%, Developer
  80%, Reviewer 40%, of a 200k-token window), expressed as a fraction of the *tier ceiling itself* (so `1.0`
  means exactly at the tier limit, not at 100% of the full window). This is context-engineer's upfront
  estimate, not an independently re-measured value — there's no live hook into actual model context usage
  at runtime. Use `null` for agents that don't consume `context-manifest.md` directly (e.g. code-reviewer,
  security-reviewer, qa-engineer read specific prior artifacts, not the manifest's Pinpoint Files) rather
  than fabricating a number.
- `agentVersion`: that agent's `version:` frontmatter field (`shared/agents/CHANGELOG.md` tracks the
  history) at the time it ran. This is what lets `agent-scorecard` and `pipeline-retrospective` correlate a
  duration or iteration-count trend with a specific prompt edit instead of just observing drift with no
  attributable cause.
- `iterations`: total number of times this agent's step was executed, including retries (structural
  validate-artifact retries and, for code-reviewer, CHANGES REQUESTED loops).
- `contractRetries`: present only for contract-bound agents (see `shared/contracts/`) — how many times
  `validate-artifact` failed before passing.
- `changesRequestedCount`: present only for `code-reviewer` — how many times the verdict was CHANGES
  REQUESTED before APPROVED.
- `durationSeconds` accumulates across retries — it's total wall-clock time spent on that agent's step
  across every attempt, not just the final successful one.

## Process
1. Read the relevant `pipeline-trace.json` (in-progress in `.claude/feature-workspace/`, or persisted in
   `docs/features/<feature-name>/` for a completed run).
2. Answer the specific question asked — total duration, slowest agent, which agent looped, etc. — by
   reading the `agents` array directly. No aggregation across multiple files here (that's
   `pipeline-retrospective`).

## Output Format
Free-form answer to the question asked, grounded in the actual JSON values — always cite the specific
`durationSeconds`/`iterations` numbers rather than vague characterizations ("code-reviewer looped 3 times,
totaling 1800s" not "code-reviewer took a while").

## Guardrails
- Never fabricate or estimate a timestamp/duration not present in the file — if the trace is missing or
  incomplete (e.g. an interrupted run), say so explicitly rather than guessing.
- This skill is read-only with respect to `pipeline-trace.json` — it does not write to it. Only
  `deliver-feature` writes trace entries, as part of its own Checkpoint steps.

## Standalone Mode
Pure local JSON file reads. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
