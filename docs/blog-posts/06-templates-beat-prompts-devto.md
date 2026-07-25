---
title: "Templates Beat Prompts: Fixing Agent Output Structure Drift"
published: false
description: "Why inlining output markdown examples inside agent prompt files causes structural drift—and how single-source-of-truth templates solved it."
tags: ai, architecture, prompt-engineering, devtools, testing
canonical_url:
cover_image:
---

If you build multi-agent systems, you have almost certainly encountered **output structure drift**.

You write a prompt for an `analyst` agent or a `security-reviewer` agent. Inside the prompt, you include a markdown example showing the expected output layout—complete with headings like `## Summary`, `### Acceptance Criteria`, and `## Task List`.

For a while, everything works. Then, subtle LLM drift creeps in:
- One prompt version outputs `## Overview` instead of `## Summary`.
- Another outputs `### Tasks` instead of `## Task List`.
- A third agent decides to wrap section titles in bold bullet points (`- **Acceptance Criteria**:`).

Downstream agents—and automated contract validation scripts like `validate-artifact`—rely on exact heading strings to parse task lists and requirements. A single drifted heading breaks the pipeline contract.

Here is how we solved agent output drift in `ai-assistant-dot-files` by extracting **13 single-source-of-truth templates** to `shared/templates/` and decoupling format definitions from agent personality prompts.

---

## The Root Cause: Inline Prompt Duplication

In early versions of our framework, every agent prompt in `shared/agents/` contained an inline markdown example of its expected output. 

For example, `shared/agents/analyst.md` contained a 50-line inline example of `analysis.md`. Meanwhile, `shared/contracts/analysis-contract.md` maintained a grep list of required headings to validate that artifact.

This violated DRY (Don't Repeat Yourself) in two ways:
1. The structural definition was duplicated across the agent prompt and the validation contract.
2. Updating an artifact format required editing both the agent prompt AND the contract file.

If an engineer updated the contract but forgot to update the prompt, the agent would keep generating the old structure—causing contract validation checks (`validate-artifact`) to fail and triggering unnecessary agent retry loops.

---

## The Solution: Extract Templates as the Single Source of Truth

In commit `6973541`, we executed a systematic refactoring across our 26-agent roster:

1. **Created `shared/templates/`**: We extracted 13 standardized template files (e.g., `analysis.template.md`, `architecture.template.md`, `qa.template.md`).
2. **Updated Agent Prompts**: We removed inline markdown examples from all agent definition files in `shared/agents/`. Agents were instructed to read their canonical template file from `shared/templates/` before generating output.
3. **Aligned Contracts**: Contract validators in `shared/contracts/` were updated to validate strictly against the template's exact heading levels and strings.

```markdown
<!-- shared/templates/analysis.template.md -->
# Feature Analysis: [Feature Name]

## Summary
[One paragraph summary]

### Acceptance Criteria
[BDD scenarios]

### Non-Functional Requirements
[Performance, Security, Scaling]
...
```

Instead of embedding 60 lines of markdown in `analyst.md`, the prompt now simply specifies:

> "Read `shared/templates/analysis.template.md` and produce your artifact at `.claude/feature-workspace/analysis.md` by filling in the bracketed `[placeholder]` markers. Preserve every heading exactly as it appears in the template."

---

## The Loader Gotcha: The `AGENT_TEMPLATE` Edge Case

During the initial extraction, we created an example agent template file named `shared/agents/AGENT_TEMPLATE.md` to guide authors creating new custom agents.

However, our automated platform projection pipeline (`scripts/generate-configs.sh`) scans `shared/agents/` to build `.cursorrules`, `.windsurfrules`, `.copilot-instructions.md`, and `AGENTS.md`. Because `AGENT_TEMPLATE.md` lived inside `shared/agents/`, the generator registered `AGENT_TEMPLATE` as a 27th active agent persona!

In commit `ca827dc`, we fixed this loader collision by moving `AGENT_TEMPLATE.md` out of the loader path to `shared/templates/agent.template.md`. 

**Rule learned**: Never place meta-templates inside auto-discovered registry directories.

---

## Results and Benefits

- **Zero Contract Drift**: Contract validation pass rates increased, and header-mismatch retries dropped to zero across pipeline runs.
- **Maintainability**: Changing an output structure now requires editing exactly one file in `shared/templates/`.
- **Token Efficiency**: Agent prompts became cleaner, shorter, and easier to maintain without bloated inline markdown examples.

---

## Summary Checklist for Multi-Agent Output Formatting

1. **Decouple formatting from personality.** Agent prompts should define *how to think*; templates should define *how to format*.
2. **Treat templates as contracts.** Ensure automated linter/validator scripts check against the exact template headings.
3. **Isolate meta-templates.** Keep template files out of auto-scanned agent discovery paths.

---

## Image Prompt

> Hero image prompt: A clean, modern blueprint schematic showing modular puzzle-piece document templates snapping precisely into a standardized digital conveyor belt, eliminating distorted text lines. High contrast dark background with cyan blueprint grid lines and glowing white typography.
