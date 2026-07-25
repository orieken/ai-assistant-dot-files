---
title: "Handoff Prompts as First-Class Artifacts: Surviving Context Exhaustion"
published: false
description: "Why we check task handoff prompts into version control—and how self-contained prompts unlock asynchronous multi-agent collaboration."
tags: ai, productivity, devtools, architecture, software-engineering
canonical_url:
cover_image:
---

Long-running AI chat sessions inevitably degrade.

You start a session with a big feature request. By step 14, the model's context window is bloated with past terminal outputs, intermediate code attempts, and stale stack traces. The model starts forgetting early constraints, hallucinating file paths, and making sloppy edits.

The naive reaction is to keep pushing through the heavy session, fighting prompt decay until the chat crashes.

In `ai-assistant-dot-files`, we established a pattern that completely eliminates context exhaustion: **Handoff Prompts as First-Class Artifacts**.

Instead of trying to finish a massive feature in a single degrading conversation, we check the next unit of work into version control as a self-contained markdown file under `docs/prompts/`.

---

## Anatomy of a First-Class Handoff Prompt

A handoff prompt is not a vague 2-line instruction like "Implement the login page."

It is a fully briefed, self-contained design document checked into git (`docs/prompts/<task-slug>.md`) that allows a fresh AI agent—or a fresh human developer—to pick up the work cold without needing prior conversation memory.

Every handoff prompt in our framework enforces 6 mandatory sections:

```markdown
# [Task Title]

## Target Repo
Absolute path to the target repository. Explicit instruction on commit landing.

## Context to Read First
Key ADRs, specifications, and prior delivery manifests to read BEFORE editing code.

## Scope of Implementation
Numbered, explicit list of tasks ordered logically by dependency.

## Guardrails & Commit Discipline
- Conventional Commit message format.
- Staging rules (never `git add -A`; stage explicit paths only).
- Do NOT push to remote.

## Escalation Criteria
Explicit conditions that require halting execution to ask a human operator.

## Report Format
Structured summary (<200 words) to output upon completion.
```

---

## Why Version-Controlled Prompts Beat Heavy Chat Sessions

### 1. Zero Context Overhead
When a fresh agent opens a handoff prompt, 100% of its context window is clean. There are no stale error tracebacks or old code attempts cluttering attention. The agent executes with maximum reasoning capability.

### 2. Async Multi-Agent Orchestration
By checking prompts into `docs/prompts/`, different agents (or subagents) can work asynchronously. Agent A completes Phase 1, drafts the Phase 2 handoff prompt in `docs/prompts/`, and exits. Agent B picks up the Phase 2 prompt in a clean context and executes it flawlessly.

### 3. Reusability and Auditability
Because handoff prompts live in `docs/prompts/` (and move to `docs/prompts/done/` upon shipping), your git history records not just *what* code changed, but *the exact instructions* that drove the change.

---

## The Handoff Lifecycle

1. **Draft**: Author `docs/prompts/<task-slug>.md`.
2. **Index**: Add entry to Active Prompts table in `docs/prompts/README.md`.
3. **Execute**: Launch a fresh chat/subagent passing the prompt file.
4. **Ship & Relocate**: Move prompt to `docs/prompts/done/<task-slug>.md` and update the index.

---

## Takeaway: Treat Prompts Like Code

Prompts are code. They belong in your repository, under version control, structured with clear interfaces and explicit guardrails.

When you elevate handoff prompts to first-class artifacts, context exhaustion disappears—and multi-agent development becomes deterministic.

---

## Image Prompt

> Hero image prompt: A modern architectural concept showing a glowing crystal baton labeled "Handoff Prompt" being passed seamlessly between two robotic arm nodes across a clean digital grid line. Cyberpunk dark mode background with cyan and violet light trails.
