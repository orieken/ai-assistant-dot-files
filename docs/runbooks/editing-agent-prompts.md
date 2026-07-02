# Runbook: Editing Agent Prompts

Agent prompts are code here — every behavior change needs a version bump, a changelog entry, and ideally a
test run before you trust it. This runbook is the full workflow; see
[CONTRIBUTING.md](../CONTRIBUTING.md) for adding a brand-new agent instead of editing an existing one.

## 1. Decide the version bump

`shared/agents/CHANGELOG.md` defines the scheme (semantic-ish, not strict SemVer):
- **Patch** (`1.0.x`): wording/clarity fixes that don't change behavior.
- **Minor** (`1.x.0`): new process step, new output section, expanded guardrail — additive, backward
  compatible.
- **Major** (`x.0.0`): changed output contract (update the matching `shared/contracts/` file too, if one
  exists), removed/renamed a process step, or changed tool access.

If you're not sure which one, err toward the larger bump — a downstream consumer (another agent,
`validate-artifact`, a human reading the changelog) should never be surprised that a "patch" silently changed
behavior.

## 2. Make the edit

Edit `shared/agents/<name>.md` directly (never a generated config — see [CONTRIBUTING.md](../CONTRIBUTING.md)).
Bump the `version:` frontmatter field to match your decision in step 1.

## 3. Update the contract, if this agent has one

If `shared/contracts/<name>-contract.md` exists and your edit changed the agent's Output Format (added,
removed, or renamed a required section), update the contract in the same commit. A stale contract means
`validate-artifact` will pass artifacts that no longer match what the agent actually produces, or fail ones
that do.

## 4. Add a `shared/agents/CHANGELOG.md` entry

Add a row under a new (or the current day's) dated heading:
```markdown
## 2026-MM-DD — <one-line reason for the change>

| Agent | Version | Change |
|---|---|---|
| your-agent | 1.0.0 -> 1.1.0 | <what changed and why, one sentence> |
```
This is required in the same commit if the pre-commit hook (`scripts/hooks/pre-commit`) is enabled — see
step 6.

## 5. Test it

If `tests/agents/<name>/` has a fixture:
1. Invoke the agent against `tests/agents/<name>/input-*` in a live Claude Code session (agent invocation
   can't be scripted — see `tests/agents/README.md` for why).
2. Save the output to `tests/agents/<name>/actual-output.md`.
3. Run `scripts/test-agents.sh` — it checks the output against `expected-patterns.txt` and the agent's
   contract (if one exists).

If there's no fixture yet and this agent is likely to regress silently on future edits, consider adding one
(see [CONTRIBUTING.md](../CONTRIBUTING.md), "Adding a new agent," step 6 — the same guidance applies to
significant edits of an existing agent).

## 6. Enable the pre-commit gate (optional, recommended)

The hook isn't wired up automatically — git doesn't track its own hooks directory by default. To enable it:
```bash
git config core.hooksPath scripts/hooks
```
or symlink it into the standard location:
```bash
ln -s ../../scripts/hooks/pre-commit .git/hooks/pre-commit
```
Once enabled, any staged `shared/agents/*.md` change is blocked unless its version was bumped *and*
`shared/agents/CHANGELOG.md` is staged with an entry mentioning that agent's name.

## 7. Run the parity check

```bash
scripts/check-parity.sh
```
An agent description or name change needs to propagate to every platform's persona roster — this confirms
it did (or flags it if `generate-configs.sh` hasn't been re-run yet).

## Common mistakes this workflow exists to catch
- Editing a generated config file directly instead of `shared/agents/` (silently overwritten next generation).
- Changing an agent's Output Format without updating its `shared/contracts/` entry (breaks `validate-artifact`
  silently — it'll either falsely fail or falsely pass future runs).
- Bumping the version without a changelog entry, or vice versa (the pre-commit hook, once enabled, catches
  both — but only if you've enabled it per step 6).
