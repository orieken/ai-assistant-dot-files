# loom CLI

The `loom` binary installs the framework's agents, skills, and rules across AI
coding platforms and serves the framework's MCP tools. Install it via
`brew install orieken/tap/loom` or
`go install github.com/orieken/loom/cmd/loom@latest`.

## Commands

| Command | Purpose |
|---|---|
| `loom install` (alias `loom init`) | Install framework content for detected platforms (`--target`, `--platform`, `--copy`, `--dry-run`, `--stack`, `--level` — see the maturity-level table in the root README; profiles are defined in `shared/levels.yaml`) |
| `loom health` | Verify installed configs match the canonical `shared/` source, and report the project's agentic maturity level (see below) |
| `loom tools status` / `loom tools install` | Report / install opt-in context tools |
| `loom mcp serve` | Serve the framework MCP tools over stdio |
| `loom run` | Execute the delivery pipeline for a feature spec (experimental — see below) |
| `loom uninstall` | Remove installed framework content |
| `loom update` | Update installed framework content |

## Running pipelines (experimental)

`loom run` is the executor skeleton decided by ADR-006 (loom executes
pipelines) and built by roadmap item M0.4. It runs the built-in linear
`deliver-feature` plan — the same 14-agent sequence as the markdown
pipeline — stage by stage, persisting durable state after every transition.

```bash
# Execute the pipeline for a spec (spawns `claude -p` per stage)
loom run --spec docs/features/user-auth/spec.md

# Continue an interrupted run from its checkpoint
loom run --spec docs/features/user-auth/spec.md --resume

# Deterministic dry run with canned artifacts (no LLM calls)
loom run --spec docs/features/user-auth/spec.md --provider mock
```

- **State**: `.claude/feature-workspace/<feature>/run-state.json` — schema-versioned,
  written atomically (temp file + rename), with per-stage status, timestamps, and
  the SHA-256 of each stage's artifact. A fresh run refuses to start over existing
  state (pass `--resume` or delete the file); `--resume` requires existing state.
- **Integrity on resume**: before the run loop trusts anything, every COMPLETED stage's artifact is
  re-hashed in Go and compared to the digest recorded when it finished. A changed or missing
  artifact demotes that stage to `STALE` — it re-runs — and the demotion cascades to every stage
  completed after it, whose own output came from content that no longer exists. The resume prints
  what it invalidated and why. There is no flag to skip verification (roadmap L2.12).
- **Interruption**: the first Ctrl-C cancels the in-flight stage and persists a
  clean `INTERRUPTED` checkpoint; a second Ctrl-C kills immediately. `--resume`
  skips completed stages and re-runs the interrupted one.
- **Providers**: `claude` (default) spawns the `claude` CLI headless per stage,
  building the prompt from the agent's `shared/agents/<agent>.md` definition; if
  the binary is missing the stage fails with a remediation message — there is no
  silent fallback. `mock` is for tests and dry runs.

- **Approval gates**: a gated stage does not start until a human approves its gate
  (roadmap L2.13). The built-in plan gates `developer` behind `confirm-design`,
  `qa-engineer` behind `confirm-security`, and `devops-engineer` behind
  `confirm-ship` — the same stops `deliver-feature` asks a human for. On a
  terminal you are asked `approve gate "X" for stage "Y"? [y/N]` at the barrier;
  otherwise the run persists `WAITING_APPROVAL`, prints the resume command, and
  exits with code **3** so scripts can tell "waiting on a human" from a failure:

  ```bash
  loom run --spec docs/features/user-auth/spec.md --resume --approve confirm-design
  ```

  `--approve` is only valid with `--resume`, and only for the gate the run is
  actually halted on. Nothing an agent returns can approve a gate.
- **Any edit resets the gate**: an approval binds to the SHA-256 of every artifact
  completed when it was given. Edit one of them and the approval is invalidated — the
  record is kept, naming what changed — and the run halts at that gate again until you
  approve the state as it now stands. A stage that re-runs to a byte-identical artifact
  keeps its approval; work completed *after* an approval belongs to the next gate.
  Approving a run whose own verification is about to reset that approval is refused up
  front, rather than recorded and immediately destroyed.

**What it does NOT do yet** (skeleton by design — see `docs/roadmaps/BUILD-ROADMAP.md`):
no `--from-phase` or rollback (L2.15), no retries or backoff, no parallelism
(L3.3), no policy evaluation (L2.16), no conditional stage routing (L3.1) —
every stage of the linear plan runs, including ones the markdown pipeline
would skip — and no telemetry emission (L3.8).

## Recording markdown-pipeline checkpoints

`loom state` is how the `deliver-feature` skill records checkpoints without a model
computing its own integrity hashes (roadmap L2.12). The model keeps routing — conditional
skips, contract retry loops, the code-reviewer loop, none of which `loom run` can do yet —
while this binary reads each artifact and hashes it.

