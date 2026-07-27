# Architecture Decision & Plan: Extract Framework-Generic MCP Patterns Into `shared/mcp-patterns/`

- **Status**: Draft (Phase A)
- **Date**: 2026-07-27
- **Target Repository**: `/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`

---

## 1. Problem Statement & Motivation

Currently, `docs/prompts/install-framework-with-mcp-bridge.md` and `docs/prompts/update-installed-framework.md` instruct executing agents to inspect tool patterns directly from `saturday-mcp/internal/tools/*.go`. 

This creates an accidental coupling: downstream framework consumers are forced to clone `saturday-monorepo` locally even when their project has no dependency on Saturday domain logic.

Extracting framework-generic Model Context Protocol (MCP) tool patterns, retrievers, analyzers, and server registration scaffolding into `shared/mcp-patterns/` removes this coupling, allowing framework installations to supply a self-contained reference implementation.

---

## 2. Component Split Analysis

`saturday-mcp` contains 22 tools. The classification between framework-generic patterns to extract and Saturday-domain-specific tools to retain is detailed below.

### 2.1 Framework-Generic (Extract into `shared/mcp-patterns/`)
Total: 6 M1 tools + 5 Analyzers/Utilities + Retrieval & Registration Infrastructure

- **Tools (6)**:
  - `analyze_complexity_tool.go` (Cyclomatic & function length analysis)
  - `check_accessibility_tool.go` (UI/HTML accessibility violation scanning)
  - `check_ubiquitous_language_tool.go` (Domain dictionary compliance scanner)
  - `verify_dependencies_tool.go` (Clean architecture layer import boundary checker)
  - `search_ki_tool.go` (LLM-as-retriever for Knowledge Items)
  - `search_docs_tool.go` (BM25 documentation search via SQLite FTS5)
- **Retrievers & Adapters**:
  - `retriever.go` (`Retriever` interface, `Reference` struct, `KICorpusRetriever`)
  - `bm25_retriever.go` (SQLite FTS5 index & query engine)
- **Analyzers & Utilities**:
  - `complexity_analyzer.go`
  - `accessibility_analyzer.go`
  - `ubiquitous_language_analyzer.go`
  - `dependency_boundary_analyzer.go`
  - `walkutil.go` (File-system walking and filter helpers)
- **Infrastructure & Test Patterns**:
  - `responses.go` (Extracted M1 response structs)
  - `testfixtures.go` (Reusable test helper patterns extracted from `testfixtures_test.go`: `writeFile`, `buildRequest`, `extractText`, `silentLogger`)
  - `registration.go` (Loop-based tool registration pattern)
  - `tool_provider.go` (Provider function pattern returning `[]domain.Tool`)

### 2.2 Saturday-Specific (Retain in `saturday-mcp`)
Total: 16 tools

- **Generators (6)**: `generate_site`, `generate_page`, `generate_flow`, `generate_steps`, `generate_element`, `generate_service`
- **Analysis & Code Modifiers (8)**: `migrate_code`, `generate_documentation`, `analyze_framework`, `validate_patterns`, `suggest_improvements`, `analyze_impact`, `analyze_performance`, `parse_test_failure`
- **Workflows (2)**: `run_tests`, `prioritize_tests`

---

## 3. Target Directory Structure

The extracted patterns will be stored in `shared/mcp-patterns/` with the following structure:

```
shared/mcp-patterns/
├── README.md                    ← Purpose, design goals, usage instructions
├── go/                          ← Go reference implementation (source-to-copy)
│   ├── tools/
│   │   ├── retriever.go         ← Retriever interface & KICorpusRetriever
│   │   ├── bm25_retriever.go
│   │   ├── analyze_complexity_tool.go
│   │   ├── check_accessibility_tool.go
│   │   ├── check_ubiquitous_language_tool.go
│   │   ├── verify_dependencies_tool.go
│   │   ├── search_ki_tool.go
│   │   ├── search_docs_tool.go
│   │   ├── responses.go         ← M1 response structs
│   │   └── testfixtures.go      ← Reusable test helpers (renamed from _test.go)
│   ├── analyzers/
│   │   ├── complexity_analyzer.go
│   │   ├── accessibility_analyzer.go
│   │   ├── ubiquitous_language_analyzer.go
│   │   ├── dependency_boundary_analyzer.go
│   │   └── walkutil.go
│   ├── server/
│   │   ├── registration.go      ← Loop-based tool registration pattern
│   │   └── tool_provider.go     ← Provider function pattern
│   └── README.md                ← Go-specific adaptation & copying instructions
└── porting-guides/
    ├── typescript.md            ← TS implementation guidance & mapping
    ├── python.md                ← Python implementation guidance & mapping
    └── java.md                  ← Java implementation guidance & mapping
```

---

## 4. Language Strategy & Decoupled Reference Model

1. **Go as Primary Reference**: Go matches `saturday-mcp`'s source implementation, providing concrete, copyable code for Go-based bridge servers.
2. **Prose Porting Guides**: Architectural patterns, interface contracts, and library mappings for TypeScript, Python, and Java will be documented in `shared/mcp-patterns/porting-guides/`.
3. **Copy (Not Move) Policy**: The extraction is a non-destructive copy from `saturday-mcp`. `saturday-mcp` retains its files and serves as a working downstream consumer. `shared/mcp-patterns/` becomes the authoritative blueprint for new bridges.
4. **Header Comments & Build Tags**: Extracted Go files will feature header comments explaining required module path adjustments and `// +build ignore` build tags to prevent local compilation conflicts within the framework repo.

---

## 5. Phased Implementation Roadmap

- **Phase A**: Draft extraction plan (this document), present Phase A report, obtain user confirmation.
- **Phase B**: Extract files into `shared/mcp-patterns/go/` across 5 atomic file-group commits (Op B1 to Op B5).
- **Phase C**: Update prompt files (`install-framework-with-mcp-bridge.md` and `update-installed-framework.md`) to reference `shared/mcp-patterns/`.
- **Phase D**: Author language porting guides (`typescript.md`, `python.md`, `java.md`).
