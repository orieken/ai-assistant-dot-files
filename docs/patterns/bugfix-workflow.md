# Bugfix Workflow Pattern

The `deliver-bugfix` skill and the incident-memory pipeline together form a closed production
feedback loop. This pattern describes how a production incident becomes a targeted bugfix and
then a permanent improvement to the framework's rule and knowledge base.

## When to use `deliver-bugfix` vs `deliver-feature`

| Signal | Route to |
|---|---|
| Known symptom, approximate location, no architectural change | `deliver-bugfix` |
| Root cause unknown, needs analyst exploration | `deliver-feature` |
| Fix touches ≥3 files across bounded contexts | `deliver-feature` (escalate) |
| Migration or contract change required | `deliver-feature` (escalate) |
| New behavior needs to be added alongside the fix | `deliver-feature` (split first) |
| Regression in a recently delivered feature | `deliver-bugfix` (check incident record) |

## Full loop from incident to improvement

```
1. Production incident fires
   /on-call — timeline, blast radius, preliminary fix
           ↓
2. Root cause analysis
   /five-whys — why chain, root cause, recommended action
           ↓
3. Incident record persisted
   docs/incidents/<YYYY-MM-DD>-<slug>.md
   - Affected Feature: docs/features/<name>/  ← cross-reference
   - Candidate Records (promote-memory format)
           ↓
4. Bug filed and fixed
   /new-feature → type=bug → /deliver-bugfix
   docs/features/<bug-slug>/
   - characterization test (Phase 1 — Reproduce)
   - implementation-notes.md (Phase 2 — Fix)
   - code-review-report.md (Phase 3 — Review)
   - qa-report.md (Phase 4 — Verify)
   docs/incidents/<slug>.md  ← "Fixed by" link added (Phase 5 — Record)
           ↓
5. Learning extraction (periodic)
   /extract-lessons
   - Step 6: incident-feature pair → "which pipeline stage missed this?"
   - Proposes rule/prompt change (human-gated)
   - Incident Candidate Records promoted via same gate as feature lessons
           ↓
6. Framework improves
   rule or prompt change → prevents the same class of bug in future deliveries
```

## The reproduce-first rule (Michael Feathers)

`deliver-bugfix` enforces this strictly:

1. Write a characterization test that **fails** against the current code.
2. Confirm the test fails with the *expected* error — not a different failure.
3. Only then implement the fix.
4. Confirm the characterization test now **passes** and no prior tests regressed.

A fix without a red test first is not a reproducible fix — it is a guess. If the characterization
test passes without changes, the bug has not been reproduced; investigate further before Phase 2.

## Escalation rule

If at any point during `deliver-bugfix` the developer discovers:

- The fix requires changes in ≥3 files across distinct bounded contexts
- A contract or schema change is needed
- A database migration is required
- The "fix" is actually new behavior that was missing

**Stop. Write a summary. Route to `/deliver-feature`.** Never scope-creep a bugfix into a
feature silently. The characterization test written in Phase 1 remains — it becomes the
acceptance criterion seed for the full feature delivery.

## Artifacts and retrieval

Bugfix artifacts land in `docs/features/<bug-slug>/` — the same archive as full feature
deliveries. This is deliberate:

- `/retrospective` sees bugfixes alongside features (important signal: did the fix close the
  incident fully?)
- `/extract-lessons` mines bugfix artifacts for recurring code-review patterns
- `analyst` and `context-engineer` can surface related bugfixes when starting a feature in
  the same bounded context

The incident record at `docs/incidents/<date>-<slug>.md` separately preserves the production
signal (severity, blast radius, timeline, five-whys chain) that the feature artifacts don't
carry.

## Relationship to existing skills

| Skill | Relationship |
|---|---|
| `deliver-feature` | Full pipeline; bugfix escalates here if scope exceeds one bounded context |
| `on-call` | Triggers the incident record; bugfix links back to it via "Fixed by" |
| `five-whys` | Provides root cause for the incident record; runs before or after bugfix |
| `extract-lessons` | Mines incident-feature pairs; consume bugfix artifacts + incident record |
| `promote-memory` | Consumes Candidate Records from incident records unchanged |
| `retrospective` | Sees bugfix artifacts from `docs/features/<bug-slug>/` same as features |

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
