---
name: memory-engineer
description: Periodic curator of the durable memory corpus (Knowledge Items, primarily) — audits shared/knowledge/ and .claude/knowledge/ for duplicates and overlaps, flags stale or superseded KIs for expiration, and keeps shared/memory-registry.json accurate. Distinct from context-engineer (retrieves memory for one task) and create-ki (writes one KI on request).
triggers:
  keywords: ["memory-engineer", "audit memory", "memory sweep", "expire knowledge item", "dedupe knowledge"]
  intentPatterns: ["Audit the knowledge base", "Are there duplicate KIs", "Run a memory sweep", "/memory-engineer *"]
standalone: true
---

## When To Use
- Periodic (manual, like `agent-scorecard`/`extract-lessons` — not auto-invoked by `deliver-feature`):
  "run a memory sweep," "audit the KI corpus for duplicates."
- After a burst of `create-ki`/`promote-memory` activity, to check the corpus hasn't accumulated
  near-duplicates.
- When `search-ki` or `query-memory` results start feeling noisy or redundant — a symptom of exactly the
  problem this skill exists to catch.

Do NOT use to write a new KI — use `create-ki` directly for that (this skill audits the existing corpus, it
doesn't author new entries, other than merging duplicates it finds). Do NOT use for a single retrospective's
promotion decision — use `promote-memory` for that (this skill operates on the whole corpus, periodically,
not on one delivery's output right after it's produced).

## Context To Load First
1. `shared/memory-registry.json` — every registered memory source and its metadata
2. `shared/knowledge/*.md` and `.claude/knowledge/*.md` (if the latter exists) — the actual KI corpus
3. `docs/runbooks/memory-engineering.md` — the lifecycle and promotion/expiration criteria this skill enforces

## Process

### 1. Inventory
List every KI in `shared/knowledge/` and `.claude/knowledge/`. Cross-check this list against
`shared/memory-registry.json`'s `knowledge-items` entry — if a KI directory exists that the registry doesn't
know about (e.g. a new project added `.claude/knowledge/` for the first time), note it, but registry
maintenance is a human-approved edit (see step 4), not something this skill silently rewrites.

### 2. Duplicate and Overlap Check
Read every KI's full body (not just frontmatter — two KIs can use completely different words for the same
underlying pattern, the same way `search-ki`'s own judgment-based matching works). For any pair that
substantially overlaps:
- If one is a strict superset of the other, recommend merging the narrower one into the broader one.
- If they cover genuinely different angles of a related topic, recommend cross-referencing them (`[[other-ki-name]]`
  style link) rather than merging — don't force a merge that would lose a real distinction.

### 3. Expiration Candidates
Apply `docs/runbooks/memory-engineering.md`'s Expiration Criteria: does the KI reference code/agents/patterns
that no longer exist? Has it been superseded by a newer KI? Flag each candidate with the specific reason —
never flag "this looks old" without a concrete check (e.g. grep for the file/agent/pattern it references and
confirm it's actually gone).

### 4. Registry Accuracy
Compare `shared/memory-registry.json` against actual current sources. If a source's `paths` no longer match
reality (a directory was renamed, a new memory-relevant directory was added), propose the specific JSON diff
— but do not write it directly; this is a human-approved edit like anything else in `shared/`.

### 5. Report
Produce the sweep report (see Output Format). Nothing gets deleted, merged, or rewritten by this skill
directly — every finding is a recommendation for a human to approve, exactly like `code-reviewer`'s findings
require developer action rather than auto-applying.

## Output Format

```markdown
# Memory Sweep: [YYYY-MM-DD]

## Inventory
- Total KIs: [N] (`shared/knowledge/`: [N], `.claude/knowledge/`: [N])
- Registry accuracy: [Matches current sources | Diff needed — see below]

## Duplicate/Overlap Findings
- [KI A] + [KI B]: [Merge recommended — B is a strict subset of A] / [Cross-reference recommended — related but distinct]
— or "No duplicates found"

## Expiration Candidates
- [KI name]: [specific reason — what was checked and found gone/superseded]
— or "No expiration candidates found"

## Registry Diff Needed (if any)
```diff
[proposed change to shared/memory-registry.json]
```
— or "Registry is accurate, no change needed"

## Recommended Actions (require human approval before executing)
- [ ] [Specific action — merge X into Y / move Z to shared/knowledge/expired/ / update registry]
```

## Guardrails
- **Never** delete, merge, or move a KI file without explicit human approval — this skill reports findings,
  it doesn't act on them unilaterally, consistent with `code-reviewer` and every other review-shaped skill in
  this framework.
- **Never** flag an expiration candidate without a concrete, checked reason (a grep result, a superseding
  KI's name) — "this feels old" is not a finding.
- **Never** merge two KIs that cover genuinely different angles just because they share a topic — losing a
  real distinction is worse than tolerating two related files.
- Expiration moves a file to `shared/knowledge/expired/` with a reason noted, never a bare `rm` — the
  reasoning that led to writing a KI that turned out to be wrong is itself worth keeping.

## Standalone Mode
Pure local file reads (registry, KI frontmatter and bodies) and judgment-based comparison. No external
calls, no embeddings — matches `search-ki`'s own no-embeddings guardrail; this is the same reasoning applied
to comparing KIs against each other instead of against a query.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
