---
title: "Memory Engineering Is a Promotion Pipeline, Not a Pile of Notes"
published: false
description: "How a small AI tooling framework handles durable memory without turning every lesson into permanent clutter."
tags: ai, architecture, knowledge-management, devtools
canonical_url:
cover_image:
---

A lot of AI memory systems start with the same temptation:

"Just save the useful thing."

That sounds harmless until the knowledge base becomes a junk drawer. Half the notes are too specific, a few
are duplicates, some are obsolete, and nobody knows which ones the agent should trust.

In `ai-assistant-dot-files`, the memory system is deliberately slower.

It uses a promotion lifecycle:

`Capture -> Candidate -> Audit -> Approve -> Index -> Retrieve -> Expire`

That lifecycle is documented in `docs/runbooks/memory-engineering.md`, and the important word is not
"capture."

It is "candidate."

## Nothing writes directly to memory

The framework has a durable memory layer: Knowledge Items in `shared/knowledge/`, ADRs in `docs/adrs/`,
the domain dictionary, team topology, a feature archive, and a registry at `shared/memory-registry.json`.

But a lesson from a delivery does not jump straight into `shared/knowledge/`.

It first becomes a Candidate Record.

That record has required fields:

- Source
- Type
- Evidence
- Tags
- Expiration condition

Then `memory-engineer` audits it:

- Is it reusable?
- Is it already covered?
- Is it too speculative?
- Does it belong as a Knowledge Item, or should it become a rule change, prompt edit, or ADR instead?

Only after that does a human approve the destination.

The design is intentionally similar to code review. Durable memory changes future behavior, so they deserve
a paper trail.

## Rejection is a feature

One of my favorite parts of the memory runbook is that it has explicit rejection rules.

Do not promote a memory when it is:

- a one-off
- already covered
- too speculative

That makes "zero candidates promoted this cycle" a healthy result, not a failure.

This is where memory engineering starts to look less like note-taking and more like gardening. The point is
not to preserve every leaf. The point is to keep the soil useful.

## Expiration matters

The lifecycle also includes expiration.

A Knowledge Item can become stale when the underlying code, agent, or pattern changes. It can be superseded
by a better KI. Or usage analytics can show that it never appears in context manifests, which may mean it is
not useful or just badly tagged.

The repo does not delete those blindly. Expired KIs move to `shared/knowledge/expired/` with a note.

That choice matters because a wrong memory is still evidence. It tells you what the team used to believe.

## Why not build the big retrieval system now?

There is a runbook for LightRAG integration at `docs/runbooks/lightrag-integration.md`.

There is intentionally no implementation.

That is not an omission. It is a YAGNI decision.

The current retrieval path is smaller: `search-ki` searches Knowledge Items and ADRs; `query-memory` works
across the broader memory registry. The repo currently has 4 portable Knowledge Items, so building a bigger
retrieval subsystem before the corpus needs it would add moving parts without solving an observed bottleneck.

The runbook exists so the future integration has a shape if the need becomes real.

That is the kind of "not yet" I trust: documented, intentional, and reversible.

## Governance by design, not vibes

The memory system fits into a larger governance model.

`docs/AGENT_REFERENCE.md` lists every one of the 24 agents and names what checks its work: a structural
contract, a downstream reviewer, a human approval gate, or an aggregate metric.

Some gaps are real and stated plainly. For example, `test-driven-developer` deliberately bypasses the full
review chain for speed. The doc does not pretend otherwise.

That same honesty shows up in memory.

The system does not claim every remembered thing is true forever.

It asks:

- Where did this come from?
- What evidence supports it?
- Who approved it?
- What would make it expire?

Those are small questions, but they change the shape of the system.

## What I would copy into another project

If you are adding memory to an AI workflow, I would start here:

1. Do not let agents write directly to durable memory.
2. Make every memory a candidate first.
3. Require evidence and an expiration condition.
4. Treat rejection as a valid outcome.
5. Periodically compress duplicates instead of making retrieval disambiguate them forever.

The hard part of memory is not remembering.

It is staying worth remembering.

## Source trail

- `docs/runbooks/memory-engineering.md` — full memory lifecycle and Candidate/Audit fields.
- `docs/runbooks/context-engineering.md` — distinction between Context, Memory, and Learning.
- `shared/memory-registry.json` — registry of durable memory sources and retrieval backends.
- `docs/AGENT_REFERENCE.md` — agent counterbalances and explicit gaps.
- `docs/features/context-engineering-framework/TODO.md` — Epic 22 memory engineering and Epic 26 documentation-manager boundary cleanup.
- `docs/runbooks/lightrag-integration.md` — documented future path, intentionally not implemented yet.

## Image prompt

Create a polished editorial hero image for a technical blog post titled "Memory Engineering Is a Promotion
Pipeline." Show a careful conveyor-belt workflow with seven labeled stations represented by icons:
Capture, Candidate, Audit, Approve, Index, Retrieve, Expire. Notes enter as rough scraps and leave as
organized glowing cards in a small knowledge garden. Style: modern vector illustration, muted indigo,
green, and amber palette, calm and precise, no brand logos, wide 16:9 composition, minimal readable text.

