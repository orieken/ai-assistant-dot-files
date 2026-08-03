# Epic 62 — Gate-Decision + Cost Telemetry

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #2 — cheap schema additions,
compounding returns). The gap: human corrections at approval gates evaporate, and nobody records
what deliveries actually cost.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

- `shared/telemetry/event-schema.md` + `event-recorder.md` (AOS Phase 1) — opt-in event log at
  `.claude/telemetry/events.jsonl`; the v3.0.0 guarantee is that NOTHING emits telemetry by
  default. Verified: the schema has no gate-decision events and no token/cost fields.
- `shared/rules/approval-gates.md` — 8 gates; `deliver-feature` + `ship-feature` (Epic 47) halt
  at them. When a human rejects or edits an artifact before approving, that signal — the richest
  feedback the system receives — is discarded. `extract-lessons` mines what agents wrote, never
  what humans fixed.
- `pipeline-trace.json` (see `shared/skills/pipeline-trace/SKILL.md`) — per-agent duration,
  status, iterations, and *estimated* `budgetUtilization`. No actual token spend, so the
  `cost-optimizer`/`quality-optimizer` opposing pair and `finops-engineer` operate on no data.

## Scope: 4 ops

**Op 1 — Gate-decision events in the schema.**
Extend `shared/telemetry/event-schema.md` with a `gate_decision` event: gate id (1–8 from
approval-gates.md), artifact path, decision (`approved` | `rejected` | `edited_then_approved`),
optional free-text reason, timestamp, feature name. Follow the schema doc's existing structure
and its secrets/PII prohibition verbatim. Bump the schema's version marker if it has one; add one
if it doesn't (note it — downstream aggregation needs it).
Commit: `feat(telemetry): gate_decision event schema (Epic 62 Op 1)`

**Op 2 — Emit gate decisions from deliver-feature.**
Update `deliver-feature`'s gate-halt instructions: WHEN telemetry is enabled (and only then —
the opt-in guarantee is non-negotiable), record a `gate_decision` event after each human
response. An "edited" detection heuristic: artifact mtime/checksum changed between present and
approve (checksums already exist in `pipeline-state.json`).
Commit: `feat(deliver-feature): emit gate_decision telemetry, opt-in (Epic 62 Op 2)`

**Op 3 — Actual token spend in pipeline-trace.**
Extend the `pipeline-trace.json` schema (owned by `shared/skills/pipeline-trace/SKILL.md`) with
per-agent `tokensIn`/`tokensOut`/`estimatedCost` fields, populated when the runtime exposes usage
(document honestly: Claude Code sessions don't always surface per-subagent token counts — record
what's available, `null` otherwise, never fabricate). Update `pipeline-retrospective` and
`agent-scorecard` skills to surface cost columns when present.
Commit: `feat(pipeline-trace): token spend fields (Epic 62 Op 3)`

**Op 4 — Close the loop to the miners.**
Update `extract-lessons` + `retrospective` skills: when a telemetry log with `gate_decision`
events exists, rejected/edited gates are first-class input ("what did humans correct, and does a
pattern exist?"). Update `shared/evaluation/README.md` to name gate-rejection rate as a
per-agent quality signal `agent-scorecard` may aggregate.
Commit: `feat(skills): mine gate decisions in extract-lessons + retrospective (Epic 62 Op 4)`

After every op: `bash scripts/health-check.sh` green; verify on a pristine repo that NO event is
emitted without telemetry enabled (the Phase 1 guarantee, re-proven).

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If the event schema has no versioning story and adding one ripples into event-recorder
  consumers, halt on Op 1 with a proposal.
- If per-agent token usage is simply unobtainable in this runtime, ship the schema fields anyway
  (`null`-populated), and say so in the report — schema first, data when available.
- If emitting from gate halts would require touching the approval-gates RULES file, halt — gates
  themselves must not change in this epic.

## Report (under 120 words)

```
Commits: <sha> x4
Schema version: <old → new | added at vX>
Opt-in guarantee re-proven: <how>
Token fields populated in practice: <yes/partial/null-only + why>
Gate-rejection mining: <which skills updated>
```

Go.
