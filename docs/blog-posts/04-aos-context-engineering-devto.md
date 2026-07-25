---
title: "Toward an AI Operating System: Context Engineering as the First Runtime Primitive"
published: false
description: "Notes from building a governed agent framework, and why the next question is not more prompts but better runtime boundaries."
tags: ai, architecture, devtools, agents
canonical_url:
cover_image:
---

After publishing the first few posts about `ai-assistant-dot-files`, the obvious next question is:

Where does this go?

The tempting answer is "more agents."

I do not think that is the right answer.

The more interesting direction is an AI Operating System: not an OS in the kernel-and-device-driver sense,
but a runtime model for governed agentic work.

That distinction matters. The goal is not to wrap every task in bigger orchestration theater. The goal is
to ask what runtime primitives agentic systems actually need if they are going to be useful, inspectable,
and safe enough to trust with real software work.

The first primitive is context.

## Why context comes first

In the current framework, Context Engineering is already doing operating-system-shaped work.

The `context-engineer` agent builds a `context-manifest.md` before the delivery pipeline starts. That
manifest scopes relevant files, prior deliveries, Knowledge Items, ADRs, and token budget pressure.

That is not just prompt hygiene.

It is resource management.

An operating system decides what a process can see, which resources it can access, and how much room it has
to work before it starts corrupting other work. Context Engineering plays a similar role for agents.

It answers:

- What should this agent see?
- What should it not see?
- Which prior knowledge is relevant?
- Which stale artifact should be summarized instead of loaded in full?
- How much of the context budget should this phase consume?

That is why I keep coming back to the phrase: treat the context window like a budget, not a junk drawer.

If an AOS exists, context is one of its schedulable resources.

## The current system is not AOS yet

The repo today is a Context Engineering Framework.

It has one canonical `shared/` layer, 24 agents, 53 skills, inter-agent contracts, a memory lifecycle, and
platform projections for Claude Code, Cursor, Windsurf, GitHub Copilot, Gemini/Antigravity, and Codex.

That is real.

The AOS idea is earlier.

The notes in `docs/aos/AOS_Governance_Design_Pack.zip` describe design seeds: capability, governance,
learning, memory engineering, context engineering, and continuous improvement. They sketch pairs like:

- Context Engineer ↔ Context Auditor
- Memory Engineer ↔ Memory Auditor
- Prompt Architect ↔ Prompt Evaluator
- Orchestrator ↔ Scheduler
- Learning Engine ↔ Forgetting Engine
- Cost Optimizer ↔ Quality Optimizer

That is not a shipped runtime.

It is a question set.

And honestly, that is the part I trust most. A premature AOS would be easy to overbuild. A useful AOS has
to start by finding which governance gaps are real.

## Governance is the operating model

The strongest thing in the current framework is not the number of agents.

It is that each important handoff has some kind of counterbalance.

`docs/AGENT_REFERENCE.md` names four kinds:

- structural contracts
- downstream agent review
- human approval gates
- aggregate or delayed metrics

That model is small, but it changes the conversation.

Instead of asking "Can the agent do the task?" we ask "What checks the agent's work?"

That is the AOS-flavored question.

An operating system is not just a place where programs run. It is a place where programs run under rules:
permissions, scheduling, memory boundaries, process isolation, accounting, cleanup.

For agents, the equivalent rules are things like:

- Does this agent have the right context and only the right context?
- Is its output structurally valid?
- Is there an independent reviewer where judgment matters?
- Are irreversible actions gated by human approval?
- Are repeated failures converted into learning rather than buried in chat history?
- Is stale or duplicated knowledge expired instead of retrieved forever?

Those questions are less flashy than "autonomous swarm."

Good. Flash is usually where the bugs breed.

## Memory is the second runtime primitive

If context is what an agent sees now, memory is what survives the run.

The current memory model uses a promotion lifecycle:

`Capture -> Candidate -> Audit -> Approve -> Index -> Retrieve -> Expire`

That lifecycle is important because it refuses to treat memory as a pile of saved notes.

