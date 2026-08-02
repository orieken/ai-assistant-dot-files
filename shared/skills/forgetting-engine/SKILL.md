---
name: forgetting-engine
description: Identifies and flags obsolete or superseded lessons and Knowledge Items for expiration, and audits the capability inventory for duplicate skills, keyword collisions, and agent+skill name pairs. Opposing-force pair with learning-engine. In AOS Phase 3, scheduled monthly via shared/hooks/scheduled-monthly.yaml (opt-in, disabled by default).
triggers:
  keywords: [forgetting engine, expire lesson, archive ki, decay memory, capability inventory, skill overlap, keyword collision, duplicate skill, inventory audit]
  intentPatterns: ["identify obsolete lessons", "run forgetting engine", "audit capability inventory", "find duplicate skills *", "check skill keyword overlap"]
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
| Capability Inventory | Human runs `/forgetting-engine capability-inventory` or asks about duplicate/overlapping skills | Scans `shared/agents/` + `shared/skills/` for collisions, description overlap, and dual name-pairs |

## Staleness Criteria

An item is a candidate for expiration if it meets **any** of the following:

1. **ADR supersession**: an ADR in `docs/adrs/` explicitly supersedes or invalidates it.
2. **Temporal staleness**: `last-referenced` frontmatter (if present) is > `stalenessThresholdMonths` months ago **and** no file in `shared/knowledge/`, `docs/features/`, or `docs/adrs/` links to it within the same window.
3. **Obsolete framework version**: the item references a framework pattern, tool, or convention that no longer exists in the current codebase (e.g., references a removed agent, a renamed skill, or a deprecated config key).

Items flagged by criteria 1 are high-confidence expiration candidates. Items flagged by 2 or 3 are moderate-confidence and should be reviewed carefully before approving.

## Context To Load First

### KI / Lesson Sweep Mode
1. `shared/knowledge/README.md` — corpus overview
2. `shared/memory-registry.json` — which sources are tracked and their metadata
3. `shared/knowledge/*.md` — KI corpus
4. `docs/lessons-learned/*.md` — lessons corpus (if directory exists)
5. `docs/adrs/*.md` — to identify supersession relationships

### Capability Inventory Mode
1. `shared/schemas/agent-frontmatter.schema.json` — field definitions (status, superseded_by)
2. `shared/schemas/skill-frontmatter.schema.json` — field definitions
3. `shared/agents/*.md` — all agent frontmatter (name, description, status, superseded_by)
4. `shared/skills/*/SKILL.md` — all skill frontmatter (name, description, triggers, status, superseded_by)

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

## Capability Inventory Audit

Run this mode when asked about duplicate skills, keyword collisions, or agent+skill name pairs. It produces PROPOSALS only — never applies changes automatically.

### Process

**CI.1 — Read all frontmatter.** Glob `shared/agents/*.md` (skip `CHANGELOG.md`) and `shared/skills/*/SKILL.md`. Extract `name`, `description`, `triggers.keywords`, and `status` for each.

**CI.2 — Detect keyword collisions.** Two or more skills share the same keyword in `triggers.keywords`. Flag as a collision. Assess severity:
  - Same keyword + similar description → HIGH: strong merge candidate.
  - Same keyword + different scopes → MEDIUM: may need disambiguation in descriptions or keyword narrowing.

**CI.3 — Detect description overlap.** Compare descriptions pairwise. Flag pairs whose descriptions convey the same job in different words (≥ 2 of 3 of: same domain, same action verb class, same audience). Assess:
  - Same domain + same action → HIGH merge candidate.
  - Same domain + complementary actions → LOW: keep both, document difference.

**CI.4 — Detect agent+skill name pairs.** List names that appear in both `shared/agents/*.md` and `shared/skills/*/SKILL.md`. For each pair, assess:
  - **Intentional wrapper**: the skill's description says it "invokes" or "wraps" the agent → mark WRAPPER.
  - **Accidental**: descriptions diverge or the skill duplicates the agent's behavior → mark COLLISION.

**CI.5 — Draft inventory proposal** in `.claude/feature-workspace/proposed-inventory-changes.md`:

```markdown
# Capability Inventory Proposal: YYYY-MM-DD

## Keyword Collisions

| Keyword | Skills Sharing It | Severity | Recommendation |
|---|---|---|---|
| complexity | analyze-complexity, complexity-check | HIGH | merge into analyze-complexity |

## Description Overlap Pairs

| Capability A | Capability B | Overlap Reason | Recommendation |
|---|---|---|---|
| analyze-complexity | complexity-check | same domain + same action verb class | merge |

## Agent+Skill Name Pairs

| Name | Agent file | Skill file | Assessment | Recommendation |
|---|---|---|---|---|
| spec-writer | shared/agents/spec-writer.md | shared/skills/spec-writer/SKILL.md | WRAPPER — skill description says "invokes the spec-writer agent" | document as wrapper convention |

## Summary
N collision(s), M overlap pair(s), P name pair(s).
Recommended: X merge(s), Y deprecation(s), Z documentation-only changes.
```

**CI.6 — Pause for human review.** Present the proposal and ask:
> "Capability inventory scan found N issues. Reply with which items to act on (e.g., 'merge analyze-complexity + complexity-check', 'document spec-writer as wrapper', 'skip all') and I'll implement only the approved ones."

Do NOT apply any change without explicit approval.

## Output Format (KI / Lesson Sweep)

See Draft Expiration Proposal template in step 4.

## Guardrails

- **Never delete files.** Archive only — move to `archive/` subdirectory.
- **Never archive without explicit human confirmation** ("approve expiration" or "confirm"). Non-negotiable regardless of invocation mode.
- If invoked via scheduled hook and no candidates are found, emit "no expirations proposed" and exit cleanly.
- Low-confidence candidates must surface a "Keep (update instead of expire)" option.
- If `archive/` directory does not exist, create it before moving the first file.
- **Capability inventory proposals are PROPOSALS only.** Never rename, delete, or modify agent or skill files without explicit per-item human approval from the CI.6 confirmation step. The scan is diagnostic, not prescriptive.

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
