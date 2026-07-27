# Update the Framework in an Installed Project

Propagate framework updates from `ai-assistant-dot-files` into a target project, detecting which install pattern the project uses and running the correct update sequence for that pattern.

## When to use

- Framework has changed (new agents, updated skills, new templates, new tools, etc.)
- Target project needs to pick up those changes
- Unsure which install pattern the project uses (traditional install / bridge / adopted saturday-mcp / mixed)

## Prerequisites

- Absolute path to the target project
- Absolute path to the framework clone (typically `/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`)
- Framework clone at the version you want to propagate (pull latest first if that's the intent)
- **Both** target project and framework at clean working trees (`git status --short` shows only intended WIP)

If either working tree is dirty in a way you don't want carried into the update, halt before starting.

## The four install patterns

The framework can live in a target project in one of four shapes. Phase A of this prompt detects which one(s) apply — a single project can have more than one.

| Pattern | How to detect | Update mechanism |
|---|---|---|
| **A. Traditional install** — symlinks | `.claude/agents/`, `.claude/skills/`, `.claude/rules/` in target project exist and are symlinks pointing into the framework clone | `install.sh` re-run (refreshes symlinks; usually no-op since symlinks track the framework live) |
| **B. Bridge install** — copied tool source in project's MCP | Target project has an MCP server (e.g., `./project/mcp/`) with source files whose content matches or descends from `shared/mcp-patterns/go/tools/*` | Manual re-copy of updated tool files + regen tests + rebuild MCP server (source was **copied**, not symlinked, so upstream changes don't flow automatically). Use `shared/mcp-patterns/` in the framework install as the diff source. |
| **C. Adopted saturday-mcp** — used as-is | Target project uses `saturday-mcp` (as a submodule, dependency, or installed binary) rather than its own MCP with bridged tools | Update the saturday-mcp version pin (submodule commit / dependency version / installed-binary rebuild); restart MCP server |
| **D. Generated platform configs** — non-symlinkable | `.cursor/rules/*.mdc`, `.cursorrules`, `.windsurfrules`, `.github/copilot-instructions.md`, `.openai.md`, `AGENTS.md` exist in target project as regenerated content | Re-run `scripts/generate-configs.sh` (in-repo tool that produces the platform-specific files from the shared sources) |

Most projects hit A + D at minimum. Bridge projects hit A + B + D. saturday-mcp-adopting projects hit A + C + D.

## Scope

### Phase A — Detect install patterns (one commit — investigation log)

Run these checks against the target project and record findings in `./project/docs/framework-update-<YYYY-MM-DD>.md`:

1. **Pattern A**: is `./project/.claude/agents/` present? Are its entries symlinks? `readlink` each to verify they resolve into a valid framework clone path
2. **Pattern B**: does `./project/mcp/` (or equivalent) exist? Do its `internal/tools/*.go` (or TS/Python equivalents) match filenames + shape from `saturday-mcp/internal/tools/`? A quick shape-check: same tool names? same input schemas? same response struct shapes?
3. **Pattern C**: is `saturday-mcp` referenced as a dependency (in `go.mod`, `package.json`, `pyproject.toml`) OR installed as a system binary reachable on `PATH`? If yes, note the pinned version
4. **Pattern D**: which platform config files exist under `.cursor/`, `.github/`, etc.? Note timestamps to see how stale they are relative to `shared/` source-of-truth
5. **What's changed upstream**: `cd /path/to/ai-assistant-dot-files && git log --oneline HEAD@{last-update-of-target}..HEAD -- shared/` — enumerate the changed shared/ files since target project last synced (approximate — if target didn't record last-sync, use last framework commit that's in target's `.claude/`)
6. **Local modifications in target**: any files under `./project/.claude/` that are NOT symlinks? Any tool files in `./project/mcp/` that diverge from their saturday-mcp counterparts beyond expected path adjustments?

Write findings to `./project/docs/framework-update-<YYYY-MM-DD>.md`. Include:
- Detected patterns (A/B/C/D — usually multiple)
- Changed upstream files (categorized: agents, skills, rules, templates, contracts, schemas, tools, blueprints)
- Local modifications in target (with paths)
- Estimated update ops

Commit (in target repo): `docs(framework): investigate update needs <YYYY-MM-DD>`.

**Pause. Get user approval on the plan before Phase B.**

### Phase B — Execute the update sequence

Run the right steps for each detected pattern. Each is a separate op / commit.

#### If Pattern A detected — refresh symlinks
```bash
cd /path/to/ai-assistant-dot-files && ./install.sh --target /path/to/project
```
Verify `readlink -f` on a few symlinks now points to the updated framework paths.
Commit (in target repo, if install creates any tracked files): `chore(framework): re-run install.sh for latest framework updates`.

#### If Pattern D detected — regenerate platform configs
```bash
cd /path/to/ai-assistant-dot-files && bash scripts/generate-configs.sh
# then copy the regenerated files into the target project OR re-run install.sh if it handles this
```
Verify `bash scripts/check-parity.sh` passes (no drift between shared sources and generated configs).
Commit (in target repo): `chore(platform-configs): regenerate for framework updates`.

#### If Pattern B detected — sync copied tool source (most work)

