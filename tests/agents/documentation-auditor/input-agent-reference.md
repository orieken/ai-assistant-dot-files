# AGENT_REFERENCE.md (excerpt — under audit)

## Available Agents

| Agent | Description | Invocation |
|---|---|---|
| analyst | Reads a feature file and produces analysis.md | First step of deliver-feature |
| developer | Implements the feature | After analyst |
| code-reviewer | Reviews implementation | After developer |
| **legacy-formatter** | Formats code output (deprecated) | Manual only |
| qa-engineer | Writes and runs tests | After code-reviewer |
| tech-writer | Updates documentation | After qa-engineer |
| devops-engineer | Handles CI/CD updates | After tech-writer |

## Skills

- `/deliver-feature` — Full pipeline from analyst to devops-engineer
- `/create-ki` — Creates a Knowledge Item
- `/epic-planner` — Plans an epic (DEPRECATED — replaced by /spec-writer in v1.9)

## Recent Changes
- Added `sre-engineer` agent (v1.8) — not listed in table above
- Added `visual-qa-engineer` agent (v1.7) — not listed in table above
- Removed `legacy-formatter` (v1.6) — still listed in table above
