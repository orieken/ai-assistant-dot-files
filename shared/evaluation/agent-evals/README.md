# Agent Eval Records

Regression records written by `scripts/run-agent-evals.sh` (Epic 61). One Markdown file per
agent per run. The harness reads these files to determine whether a new run is a regression,
stable, or improved versus the previous baseline.

---

## Schema

Each record file is named `<agent>-eval-<YYYY-MM-DD>.md` and contains these machine-readable
fields (bold key, value on same line — bash `grep` extracts them without a parser):

```markdown
# Agent Eval: <agent> — <YYYY-MM-DD>

**Agent version**: <semver from shared/agents/<agent>.md frontmatter>
**Model used**: <concrete model ID, e.g. claude-sonnet-4-6>
**Fixture**: tests/agents/<agent>/<input-file>
**Run mode**: full | pattern-only | no-judge

## Pattern Grade

(Per-pattern PASS/FAIL lines from terminal run — written for human readability.)

**Pattern overall**: PASS | FAIL | SKIP

## Rubric Grade

(Per-criterion PASS/FAIL lines from the Haiku judge — written for human readability.)

**Rubric overall**: PASS | FAIL | SKIP

## Regression Delta

Compared against: <previous record file name, or "no baseline — first recorded eval">

**Overall delta**: REGRESSION | STABLE | IMPROVED | BASELINE
```

### Field semantics

| Field | Type | Values | Used for |
|---|---|---|---|
| `Agent version` | string | semver | Tracking prompt version under test |
| `Model used` | string | model ID | Identifying which model produced the output |
| `Fixture` | string | repo-relative path | Traceability |
| `Run mode` | string | `full` \| `pattern-only` \| `no-judge` | Determines which grades are meaningful |
| `Pattern overall` | enum | `PASS` \| `FAIL` \| `SKIP` | **Regression detection input** |
| `Rubric overall` | enum | `PASS` \| `FAIL` \| `SKIP` | **Regression detection input** |
| `Overall delta` | enum | see below | **Regression detection output** |

### Delta values

| Value | Meaning |
|---|---|
| `BASELINE` | No prior record exists for this agent — first eval, no comparison possible |
| `STABLE` | Both pattern and rubric at same grade as previous (or SKIP vs SKIP) |
| `REGRESSION` | Previous=PASS, current=FAIL on pattern OR rubric — harness exits 1 |
| `IMPROVED` | Previous=FAIL, current=PASS on pattern OR rubric (human review encouraged) |

### Regression rule

A REGRESSION is flagged whenever **a dimension that previously passed now fails**:

```
prev Pattern overall = PASS  AND  current Pattern overall = FAIL  →  REGRESSION
prev Rubric overall  = PASS  AND  current Rubric overall  = FAIL  →  REGRESSION
```

`SKIP` fields (from `--pattern-only` or no API key) are never treated as regressions — a
missing grade cannot regress. A REGRESSION causes `run-agent-evals.sh` to exit 1 and lists
the affected agents on stderr.

---

## Retention and multiple runs per day

If the harness is run more than once on the same day, the second run overwrites the first
(same filename). Keep the most recent run for that day. For a permanent historical archive,
run with a timestamp in a CI artifact rather than committing all records.

Records committed to this directory serve as the **regression baseline** — the harness reads
the most recent file per agent (sorted by filename date descending). Do not delete a record
without a replacement unless you intend to reset the baseline for that agent.

---

## Relationship to `agent-eval` skill

The `agent-eval` skill (interactive, single-agent) writes to `docs/agent-metrics/evals/`
today. Both locations act as valid sources of regression baselines. The harness only reads
`shared/evaluation/agent-evals/`; the skill will be updated in a follow-on to also write
here so both sources converge on the same baseline pool.

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
