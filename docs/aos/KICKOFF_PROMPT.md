# Kickoff Prompt: AOS (v3) Design Conversation

Copy everything below the line into a fresh Claude Code session, on the `v3-aos` branch of this repo.

---

```
We're starting a design conversation for "AOS" -- a possible v3 evolution of the ai-assistant-dot-files
repo. Don't write code or make repo changes yet; this is a scoping and architecture conversation first.

## What already exists (v2, shipped, on `main`)

This repo already has a fully working "Context Engineering Framework": a canonical shared/ layer (24
agents, 47+ skills, rules) generated or symlinked into 6 AI coding tools (Claude Code, Cursor, Windsurf,
GitHub Copilot, Gemini/Antigravity, OpenAI Codex), a 14-agent feature-delivery pipeline
(shared/skills/deliver-feature/SKILL.md) with inter-agent contracts (shared/contracts/) validated at every
handoff, agent versioning + changelog discipline, pipeline observability (traces, scorecards,
retrospectives), and -- most relevant here -- a working Memory Engineering system: shared/memory-registry.json,
docs/runbooks/memory-engineering.md (Capture -> Candidate -> Audit -> Approve -> Index -> Retrieve -> Expire
lifecycle), and three skills (memory-engineer, promote-memory, query-memory).

Read these first to understand current state before proposing anything new:
1. docs/AGENT_REFERENCE.md -- every agent's Role, Counterbalance, and explicit Gap (this already covers a
   lot of what "governance checks and balances" might mean -- read it before assuming that concept is missing)
2. docs/runbooks/memory-engineering.md -- the memory lifecycle already built
3. docs/features/context-engineering-framework/TODO.md -- Phase 9 section specifically, the most recent
   work (Memory Engineering, contract closure, this AGENT_REFERENCE.md doc, Team Topology, ci-check.sh)
4. docs/runbooks/context-engineering.md -- the Context/Memory/Learning mental model everything else builds on

## What this branch contains (v3 prototyping material, not yet designed)

docs/aos/ has three sets of raw, terse seed notes -- treat these as a starting point to interrogate, not a
spec to implement as-is:

- AOS_v1_Future_Architecture_Suggestions.md -- an OS-metaphor architecture: Kernel (agent lifecycle,
  contracts, orchestration, event bus, state transitions, handoffs), Memory Manager, Scheduler (parallel
  execution, retries, priorities), Security Manager, Package Manager (agent/skill/rule distribution and
  versioning), Telemetry, Knowledge System, Plugin System (LightRAG/GraphRAG/MCP/vector stores), Developer
  Experience.
- AOS_Governance_Design_Pack/ (9 short files) -- the more concrete of the two:
  - A "governance checks and balances" pattern: every builder role paired with an auditor/evaluator role
    (Context Engineer <-> Context Auditor, Memory Engineer <-> Memory Auditor, Knowledge Curator <-> Knowledge
    Auditor, Prompt Architect <-> Prompt Evaluator, Agent Designer <-> Agent Evaluator, Rule Author <-> Rule
    Auditor, Pattern Author <-> Pattern Reviewer, Tool Builder <-> Tool Validator, Documentation Writer <->
    Documentation Auditor, Retrieval Engine <-> Retrieval Evaluator, Memory Expansion <-> Memory Compression,
    Learning Engine <-> Forgetting Engine, Cost Optimizer <-> Quality Optimizer)
  - A suggested `.project-ai/` directory layout (agents/skills/rules/hooks/contracts/context/memory/rag/
    orchestration/evaluation/telemetry/prompts)
  - New fitness-function categories: context precision, retrieval quality, token efficiency, memory quality,
    architecture health, entropy
  - An "Entropy Manager" concept: whole-repo hygiene -- dedup knowledge, detect stale docs, detect unused
    agents, merge overlapping rules
- memory_engineering_prompts/ (10 files) -- this is the *original* seed material v2's actual Memory
  Engineering epic was built from. It has already been implemented and shipped in v2 (see above). Treat this
  folder as "done, superseded" reference material, not an open to-do list -- don't rebuild what's already
  there.

## A decision already made once, deferred to here

Mid-v2, while writing docs/AGENT_REFERENCE.md, the "counterbalance" pairing pattern from the Governance
Design Pack came up directly. The explicit choice at the time was to *document* each agent's existing
counterbalance (contract, downstream reviewer, approval gate, or delayed metric) rather than *build* new
dedicated counterbalance roles/agents -- and to defer the "should we build dedicated auditor roles" question
to this v3 conversation instead of deciding it as a v2 side-quest. This is that conversation.

## What I want out of this session

Not code. A scoped analysis of:
1. Given what v2 already covers (see AGENT_REFERENCE.md's existing counterbalance categorization), which of
   the Governance Design Pack's 13 role-pairs represent a real, currently-uncovered gap versus which are
   already handled by an existing contract/reviewer/gate/metric under a different name?
2. Is AOS a v3 evolution of *this* repo, or does the OS metaphor (Kernel/Scheduler/Security Manager/Package
   Manager) imply a genuinely different kind of project that shouldn't live inside ai-assistant-dot-files at
   all?
3. If it's worth pursuing here: what's the smallest real slice worth building first, and does it belong on
   this `v3-aos` branch as continued iteration, or as a fresh spec via the existing `new-feature`/
   `deliver-feature` pipeline once scoped?

Do all work on the `v3-aos` branch (or a new branch off it) -- don't touch `main`, which is the shipped,
CI-green v2 framework.
```

---
*Prototyping material for a possible v3 evolution of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework. Not part of the shipped v2 framework on `main`.*
