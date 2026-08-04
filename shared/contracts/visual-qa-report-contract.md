# Contract: visual-qa-report.md

**Produced by**: visual-qa-engineer
**Consumed by**: sre-engineer, tech-writer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## Summary`
- `## Visual Regression`
- `## Heatmap Coverage`
- `## Cold Spots`
- `## Recommendations`
- `## Notes for QA`

## Validation Rules

`validate-artifact` checks presence of every heading above, plus:

- `## Summary` must contain one of: `PASS`, `FAIL`, or `UNCONFIGURED` — an empty or missing verdict is a FAIL.
- `## Heatmap Coverage` must contain `Coverage Score:` — or the literal `UNCONFIGURED` if heatmap data was absent.
- If `## Summary` contains `FAIL`, `## Recommendations` must be non-empty (cannot be "None") — a FAIL with no recommendations means the agent didn't complete its analysis.
- If `## Visual Regression` contains `diffs detected` or `screenshot diff`, at least one named scenario must appear in `## Visual Regression` — vague "diff found" without a scenario name is a FAIL.

## Pipeline Behavior

This artifact is conditional — `visual-qa-engineer` is skipped on non-UI features or when neither heatmap data nor visual baselines are present. When skipped, no contract validation step runs.

A `FAIL` verdict blocks the pipeline (same as `qa-report.md`'s `Failed: 0` rule). An `UNCONFIGURED` verdict allows the pipeline to proceed — it signals that Saturday heatmap instrumentation hasn't been adopted yet, not that something is broken.

## Retrieval Frontmatter (WARN)

Pipeline artifacts should include a YAML frontmatter block at the very top of the file. Missing or incomplete retrieval frontmatter triggers a **WARN** from `validate-artifact` — not a FAIL. Existing archived artifacts without frontmatter are unaffected.

```yaml
---
feature: "<feature-name>"             # kebab-case slug derived from the feature file name
bounded_context: "<context>"          # owning bounded context (from DOMAIN_DICTIONARY.md domain list)
domain_terms: []                      # canonical terms from DOMAIN_DICTIONARY.md used in this feature
files_touched: []                     # repo-relative paths of files created or modified
issue_refs: []                        # ticket/issue references (e.g., PROJ-123, #456)
linked_adrs: []                       # repo-root-relative paths to referenced ADRs
linked_kis: []                        # repo-root-relative paths to referenced Knowledge Items
---
```

Once frontmatter adoption is visible across a project's feature archive, this check will be promoted to FAIL in a future release.
