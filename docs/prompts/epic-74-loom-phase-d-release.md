# Epic 74 Phase D — loom CLI: goreleaser + Homebrew Release Pipeline

Parent epic: `docs/prompts/epic-74-loom-cli.md`
Prerequisite: Phase C (`epic-74-loom-phase-c-subcommands.md`) must be complete.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Expected state coming in

Phase C shipped `loom uninstall`, `loom version` (with ldflags vars), and `loom health`.
The binary builds cleanly and lint passes. Version vars are set to `"dev"` for local builds.

## ⚠️ Human action required before this phase can fully complete

The GitHub Actions release workflow needs a `HOMEBREW_TAP_GITHUB_TOKEN` secret — a GitHub PAT
with `repo` write scope on the Homebrew tap repo. This must be added to the repository's
Actions secrets by the repo owner before the first release tag push will succeed. Flag this
at the end of the report and do not block on it — wire the pipeline regardless; the secret
can be added later.

## What to build

### Op 1 — goreleaser configuration

Create `.goreleaser.yaml` at the repo root. Before writing it, run:

```bash
git tag --sort=-v:refname | head -3   # confirm latest version tag format
ls shared/VERSION 2>/dev/null          # confirm VERSION file from Phase C
```

**`.goreleaser.yaml`:**

```yaml
version: 2
project_name: loom

before:
  hooks:
    - go mod tidy

builds:
  - id: loom
    main: ./cmd/loom
    binary: loom
    env:
      - CGO_ENABLED=0
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X github.com/orieken/loom/cmd/loom/cmd.version={{.Version}}
      - -X github.com/orieken/loom/cmd/loom/cmd.commit={{.Commit}}
      - -X github.com/orieken/loom/cmd/loom/cmd.date={{.Date}}

archives:
  - id: loom
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"

brews:
  - name: loom
    repository:
      owner: orieken
      name: homebrew-tap   # VERIFY: confirm this matches the actual tap repo name
    commit_author:
      name: loom-release-bot
      email: noreply@github.com
    homepage: https://github.com/orieken/loom
    description: "AI assistant framework installer — agents, skills, and rules for Claude, Cursor, and more"
    license: MIT
    test: |
      system "#{bin}/loom", "version"
    install: |
      bin.install "loom"

release:
  github:
    owner: orieken
    name: ai-assistant-dot-files
  draft: false
  prerelease: auto
```

**IMPORTANT**: verify the tap repo name before committing. Run:
```bash
# Ask the user or check GitHub — the tap repo is where testsmith's formula lives.
# Common names: homebrew-tap, homebrew-orieken, tap
```

If the tap repo name cannot be confirmed, write `homebrew-tap` as a placeholder and flag it
prominently in the report.

Commit Op 1 as: `chore(loom): add goreleaser configuration`

---

### Op 2 — GitHub Actions release workflow

Create `.github/workflows/loom-release.yml`. Pin ALL action versions to full commit SHAs —
never use mutable tags like `@v4`. Look up current SHAs:

```bash
# Find current SHAs for the actions you need at their latest stable release
# actions/checkout, actions/setup-go, goreleaser/goreleaser-action
# Check https://github.com/<owner>/<repo>/releases for the latest tag, then:
# git ls-remote https://github.com/<owner>/<repo> refs/tags/<tag>
```

**`.github/workflows/loom-release.yml`:**

```yaml
name: Release loom
on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    name: Release loom binary
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@<SHA>   # pin to full SHA
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@<SHA>   # pin to full SHA
        with:
          go-version-file: go.mod

      - name: Run goreleaser
        uses: goreleaser/goreleaser-action@<SHA>   # pin to full SHA
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

Commit Op 2 as: `chore(loom): add GitHub Actions release workflow`

---

### Op 3 — Verify goreleaser config locally

Run goreleaser in snapshot mode to confirm the config is valid without publishing:

```bash
# Install goreleaser if not present
brew install goreleaser

# Snapshot build (no publish, no tag required)
goreleaser release --snapshot --clean

# Verify artifacts were produced
ls dist/
# Should see: loom_*_darwin_arm64.tar.gz, loom_*_linux_amd64.tar.gz, etc.
```

If goreleaser reports config errors, fix them before committing. Do not commit a broken config.

This op produces no commit — it is verification only.

---

## Guardrails

- Action SHAs must be pinned to full 40-character commit hashes — see
  `shared/rules/iac-conventions.md` GitHub Actions guardrails.
- The `GITHUB_TOKEN` permissions block must be present — `contents: write` only.
- Do not add `pull_request_target` triggers.
- Do not store any credentials in the workflow file — only `${{ secrets.* }}` references.
- `golangci-lint run ./...` must still pass after adding `.goreleaser.yaml` (it should be
  unaffected, but verify).

## Escalation

Stop and report if:
- The tap repo name cannot be confirmed — use `homebrew-tap` as placeholder and flag loudly.
- `goreleaser release --snapshot` fails with config errors that cannot be resolved.
- The `go.mod` module path set in Phase A does not match the ldflags `-X` paths — fix the
  ldflags to match the actual module path before committing.

## Report

On completion:
- The tap repo name used (confirmed or placeholder)
- Output of `goreleaser release --snapshot` (last 20 lines)
- List of artifact files in `dist/` after snapshot
- The full SHA for each pinned GitHub Action
- Commit hash for Op 1 and Op 2
- ⚠️ Reminder: `HOMEBREW_TAP_GITHUB_TOKEN` secret must be added to repo Actions secrets
  before the first `v*` tag push will publish to the tap
