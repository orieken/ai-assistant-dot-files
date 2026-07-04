---
name: promote-memory
description: Evaluates one delivery's retrospective.md immediately after it's produced and decides whether anything in it is worth promoting to a Knowledge Item, an ADR, a rule change, or a lessons-learned entry — or explicitly rejecting. Prevents memory bloat by having real rejection criteria, not promoting by default.
triggers:
  keywords: ["promote-memory", "promote this retrospective", "is this worth remembering"]
  intentPatterns: ["Should this go into a KI", "Promote findings from *", "/promote-memory *"]
standalone: true
---

## When To Use
Right after a `retrospective.md` is produced (auto-invoked every 5th delivery per `deliver-feature`, or
manual via `/retrospective`) — evaluate that single retrospective's "What Went Poorly" / "What To Improve" /
"Patterns Identified" sections for promotion-worthy content.

Do NOT use for cross-delivery pattern detection (a finding that only becomes clear after seeing it recur
across many deliveries) — use `extract-lessons` for that; it's periodic and looks at many retrospectives at
once, this skill looks at exactly one, immediately. Do NOT use to audit the existing KI corpus for
duplicates — use `memory-engineer` for that.

## Context To Load First
1. The specific `retrospective.md` being evaluated
2. `docs/runbooks/memory-engineering.md` — Promotion Rules and the Memory Contract fields
3. `shared/knowledge/*.md` (frontmatter only, via `search-ki`'s pre-filter) — to check for an existing
   near-duplicate before recommending a new KI

## Process

### 1. Extract Candidates
Read the retrospective's "What Went Poorly," "What To Improve," and "Patterns Identified" sections. For each
distinct point, treat it as a candidate for evaluation — don't bundle multiple unrelated findings into one
candidate.

### 2. Apply Promotion Rules
For each candidate, check `docs/runbooks/memory-engineering.md`'s Promotion Rules: is it reusable (would help
a *different* future feature), non-obvious (not already inferable from `ARCHITECTURE_RULES.md`), and
actionable (a concrete "when X, do Y because Z")? If any of these fail, the candidate is a reject —
say so plainly rather than promoting a low-value entry just because it was mentioned.

### 3. Check for an Existing Match
Invoke `search-ki` (or check frontmatter directly) before recommending a new KI — if an existing KI already
covers this, recommend updating that KI instead of creating a near-duplicate.

### 4. Decide the Destination
For each surviving candidate, decide: **KI** (a reusable pattern/fix), **ADR-worthy** (a decision with
real alternatives considered, not just a fix), **Rule-change-worthy** (something that should become a
`shared/rules/` guardrail, not a one-off note), or **Lesson** (worth recording in `docs/lessons-learned/`
even though it's not yet a confirmed recurring pattern — the threshold `extract-lessons` would need to act on).

### 5. Produce the Candidate Record
Write each surviving candidate using the Memory Contract's Candidate fields (Source, Type, Evidence, Tags,
Expiration condition) from `docs/runbooks/memory-engineering.md` — this is what a human reviews before
`create-ki`/`adr`/a rule edit actually happens.

## Output Format

```markdown
# Memory Promotion: [Feature Name] retrospective

## Candidates Evaluated
- [Finding 1]: [Promote as KI / Promote as ADR / Promote as Rule change / Promote as Lesson / Reject — reason]
- [Finding 2]: ...

## Candidate Records (for surviving candidates only)

### Candidate: [short title]
- **Source**: docs/features/<name>/retrospective.md, "[section name]"
- **Type**: KI | ADR-worthy | Rule-change-worthy | Lesson
- **Evidence**: [exact quote or close paraphrase from the retrospective]
- **Tags**: [proposed frontmatter tags, if Type is KI]
- **Expiration condition**: [what would make this stop being true]
- **Existing overlap checked**: [KI name checked, found no overlap] / "None found in corpus"

## Rejected (with reason)
- [Finding]: [Why it doesn't meet the promotion bar]
— or "Nothing rejected — all candidates promoted" / "No candidates found — nothing in this retrospective met the bar"
```

## Guardrails
- **Never** promote by default — a retrospective producing zero promotable candidates is a normal, healthy
  outcome, not a failure of this skill.
- **Never** create the actual KI/ADR/rule-edit file directly — this skill produces the Candidate Record for
  human review; `create-ki`, `adr`, or a direct rule edit happens after approval, same Approve gate as any
  other `shared/` change.
- **Never** promote something already covered by an existing KI without at least checking — recommend a
  merge/update instead.

## Standalone Mode
Pure local file reads and judgment. No external calls, no embeddings.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
