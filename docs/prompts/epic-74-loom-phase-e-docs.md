# Epic 74 Phase E — loom CLI: Documentation Updates

Parent epic: `docs/prompts/epic-74-loom-cli.md`
Prerequisite: Phase D (`epic-74-loom-phase-d-release.md`) must be complete.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Expected state coming in

Phases A–D are complete. `loom` builds and passes lint. goreleaser config and the GitHub
Actions release workflow are committed. The binary is not yet published (no `v*` tag pushed).
`install.sh` is still the documented primary install path everywhere.

## What to build

### Op 1 — Update root `README.md`

Read `README.md` in full before editing.

Add an **Installation** section near the top (after the project title/description, before any
other sections). Make it the primary install path:

```markdown
## Installation

```bash
# macOS / Linux
brew install orieken/tap/loom
loom install

# Windows
winget install orieken.loom
loom install
```

`loom install` detects which AI platforms are present in your project (Claude Code, Cursor,
Windsurf, etc.) and installs the framework to each one automatically.

### Manual install (repo contributors)

Clone the repo and run `./install.sh` — see [Contributing](CONTRIBUTING.md) for the full setup.
```

Also update any existing "Getting Started" or "Setup" sections that currently describe
`./install.sh` as the primary path — demote them to "For contributors" or "Advanced" without
deleting them.

Commit Op 1 as: `docs(readme): add loom as primary install path`

---

### Op 2 — Update `docs/ONBOARDING.md`

Read `docs/ONBOARDING.md` in full before editing.

Replace the "clone and run install.sh" quick-start steps with the loom path:

```markdown
## Quick start

1. Install loom: `brew install orieken/tap/loom` (macOS/Linux) or `winget install orieken.loom` (Windows)
2. In your project root: `loom install`
3. Open your AI tool — agents, skills, and rules are ready.
```

Keep the existing repo-clone instructions under a clearly labeled **"Contributing to the framework"**
or **"Advanced setup"** heading — do not delete them.

Commit Op 2 as: `docs(onboarding): make loom the primary install path`

---

### Op 3 — Create `cmd/loom/README.md`

Write a focused reference document for the loom binary itself. This is the file a user lands
on when they navigate to `cmd/loom/` in the repo.

Required sections:

**Overview** — one paragraph: what loom is, that it embeds the framework, distributed via Homebrew.

**Subcommand reference** — a table for each subcommand with its flags:

```markdown
### loom install

Install the framework to all detected AI platforms in the current project.

| Flag | Default | Description |
|---|---|---|
| `--target <path>` | `.` | Project root to install into |
| `--platform <name>` | (all detected) | Install for one platform only |
| `--stack <list>` | (all rules) | Comma-separated language stacks |
| `--copy` | `false` | Use copies instead of symlinks |
| `--with-configs` | `false` | Also install linter fitness-function configs |
| `--with-mcp` | `false` | Scaffold MCP server source into the target |
| `--dry-run` | `false` | Print planned actions without writing files |
```

(Repeat the table pattern for `uninstall`, `health`, `version`.)

**Platform detection** — list which marker file/directory triggers detection for each platform.

**Embedded content versioning** — explain that the framework version is baked into the binary
at build time; users upgrade the framework by upgrading loom via `brew upgrade loom`.

**Building from source** — for contributors:
```bash
go build -o loom ./cmd/loom
./loom install --dry-run
```

Commit Op 3 as: `docs(loom): add cmd/loom/README.md subcommand reference`

---

### Op 4 — Update `docs/prompts/README.md` to mark Epic 74 as shipped

When all prior ops are committed, update the prompts README:
- Strike through the Epic 74 row
- Add commit hash range for Phase A through E
- Add a row in the Completed Prompts table

Also move the five phase prompt files to `docs/prompts/done/`:
```bash
mv docs/prompts/epic-74-loom-phase-a-scaffold.md docs/prompts/done/
mv docs/prompts/epic-74-loom-phase-b-install.md docs/prompts/done/
mv docs/prompts/epic-74-loom-phase-c-subcommands.md docs/prompts/done/
mv docs/prompts/epic-74-loom-phase-d-release.md docs/prompts/done/
mv docs/prompts/epic-74-loom-phase-e-docs.md docs/prompts/done/
mv docs/prompts/epic-74-loom-cli.md docs/prompts/done/
```

Commit Op 4 as: `docs(prompts): mark Epic 74 shipped, move loom prompts to done/`

---

## Guardrails

- Each Op is a separate commit.
- Do not delete any existing documentation — only demote or supplement.
- `install.sh` must remain documented as the contributor/CI path in both README and ONBOARDING.
- Do not invent features or behaviors not implemented in Phases A–D — document only what exists.

## Escalation

Stop and report if:
- `README.md` or `ONBOARDING.md` have a significantly different structure than expected,
  making the section placement ambiguous — describe the structure and ask before editing.

## Report

On completion:
- Which sections in README.md were added/changed
- Which sections in ONBOARDING.md were added/changed
- Commit hash for each Op
- Confirmation that all 6 loom prompt files are in `docs/prompts/done/`
