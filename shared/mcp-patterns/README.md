# MCP Patterns (`shared/mcp-patterns/`)

Guidance for building or extending a Model Context Protocol server that speaks this framework's
tools, in a language other than Go.

## The Go "reference implementation" is gone, and that is the point

This directory used to carry `go/` — about 1,200 lines of Go copied out of `shared/mcp/internal/`,
tagged `//go:build ignore` so nothing compiled it, presented as templates to copy. Every one of the
eleven shared files had drifted from the original it was copied from (`retriever.go` by 186 lines,
`bm25_retriever.go` by 106), including a concurrency defect its own comment documented. Code that
cannot be built cannot be tested, and code that cannot be tested drifts silently and ships the drift
downstream.

It is deleted (roadmap **M0.5**). Copy-paste is the wrong distribution mechanism for a Go library
when a supported one exists.

## If you are writing Go: import, don't copy

`register.Frameworks(logWriter)` returns the framework's built-in tool registrations as a
transport-free `tools.Registry`. Merge it into your own registry and adapt it to whatever MCP
library and version your server uses:

```go
import "github.com/orieken/loom/shared/mcp/register"

registry := register.Frameworks(nil) // nil logWriter = os.Stderr
```

Nothing in the returned types references `mcp-go` or loom internals, so you are not pinned to loom's
library version. `examples/embedding/` is a complete working server built this way, and it is
compiled in CI — so if the embedding path breaks, a build fails rather than a copy silently rots.

If you would rather not embed anything, `loom mcp serve` runs the same tools over stdio.

## If you are porting to another language

`porting-guides/` covers TypeScript, Python and Java. Read the **real** implementation as your
source — `shared/mcp/internal/tools/` and `shared/mcp/internal/analyzers/` — not a copy of it. That
code is compiled and tested on every CI run, which is exactly the property the deleted templates
lacked.

## Structure

- `porting-guides/` — language porting guides (`typescript.md`, `python.md`, `java.md`)
