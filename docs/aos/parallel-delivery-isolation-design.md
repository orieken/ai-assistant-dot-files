# Parallel-Delivery Workspace Isolation — Design (Epic 63 Phase A)

**Status**: Complete — Phase A design approved; Phase B implemented in commits `10c41bd`→`77ad560` (v3.3.8).
**Source**: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b / `docs/prompts/done/epic-63-parallel-delivery-isolation.md`

---

## 1. Ruling: Isolation Model

**Decision: Option (a) — named workspace-per-feature**

Workspace path becomes `.claude/feature-workspace/<feature-name>/` (slugified, matching the
feature name already used as the key in `pipeline-state.json`).

**Rationale:**

Option (b) (git-worktree-per-feature) gives stronger filesystem-level isolation but introduces
worktree lifecycle management that must be surfaced to the user, and any tool or script that
derives the repo root via `git rev-parse --show-toplevel` or relative `../../` paths would silently
resolve to the wrong directory when run from inside a worktree that is not at the canonical
project root. The `shared/` and `docs/` reads that agents make during delivery are all relative
to the project root — those reads would still work inside a worktree, but the accidental discovery
cost of the first tool that breaks is high relative to the gain.

Option (a) is sufficient for the primary unlock: two features (or two operators) concurrently
in-flight no longer block each other at the `pipeline-state.json` singleton. Staged-file
collisions in the shared git index remain theoretically possible when two deliveries modify the
same source file, but: (i) different features rarely edit the same source file; (ii) when they
do, git itself surfaces the conflict at commit time; (iii) `deliver-feature` Gate #2
(git-commit) is already a human gate and that human reviews `git diff` before confirming — so
the conflict is caught, not silently overwritten.

Option (c) (both-layered) adds complexity without solving a problem option (a) leaves open.

---

## 2. Backward Compatibility — Detection and Migration

Any existing install that has a delivery in-flight under the old singleton layout must resume
cleanly. Detection rule:

```
IF .claude/feature-workspace/pipeline-state.json EXISTS
   AND .claude/feature-workspace/ contains no <feature-name>/ subdirectories
   → LEGACY singleton workspace detected
```

Migration behavior (added to deliver-feature Phase 0, step 3):

1. Read `.claude/feature-workspace/pipeline-state.json` to extract the `feature` field.
2. Move the flat workspace to `.claude/feature-workspace/<feature>/` — rename the directory
   in place (or create the subdirectory, move all files, and remove the now-empty root-level
   artifacts).
3. Proceed normally from the new path. Log a `workspace.migrated` telemetry event.
4. If the `feature` field is absent from the legacy state file, use `default` as the feature
   name and inform the user.

This migration is non-destructive (move, not copy-then-delete) and idempotent: if the
subdirectory already exists, the migration step is skipped.

---

## 3. Cross-Cutting Inventory

Total files requiring operational edits in Phase B: **20**

> ⚠️ **Escalation flag**: this exceeds the 15-file threshold defined in the epic's escalation
> clause. Phase B must be sequenced carefully — see §5 for the recommended commit order.

### 3a. Skills (5 files)

Each needs: hardcoded `.claude/feature-workspace/` replaced with a reference to the active
workspace path (resolved from the feature name passed by `deliver-feature`). The path-resolution
convention is defined in `deliver-feature/SKILL.md` and passed to each agent/skill via the
orchestration prompt preamble.

| File | Change needed |
|---|---|
| `shared/skills/deliver-feature/SKILL.md` | Add workspace-path resolution section (Phase 0); update every path reference (~12 occurrences) |
| `shared/skills/resume-pipeline/SKILL.md` | Accept feature-name parameter; resolve workspace path the same way; update state/history paths |
| `shared/skills/pipeline-trace/SKILL.md` | Resolve `pipeline-trace.json` from named workspace, not singleton |
| `shared/skills/context-engineer/SKILL.md` | Write `context-manifest.md` to named workspace path |
| `shared/skills/validate-artifact/SKILL.md` | Resolve all 13 artifact paths from named workspace |

### 3b. Agents (15 files)

Each agent's pre-read block and write instruction hardcodes `.claude/feature-workspace/<artifact>`.
The change is uniform: prefix the workspace path variable. Phase B will use a sed-equivalent
batch update against the common path pattern, then verify each file individually.

