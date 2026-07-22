# AOS Migration Guide

Companion doc to `docs/aos/migration-plan.md`. That file is the
architectural plan; this file is the operator-facing "what do I actually
do" — one section per phase, tracking what's live and what's still to come.

**Status at v3.0.0**: Phase 1 shipped. Nothing forces adoption yet. Every
AOS layer is opt-in, and the default install path is unchanged from v2.x.
Later phases will add more layers and eventually a `--full` install mode.

---

## Upgrading from v2.x

**Nothing to do. All AOS additions are opt-in.**

A v2.x install upgraded to v3.0.0 without invoking any new AOS-specific
capability behaves identically to v2.x. Every agent's version, every skill's
behavior, every rule, every contract, and every blueprint stays the same.
`deliver-feature` still orchestrates the same pipeline in the same order and
writes the same artifacts to the same paths. `validate-artifact` still runs
the same structural checks. `health-check` still passes on the same criteria
(it now also prints an "AOS Layers" inventory at the bottom, but that section
never fails — presence or absence of AOS layers is informational).

You may re-run `install.sh` after pulling v3.0.0 to pick up the new files
in `shared/telemetry/`, `shared/evaluation/`, and the `memory-auditor` agent —
they'll be symlinked/copied into your `.claude/` (or `.cursor/`, or platform
equivalent) alongside everything you already had. If you skip re-running
`install.sh`, everything you already had keeps working; you just don't have
the new files visible in your platform config until you refresh.

---

## Opting into telemetry

**Not yet — Phase 2 wires it.**

The `shared/telemetry/` layer landed in v3.0.0 with the schema (`event-schema.md`),
the recorder skill (`event-recorder.md`), and the layer overview (`README.md`).
But no producer in the framework emits events yet — `deliver-feature`,
`validate-artifact`, and every existing agent still behave exactly as they did
in v2.x. There is nothing to opt into today because there is nothing to opt
into against.

Phase 2 will add the `shared/hooks/` layer, at which point telemetry-emission
becomes a hook-driven, per-project opt-in. Phase 3 wires the continuous
evaluation triggers that consume the emitted events.

If you want to see the design ahead of time:
- `shared/telemetry/README.md` — layer overview + retention convention
- `shared/telemetry/event-schema.md` — the event types and fields
- `shared/telemetry/event-recorder.md` — the append-only recorder skill
- `docs/aos/migration-plan.md` — Phase 2 and 3 sections

---

## Opting into memory-auditor

**Available now — invoke on demand.**

`memory-auditor` is the first AOS counter agent, paired with the existing
`memory-engineer` skill. It is read-only (tools: `Read, Glob, Grep`) and
never modifies KIs — it produces a findings report you or `memory-engineer`
then act on with explicit approval.

To invoke it in Claude Code:

```
> Use the memory-auditor agent to audit shared/knowledge/ and .claude/knowledge/.
```

Or, if your platform supports slash-shortcut invocation of subagents,
whatever the platform-specific spelling is (`@memory-auditor`, etc. — see
your platform's own docs).

What it checks:

- **Schema compliance** — every KI has `name`, `tags`, `domain`, `created`
  in its frontmatter, matching the format in `shared/knowledge/README.md`.
- **Exact duplicates** — two or more KIs share the same frontmatter `name:`.
- **Semantic-overlap candidates** — pairs of KIs whose bodies cover
  substantially the same subject in different words. These are surfaced for
  human/memory-engineer judgment — the auditor never merges anything.
- **Stale-metadata candidates** — KIs older than 6 months with zero
  references anywhere in the corpus. Also surfaced for human judgment —
  stale is not the same as wrong.

Findings default to stdout. To write to a file instead:

```
> Use the memory-auditor agent and save the report to .claude/audits/memory-audit-2026-07-22.md.
```

`.claude/audits/` is created if missing. Absence of this directory in your
project is fine — the auditor creates it on first file-write.

### Difference vs `memory-engineer`

- `memory-engineer` (skill under `shared/skills/memory-engineer/`) — the
  producer-shaped curator. Can propose merges, cross-references, registry
  edits, and expiration decisions. Findings still require human approval;
  it does not act unilaterally.
- `memory-auditor` (agent under `shared/agents/memory-auditor.md`) — the
  read-only counter. Reports what's wrong; never proposes edits. Cheap
  enough to invoke as a scheduled or hook-driven check without worrying
  about unwanted edit recommendations.

Run both when convenient. When both have run recently, treat the
`memory-engineer` sweep as the more comprehensive judgment; the auditor
is the fast, hook-friendly first pass.

---

## What's coming next

**Phase 2 (v3.1)** — see `docs/aos/migration-plan.md`, Phase 2 section:

- 10 more counter agents (`context-auditor`, `knowledge-auditor`,
  `prompt-evaluator`, `agent-evaluator`, `rule-auditor`, `pattern-reviewer`,
  `tool-validator`, `documentation-auditor`, `retrieval-evaluator`,
  `privacy-auditor`), all following the `memory-auditor` shape.
- 4 opposing-force skill pairs (Memory Expansion / Compression, Learning /
  Forgetting, Cost / Quality, Scheduler).
- `shared/hooks/` layer with an event → skill/agent config schema and
  example hooks.
- `validate-artifact` gets an optional post-structural-check invocation of
  the corresponding counter agent (opt-in via config; default remains
  structural only).

**Phase 3 (v3.2)** — orchestration runtime, RAG layer, Learning/Forgetting
engines, and the Trinity-native workflow refactor of `deliver-feature` and
`test-driven-developer`. This is also where telemetry emission gets wired
and `install.sh` gains `--base` (default, current behavior) and `--full`
modes.

**Phase 4 (v3.3)** — policy layer for auto-approval on documented safe
paths (doc-only changes, pure refactors, test additions).

Each phase ships independently and is independently revertable. Each phase
holds the same backward-compat guarantee: install it without opting in,
and nothing changes from the phase before.

---

## Related

- `docs/aos/migration-plan.md` — the architectural plan this guide operationalizes
- `docs/aos/AOS_Governance_Design_Pack/` — the underlying vision and 15
  governance pairs
- `shared/agents/CHANGELOG.md` — the v3.0.0 entry and every future phase entry

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
