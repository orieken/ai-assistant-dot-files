# Rules Excerpt (under audit)

## From shared/rules/testing-conventions.md
> NEVER write feature code without tests.
> Test coverage MUST be >= 85%.

## From shared/rules/typescript-conventions.md
> `strict: true` in `tsconfig.json`, no raw `any`

## From shared/rules/architecture-guardrails.md (#4)
> No raw `any` types allowed in TypeScript.
> If you genuinely don't know the type, use `unknown` with runtime narrowing/Zod validation.

## From shared/rules/design-principles.md
> Cyclomatic complexity < 7

## From CLAUDE.md
> Cyclomatic complexity < 7 on every function or method — no exceptions
> ESLint complexity rule: max 6 (enforces the < 7 rule)

## Potential issue
`shared/rules/go-conventions.md` references `DOMAIN_DICTIONARY.md` for ubiquitous language alignment, but:
- `design-principles.md §6` also references `DOMAIN_DICTIONARY.md`
- The path referenced in go-conventions.md is `docs/DOMAIN_DICTIONARY.md`
- The path referenced in design-principles.md is just `DOMAIN_DICTIONARY.md` (no path prefix)

## Un-indexed rule file
`shared/rules/iac-conventions.md` exists on disk but is NOT listed in any rules index or CLAUDE.md load order.
