# The Twelve-Factor App

Covers an axis nothing else in `docs/patterns/` addresses: not object design (Gang of Four), not code
layering (Clean Architecture), not domain modeling (DDD), not testing (Saturday/Sunday) — how a service
should actually be built to run well in production. Directly relevant to `devops-engineer`'s and
`sre-engineer`'s scope. Several factors already have partial coverage scattered across existing rules;
noted per factor rather than restated.

## I. Codebase

One codebase tracked in version control, many deploys. This repo's own `shared/` canonical-source model
(one definition, generated or symlinked into every platform/environment) is this principle applied to
configuration rather than application code — see `docs/ARCHITECTURE.md`.

## II. Dependencies

Explicitly declare and isolate dependencies — never rely on system-wide packages implicitly being
present. Each language's convention file already names the canonical package manager per language
(`shared/rules/go-conventions.md`'s Go modules, `python-conventions.md`'s `uv`, `csharp-conventions.md`'s
NuGet with Central Package Management, `typescript-conventions.md`'s pnpm) — all of which enforce this
factor by construction.

## III. Config

Store config in environment variables, never in code. Already a hard guardrail in this repo:
`shared/rules/architecture-guardrails.md` #3, "No Hardcoded Secrets" — "Use `.env` placeholders mapped to
secure vaults."

**Trade-offs**: Env-var config is harder to validate at startup than a typed config object — a typo in a
variable name fails silently until the code path that reads it actually runs, unless something validates
the full expected set of env vars eagerly at boot.

## IV. Backing Services

Treat backing services (databases, queues, caches, third-party APIs) as attached resources, swappable via
config without a code change. Structurally the same discipline as this repo's own "every external
dependency hides behind an interface" rule (`shared/rules/architecture-guardrails.md` #1's spirit, and
`CLAUDE.md`'s Non-Negotiable Rules) — a backing service swap should be a config change plus an adapter
swap, not a rewrite of business logic.

## V. Build, Release, Run

Strictly separate the build stage (compile/bundle), release stage (combine build + config), and run stage
(execute) — never modify code at runtime, and make every release a specific, referenceable artifact.

**Trade-offs**: Strict separation means a "hotfix directly on the running instance" is never an option,
which is uncomfortable during an incident but is exactly what makes a release reproducible and a rollback
possible — the alternative is a running system nobody can reconstruct from source control.

## VI. Processes

Execute the app as one or more stateless, share-nothing processes. Any state that needs to persist
belongs in a backing service (a database, a cache), not in memory or on local disk between requests.

**Trade-offs**: Statelessness means every request must carry or look up everything it needs — no
"remember what the last request did" shortcuts. That's what makes horizontal scaling and process restarts
safe without a coordination mechanism.

## VII. Port Binding

The app is self-contained and exports services via port binding — it doesn't rely on runtime injection of
a webserver into an execution environment; it *is* the webserver, listening on a port it owns.

## VIII. Concurrency

Scale out via the process model — add more processes/instances, rather than growing threads within one
process without bound. Pairs directly with Factor VI (Processes): this only works cleanly if processes
are actually stateless.

## IX. Disposability

Maximize robustness with fast startup and graceful shutdown — a process should be killable at any moment
without corrupting state, and should start fast enough that scaling up or recovering from a crash is
cheap.

**Trade-offs**: Graceful shutdown handling (finish in-flight requests, release connections cleanly) is
extra code that does nothing in the happy path — it only matters during a deploy or a crash, which is
exactly why it's easy to skip and exactly when skipping it hurts most.

## X. Dev/Prod Parity

Keep development, staging, and production as similar as possible — same backing services, same
dependency versions, minimal time gap between writing code and deploying it. This repo's own
`scripts/ci-check.sh` — which runs the framework's own checks inside a Docker container matching the
actual CI runner, specifically because a local macOS run had already once passed while the same code
failed in CI — is this factor applied to development tooling itself: don't trust "it works on my
machine" when the machine isn't the same as where it actually has to run.

## XI. Logs

Treat logs as event streams, not files the app manages itself — write to stdout, let the execution
environment handle routing/storage/rotation. Partially covered already:
`shared/rules/architecture-guardrails.md`'s Error Handling guidance and every language convention file's
"structured logging with low-cardinality message strings" rule cover log *content*; this factor is about
log *handling* — the app shouldn't be opening its own log files or managing its own rotation.

## XII. Admin Processes

Run one-off admin/maintenance tasks (migrations, console scripts) as one-off processes in an environment
identical to the app's regular long-running processes — never as ad hoc scripts run somewhere
different from where the app actually runs. `db-migration`'s expand/contract migrations
(`shared/rules/architecture-guardrails.md` #2) are exactly this kind of admin process — versioned,
run through the same pipeline as everything else, not a manual one-off against a database console.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
