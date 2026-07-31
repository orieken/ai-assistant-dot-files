# Epic 43 — Standardized MCP Tool Packaging (`shared/mcp/` exporter)

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 1; re-confirmed open (lowest
priority of the remaining standalone epics) by `docs/audits/framework-gap-audit-2026-07-31.md`
§ F5. The 07-25 audit noted this epic needed scope clarification before drafting — Phase A below
IS that clarification; do not skip it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

The gap as originally stated: no canonical way to package framework capabilities
(analyze-complexity, check-accessibility, search-ki, etc.) as standalone MCP servers consumable
by Cursor MCP, Claude Desktop, or custom runtimes. What exists since then:

- `shared/mcp-patterns/` — framework-generic MCP *patterns* extracted from `saturday-mcp`
  (retrievers, analyzers, walkutil, 6 M1 tools, response structs, registration pattern). Patterns,
  not an exporter.
- `mcp-add` skill — retrofits/extends an *existing* MCP server to framework conventions.
- `shared/blueprints/mcp-server.md` — the authoritative MCP server pattern source.
- `docs/prompts/done/install-framework-with-mcp-bridge.md` — the manual per-project bridge flow
  (Option B: copy source, keep MCP sovereign) this epic would partially automate.
- Backend convention: Go (root `CLAUDE.md` technology stack; `shared/mcp-patterns/` is
  Go-shaped).

## Scope

**Phase A — Scope ruling (one commit, then PAUSE for user approval):**

Answer, in a design note committed as
`docs(mcp): scope ruling for shared/mcp exporter (Epic 43 Phase A)`:

1. **What is actually being exported?** Framework skills are markdown prompts, not code. An "MCP
   server exposing framework skills" means one of: (a) a generator emitting a Go server scaffold
   whose tools wrap the deterministic scripts (`analyze-complexity` → `scripts/` logic, `search-ki`
   → lexical search) — real code, no LLM; (b) a static `mcp.json`/manifest generator pointing
   existing clients at an already-built server (thin); or (c) declaring `saturday-mcp` + the
   bridge prompt the permanent answer and closing this epic as superseded. Recommend one with
   rationale.
2. **Which capabilities are exportable?** Only the ones with deterministic, non-LLM cores. List
   them explicitly (the 6 M1 tools in `shared/mcp-patterns/` are the starting candidates).
3. **Where does generated output live?** (`shared/mcp/` template + per-project generation via
   `install.sh` flag, vs. a checked-in reference server.)

**Phase B — Implementation (after approval; scope depends on the Phase A ruling):**

If (a): `shared/mcp/` generator + template following `shared/blueprints/mcp-server.md` Clean
Architecture layers; wire an `install.sh --with-mcp` (or similar) flag; docs + README counts via
`scripts/check-inventory-drift.sh`; `bash scripts/health-check.sh` green after every commit.
If (b): manifest generator + docs only.
If (c): close the epic — update the two audits' checklists + this prompt to `done/` with a
"superseded" note.

One commit per op, enumerated in your Phase B plan before executing.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If Phase A ruling is (c), that IS the deliverable — don't build anything to justify the epic's
  existence.
- If the exporter would need to duplicate logic that lives in `saturday-mcp` proper (not in
  `shared/mcp-patterns/`), halt — that re-creates the coupling `add-mcp-patterns-directory.md`
  deliberately broke.
- Any new external dependency (Go modules beyond stdlib + what mcp-patterns already uses): halt
  for approval.

## Report (under 150 words)

```
Phase A commit: <sha>
Ruling: <a|b|c> — <one-line rationale>
Exportable capabilities: <list or "n/a">

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, check-inventory-drift <pass>
```

Go.
