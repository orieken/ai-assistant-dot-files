# Epic 50 — Rust / Systems Programming Conventions

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 3.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

`shared/rules/` has language conventions for TypeScript, Go, Python, Java, C#. No Rust coverage; framework can't guide Rust systems programming work today.

## Scope

**One commit.** Create `shared/rules/rust-conventions.md` matching the shape of existing convention files (read `shared/rules/go-conventions.md` for the closest structural analog — Rust and Go share a lot of philosophical DNA).

Sections:

- **Complexity**: match framework's `< 7` cyclomatic budget via `clippy::cognitive_complexity` or `rustc`'s built-in max-function-length lints
- **File naming**: snake_case files, one public top-level type per file for large types, module-per-directory (`mod.rs` or `<name>.rs`)
- **Testing**: 
  - Unit: `#[cfg(test)]` mod tests inline
  - Property-based: `proptest` for invariant testing
  - Mocking: `mockall` for trait-based test doubles
  - Table-driven: `rstest` fixtures + parameterized tests
- **Dependency management**: `cargo`, pinned versions in `Cargo.lock`, `deny` for supply-chain lints, `cargo audit` for known CVEs
- **Error handling**: `?` operator + `Result<T, E>` + `thiserror` for library errors + `anyhow` for application errors. Never `unwrap()` in library code (only test/example)
- **Async**: `tokio` as default runtime; `async fn` traits via `async-trait` (until native support stabilizes); explicit `.await` for readability
- **Clean Architecture layer separation**: same as other languages — Domain (pure Rust types + trait definitions) / Use-Cases / Adapters (crate-per-external-boundary) / Frameworks (bin crates depend on lib crates, never inverse)
- **Unsafe**: `#![forbid(unsafe_code)]` at crate root unless there's an ADR-worthy reason otherwise
- **Memory safety fitness function**: `cargo clippy -- -D warnings` in CI

Commit: `docs(rules): add rust-conventions.md (Epic 50)`.

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- If async-trait situation has changed materially by the time this executes (Rust async traits stabilization) — halt, update the recommendation.
- If `tokio` isn't the right default for a specific workload (embedded, no-std) — note as an ADR-worthy sub-decision.

## Report (under 100 words)

```
Commit: <sha>
Sections covered: <list>
Cross-references: <list — probably devops-engineer for cargo/CI wiring>
Health-check green: yes
```

Go.
