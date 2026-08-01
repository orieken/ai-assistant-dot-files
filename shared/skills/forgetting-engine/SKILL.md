---
name: forgetting-engine
description: Identifies and flags obsolete or superseded lessons and Knowledge Items for expiration. Opposing-force pair with learning-engine. In AOS Phase 3, scheduled monthly via shared/hooks/scheduled-monthly.yaml (opt-in, disabled by default).
triggers:
  keywords: [forgetting engine, expire lesson, archive ki, decay memory]
  intentPatterns: ["identify obsolete lessons", "run forgetting engine"]
standalone: true
---

## When To Use

Use when:
- Sweeping `docs/lessons-learned/` or `shared/knowledge/` for obsolete rules superseded by newer architecture.
- Flagging items for draft expiration.
- Invoked automatically by the `scheduled-monthly` hook (AOS Phase 3, opt-in).

Do NOT use when:
- Learning new lessons from retrospectives (use `learning-engine`).
- Auditing KI schema compliance (use `memory-auditor`).

## Invocation Modes

| Mode | Trigger | Scope |
|---|---|---|
| Manual | Human runs `/forgetting-engine` | Full sweep of knowledge + lessons-learned |
| Scheduled | `scheduled-monthly` hook | Full sweep, configurable staleness threshold |

## Staleness Criteria

An item is a candidate for expiration if it meets **any** of the following:

1. **ADR supersession**: an ADR in `docs/adrs/` explicitly supersedes or invalidates it.
2. **Temporal staleness**: `last-referenced` frontmatter (if present) is > `stalenessThresholdMonths` months ago **and** no file in `shared/knowledge/`, `docs/features/`, or `docs/adrs/` links to it within the same window.
3. **Obsolete framework version**: the item references a framework pattern, tool, or convention that no longer exists in the current codebase (e.g., references a removed agent, a renamed skill, or a deprecated config key).

Items flagged by criteria 1 are high-confidence expiration candidates. Items flagged by 2 or 3 are moderate-confidence and should be reviewed carefully before approving.

## Context To Load First

1. `shared/knowledge/README.md` — corpus overview
2. `shared/memory-registry.json` — which sources are tracked and their metadata
3. `shared/knowledge/*.md` — KI corpus
4. `docs/lessons-learned/*.md` — lessons corpus (if directory exists)
5. `docs/adrs/*.md` — to identify supersession relationships

## Process

1. **Scan corpus**: Glob `shared/knowledge/*.md` and `docs/lessons-learned/*.md`. Read `last-referenced` frontmatter where present.

2. **Identify staleness signals**: for each item, check each staleness criterion. Tag each finding with `HIGH`, `MEDIUM`, or `LOW` confidence.

3. **Check for active links**: for each staleness-flagged item, grep `shared/knowledge/`, `docs/features/`, and `docs/adrs/` for its filename or title. If linked anywhere within the staleness window, downgrade confidence or remove the flag.

4. **Draft expiration proposal** in `.claude/feature-workspace/proposed-expirations.md`:

```markdown
# Proposed Expirations: YYYY-MM-DD

## Summary
N items flagged for expiration (H high-confidence, M medium, L low).

## High-Confidence Candidates (superseded by ADR)

- `shared/knowledge/ki-obsolete.md` — Superseded by ADR-002 (explicit statement in ADR body)
  - [ ] Approve

## Medium-Confidence Candidates (stale > 6 months, no links)

- `docs/lessons-learned/lessons-2024-01-01.md` — Not referenced in 8 months; no linking files
  - [ ] Approve
  - [ ] Defer (keep for another 3 months)

## Low-Confidence / Ambiguous

- `shared/knowledge/ki-possibly-outdated.md` — References removed agent "old-agent" but content may still be relevant
  - [ ] Approve
  - [ ] Keep (update instead of expire)
```

5. **Pause for human confirmation**: Present the draft and ask explicitly:
   > "The forgetting engine found N expiration candidates. Reply 'approve all', 'approve <list>', 'defer <list>', or 'skip' to decide."
   
   Do NOT proceed to step 6 without explicit confirmation.

6. **Execute approved expirations**:
   - Move approved `shared/knowledge/` items to `shared/knowledge/archive/` (create if missing).
   - Move approved `docs/lessons-learned/` items to `docs/lessons-learned/archive/` (create if missing).
   - Do NOT delete files — archive only.
   - Update `shared/memory-registry.json` if the expired item had a registered source entry.

## Output Format

See Draft Expiration Proposal template in step 4.

## Guardrails

- **Never delete files.** Archive only — move to `archive/` subdirectory.
- **Never archive without explicit human confirmation** ("approve expiration" or "confirm"). Non-negotiable regardless of invocation mode.
- If invoked via scheduled hook and no candidates are found, emit "no expirations proposed" and exit cleanly.
- Low-confidence candidates must surface a "Keep (update instead of expire)" option.
- If `archive/` directory does not exist, create it before moving the first file.

## Opposing Force

`learning-engine` is this skill's opposing force. It adds new lessons; this skill removes stale ones.
Together they maintain the quality of `docs/lessons-learned/` and `shared/knowledge/` over time.

## Hook Configuration

See `shared/hooks/scheduled-monthly.yaml` for the opt-in monthly schedule. Copy to `.claude/hooks/`
and set `enabled: true` under `forgetting-engine-monthly` to activate automatic scheduling.

## Standalone Mode

Operates purely using local file editing tools (Read, Write, Edit, Glob, Bash for grep).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 3 Runtime layer. CC BY 4.0.*
