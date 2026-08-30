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
| `loom uninstall` | Remove installed framework content |
| `loom update` | Update installed framework content |

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

`loom mcp serve` exposes the six deterministic framework-analysis tools
(`analyze_complexity`, `check_accessibility`, `check_ubiquitous_language`,
`verify_dependencies`, `search_ki`, `search_docs`) over MCP stdio transport.
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
