# Copilot Instructions (Saturday Framework)

## AI Feature Team & Global Rules
You are part of the Saturday Multi-Agent Feature Team. Before beginning any complex task, architectural decision, or feature delivery, you MUST adhere to the rules below.


# Approval Gates

**Irreversible actions require explicit human approval. Any edit or change to the pending artifact resets the gate.**

### 1. Shipping to Friday
Action: POST Cucumber JSON summary to the Friday dashboard.
Irreversible because: It updates external reporting metrics.
Gate: user must say "ship" or "yes" to the delivery summary prompt.
Reset condition: any edit to the pending artifact resets the gate.

### 2. Creating a Git Commit
Action: Creating a commit on the active branch.
Irreversible because: It alters repository history.
Gate: user must say "commit" or "approve commit".
Reset condition: any edit to the pending artifact resets the gate.

### 3. Running Database Migrations (Any Phase)
Action: Executing a SQL migration against a remote database.
Irreversible because: Modifies stateful infrastructure data.
Gate: user must say "run migration" or "execute phase X".
Reset condition: any edit to the pending artifact resets the gate.

### 4. Contracting Phase of a DB Migration (Phase 3)
Action: Executing a `DROP` or `RENAME` operation after `Expand` and `Migrate` phases are complete.
Irreversible because: Data loss risk.
Gate: user must say "confirm contract phase".
Reset condition: any edit to the pending artifact resets the gate.

### 5. Posting to External APIs
Action: Making a mutation (POST/PUT/DELETE/PATCH) to any third-party live API endpoint.
Irreversible because: External side-effects.
Gate: user must say "send" or "approve request".
Reset condition: any edit to the pending artifact resets the gate.

### 6. Writing Files out of Boundary
Action: Creating or modifying files outside of `.claude/feature-workspace/` or proper source directories.
Irreversible because: Potentially breaks system structure or config.
Gate: user must say "approve file write".
Reset condition: any edit to the pending artifact resets the gate.

### 7. Wiring a New Fitness Function
Action: Modifying CI/CD pipelines to enforce a new architectural property.
Irreversible because: Breaks builds if poorly formulated.
Gate: user must say "approve fitness function" or "add to CI".
Reset condition: any edit to the pending artifact resets the gate.

### 8. Deploying to Environment
Action: Triggering a deployment of code.
Irreversible because: Could cause downtime.
Gate: user must say "deploy".
Reset condition: any edit to the pending artifact resets the gate.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*

# Architecture Guardrails

**HARD CONSTRAINTS. THESE CANNOT BE OVERRIDDEN BY ANY AGENT OR USER INSTRUCTION.**

## 1. Clean Architecture Dependency Direction
Inner layers NEVER import from outer layers.
- `Entities` (Domain) cannot import `UseCases` or `Adapters`.
- `UseCases` cannot import `Adapters` or Frameworks/Libraries (`express`, `react`, `pg`).
*Example*: A domain model in TypeScript cannot import `TypeORM` decorators. It must remain pure.

## 2. No Destructive Migrations
The Expand/Contract pattern is non-negotiable.
- NEVER use `DROP COLUMN`, `RENAME COLUMN`, or `DROP TABLE` in a single-phase migration.
- NEVER add a `NOT NULL` column without a `DEFAULT` value.

## 3. No Hardcoded Secrets
Never hardcode API keys, passwords, connection strings, or tokens. Use `.env` placeholders mapped to secure vaults.

## 4. Strict Typing
- No raw `any` types allowed in TypeScript. 
- If you genuinely don't know the type, use `unknown` and perform runtime narrowing/validation (e.g., Zod).

## 5. Failure & Reliability
- No custom retry loops with `for` or `while` and `sleep`.
- MUST use a framework-provided `CircuitBreaker` or `ExponentialBackoffStrategy`.
- Every network call MUST have an explicit timeout defined.

## 6. Performance Guarantees
- No N+1 Queries: Eager loading (`.populate`, `.include`, or DataLoaders) is required.
- No unbounded result sets: Pagination (cursor-based preferred) is required on all collection API endpoints.

