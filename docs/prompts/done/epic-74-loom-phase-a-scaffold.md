# Epic 74 Phase A — loom CLI: Go Module Scaffold

Parent epic: `docs/prompts/epic-74-loom-cli.md`
Run this phase first. Phase B (install logic) must not start until this phase is reviewed and approved.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Background

`loom` is a Go CLI that embeds the AI framework's `shared/` content at build time and installs
it to AI platform locations (Claude Code, Cursor, Windsurf, etc.) without requiring users to
clone the repo. It will be distributed via `brew install orieken/tap/loom`.

Phase A is scaffold only — directory structure, embed wiring, Cobra CLI with stub subcommands.
No file-writing logic yet.

## Current repo state

The repo already has a Go module at `shared/mcp/` — do not merge into it or modify it.
`shared/` contains 39 agents, 69 skills, 14 rules, plus configs, contracts, and schemas.
`install.sh` at the root handles current installation — do not modify it.

## What to build

### Module placement decision

The `//go:embed` directive can only reference paths relative to the `.go` file containing it,
and only within the same module. Since `shared/` is at the repo root, the module must also
be at the repo root OR the embed file must live adjacent to `shared/`.

**Preferred approach**: create a dedicated `go.mod` at the repo root for `loom` only:
- Module: `github.com/orieken/loom`
- Keep it separate from `shared/mcp/go.mod` (that module stays untouched)
- `cmd/loom/main.go` is the entrypoint

If this creates a conflict with `shared/mcp/go.mod` (e.g. Go tooling confusion with two
modules), fall back to placing `go.mod` inside `cmd/loom/` and the embed file at
`cmd/loom/embed.go` using a `../../shared` relative path — but test that `go build` resolves
the embed correctly before committing. Document the choice with a comment in `embed.go`.

### Directory structure

```
cmd/loom/
├── main.go
└── cmd/
    ├── root.go
    ├── install.go
    ├── uninstall.go
    ├── update.go
    ├── health.go
    └── version.go
embed.go          # at repo root alongside shared/ (or cmd/loom/embed.go — see above)
go.mod            # github.com/orieken/loom
go.sum
```

### embed.go

```go
package main

import "embed"

// frameworkFS holds all shared framework content baked in at compile time.
//go:embed all:shared/agents all:shared/skills all:shared/rules all:shared/configs all:shared/contracts all:shared/schemas
var frameworkFS embed.FS
```

Use explicit subdirectory patterns rather than `all:shared` to avoid accidentally embedding
large or binary files. Adjust the list if additional subdirectories are needed.

### cmd/loom/main.go

```go
package main

import "github.com/orieken/loom/cmd/loom/cmd"

func main() {
    cmd.Execute()
}
```

### cmd/loom/cmd/root.go

Cobra root command:
- Binary name: `loom`
- Short description: `"AI assistant framework installer"`
- Long description: brief paragraph explaining loom installs agents, skills, and rules for
  Claude Code, Cursor, Windsurf, and other AI platforms
- Register all five subcommands: `install`, `uninstall`, `update`, `health`, `version`
- `--version` flag prints the version string (hardcode `"0.1.0-dev"` for Phase A)

### All subcommands (Phase A stubs)

Each prints `"not yet implemented"` and exits 0. No flags needed on stubs.

```go
// install.go example
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install the framework to detected AI platforms",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("not yet implemented")
        return nil
    },
}
```

### go.mod dependencies

Phase A needs only one external dependency:
- `github.com/spf13/cobra` — CLI framework

Run `go mod tidy` after adding it.

## Guardrails

- No `any` / `interface{}` types.
- `golangci-lint run ./...` must pass before committing. Install with `brew install golangci-lint`
  if not present. Fix all reported issues — do not suppress with `//nolint` without a comment.
- Do not modify any file outside `cmd/loom/`, `embed.go`, `go.mod`, `go.sum`.
- `shared/mcp/go.mod` must remain untouched and `go build ./shared/mcp/...` must still pass.

## Verify before committing

```bash
go build -o loom ./cmd/loom
./loom --help           # shows usage + all 5 subcommands listed
./loom install          # prints "not yet implemented"
./loom --version        # prints version string
golangci-lint run ./... # clean (run from same directory as go.mod)
rm loom                 # clean up the test binary
```

## Commit

Single commit after all checks pass:
```
feat(loom): scaffold Go CLI module with embed and cobra subcommands
```
Stage only: `cmd/loom/`, `embed.go`, `go.mod`, `go.sum`

## Escalation

Stop and report if:
- `//go:embed` cannot resolve `shared/` from any sensible module placement — present the
  options and wait for direction before committing.
- `golangci-lint` reports issues you cannot resolve without changing the design.
- `shared/mcp/go.mod` breaks after adding the root `go.mod`.

## Report

On completion:
- Where `go.mod` ended up and why (the embed path constraint decision)
- Full output of `./loom --help`
- Lint result (pass or issues found + fixed)
- Commit hash
- Any open questions before Phase B begins