```bash
loom state record --spec docs/features/user-auth/spec.md --stage analyst   --artifact .claude/feature-workspace/user-auth/analysis.md
loom state verify --spec docs/features/user-auth/spec.md   # exits non-zero on any mismatch
loom state approve --spec docs/features/user-auth/spec.md --gate confirm-design
loom state show --spec docs/features/user-auth/spec.md [--json]
```

- No subcommand accepts a caller-supplied digest. A caller that could hand `loom` a hash
  could hand it a wrong one, which is the failure mode being closed.
- Stages carry a **sequence** assigned when first recorded and preserved when re-recorded, so
  a CHANGES REQUESTED loop keeps its place in the run and `verify` knows what is downstream of
  an edited artifact. The markdown pipeline has no fixed plan, so recording order is the only
  ordering available.
- `state approve` **records** a human's approval for audit; it does not enforce the gate.
  Enforcement exists only for stages run by `loom run` (see above). `state verify` does
  report an approval as INVALIDATED once a bound artifact changes, and exits non-zero —
  detection, which is strictly better than the nothing that came before it, and still not
  a barrier.
- `loom state timeline --spec <spec> [--json]` prints the run's event log (see below).
- The two pipelines refuse each other's state files: `loom run` will not resume a run recorded
  by `loom state`, and vice versa. They route differently, so resuming across them would replay
  the wrong work.

## Typed pipeline state

Seven artifacts of the built-in plan are **typed state** instead of markdown (roadmap L2.9):
`analysis`, `architecture`, `route`, `review`, `implementation-notes`, `security-report`, and
`qa-report`. Every hop between them carries fields, not prose.

- **Where**: `.claude/feature-workspace/<feature>/state/<stage>.json`. That document IS the
  stage's artifact, so digest verification and the staleness cascade cover it exactly as they
  cover any other artifact.
- **Schemas**: `shared/schemas/pipeline/*.schema.json`, generated from `internal/state/` by
  `go run ./cmd/gen-schemas`. Never hand-edit them; a test fails when they drift from the
  structs. The schema is inlined into the stage prompt, so a typed run does not depend on the
  framework being installed in the target project.
- **Markdown is a view**: `analysis.md`, `implementation-notes.md`, `qa-report.md` and the rest
  are *rendered* from state under the filenames every contract and downstream agent already
  expects — including the retrieval frontmatter. This is what keeps the typed hop invisible to
  the stages that are still untyped: `sre-engineer`, `devops-engineer` and `accessibility-engineer`
  read `implementation-notes.md` and cannot tell the difference. The view is derived and not
  digest-tracked: annotate it freely, nothing downstream reads it and nothing will demote a stage
  because you did.
- **Projections**: a consuming stage receives the fields its contract needs, not the whole
  upstream document. They are keyed by `(consuming stage, upstream kind)`, because what a stage
  reads and what it writes vary independently — the architect and the qa-engineer both read the
  analysis and get demonstrably different slices of it. A stage with several upstreams, like
  `qa-engineer` with three, gets one labelled block each, so provenance survives and same-named
  fields cannot silently collide.
- **Contract rules become invariants**: two things the contracts already asserted are now checked
  when state loads rather than grepped afterwards — a `qa-report` with a non-zero failed count,
  and a CRITICAL or HIGH security finding with no fix applied, are validation errors that fail the
  stage. A red suite cannot become completed state.
- **No LLM on the data path**: `qa-engineer` and `tech-writer` used to read a model-written summary
  of `analysis.md`. They now receive projections of it (roadmap L2.10).
- **Invalid output fails loudly**: a response that is not a single JSON object (raw, or one
  fenced block) fails the stage, as does one missing a required field or inventing a field the
  schema does not have. There is no repair or retry loop.

The eight remaining artifacts — the end-of-pipeline reports and `context-manifest` — still pass
markdown, unchanged. Nothing evaluates a condition over them yet.

## Routing: which stages actually run

Immediately after the analyst — the earliest point the facts exist — the executor
computes the **Delivery Route** from the typed analysis and records it as
`route.md` plus `state/router.json` (roadmap L3.0). The run prints a one-line
summary and halts at the design gate:

```
routed 7 of 12 stages — skipped architect, performance-engineer, data-engineer,
  accessibility-engineer, devops-engineer (see .../route.md)
Halted at gate "confirm-design" before stage "developer" — approval required.
```

- **Predicates, not prose**: a context crossing or migration or new dependency or
  performance threshold summons the architect; a threshold summons performance; an
  expand/contract migration summons data; an accessibility requirement (which the
  analysis contract makes mandatory for any UI) summons accessibility; a non-empty
  DevOps task list summons devops. Each decision records its reason.
- **The route is approved with the design.** It completes before `confirm-design`,
  so the approval binds its digest: edit `route.md`'s source to force a stage back
  in and the gate resets, exactly as editing any other approved artifact does.
