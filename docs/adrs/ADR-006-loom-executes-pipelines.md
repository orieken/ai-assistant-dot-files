# ADR-006: Loom Executes Pipelines

## Status

Accepted

## Date

2026-08-29

## Context

The 2026-08-29 architectural audit (`docs/roadmaps/architectural-audit-2026-08-29.md`) reached a
blunt verdict: `loom` today is a prompt distribution system, not an agentic framework. Of ~8.9k
lines of Go, the executable surface is a Cobra file-installer and an MCP tool server. Every
runtime subsystem the documentation describes — orchestration, checkpointing, retry governance,
telemetry recording, policy evaluation, retrieval tiers — is a markdown file instructing a host
platform's LLM to role-play that subsystem. The ambiguity is load-bearing: two prior roadmaps
(`agy.md`, the maturity TODO) independently proposed "build a kernel" without either committing
to it, and twenty roadmap items target an `internal/orchestrator/` package that does not exist.

Roadmap item **M0.1** requires this question to be answered explicitly rather than drifted past:

> Does `loom` **execute** pipelines, or does it only **validate and distribute** content that a
> host runtime executes?

Both are defensible products. Only one can be the roadmap.

## Decision

**Loom executes pipelines.** A minimal Go executor (`internal/orchestrator/`, roadmap item M0.4)
owns the run loop: load a plan, execute stages in order via a provider adapter, persist durable
state, and stop. Gate enforcement (L2.13), typed inter-stage state (L2.9), policy evaluation
(L2.16), and telemetry emission (L3.8) subsequently attach to that executor rather than living
as prose.

Two boundaries fix the shape of this decision:

1. **Host platforms become providers, not the runtime.** Claude Code, Cursor, and every other
   installed platform are one way to invoke an agent — a `Provider` implementation the executor
   calls — never the owner of run-loop semantics, retry ceilings, or gate enforcement.
2. **The MCP relationship is client-side only** (fixed by epic 75's distribution strategy): the
   orchestration kernel is `loom run` acting as an MCP *client*. The pipeline itself is never
   exposed as an MCP tool someone else calls.

## Alternative Considered: Validate and Distribute

This alternative deserves a serious statement, because it is a legitimate product — arguably the
product loom already is. In this model, loom's Go surface stays what it is today: an installer,
a health checker, a parity checker, and an MCP server exposing read-only analysis tools. The
markdown content — agents, skills, rules, contracts — is the deliverable, and host platforms
(Claude Code first among them) supply all execution semantics. Loom's value proposition becomes
content quality, cross-platform sync, and validation depth, and the honest comparison set is
prompt-library and dotfiles-manager tooling, not agent frameworks.

Choosing it would delete most of the build roadmap. Everything in the KERNEL workstream goes:
the executor skeleton (M0.4), typed graph state (L2.9), executor-owned pipeline state (L2.12),
gates as process interrupts (L2.13), gate-reset enforcement (L2.14), real resume (L2.15), the
deterministic policy evaluator (L2.16), the planner/router (L3.1), real parallelism (L3.3), the
hook executor (L3.10), and effectively all of Milestone 4 (bounded Reflexion, budget governor,
eval-gated prompt promotion, MCP client runtime). Telemetry (L3.8) shrinks to instrumenting the
MCP server only, since there is no run loop to trace. The roadmap reduces to content quality,
tool hardening (L2.1–L2.8), and distribution — and the documentation must stop using the words
"orchestration," "checkpoint," and "governed pipeline" for behavior loom cannot observe.

The reason it is not chosen: every guarantee the framework sells — retry ceilings, checkpoint
integrity, gate enforcement, auditability — degrades to a suggestion when the enforcement
mechanism is prose interpreted by the same model being governed. The audit's H1 (unbounded
review loop), H6 (incompatible telemetry events), and H7 (LLM-interpreted approval policy) are
all instances of one root cause: no process owns the loop. Validate-and-distribute cannot close
them, only document them.

## Consequences

- Milestones 1–4 of `BUILD-ROADMAP.md` become real work items rather than aspiration; M0.4 (the
  executor skeleton) is unblocked once this ADR is accepted.
- Host platforms are demoted from "the runtime" to one `Provider` among several. The first
  providers are `claude` (spawning `claude -p` headless) and a deterministic `mock` for tests.
- Run-loop semantics — stage ordering, state persistence, resume, timeouts, and eventually
  gates — move from markdown instructions into Go, where they can be tested and can fail.
- The existing markdown pipeline (`deliver-feature` as a Tier-1 skill) keeps working unchanged
  during the transition; the executor ships running the same linear agent sequence as its
  built-in default plan, so behavior is preserved while the substrate changes underneath.
- `README.md` and `docs/ARCHITECTURE.md` must describe unshipped runtime subsystems with honest
  status markers ("specified, ships with LX.Y") rather than in the present tense — done in the
  same change that proposes this ADR.
- If this ADR is instead rejected in favor of validate-and-distribute, epic 76 ends at Phase A
  and `BUILD-ROADMAP.md` is rewritten as a content-quality roadmap.

### Verification

Fitness functions arrive with the implementation items, not this ADR: M0.2 puts the Go build,
tests, and lint in CI; M0.4's done-when criterion (`loom run` executes ≥3 real stages, persists
state, resumes after SIGINT) is the first executable proof of the decision. Until M0.4 lands,
this decision is **judgment-only**: its observable artifact is that `README.md` and
`docs/ARCHITECTURE.md` no longer claim unimplemented subsystems in the present tense.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context
Engineering Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy
or adapt this file, please keep this attribution.*
