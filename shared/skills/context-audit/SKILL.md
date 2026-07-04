---
name: context-audit
description: Analyzes a completed pipeline run (or the current conversation) for context waste — pinned files that were never actually referenced in the output, duplicate information across loaded files, and large files loaded without a line-range constraint. Distinct from pipeline-retrospective (timing trends) and agent-scorecard (quality trends) — this one is specifically about whether the context budget was spent well.
triggers:
  keywords: ["context-audit", "context waste", "wasted context", "audit context"]
  intentPatterns: ["Audit context for *", "Was context wasted on *", "/context-audit *"]
standalone: true
---

## When To Use
- After a delivery, to check whether `context-engineer`'s manifest actually earned its keep — were the
  pinned files used, or was budget spent on files nobody referenced?
- Mid-conversation, when a session feels sluggish or repetitive, to check the current session's own tool-call
  history for the same waste patterns (full-file reads on large files, duplicate reads, files read but never
  mentioned again).

Do NOT use for timing/duration analysis — that's `pipeline-trace`/`pipeline-retrospective`. Do NOT use for
whether an agent's *output* was good — that's `agent-scorecard`. This skill only judges context spend.

## Context To Load First
1. For a pipeline run: `docs/features/<feature-name>/context-manifest.md` and the artifacts produced by the
   agents that consumed it (`analysis.md`, `implementation-notes.md`)
2. For the current conversation: the session's own tool-call history (Read/Grep/Glob calls and their
   arguments — line ranges used or not, files touched more than once)

## Process

### 1. Unused pins
For each file in `context-manifest.md`'s "2. Pinpoint Files (To Keep Open)," check whether the consuming artifact
(`analysis.md` for analyst, `implementation-notes.md` for developer) actually references it — cites a
specific finding, pattern, or interface from it, not just a generic mention. A pinned file with zero
identifiable influence on the output is waste: it consumed budget context-engineer estimated as necessary,
but wasn't.

### 2. Duplicate information
Check whether two or more pinned files describe overlapping ground (e.g. two files that both define the
same interface, or a file and a KI that say the same thing). Flag it — one of them should have been pruned
or the KI should have been referenced instead of the raw file.

### 3. Unconstrained large reads
Cross-reference `context-manifest.md`'s Pinpoint Files against actual file line counts. `context-engineer`'s
own guardrail requires range-constraining anything over 500 lines — flag any pin over 500 lines with no
`#Lxx-Lyy` range as a guardrail violation, not just inefficiency.

### 4. Session-level audit (conversation mode)
When auditing the live conversation instead of a persisted pipeline run: scan the tool-call history for the
same three patterns — full reads of files over 500 lines with no `offset`/`limit`, the same file read more
than once without the content changing, and files read but never referenced again in any subsequent
response.

### 5. Estimate wasted tokens
For each flagged item, estimate the token cost using the same heuristic `context-engineer` uses (~line count
× 8 chars/line ÷ 4 chars/token) so the waste is quantified, not just qualitative.

## Output Format
```markdown
# Context Audit: [Feature Name or "Current Session"]

## Unused Pins
- [file] -- pinned for [stated reason], not referenced in [artifact] -- ~[N] tokens

## Duplicate Information
- [file A] and [file B] -- both describe [overlapping content] -- one should have been pruned

## Unconstrained Large Reads
- [file] ([N] lines) -- pinned without a line range, violates context-engineer's >500-line guardrail

## Estimated Total Waste
~[N] tokens across [M] flagged items ([X]% of the manifest's estimated total)

## Recommendation
[Specific: which pin to drop next time, which guardrail was violated and should be enforced going forward]
```

## Guardrails
- **Never** flag a pin as unused without checking the artifact reasonably thoroughly first — a file can
  legitimately inform an agent's reasoning without being explicitly cited by name.
- **Never** re-estimate or fabricate a token count without the same heuristic `context-engineer` uses —
  consistency matters more than precision here, since this skill's whole point is comparing against that
  estimate.
- This is a read-only analysis — it does not modify `context-manifest.md` or any pipeline artifact.

## Standalone Mode
Pure local file/history reads. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
