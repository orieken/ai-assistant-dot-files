# Implement Automated Delivery Tier A (Interim Slice)

Implement the first interim slice of policy-driven graduated automation for the `/deliver-feature` orchestrator skill (`shared/skills/deliver-feature/SKILL.md`), following the design in `docs/aos/automated-delivery-design.md`.

This interim slice introduces Tier A (Auto-continue on green) and Tier B (Auto-retry on contract failure) directly into the skill without requiring full AOS Phase 4 runtime infrastructure.

## Target Repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` — this IS the git repo. Do NOT push.

## Context to Read First

- `docs/aos/automated-delivery-design.md` — the authoritative design doc for policy evaluation, gate classification, and escalation rules.
- `shared/skills/deliver-feature/SKILL.md` — the current orchestrator skill to be extended.
- `shared/rules/approval-gates.md` — the 8 non-negotiable approval gates.
- `shared/telemetry/event-schema.md` — event schema for telemetry logging.

## Scope of Implementation

### 1. Delivery Policy File Parsing
Extend Phase 0 of `shared/skills/deliver-feature/SKILL.md` to check for `.claude/delivery-policy.yaml`:
- If `.claude/delivery-policy.yaml` exists: read and parse policy configuration.
- If `.claude/delivery-policy.yaml` is missing: default to standard human prompt mode (100% backward compatible).

### 2. Interim Tier B Auto-Retry Loop
When `validate-artifact` returns `FAIL` for any pipeline stage:
- If `attempts < maxContractRetries` (default $N=3$): re-invoke the producing agent automatically, passing the specific contract violations list as context.
- Log event `contract.retry` to `.claude/telemetry/events.jsonl`.
- If violations persist after $N$ attempts: halt pipeline and prompt human operator.

### 3. Interim Tier A Auto-Proceed Evaluation
Replace prose `PAUSE` checkpoints for policy-eligible stages (Phase 1 discovery, developer review pass, QA/docs/devops checks) with policy evaluation:
- If stage policy `autoProceed: true` AND contract validation == `PASS` AND diff lines < `maxDiffLines`:
  - Log `policy.evaluated` event (`decision: AUTO_PROCEED`).
  - Proceed immediately to next pipeline step without pausing.
- Else:
  - Log `policy.evaluated` event (`decision: PAUSE_HUMAN`).
  - Execute standard human PAUSE prompt.

### 4. Non-Negotiable Human Gate Enforcement
Ensure Gates #1 (Friday ship), #3-4 (DB migrations), #5 (External API mutations), and #8 (Deploy) ALWAYS require explicit human confirmation regardless of policy file settings.

## Guardrails & Discipline

- Maintain 100% backward compatibility when `.claude/delivery-policy.yaml` is absent.
- Update `scripts/health-check.sh` or `check-parity.sh` if necessary to verify policy parsing.
- Commit discipline: explicit paths only, Conventional Commits.
- Do NOT push.

## Report Format (under 200 words)

```
STATUS: complete | stopped-at-<reason>

Commits landed:
  <sha> feat(deliver-feature): add delivery-policy parsing and Tier A/B automation
  <sha> test(deliver-feature): add verification test for policy auto-proceed

Health-check state: <n> passed, <n> warned, <n> failed
Parity-check state: <pass | drift>
```

Go.