- **Two stages are never routed around.** `code-reviewer` and `security-reviewer`
  always run — an unnecessary review wastes an invocation, a skipped one does not
  fail so cheaply. `visual-qa-engineer` also always runs: its condition is about
  the environment, not the feature (see ADR-007).
- **A gate survives its stage being skipped.** `confirm-ship` guards
  `devops-engineer`; routing devops out still stops at the gate, because the
  checkpoint is about reaching that point in the run.
- `loom state show` lists what was routed out and why.

## The review loop

`developer → code-reviewer` is a **bounded loop** declared in plan data (roadmap
L2.17): the executor reads the reviewer's typed verdict, sends the developer back
when it asks for changes, and counts the rounds.

```
review loop: changes requested — round 2 of 3, re-running from developer
```

- **The condition is a field, not a grep.** `review-approved` reads
  `ReviewState.Verdict`. The rendered `code-review-report.md` still carries the
  bolded literal for the markdown pipeline, but nothing parses prose to decide
  whether to loop.
- **Every round is retained.** `.iterations/<stage>.<n>.<ext>`, each with its own
  digest recorded in run state — so what round two actually said is verifiable,
  not merely remembered. `loom state show` reports the round count.
- **The bound is three.** The markdown pipeline's version of this loop was
  unbounded ("repeat until APPROVED"); the executor's has a number, because a
  loop it cannot bound is one it cannot safely run.
- **Exhausting the bound halts at `confirm-unresolved-review`.** Not a failure
  and not a silent pass: a human accepts the outstanding findings or stops the
  run, and that approval binds the artifacts like any other.
- A gate *before* the loop is approved once and does not re-halt each round.

## Run event timeline

Both pipelines append to `.claude/feature-workspace/<feature>/run-events.jsonl` — one JSON
object per line, written by this binary with timestamps taken from the clock at the moment
each transition happens. Stage durations are therefore *measured* by subtracting two
timestamps, not estimated by a model after the fact.

```bash
loom state timeline --spec docs/features/user-auth/spec.md
#  0s  run.started
#  0s  stage.started        analyst
#  4s  stage.completed      analyst
#  4s  gate.waiting         developer confirm-design
# 2m1s gate.approved        confirm-design tty
```

Recorded kinds: `run.started`, `run.completed`, `stage.started`, `stage.completed`,
`stage.failed`, `stage.interrupted`, `stage.stale`, `gate.waiting`, `gate.approved`.

- **Append-only.** The file is never rewritten or truncated; each event is one write. A
  process killed mid-write can leave a torn final line, which readers skip rather than
  failing on.
- **Not telemetry.** This is a local audit log. OpenTelemetry emission is roadmap L3.8, and
  this file does not replace `pipeline-trace.json`, whose `budgetUtilization` and iteration
  counts are still model-written estimates.

## Maturity level report

`loom health` ends with an agentic-maturity assessment driven by
`shared/levels.yaml`: it reports the highest level whose *entire* mechanical
evidence set passes, the passing evidence, and a concrete gap checklist for
the next level. Evidence is strictly mechanical — installed bundles on disk,
an MCP server that is configured in `.mcp.json` **and** actually answers
`tools/list` when spawned, a non-empty telemetry stream, policy files present.
Documentation-only bundles (`docsOnly` in `levels.yaml`) never count as
evidence, and a level whose enforcement bundles are all gated on unlanded
roadmap items is reported as not attainable yet rather than pretended into
reach.

```
Agentic maturity (shared/levels.yaml):
  Level 1 — Foundational prompts
  ✓ bundle "core-rules" installed
  ✓ bundle "agents" installed
  gaps to Level 2:
    ✗ MCP server not configured: read .mcp.json: no such file
    ✗ bundle "workflows" not fully installed (missing workflows)
```

## MCP server

`loom mcp serve` exposes the deterministic framework tools
(`analyze_complexity`, `check_accessibility`, `check_ubiquitous_language`,
`verify_dependencies`, `search_ki`, `search_docs`, `validate_artifact`)
over MCP stdio transport.
Structured JSON logs go to stderr, or to a file with `--log-file <path>` —
never stdout, which carries the MCP wire protocol. The server runs until the
client closes stdin or the process receives SIGINT.

Register it with Claude Code:

```bash
claude mcp add loom -- loom mcp serve
```

Or add it to `.mcp.json` (project) / `~/.claude/mcp.json` (global):

```json
{
  "mcpServers": {
    "loom": {
      "command": "loom",
      "args": ["mcp", "serve"],
      "env": {
        "AI_ASSISTANT_DOTFILES_PATH": "/absolute/path/to/loom-checkout"
      }
    }
  }
}
```

`AI_ASSISTANT_DOTFILES_PATH` is only required for `search_ki` (it points the
tool at the Knowledge Item corpus); the other five tools need no
configuration. See [shared/mcp/README.md](../../shared/mcp/README.md) for the
full tool reference and environment variables.
