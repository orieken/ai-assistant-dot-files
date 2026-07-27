# Shared MCP Patterns (`shared/mcp-patterns/`)

This directory contains framework-generic Model Context Protocol (MCP) tool patterns, retrievers, analyzers, and server registration abstractions.

## Purpose

When downstream projects build or extend an MCP bridge server to interface with the Context Engineering Framework, they can copy reference implementations and patterns directly from this directory rather than requiring a separate checkout of `saturday-mcp`.

## Structure

- `go/`: Reference implementation in Go (source files with header instructions and build tags to copy from).
  - `go/tools/`: 6 M1 framework tools (`analyze_complexity`, `check_accessibility`, `check_ubiquitous_language`, `verify_dependencies`, `search_ki`, `search_docs`), `retriever.go`, `bm25_retriever.go`, response structs, and test fixtures.
  - `go/analyzers/`: Concrete code and document analyzers (`complexity_analyzer`, `accessibility_analyzer`, `ubiquitous_language_analyzer`, `dependency_boundary_analyzer`, `walkutil`).
  - `go/server/`: Loop-based tool registration (`registration.go`) and provider function patterns (`tool_provider.go`).
- `porting-guides/`: Language porting guides for non-Go implementations (`typescript.md`, `python.md`, `java.md`).

## Copying & Adapting Patterns

1. Copy the relevant files into your project's internal directory structure.
2. Remove the `//go:build ignore` build tags from copied files.
3. Update package names and module import paths (`<YOUR_MODULE>/*`) to match your codebase.
