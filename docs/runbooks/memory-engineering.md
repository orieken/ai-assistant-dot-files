# Runbook: Memory Engineering

Deep-dive companion to [context-engineering.md](context-engineering.md) section 2 ("Memory") — that file is
the map of every memory mechanism and how they differ; this file is the operational detail for the
lifecycle a piece of memory (a Knowledge Item, in practice — this lifecycle doesn't apply to ADRs, the
feature archive, or DOMAIN_DICTIONARY.md, which are append-only or human-maintained records, not curated
memory — see `shared/memory-registry.json`'s per-source notes for why each source is or isn't in scope here).

## Context vs. Memory vs. Learning (recap)
- **Context**: what's loaded into an agent's window for the current turn/run.
- **Memory**: durable knowledge that outlives any single run — this runbook.
- **Learning**: feedback loops that change future agent behavior based on past runs.

Memory Engineering is specifically about keeping the *durable* layer (mainly Knowledge Items) healthy over
time — not too sparse to be useful, not so bloated with duplicates and stale entries that a search through
it stops being cheap or trustworthy.

## The Lifecycle

```
Capture -> Candidate -> Audit -> Approve -> Index -> Retrieve -> Expire
```

1. **Capture** — a pattern, bug fix, or decision happens somewhere: a delivery's `retrospective.md`, a
   recurring finding `extract-lessons` notices across many deliveries, or a human directly asking to
   remember something.
2. **Candidate** — the raw capture becomes a structured candidate with the fields in "Memory Contract"
   below. Produced by `promote-memory` (single retrospective, immediate) or `extract-lessons` (recurring
   pattern across many deliveries, periodic) — never written directly to `shared/knowledge/` without passing
   through a candidate first, even for a single-line fix, so there's always a paper trail of *why* something
   was captured.
3. **Audit** — `memory-engineer` checks the candidate against the existing corpus: is this a duplicate of an
   existing KI (in which case, update the existing one rather than create a near-copy)? Is it too narrow or
   too speculative to be reusable (reject)? Does it belong as a KI, or does it actually belong as a rule
   change, a prompt edit, or a new ADR instead (redirect)?
4. **Approve** — a human confirms the audit's recommendation before anything is written, same as any other
   commit in this repo — there is no separate approval-gate ceremony beyond the normal "Creating a Git
   Commit" gate (`shared/rules/approval-gates.md` #2). Memory promotion is not irreversible or external-facing
   enough to warrant its own gate.
5. **Index** — the approved KI is written via `create-ki` (or updated in place if it's a duplicate merge),
   with frontmatter tags/domain set so `search-ki`'s pre-filter can find it later.
6. **Retrieve** — `search-ki` (KIs + ADRs) or `query-memory` (the full registry) surfaces it during
   `context-engineer`'s Proactive RAG step, or on direct request.
7. **Expire** — `memory-engineer` periodically flags KIs that are candidates for expiration (see below);
   a human confirms actual deletion/archival via the normal commit gate, same as Approve.

## Memory Contract (Candidate and Audit fields)

Follows the same spirit as `shared/contracts/` — a required-field check, not free-form prose.

**Candidate fields** (produced by `promote-memory` / `extract-lessons`):
- **Source**: which delivery/retrospective/pattern this came from
- **Type**: KI | ADR-worthy | Rule-change-worthy | Lesson | Reject
- **Evidence**: the specific finding/quote that justifies capturing this
- **Tags**: candidate frontmatter tags for the resulting KI, if Type is KI
- **Expiration condition**: what would make this stop being true (e.g. "if `security-reviewer`'s prompt
  changes to catch this class of finding by default, this KI is redundant" — not a date, a condition)

**Audit fields** (produced by `memory-engineer`):
- **Verdict**: Approve as new | Approve as merge into `<existing KI>` | Reject | Redirect to `<rule/ADR/prompt edit>`
- **Scores**: reusability (would this help a *different* future feature, not just this one?), specificity
  (concrete enough to act on, not so generic it passes trivially)
- **Duplicates**: which existing KIs were checked and found not to overlap (or found to overlap — the merge case)
- **Final Destination**: exact file path the candidate will land at, if approved

## Promotion Rules

Promote to a KI when the pattern is:
- **Reusable**: would apply to a different feature in a different bounded context, not just this one
- **Non-obvious**: something a competent agent wouldn't already know or infer from `ARCHITECTURE_RULES.md`
- **Actionable**: phrased as "when X, do Y because Z," not just an observation

Do NOT promote when the pattern is:
- **One-off**: specific to this exact feature's constraints, won't recur
- **Already covered**: a KI or rule already says this — update the existing one instead of creating a near-duplicate
- **Too speculative**: "this might matter someday" — that's exactly the premature-generalization anti-pattern
  this framework's own `shared/rules/design-principles.md` already warns against

## Expiration Criteria

A KI is a candidate for expiration when any of:
- The pattern it documents no longer applies (the code/agent/pattern it references was removed or fully
  rewritten)
- It's been superseded by a newer, more complete KI covering the same ground
- `extract-lessons`' KI-usage-analytics check (see its own SKILL.md) shows it's never been referenced in any
  `context-manifest.md` across the available delivery history, *and* a human confirms it's not just an
  under-tagged KI that search isn't finding

Expiring a KI means moving it to `shared/knowledge/expired/` (not deleting outright) with a one-line note on
why, so the reasoning that led to writing it in the first place isn't lost — a KI that turns out to be wrong
is itself a lesson.

## Compression

When two or more KIs cover overlapping ground (found during Audit, or during a periodic `memory-engineer`
sweep), merge them into one file rather than leaving near-duplicates for `search-ki` to have to
disambiguate between at query time. Compression is a `memory-engineer` judgment call, not automatable —
don't build a script that auto-merges files based on text similarity; a human/LLM has to confirm the
merged version doesn't lose a distinction that actually mattered.

## Validation

`scripts/health-check.sh` includes a Memory Registry check: the registry parses as valid JSON, every path it
references exists, and no two KIs share an exact frontmatter `name:`. Run it after any `memory-engineer` or
`create-ki` change, same as any other `shared/` edit.

## Definition of Done (for a memory-engineer sweep)
- [ ] Every KI in `shared/knowledge/` and `.claude/knowledge/` was checked for duplicates against every other KI
- [ ] Any expiration candidates were flagged with a reason, not silently deleted
- [ ] `shared/memory-registry.json` still accurately lists every source (add a new entry if a new source type
      was introduced since the last sweep)
- [ ] `scripts/health-check.sh` passes the Memory Registry check
- [ ] Findings are reported for human approval before any file is deleted, merged, or moved — see Approve above

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
