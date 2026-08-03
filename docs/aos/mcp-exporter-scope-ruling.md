# MCP Exporter Scope Ruling (Epic 43 Phase A)

**Decision date**: 2026-08-02
**Status**: Phase A complete — awaiting Phase B approval

---

## Question 1: What is actually being exported?

**Ruling: Option (a) — checked-in reference scaffold (not a generator, not option b, not c).**

The framework already has `shared/mcp-patterns/go/` — a full set of Go source files
(6 tools, 4 analyzers, server wiring) extracted from `saturday-mcp`, tagged with
`//go:build ignore` so they do not compile inside this repo. These are "copy-me"
templates, not a runnable server.

The gap: no pre-assembled, compilable reference server exists. Someone who wants to run
framework tools via MCP today must:
1. Copy individual files from `shared/mcp-patterns/go/` by hand
2. Remove the `//go:build ignore` tags
3. Wire a `go.mod`, `cmd/`, and `main.go` from scratch
4. Replace placeholder import paths

**Option (a)** delivers: a `shared/mcp/` directory that IS that assembled module —
checked in, compilable, ready to copy into a project or run directly.

**Why not (b)?** A static `mcp.json` manifest has nothing to point at — there is no
pre-built binary in this repo, and generating a pointer to a non-existent binary adds no
value.

**Why not (c)?** The `install-framework-with-mcp-bridge.md` prompt (in `done/`) covers
the "existing MCP server — bridge into it" case. The scaffold covers the orthogonal case:
"no existing MCP server — give me a standalone one." These are complementary, not
duplicates. The bridge prompt also still requires an agent to manually execute 10 ops;
a scaffold reduces that to "copy directory, update module path, `go build`."

**Escalation check (duplication of saturday-mcp proper)**: The scaffold uses only
`shared/mcp-patterns/go/` — which was extracted from `saturday-mcp` specifically to
enable this use case. No logic from `saturday-mcp` outside `shared/mcp-patterns/` is
needed. Escalation condition does not trigger.

**External dependencies**: `github.com/mark3labs/mcp-go` (the MCP Go SDK) and
`github.com/invopop/jsonschema` — both already referenced in `shared/mcp-patterns/go/`
files. No new external dependency is introduced.

---

## Question 2: Which capabilities are exportable?

Only those with a **deterministic, non-LLM core** in `shared/mcp-patterns/go/`:

| Tool | Analyzer core | Input | Output |
|---|---|---|---|
| `analyze_complexity` | `complexity_analyzer.go` — Go AST cyclomatic complexity + LOC | file/dir path | violations JSON |
| `check_accessibility` | `accessibility_analyzer.go` — HTML/ARIA semantic rules | file/dir path | violations JSON |
| `check_ubiquitous_language` | `ubiquitous_language_analyzer.go` — DOMAIN_DICTIONARY term scan | file/dir path | term mismatches JSON |
| `verify_dependencies` | `dependency_boundary_analyzer.go` — Go/TS import layer checker | file/dir path | boundary violations JSON |
| `search_ki` | `retriever.go` + `bm25_retriever.go` — SQLite FTS5 BM25 search | query string | KI excerpts |
| `search_docs` | `bm25_retriever.go` — same engine, docs corpus | query string | doc excerpts |

`KICorpusRetriever` (LLM-as-retriever) is in `retriever.go` but is excluded from M1 —
it requires an LLM call (violates "no LLM" constraint for this scope). The BM25 path is
the exportable tier.

---

## Question 3: Where does generated output live?

**`shared/mcp/`** — checked-in reference module with its own `go.mod`.
Module path: `github.com/orieken/ai-assistant-dotfiles/mcp` (or a local-friendly path
the README instructs to replace).

Usage modes:
- **Standalone**: `cd shared/mcp && go build ./cmd/mcp-server` → single binary, point
  any MCP client at it via stdio.
- **Copied into project**: copy `shared/mcp/` to `<project-root>/<name>-mcp/`, update
  module path in `go.mod`, `go build`. Equivalent to running `install-framework-with-mcp-bridge.md`
  Phase C Ops 1–8 automatically.
- **`install.sh --with-mcp`** (Phase B Op 7): copies the scaffold to a target project,
  updates the module path to `github.com/<org>/<project>-mcp`, and runs `go mod tidy`.

The scaffold does NOT generate into a per-feature workspace — it lives at `shared/mcp/`
permanently as a maintained reference.

---

## Phase B Commit Sequence (for approval)

If approved, one commit per operation:

| Op | Commit message | What |
|---|---|---|
| 1 | `feat(mcp): add shared/mcp module skeleton (Epic 43)` | `shared/mcp/go.mod`, `go.sum`, `.gitignore`, `.env.example`, `.golangci.yml` |
| 2 | `feat(mcp): add cmd/mcp-server entrypoint (Epic 43)` | `cmd/mcp-server/main.go` — wires stdio transport, registers tool provider |
| 3 | `feat(mcp): copy analyzers from mcp-patterns (Epic 43)` | `internal/analyzers/` — 4 files + walkutil, build tags removed |
| 4 | `feat(mcp): copy tools foundation from mcp-patterns (Epic 43)` | `internal/tools/retriever.go`, `bm25_retriever.go`, `responses.go`, `schemas.go`, `testfixtures.go` |
| 5 | `feat(mcp): copy 6 M1 tools from mcp-patterns (Epic 43)` | `internal/tools/analyze_complexity_tool.go` + 5 peers, import paths updated |
| 6 | `feat(mcp): add server wiring from mcp-patterns (Epic 43)` | `internal/server/registration.go`, `tool_provider.go` |
| 7 | `feat(mcp): add install.sh --with-mcp flag (Epic 43)` | `install.sh` update + `scripts/test-install.sh` update |
| 8 | `docs(mcp): add shared/mcp README and usage guide (Epic 43)` | `shared/mcp/README.md`, prose counts in root `README.md` |
| 9 | `docs(mcp): mark Epic 43 complete (Epic 43)` | move `epic-43-mcp-exporter.md` → `docs/prompts/done/` |

Verification after Op 9: `bash scripts/health-check.sh` green,
`bash scripts/check-inventory-drift.sh` green, `cd shared/mcp && go build ./...` green.

---

## Risks and Constraints

- `shared/mcp/go.sum` will be a tracked file (~100–200 lines for 2 external dependencies).
  This is normal for a Go module — no different from any project's lockfile.
- The scaffold's `go.mod` pins `github.com/mark3labs/mcp-go` at the version currently
  used by saturday-mcp (to be confirmed at Phase B Op 1 via `grep` of the extracted files).
- `shared/mcp/` uses the same `//go:build ignore` removal strategy as `install-framework-with-mcp-bridge.md`
  Phase C — the patterns' build tags exist to keep them from compiling inside *this* repo;
  the scaffold removes them and becomes its own module.
- `KICorpusRetriever` (LLM-as-retriever) is scaffolded as a no-op struct with a clear
  comment that it requires an LLM backend — it is not wired into the default tool provider.
  Teams can enable it if they have an appropriate LLM adapter.