| Agent | Artifacts referenced |
|---|---|
| `shared/agents/analyst.md` | `context-manifest.md`, `analysis.md` |
| `shared/agents/architect.md` | `analysis.md`, `architecture-notes.md`, `rfc-*.md` |
| `shared/agents/developer.md` | `analysis.md`, `context-manifest.md`, `implementation-notes.md` |
| `shared/agents/code-reviewer.md` | `analysis.md`, `architecture-notes.md`, `implementation-notes.md`, `code-review-report.md` |
| `shared/agents/security-reviewer.md` | `analysis.md`, `architecture-notes.md`, `implementation-notes.md`, `code-review-report.md`, `security-report.md` |
| `shared/agents/performance-engineer.md` | `analysis.md`, `architecture-notes.md`, `performance-report.md` |
| `shared/agents/data-engineer.md` | `analysis.md`, `architecture-notes.md`, `data-engineering-notes.md` |
| `shared/agents/accessibility-engineer.md` | `analysis.md`, `implementation-notes.md`, `accessibility-report.md` |
| `shared/agents/sre-engineer.md` | `analysis.md`, `implementation-notes.md`, `observability-report.md` |
| `shared/agents/qa-engineer.md` | `analysis.md`, `implementation-notes.md`, `qa-report.md` |
| `shared/agents/visual-qa-engineer.md` | `qa-report.md`, `visual-qa-report.md` |
| `shared/agents/tech-writer.md` | `analysis.md`, `implementation-notes.md`, `qa-report.md`, `docs-report.md` |
| `shared/agents/devops-engineer.md` | `analysis.md`, `implementation-notes.md`, `docs-report.md`, `devops-report.md` |
| `shared/agents/privacy-auditor.md` | globs `*.md` in workspace |
| `shared/agents/context-auditor.md` | `context-manifest.md`, downstream artifacts |

### 3c. Documentation (not blocking Phase B)

The following files reference `feature-workspace` in explanatory prose only and do not drive
runtime behavior. They need updates in a docs-only pass but are not on the Phase B critical path:

`docs/ARCHITECTURE.md`, `docs/THREAT_MODEL.md`, `docs/aos/automated-delivery-design.md`,
`docs/aos/governance-pairs.md`, `docs/aos/migration-guide.md`, `docs/aos/migration-plan.md`,
`docs/patterns/enterprise-integration-patterns.md`, `docs/runbooks/context-engineering.md`,
`shared/DOMAIN_DICTIONARY.md`, `shared/orchestration/audit-composition-pattern.md`,
`shared/orchestration/interface.md`, `shared/orchestration/pipeline-schema.md`,
`shared/platform-registry.json`, `shared/templates/agent.template.md`,
`shared/workflows/feature-delivery-workflow.md`, `shared/hooks/scheduled-monthly.yaml`,
and ~20 additional skill files whose references are illustrative examples, not execution paths.

---

## 4. Concurrency Constraints

### Safe (no change needed)

| Concern | Assessment |
|---|---|
| `pipeline-state.json` collision | Resolved by named workspaces — each feature has its own file |
| `pipeline-trace.json` collision | Same — per named workspace |
| `delivery-policy.yaml` | Read-only config singleton; concurrent reads are safe |
| Friday Gate #1 | Independent per delivery; two concurrent deliveries each prompt their own ship gate |
| `context-manifest.md`, `analysis.md`, all artifacts | Per named workspace after Phase B |

### Acceptable risk (document, no fix)

| Concern | Assessment |
|---|---|
| `events.jsonl` concurrent appends | Append-only JSONL; two in-flight runs write interleaved lines. Each event carries `pipeline_id`, so queries remain correct. Individual JSON line corruption from concurrent byte-level writes is theoretically possible but extremely unlikely in interactive Claude usage (sequential LLM turns, not parallel OS processes). Accept the risk; add a note to `shared/telemetry/README.md`. |

### Real constraint (Phase B must document)

| Concern | Assessment |
|---|---|
| Git index sharing | With option (a), two concurrent deliveries share one git working tree and one index. `deliver-feature` Gate #2 (git-commit) is a human gate — the human reviewing `git diff` will see staged hunks from both deliveries if they overlap. **Constraint**: two concurrent deliveries MUST NOT attempt Gate #2 simultaneously. Users running parallel deliveries must coordinate their commit gates sequentially. Document this in the deliver-feature SKILL.md concurrency section and the runbook. |
| Test runner ports/tmp | Two test suite runs in the same working tree can collide on ports or `tmp/` artifacts. This is user-space and outside framework scope — document as a known constraint, not something the framework fixes. |

