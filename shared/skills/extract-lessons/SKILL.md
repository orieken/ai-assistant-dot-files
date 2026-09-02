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
1. `docs/incidents/*.md` — all incident records (Epic 67: new first-class input)
2. `docs/features/*/security-report.md` — all past deliveries
2. `docs/features/*/code-review-report.md` — all past deliveries
3. `docs/features/*/architecture-notes.md` — all past deliveries
4. `docs/features/*/context-manifest.md` — for the KI usage tally
5. `shared/knowledge/*.md` and `.claude/knowledge/*.md` — the current KI corpus
6. `.claude/rules/approval-gates.md` — this skill's promotion step is gated by it, not exempt from it
7. `run-state.json` and `run-events.jsonl` for each feature delivered by `loom run` (their absence
   is not an error) — the `artifact.corrected` entries, with their retained diffs. These are
   first-class corrective signals: a human editing an artifact at a gate is the strongest
   real-world evidence that the producing agent needs improvement, and the record already names
   which agent.

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

### 5. Gate decision patterns -> candidate agent improvement

**Two sources, and only one of them exists today.**

**5a. `run-events.jsonl` (executor runs — this is the one that has data).** For any feature delivered
by `loom run`, read `.claude/feature-workspace/<feature>/run-events.jsonl` and
`docs/features/<feature>/run-state.json` for `artifact.corrected` entries (roadmap L4.5). Each is a
human editing a stage's output at a gate, already attributed to the **producing agent** — no
inference from gate ownership needed — with a retained unified diff under
`.approved/<gate>/corrections/`. Read the diff: it is a labelled example of what that agent got
wrong and what a human thought it should have said, which is far stronger evidence than a reason
string. Group by agent, and apply the same 3+-occurrences-across-distinct-features bar below.

Two caveats to carry into any hypothesis. The correction was **advisory** — the pipeline did not
adopt the human's text, so do not assume the delivered artifact contains it. And a single human's
edits are one person's judgement, not a verdict; the recurrence bar exists precisely because one
correction is an anecdote.

**5b. The markdown pipeline records nothing.** A delivery run by a host platform's LLM rather than by
`loom run` leaves no correction record. `.claude/telemetry/events.jsonl` and its `gate_decision`
events were specified, never emitted, and retired in roadmap L3.9 — do not look for that file, and
do not treat its absence as telemetry being switched off. For those deliveries, gate-decision
analysis is unavailable, and saying so in the output is better than inferring it from artifacts.

Apply the same bar to step 5a's corrections: 3+ occurrences for the same agent across distinct
features before drafting a hypothesis.
- Draft a hypothesis for what consistent mistake caused the human to correct that agent, reading the
  retained diffs for recurring language.
- If a hypothesis is credible, draft it as a candidate agent-prompt improvement (same gating as step
  2: present, require explicit confirmation, don't auto-apply).
- Record every correction pattern in the lessons-learned output regardless of promotion decision.

### 6. Incident-feature pair analysis -> candidate pipeline improvement
Scan every `docs/incidents/*.md` record. For each incident whose **Affected Feature** field links to a
`docs/features/<name>/` delivery:
- Load the linked feature's pipeline artifacts (`analysis.md`, `implementation-notes.md`,
  `security-report.md`, `code-review-report.md`) alongside the incident record.
- Ask: **"Which pipeline stage should have caught the bug or architectural gap that caused this incident?"**
  Frame answers as "the `[agent]` stage should have flagged X" rather than "the developer made a mistake."
- If the answer points to a consistent gap (a stage that was skipped, a check that wasn't thorough enough,
  or an edge case not covered by the QA contract), draft it as a proposed rule/prompt change — same gating
  as Step 2: present the draft and require explicit confirmation before applying anything.
- If the incident's **Candidate Records** section already contains `promote-memory`-format candidates,
  those flow directly into the output's "Incident Candidates" table — no re-analysis needed, just
  consolidate and update status.
- A single incident is not a pattern — only draft a rule/prompt change if the same pipeline gap surfaces
  in 2+ incident-feature pairs (lower bar than the 3-occurrence rule in Steps 1-2 because production
  incidents carry stronger signal than code-review findings alone).

### 7. Write the record
Write `docs/lessons-learned/lessons-[YYYY-MM-DD].md` with every finding from steps 1-6, **regardless of
whether the user approves any promotion** — the lessons file is the permanent record that a pattern was
noticed; promotion to a rule/prompt/KI is a separate, gated action tracked in the same file's status column.

## Output Format
Write `docs/lessons-learned/lessons-[YYYY-MM-DD].md`:
```markdown
# Lessons Extracted: [YYYY-MM-DD]

## Scope
- Deliveries scanned: [N] — [docs/features/ subdirectories included]
- Incidents scanned: [N] — [docs/incidents/ records included]
- Incident-feature pairs matched: [N]

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

## Gate Corrections
Corrections recorded by the executor (`artifact.corrected`, from step 5a — omit the table when a
feature was not delivered by `loom run`):
| Agent | Feature | Gate | Diff | What the human changed |
|---|---|---|---|---|
| [agent] | [feature] | [gate] | +N/-M | [one line from reading the diff] |

Recurring patterns (3+ corrections for the same agent across distinct features):
| Agent | Count | Features | Hypothesis | Proposed Action | Status |
|---|---|---|---|---|---|
| [agent] | [N] | [feature list] | [draft hypothesis] | [prompt edit / rule change / none] | Proposed / Approved / Declined |

State explicitly when a delivery in scope came from the markdown pipeline: it records no corrections,
so its absence from this table means "not observable", not "nothing was corrected".

## Incident-Feature Pairs (Step 6)
| Incident | Feature | Pipeline Gap Identified | Proposed Change | Status |
|---|---|---|---|---|
| [docs/incidents/slug.md] | [docs/features/name/] | [agent stage / check] | [draft rule or prompt change] | Proposed / Approved / Declined |
— or "No matched incident-feature pairs in docs/incidents/"

## Incident Candidates (promote-memory format from incident Candidate Records)
| Incident | Candidate Title | Type | Status |
|---|---|---|---|
| [slug] | [short title] | KI / ADR-worthy / Rule-change-worthy / Lesson | Proposed / Approved / Declined |
— or "None"

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
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
