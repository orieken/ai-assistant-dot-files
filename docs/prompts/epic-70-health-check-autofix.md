# Epic 70 — `health-check --fix` + `install.sh --auto-sync`

Source: `docs/audits/framework-audit-2026-08-07.md` §3 item 2.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Prior context

`scripts/health-check.sh` runs 273 checks across parity, inventory, schema, versioning,
and structural constraints. It is entirely read-only — it reports PASS/WARN/FAIL but never
remediates. When FAIL checks appear, the user must manually identify and run the right
generator or repair script.

Many FAIL-class checks are deterministic: the fix is always the same generator re-run:
- Parity drift → `bash scripts/generate-configs.sh`
- Stale CODEMAP → `bash scripts/generate-codemap.sh`
- Missing platform configs → `bash scripts/generate-configs.sh`
- Count drift in prose docs → `bash scripts/check-inventory-drift.sh` (WARNs, not fixes)

Others require human judgment:
- A missing ADR section (what to write?)
- A broken KI link (which KI was intended?)
- A new agent without a fixture (fixture content is non-deterministic)

The framework also has `install.sh` which sets up a target project; parity drift in the
installed project requires a manual re-run today.

## Scope

**Phase A — Classification (one commit, then PAUSE for user approval):**

Draft and commit as `docs(scripts): classify health-check checks for auto-fix (Epic 70 Phase A)`:

- Read `scripts/health-check.sh` in full. For every FAIL-level check (not WARN):
  - Label it AUTO-FIXABLE if the remediation is a deterministic, idempotent shell command
    already present in the repo (e.g., re-running a generator).
  - Label it HUMAN-REQUIRED if the remediation requires authoring content, making a
    structural decision, or choosing between alternatives.
- Produce `docs/runbooks/health-check-autofix-classification.md`: a table of check name,
  FAIL condition, classification, and the exact fix command for AUTO-FIXABLE ones.
- If fewer than 5 checks are AUTO-FIXABLE, halt and flag: the `--fix` feature may not
  justify its implementation cost at this check count.

**Phase B — Implementation (after approval; one commit per op):**

Op 1 — `feat(scripts): health-check --fix flag (Epic 70 Op 1)`:
- Add `--fix` flag to `scripts/health-check.sh`.
- When `--fix` is active, AUTO-FIXABLE checks that fail run their registered fix command
  and recheck once. HUMAN-REQUIRED checks that fail emit `[NEEDS HUMAN] <check-name>` and
  are skipped.
- `--fix` must not modify any file that was not flagged by a failing check.
- Final summary line: `Fixed N, Needs-human M, Still-failing K`.
- `bash scripts/health-check.sh --fix` must itself pass after fixing fixable items.

Op 2 — `feat(scripts): install.sh --auto-sync flag (Epic 70 Op 2)`:
- Add `--auto-sync` flag to `scripts/install.sh`.
- When `--auto-sync` is active: after install, run `scripts/check-parity.sh`; on any
  parity divergence, automatically invoke `scripts/generate-configs.sh` and recheck.
- Log each auto-sync action: `[AUTO-SYNC] regenerated platform configs`.
- Emit a warning (not error) if post-sync parity still diverges (hand-off to human).

Op 3 — `feat(health-check): --fix integration test (Epic 70 Op 3)`:
- Add a test case in `scripts/test-install.sh` (or a new `scripts/test-autofix.sh`) that:
  1. Deliberately introduces a parity drift (delete one generated config).
  2. Runs `health-check --fix`.
  3. Asserts the config is restored and health-check exits 0.
- Wire as an optional check: skip if the environment flag `SKIP_AUTOFIX_TEST=1` is set.

Op 4 — `docs(runbooks): document --fix and --auto-sync (Epic 70 Op 4)`:
- Update `docs/RUNBOOKS.md` with: when to use `--fix`, what it can and cannot fix, and
  the `--auto-sync` install workflow.
- Cross-reference `docs/runbooks/health-check-autofix-classification.md`.

After every commit: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- If Phase A finds fewer than 5 AUTO-FIXABLE checks, halt after Phase A and deliver
  the classification document as the sole deliverable. The `--fix` flag is not worth
  building for 1–2 fixable cases.
- If any AUTO-FIXABLE check's fix command is not idempotent (running it twice produces
  different output), reclassify it as HUMAN-REQUIRED.
- If the `--fix` flag cannot be cleanly added to `health-check.sh` without restructuring
  the check loop (cyclomatic complexity would exceed 7 in the patched version), Extract
  the remediation logic into a separate `scripts/health-check-fix.sh` script that `--fix`
  delegates to.

## Report (under 200 words)

```
Phase A commit: <sha>
Phase A findings:
  - AUTO-FIXABLE checks: <N>
  - HUMAN-REQUIRED checks: <M>
  - Halt recommended: <yes | no>

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, --fix self-test <pass>, install --auto-sync <pass>.
```

Go.
