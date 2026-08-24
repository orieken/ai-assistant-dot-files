# Epic 74 — loom: Homebrew-Distributed Framework CLI

Source: design discussion 2026-08-16. Supersedes `install.sh`/`uninstall.sh` as the primary
distribution mechanism for the ai-assistant-dot-files framework.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Prior context

The framework currently installs via `install.sh` which symlinks or copies `shared/` content
(39 agents, 69 skills, 14 rules) into AI platform locations (9 platforms: claude-code, cursor,
windsurf, github-copilot, gemini, openai-codex, jetbrains, roo-code, cline). Users must clone
the repo first — there is no package manager install path.

Oscar already distributes `testsmith` via a personal Homebrew tap. The goal is to add `loom`
to the same tap so users can:

```bash
brew install orieken/tap/loom   # macOS / Linux
winget install orieken.loom     # Windows
loom install                    # writes framework to all detected AI platforms
```

No repo clone required. The binary embeds all `shared/` content at build time.

## Design decisions

| Decision | Choice | Rationale |
|---|---|---|
| Language | Go | Repo already has `shared/mcp/` Go module; single-binary distribution |
| Embed strategy | `//go:embed shared/` | All framework files baked in at compile time; no runtime download |
| Subcommands | `install`, `uninstall`, `update`, `version`, `health` | Mirrors current `install.sh` capabilities |
| Release tooling | [goreleaser](https://goreleaser.com/) | Standard Go release tool; generates Homebrew formula automatically |
| Targets | macOS arm64/amd64, Linux arm64/amd64, Windows amd64 | Matches testsmith pattern |
| Homebrew | Oscar's existing tap repo | Add alongside testsmith, no new tap repo needed |
| Winget | Submit to `microsoft/winget-pkgs` | Optional stretch goal — Homebrew first |
| `install.sh` fate | Kept as fallback | Repo-cloners and CI still use it; deprecated in README, not deleted |

## Scope

### Phase A — Go module scaffold (one commit, then PAUSE for review)

Create `cmd/loom/` inside this repo (sibling to `shared/mcp/`):

```
cmd/loom/
├── main.go              # CLI entrypoint, subcommand dispatch
├── embed.go             # //go:embed directives
├── cmd/
│   ├── install.go       # loom install [flags]
│   ├── uninstall.go     # loom uninstall
│   ├── update.go        # loom update
│   ├── health.go        # loom health
│   └── version.go       # loom version
├── internal/
│   ├── platform/        # platform detection + install logic per platform
│   ├── embed/           # embedded FS access helpers
│   └── symlink/         # symlink vs copy logic (mirrors install.sh behavior)
go.mod                   # module: github.com/orieken/loom
go.sum
```

**embed.go** — embed all framework content:
```go
package loom

import "embed"

//go:embed shared/agents shared/skills shared/rules shared/configs shared/contracts shared/schemas
var FrameworkFS embed.FS
```

Note: `shared/mcp/` is NOT embedded — it is a separate Go module the user opts into with
`loom install --with-mcp`, which scaffolds the MCP server source into the target project.

**main.go** — Cobra CLI with subcommands. Use
[`cobra`](https://github.com/spf13/cobra) + [`viper`](https://github.com/spf13/viper) for
flags and config, consistent with the Go ecosystem standard.

Commit Phase A as: `feat(loom): scaffold Go CLI module with embed and cobra subcommands`

**PAUSE — present the module structure for approval before implementing subcommand logic.**

---

### Phase B — `loom install` (after Phase A approval)

Implement `loom install` to replicate `install.sh` behavior. Flags:

| Flag | Type | Default | Description |
|---|---|---|---|
| `--platform <name>` | string | (all detected) | Install for one platform only |
| `--stack <list>` | string | (all rules) | Comma-separated language stacks (go, typescript, python, java, kotlin, swift, rust, csharp) |
| `--copy` | bool | false | Use file copies instead of symlinks (required on Windows/WSL) |
| `--with-configs` | bool | false | Also write linter configs from `shared/configs/` to target |
| `--with-mcp` | bool | false | Scaffold MCP server source into `<target>/<project>-mcp/` |
| `--dry-run` | bool | false | Print actions without writing files |
| `--target <path>` | string | `./` | Project root to install into |

Platform detection: check for known config dirs/files (`.claude/`, `.cursor/`, `.gemini/`,
`.github/copilot-instructions.md`, etc.) in the target path. Install to all detected
platforms by default; `--platform` restricts to one.

Install logic per platform must mirror `install.sh`'s existing behavior exactly — read
`install.sh`'s `install_claude_code()`, `install_cursor()`, `install_windsurf()`, etc.
functions and port them to Go. Do not invent new behavior in Phase B.

Commit as: `feat(loom): implement install subcommand with platform detection`

---

### Phase C — `loom uninstall`, `loom version`, `loom health`

**`loom uninstall`**: reverse of install — remove symlinks or copied files written by loom.
Accept `--platform` to target one platform. Never delete files loom did not write (detect via
a `.loom-manifest.json` written to the target during install that records every path touched).

**`loom version`**: print the embedded framework version (read from a `shared/VERSION` file
or the binary's build-time ldflags) and the loom binary version.

**`loom health`**: run a subset of `scripts/health-check.sh` checks that are safe to port to
Go — primarily: verify symlinks are unbroken, verify agent/skill counts match embedded counts,
check for stale `framework-install.json` markers. Delegate complex shell checks (golangci-lint,
parity scripts) to a message directing the user to the repo's `scripts/health-check.sh` for
full verification.

Commit each as its own commit: `feat(loom): implement uninstall/version/health subcommands`

---

### Phase D — goreleaser + GitHub Actions release pipeline

Create `.goreleaser.yaml` at the repo root:

```yaml
version: 2
project_name: loom
builds:
  - id: loom
    main: ./cmd/loom
    binary: loom
    env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
brews:
  - name: loom
    repository:
      owner: orieken
      name: homebrew-tap       # adjust to match the actual tap repo name
    homepage: https://github.com/orieken/loom
    description: "AI assistant framework installer — agents, skills, and rules for Claude, Cursor, and more"
    license: MIT
archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
checksum:
  name_template: "checksums.txt"
release:
  github:
    owner: orieken
    name: ai-assistant-dot-files
```

Create `.github/workflows/release.yml`:

```yaml
name: Release loom
on:
  push:
    tags: ['v*']
permissions:
  contents: write
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@<pin-to-sha>
        with: {fetch-depth: 0}
      - uses: actions/setup-go@<pin-to-sha>
        with: {go-version-file: cmd/loom/go.mod}
      - uses: goreleaser/goreleaser-action@<pin-to-sha>
        with: {version: latest, args: release --clean}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

**HOMEBREW_TAP_GITHUB_TOKEN** must be a PAT with write access to the tap repo. Add it to the
repo's secrets — this requires human action; flag it in the report.

Pin all action versions to full commit SHAs (see `shared/rules/iac-conventions.md` GitHub
Actions guardrails). Do not use mutable tags like `@v4`.

Commit as: `chore(loom): add goreleaser config and GitHub Actions release pipeline`

---

### Phase E — README and documentation

Update `README.md`:
- Add an "Installation" section at the top: `brew install orieken/tap/loom` / `loom install`
- Mark `install.sh` as "for repo contributors and CI" rather than the primary path
- Add a "Quick start" showing the three-command flow: install loom → `loom install` → done

Update `docs/ONBOARDING.md`:
- Replace the "clone and run install.sh" steps with the loom path as primary
- Keep the repo-clone path as "contributing / advanced"

Add `cmd/loom/README.md` with: subcommand reference, flag table, platform detection logic,
how embedded content is versioned.

Commit as: `docs(loom): update README and ONBOARDING for loom install path`

---

## Guardrails

- One commit per Phase — do not combine Phase A scaffold with Phase B logic.
- `install.sh` is NOT modified or deleted in this epic. Loom is additive.
- No `any` types in Go (`shared/rules/architecture-guardrails.md` #4).
- All network calls (if any) must have explicit timeouts.
- Action SHAs in GitHub Actions workflows must be pinned — never `@v4` or `@main`.
- `golangci-lint run ./...` must pass in `cmd/loom/` before each commit
  (same standard as `shared/mcp/`).
- The `HOMEBREW_TAP_GITHUB_TOKEN` secret is a human action — pause and report at Phase D
  if it has not been set.

## Escalation

Stop and report if:
- The tap repo name is not `homebrew-tap` — get the correct name before writing goreleaser config.
- `shared/` embed size exceeds 50MB (run `du -sh shared/` to check before Phase A).
- Phase A structure review is rejected — do not implement Phase B without approval.

## Report format

After each Phase, report:
- Files created/modified
- `golangci-lint` result (pass / any warnings)
- Open questions or human actions required (secrets, tap repo name, etc.)
- Commit hash
