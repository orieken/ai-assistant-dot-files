---
name: dependency-update
description: Safe, structured process for updating packages in a pnpm monorepo.
triggers:
  keywords: ["update", "packages", "outdated", "audit", "dependabot"]
  intentPatterns: ["Update dependencies", "Update *", "Check for outdated packages", "Security advisories in dependencies", "Bump * to *"]
standalone: true
---

## When To Use
For minor and patch version updates. Do NOT use for major version updates with breaking changes, which should go through the full feature delivery pipeline.

## Context To Load First
1. `ARCHITECTURE_RULES.md`
2. `pnpm-workspace.yaml`

## Process
1. Run security audit first: `pnpm audit` — fix any Critical or High CVEs before anything else
2. Check for outdated packages: `pnpm outdated`
3. For each package to update, check the changelog before updating:
   - Find the package's CHANGELOG.md or GitHub releases
   - Read entries between current version and target version
   - Look for: breaking changes hidden in minor/patch, deprecation notices, peer dependency requirement changes, behavior changes in edge cases
4. Update packages one at a time (never batch major updates):
   `pnpm update [package]@[version] --filter [workspace]`
5. After each update: run `tsc --noEmit`, `pnpm build`, `pnpm -r run test`
6. If tests fail: isolate whether it's the update or a pre-existing issue
7. Regenerate lock file: `pnpm install`
8. Produce Dependency Update Report

## Output Format
`.claude/feature-workspace/dependency-update-report.md`

```markdown
# Dependency Update Report

Date: [today]

## Security Audit
- `pnpm audit` result: [N critical, N high, N medium]
- CVEs fixed: [list with CVE numbers]

## Updates Applied
| Package | From | To | Type | Changelog reviewed | Test result |
|---|---|---|---|---|---|

## Changelog Notes
- `[package] [version]`: [what changed that's relevant to this codebase]

## Peer Dependency Conflicts
- [Conflict]: [Resolution applied] / "None"

## Lock File
- Regenerated: ✅
- Committed: yes / pending

## Build & Test Results
- `tsc --noEmit`: ✅ / ❌ [errors]
- `pnpm build`: ✅ / ❌ [errors]
- Test suite: ✅ / ❌ [N failures]

## Skipped Updates
| Package | Current | Latest | Reason skipped |
|---|---|---|---|

## Recommended Next Steps
- [Packages that need major updates — with ADR suggestion]
- [Packages with known issues to watch]
```

## Guardrails
- Always run `pnpm audit` first — do not skip security check
- Never batch major version updates — each one needs individual analysis
- Always read the changelog before updating a package, however briefly
- Never commit a broken lock file
- Never use `--force` to resolve peer dependency conflicts silently

## Standalone Mode
Works without external systems. All commands are local pnpm operations.
