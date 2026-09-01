# ADR-007: Where UI Tests Live, and What They Hand the Pipeline

## Status

Proposed

## Date

2026-08-31

## Context

`visual-qa-engineer` consumes two things a UI test run produces: interaction heatmaps
(`heatmap-data/*.json`, written by `@orieken/saturday-playwright-heatmap` during the
qa-engineer's run) and Playwright screenshot baselines. Its own preconditions say it runs
"when heatmap instrumentation or Playwright visual snapshots are present" — a check for
directories on disk, evaluated by a model at the moment the stage would run.

Three things have made that arrangement's ambiguity load-bearing.

**The executor now routes.** Roadmap item L3.0 (epic 81) computes which stages run from typed
`AnalysisState` before the design gate. Every other conditional stage has an analysis-derived
predicate — crossings, migrations, dependencies, SLAs. `visual-qa-engineer` does not: whether
heatmap data exists is a fact about the *environment*, not about the feature. No predicate over
the analysis can decide it, so the router either always includes the stage or needs a different
kind of fact to read.

**A manifest would change the stage's shape.** `saturday-playwright-heatmap`'s `scanner.ts`
already enumerates every visible interactable element on a page with a selector resolved in a
stability order (`#id` → `[data-testid]` → `tag.classes` → bare tag). Persisted per route rather
than per scenario, that becomes a manifest a UI-testing agent can *plan* against instead of only
measuring after the fact. But a manifest is a file in one world and a published, versioned
artifact in another — and which one it is depends entirely on where the tests live.

**Saturday suites frequently live in their own repository.** The pattern this framework
documents — Site-Centric `BaseSite`/`BasePage`/`BaseFlow`, Cucumber scenarios as business
language — is often owned by a QA group and versioned separately from the application it
exercises. `saturday-monorepo` itself is a separate repository from any app under test.

The question this ADR settles is therefore not a preference about repository layout. It is:
**what does the delivery pipeline assume about UI tests, and what contract carries their output
back to it?**

## Decision

*(Proposed — the two options below are the decision to be made. This section records the
recommendation; it is not Accepted until a human says so.)*

**Recommended: support both, and make the artifact contract the only thing the pipeline depends
on.**

The pipeline must not assume it can run the UI suite. It must assume only that a **UI evidence
bundle** is available, produced by whoever owns the tests:

- `manifest.json` — interactable elements per route, each with its selector and the selector's
  stability tier (`id` | `testid` | `class` | `tag`), plus the app version or commit it was
  scanned against.
- `coverage.json` — which manifest entries were exercised, derived from heatmap capture.
- `baselines/` — screenshot baselines with the same version key.

Where the tests live then becomes a *sourcing* question with two supported answers:

1. **Co-located** — the bundle is produced in-repo during the qa-engineer's run, and
   `visual-qa-engineer` reads it directly from the workspace. This is the default and the only
   path the executor can currently run end to end.
2. **Separate repository** — the bundle is published as a versioned artifact by the test repo's
   own CI and fetched by version key. The delivery pipeline never runs the suite; it consumes
   evidence and refuses to consume evidence whose version key does not match what was built.

`visual-qa-engineer`'s precondition becomes "a UI evidence bundle for this version is available",
which is a checkable fact in both worlds — and one the router can read.

## Consequences

**What this buys.**

- The routing problem is resolved without special-casing: the stage's condition becomes a fact
  about an artifact, the same kind of fact every other stage's condition already is.
- A stale bundle becomes detectable rather than silently trusted. The version key is what makes
  "these baselines describe a different build" a refusable condition, and it composes with the
  digest machinery from L2.12 — a bundle in the workspace is an artifact like any other.
- The manifest gains a second use beyond speed: the stability tier of each selector is a
  testability signal (`"37% of interactables on /checkout have no stable selector"`), which
  overlaps with what `accessibility-engineer` already looks for.

**What it costs.**

- A bundle format is a contract, and contracts have to be versioned and documented. This is real
  work — it is a schema, a producer, and a consumer, not a convention.
- The separate-repository path cannot be enforced by the executor. Fetching an artifact by
  version key is a supply-chain step, and nothing in the current gate machinery covers "is this
  bundle the one CI produced". Treat that path as documented-but-unenforced until it is.
- Two supported topologies is more surface than one. The mitigation is that only the *sourcing*
  differs; everything downstream reads one format.

**What is explicitly not decided here.**

- Whether `loom` gains a manifest-generation tool, and whether it wraps `saturday-playwright-heatmap`
  or reimplements the scan. That is an epic, and it should follow this ADR rather than precede it.
- Whether the bundle is fetched by `loom` or by the project's CI before `loom run` starts.
- Anything about non-UI test evidence. API and unit results have their own reporting path
  (`shared/rules/testing-conventions.md`) and are out of scope.

## Alternatives Considered

**Assume co-located tests only.** Simplest, and it matches what the executor can run today: one
topology, no bundle format, `visual-qa-engineer` keeps reading directories. Rejected because it
makes the framework's own documented Saturday pattern unsupported — a Saturday suite in its own
repository would have no path into the pipeline at all, and the framework would be recommending a
structure it cannot deliver against.

**Assume separate repositories only.** Forces the artifact contract from day one and keeps the
pipeline honest about not running the suite. Rejected as over-general: most projects adopting this
framework start with tests beside their source, and requiring a publish-and-fetch cycle for a
directory that is already on disk is ceremony that would get worked around.

**Leave it undecided and let each project figure it out.** The status quo. Rejected because it is
what produced the current situation: a precondition that reads the filesystem, a stage the router
cannot reason about, and no way to tell a stale baseline from a current one.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
