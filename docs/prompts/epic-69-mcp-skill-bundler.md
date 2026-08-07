# Epic 69 — MCP Skill Auto-Bundler

Source: `docs/audits/framework-audit-2026-08-07.md` §3 item 1; extends Epic 43 (`v3.3.12`).

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Prior context

Epic 43 (`v3.3.12`) shipped `shared/mcp/` — a Go module with `cmd/mcp-server/main.go`
containing 6 M1 tools hardcoded by hand: `search-ki`, `analyze-complexity`,
`check-ubiquitous-language`, `team-topology-check`, `list-agents`, `resolve-model-tier`.
`install.sh --with-mcp` compiles and links the binary.

The problem: 69 skills exist in `shared/skills/*/SKILL.md` and none are auto-discovered.
Adding a new skill to the MCP server today means manually editing `main.go`. The audit
calls for auto-bundling framework skills into the MCP server without manual wiring.

Existing surface:
- `shared/skills/*/SKILL.md` — each file has YAML frontmatter (`name`, `description`,
  `standalone_mode`, etc.; see `shared/contracts/skill-frontmatter-contract.md`)
- `shared/schemas/skill-frontmatter.schema.json` — JSON schema for frontmatter validation
- `shared/mcp/` — Go module with `go.mod`, `cmd/mcp-server/`, `analyzers/`, `tools/`
- `scripts/health-check.sh` — the fitness function harness

## Scope

**Phase A — Design (one commit, then PAUSE for user approval):**

Draft and commit as `docs(mcp): design skill-bundler schema (Epic 69 Phase A)`:

- Decide: which SKILL.md frontmatter fields map to MCP tool name / description / input
  schema? Propose a new optional frontmatter field (e.g., `mcp_bundleable: true`) or an
  opt-in list approach. Both options must be sketched with rationale.
- Decide: output format — generate a Go source file (`shared/mcp/generated/tools_gen.go`)?
  Or a JSON registry that `main.go` reads at startup? Tradeoff: generated source requires a
  compile step; runtime registry requires a Go embed.
- Decide: which existing 6 M1 tools map to which skills (if any), and whether they should
  be migrated to the auto-bundled approach or stay as hand-crafted implementations.
- Define the fitness function: `scripts/check-mcp-drift.sh` must detect when a skill
  marked `mcp_bundleable` is not reflected in the registered tool set.
- Write the design as `shared/mcp/skill-bundler-design.md` (≤ 3 pages).

**Phase B — Implementation (after approval; one commit per op):**

Op 1 — `feat(mcp): add mcp_bundleable field to skill frontmatter schema (Epic 69 Op 1)`:
- Add `mcp_bundleable` (boolean, optional, default false) to
  `shared/schemas/skill-frontmatter.schema.json` and document in
  `shared/contracts/skill-frontmatter-contract.md`.
- Mark 5–10 broadly-applicable skills as `mcp_bundleable: true` as a seed set (e.g.,
  `search-ki`, `analyze-complexity`, `check-ubiquitous-language`).

Op 2 — `feat(scripts): generate-mcp-tools.sh skill-to-MCP bundler (Epic 69 Op 2)`:
- `scripts/generate-mcp-tools.sh`: reads all SKILL.md files where `mcp_bundleable: true`,
  emits output per the format chosen in Phase A.
- Must be idempotent (re-runnable without duplicate output).
- Dry-run mode: `--dry-run` prints what would change, exits 0.

Op 3 — `feat(mcp): generated tool registration from bundler output (Epic 69 Op 3)`:
- Wire the bundler output into `shared/mcp/` so `cmd/mcp-server` exposes bundled tools.
- `make generate` or equivalent regenerates the output.
- Verify the server compiles: `cd shared/mcp && go build ./...`.

Op 4 — `feat(scripts): check-mcp-drift.sh + health-check integration (Epic 69 Op 4)`:
- `scripts/check-mcp-drift.sh`: counts skills marked `mcp_bundleable: true`, counts
  registered MCP tools in the generated output, FAILs if they diverge.
- Wire as a FAIL-level check in `scripts/health-check.sh`.
- `bash scripts/health-check.sh` must pass.

Op 5 — `docs(mcp): update README and usage guide (Epic 69 Op 5)`:
- Update `shared/mcp/README.md` with: bundler usage, how to mark a skill as bundleable,
  how to regenerate, and the `--with-mcp` install flow.

After every commit: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- If Phase A concludes that auto-generated tool wrappers can't produce correct MCP input
  schemas without per-skill manual parameter definitions (beyond what frontmatter provides),
  halt and propose a lighter alternative: a skill-registry JSON that the server reads to
  expose skills as single-parameter `{ "input": string }` tools.
- If the 6 existing M1 tools would need to be rewritten to fit the auto-bundled shape and
  the rewrite is non-trivial, keep them as hand-crafted and document the two-tier
  coexistence in `shared/mcp/README.md` rather than forcing migration.
- If `go build` fails on the generated output for any non-trivial reason, halt with the
  exact error before proceeding to Op 4.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A rulings:
  - Frontmatter field approach: <mcp_bundleable field | opt-in list | other>
  - Output format: <generated Go source | runtime JSON registry>
  - M1 tool migration: <migrate all | keep hand-crafted + add auto-bundled>
  - check-mcp-drift check level: <FAIL | WARN>

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, go build <pass>, mcp-drift check <pass>.
```

Go.
