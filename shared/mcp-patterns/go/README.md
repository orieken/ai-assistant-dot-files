# Go Reference Implementation for MCP Patterns (`shared/mcp-patterns/go/`)

This directory contains the reference Go implementation of framework-generic MCP tools, analyzers, retrievers, and registration scaffolding extracted from `saturday-mcp` (commit `0e5549125b7129e2b308df09d99e10d0b29a41bb`).

## Overview of Components

### `tools/`
- `retriever.go`: `Retriever` interface & `KICorpusRetriever` (LLM-as-retriever tier 1)
- `bm25_retriever.go`: SQLite FTS5 BM25 retrieval engine (tier 2)
- `analyze_complexity_tool.go`: Cyclomatic complexity and function length checker tool
- `check_accessibility_tool.go`: Semantic-HTML & ARIA scanner tool
- `check_ubiquitous_language_tool.go`: Domain dictionary term scanner tool
- `verify_dependencies_tool.go`: Clean Architecture layer boundary checker tool
- `search_ki_tool.go`: Knowledge Item search tool
- `search_docs_tool.go`: Documentation BM25 search tool
- `responses.go`: Response structs for all 6 tools
- `schemas.go`: Shared input schema helpers
- `testfixtures.go`: Reusable test helpers (`WriteFile`, `BuildRequest`, `ExtractText`, `SilentLogger`)

### `analyzers/`
- `complexity_analyzer.go`
- `accessibility_analyzer.go`
- `ubiquitous_language_analyzer.go`
- `dependency_boundary_analyzer.go`
- `walkutil.go`

### `server/`
- `registration.go`: Registration loop over `[]domain.Tool`
- `tool_provider.go`: Provider function pattern returning `[]domain.Tool`

## Usage Instructions

All Go reference files include `//go:build ignore` so they do not attempt to compile inside `ai-assistant-dot-files` directly. When copying files into your Go project:
1. Delete the `//go:build ignore` and `// +build ignore` lines at the top of each file.
2. Replace `package tools` / `package analyzers` / `package server` if your internal package names differ.
3. Replace placeholder import paths (`github.com/orieken/saturday-mcp/...` or `<YOUR_MODULE>/...`) with your project's module path.
