# Epic 74 Phase C — loom CLI: `uninstall`, `version`, `health` Subcommands

Parent epic: `docs/prompts/epic-74-loom-cli.md`
Prerequisite: Phase B (`epic-74-loom-phase-b-install.md`) must be complete.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Expected state coming in

Phase B shipped `loom install` with platform detection, manifest writing, symlink/copy logic,
stack filtering, and `--dry-run`. `.loom-manifest.json` is written to the target on install.
The three remaining subcommands (`uninstall`, `version`, `health`) still print
"not yet implemented".

## What to build

### Op 1 — `loom uninstall`

Remove everything loom wrote during `loom install`. Read the `.loom-manifest.json` in the
target directory to determine what to remove — never delete files that loom did not write.

**Flags:**

| Flag | Type | Default | Description |
|---|---|---|---|
| `--target <path>` | string | `"."` | Project root to uninstall from |
| `--platform <name>` | string | `""` (all) | Remove only one platform's files |
| `--dry-run` | bool | `false` | Print what would be removed without deleting |

**Behavior:**

1. Read `.loom-manifest.json` from `--target`. If it does not exist, print an error:
   `"loom: no manifest found at <target>/.loom-manifest.json — was loom install run here?"`
   and exit non-zero.
2. For each path in `manifest.paths`, remove the symlink or file/directory. Skip any path
   that no longer exists (already manually removed) with a warning.
3. After removal, delete `.loom-manifest.json` itself (unless `--dry-run`).
4. If `--platform` is set, only remove paths that belong to that platform (use the platform
   prefix in the path to determine ownership, e.g. `.claude/` → `claude-code`).

**Output:**

```
loom: uninstalling framework from .

  ✓ removed .claude/agents
  ✓ removed .claude/skills
  ✓ removed .cursor/agents

3 paths removed. Manifest deleted.
```

Commit Op 1 as: `feat(loom): implement uninstall subcommand using manifest`

---

### Op 2 — `loom version`

Print the embedded framework version and the loom binary version.

**Implementation:**

Use Go build-time ldflags to inject the version at release time. Define package-level vars
in `cmd/loom/cmd/version.go`:

```go
var (
    version = "dev"     // set by -ldflags "-X main.version=..."
    commit  = "none"    // set by -ldflags "-X main.commit=..."
    date    = "unknown" // set by -ldflags "-X main.date=..."
)
```

The goreleaser config (Phase D) will inject these. For local dev builds, they remain `"dev"`.

Read the embedded framework version from `shared/VERSION` if it exists; otherwise report
`"embedded"`. If `shared/VERSION` does not yet exist, create it at the repo root with content
`v3.3.14` (the current framework version — check `git tag --sort=-v:refname | head -1` to
confirm the latest tag).

**Output:**

```
loom v0.1.0 (commit abc1234, built 2026-08-16)
framework v3.3.14 (embedded)
```

Commit Op 2 as: `feat(loom): implement version subcommand with build-time ldflags`

---

### Op 3 — `loom health`

Run a lightweight health check against a target directory to verify the loom install is intact.

**Flags:**

| Flag | Type | Default | Description |
|---|---|---|---|
| `--target <path>` | string | `"."` | Project root to check |

**Checks to implement (in order):**

1. **Manifest exists** — `.loom-manifest.json` is present in target. FAIL if missing.
2. **Manifest version matches binary** — the `version` field in the manifest matches the
   embedded framework version. WARN if they differ (upgrade available).
3. **Paths intact** — every path listed in `manifest.paths` exists. FAIL for each missing path.
4. **Symlinks unbroken** — for every symlinked path, the target of the symlink resolves.
   FAIL for broken symlinks.
5. **Agent count** — count `.md` files in the installed agents directory (excluding CHANGELOG).
   WARN if the count differs from the embedded count.

For checks that cannot be safely ported from `scripts/health-check.sh` (golangci-lint, parity
scripts, MCP drift), print a notice:

```
note: for full framework health checks, run scripts/health-check.sh in the framework repo
```

**Output format:**

```
loom health: checking .

  ✓ manifest found (v3.3.14)
  ✓ all 3 installed paths intact
  ✓ symlinks unbroken
  ⚠ agent count mismatch: manifest has 39, found 38 (one may have been manually removed)

1 warning. Run `loom install` to repair.
```

Exit 0 if no FAILs; exit 1 if any FAIL.

Commit Op 3 as: `feat(loom): implement health subcommand with manifest and path checks`

---

## Guardrails

- Each Op is a separate commit — do not combine.
- No `any` / `interface{}` types.
- All errors must be returned and handled — no silent swallows.
- `golangci-lint run ./...` must pass before each commit.

## Verify before each commit

```bash
go build -o loom ./cmd/loom

# Op 1 verify
./loom install --target /tmp/test-target --copy
./loom uninstall --target /tmp/test-target --dry-run  # shows what would be removed
./loom uninstall --target /tmp/test-target             # removes files
ls /tmp/test-target/.loom-manifest.json 2>/dev/null || echo "manifest removed ✓"

# Op 2 verify
./loom version   # prints version strings

# Op 3 verify
./loom install --target /tmp/test-target --copy
./loom health --target /tmp/test-target    # all checks pass
rm /tmp/test-target/.claude/agents -rf
./loom health --target /tmp/test-target    # shows FAIL for missing path

golangci-lint run ./...
rm loom && rm -rf /tmp/test-target
```

## Escalation

Stop and report if:
- The manifest path structure from Phase B makes it ambiguous which paths belong to which
  platform — fix the manifest schema before implementing `--platform` uninstall.
- `shared/VERSION` does not exist and the latest git tag is unclear.

## Report

After all three ops:
- Sample `loom version` output
- Sample `loom health` output (passing and failing case)
- Sample `loom uninstall --dry-run` output
- Lint result
- Commit hash for each op
