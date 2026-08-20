# Docs Cleanup — Annotate THREAT_MODEL.md with Implementation Status

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/THREAT_MODEL.md` was produced by Epic 65 (Op 1 of 2). It still presents itself as an Op 1
document and labels its mitigations as "Op 2 candidates" — but multiple Epic 65 Op 2 controls
have since been implemented:

- `shared/rules/memory-trust-boundary.md` — KI/ADR-as-instructions rule (shipped)
- Spec-ingestion caution block in `shared/agents/analyst.md` (shipped)
- Sync provenance (`sync_commit_sha`) on pulled KIs via `sync-memory.sh` (shipped)
- Hook example defaults to `enabled: false` in `shared/hooks/` (shipped)

The file should be annotated to reflect this — not rewritten — so readers can tell what is live
versus what remains aspirational.

## Scope

**Op 1 — Verify which Op 2 controls are present:**

Check each claimed implementation:
```bash
ls shared/rules/memory-trust-boundary.md
grep -n "spec.*ingestion\|untrusted input\|spec content" shared/agents/analyst.md | head -10
grep -n "sync_commit_sha" scripts/sync-memory.sh 2>/dev/null | head -5
grep -n "enabled.*false\|false.*enabled" shared/hooks/*.yaml 2>/dev/null | head -5
```

Note the result of each check before making any changes.

**Op 2 — Update the document banner:**

The opening note currently reads "Op 1 of 2. Mitigations that require rule or schema changes
are gated at Op 2 — findings here are candidates." Update it to reflect that Op 2 controls
have been implemented. Suggested replacement:

> Op 1 + Op 2 complete. Controls marked "Implemented" below are live as of Epic 65.
> Remaining "proposed" items require future work.

Adjust wording to accurately reflect what Op 1 check found above.

**Op 3 — Annotate individual mitigation blocks:**

For each "Proposed mitigation (Op 2 candidate)" block:
- If the control is verified present (Op 1 check above): change the label to
  `**Implemented (Epic 65 Op 2)**` and add a one-line file reference.
- If the control is NOT present or partially present: leave the existing label and add
  `**Status: Not yet implemented**` so the gap is explicit rather than implied.

Do not change the threat descriptions or risk ratings — only the mitigation status labels.

**Op 4 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Update `docs/THREAT_MODEL.md`" item as `[x]`.

## Guardrails

- Conventional commit: `docs(threat-model): annotate Op 2 implementation status`
- Stage only `docs/THREAT_MODEL.md` and `docs/TODO.md`.
- Do not rewrite threat descriptions, STRIDE categories, or risk ratings.
- Only change status labels and the banner. This is an annotation pass, not a content rewrite.

## Escalation

Stop and report if:
- The Op 1 verification finds that fewer than 3 of the 4 claimed controls are present — the
  banner update wording will need adjustment and human review before committing.
- The document structure has changed substantially from what the TODO describes.

## Report

On completion, confirm:
- Which of the four controls were verified as present
- How many mitigation blocks were updated to "Implemented" vs. left as "Not yet implemented"
- The updated banner text
- Commit hash
