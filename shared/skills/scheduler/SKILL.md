---
name: scheduler
description: Orchestrates scheduled or hook-driven pipeline runs (cron triggers, automated memory audits, or periodic health checks).
triggers:
  keywords: [scheduler, cron trigger, run scheduled task, periodic audit]
  intentPatterns: ["schedule pipeline run", "execute cron task"]
standalone: true
---

## When To Use

Use when:
- Executing recurring framework checks (e.g., weekly `memory-auditor` or `pipeline-retrospective` runs).
- Configuring background event schedules using the system `schedule` tool.

Do NOT use when:
- Running interactive single-feature deliveries (use `deliver-feature`).

## Context To Load First

1. `shared/hooks/README.md`
2. `shared/telemetry/`

## Process

1. **Read Schedule Specification**: Identify target event, cadence, and action script/agent.
2. **Configure Background Timer**: Invoke system `schedule` tool with explicit prompt or cron expression.
3. **Log Scheduled Event**: Append schedule event log to `.claude/telemetry/events.jsonl`.
4. **Report Schedule Configuration**: Output status confirmation report.

## Output Format

```markdown
# Scheduler Status Report

## Active Timers / Cron Schedules
- Event: `weekly-memory-audit` | Cadence: `0 0 * * 1` | Action: `memory-auditor`
```

## Guardrails

- Never launch mutating background commands without explicit human approval.

## Standalone Mode

Delegates timer management to the `schedule` tool or outputs standard shell instructions.