In an AOS model, memory needs governance the same way context does.

Not every observation deserves to become durable. Not every durable item deserves to live forever. Not every
retrieved item deserves to enter the active context window.

The AOS notes call out the pair:

Memory Engineer ↔ Memory Auditor

That feels right. One side promotes and organizes. The other side asks whether the memory is reusable,
non-duplicative, supported by evidence, and still true.

Memory without forgetting is just entropy with a search box.

## Entropy management might be the underrated subsystem

One of the AOS notes sketches an `Entropy Manager`.

Its job:

- remove duplicate knowledge
- detect stale docs
- detect unused agents
- merge overlapping rules
- reduce repository entropy

That is not glamorous.

It might be essential.

Any agent framework that learns will also accumulate. It will accumulate rules, prompts, skills, memories,
exceptions, platform quirks, and "temporary" workarounds that quietly become permanent.

Without an entropy-management function, learning turns into clutter.

The forgotten Cursor symlink story is a small example. The repo already had `.cursor/agents` and
`.cursor/skills` symlinked to `shared/`, but the decision was not documented, checked, or integrated into
the platform model. The fix was not only to make it work. The fix was to make it maintained.

That is AOS territory too.

Not "can the system do the thing once?"

"Can the system preserve the reason it does the thing?"

## What I do not want AOS to become

There are a few traps I want to avoid.

First: AOS should not become a cooler name for a giant prompt library.

Second: it should not pretend every tool has the same capabilities. The current framework already learned
that lesson through its platform tier system. Claude Code, Cursor, Copilot, Gemini, Windsurf, and Codex do
not expose the same runtime primitives.

Third: it should not optimize for maximum autonomy by default.

Autonomy without counterbalances is not maturity. It is just speed with a longer blast radius.

The goal is governed agency: more work can happen through agents because the system knows where to place
boundaries, reviews, summaries, approvals, and forgetting.

## The direction

If the current framework answers:

"How do we define agents, skills, rules, context, memory, and governance once, then project them into many
AI coding tools?"

Then AOS asks:

"What runtime model lets those agents operate with explicit context, memory, scheduling, permissions,
fitness functions, and entropy control?"

That is the north star.

But the path there should stay grounded:

1. Improve Context Engineering first.
2. Add Context Auditing before adding more autonomy.
3. Keep Memory Engineering evidence-based and approval-driven.
4. Track fitness functions like context precision, retrieval quality, token efficiency, memory quality, and entropy.
5. Treat every new subsystem as guilty until it proves it closes a real gap.

The point is not to build an AI Operating System because the metaphor sounds good.

The point is to discover which parts of software delivery become safer and more comprehensible when agents
run inside a governed runtime instead of a chat transcript.

That is the direction I want to explore.

## Source trail

- `README.md` — current framework shape: canonical `shared/` layer, 24 agents, 53 skills, six platform targets.
- `docs/runbooks/context-engineering.md` — Context/Memory/Learning distinction and context manifest role.
- `docs/runbooks/memory-engineering.md` — memory promotion lifecycle.
- `docs/AGENT_REFERENCE.md` — counterbalance model for every current agent.
- `docs/aos/AOS_Governance_Design_Pack.zip` — exploratory AOS notes, especially `00-AOS-Vision.md`, `01-Governance-Checks-and-Balances.md`, `02-Context-Governance.md`, `08-Fitness-Functions.md`, and `09-Entropy-Manager.md`.
- `docs/features/context-engineering-framework/TODO.md` — Epic 30 forgotten symlink story and the concrete drift-prevention fix.

## Image prompt

Create a polished editorial hero image for a technical blog post titled "Toward an AI Operating System."
Show a futuristic but grounded control room where "Context" is the central resource dashboard: token
budgets, memory cards, agent processes, review gates, and entropy cleanup are arranged like operating
system panels. The image should feel like governed infrastructure rather than sci-fi chaos. Style: modern
vector illustration, dark navy background, teal and amber highlights, subtle grid, no real logos, no dense
text, wide 16:9 composition.

