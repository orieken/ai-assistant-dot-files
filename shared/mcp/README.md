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

```bash
# Build the server
go build -o mcp-server ./cmd/mcp-server

# Run — communicates over stdin/stdout per MCP stdio transport
./mcp-server
```

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
    "ai-assistant-dotfiles": {
      "command": "/absolute/path/to/mcp-server",
      "env": {
        "AI_ASSISTANT_DOTFILES_PATH": "/absolute/path/to/ai-assistant-dot-files"
      }
    }
  }
}
```

Or use the install helper (copies the scaffold into your project and runs `go mod tidy`):

```bash
./install.sh --project /path/to/my-project --with-mcp
```

## Installing into a downstream project

If you already have an MCP server, see
[`shared/skills/install-framework-with-mcp-bridge/SKILL.md`](../skills/install-framework-with-mcp-bridge/SKILL.md)
for the bridge-prompt approach (no Go required — drop a
`shared/mcp-patterns/go/` tool call into your existing server).

If you do not have an MCP server, `./install.sh --project <path> --with-mcp`
copies this scaffold to `<path>/<project-name>-mcp/`, runs `go mod tidy`,
and leaves a ready-to-build Go module.

## Layout

```
shared/mcp/
├── cmd/mcp-server/main.go       # stdio entrypoint
├── internal/
│   ├── analyzers/               # complexity, accessibility, language, deps analyzers
│   ├── domain/tool.go           # Tool interface
│   ├── logging/logger.go        # slog-backed Logger
│   ├── server/
│   │   ├── handler.go           # Handler struct
│   │   ├── registration.go      # AddTool loop
│   │   └── tool_provider.go     # buildFrameworkTools wiring
│   └── tools/                   # 6 MCP tool implementations + retriever
├── go.mod                       # module github.com/orieken/ai-assistant-dotfiles/mcp
├── .golangci.yml                # gocyclo cap at 7
└── .env.example
```

## Dependencies

- [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) v0.57.0 — MCP SDK
- [`invopop/jsonschema`](https://github.com/invopop/jsonschema) v0.14.0 — output schema reflection
- [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) v1.55.0 — pure-Go sqlite for BM25 retrieval
