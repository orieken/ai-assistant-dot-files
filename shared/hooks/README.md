# AOS Hooks Layer (`shared/hooks/`)

The AOS Hooks layer provides an event-driven interceptor system for the framework. Hooks allow skills, counter-agents, and telemetry recorders to react to pipeline lifecycle events automatically.

> **Opt-In Guarantee**: Hooks are strictly opt-in. If no hook configuration file exists in `.claude/hooks/` or project config, the framework operates in standard synchronous mode without executing hook listeners.

---

## Event Catalog

Hooks trigger on specific pipeline events. These trigger names are the hook layer's own vocabulary and are not the executor's event types — for those, see the generated `shared/schemas/telemetry/run-event-types.md`. Nothing dispatches hooks today; the executor for them is roadmap L3.10.

| Event Name | Trigger Condition | Common Action |
|---|---|---|
| `on-artifact-write` | Any agent writes or updates a workspace artifact | Record telemetry event / audit PII |
| `on-validation-pass` | `validate-artifact` returns `PASS` for a stage | Invoke corresponding counter auditor |
| `on-validation-fail` | `validate-artifact` returns `FAIL` for a stage | Log contract retry event |
| `on-ki-created` | `create-ki` or `memory-expansion` writes a new KI | Invoke `knowledge-auditor` |
| `on-pipeline-complete` | Feature delivery pipeline finishes Phase 4 | Trigger retrospective / scorecard sweep |
| `on-inventory-change` | A file under `shared/agents/` or `shared/skills/` is added, modified, or removed | Invoke `documentation-auditor` to check for prose staleness |

---

## Index-Freshness Story (ADR-002 rebuild hook — pending saturday-mcp M2)

[ADR-002](../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md) names index staleness as a
known consequence and proposes an event-driven rebuild story: "cache invalidation on doc edits is a
real concern for the vector tier" — the hooks layer is the correct place to wire this without
requiring manual `/reindex` runs after every delivery.

The planned hook pair (awaiting saturday-mcp M2 delivery of a `reindex` MCP tool):

| Hook file | Event | Action |
|---|---|---|
| `examples/on-artifact-written-reindex.yaml` | `on-artifact-write` (filtering to `docs/features/`, `docs/adrs/`) | `skill: reindex` — upsert the written path into the BM25/vector index |
| `examples/on-ki-created-reindex.yaml` | `on-ki-created` (targeting `shared/knowledge/`, `.claude/knowledge/`) | `skill: reindex` — upsert the new KI into the index |

**Why deferred**: saturday-mcp M1 ships `search_docs` and `search_ki` but no `reindex` entry point.
Creating hook example files targeting a non-existent tool would produce silent failures when enabled.
These examples ship alongside the M2 `reindex` tool, not before it.

**Manual fallback until M2**: `/reindex` (run after any bulk write to docs/ or shared/knowledge/).
This is the "install-time build + rebuild story" already documented in `docs/runbooks/lightrag-integration.md`.

**Proposed M2 tool interface** (for saturday-mcp implementors):
```yaml
# shared/hooks/examples/on-artifact-written-reindex.yaml (M2 placeholder — not yet active)
version: "1.0"
hooks:
  - id: "bm25-reindex-on-artifact-write"
    event: "on-artifact-write"
    enabled: false   # Enable after saturday-mcp M2 ships reindex tool
    filter:
      pathPrefixes: ["docs/features/", "docs/adrs/", "docs/patterns/"]
    action:
      type: "skill"
      target: "reindex"
      args:
        scope: "docs"
```

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
