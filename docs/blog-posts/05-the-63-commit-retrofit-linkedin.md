Most Model Context Protocol (MCP) servers start as simple inline handlers.

Then you add 15 tools, direct sub-process calls, and un-schema'd JSON maps—and suddenly your server handler is an 880-line God Object.

While refactoring `saturday-mcp` (the Go MCP server powering our multi-agent framework), we executed a 63-commit retrofit to bring Clean Architecture to our MCP tool layer.

The highlights:

- **880 LOC → 93 LOC** in `handler.go` (an 89% reduction in core server code)
- **86.7% unit test coverage** across our newly extracted `internal/tools/` package
- **100% test coverage** on filesystem adapters and 97.2% on sub-process test runners
- **Zero SDK imports** inside our core domain interfaces (`Tool`, `Persona`, `Workflow`)

The key takeaway: MCP tools are just entrypoint adapters. When you decouple business logic from MCP SDK request/response shapes, adding a new tool becomes as simple as dropping a single Go file into a directory.

Full technical breakdown and commit-by-commit refactoring pattern: TODO_DEVTO_URL
