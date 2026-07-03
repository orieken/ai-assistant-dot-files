---
name: summarize-artifact
description: Produces a ~200-word summary of any pipeline artifact (analysis.md, architecture-notes.md, etc.) for downstream agents that need its gist, not its full text — the mechanism behind "context decay" in deliver-feature, where artifacts from 2+ phases prior get read as a summary instead of loaded in full.
triggers:
  keywords: ["summarize-artifact", "summarize this artifact", "context decay"]
  intentPatterns: ["Summarize *.md for context", "/summarize-artifact *"]
standalone: true
---

## When To Use
- A downstream agent needs the gist of an artifact produced 2+ phases earlier in the pipeline (see
  `deliver-feature/SKILL.md`, "Context Decay") — e.g. `qa-engineer` or `tech-writer` needing `analysis.md`'s
  scope without re-reading every acceptance criterion verbatim, since they already have
  `implementation-notes.md` (which restates the decisions that matter for their job).
- Standalone: any time a human wants a quick gist of a long artifact before deciding whether to read it in full.

Do NOT use when the agent's task actually depends on exact wording — e.g. `code-reviewer` checking
`implementation-notes.md`'s Self-Review Checklist needs the literal checked items, not a paraphrase. Context
decay applies to *older* artifacts whose broad strokes still matter but whose exact wording usually doesn't;
it never applies to the artifact an agent is immediately reviewing.

## Context To Load First
1. The artifact to summarize (path given by the caller)

## Process
1. Read the artifact in full once.
2. Identify what actually matters to a downstream reader: the decision/outcome, not the reasoning that led
   to it (the reasoning is preserved in the full file for anyone who needs to dig in).
3. Write a summary of roughly 200 words (150-250 is fine; don't pad to hit exactly 200) covering:
   - What the artifact concluded (scope, decision, or result — not process)
   - Anything a downstream agent would get wrong by *not* knowing it
   - A pointer back to the full file for anyone who needs the detail
4. Do not summarize a summary — always summarize from the original artifact, so quality doesn't degrade
   across repeated compressions.

## Output Format
```markdown
## Summary: [artifact filename]
[~200 words]

Full artifact: [path], in case the detail matters for your specific task.
```

## Guardrails
- **Never** drop a constraint that would cause downstream work to violate it (e.g. a Non-Functional
  Requirement's SLA, a security constraint like "must not reveal user enumeration") just to hit a word
  count — correctness of the summary matters more than length.
- **Never** summarize `implementation-notes.md`'s Self-Review Checklist or `code-review-report.md`'s Design
  Score for `code-reviewer`/`security-reviewer` — those need exact values, not gists (they're the artifact
  currently being reviewed, not an aging one).
- This is read-only — it produces a summary, it does not edit the original artifact.

## Standalone Mode
Pure local file read + summarization. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
