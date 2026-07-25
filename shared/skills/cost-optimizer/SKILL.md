---
name: cost-optimizer
description: Recommends model and agent cost optimizations when quality evaluation metrics permit. Opposing-force pair with quality-optimizer.
triggers:
  keywords: [cost optimizer, optimize cost, reduce tokens, model selection]
  intentPatterns: ["recommend cost optimizations", "run cost optimizer"]
standalone: true
---

## When To Use

Use when:
- Analyzing token consumption and model selection across agent runs.
- Recommending cheaper model mappings for simple pipeline tasks.

Do NOT use when:
- Elevating model fidelity for complex tasks (use `quality-optimizer`).

## Context To Load First

1. `.claude/telemetry/`
2. `shared/agents/`

## Process

1. **Read Pipeline Telemetry**: Glob `.claude/telemetry/*.jsonl`.
2. **Evaluate Token & Model Cost**: Find stages where fast/cheaper models achieve 100% contract pass rates.
3. **Draft Cost Optimization Report**: Produce `.claude/feature-workspace/cost-report.md`.
4. **Pause for Human Confirmation**: Present recommendations to human operator.

## Output Format

```markdown
# Cost Optimization Recommendations

## Proposed Model Downgrades
- Stage `tech-writer`: Recommend switching from `opus` to `haiku/sonnet` (Contract Pass Rate: 100%).
```

## Guardrails

- Never alter agent model mappings without human confirmation ("approve cost optimizations" or "confirm").

## Standalone Mode

Operates purely using local file editing tools.
