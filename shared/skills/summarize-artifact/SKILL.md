---
name: summarize-artifact
description: Produces a ~200-word summary of any pipeline artifact (analysis.md, architecture-notes.md, etc.) for downstream agents that need its gist, not its full text — the mechanism behind "context decay" in deliver-feature, where artifacts from 2+ phases prior get read as a summary instead of loaded in full. Supports --persist <feature-name> to write the summary to docs/features/<feature-name>/summary.md as a retrieval surrogate.
triggers:
  keywords: ["summarize-artifact", "summarize this artifact", "context decay"]
  intentPatterns: ["Summarize *.md for context", "/summarize-artifact *", "/summarize-artifact --persist *"]
standalone: true
---

## Not the Inter-Stage Data Path

This skill used to be how `qa-engineer` and `tech-writer` read `analysis.md`. It is not any more.

Seven pipeline artifacts are **typed state** (`analysis`, `architecture`, `route`, `review`,
`implementation-notes`, `security-report`, `qa-report` — see any of their contracts' Typed State
sections). For those, a consuming stage receives a **projection**: the fields its contract declares,
selected deterministically, with nothing paraphrased and nothing silently dropped. That is roadmap
L2.10's answer to context decay, and it is strictly better than a summary here — it costs no model
call, it is reproducible, and what it omits is declared in code rather than decided by a summarizer.

So: never call this skill to hand one stage's output to the next stage when the artifact is typed.
Two things it is still for:

1. **Untyped artifacts.** Eight of the fifteen artifacts still pass as markdown — the end-of-pipeline
   reports and `context-manifest`. Context decay still applies to them, unchanged.
2. **The retrieval surrogate** (`--persist`, `deliver-feature/SKILL.md` step 37a) and any other
   human-facing gist. The surrogate is not a handoff between stages; it is an indexing target for the
   retrieval tier, and prose is the point rather than a compromise.

## Flags

| Flag | Effect |
|---|---|
| *(none)* | Ephemeral mode. Summary is written to output only — not persisted to disk. Default, unchanged. |
| `--persist <feature-name>` | Surrogate mode. Summary is written to `docs/features/<feature-name>/summary.md` as a retrieval surrogate AND returned in output. The surrogate is what future BM25/vector retrieval tiers should index first. |

## When To Use
- A downstream agent needs the gist of an **untyped** artifact produced 2+ phases earlier in the
  pipeline (see `deliver-feature/SKILL.md`, "Context Decay") — one of the eight artifacts that still
  pass as markdown, whose broad strokes matter but whose exact wording usually doesn't.
- Standalone: any time a human wants a quick gist of a long artifact before deciding whether to read it in full.
- `deliver-feature` calls this with `--persist` after all artifacts are persisted to `docs/features/<name>/` to create a durable retrieval surrogate for the entire feature delivery.

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
5. **If `--persist <feature-name>` is set**: write the output to `docs/features/<feature-name>/summary.md`
   using the Surrogate Output Format below. Do not write to a temp path first — write directly to the
   feature archive. If the file already exists, overwrite it (a re-delivery supersedes the previous summary).

## Output Format (ephemeral mode)
```markdown
## Summary: [artifact filename]
[~200 words]

Full artifact: [path], in case the detail matters for your specific task.
```

## Surrogate Output Format (--persist mode)
The persisted `summary.md` file uses a distinct header that signals its retrieval role to any future
indexing tier. Future BM25/vector tiers should index `docs/features/*/summary.md` files first — this
header makes the file's purpose machine-readable.

```markdown
<!--
retrieval-surrogate: true
feature: <feature-name>
source-artifact: docs/features/<feature-name>/analysis.md
generated: <ISO-date>
index-first: true
-->

# Feature Summary: [Feature Name]

[~200-word summary]

Full artifact set: docs/features/<feature-name>/
```

## Guardrails
- **Never** drop a constraint that would cause downstream work to violate it (e.g. a Non-Functional
  Requirement's SLA, a security constraint like "must not reveal user enumeration") just to hit a word
  count — correctness of the summary matters more than length.
- **Never** summarize `implementation-notes.md`'s Self-Review Checklist or `code-review-report.md`'s Design
  Score for `code-reviewer`/`security-reviewer` — those need exact values, not gists (they're the artifact
  currently being reviewed, not an aging one).
- **In surrogate mode**: summarize from `analysis.md` (the artifact with the broadest scope), not from a
  later-stage artifact like `delivery-summary.md` — the summary's purpose is to describe what the feature
  *does* for retrieval, not to replay the pipeline's conclusion.
- This skill writes only to `docs/features/<feature-name>/summary.md` when `--persist` is set. It never
  edits the original artifact.

## Standalone Mode
Pure local file read + summarization. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
