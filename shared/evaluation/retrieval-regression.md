# Retrieval Regression Set

Format for capturing real-query retrieval regression cases for the lexical tier
(`search-ki`, `query-memory`). Cases are proposed by `retrieval-evaluator` from
telemetry and approved by a human before entering this file — same discipline as
`learning-engine`'s draft-KI flow.

This graduates ADR-002's judgment-only retrieval fitness function toward mechanical
verification: using queries people actually asked, not synthetic benchmarks.

---

## Telemetry dependency (PENDING — schema extension required)

Auto-population from telemetry requires a `retrieval.queried` event type not yet
in `shared/telemetry/event-schema.md` (schema v1.1.0 as of Epic 60). Until that
event type is added and emitting, `retrieval-evaluator` must propose cases manually
from memory (i.e., from prior known query failures surfaced in evaluation reports).

**Proposed schema extension** (for human review before implementing — schema changes
ripple to all telemetry consumers):

```json
{
  "timestamp": "2026-08-04T00:00:00.000Z",
  "event_type": "retrieval.queried",
  "agent_or_skill_name": "search-ki",
  "artifact_path": null,
  "outcome": "hit | miss",
  "metadata": {
    "query": "the literal query string",
    "hits": 3,
    "top_hit": "shared/knowledge/some-ki.md",
    "chosen": "shared/knowledge/some-ki.md",
    "corpus": "ki | adr | feature-archive | domain-dictionary"
  }
}
```

Fields needed for regression set automation: `query` (required), `hits` (required),
`chosen` (optional — the result the agent actually loaded). A `hits: 0` case is
auto-candidate for a missing-KI finding.

---

## Case format

```markdown
### Case: <short-slug>

**Query**: "<exact query string>"
**Corpus**: ki | adr | feature-archive | feature-archive-summaries | domain-dictionary
**Must appear in top-5**: <file path>
**Acceptable alternatives**: <file path>, <file path>  (optional)
**Source**: telemetry:<ISO-date> | manual:<ISO-date>
**Notes**: why this case matters (what would break if it regressed)
```

---

## Running the regression set

`retrieval-evaluator` uses these steps when asked to "run the retrieval regression set":

1. Read this file and load all approved cases (those with a `### Case:` heading).
2. For each case, invoke the appropriate retrieval skill (`search-ki` or `query-memory`)
   with the recorded query.
3. Check whether the "Must appear in top-5" reference is in the result set.
4. Record PASS / FAIL per case.
5. Report results in the standard output format below.

The evaluator never modifies this file — all case additions require human approval.

---

## Output format

```markdown
# Retrieval Regression Run: [YYYY-MM-DD]

## Summary
- Cases run: N
- Passed: N
- Failed: N

## Failures
| Case | Query | Expected | Got (top-5) |
|---|---|---|---|
| <slug> | "<query>" | <expected-file> | <actual top-5 list, or "no hits"> |

## Passes
(omit if all pass — keep output short)
| Case | Query | Expected |
|---|---|---|

## Proposed new cases (from telemetry misses)
(only populated if retrieval.queried events are available)
- Query "<query>" returned 0 hits — candidate for create-ki or tag update
```

---

## Seed cases (manually added — 0 from telemetry, schema extension pending)

No seed cases yet. Once the `retrieval.queried` event type is added to
`shared/telemetry/event-schema.md` and emitting, `retrieval-evaluator` can propose
cases from actual zero-hit queries. Until then, cases can be added manually by
the framework team when a known retrieval failure is identified.

### Example (reference — not a real case, shows format)

```markdown
### Case: context-engineer-wiring

**Query**: "why isn't context-engineer running automatically before analyst"
**Corpus**: ki
**Must appear in top-5**: shared/knowledge/context-engineer-must-be-wired-into-pipeline.md
**Source**: manual:2026-08-04
**Notes**: This KI was written specifically because the pipeline had a dead-capability
  bug. Regression here would mean the KI became unretrievable, defeating its purpose.
```

---

## Governance

- Cases in this file are approved — `retrieval-evaluator` runs them as-is.
- Cases under "Proposed new cases" in an evaluator report are drafts — a human must
  copy them here (adding source/notes) before they count as regression tests.
- Removing a case requires a commit with a comment explaining why the case is no
  longer relevant (e.g., the KI was deprecated and removed).

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
