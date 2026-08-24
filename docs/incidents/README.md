# Incident Records

Permanent incident records written by `shared/skills/on-call/SKILL.md` and
`shared/skills/five-whys/SKILL.md` (Epic 67). One Markdown file per production incident.

## Naming convention

```
docs/incidents/<YYYY-MM-DD>-<kebab-slug>.md
```

Example: `docs/incidents/2026-08-04-payment-webhook-timeout.md`

## Schema

See the "Incident Record Schema" section in `shared/skills/five-whys/SKILL.md` for the canonical
field definitions. Key fields:

| Field | Purpose |
|---|---|
| `Affected Feature` | Link to `docs/features/<name>/` when the incident traces to a delivered feature — this is the cross-reference `extract-lessons` uses to mine incident-feature pairs |
| `Five Whys Chain` | Root cause chain (populated by `/five-whys` session; "Pending" if session has not run) |
| `Candidate Records` | `promote-memory`-format candidates for KI/rule/lesson promotion (human-gated via existing approval machinery) |

## Relationship to the memory pipeline

```
Production incident
  → /on-call (response + timeline)
  → /five-whys (root cause chain)
  → docs/incidents/<date>-<slug>.md  ← this directory
        ↓
  extract-lessons mines incident-feature pairs ("which pipeline stage should have caught this?")
        ↓
  promote-memory gates promote-worthy Candidate Records → KIs / rules / lessons-learned
```

Incident records are registered in `shared/memory-registry.json` as a source type. The
`docs/incidents/` path is indexed lexically the same way `docs/features/` is.

## Retention

Records are permanent — never deleted, only updated with a "Status: Resolved" line when the
incident closes. If a record linked a feature that is later superseded, update the `Affected
Feature` link to point to the successor.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