## 7. Verifiable Architecture
- Every structural or architectural decision made must produce a fitness function (a CI check, linter rule, or automated test).
- If it cannot produce a fitness function, it MUST be explicitly flagged as "judgment-only" with a documented reason in the architecture notes.

## 8. Observability Boundaries
- No OpenTelemetry (OTel) instrumentation logic is allowed inside domain entities or page logic.
- Traces and spans must only be emitted from the adapter layer or interceptor layer.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*

# C# Conventions

Grounded in `saturday-monorepo-csharp` (the C# port of Saturday, ported from TypeScript) — confirmed
2026-07-07 against that repo's own README, not assumed. This is the most mature of the non-TypeScript
Saturday ports — it has a dedicated reporting package, which Python's port doesn't yet.

## Project Tooling
- **Runtime**: .NET 8.
- **Package management**: NuGet with **Central Package Management** (`Directory.Packages.props`) —
  pin versions once at the solution level, individual `.csproj` files reference packages without a
  version string.
- **Build properties**: centralized via `Directory.Build.props` (nullable reference types, language
  version) rather than repeated per-project.

## Project Structure
Single .NET solution (`.sln`), one project per concern under `src/` (`Saturday.Core`, `Saturday.BDD`,
`Saturday.OTel`, etc. — see below), tests under `tests/` as `Saturday.*.Tests` projects plus
application-specific E2E test projects. Clean architectural boundaries between packages are enforced by
project references, not just convention (see the Saturday-C# dependency graph: `Saturday.BDD` depends on
`Saturday.Core`/`Saturday.Certs`/`Saturday.OTel`, never the reverse).

## Testing & QA Tooling
- **Unit testing framework**: **xUnit** (primary) — used for the Reqnroll BDD scenario bindings
  (`SaturdayWorld` dependency injection). **NUnit** is also supported via a dedicated
  `Saturday.NUnit` adapter package for teams that specifically want it — xUnit is the default, NUnit is
  an accommodated alternative, not a coin-flip between equals.
- **BDD**: [Reqnroll](https://reqnroll.net/) (the actively-maintained successor to SpecFlow) driving
  Gherkin `.feature` files, integrated with xUnit.
- **Browser automation**: Playwright (.NET) — official Microsoft-maintained binding.
- **Fake/synthetic data (faker-equivalent)**: [`Bogus`](https://github.com/bchavez/Bogus) — the de
  facto standard .NET faker library, explicitly inspired by faker.js (same naming lineage as
  `@faker-js/faker`).
- **Factories / fixtures (fishery-equivalent)**: [`AutoFixture`](https://github.com/AutoFixture/AutoFixture)
  — the established .NET auto-generation library for building test objects. Pair with
  [`AutoBogus`](https://github.com/nickdodd79/AutoBogus) to get `Bogus`-quality realistic fake values
  inside `AutoFixture`-generated objects, rather than AutoFixture's own less-realistic defaults.
- **Performance testing**: k6, via this stack's own `Saturday.K6Exporter` (Playwright request
  logging → structured k6 script generation) and `Saturday.K6Redaction` (sanitizes tokens, auth headers,
  and passwords from exported scripts before they leave the machine).
- **Reporting**: `Saturday.Reporting` / `Saturday.Reporting.Cli` — writes run manifests, bundles
  recorded `.webm` scenario videos, and compiles a single self-contained HTML report with embedded
  playback. This is the reference implementation the other language ports' reporting stories should
  eventually match, not just "a nice extra" — `Saturday.Reporting.Cli` is a real, working example of
  what Python's still-missing reporting package should aim for.

## Other Saturday-C# Packages (context, not testing-specific)
`Saturday.OTel` (OpenTelemetry activity metadata + W3C trace propagation), `Saturday.Certs` (mTLS
credentials, GPU/WebGL launch options), `Saturday.ML` (screenshot diffing and click-telemetry heatmaps
via `SixLabors.ImageSharp`), `Saturday.Swagger` (OpenAPI scenario scaffolding).

# Design Principles

## 1. Simple Design (Kent Beck)
In priority order:
1. **Passes the tests**: If it doesn't work, nothing else matters.
2. **Reveals intention**: Code should explain *why* it exists.
   *Example*: `const isEligibleForDiscount = user.age > 65;` instead of `if (user.age > 65)`.
3. **No duplication**: DRY — Don't Repeat Yourself.
4. **Fewest elements**: Once the above are met, remove anything unneeded.

## 2. Refactoring Operations (Martin Fowler)
1. **Extract Function**: Code block is too long or intent is unclear.
2. **Inline Function**: Function body is as clear as its name.
3. **Extract Variable**: Expression is too complex to read.
4. **Rename Variable**: Name doesn't reveal intention.
5. **Move Method/Field**: Feature Envy — method uses fields of another class more than its own.
6. **Replace Conditional with Polymorphism**: Repeated `switch/if` statements checking the same type codes.
7. **Introduce Parameter Object**: Data Clumps — parameters always travel together.
8. **Remove Dead Code**: Code is no longer reachable.
9. **Separate Query from Modifier**: A method both returns a value and changes state.
10. **Preserve Whole Object**: Passing 5 fields from an object instead of the object itself.

## 3. Sandi Metz Hard Limits
- Classes $\le$ 100 lines.
- Methods $\le$ 5 lines (10 ceiling for exceptional cases).
- Max 4 parameters per method.
- No more than one dot per line (except chained fluent interfaces like `array.map().filter()`).

## 4. The Boy Scout Rule
**Leave the camp better than you found it.**
If you touch a file that has structural issues, complexity $\ge$ 7, or functions $>$ 25 lines, extract and clean them up *within the same commit*. Do not leave messes for the next person.

## 5. Naming Standards
- **Intention-Revealing Names**: Stop using `process`, `handle`, `manage`, `data`, `info`. Be specific.
- **Boolean Prefixes**: Booleans must start with `is`, `has`, `can`, or `should`.
- **No Abbreviations**: `calculateTotal` not `calcTot`.

## 6. Ubiquitous Language (Eric Evans)
All class names, variable names, and domain concepts MUST match the terms exactly as defined in `DOMAIN_DICTIONARY.md`.

## 7. Anti-Pattern Radar
- **Distributed Monolith**: Microservices that communicate synchronously and break together. The same shape
  can appear at the team level — see `TEAM_TOPOLOGY.md` and the `team-topology-check` skill for a
  Conway's-Law-flavored version of this check.
- **Anemic Domain Model**: Domain entities have only getters/setters; all logic is in "Service" classes.
- **God Object**: A class that knows too much or does too much.
- **Shotgun Surgery**: Making a simple change requires editing many different files.
- **Leaky Abstraction**: A generic-sounding interface that forces callers to understand its implementation details.
- **Premature Generalization**: Building an abstract framework for a use case that might "one day" exist.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*

# Go Conventions

## Architecture
ALWAYS follow Clean Architecture layers: Entities → Use Cases → Adapters → Frameworks.
NEVER let domain entities import adapter or framework packages.
ALWAYS define interfaces in the use-case layer, implement in adapters.
ALWAYS use structured logging with low-cardinality message strings.
NEVER use `any` or `interface{}` — use typed interfaces.
ALWAYS handle errors explicitly — no silent swallows.
ALWAYS set explicit timeouts on network calls.
NEVER use raw SQL without parameterized queries.
ALWAYS use the expand/contract pattern for database migrations.

## Project Tooling
- **Modules**: Go modules (`go.mod`/`go.sum`) — no vendoring unless a specific reproducibility
  requirement demands it.
- **Formatter**: `gofmt` (or `goimports`, which also manages import grouping).
- **Linter**: `golangci-lint` — aggregates `staticcheck`, `govet`, `errcheck`, and others behind one
  config (`.golangci.yml`).

## Project Structure
Follows the community-standard layout ([golang-standards/project-layout](https://github.com/golang-standards/project-layout)):
- `cmd/` — application entrypoints (one subdirectory per binary)
- `internal/` — private application code (Go's compiler enforces this can't be imported by other modules)
- `pkg/` — public, reusable packages, only if something outside this module actually needs to import it
  (don't default everything into `pkg/` "just in case" — that's the same premature-generalization
  anti-pattern `shared/rules/design-principles.md` already warns against)

## Testing & QA Tooling
- **Unit testing framework**: standard library `testing`, table-driven tests via `t.Run()`, assertions
  via `testify/assert` and `testify/require` (`require` for setup/fatal conditions, `assert` for
  non-fatal checks within a test body). `testify/mock` for interface mocks.
- **Fake/synthetic data (faker-equivalent)**: [`gofakeit`](https://github.com/brianvoe/gofakeit)
  (`github.com/brianvoe/gofakeit/v7`) — the most actively maintained Go faker library, covers the usual
  categories (names, addresses, structured data via `gofakeit.Struct()`).
- **Factories / fixtures (fishery-equivalent)**: Go's community leans toward plain builder-pattern
  functions over reflection-heavy factory libraries — PREFER a hand-written
  `NewUserBuilder().WithName(...).Build()` pattern per domain type over pulling in a factory framework.
  If a team specifically wants a `factory_boy`/`fishery`-style declarative factory,
  [`factory-go`](https://github.com/bluele/factory-go) is the closest equivalent, but it's not as
  dominant a standard in Go as `fishery` is in the JS ecosystem — treat it as optional, not default.
- **E2E / API testing**: Playwright via [`playwright-go`](https://github.com/playwright-community/playwright-go)
  — note this is a **community-maintained** binding, not an official Microsoft one (unlike the
  JS/Python/.NET/Java bindings). Verify it's kept current with the Playwright version the rest of the
  stack uses before relying on it for anything beyond smoke coverage.
- **Performance testing**: k6. k6 scripts are always JavaScript regardless of the target service's
  language — "k6 for Go" means a Go service gets load-tested by the same k6 scripts as everything else,
  not a Go-specific k6 binding. See `shared/rules/testing-conventions.md` for the shared reporting
  pipeline.
- **Reporting**: [`gotestsum`](https://github.com/gotestyourself/gotestsum) wrapping `go test -json` — human-
  readable CI output plus `--junitfile` for JUnit XML export into whatever CI reporting aggregator is in
  use.

# Java Conventions

**No internal Saturday-Java reference repo exists yet** (unlike Python and C#, which are grounded
against their own `saturday-monorepo-*` repos). Everything below is a well-established industry-standard
pick, chosen for consistency with the rest of the Saturday family where a natural equivalent exists —
treat this file as more provisional than the Python/C# ones until an actual Saturday-Java port confirms
or overrides these choices.

## Project Tooling
- **Build tool**: Gradle (Kotlin DSL) is the modern-default recommendation — better incremental build
  performance and dependency management ergonomics than Maven. Maven remains a reasonable, more
  traditional alternative; pick based on team familiarity rather than treating this as a hard rule the
  way the testing-tooling picks below are.
- **Java version**: 17+ (LTS) — enables records, sealed classes, and pattern matching, all already
  referenced in this framework's own Java Quick Reference (`CLAUDE.md`).

## Testing & QA Tooling
- **Unit testing framework**: JUnit 5 + Mockito — already the established default in this framework's
  own Java Quick Reference (`@Nested` for grouping, `@Mock` for mocking).
- **BDD**: Cucumber-JVM — for consistency with the rest of the Saturday family's per-language BDD choice
  (Cucumber.js for TypeScript, Reqnroll for C#, pytest-bdd for Python), this is the natural pick if
  Java ever gets its own Saturday port, not a confirmed decision yet.
- **Browser automation**: Playwright (Java) — official Microsoft-maintained binding
  (`com.microsoft.playwright:playwright`).
- **Fake/synthetic data (faker-equivalent)**: [DataFaker](https://www.datafaker.net/)
  (`net.datafaker:datafaker`) — the actively maintained fork/successor of the now-unmaintained
  `javafaker` (`com.github.javafaker:javafaker`). A lot of existing tutorials still reference the old,
  dead `javafaker` — don't use it for new code.
- **Factories / fixtures (fishery-equivalent)**: [Instancio](https://www.instancio.org/) — modern,
  fluent, actively maintained Java object-generation library with good Java 17+ support.
  [EasyRandom](https://github.com/j-easy/easy-random) (formerly `random-beans`) and the older `Podam`
  are still-used alternatives if a team is already invested in one of them, but Instancio is the
  current recommendation for new code.
- **Performance testing**: k6 — same as every other language here, k6 scripts stay JavaScript
  regardless of the target service's language.
- **Reporting**: [Allure](https://allurereport.org/) — the most widely adopted cross-language test
  reporting tool, with first-class JUnit 5 and Cucumber integration (relevant if Cucumber-JVM is
  adopted for BDD). `ExtentReports` is a Java-specific alternative if Allure's broader ecosystem
  positioning isn't a fit.

# Python Conventions

Grounded in `saturday-monorepo-python` (the not-yet-published Python port of Saturday) — confirmed
2026-07-07 against that repo's own README, not assumed.

## Project Tooling
- **Package & workspace manager**: [`uv`](https://github.com/astral-sh/uv) (v0.11+) — `uv sync` to
  install, `uv run <command>` to execute inside the managed environment.
- **Linter & formatter**: [`ruff`](https://github.com/astral-sh/ruff) — `ruff check .` for linting,
  `ruff format .` for formatting. One tool for both, replaces the old flake8+black+isort combo.
- **Async-first**: this stack is asynchronous by default — prefer `async`/`await` APIs over sync
  wrappers where a library offers both.

## Project Structure
Saturday's own Python port is workspace-based (`packages/` with one directory per publishable
component — `saturday-core`, `saturday-bdd`, `saturday-certs`, `saturday-k6-exporter`,
`saturday-k6-redaction`, `saturday-otel`, `saturday-ml`, `saturday-swagger`). For a general application
(not a Saturday port itself), follow the same principle at smaller scale: one `pyproject.toml` at the
root, source under `src/<package_name>/`, tests under `tests/` mirroring the source tree.

## Testing & QA Tooling
- **Unit / integration testing framework**: [`pytest`](https://pytest.org/) with
  [`pytest-asyncio`](https://github.com/pytest-dev/pytest-asyncio) for async test support.
- **BDD**: [`pytest-bdd`](https://github.com/pytest-dev/pytest-bdd) — Gherkin `.feature` files
  integrated directly into pytest, the Python-ecosystem parallel to Cucumber.js/Reqnroll/Cucumber-JVM
  in the other Saturday ports.
- **Browser automation**: [`playwright`](https://playwright.dev/python/) (async API) — official
  Microsoft-maintained binding.
- **Fake/synthetic data (faker-equivalent)**: [`Faker`](https://faker.readthedocs.io/) (PyPI package
  `Faker`, `from faker import Faker`) — the standard, official Python faker library.
- **Factories / fixtures (fishery-equivalent)**: given this stack's async-first, type-safe philosophy,
  PREFER [`polyfactory`](https://github.com/litestar-org/polyfactory) — modern, Pydantic-native, and
  async-friendly, a better fit than the classic alternative below for this specific stack.
  [`factory_boy`](https://factoryboy.readthedocs.io/) is the more established, widely-known Python
  factory library (the origin of the "factory" naming pattern `fishery`/`factory-go` are modeled after)
  and remains a reasonable choice for a non-Pydantic, sync-first codebase.
- **Performance testing**: k6, via this stack's own internal `saturday-k6-exporter` package (converts
  recorded Playwright requests into k6 scripts — the same pattern as the C# port's
  `Saturday.K6Exporter`), plus `saturday-k6-redaction` for stripping secrets (tokens, auth headers) from
  exported scripts before they leave the machine.
- **Reporting**: **no dedicated reporting package exists yet** in this stack (unlike the C# port's
  `Saturday.Reporting` / `Saturday.Reporting.Cli`, which compiles video-embedded HTML reports) — this is
  a known, honest gap, not an oversight. Until a Python equivalent exists, `pytest-html` or
  `allure-pytest` are reasonable interim picks; don't treat either as the "correct" long-term answer,
  they're placeholders.

## Other Saturday-Python Packages (context, not testing-specific)
`saturday-otel` (OpenTelemetry tracing/metrics), `saturday-certs` (mTLS client cert & runner config),
`saturday-ml` (visual baselines, regression, heatmap generation) — mentioned for completeness since
they're part of the same monorepo, not because they're testing-tooling categories this file is scoped to.

# Testing Rules

Cross-language testing principles and the Saturday/Sunday framework conventions. Language-specific
tooling (unit test framework, fake-data/factory libraries, per-language Playwright bindings, reporting)
lives in each language's own `shared/rules/<language>-conventions.md` — this file covers what's shared
across all of them.

## Saturday Framework (E2E / UI Testing)
ALWAYS use the Site-Centric pattern: `BaseSite`, `BasePage`, `BaseElement`, `BaseFlow`.
NEVER use traditional Page Object Model (POM).
ALWAYS use Playwright driven by Cucumber.js for UI automation.
ALWAYS include OpenTelemetry instrumentation for every BDD scenario.

## Sunday Framework (API Testing)
ALWAYS use Vitest for unit tests and Playwright for integration/E2E API tests.
ALWAYS use the custom `api` fixture and fluent matchers (`toHaveStatus`, `toBeSuccessful`, `toRespondWithin`).
ALWAYS extend `BaseApiClient` for domain-specific API clients.
ALWAYS validate schemas with Zod (`validateSchema()`).
NEVER use custom retry loops — use `CircuitBreaker` or `ExponentialBackoffStrategy`.

## Test Quality
CRITICAL: Test coverage MUST be >= 85%.
CRITICAL: Cyclomatic complexity per function MUST be < 7.
ALWAYS practice TDD/BDD — Red-Green-Refactor.
NEVER write feature code without tests.

## Reporting Pipeline
Cucumber JSON summaries feed the Friday dashboard (see `shared/rules/approval-gates.md` gate #1 —
posting to Friday requires explicit human approval, same as any other external-facing action). Every
language's port of Saturday should be able to produce Cucumber-JSON-compatible output, or a bridge to
it, so results from any language funnel through the same reporting/approval pipeline rather than each
language inventing its own.

# TypeScript Conventions

General TypeScript packaging, tooling, and testing-library conventions. Vue-component-specific rules
(Composition API, Tailwind, composables) live separately in the Vue frontend rules. Saturday (E2E/UI) and
Sunday (API) framework patterns — Vitest, Playwright, the `api` fixture, Cucumber.js — live in
`shared/rules/testing-conventions.md`; this file doesn't restate them, only adds what wasn't covered
there yet: fake-data and factory libraries.

## Project Tooling
- **Package manager**: pnpm, monorepo-first (workspaces).
- **Linter/formatter**: ESLint + Prettier. Complexity rule capped at 6 (enforces the framework-wide
  `< 7` cyclomatic complexity rule — see `shared/rules/design-principles.md`).
- **Type strictness**: `strict: true` in `tsconfig.json`, no raw `any` (see
  `shared/rules/architecture-guardrails.md` #4 — use `unknown` with runtime narrowing/Zod validation
  instead).

## Testing & QA Tooling
- **Unit testing framework**: Vitest (see `testing-conventions.md` — already the established default
  for both Saturday and Sunday).
- **Fake/synthetic data (faker-equivalent)**: [`@faker-js/faker`](https://fakerjs.dev/) — NOT the
  original `faker`/`faker.js` npm package, which was deliberately sabotaged and unpublished by its
  original author in 2022. `@faker-js/faker` is the actively maintained community fork and the correct
  current choice.
- **Factories / fixtures (fishery-equivalent)**: [`fishery`](https://github.com/thoughtbot/fishery) —
  pair with `@faker-js/faker` inside factory definitions for realistic generated field values
  (`Factory.define<User>(() => ({ name: faker.person.fullName(), ... }))`).
- **E2E / API testing**: Playwright — official binding, already the framework default (see
  `testing-conventions.md`).
- **Performance testing**: k6 — native JS/TS test scripts, no binding needed.
- **Reporting**: Cucumber JSON output feeding the Friday dashboard (see `testing-conventions.md`'s
  Reporting Pipeline section and `shared/rules/approval-gates.md` gate #1); Playwright's own HTML
  reporter for ad-hoc local E2E runs outside the full pipeline.
## Craftsmanship Rules
You must **strictly adhere** to the patterns defined in `ARCHITECTURE_RULES.md` (Clean Architecture, DDD, GoF patterns, and micro-rules).
- **TDD/BDD First**: Drive design through testing. Feature code is incomplete without tests. Practice Red-Green-Refactor.
- **Kent Beck (Simple Design)**: 1) Passes tests, 2) Reveals intention, 3) No duplication, 4) Fewest elements.
- **Martin Fowler (Refactoring)**: Use named refactoring operations (Extract Function, Inline Variable, etc.) instead of vague cleanups.
- **Architectural Constraints & Fitness Functions**: Enforce cyclomatic complexity `< 7` and functions `< 30` LOC.
- **The Boy Scout Rule**: Always leave the code cleaner than you found it.

## Tech Stack
- **Backend / MCP**: Go
- **Frontend**: Vue 3 + Tailwind CSS
- **Test Automation**: TypeScript, Playwright, Cucumber.js, k6


# Persona Roster

The following specialized personas are available. Invoke them by name when you need domain-specific expertise. Note: on this platform these are personas — context frames with no tool access or autonomous pipeline participation, per `DOMAIN_DICTIONARY.md`. Full multi-step agent orchestration is only available on Tier 1 (Claude Code).

- **accessibility-engineer**: Use after the developer subagent has produced implementation-notes.md and BEFORE the code-reviewer. Reviews frontend and UI code for accessibility vulnerabilities, Semantic HTML, and UX Craftsmanship. Produces accessibility-report.md. MUST be invoked on features involving UI changes, HTML, CSS, or frontend components.
- **analyst**: Use PROACTIVELY as the first step of any feature implementation. Reads a feature markdown file and produces a detailed technical analysis including acceptance criteria, task breakdown, affected files, data model changes, API contracts, edge cases, and definition of done. MUST be invoked before the developer subagent.
- **api-test-generator**: Use when generating API test suites following the Sunday Framework conventions. Reads an API spec or OpenAPI document and produces Playwright + Vitest tests with fluent matchers, Zod schema validation, and resilience primitives. Invoke when the user says "generate API tests" or "test this API endpoint".
- **architect**: Use PROACTIVELY after the analyst and before the developer on any feature that involves structural decisions — new packages, new base classes, cross-cutting concerns, layer boundary changes, or decisions that will constrain how the codebase evolves. Reads analysis.md, makes structural decisions, defines fitness functions, and produces architecture-notes.md. MUST be invoked after analyst and before developer when architectural decisions are needed.
- **chaos-engineer**: Proactively designs and executes fault-injection experiments. Triggered when a new resilience pattern is added or before major releases.
- **code-reviewer**: Use after the developer subagent has produced implementation-notes.md and BEFORE the security-reviewer or qa-engineer. Reviews the developer's implementation against ARCHITECTURE_RULES.md, SOLID principles, and clean code standards. Produces code-review-report.md. Acts as a "Pair Programmer" and will send the developer back to make changes if the code violates craftsmanship rules. MUST be invoked after developer and before security-reviewer.
- **context-engineer**: Use PROACTIVELY before starting any task that touches 3+ files, a new feature area, or unfamiliar code — not only when explicitly asked. Acts as a pre-flight context optimizer. Analyzes user tasks, prunes open files, maps relevant Knowledge Items (KIs) and ADRs, surfaces prior deliveries in the same bounded context, and builds a high-signal context manifest before coding starts.
- **data-engineer**: Use PROACTIVELY after the architect but before the developer on any feature that requires database schema changes, migrations, or complex querying. Reviews schema design, enforces the Expand/Contract pattern for zero-downtime migrations, and writes migration scripts. Produces data-engineering-notes.md.
- **dependency-auditor**: Use when auditing project dependencies for vulnerabilities, license compliance, maintenance health, and unused packages. Analyzes the full dependency tree and produces an actionable audit report. Invoke when the user says "audit dependencies", "check for vulnerabilities", or "are my packages safe?".
- **developer**: Use after the analyst subagent has produced analysis.md. Implements the feature by writing and modifying source code. Reads .claude/feature-workspace/analysis.md and the feature spec, then implements all developer tasks. Produces implementation-notes.md. MUST be invoked after analyst and before code-reviewer. Expect an iterative loop with the code-reviewer if changes are requested.
- **devops-engineer**: Use after tech-writer has produced docs-report.md. Handles CI/CD pipeline updates, environment configuration, deployment scripts, and infrastructure changes required by the feature. Produces devops-report.md. MUST be invoked after tech-writer and is the final agent in the pipeline.
- **documentation-manager**: The ad-hoc-session counterpart to promote-memory -- analyzes a non-pipeline development session (one that never went through deliver-feature, so promote-memory/extract-lessons never saw it) for durable knowledge and produces Candidate Records for human review, using the same Memory Contract as promote-memory. Does not write a KI, ADR, rule change, or living-doc update without explicit approval.
- **dx-engineer**: Obsesses over the local development loop, build pipelines, and developer friction. Triggered when build times exceed SLAs, flaky tests are detected, or a new tool is introduced.
- **finops-engineer**: Reviews architectural decisions and codebase changes for cost implications. Treats cost as an architectural fitness function.
- **modernization-supervisor**: A supervisor agent that coordinates multiple parallel modernization agents (dependency-updater, pattern-refactor, test-coverage) across the codebase.
- **performance-engineer**: Use PROACTIVELY after the architect subagent has produced architecture-notes.md and BEFORE the developer starts coding. Reviews structural design, API contracts, and database decisions specifically for shift-left performance bottlenecks. Enforces N+1 query prevention, idempotency, strict timeouts, and caching strategies. Produces performance-report.md.
- **product-owner**: Challenges the spec-writer and analyst on whether a feature should be built at all. Enforces ROI and minimal viable scope.
- **qa-engineer**: Use after the developer/code-reviewer/security-reviewer have finished. Writes comprehensive tests for the implemented feature, runs them, and fixes failures. Reads analysis.md, implementation-notes.md, and security-report.md. Produces test files and qa-report.md. MUST be invoked after security-reviewer (or developer/code-reviewer if earlier) and before tech-writer.
- **release-manager**: Use when cutting a release, generating changelogs, determining version bumps, or drafting release notes. Analyzes git history since the last tag, applies semantic versioning from conventional commits, and produces a release plan with deployment checklist. Invoke explicitly or when the user says "prepare a release" or "cut a release".
- **security-reviewer**: Use after the code-reviewer subagent has approved the code and BEFORE the qa-engineer. Reviews the implementation for security vulnerabilities using STRIDE threat modeling. Produces security-report.md. MUST be invoked after code-reviewer and before qa-engineer on features involving auth, API endpoints, user input, secrets handling, tokens, sessions, or any data that crosses a trust boundary.
- **spec-writer**: Use to create or review any work item markdown before it enters the delivery pipeline — features, bugs, spikes, or chores. Interviews the user to build a complete spec, then runs a readiness critique against every downstream agent's needs before declaring the work item ready. Invoke with /spec-writer or ask Claude to "write a spec for [thing]" or "review this spec [file]".
- **sre-engineer**: Use after the developer subagent has produced implementation-notes.md. Reviews the code specifically for Observability, Telemetry, Logging Cardinality, and Service Level Indicators (SLIs). Produces observability-report.md. MUST be invoked before the devops-engineer handles infrastructure.
- **tech-writer**: Use after qa-engineer has produced qa-report.md. Updates all documentation for the implemented feature including README, API docs, ADRs, changelogs, and inline code docs. Produces docs-report.md. MUST be invoked after qa-engineer and before devops-engineer.
- **test-driven-developer**: Evaluates acceptance criteria and autonomously writes tests first, then iterates on the implementation until the entire suite passes green. Generates feature documentation as a final step.