**Diff source selection**: read from `<framework_clone>/shared/mcp-patterns/go/tools/` (or `analyzers/`, `server/`).

For each tool file that changed upstream:

1. Diff the target's copy against the current reference source:
   ```bash
   diff /path/to/project/mcp/internal/tools/<tool>.go \
        /path/to/ai-assistant-dot-files/shared/mcp-patterns/go/tools/<tool>.go
   ```
2. Identify what changed upstream + what was intentionally adapted downstream (path defaults, monorepo `packagePath` handling, etc.)
3. Re-apply the upstream changes to the target's copy PRESERVING the downstream adaptations
4. Update the corresponding test file the same way
5. Re-run `go build ./... && go test ./...` (or language equivalent) — must be green

Commit per tool (in target repo): `feat(mcp): sync <tool> with saturday-mcp <sha> (bridge update)`.

For the retriever adapters, analyzers, and shared helpers:
- Same diff-and-merge approach
- These change less frequently than tools; verify carefully because they underpin every bridged tool

For new tools that landed in `saturday-mcp` since the last update:
- Ask user first — should this project bridge the new tools too? A bridge project may deliberately not want every new saturday-mcp tool
- If yes, follow the per-tool ops from `install-framework-with-mcp-bridge.md`'s Phase C

#### If Pattern C detected — update saturday-mcp version
- Submodule: `git submodule update --remote` in target project
- Dependency: bump version in `go.mod` / `package.json`
- Installed binary: rebuild + reinstall

Restart the MCP server. Verify tool discovery lists the expected tools (count matches saturday-mcp's current tool count).
Commit (in target repo): `chore(mcp): update saturday-mcp to <version> for framework updates`.

### Phase C — Contract-drift check (one commit — verification report)

The framework may have changed CONTRACTS between framework versions (e.g., a new required section added to `shared/contracts/analysis-contract.md`). Projects that have DELIVERED features under an older contract won't retroactively meet the new one.

1. Enumerate contract changes: `cd /path/to/ai-assistant-dot-files && git log --oneline HEAD@{last-sync}..HEAD -- shared/contracts/`
2. For each changed contract, run `/validate-artifact` (via skill or manual grep) against representative existing artifacts in `./project/docs/features/*/*.md` (if any)
3. Report any artifacts that would newly fail under the new contract — these are NOT regressions in the target project, they're a signal the project should either (a) update those old artifacts to the new shape or (b) explicitly note them as "authored under contract vN" and accept the drift

Write findings to `./project/docs/framework-update-<YYYY-MM-DD>.md` (append a Phase C section).

Commit (in target repo): `docs(framework): contract-drift check post-update`.

## Discipline (non-negotiable)

- One commit per op.
- Conventional Commits.
- **NEVER `git add -A`** in the target project — bridge updates in particular touch multiple files and it's easy to sweep in unrelated WIP.
- `git status --short` after staging, before every commit.
- Build + test green per commit for any pattern that touches code (Pattern B always; Pattern C when restarting).
- Do NOT push in the target project — human step.

## Escalation criteria

Stop and report if:
- Symlinks in the target's `.claude/` point to a DIFFERENT framework clone path than the one the user specified — halt, describe (target may have multiple framework installs, needs decision)
- Pattern B diffs reveal downstream modifications that CONFLICT with upstream changes (both changed the same lines) — halt for each conflict, describe both sides, ask
- A tool file in the target has been renamed or moved in a way that doesn't match the current saturday-mcp shape — halt, describe
- Any new tool file in saturday-mcp since last sync — halt, ask if it should be bridged (default: don't auto-bridge new tools; that's a Phase C decision from the bridge install prompt)
- Pattern C detected AND the pinned version bump would cross a major saturday-mcp version — halt, may need review of breaking changes
- Phase C contract-drift check reveals existing artifacts newly failing — halt with the list, DO NOT auto-modify persisted artifacts (they're historical records)

## Report format (per phase, under 300 words each)

### Phase A report
```
Target project: <absolute path>
Detected patterns: <A/B/C/D + brief description>
Upstream changes since last sync: <n> files in shared/, categorized:
  - agents: <n>
  - skills: <n>
  - rules: <n>
  - templates: <n>
  - contracts: <n>
  - schemas: <n>
  - blueprints: <n>
  - tools (in saturday-mcp/): <n>

Local modifications in target: <list, or "none">
Contract changes worth Phase C attention: <list>
Estimated Phase B ops: <n>

Investigation commit: <sha>
Recommended next: proceed to Phase B with these ops
```

### Phase B report
```
Ops executed (per pattern):
  Pattern A: <commit sha, or N/A>
  Pattern D: <commit sha, or N/A>
  Pattern B: <n commits — list SHAs>
  Pattern C: <commit sha, or N/A>

Build + test green after each: verified
Tool discovery lists expected count: <n>

Anything skipped or halted: <list with reasons>
```

### Phase C report
```
Contracts changed since last sync: <n>
Persisted artifacts checked: <n>
Newly-failing artifacts: <list, or "None">
Recommendation per newly-failing artifact:
  - <path>: update to new contract | accept drift + annotate | rewrite

Contract-drift check commit: <sha>
```

Go.
