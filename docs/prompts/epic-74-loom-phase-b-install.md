# Epic 74 Phase B — loom CLI: `loom install` Subcommand

Parent epic: `docs/prompts/epic-74-loom-cli.md`
Prerequisite: Phase A (`epic-74-loom-phase-a-scaffold.md`) must be complete and approved.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Expected state coming in

Phase A shipped `cmd/loom/` with a working Cobra CLI and stub subcommands. The `frameworkFS`
embed var is wired and `go build` produces a working binary. `loom install` currently prints
"not yet implemented".

## What to build

Implement `loom install` to replicate `install.sh` behavior. Read `install.sh` in full before
writing any Go code — every install behavior must be ported faithfully, not reimagined.

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--platform <name>` | string | `""` (all detected) | Install for one platform only. Valid values: `claude-code`, `cursor`, `windsurf`, `github-copilot`, `gemini`, `openai-codex`, `jetbrains`, `roo-code`, `cline` |
| `--stack <list>` | string | `""` (all rules) | Comma-separated language stacks: `go`, `typescript`, `python`, `java`, `kotlin`, `swift`, `rust`, `csharp` |
| `--copy` | bool | `false` | Use file copies instead of symlinks (required on Windows/WSL) |
| `--with-configs` | bool | `false` | Also write `shared/configs/` linter configs to target |
| `--with-mcp` | bool | `false` | Scaffold `shared/mcp/` reference source into `<target>/<project>-mcp/` |
| `--dry-run` | bool | `false` | Print planned actions without writing any files |
| `--target <path>` | string | `"."` | Project root to install into |

### Internal package structure

Add these packages under `cmd/loom/internal/`:

```
cmd/loom/internal/
├── platform/
│   ├── detect.go       # detect which AI platforms are present in the target dir
│   ├── claude.go       # install logic for claude-code
│   ├── cursor.go       # install logic for cursor
│   ├── windsurf.go     # install logic for windsurf
│   ├── copilot.go      # install logic for github-copilot
│   ├── gemini.go       # install logic for gemini
│   ├── codex.go        # install logic for openai-codex
│   ├── jetbrains.go    # install logic for jetbrains
│   ├── roocode.go      # install logic for roo-code
│   └── cline.go        # install logic for cline
├── fs/
│   ├── link.go         # symlink or copy a file/dir (respects --copy flag)
│   └── embed.go        # helpers for reading from frameworkFS
└── manifest/
    └── manifest.go     # write/read .loom-manifest.json tracking installed paths
```

### Platform detection (`detect.go`)

Detect platforms by checking for known markers in the target directory:

| Platform | Detection marker |
|---|---|
| `claude-code` | `.claude/` directory exists OR `CLAUDE.md` exists |
| `cursor` | `.cursor/` directory exists |
| `windsurf` | `.windsurf/` directory exists |
| `github-copilot` | `.github/copilot-instructions.md` exists OR `.github/` exists |
| `gemini` | `.gemini/` directory exists |
| `openai-codex` | `.codex/` directory exists OR `AGENTS.md` exists at root |
| `jetbrains` | `.aiassistant/` directory exists OR `.junie/` exists |
| `roo-code` | `.roo/` directory exists OR `.roomodes` file exists |
| `cline` | `.cline/` directory exists OR `.clinerules/` exists |

Return all detected platforms. If `--platform` is set, validate it is a known name and skip
detection — use it directly.

### Install logic per platform

Read `install.sh` functions `install_claude_code()`, `install_cursor()`, etc. and port them
to Go. The Go implementation must produce identical file outcomes to the bash script. Key
behaviors to preserve:

- **Symlink vs copy**: when `--copy` is false, create symlinks pointing into `frameworkFS`
  extracted temp paths, OR extract to a stable local cache dir (e.g.
  `~/.loom/cache/<version>/`) and symlink from there. Do not symlink into the binary itself.
  When `--copy` is true, write files directly.
- **Stack filtering**: when `--stack` is set, only install rule files whose names match the
  requested stacks (e.g. `go-conventions.md`, `typescript-conventions.md`). Always include
  core rules (`architecture-guardrails.md`, `design-principles.md`, `approval-gates.md`,
  `testing-conventions.md`).
- **Dry run**: log every action with a `[dry-run]` prefix; write nothing to disk.

### Manifest (`.loom-manifest.json`)

After a successful install, write `.loom-manifest.json` to the target root:

```json
{
  "version": "0.1.0",
  "installedAt": "2026-08-16T12:00:00Z",
  "platforms": ["claude-code", "cursor"],
  "paths": [
    ".claude/agents",
    ".claude/skills",
    ".cursor/agents"
  ]
}
```

This file is used by `loom uninstall` (Phase C) to know exactly what to remove.

### Output

Print a summary on completion:

```
loom: installing framework v0.1.0

  ✓ claude-code  → .claude/agents, .claude/skills, .claude/rules
  ✓ cursor       → .cursor/agents, .cursor/skills

2 platforms, 39 agents, 69 skills, 14 rules installed.
Manifest written to .loom-manifest.json
```

On dry run, prefix each line with `[dry-run]` and omit the manifest line.

## Guardrails

- No `any` / `interface{}` types.
- Every filesystem operation must check and return errors — no silent swallows.
- All file paths constructed from user input (`--target`) must be cleaned with `filepath.Clean`
  and validated to prevent path traversal.
- `golangci-lint run ./...` must pass before committing.
- Do not modify `install.sh` — loom is additive in this phase.

## Verify before committing

```bash
# Build
go build -o loom ./cmd/loom

# Dry run against a temp dir
mkdir /tmp/test-target && ./loom install --target /tmp/test-target --dry-run
# Should print planned actions, write nothing

# Real install into temp dir (copy mode to avoid symlink permission issues in CI)
./loom install --target /tmp/test-target --copy
# Should write files; check .loom-manifest.json exists

# Platform filter
./loom install --target /tmp/test-target --platform claude-code --dry-run

# Stack filter
./loom install --target /tmp/test-target --stack go --dry-run

# Lint
golangci-lint run ./...

# Cleanup
rm loom && rm -rf /tmp/test-target
```

## Commit

Single commit after all checks pass:
```
feat(loom): implement install subcommand with platform detection and manifest
```
Stage only files under `cmd/loom/`.

## Escalation

Stop and report if:
- `install.sh` contains platform logic that cannot be cleanly ported to Go without significant
  redesign — describe the gap and wait for direction.
- The symlink strategy for embedded files needs a design decision (temp extract vs cache dir).
- `golangci-lint` reports issues that require a design change.

## Report

On completion:
- Sample `--dry-run` output against a fresh temp directory
- Content of `.loom-manifest.json` after a real `--copy` install
- Lint result
- Commit hash
- Any platform logic that diverges from `install.sh` and why
