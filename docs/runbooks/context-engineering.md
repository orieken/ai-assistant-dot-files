# Context Engineering Guide for AI Agents

This runbook outlines the core principles and techniques of **Context Engineering** in our multi-agent software development ecosystem. It serves as a guide for both human developers and autonomous agents to manage context windows, maximize the Signal-to-Noise Ratio (SNR), and prevent context drift.

---

## When to Use
- Before launching a new agent pipeline or feature task.
- When an agent becomes slow, starts repeating commands, or shows signs of "context drift" (forgetting guidelines or earlier instructions).
- To optimize token consumption and reduce execution latency.

---

## The Context Challenge
Large Language Models (LLMs) operate on a finite context window. While modern models have large windows (e.g., 200k+ tokens), their reasoning accuracy and instruction-following capability degrade as the window fills with noise. 
- **Context Drift**: As conversation history grows, earlier instructions, architecture rules, and specs lose influence over the model's responses.
- **Signal-to-Noise Ratio (SNR)**: Unrelated files, redundant logs, and full directory listings reduce the density of relevant information, causing errors.

---

## Context Taxonomy
Context is categorized into six dimensions, ordered from permanent/static to ephemeral/dynamic:

| Layer | Name | Source | Lifetime | Purpose |
|---|---|---|---|---|
| **1** | **System Context** | `CLAUDE.md`, `.cursorrules` | Session-long | Sets the agent's identity, role, and rules of engagement. |
| **2** | **Rule Context** | `ARCHITECTURE_RULES.md` | Session-long | Defines core technical laws (TDD, Clean Architecture). |
| **3** | **Knowledge Context** | Knowledge Items (KIs), ADRs | Demand-driven | Cached repository-specific context (known bugs, APIs). |
| **4** | **Task/Goal Context** | Feature spec (`features/*.md`), plan | Task-long | Outlines target requirements and acceptance criteria. |
| **5** | **Historical Context** | Thread history, `docs/features/` | Ephemeral | Tracks what has been done in the current run/past runs. |
| **6** | **Runtime Context** | Open files, cursor, compiler/lints | Real-time | Code files, line highlights, and tool outputs. |

---

## Core Principles of Context Engineering

### 1. The Principle of Least Context (Least Privilege)
Agents must only load the information necessary for their immediate sub-task.
- **Do NOT** read entire source directories recursively if you only need to understand an interface.
- **Do NOT** read a whole 1000-line file if you only need to modify a 50-line class. Use range-based file reads (`view_file` with `StartLine`/`EndLine`).
- **Do NOT** keep files open that are unrelated to the current work.

### 2. Proactive RAG (Check Knowledge Items First)
Before performing independent analysis or proposing a design:
- Look up existing **Knowledge Items (KIs)** under the local `<appDataDir>/knowledge/` directory.
- Verify whether the pattern, bug, or configuration has already been solved or documented.
- Check **Architecture Decision Records (ADRs)** to understand why a technology or pattern was selected.

### 3. State Externalization (Offload Reasoning Memory)
Do not force the agent to track long checklists and detailed states inside the active conversation context.
- Maintain a separate [task.md](file:///Users/oscarrieken/.gemini/antigravity-ide/brain/e9680908-8c1b-480d-8c4d-6ab41d73a1fc/task.md) and [implementation_plan.md](file:///Users/oscarrieken/.gemini/antigravity-ide/brain/e9680908-8c1b-480d-8c4d-6ab41d73a1fc/implementation_plan.md) as living artifacts.
- Update these files to reflect progress, then refer to them briefly rather than dumping full status updates in the chat.
- Use scratch files (`<appDataDir>/brain/<conversation-id>/scratch/`) for temporary scripts or dump data.

### 4. Dynamic Context Loading (Just-in-Time)
Load files when they are required, and unload or summarize them once processed.
- Use targeted tools like `grep_search` with filter parameters to find exact symbols.
- If a subtask is complete, synthesize the result into a concise summary and proceed, allowing the older detailed files to naturally age out of the immediate reasoning window.

### 5. Isolation of Subagent Boundaries
Multi-agent systems must enforce hard boundaries between agents to keep individual context windows clean.
- Spawning a subagent (e.g. `@security-reviewer`) creates a clean slate.
- The subagent is provided only with the target code and security rules. It is not bogged down by database schema discussions, CSS layouts, or CLI flags.
- The subagent returns a structured report (`security-report.md`) containing the synthesis of its findings. The orchestrator reads only this report, not the subagent's internal logs.

---

## Verification
To verify if context is optimized:
- Check that the active open files in the IDE only belong to the component currently being modified.
- Verify that terminal/command outputs in the chat history are concise (use filters, `head`, or limit counts like `git log -n 5`).
- Ensure no duplicate copies of files or guidelines exist within the context.
