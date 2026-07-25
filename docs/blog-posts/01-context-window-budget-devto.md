---
title: "Treat the Context Window Like a Budget, Not a Junk Drawer"
published: false
description: "What we learned building a context-engineering workflow for a 24-agent AI coding framework."
tags: ai, productivity, architecture, devtools
canonical_url:
cover_image:
---

Most AI coding workflows treat context as something that happens accidentally.

You open a few files. Paste a stack trace. Ask the model to inspect a directory. Then another. Then the
chat grows heavy, the model starts missing earlier instructions, and everyone pretends the problem is
"the model got weird."

In `ai-assistant-dot-files`, I wanted to treat the context window as a budget.

Not a vibe. Not a giant bucket. A budget.

The repo now ships a Context Engineering Framework that defines one canonical set of agents, skills, and
rules in `shared/`, then projects them into six AI coding tools: Claude Code, Cursor, Windsurf, GitHub
Copilot, Gemini/Antigravity, and OpenAI Codex. The current repo has 24 agents, 53 skills, 13 inter-agent
contracts, and 6 platform targets.

The core idea is simple: before an agent does serious work, another agent should decide what belongs in
the room.

## The context-engineer agent

The framework has a dedicated `context-engineer` agent. Its job is not to implement anything.

Its job is to produce a `context-manifest.md` before the rest of the pipeline starts.

That manifest scopes the bounded context, identifies relevant files, surfaces Knowledge Items and ADRs,
notes prior related deliveries, and estimates token budget pressure for the downstream agents.

This matters because the feature-delivery pipeline is not one prompt. It is a sequence:

- spec-writing
- product review
- context engineering
- analysis
- architecture
- performance review
- data review
- development
- code review
- accessibility review
- security review
- QA
- SRE review
- documentation
- DevOps

If the first few steps load the wrong material, every later agent pays for it.

The important design move is that the context manifest is itself a governed artifact. It has a contract in
`shared/contracts/context-manifest-contract.md`, and the `validate-artifact` skill checks that required
sections are present before the pipeline moves forward.

Context is not just "whatever the chat accumulated."

It is an explicit handoff.

## Context, memory, and learning are different problems

One of the most useful distinctions in the repo lives in `docs/runbooks/context-engineering.md`:

- Context is what is loaded into the model right now.
- Memory is durable knowledge that outlives the current run.
- Learning is a feedback loop that changes future behavior.

People often flatten all three into "RAG."

That loses important design pressure.

Context is a working set. It should be small, relevant, and current.

Memory is a durable corpus. It should be curated, searchable, and allowed to expire.

Learning is a loop. It should turn repeated delivery evidence into changed rules, changed prompts, or new
Knowledge Items.

The `context-engineer` reads memory to build better context, but it does not automatically rewrite memory.
That separation keeps a bad or noisy run from polluting the durable layer.

## Context decay

The framework also uses "context decay."

An artifact more than two pipeline phases old should be read as a summary, not in full. The target is a
small gist, roughly 200 words, produced through the `summarize-artifact` skill.

That is an intentionally boring mechanism, and that is why I like it.

Most context-window failures do not need a magic retrieval system. They need fewer stale artifacts loaded
verbatim.

If the developer is five phases downstream from the analyst, they probably need the current acceptance
criteria, edge cases, and constraints. They do not need every sentence of the analyst's intermediate
reasoning still floating in the model's attention.

## Platform reality changes the design

The repo does not pretend every AI coding tool has the same capabilities.

`docs/ARCHITECTURE.md` defines a capability tier system. Claude Code has full agent orchestration. Cursor
now has real `.cursor/agents/` and `.cursor/skills/` loading, but its rule files still need fully inlined
content. Windsurf and Copilot get persona/rule projections. Gemini/Antigravity reads `AGENTS.md` and has
confirmed skill invocation. Codex gets an inlined `.openai.md`.

That means the framework has to distinguish "agent" from "persona."

An agent can have tool access and participate in a multi-step process. A persona is a context frame: useful,
but not autonomous.

The practical result is a `shared/` canonical layer, plus generation and parity checks for each platform.
`scripts/check-parity.sh` exists because hand-copying instructions across tools is how drift wins.

## What I would copy into another project

If you do not need a 24-agent framework, I would still steal these ideas:

1. Write down the difference between Context, Memory, and Learning.
2. Add a pre-flight context manifest before complex work starts.
3. Treat old artifacts as summaries by default.
4. Make context handoffs structural, not conversational.
5. Add a parity check for any instruction copied across multiple tools.

The context window is not just a bigger prompt.

It is a scarce design surface.

Use it like one.

## Source trail

- `README.md` — canonical `shared/` layer, 24-agent roster, 53-skill catalog, six platform targets.
- `docs/ARCHITECTURE.md` — capability tiers, platform projection model, and context flow.
- `docs/runbooks/context-engineering.md` — Context/Memory/Learning distinction, context decay, manifest role.
- `docs/AGENT_REFERENCE.md` — `context-engineer` role and counterbalances.
- `docs/features/context-engineering-framework/TODO.md` — Epic 5 contract work and Epic 23 contract closure.

## Image prompt

Create a clean editorial hero image for a technical blog post titled "Treat the Context Window Like a
Budget." Show a calm, modern workspace where a glowing rectangular "context window" is divided into
carefully labeled budget categories: Rules, Task, Memory, Runtime, and History. Nearby, loose noisy
documents are being filtered out by a small mechanical sieve. Style: polished vector illustration,
dark-mode friendly, deep navy and teal palette, subtle warm highlights, no text except short category
labels, no logos, wide 16:9 composition.

