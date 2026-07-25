# AOS Hooks Layer (`shared/hooks/`)

The AOS Hooks layer provides an event-driven interceptor system for the framework. Hooks allow skills, counter-agents, and telemetry recorders to react to pipeline lifecycle events automatically.

> **Opt-In Guarantee**: Hooks are strictly opt-in. If no hook configuration file exists in `.claude/hooks/` or project config, the framework operates in standard synchronous mode without executing hook listeners.

---

## Event Catalog

Hooks trigger on specific pipeline events (aligned with `shared/telemetry/event-schema.md`):

| Event Name | Trigger Condition | Common Action |
|---|---|---|
| `on-artifact-write` | Any agent writes or updates a workspace artifact | Record telemetry event / audit PII |
| `on-validation-pass` | `validate-artifact` returns `PASS` for a stage | Invoke corresponding counter auditor |
| `on-validation-fail` | `validate-artifact` returns `FAIL` for a stage | Log contract retry event |
| `on-ki-created` | `create-ki` or `memory-expansion` writes a new KI | Invoke `knowledge-auditor` |
| `on-pipeline-complete` | Feature delivery pipeline finishes Phase 4 | Trigger retrospective / scorecard sweep |

---

## Directory Structure

- `shared/hooks/README.md`: Overview and event catalog.
- `shared/hooks/hooks-schema.md`: YAML/JSON schema for valid hook definitions.
- `shared/hooks/examples/`: Reference hook configurations for project opt-in.
