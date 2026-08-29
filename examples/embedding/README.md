# Embedding example

A minimal MCP server that embeds loom's six framework tools **plus one custom
tool** (`echo`) through the public embedding API (roadmap D.2):

- [`echotool/echo_tool.go`](echotool/echo_tool.go) — the custom capability.
  Imports only the stdlib and `github.com/orieken/loom/tools`: no loom
  internals, no MCP library.
- [`main.go`](main.go) — merges `register.Frameworks(nil)` into its own
  `tools.Registry`, then adapts the registry to **this module's own** `mcp-go`
  dependency via the ~30-line `wire*` helpers.

This module has its own `go.mod` (with a `replace` pointing at the repo root),
so it builds exactly the way an external consumer would:

```bash
cd examples/embedding
go build ./...
```

See the "Embedding loom's tools" section of
[`shared/mcp/README.md`](../../shared/mcp/README.md) for the full walkthrough
and the semver compatibility promise.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
