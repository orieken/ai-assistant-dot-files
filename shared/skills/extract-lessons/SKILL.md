---
name: extract-lessons
description: Scans pipeline artifacts across many past deliveries for recurring security findings, recurring code-review rejections, and architecture decisions worth promoting to a Knowledge Item — then drafts proposed rule/prompt changes for human approval and records everything (approved or not) in docs/lessons-learned/.
triggers:
  keywords: ["extract-lessons", "lessons learned", "recurring findings", "promote to rule"]
  intentPatterns: ["Extract lessons from recent deliveries", "What patterns keep recurring?", "/extract-lessons *"]
standalone: true
---

## When To Use
Periodically, once enough deliveries exist to see a pattern (this skill itself refuses to claim a pattern
from fewer than 3 occurrences — see Guardrails). Natural cadence: after a few `/retrospective` runs or
alongside `agent-scorecard`'s monthly cycle, though nothing auto-invokes this — it's a deliberate,
manually-triggered analysis, not a background job.

Do NOT use this to make a single delivery's observations — that's `/retrospective`. This skill only speaks
in terms of things that recur *across* deliveries.

## Context To Load First
1. `docs/features/*/security-report.md` — all past deliveries
2. `docs/features/*/code-review-report.md` — all past deliveries
3. `docs/features/*/architecture-notes.md` — all past deliveries
4. `docs/features/*/context-manifest.md` — for the KI usage tally
5. `shared/knowledge/*.md` and `.claude/knowledge/*.md` — the current KI corpus
6. `.claude/rules/approval-gates.md` — this skill's promotion step is gated by it, not exempt from it

## Process

### 1. Recurring security findings -> candidate guardrail
Scan every `security-report.md`'s `## Findings` section. Group by threat category + rough description
similarity (e.g. "user enumeration via differing error messages" appearing in 3+ distinct features). For
any pattern at 3+ occurrences:
- Draft a proposed addition to the relevant file in `shared/rules/` (usually
  `architecture-guardrails.md`) stating the guardrail as a rule, not a finding.
- **Do not write it.** Per `approval-gates.md` Gate #7 ("Wiring a New Fitness Function"), present the draft
  and require explicit "approve fitness function" or "add to CI" before touching `shared/rules/`. This
  applies even though the pattern is well-evidenced — the gate exists specifically so a promotion isn't a
  drive-by edit to a file every agent treats as session-long law.

### 2. Recurring code-review rejections -> candidate prompt improvement
Scan every `code-review-report.md`'s `## Feedback for the Developer` CHANGES REQUESTED sections. Group by
named refactoring operation or architectural violation type. For any pattern at 3+ occurrences:
- Draft a proposed edit to the relevant agent's prompt (usually `developer.md`'s guardrails, or the
  specific rule the violation slipped past).
- **Do not write it.** An agent prompt edit requires a version bump + `shared/agents/CHANGELOG.md` entry
  (Epic 8) in the same commit — present the draft, the proposed version bump, and the changelog line, and
  wait for explicit confirmation before applying any of it.

### 3. Recurring architecture decisions -> candidate KI
Scan every `architecture-notes.md`'s `## Structural Decisions`. If the same decision (same problem, same
resolution) appears in 3+ features instead of being referenced from a shared source, that's a sign it
should be a KI so future architects reference it instead of re-deciding it. Invoke `create-ki` to draft it
(searches for a duplicate first) — creating a KI is not gated the way rules/prompts are (see
`approval-gates.md`; it isn't listed there), but still confirm the draft content with the user before
`create-ki` writes it, since a bad KI actively misleads future agents.

### 4. KI usage analytics
Tally, across all `context-manifest.md` files, how many times each KI in `shared/knowledge/` and
`.claude/knowledge/` was actually listed under "4. Knowledge Items & ADRs (To Load)." Flag:
- KIs referenced 0 times across all available history — candidates for removal or re-tagging (maybe its
  tags don't match how tasks actually get described).
- KIs referenced frequently — evidence the KI system is paying for itself.

### 5. Write the record
Write `docs/lessons-learned/lessons-[YYYY-MM-DD].md` with every finding from steps 1-4, **regardless of
whether the user approves any promotion** — the lessons file is the permanent record that a pattern was
noticed; promotion to a rule/prompt/KI is a separate, gated action tracked in the same file's status column.

## Output Format
Write `docs/lessons-learned/lessons-[YYYY-MM-DD].md`:
```markdown
# Lessons Extracted: [YYYY-MM-DD]

## Scope
- Deliveries scanned: [N] — [docs/features/ subdirectories included]

## Recurring Security Findings
| Pattern | Occurrences | Features | Proposed Guardrail | Status |
|---|---|---|---|---|
| [pattern] | [N] | [feature names] | [draft rule text, or link to it below] | Proposed / Approved / Declined |

## Recurring Code-Review Rejections
| Pattern | Occurrences | Features | Proposed Prompt Change | Status |
|---|---|---|---|---|
| [pattern] | [N] | [feature names] | [agent + draft edit summary] | Proposed / Approved / Declined |

## Architecture Decisions -> KI Candidates
| Decision | Occurrences | Features | KI Status |
|---|---|---|---|
| [decision] | [N] | [feature names] | Created (link) / Declined / Duplicate of existing KI (link) |

## KI Usage Analytics
| KI | Times Referenced | Last Referenced | Note |
|---|---|---|---|
| [ki name] | [N] | [feature name / "never"] | [flag if 0] |

## Declined / Deferred
[Anything the user explicitly declined to promote, and why — so it isn't silently re-proposed identically next run]
```

## Guardrails
- **Never** claim a recurring pattern from fewer than 3 distinct feature occurrences — say "only seen once
  or twice, not yet a pattern" instead.
- **Never** write to `shared/rules/` or bump/edit an agent's prompt without the explicit confirmation
  `approval-gates.md` requires for that action — this skill drafts, it does not apply, for those two cases.
- **Always** write the lessons-learned record even when nothing is approved — a declined promotion is still
  useful history (and prevents re-proposing the identical thing next run without acknowledging it was
  already declined).
- This is additive: never edit or delete a prior `docs/lessons-learned/` file.

## Standalone Mode
Pure local file reads/writes; the two gated promotion paths pause for human input exactly like any other
approval-gated action in this framework — no external services required.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
