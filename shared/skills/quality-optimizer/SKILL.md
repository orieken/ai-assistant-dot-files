---
name: quality-optimizer
description: Identifies pipeline stages with contract retries or quality degradation and recommends higher-fidelity models or additional review steps. Opposing-force pair with cost-optimizer.
triggers:
  keywords: [quality optimizer, optimize quality, elevate model, quality fidelity]
  intentPatterns: ["recommend quality optimizations", "run quality optimizer"]
standalone: true
---

## When To Use

Use when:
- Analyzing pipeline retries, code review rejections, or contract failures.
- Recommending higher-fidelity model tiers for complex reasoning tasks.

Do NOT use when:
- Reducing model cost (use `cost-optimizer`).

## Context To Load First

1. `.claude/telemetry/`
2. `shared/evaluation/`

## Process

1. **Read Pipeline Retries**: Glob `.claude/telemetry/*.jsonl`.
2. **Identify Quality Bottlenecks**: Find stages with high retry counts ($N \ge 2$) or code review rejections.
3. **Draft Quality Recommendation**: Produce `.claude/feature-workspace/quality-report.md`.
4. **Pause for Human Confirmation**: Present model upgrade recommendations to human operator.

## Output Format

```markdown
# Quality Optimization Recommendations

## Proposed Model Upgrades
- Stage `architect`: High contract retry rate ($N=2$). Recommend upgrading model to `opus`.
```

## Guardrails

- Never alter agent configuration without human confirmation ("approve quality upgrades" or "confirm").

## Standalone Mode

Operates purely using local file editing tools.
