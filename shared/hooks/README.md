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
| `on-inventory-change` | A file under `shared/agents/` or `shared/skills/` is added, modified, or removed | Invoke `documentation-auditor` to check for prose staleness |

---

---

## Security Constraints

### `action.type: "script"` Is Privileged — Requires Explicit Review

Hook actions of type `"agent"` and `"skill"` run within the AI tool's normal permission model.
`action.type: "script"` is fundamentally different: it executes an **arbitrary shell script** in
the AI tool's process environment with no framework-level sandboxing. Any tool or file the AI
tool can touch, the script can touch.

**Before enabling any hook with `action.type: "script"`**:

1. The script path must be under the project's version-controlled tree (no `/tmp/`, no `~/.local/`).
2. The script must be code-reviewed by a human with the same scrutiny as any other shell script
   that runs in CI — not treated as a "just a hook" afterthought.
3. If `passContext: true` is set, the script receives pipeline context (feature name, workspace
   paths, potentially artifact content). Treat this as granting the script read access to all
   pipeline artifacts. Scope it explicitly if the script only needs a subset.
4. Never copy a `type: "script"` hook file from an untrusted source without reviewing the
   `target` path and the script itself.

**Fitness function**: `scripts/health-check.sh` warns when any enabled hook in `.claude/hooks/`
uses `type: "script"`. The warning is informational — it means the hook was reviewed and
intentionally enabled — not a blocker, but it requires the reviewer to confirm it explicitly.

### What Hooks Must Never Do

- Write outside `.claude/feature-workspace/` or `docs/features/` without a separate Gate #6
  (out-of-boundary write) approval.
- Exfiltrate artifact content, credentials, or file paths to external services.
- Modify `shared/rules/`, agent prompts, or `shared/schemas/` — those changes require a human
  commit gate, not an automated hook.
- Chain into other hooks in a way that creates an unbounded execution loop.

### Adding a New Hook to the Examples Directory

Before adding a new `.yaml` to `shared/hooks/examples/`:

1. Set `enabled: false` (default for examples — a developer who opts in sets it to `true`
   explicitly in their `.claude/hooks/` copy).
2. If `action.type` is `"script"`, the PR description must include a security rationale.
3. The hook must appear in the Event Catalog above with its trigger condition documented.

---

## Directory Structure

- `shared/hooks/README.md`: Overview, event catalog, and security constraints (this file).
- `shared/hooks/hooks-schema.md`: YAML/JSON schema for valid hook definitions.
- `shared/hooks/examples/`: Reference hook configurations for project opt-in (set `enabled: false` by default).
