# loom CLI

The `loom` binary installs the framework's agents, skills, and rules across AI
coding platforms and serves the framework's MCP tools. Install it via
`brew install orieken/tap/loom` or
`go install github.com/orieken/loom/cmd/loom@latest`.

## Commands

| Command | Purpose |
|---|---|
| `loom install` (alias `loom init`) | Install framework content for detected platforms (`--target`, `--platform`, `--copy`, `--dry-run`, `--stack`, `--level` — see the maturity-level table in the root README; profiles are defined in `shared/levels.yaml`) |
| `loom health` | Verify installed configs match the canonical `shared/` source, and report the project's agentic maturity level (see below) |
| `loom tools status` / `loom tools install` | Report / install opt-in context tools |
| `loom mcp serve` | Serve the framework MCP tools over stdio |
| `loom run` | Execute the delivery pipeline for a feature spec (experimental — see below) |
| `loom uninstall` | Remove installed framework content |
| `loom update` | Update installed framework content |

## Running pipelines (experimental)

`loom run` is the executor skeleton decided by ADR-006 (loom executes
pipelines) and built by roadmap item M0.4. It runs the built-in linear
`deliver-feature` plan — the same 14-agent sequence as the markdown
pipeline — stage by stage, persisting durable state after every transition.

```bash
# Execute the pipeline for a spec (spawns `claude -p` per stage)
loom run --spec docs/features/user-auth/spec.md

# Continue an interrupted run from its checkpoint
loom run --spec docs/features/user-auth/spec.md --resume

# Deterministic dry run with canned artifacts (no LLM calls)
loom run --spec docs/features/user-auth/spec.md --provider mock
```

- **State**: `.claude/feature-workspace/<feature>/run-state.json` — schema-versioned,
  written atomically (temp file + rename), with per-stage status, timestamps, and
  the SHA-256 of each stage's artifact. A fresh run refuses to start over existing
  state (pass `--resume` or delete the file); `--resume` requires existing state.
- **Interruption**: the first Ctrl-C cancels the in-flight stage and persists a
  clean `INTERRUPTED` checkpoint; a second Ctrl-C kills immediately. `--resume`
  skips completed stages and re-runs the interrupted one.
- **Providers**: `claude` (default) spawns the `claude` CLI headless per stage,
  building the prompt from the agent's `shared/agents/<agent>.md` definition; if
  the binary is missing the stage fails with a remediation message — there is no
  silent fallback. `mock` is for tests and dry runs.

- **Approval gates**: a gated stage does not start until a human approves its gate
  (roadmap L2.13). The built-in plan gates `developer` behind `confirm-design`,
  `qa-engineer` behind `confirm-security`, and `devops-engineer` behind
  `confirm-ship` — the same stops `deliver-feature` asks a human for. On a
  terminal you are asked `approve gate "X" for stage "Y"? [y/N]` at the barrier;
  otherwise the run persists `WAITING_APPROVAL`, prints the resume command, and
  exits with code **3** so scripts can tell "waiting on a human" from a failure:

  ```bash
  loom run --spec docs/features/user-auth/spec.md --resume --approve confirm-design
  ```

  `--approve` is only valid with `--resume`, and only for the gate the run is
  actually halted on. Nothing an agent returns can approve a gate.

**What it does NOT do yet** (skeleton by design — see `docs/roadmaps/BUILD-ROADMAP.md`):
no gate-reset enforcement — editing an artifact after approving does not reset
the gate (L2.14) — no retries or backoff, no parallelism (L3.3), no policy
evaluation (L2.16), no conditional stage routing (L3.1) — every stage of the
linear plan runs, including ones the markdown pipeline would skip — and no
telemetry emission (L3.8).

## Maturity level report

`loom health` ends with an agentic-maturity assessment driven by
`shared/levels.yaml`: it reports the highest level whose *entire* mechanical
evidence set passes, the passing evidence, and a concrete gap checklist for
the next level. Evidence is strictly mechanical — installed bundles on disk,
an MCP server that is configured in `.mcp.json` **and** actually answers
`tools/list` when spawned, a non-empty telemetry stream, policy files present.
Documentation-only bundles (`docsOnly` in `levels.yaml`) never count as
evidence, and a level whose enforcement bundles are all gated on unlanded
roadmap items is reported as not attainable yet rather than pretended into
reach.

```
Agentic maturity (shared/levels.yaml):
  Level 1 — Foundational prompts
  ✓ bundle "core-rules" installed
  ✓ bundle "agents" installed
  gaps to Level 2:
    ✗ MCP server not configured: read .mcp.json: no such file
    ✗ bundle "workflows" not fully installed (missing workflows)
```

## MCP server

`loom mcp serve` exposes the deterministic framework tools
(`analyze_complexity`, `check_accessibility`, `check_ubiquitous_language`,
`verify_dependencies`, `search_ki`, `search_docs`, `validate_artifact`)
over MCP stdio transport.
Structured JSON logs go to stderr, or to a file with `--log-file <path>` —
never stdout, which carries the MCP wire protocol. The server runs until the
client closes stdin or the process receives SIGINT.

Register it with Claude Code:

```bash
claude mcp add loom -- loom mcp serve
```

Or add it to `.mcp.json` (project) / `~/.claude/mcp.json` (global):

```json
{
  "mcpServers": {
    "loom": {
      "command": "loom",
      "args": ["mcp", "serve"],
      "env": {
        "AI_ASSISTANT_DOTFILES_PATH": "/absolute/path/to/loom-checkout"
      }
    }
  }
}
```

`AI_ASSISTANT_DOTFILES_PATH` is only required for `search_ki` (it points the
tool at the Knowledge Item corpus); the other five tools need no
configuration. See [shared/mcp/README.md](../../shared/mcp/README.md) for the
full tool reference and environment variables.
