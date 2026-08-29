# Epic 75 — Distribution & Adoption: MCP as the Portable Surface, Maturity Levels as Product

Source: distribution-strategy discussion 2026-08-29. Operationalizes roadmap items **D.1–D.5**
(BUILD-ROADMAP.md, "Workstream: PLATFORM — Distribution & Adoption"). Read that section first —
it carries the Problem/Fix/Done-when detail; this file is the executable handoff.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — the D.1–D.5 items and their dependency edges. **Respect
   Blocked-by**: D.2 needs M0.3 + L2.4 shipped; D.5 additionally needs L2.9/L2.11/L2.12 for its
   later tools. Before starting a phase, verify its blockers are actually done (check git log and
   the roadmap's checkboxes) — if a blocker is unshipped, STOP and report instead of improvising
   around it.
2. `docs/roadmaps/architectural-audit-2026-08-29.md` — findings H1–H11; the *why*.
3. The strategy, in three sentences, so no phase drifts from it:
   - **MCP is the portable surface; markdown is the content.** Executable capabilities ship as MCP
     tools usable from any MCP host; the dotfile export becomes the Level 1 convenience layer.
   - **Both distribution modes from one Go module**: `loom mcp serve` (standalone) AND a semver'd
     public embedding package (`Tool` interface + registry). Extension = implement + `Register`.
   - **MCP exposes tools and resources only.** The orchestration kernel is `loom run` acting as an
     MCP *client* — the pipeline itself is never a tool someone else calls.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Server distribution | `loom mcp serve` subcommand in the existing `loom` binary | One artifact in the Homebrew tap; `shared/mcp/cmd/mcp-server/` retired after one release cycle |
| Transport | stdio now; network waits for L2.8 | Don't invent auth ad hoc |
| Embedding API | Public package with transport-free `Tool` + `ToolRegistration` + `Registry.Merge`; no `internal/` or `mcp-go` types in its signatures | Consumers must not inherit our mcp-go version pin |
| Extension mechanism | Compile-time Go embedding first | Typed and simple; subprocess plugins only on demonstrated demand |
| Level profiles | Data-driven (`shared/levels.yaml`), consumed by installer and health | Levels are product, not prose |
| Level reporting | `loom health` infers level from mechanical evidence only | Never report a level whose enforcement layer isn't installed and answering |
| New MCP tools | Only behind landed code; read-only first; mutating pipeline state stays executor-exclusive | No prompt-as-runtime one layer down |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase (per-epic commit discipline), then PAUSE
  for human review before the next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any`, errors handled explicitly, interfaces defined at the consumer.
- Every phase ends by running: `go build ./...`, `go test ./...`, `golangci-lint run ./...` in
  each touched module, plus `scripts/health-check.sh` if `shared/` content changed.
- Update docs in the same phase that changes behavior (`README.md`, `shared/mcp/README.md`,
  `cmd/loom/README.md`) — do not defer to a docs-later phase.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — `loom mcp serve` (roadmap D.1) — UNBLOCKED

**Goal**: the brew-installed `loom` binary serves the six framework MCP tools over stdio.

1. Decide module unification: prefer merging `shared/mcp`'s module into the root `go.mod`
   (single module, single goreleaser build). If merge is disproportionate (dependency conflicts),
   use a `go.work` + `replace` and record why in the commit body. Either way the released binary
   embeds the server.
2. Add `cmd/loom/cmd/mcp.go` (parent `loom mcp` command) and `mcp_serve.go`:
   `loom mcp serve [--log-file <path>]` constructs the mcp-go stdio server, calls
   `register.FrameworkTools`, and blocks until EOF/SIGINT. Structured JSON logs to stderr or
   `--log-file`, never stdout (stdout is the MCP wire).
3. Keep `shared/mcp/cmd/mcp-server/` building; add a deprecation note to its `main.go` doc comment
   and `shared/mcp/README.md` pointing at `loom mcp serve`.
4. Update the goreleaser config so the one binary ships; verify `goreleaser build --snapshot`
   locally if goreleaser is installed, otherwise note it for CI.
5. Tests: a smoke test that spawns the built binary, performs an MCP `initialize` + `tools/list`
   handshake over stdio pipes, and asserts all six tool names are present.
6. Docs: `cmd/loom/README.md` gains an "MCP server" section with a sample `.mcp.json` /
   `claude mcp add` snippet.

**Done when** (D.1): `loom mcp serve` answers `tools/list` with all six tools, one binary per
platform ships. **Commit** (`feat(loom): serve framework MCP tools from the loom binary`), report,
PAUSE.

## Phase B — Public embedding package (roadmap D.2) — BLOCKED BY M0.3, L2.4

**Verify first**: the domain layer no longer imports `mcp-go` types (M0.3) and the tool registry
exists (L2.4). If either is missing, STOP and report — this phase must not re-implement them.

1. Create the public package (working name `tools/` at module root; final import path
   `github.com/orieken/loom/tools` — align with whatever module path Phase A settled). It exports:
   the transport-free `Tool` interface, `ToolRegistration` (timeout, retry class, permission
   scope), `Registry` with `Register` and `Merge`, and a `Frameworks()` constructor returning the
   built-in registrations.
2. Nothing in the public package's exported signatures may reference `internal/` packages,
   `mark3labs/mcp-go`, or `invopop/jsonschema`. Enforce with a small fitness-function test that
   parses the package's exported API (e.g. `go/packages`) and fails on forbidden imports in
   exported types.
3. Rewrite `register.FrameworkTools` as a thin compatibility wrapper over the new API; mark
   deprecated in its doc comment with the replacement spelled out.
4. Build an out-of-repo example (`examples/embedding/` with its own `go.mod` + `replace`): a
   minimal MCP server that merges loom's registry and adds one custom `Tool`. CI must build it.
5. Documentation: `shared/mcp/README.md` gains an "Embedding loom's tools" section with the
   example inlined; state the semver compatibility promise explicitly.
6. Do NOT tag a release version in this phase — tagging is a human decision (escalate with a
   recommended version).

**Done when** (D.2): the example project registers a custom tool with no `internal/` or `mcp.*`
import, and CI builds it. **Commit** (`feat(mcp): public embedding API for framework tools`),
report, PAUSE.

## Phase C — Level profiles: `loom init --level N` (roadmap D.3) — BLOCKED BY Phase A

1. Author `shared/levels.yaml`: four profiles, each listing the rule files, agent files, skill
   dirs, and config actions it installs. L1 = core rules (guardrails, approval gates, memory trust
   boundary) + agents/skills prompts. L2 = + MCP server registration + workflow YAML (+ executor
   once M0.4 exists — gate that entry on a `requires:` field so the installer skips features whose
   implementation hasn't landed, with a warning). L3/L4 entries mirror the roadmap milestones and
   will mostly carry `requires:` gates today — that's expected and honest.
2. Split `shared/rules/` into always-on core vs. on-demand modules (language conventions, IaC,
   testing stacks move to on-demand). Preserve every file's content — this is a re-bucketing, not
   a rewrite. Update the exporters/installer so only the core set is injected by default; document
   how a project opts specific on-demand modules in.
3. Wire `--level N` into `loom install` (and an `init` alias if trivially cheap): selects the
   bundle from `levels.yaml`. Default with no flag = current full-install behavior, unchanged.
4. Measure and record the injected-context size per level (a test asserting the L1 core bundle
   stays under a documented ceiling — pick the ceiling from the actual measured size + headroom,
   record it in `levels.yaml`).
5. Update `README.md`'s install section with the level ladder table.

**Done when** (D.3): `loom install --level 1` installs only the core bundle under the documented
ceiling; `--level 2` adds exactly the L2 delta. **Commit**
(`feat(loom): maturity-level install profiles`), report, PAUSE.

## Phase D — `loom health` level report (roadmap D.4) — BLOCKED BY Phase C

1. Add a level-assessment check set to `cmd/loom/cmd/health_checks.go`, driven by
   `shared/levels.yaml`: bundle installed? MCP server configured in the platform config AND
   answering `tools/list`? workflow definitions present? telemetry stream present? executor in
   use? Each check yields evidence strings.
2. Inference rule: report the highest level whose *entire* mechanical evidence set passes; list
   the failing evidence for level N+1 as the gap checklist. Documentation-only presence never
   satisfies an evidence item.
3. Output section in `health_output.go`: current level, passing evidence, next-level gaps.
4. Table-driven unit tests covering inference for all four levels plus the "docs present but
   server dead → still Level 1" case.

**Done when** (D.4): fresh L1 install → `loom health` prints "Level 1" + concrete L2 gaps; tests
cover all levels. **Commit** (`feat(loom): health reports agentic maturity level`), report, PAUSE.

## Phase E — Capability tools on the MCP surface (roadmap D.5) — BLOCKED BY Phases A+B, then per-tool

Each tool below has its own blocker; implement only those whose blocker has landed, in this order,
one commit per tool:

| Tool | Read/write | Blocked by | Behavior |
|---|---|---|---|
| `validate_artifact` | read | L2.11 | Structural contract validation of an artifact file against `shared/contracts/`; returns typed violations |
| `pipeline_state` | read | L2.12 | Read a run's state (stages, checksums, gate status) for a named feature workspace |
| `query_telemetry` | read | L3.9 | Filtered read over `.claude/telemetry/events.jsonl` (event type, agent, time range) |
| `evaluate_policy` | read | L2.16 | Dry-run policy evaluation for a gate + context; returns proceed/halt/require-human + matched policies |

Rules: register via the L2.4 registry with explicit timeout and permission scope; input schemas
enforced server-side (L2.1 pattern); no tool mutates pipeline state — if a design pressure appears
to want that, STOP and escalate (it belongs in the executor).

**Done when** (D.5): an MCP host with no loom markdown installed validates an artifact and reads
pipeline state via tool calls alone. **Commit per tool**, report, PAUSE.

---

## Report format (end of every phase)

```
## Epic 75 Phase <X> Report
- Roadmap item: D.<n> — <title>
- Blockers verified: <list, with evidence (commit SHAs / files)>
- Commits: <sha> <subject>
- Build/lint/test: go build PASS|FAIL · golangci-lint PASS|FAIL · go test PASS|FAIL (counts)
- Done-when criterion: <restate it> — MET | NOT MET (why)
- Escalations / open questions: <list or "none">
- Next phase blocked by: <what must land first>
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
