# Epic 64 — Ship the Linter Configs the Conventions Prescribe (`shared/configs/`)

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #4). The gap: every
convention file names an enforcement tool and a cap, but the framework ships zero config files —
installed projects hand-author their fitness functions from prose, defeating
`shared/rules/architecture-guardrails.md` #7.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

Caps prescribed across `shared/rules/*-conventions.md` (READ EACH — this table is from the audit,
verify before writing configs):

| Language | Tool | Cap |
|---|---|---|
| TypeScript | ESLint `complexity` | 6 |
| Go | golangci-lint (gocyclo/revive) | complexity < 7 |
| Python | ruff (+ radon/flake8-cognitive-complexity per docs copy) | < 7 |
| Kotlin | detekt `ComplexMethod` (+ ktlint format) | 6 |
| Swift | SwiftLint `cyclomatic_complexity` | 6 |
| Rust | clippy `cognitive_complexity` | 6 |
| Java | Checkstyle | < 7 |
| C# | (check `csharp-conventions.md` — no cap tool named; skip if absent) | — |

Also prescribed but unshipped: `strict: true` tsconfig posture, `#![forbid(unsafe_code)]`,
`cargo-deny` allow-lists. `install.sh` currently installs rules/agents/skills but no tool
configs. `scripts/check-inventory-drift.sh` greps prose for counts — the same pattern can
cross-check caps.

## Scope: 3 ops

**Op 1 — `shared/configs/` with one config per prescribed tool.**
Minimal, drop-in files: `.eslintrc.framework.json` (or flat-config equivalent — pick the current
ESLint default), `.golangci.yml`, `ruff.toml`, `detekt.yml`, `.swiftlint.yml`,
`clippy.toml`-or-documented-lint-attrs (clippy's cap lives in code attrs — ship a
`rust-lints.md` snippet instead if a config file can't express it; be honest), `checkstyle.xml`.
Each file: header comment naming its source convention file and the cap it enforces. A
`shared/configs/README.md` maps config → convention → cap. Do NOT invent rules beyond what the
conventions prescribe — these are fitness-function floors, not full lint policies teams must
adopt wholesale.
Commit: `feat(configs): ship linter fitness-function configs (Epic 64 Op 1)`

**Op 2 — `install.sh --with-configs` (opt-in flag).**
Copies (never symlinks — teams will edit them) the relevant configs into a target project,
skipping any that already exist (idempotency discipline matches the existing `--copy` mode
behavior). Update install verification (`scripts/test-install.sh`) with one case.
Commit: `feat(install): --with-configs flag (Epic 64 Op 2)`

**Op 3 — Cap-drift cross-check.**
Extend `check-inventory-drift.sh` (or a sibling script wired into health-check the same way):
parse each shipped config's cap value and each convention file's stated cap; FAIL on mismatch.
Deterministic parsing only — if a config format can't be parsed without new dependencies, WARN
and document (the `6c422cb` lesson: never false-FAIL on missing tooling).
Commit: `feat(health-check): cap-drift check between configs and conventions (Epic 64 Op 3)`

After every op: `bash scripts/health-check.sh` + `bash scripts/check-inventory-drift.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- A convention names a cap ambiguously (e.g., "< 7" vs "capped at 6" for the same tool) — halt
  only if genuinely contradictory; "< 7" and "max 6" are the same rule, pick 6 and note it.
- A tool's current config format has diverged from what the convention doc describes (e.g.,
  ESLint flat config vs `.eslintrc`) — follow the tool's CURRENT reality, note the convention doc
  needs a touch-up, and make that touch-up in Op 1's commit.
- Shipping a config would exceed "fitness-function floor" into opinionated-full-policy territory
  — trim to the floor and note what was excluded.

## Report (under 120 words)

```
Commits: <sha> x3
Configs shipped: <list>
Caps verified against conventions: <n>/<n> matched (mismatches fixed: <list | none>)
install.sh --with-configs: <tested how>
Cap-drift check: <pass; parsing method per format>
```

Go.