---

## 5. Phase B Commit Sequence

(Pending human approval of this Phase A design)

```
1. feat(aos): add workspace-path resolution helper to deliver-feature (Epic 63)
   → shared/skills/deliver-feature/SKILL.md — add §"Workspace Path Resolution", update all paths

2. feat(aos): isolate resume-pipeline to named workspace (Epic 63)
   → shared/skills/resume-pipeline/SKILL.md

3. feat(aos): isolate pipeline-trace to named workspace (Epic 63)
   → shared/skills/pipeline-trace/SKILL.md

4. feat(aos): isolate context-engineer and validate-artifact to named workspace (Epic 63)
   → shared/skills/context-engineer/SKILL.md, shared/skills/validate-artifact/SKILL.md

5. feat(aos): update all agent pre-read blocks to named workspace (Epic 63)
   → shared/agents/*.md (15 files, batch)

6. docs(aos): update prose references to named workspace convention (Epic 63)
   → docs/*, shared/DOMAIN_DICTIONARY.md, shared/templates/, shared/workflows/

7. docs(aos): add concurrency constraints and telemetry note (Epic 63)
   → shared/telemetry/README.md, docs/runbooks/context-engineering.md
```

After every commit: `bash scripts/health-check.sh` must pass green.

Final verification: two scratch deliveries with distinct feature names progressing through
Phase 0 + one artifact each, confirming no cross-contamination of `pipeline-state.json`,
`context-manifest.md`, or `.history/` entries between the two runs.

---

## 6. Fitness Functions

Per `shared/rules/architecture-guardrails.md` §7, every architectural decision must produce a
measurable constraint:

- **Pipeline-orchestrator grep**: the following command must return zero matches after Phase B —
  it scopes to the four orchestrator skills and the pipeline agents that are directly invoked
  by `deliver-feature` or `deliver-atdd`:

  ```
  grep -r '\.claude/feature-workspace/[^<]' \
    shared/skills/deliver-feature/ \
    shared/skills/deliver-atdd/ \
    shared/skills/resume-pipeline/ \
    shared/skills/pipeline-trace/ \
    shared/skills/context-engineer/ \
    shared/skills/validate-artifact/ \
    shared/skills/orchestrate/ \
    shared/agents/analyst.md \
    shared/agents/architect.md \
    shared/agents/developer.md \
    shared/agents/code-reviewer.md \
    shared/agents/security-reviewer.md \
    shared/agents/performance-engineer.md \
    shared/agents/data-engineer.md \
    shared/agents/accessibility-engineer.md \
    shared/agents/sre-engineer.md \
    shared/agents/qa-engineer.md \
    shared/agents/visual-qa-engineer.md \
    shared/agents/tech-writer.md \
    shared/agents/devops-engineer.md \
    shared/agents/privacy-auditor.md \
    shared/agents/context-auditor.md \
    shared/agents/context-engineer.md
  ```

  Two expected exceptions exist in `deliver-feature/SKILL.md`'s "Workspace Path Resolution"
  section: the two lines that describe the legacy singleton detection (`look for a flat
  .claude/feature-workspace/pipeline-state.json`) and the migration step (`Move all files
  from .claude/feature-workspace/ into`). These intentionally reference the old flat path;
  they are not artifact paths that need to change.

  Standalone skills (e.g., `five-whys`, `adr`, `threat-model`) and non-pipeline agents
  (e.g., `chaos-engineer`, `release-manager`) use `.claude/feature-workspace/` as a scratch
  area for ad-hoc output and are intentionally excluded from this check — they are invoked
  without a feature-name context and do not participate in `deliver-feature`'s concurrency model.

- **Health check**: existing `scripts/health-check.sh` must remain green at every commit.
- **Legacy path**: manual verification that a delivery started with the pre-Phase-B layout can
  resume after the workspace migration step runs.

---

*Epic 63 — complete. Phase B shipped 2026-08-02 in commits `10c41bd`→`77ad560`, tagged `v3.3.8`.*
