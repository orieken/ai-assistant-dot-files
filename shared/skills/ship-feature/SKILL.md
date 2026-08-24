---
name: ship-feature
description: Automates branch creation, Conventional Commit assembly, PR compilation, and optional release tagging — with human approval gates at every irreversible step.
triggers:
  keywords: ["create pr", "open pr", "pull request", "release tag", "ship feature", "ship this"]
  intentPatterns: ["Ship * to a PR", "Create a pull request for *", "Open a PR for *", "Tag a release for *", "I'm done with *, ship it", "/ship-feature *"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use when the user is ready to move completed local work into version control and open a PR. Works
standalone against any git branch — no prior `deliver-feature` run is required (see Standalone Mode
for graceful fallback when feature artifacts are absent).

Do NOT use when:
- The user only wants to run tests, review code, or generate docs — use the relevant dedicated skill.
- The current branch is `main` or `master` and no `--force-main` flag is given — the skill refuses to commit directly to the default branch.
- There are uncommitted changes the user hasn't reviewed yet — surface them and ask before staging.

## Context To Load First
1. `shared/rules/approval-gates.md` — gate #2 (commit) and gate #5 (external API mutations, including PR creation and `gh release create`).
2. `docs/features/<feature-name>/` if it exists — spec summary, `qa-report.md`, `retrospective.md`.
3. Root `CLAUDE.md` § Git Commits — Conventional Commits format rules.
4. `git status` and `git log --oneline -10` — understand current branch and pending changes.

## Process

### Phase 0 — Discover
1. Run `git status` to surface any uncommitted changes. If unexpected untracked or modified files exist,
   list them and ask the user to confirm they belong in this commit before staging anything.
2. Determine `<feature-name>`: if invoked as `/ship-feature <name>`, use that. Otherwise derive from
   the current branch name (`feat/<name>` → `<name>`), or from the most recent `docs/features/*/`
   directory modified in the current session.
3. Detect the default branch: `git remote show origin | grep 'HEAD branch'` (fallback: check for
   `main` then `master`). Store as `<default-branch>`.
4. Derive remote URL and repo slug: `git remote get-url origin` + `gh repo view --json nameWithOwner`
   (no hardcoded values).

### Phase 1 — Branch
5. Check the current branch name. If it is `<default-branch>`, create and switch to
   `feat/<feature-name>` (or `fix/<feature-name>` if the user signals a fix). Announce: "Created
   branch `feat/<feature-name>` from `<default-branch>`."
6. If already on a non-default branch, confirm the branch name with the user before proceeding (a
   one-line acknowledgment is fine — no gate halt needed for this step; it is reversible).

### Phase 2 — Commit (Gate #2)
7. Identify files to stage: prefer explicit paths from `implementation-notes.md`'s "Files Modified"
   section if present; otherwise list all modified tracked files from `git status` and ask the user
   to confirm the staging list. **Never run `git add -A` or `git add .`.**
8. Draft the Conventional Commit message:
   - First line: `<type>(<scope>): <imperative subject under 72 chars>`
   - Body (if needed): explain *why*, not *what*; reference the feature name or ticket.
   - Types: `feat`, `fix`, `refactor`, `test`, `chore`, `docs` (from root `CLAUDE.md`).
9. **HALT — Gate #2.** Present:
   ```
   Gate #2 — Git Commit
   Files to stage:
     <list of explicit paths>
   Commit message:
     <type>(<scope>): <subject>

     <body if any>
   Approval word: "commit" or "approve commit"
   ```
   Wait for explicit "commit" or "approve commit" before running any `git add` or `git commit`.
10. On approval: stage the listed paths, create the commit, show the resulting `git log --oneline -1`.
10a. Check whether `<branch>` exists on origin: `git ls-remote --exit-code origin <branch>`. If it
    does not, run `git push -u origin <branch>` and report the push before the Gate #5 halt.
    Never push to `<default-branch>` — only the feature branch.

### Phase 3 — PR (Gate #5)
11. Compile the PR body. Priority order for source material:
    - `docs/features/<feature-name>/delivery-summary.md` → summary paragraph.
    - `docs/features/<feature-name>/qa-report.md` → acceptance criteria coverage table (first 10 rows).
    - `docs/features/<feature-name>/retrospective.md` → link to retrospective.
    - Fallback (no feature artifacts): `git log <default-branch>..HEAD --oneline` as the body (see Standalone Mode).
12. Draft PR title: `<type>(<scope>): <subject>` matching the commit, ≤ 70 characters.
13. **HALT — Gate #5.** Present:
    ```
    Gate #5 — PR Creation (external API mutation)
    Title: <title>
    Body:
    ---
    <full compiled body>
    ---
    Command that will run:
      gh pr create --title "<title>" --body "<body>"
    Approval word: "send" or "approve request"
    ```
    Wait for explicit "send" or "approve request" before running `gh pr create`.
14. On approval: run `gh pr create --title "<title>" --body "$(cat <<'EOF' ... EOF)"`. Print the
    resulting PR URL. If `gh` is not installed, print the command for the user to run and note
    degraded mode (see Standalone Mode).

### Phase 4 — Release (optional, `--release` flag only)
15. If `--release` was NOT passed, stop here and report success.
16. Invoke the `release-manager` agent to analyze `git log` since the last tag, determine the semver
    bump, and draft release notes. Read the agent's output before continuing.
17. **HALT — Gate #5.** Present:
    ```
    Gate #5 — Release Tag + GitHub Release (external API mutation)
    Proposed tag:    v<major>.<minor>.<patch>
    Release notes:
    ---
    <release-manager output>
    ---
    Commands that will run:
      git tag v<major>.<minor>.<patch>
      gh release create v<major>.<minor>.<patch> --notes "<notes>"
    Approval word: "send" or "approve request"
    ```
    Wait for explicit "send" or "approve request".
18. On approval: run `git tag`, then `gh release create`. Print the resulting release URL.

## Output Format

At the end of a successful run, print a summary:

```
## ship-feature report

Branch:  feat/<feature-name>  (created from <default-branch> | already existed)
Commit:  <sha> <subject>  [Gate #2 — approved]
PR:      <url>             [Gate #5 — approved]
Release: <url> | skipped   [Gate #5 — approved | --release not passed]

Standalone fallback: <"used git-log body — no feature artifacts found" | "used delivery-summary.md">
```

## Guardrails

- **Never commit to the default branch directly.** If the current branch is `<default-branch>`, create
  a feature branch first. Refuse entirely if `--force-main` is not passed AND the branch is the
  default branch.
- **Gate #2 halts are non-negotiable.** Any edit to the staged file list or commit message after
  presenting the gate resets the gate — re-present the full updated view and wait again.
- **Gate #5 halts are non-negotiable.** Any edit to the PR title, body, or `gh` command after
  presenting the gate resets the gate. Gates CANNOT be bypassed by this skill or any caller.
- **Never run `git add -A` or `git add .`** — stage only explicit paths confirmed by the user.
- **Never hardcode remote names, repo slugs, or branch names** — derive all of them from `git remote`
  and `gh repo view` at runtime.
- **Never push directly to the default-branch remote.** Pushing the feature branch to origin is
  required before `gh pr create` (handled in step 10a) — that is the only permitted push.
- **Every gate halt must name the gate number and quote the exact approval word it is waiting for.**
- **`git tag` is local and reversible; `gh release create` is an external mutation (Gate #5).** They
  are always presented together in a single Gate #5 halt — approving one implicitly approves the other
  in the same response.

## Standalone Mode

This skill degrades gracefully when external tools or feature artifacts are absent:

| Missing resource | Degraded behavior |
|---|---|
| `docs/features/<feature-name>/` | PR body falls back to `git log <default-branch>..HEAD --oneline`. Report notes "standalone fallback: no feature artifacts found". |
| `gh` CLI not installed | Skip `gh pr create`; print the exact command for the user to run. Note: "gh not available — copy the command above to open the PR manually." Same for `gh release create`. |
| `release-manager` agent unavailable | Halt with: "release-manager unavailable — provide version bump and release notes manually, then re-run with `--release`." |
| `git remote` absent (local-only repo) | Halt with: "No remote found — push the repo first, then re-run /ship-feature." |

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
