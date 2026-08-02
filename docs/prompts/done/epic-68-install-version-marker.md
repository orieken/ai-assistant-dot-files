# Epic 68 — Framework-Version Marker in Installed Projects

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #8 — smallest). The gap:
`install.sh` writes no version record into installed projects, so detecting drift between an
install and upstream is manual archaeology.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

- `install.sh` supports `--global` (symlink) and `--project <path>` (copy) modes across 10+
  platforms, is idempotent, prints a verification summary — but records nothing about WHAT
  version it installed.
- Symlink installs track upstream automatically (drift is only uncommitted upstream changes);
  copy installs freeze at install time — those are the ones that rot silently. The marker matters
  most for copy mode, but write it in both (mode is part of the record).
- `docs/prompts/done/update-installed-framework.md` — the update flow that currently
  auto-detects which of 4 install patterns apply by inspecting the filesystem; a marker makes
  that detection reliable instead of forensic.
- Repo versions: git tags (v3.0.0 latest tagged; v3.1.0 pending human tag), `shared/agents/CHANGELOG.md`.

## Scope: 3 ops

**Op 1 — Write the marker.**
`install.sh` writes `.claude/framework-install.json` into the target (project mode) or the
global config root (global mode): source repo path/URL, git tag + commit SHA at install time,
ISO date, mode (symlink/copy), platforms installed, flags used. Update `uninstall.sh` to remove
it (and `test-install.sh` to assert it exists post-install with valid JSON).
Commit: `feat(install): write framework-install.json version marker (Epic 68 Op 1)`

**Op 2 — Teach the readers.**
- `health-check.sh`: when run in a directory that has a marker but is NOT the framework repo
  itself, add an INFO/WARN line — installed version vs. source repo's current HEAD/tag (only
  when the source path in the marker resolves; SKIP silently otherwise — no false noise on
  machines where the source repo moved).
- `docs/prompts/update-installed-framework.md` (create an updated copy in active prompts if only
  the `done/` version exists — check first): marker becomes detection input #1, filesystem
  forensics the fallback.
Commit: `feat(health-check): report install-vs-upstream drift from marker (Epic 68 Op 2)`

**Op 3 — Docs.**
README install section + `docs/MIGRATION.md`: the marker's existence, format, and that pre-marker
installs are detected by the legacy forensic path (no forced migration).
Commit: `docs(install): document version marker (Epic 68 Op 3)`

After every op: `bash scripts/health-check.sh` green in THIS repo (which is not an install —
verify the marker check doesn't misfire here), and `bash scripts/test-install.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If a marker convention already half-exists (grep install.sh/uninstall.sh for any state file
  first), extend it rather than adding a second — halt only on conflict.
- If global-mode marker placement is ambiguous across platforms, marker project-mode only and
  note the deferral.

## Report (under 100 words)

```
Commits: <sha> x3
Marker path + fields: <summary>
test-install: <pass, marker asserted>
health-check in framework repo: <no misfire confirmed>
health-check in scratch install: <drift line shown when versions differ>
```

Go.
