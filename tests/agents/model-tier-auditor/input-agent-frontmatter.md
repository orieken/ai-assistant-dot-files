# Agent Frontmatter Samples (under audit)

## analyst.md
```yaml
name: analyst
model_tier: reasoning
description: Reads a feature file and produces analysis.md
```

## developer.md
```yaml
name: developer
model_tier: standard
description: Implements the feature
```

## context-engineer.md
```yaml
name: context-engineer
description: Pre-flight context optimizer
# model_tier: missing entirely
```

## agent-evaluator.md
```yaml
name: agent-evaluator
model_tier: budget
description: Read-only; runs golden-file evaluations
```

## sre-engineer.md
```yaml
name: sre-engineer
model_tier: observability
description: Reviews for OTel and SLIs
# "observability" is not a valid model_tier enum value
```

## chaos-engineer.md
```yaml
name: chaos-engineer
model_tier: standard
description: Designs fault-injection experiments
# chaos-engineer is write-heavy and autonomously executes bash — "standard" may be too low for its operational profile
```
