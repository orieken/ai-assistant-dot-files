# Framework Linter Configs (`shared/configs/`)

Drop-in fitness-function floors for each language the framework supports. Every file enforces
the single complexity cap prescribed in the corresponding convention file — nothing more.
Teams copy (not symlink) these files and may extend them freely after installation.

## Config → Convention → Cap

| Config file | Language | Tool | Cap | Source convention |
|---|---|---|---|---|
| `eslint.framework.config.js` | TypeScript | ESLint v9+ (flat config) | complexity ≤ 6 | `shared/rules/typescript-conventions.md` |
| `.golangci.yml` | Go | golangci-lint (gocyclo + revive) | complexity ≤ 6 | `shared/rules/go-conventions.md` |
| `ruff.toml` | Python | ruff (McCabe C901) | complexity ≤ 6 | `shared/rules/python-conventions.md` |
| `detekt.yml` | Kotlin | detekt (ComplexMethod) | complexity ≤ 6 | `shared/rules/kotlin-conventions.md` |
| `.swiftlint.yml` | Swift | SwiftLint (cyclomatic_complexity) | complexity ≤ 6 | `shared/rules/swift-conventions.md` |
| `clippy.toml` | Rust | Clippy (cognitive_complexity) | complexity ≤ 6 | `shared/rules/rust-conventions.md` |
| `checkstyle.xml` | Java | Checkstyle (CyclomaticComplexity) | complexity ≤ 6 | `shared/rules/java-conventions.md` |

**C# is not included**: `shared/rules/csharp-conventions.md` names no complexity enforcement
tool (unlike all other language conventions). When a C# complexity tool is wired in, add a
config here and update the table above.

## Why "cap ≤ 6" everywhere

The framework-wide rule is "cyclomatic complexity < 7" (see `shared/rules/design-principles.md`
and the Non-Negotiable Rules in `CLAUDE.md`). Every tool here flags functions that exceed 6 —
"< 7" and "max 6" are the same rule expressed two ways.

## Installation

These files are copied (not symlinked) by `install.sh --with-configs`, so teams own and can
extend them after install:

```bash
./install.sh --project /path/to/my-project --with-configs
```

Or copy manually:

```bash
cp shared/configs/.golangci.yml /path/to/my-go-project/
cp shared/configs/ruff.toml /path/to/my-python-project/
# etc.
```

## Cap-drift check

`scripts/check-cap-drift.sh` (wired into `scripts/health-check.sh`) verifies that every config
file's cap matches its source convention file. Run it standalone:

```bash
bash scripts/check-cap-drift.sh
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
