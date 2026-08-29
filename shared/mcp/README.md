# ai-assistant-dotfiles MCP Server

A reference scaffold exposing the six M1 framework-analysis tools as an
[MCP](https://modelcontextprotocol.io/) server over stdio transport.

## Tools

| Tool | Description |
|---|---|
| `analyze_complexity` | Cyclomatic complexity + LOC per function against framework thresholds (< 7 / < 30) |
| `check_accessibility` | Semantic-HTML / ARIA violations in HTML, Vue, JSX, TSX, Svelte files |
| `check_ubiquitous_language` | Synonym violations against a `DOMAIN_DICTIONARY.md` |
| `verify_dependencies` | Clean Architecture layer-boundary violations in Go and TypeScript imports |
| `search_ki` | Lexical-ranked search of the framework's Knowledge Items and ADRs |
| `search_docs` | BM25 search (sqlite-fts5) of the installed project's `docs/` corpus |

All six tools are deterministic and stateless — no LLM is required.

## Quick start

The supported way to run this server is the `loom` binary itself:

```bash
brew install orieken/tap/loom   # or: go install github.com/orieken/loom/cmd/loom@latest
loom mcp serve                  # stdio transport; logs to stderr or --log-file
```

> **Deprecated**: the standalone `cmd/mcp-server` entrypoint is retained for
> one release cycle only. It still builds from the repo root
> (`go build ./shared/mcp/cmd/mcp-server`), but new setups should use
> `loom mcp serve`. Since the module merge into `github.com/orieken/loom`,
> this directory is part of the root Go module and no longer carries its own
> `go.mod`.

## Configuration

Copy `.env.example` to `.env` and set:

| Variable | Purpose | Default |
|---|---|---|
| `AI_ASSISTANT_DOTFILES_PATH` | Root of this framework checkout; enables `search_ki` corpus | _(required for search_ki)_ |
| `DOCS_FTS_PATH` | Absolute path to the FTS5 sqlite index for `search_docs` | `.claude/rag/docs-fts5.sqlite` if `.claude/` exists |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OpenTelemetry collector endpoint | _(optional)_ |

## Wire into Claude Code

Add to `.claude/mcp.json` (or `~/.claude/mcp.json` for global):

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

Or register it with the Claude Code CLI:

```bash
claude mcp add loom -- loom mcp serve
```

## Installing into a downstream project

If you already have an MCP server, see
[`shared/skills/install-framework-with-mcp-bridge/SKILL.md`](../skills/install-framework-with-mcp-bridge/SKILL.md)
for the bridge-prompt approach (no Go required — drop a
`shared/mcp-patterns/go/` tool call into your existing server).

If you do not have an MCP server, you don't need a scaffold anymore — install
the `loom` binary and point your MCP host at `loom mcp serve`. The
`--with-mcp` scaffold copy (`install.sh` / `loom install`) still works but now
produces reference source only: since the module merge it has no `go.mod` of
its own and is not standalone-buildable. It is deprecated alongside
`cmd/mcp-server`.

## Layout

```
shared/mcp/
├── cmd/mcp-server/main.go       # standalone stdio entrypoint (deprecated — use `loom mcp serve`)
├── register/register.go         # FrameworkTools — embedding entry point
├── internal/
│   ├── analyzers/               # complexity, accessibility, language, deps analyzers
│   ├── domain/tool.go           # transport-free Tool interface (stdlib-only, enforced by test)
│   ├── logging/logger.go        # slog-backed Logger
│   ├── server/
│   │   ├── handler.go           # Handler struct
│   │   ├── mcp_adapter.go       # domain ↔ mcp-go wire-type conversion (only place mcp types touch tools)
│   │   ├── registration.go      # AddTool loop
│   │   └── tool_provider.go     # buildFrameworkTools wiring
│   └── tools/                   # 6 MCP tool implementations + retriever
└── .env.example
```

These packages live in the root Go module (`github.com/orieken/loom`) under
`shared/mcp/`; lint config (`.golangci.yml`, gocyclo cap at 7) sits at the
repo root.

## Dependencies

- [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) v0.57.0 — MCP SDK
- [`invopop/jsonschema`](https://github.com/invopop/jsonschema) v0.14.0 — output schema reflection
- [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) v1.55.0 — pure-Go sqlite for BM25 retrieval
