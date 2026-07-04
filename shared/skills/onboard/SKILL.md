---
name: onboard
description: An interactive tour for someone new to this framework — explains rules/agents/skills, shows how to invoke an agent and trigger a skill, walks through running a full pipeline, and lists the approval gates that pause for human confirmation. Invoked after install.sh --tour, or any time a new user asks "what is this" / "how do I use this."
triggers:
  keywords: ["onboard", "tour", "getting started", "how does this work", "what is this framework"]
  intentPatterns: ["Give me a tour", "How do I use this framework", "I'm new here", "/onboard"]
standalone: true
---

## When To Use
- Right after `install.sh --tour` completes (the flag prints a pointer to this skill — installing itself
  can't invoke an AI skill, since it's a plain shell script; see `install.sh`'s tail for exactly what it
  prints).
- Any time a new user or a fresh AI session needs the "what is this and how do I use it" walkthrough,
  without reading `README.md` and `docs/ARCHITECTURE.md` cover to cover first.

Do NOT use for a specific how-to (adding an agent, editing a prompt) — point to
`docs/CONTRIBUTING.md` or the relevant runbook instead. This skill is breadth-first orientation, not a
reference manual.

## Context To Load First
1. `README.md` — for accurate agent/skill counts and the platform matrix
2. `DOMAIN_DICTIONARY.md` — for the exact Rule/Agent/Persona/Skill definitions used below
3. `.claude/rules/approval-gates.md` — for the gate list in step 5

## Process

Walk through these five things, in order, adjusting depth to what the user already seems to know (don't
re-explain what they've clearly already grasped from their own question):

### 1. The three context layers
- **Rules** (`shared/rules/`) — governance documents every agent obeys: architecture guardrails, design
  principles, approval gates. Session-long, never optional.
- **Agents** (`shared/agents/`, 24 of them) — specialists with tool access and a defined process: analyst,
  developer, code-reviewer, security-reviewer, and so on. Only fully real on Tier 1 (Claude Code) — on
  Cursor/Copilot/etc. they're **personas**: the same knowledge, but a context frame with no autonomous
  tool-driven workflow, since those tools can't run one.
- **Skills** (`shared/skills/`, 53 of them) — on-demand capabilities triggered by keyword or slash command:
  `/deliver-feature`, `/complexity-check`, `/threat-model`, and so on.

### 2. How to invoke an agent
On Claude Code, agents are real subagents:
```
> @analyst please read features/my-feature.md and produce an analysis
```
On any other tool, reference the file directly, since there's no native agent mechanism:
```
Act exactly as described in shared/agents/code-reviewer.md. Review my current changes.
```

### 3. How to trigger a skill
Slash command, or just describe what you want in the words from that skill's trigger patterns:
```
> /complexity-check src/
> /threat-model src/api/checkout.ts
```

### 4. How to run a full pipeline
```
> /new-feature "password reset via email"        # interview -> features/password-reset.md
> /deliver-feature features/password-reset.md     # full agent sequence, see README's pipeline diagram
```
Point to `shared/templates/my-first-feature.md` here — it's a pre-written, already-complete feature spec
meant to be run through `/deliver-feature` immediately, without needing `/new-feature` first, so a new user
sees the whole pipeline in action before writing their own spec.

### 5. Approval gates
List these from `.claude/rules/approval-gates.md` — the pipeline (and any agent) pauses and requires
explicit confirmation before: creating a git commit, shipping to Friday, running a database migration,
the Contract phase of a migration, posting to an external API, writing outside the feature workspace/source
directories, wiring a new fitness function (rule change), or deploying. Any edit to a pending artifact resets
the gate — so approving once doesn't mean approving forever.

## Output Format
Free-form, conversational — this is a tour, not a report. Keep each of the five sections to a few sentences
unless the user asks to go deeper on one of them. End by pointing at `shared/templates/my-first-feature.md`
as the concrete next action ("run `/deliver-feature shared/templates/my-first-feature.md` to see the whole
pipeline for real").

## Guardrails
- **Never** overwhelm a new user with the full 24-agent roster or 53-skill catalog in the tour itself —
  point to `README.md` for the exhaustive list, cover the concepts here.
- **Never** run `/deliver-feature` or any other pipeline step *for* the user during the tour — describe it,
  let them decide when to actually run it (the tour is orientation, not an unsolicited pipeline execution).

## Standalone Mode
Pure conversational walkthrough referencing local files. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
