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

## Phase B audit carry-forwards

These issues were found and partially fixed during the Phase B code review. Finish them before
starting Op 1.

### 1. Global `installArgs` coupling (pre-flight refactor commit required)

`installExtras`, `writeManifest`, and `reportInstall` in `cmd/loom/cmd/install_run.go` each
read the package-level `installArgs` var directly instead of using parameters. This makes them
untestable without Cobra flag parsing and will make Phase C subcommands harder to build
correctly if the pattern spreads.

Fix before Op 1 in a separate commit (`refactor(loom): decouple install helpers from flag globals`):

- Add `withConfig`, `withMCP`, and `isDryRun` fields to `installRequest` in
  `cmd/loom/cmd/install_options.go`. Populate them from `installArgs` inside `prepareInstall`.
- Update `installExtras` to read from `installRequest` instead of `installArgs`.
- Update `writeManifest` to take `isDryRun bool` as a parameter (or use `installRequest`).
- Update `reportInstall` to use `output.isDryRun` (already on the `installOutput` struct)
  instead of `installArgs.isDryRun`.
- No behavior change — this is a pure structural refactor. All existing tests must still pass.

### 2. Manifest schema missing per-platform path ownership

`.loom-manifest.json` records a flat `paths` array with no record of which platform owns each
path. The `--platform` filter on `loom uninstall` (Op 1 below) cannot work correctly from a
flat list — inferring ownership from path prefixes (`.claude/` → `claude-code`) is fragile and
will break if a platform writes to an unexpected directory.

Fix the manifest schema as part of Op 1 before implementing uninstall:

Update `cmd/loom/internal/manifest/manifest.go` to replace the flat `Paths []string` field
with a per-platform structure:

```go
type PlatformRecord struct {
    Name  string   `json:"name"`
    Paths []string `json:"paths"`
}

type Manifest struct {
    Version     string           `json:"version"`
    InstalledAt time.Time        `json:"installedAt"`
    Platforms   []PlatformRecord `json:"platforms"`
}
```

Update `writeManifest` in `install_run.go` to build `[]PlatformRecord` from `[]platform.Result`
instead of collecting a flat path list. The `uniquePaths` helper and `paths` field in
`installRequest` can be removed. Update `manifest_test.go` to match the new schema.

This is a breaking change to the manifest format — that is acceptable because no released
version of loom exists yet.

## What to build

### Op 1 — `loom uninstall`

**Prerequisite:** complete the pre-flight refactor commit AND the manifest schema fix (both
described in the carry-forwards section above) before writing any uninstall logic.

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
2. For each `PlatformRecord` in `manifest.platforms` (using the updated schema from the
   carry-forwards section), remove the paths listed for that platform. Skip paths that no
   longer exist with a warning.
3. After removal, delete `.loom-manifest.json` itself (unless `--dry-run`).
4. If `--platform` is set, filter to only the `PlatformRecord` whose `name` matches — no
   path-prefix guessing needed because the schema records ownership explicitly.

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

Use Go build-time ldflags to inject the version at release time. Replace the `const version`
in `cmd/loom/cmd/root.go` with package-level vars in a new file
`cmd/loom/cmd/buildinfo.go` so all subcommands and the root `--version` flag share one source
of truth:

```go
// cmd/loom/cmd/buildinfo.go
package cmd

// Build metadata injected by goreleaser via -ldflags.
// Values remain "dev"/"none"/"unknown" for local builds.
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)
```

The ldflags paths use the full package import path, not `main`:

```
-ldflags "-X github.com/orieken/loom/cmd/loom/cmd.version=v0.1.0 \
           -X github.com/orieken/loom/cmd/loom/cmd.commit=abc1234 \
           -X github.com/orieken/loom/cmd/loom/cmd.date=2026-08-19"
```

Remove `const version = "0.1.0-dev"` from `root.go` — it is replaced by the var above. Update
`rootCmd.Version` to use the var; it already does via the `version` identifier so no other call
site changes are needed.

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
3. **Paths intact** — every path listed across all `manifest.platforms[*].paths` exists.
   FAIL for each missing path.
4. **Symlinks unbroken** — for every symlinked path across all platform records, the symlink
   target resolves. FAIL for broken symlinks.
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

- The pre-flight refactor commit must land before Op 1. Do not start Op 1 until
  `installExtras`, `writeManifest`, and `reportInstall` no longer read `installArgs` directly.
- Each Op is a separate commit — do not combine Ops with each other or with the pre-flight.
- No `any` / `interface{}` types.
- All errors must be returned and handled — no silent swallows.
- Helper functions must receive their inputs as parameters — never read package-level flag
  variables (`installArgs`, `uninstallArgs`, etc.) from inside a helper. Only `RunE` functions
  may read flag vars; they pass values down as parameters or request structs.
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
- `shared/VERSION` does not exist and the latest git tag is unclear.
- The pre-flight refactor reveals that decoupling `installExtras`/`writeManifest`/`reportInstall`
  from `installArgs` requires changing the `platform.Result` type or `installRequest` struct in
  a way that breaks the platform installer interface — describe the conflict before proceeding.

## Report

After all three ops:
- Sample `loom version` output
- Sample `loom health` output (passing and failing case)
- Sample `loom uninstall --dry-run` output
- Lint result
- Commit hash for each op
