# Runbook: Scaling Cross-Feature Learning

`context-engineer` currently finds prior deliveries in the same bounded context by grepping every
`docs/features/*/analysis.md` for a matching `**Owning Context**` entry (see
`shared/agents/context-engineer.md`, step 6). That's the right amount of engineering for a small archive —
don't build the index below until you've actually felt the grep approach strain. This runbook is here for
when that day comes.

## When to build this
Don't build it preemptively. Build it when one of these starts being true:
- The feature archive (`docs/features/`) has grown past **roughly 15-20 delivered features** — a grep across
  that many `analysis.md` files on every `context-engineer` run is still fast, but it's the point where a
  direct index starts being clearly worth the upkeep.
- You notice **false positives**: two unrelated bounded contexts happen to share enough vocabulary that the
  `**Owning Context**` grep starts matching things it shouldn't (e.g. "Identity" appears as a tangential
  concern in a `billing` feature's analysis, not just in real `Identity & Access` features).
- `context-engineer` visibly slows down or a delivery's context-manifest generation starts timing out.

If none of these are true yet, the current grep-based step needs no changes — adding the index earlier just
means one more artifact to keep in sync for a problem you don't have yet.

## What to build
A single file, `docs/features/_index-by-context.md`, structured as one section per bounded context:

```markdown
# Feature Index by Bounded Context

Maintained automatically by deliver-feature's Phase 4 (persistence step) — do not hand-edit; edits will be
overwritten on the next delivery. See docs/runbooks/scaling-cross-feature-learning.md for why this exists.

## Identity & Access
| Feature | Delivered | Key Lesson (from retrospective.md) |
|---|---|---|
| [user-auth](user-auth/) | 2026-03-12 | Missed the user-enumeration edge case on first pass |
| [password-reset](password-reset/) | 2026-05-02 | No retrospective.md for this one |

## Billing
| Feature | Delivered | Key Lesson (from retrospective.md) |
|---|---|---|
| [invoice-generation](invoice-generation/) | 2026-04-01 | Underestimated the Stripe webhook retry complexity |
```

Bounded context names must exactly match what `DOMAIN_DICTIONARY.md` and `analysis.md`'s `**Owning
Context**` field use — this is a lookup table, not a search index, so exact string matching is what makes it
fast. If a feature's owning context is misspelled or uses a synonym, it silently won't be found; this is
exactly the kind of drift `check-ubiquitous-language` (`shared/skills/check-ubiquitous-language/`) already
exists to catch — run it periodically against `docs/features/*/analysis.md` once this index exists.

## How to wire it in
1. Update `deliver-feature/SKILL.md`'s Phase 4 persistence step (alongside the existing `docs/features/README.md`
   update) to add or update this feature's row under its bounded context's section in
   `docs/features/_index-by-context.md`. Pull the "Key Lesson" cell from the feature's `retrospective.md`
   "What Went Poorly" or "What To Improve" section if one exists at delivery time (it usually won't yet,
   since `/retrospective` runs after `deliver-feature` completes — see step 3 below for the update path).
2. Update `retrospective`'s (`shared/skills/retrospective/SKILL.md`) output step to also update that
   feature's row in the index with the real lesson once the retrospective is written — the index's "Key
   Lesson" column will otherwise stay blank for every feature until someone runs `/retrospective` on it.
3. Update `shared/agents/context-engineer.md`'s step 6 to read `docs/features/_index-by-context.md` and look
   up the current bounded context's section directly, instead of grepping every `analysis.md`. Keep the grep
   as a fallback for any feature not yet in the index (e.g. delivered before the index existed, or delivered
   between index updates).
4. Bump `context-engineer`'s version, add a `shared/agents/CHANGELOG.md` entry, and add a parity/health check
   (`scripts/health-check.sh`) verifying the index's bounded-context section names still match
   `DOMAIN_DICTIONARY.md`'s actual bounded contexts — this is exactly the kind of drift that silently breaks
   a lookup table over time.

## What NOT to do
- Don't make this a vector/embeddings index. The whole point of a bounded-context index is exact-match
  lookup on a structured field that's already being written into every `analysis.md` — that's cheaper and
  more precise than semantic search for this specific problem. If you later find yourself wanting fuzzy
  "find features kind of like this one" matching across *different* bounded contexts, that's a different
  problem than this runbook solves, and the `search-ki` skill's own guardrails against premature embeddings
  (see its SKILL.md) apply here too.
- Don't drop the grep-based fallback once the index exists. A feature delivered in the gap between two index
  updates (or before the index was introduced) should still be findable.
