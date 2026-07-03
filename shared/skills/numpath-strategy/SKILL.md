---
name: numpath-strategy
description: Enforces the NumPath Build → Test → Document loop after any feature build or experiment. Triggers when a feature is complete, a phase milestone is reached, or the user asks to "capture", "document", "write up", or "take a research note". Produces a blog post draft, optional research note, and optional paper outline update.
triggers:
  keywords: ["capture", "document", "write up", "post", "blog", "paper", "phase complete", "milestone", "research note", "note this"]
  intentPatterns:
    - "write a blog post about *"
    - "document what we built"
    - "capture the design decisions"
    - "update the research paper"
    - "phase * is done"
    - "feature * is complete"
    - "note this observation"
    - "take a research note"
standalone: true
---

## When To Use

Use this skill after completing any of the following:
- A Phase milestone (Phase 1 DoD met, Phase 2 BKT shipped, etc.)
- A significant feature (adaptive engine, BKT modeler, teacher dashboard, LLM insight generator)
- An architectural decision that produced an ADR
- An unexpected observation during a build or test session (use Step 2.5 only)
- An experiment result worth capturing

Do NOT use for: in-progress work, minor bug fixes, refactors with no user-facing change.

## Context To Load First

1. `docs/strategy.md` — the governing workflow document
2. `docs/phd-roadmap/MASTER_PROMPT.md` — project identity and current phase
3. The ADR(s) or feature spec relevant to what was just built
4. `docs/posts/ROADMAP.md` — blog post plan (to find the matching planned post)
5. `docs/posts/` — scan for existing posts (pick up numbering, avoid duplication)
6. `docs/notes/` — scan for existing research notes (pick up numbering)
7. `docs/papers/` — scan for existing paper drafts (to update, not replace)

---

## Process

### Step 1 — Identify what was built or observed
Ask (or infer from context):
- What feature or milestone was completed? Or what unexpected behavior was observed?
- What design decisions were made that aren't obvious from the code?
- Were there any surprising trade-offs, error patterns, or model behaviors?
- What metrics or observations were captured?

---

### Step 2 — Write the blog post draft
*(Skip if this is an observation-only session with no shippable feature)*

File: `docs/posts/NN-<kebab-title>.md` (NN = next number in sequence)

Check `docs/posts/ROADMAP.md` first — if a planned post matches, use its working title.

Posts are published to **dev.to**. Use this exact frontmatter block, then the post body:

```markdown
---
title: "Post Title — written for a developer audience, not academic"
description: "One sentence shown in previews and SEO. What was built and why it matters."
tags:
    - numpath
    - adaptive-learning
    - dyscalculia
    - [1–2 tech tags relevant to this post, e.g. python, vue, bayesian]
series: "Building NumPath"
published: false
canonical_url: ""
devto_id: 0
---

## What We Built
[2–3 sentences. What the feature does and why it exists.]

## The Design Decision
[The core decision: what was chosen and what was rejected. Reference the ADR if one exists as a link.]

## Why It Matters for the Research
[1 paragraph connecting the engineering choice to the research question.]

## What We Learned
[Honest reflection: what worked, what surprised us, what we'd do differently.]

## What's Next
[One sentence on the next step.]

## Key Takeaways
- [Bullet 1 — specific technical insight]
- [Bullet 2 — research implication]
- [Bullet 3 — what you'd do differently]
```

Rules for dev.to format:
- `tags` max 4 entries — always include `numpath` and at least one tech tag
- `series` is always `"Building NumPath"` unless the post belongs to a sub-series (e.g. `"NumPath Research"` for experiment/paper posts)
- `published: false` always — user sets to `true` manually before posting
- `devto_id: 0` always — dev.to writes the ID back after first publish
- Title uses sentence case, no trailing period

**GEO checklist** (apply before showing draft to user):
- [ ] Title contains the primary concept (e.g. "Bayesian Knowledge Tracing", "adaptive engine")
- [ ] First paragraph answers: what problem does this solve?
- [ ] Each `##` section can stand alone as a quoted excerpt
- [ ] Includes at least one code block or diagram
- [ ] Ends with a **Key Takeaways** section (3 bullets — AI engines extract this for summaries)
- [ ] `published: false` is set

Show the draft + GEO checklist to the user for review before writing to disk.
Gate: user must say **"approve post"** or **"publish post"** before the file is written.

---

### Step 2.5 — Research note (if applicable)
Write a research note when:
- Something unexpected happened during the session (model behaved oddly, accuracy spiked, error pattern emerged)
- A hypothesis formed that isn't yet validated by data
- An observation is too raw for a blog post but worth keeping

File: `docs/notes/NN-<kebab-title>.md` (NN = next number in sequence, see `docs/notes/README.md`)

Use the template at `docs/templates/research-note.md`.

Fill in: What We Observed, Hypothesis, Metrics to Track, Open Questions, Impact on Research, Future Work.

Show the draft to the user.
Gate: user must say **"approve note"** before the file is written.

After writing: update the index table in `docs/notes/README.md`.

---

### Step 3 — Update the research paper outline (if applicable)
If the work touches: experiment design, BKT/DKT model behavior, teacher insights, or student outcomes — update `docs/papers/01-numpath-rct-outline.md`.

Add a note under the relevant section:
```
> [YYYY-MM-DD] Phase N: [one sentence on what was learned or confirmed]
```

Show the diff to the user before writing.
Gate: user must say **"approve paper update"**.

---

### Step 4 — Check ADR coverage
If a new architectural decision was made but has no ADR:
- Prompt: "This looks like a decision worth recording. Should I write an ADR for [decision]?"
- If yes: invoke the `adr` skill

---

### Step 5 — Update the roadmap and post roadmap
- Mark completed items in `docs/architecture/mvp-feature-spec.md`
- Mark the matching post in `docs/posts/ROADMAP.md` as `in-progress` or `drafted`

---

## Output Format

```
docs/posts/NN-<kebab-title>.md           ← blog post draft (with GEO applied)
docs/notes/NN-<kebab-title>.md           ← research note (if applicable)
docs/papers/01-numpath-rct-outline.md    ← updated (if applicable)
docs/adrs/ADR-NNN-*.md                   ← new ADR (if triggered)
docs/notes/README.md                     ← index updated (if note written)
docs/posts/ROADMAP.md                    ← post status updated
```

---

## Guardrails

- Never write any artifact to disk without user approval
- Approval gates: "approve post", "approve note", "approve paper update"
- Never invent metrics — only record what was actually observed or measured
- Blog posts are first-person developer voice, not passive-voice academic prose
- Research notes are raw — incomplete sentences and open questions are fine
- Paper outlines are academic structure only — no prose, just section notes
- Any edit to a pending draft resets the approval gate (per approval-gates.md)

---

## Standalone Mode

All logic is in this file. No MCP servers required. File reads use the Read tool. File writes use the Write tool with user approval.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
