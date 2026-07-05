# Agent Reference: Roles and Counterbalances

Every agent in this framework produces an output that *something* checks — another agent's review, a
structural contract, a human approval gate, an aggregate metric measured after the fact, or (stated plainly
where true) nothing yet. This doc makes that explicit for all 24 agents, one at a time, instead of leaving
it scattered implicitly across `deliver-feature/SKILL.md`, `shared/contracts/`, and each agent's own file.

This is documentation of what already exists in v2 today — it does not introduce new agents or roles. Where
an agent has no real counterbalance, that's stated as a gap, not filled in with something invented to make
the table look complete. (A larger, formal "checks and balances" model — pairing every role with a dedicated
auditor counterpart — is being prototyped separately for v3/AOS; see `docs/aos/`. This doc is about what v2
actually has running today.)

## How to read "Counterbalance"

Four different *kinds* of check show up across these 24 agents, and they're not interchangeable:

| Kind | What it catches | Example |
|---|---|---|
| **Structural contract** | Missing/malformed sections in the artifact itself | `validate-artifact` against `shared/contracts/*.md` |
| **Downstream agent review** | Wrong or low-quality content, live, on this delivery | `code-reviewer` reviewing `developer`'s output |
| **Human approval gate** | Irreversible or consequential actions | `shared/rules/approval-gates.md`'s 8 gates |
| **Aggregate/delayed metric** | Systemic drift across many deliveries, not visible in any one | `agent-scorecard`'s per-agent quality scores |

A structural contract can't catch bad judgment; a downstream review can't catch a slow, aggregate decline in
quality; a human gate only fires for the specific actions it names. Most agents below are checked by more
than one kind — that's by design, not redundancy.

---

## Pipeline Agents (in `deliver-feature` order)

### 1. `spec-writer`
**Role**: Interviews the user one question at a time to draft a feature spec, then immediately critiques its
own draft for downstream readiness (Write Mode → Review Mode) against every later agent's actual needs.
**Counterbalance**: Self-critique is built into its own process, not external — but `product-owner` reviews
the resulting spec next for value/scope, and `analyst` is the first agent to actually try to *use* the spec
technically; a spec that was too vague surfaces as a gap in `analysis.md`.
**Gap**: No agent double-checks whether spec-writer's own readiness critique was accurate — if it declared
READY prematurely, that's only caught downstream, not immediately.

### 2. `product-owner`
**Role**: Challenges whether a feature should be built at all — value, scope, simpler alternatives — before
any analysis or code happens. A DO NOT BUILD verdict blocks the pipeline.
**Counterbalance**: **Human approval gate** — a DO NOT BUILD or REDUCE SCOPE verdict requires explicit human
override to proceed, per its own guardrail.
**Gap**: Nothing checks whether product-owner's *own* pushback is well-calibrated over time (too aggressive,
killing legitimate work, vs. too lenient) — no aggregate metric exists for this agent the way `agent-scorecard`
tracks code-reviewer and security-reviewer.

### 3. `context-engineer`
**Role**: Pre-flight context optimizer — scopes the bounded context, prunes irrelevant files, surfaces
relevant KIs/ADRs and prior deliveries, estimates token budget, before analyst reasons independently.
**Counterbalance**: **Structural contract** (`shared/contracts/context-manifest-contract.md` via
`validate-artifact`, step 7) checks the manifest has all 7 required sections and never reports a token
budget without an estimate. Separately, the `context-audit` skill checks *after the fact* whether pinned
files were actually used — a different kind of check (waste analysis, not correctness).
**Gap**: None significant — this is one of the more thoroughly checked agents, deliberately, since everything
downstream depends on it.

### 4. `analyst`
**Role**: Turns a feature spec into acceptance criteria, task breakdown, data model/API changes, bounded
context mapping, and a definition of done — the first technical translation of the spec.
**Counterbalance**: **Structural contract** (`analysis-contract.md`) + a **human approval gate** (PAUSE, step
10) before any code is written. `architect`/`developer` downstream also surface gaps if the analysis was
wrong (e.g., a missing edge case shows up as rework).
**Gap**: None significant.

### 5. `architect`
**Role**: Structural decisions, fitness functions, layer boundaries, Team Topology fit — only invoked when a
feature actually needs structural judgment, not a rubber stamp on every feature.
**Counterbalance**: **Structural contract** (`architecture-contract.md`) + a **human approval gate** (PAUSE if
an RFC was written — expensive/permanent decisions require explicit acknowledgment) + `team-topology-check`
for Conway's-Law-shaped crossing mismatches. `code-reviewer` downstream also flags it if the implementation
didn't actually follow the architecture.
**Gap**: None significant.

### 6. `performance-engineer`
**Role**: Shift-left performance review of the *architecture* before implementation starts — idempotency,
timeouts, N+1 prevention, caching strategy.
**Counterbalance**: **Structural contract** (`performance-contract.md`, requires a Status line per risk
category). `qa-engineer` is the practical, downstream check on whether the mandated SLAs are actually met.
**Gap**: Nothing verifies performance-engineer's *estimates* were realistic (e.g., a caching recommendation
that turns out not to help) — no load-testing feedback loop exists yet.

### 7. `data-engineer`
**Role**: Schema design and zero-downtime (Expand/Contract) migrations for features touching the database.
**Counterbalance**: **Structural contract** (`data-engineering-contract.md`, requires a declared Phase) + the
`validate-migrations` skill (invoked directly by this agent, rejects destructive operations) + `devops-engineer`
independently invokes `validate-migrations` again before deploy.
**Gap**: None significant — migrations are one of the most heavily gated artifact types in the pipeline.

### 8. `developer`
**Role**: Implements the feature via TDD in an isolated worktree, following the analysis and architecture.
**Counterbalance**: **Structural contract** (`implementation-contract.md`) + **downstream agent review** from
`code-reviewer` (blocks on CHANGES REQUESTED, including an explicit "Test Design Review" of developer's own
TDD tests), `security-reviewer`, and `accessibility-engineer` — the most heavily reviewed handoff in the
pipeline by agent count.
**Gap**: None significant for the code itself. See `code-reviewer` below for a different, related gap.

### 9. `code-reviewer`
**Role**: Reviews the implementation against Clean Architecture, SOLID, Sandi Metz limits, and Fowler smells;
blocks the pipeline until APPROVED.
**Counterbalance**: **Structural contract** (`review-contract.md`, requires the literal `APPROVED`/`CHANGES
REQUESTED` string) + **aggregate metric** (`agent-scorecard`'s first-pass-acceptance-rate, flags if this
agent is systematically too lenient — rubber-stamping — or too strict, floor < 50%).
**Gap**: No agent re-reviews a single code-reviewer verdict *live* — the check is statistical and delayed
(next month's scorecard), not immediate. A rubber-stamped APPROVED on one specific delivery isn't caught
until the pattern shows up in aggregate.

### 10. `accessibility-engineer`
**Role**: Reviews UI changes for semantic HTML, WCAG compliance, keyboard navigation — fixes objective
violations directly. Only invoked when the feature touches UI.
**Counterbalance**: **Structural contract** (`accessibility-contract.md`, requires all 4 evaluation categories
present) + `qa-engineer` independently runs its own accessibility check (`[A11Y]` findings) as a second,
differently-sourced verification.
**Gap**: None significant — this is one of the few agents with a genuine second, independent check (qa-engineer's
own a11y pass) rather than just a structural contract.

### 11. `security-reviewer`
**Role**: STRIDE threat model of the implementation; fixes Critical/High findings directly rather than just
recommending them.
**Counterbalance**: **Structural contract** (`security-contract.md`) + **human approval gate** (explicit "fix
confirmed" required before QA starts on any Critical finding) + `qa-engineer` tests the specific security
scenarios security-reviewer calls out in "Notes for QA" + **aggregate metric** (`agent-scorecard`'s security
true-positive-rate proxy).
**Gap**: `agent-scorecard`'s own docs flag this explicitly — the TPR metric is a proxy, not a confirmed
true/false-positive rate, since there's no mechanism yet for a human or later agent to formally dispute a
finding after the fact.

### 12. `qa-engineer`
**Role**: Writes and runs tests covering every acceptance criterion and edge case; fixes bugs it finds.
**Counterbalance**: **Structural contract** (`qa-contract.md`, requires `Failed: 0`) + the `run-tests` skill
enforces the 85% coverage threshold mechanically, not just by qa-engineer's self-report.
**Gap**: Nothing independently reviews whether qa-engineer's tests are actually *meaningful* (verify real
behavior, not tautological assertions) the way `code-reviewer`'s "Test Design Review" does for developer's
TDD tests specifically — qa-engineer's own separate test suite has no equivalent second reviewer.

### 13. `sre-engineer`
**Role**: Observability review — defines SLIs, fixes high-cardinality logging, checks OTel span coverage and
PII hygiene.
**Counterbalance**: **Structural contract** (`observability-contract.md`).
**Gap**: Nothing verifies after deployment whether the SLIs sre-engineer defined were actually achievable or
meaningful in production — that feedback loop doesn't exist yet (it would need to compare defined SLIs
against real monitoring data, which is outside this framework's current scope).

### 14. `tech-writer`
**Role**: Updates all documentation for the delivered feature — CHANGELOG, README, ADRs, API docs.
**Counterbalance**: **Structural contract** (`docs-contract.md`) + tech-writer's own rule ("do not invent
behavior — only document what was actually implemented, verify with the code") is a self-verification
requirement, not an external check.
**Gap**: No agent fact-checks the resulting documentation against the actual implementation after the fact —
this relies entirely on tech-writer's own discipline. See `documentation-manager` below for a related,
partially-overlapping capability.

### 15. `devops-engineer`
**Role**: CI/CD, environment config, deployment artifacts — the final pipeline agent.
**Counterbalance**: **Structural contract** (`devops-contract.md`) + `validate-migrations` (invoked directly,
same as data-engineer) + **human approval gate** (Gate #8, "Deploying to Environment" — deploy requires
explicit "deploy" confirmation; Gate #2 covers the commit itself).
**Gap**: None significant — deploy is one of the most consequential actions in the whole pipeline and is
gated accordingly.

---

## Standalone / On-Demand Agents (not part of `deliver-feature`)

These aren't invoked automatically by the pipeline — they're triggered directly by a human or by a specific
condition (build time SLA breach, release request, etc.). None of them have a `shared/contracts/` entry,
since `validate-artifact` only gates pipeline handoffs.

### 16. `dependency-auditor`
**Role**: Audits the full dependency tree for vulnerabilities, license compliance, and unused packages.
**Counterbalance**: **Human approval gate** on *acting* on findings ("never auto-upgrade without explicit
user approval") — but this only gates the response, not the audit's own accuracy.
**Gap**: No agent or process double-checks the audit's findings themselves (e.g., a CVE dismissed as
inapplicable, or a license miscategorized). The report is trusted as read.

### 17. `release-manager`
**Role**: Analyzes git history since the last tag, applies semantic versioning from conventional commits,
drafts a release plan with rollback procedure.
**Counterbalance**: **Human approval gate** (its own rules: never tag without explicit "tag it"/"approve
release," never force-push, never skip CI) + its own checklist cross-checks that no open Critical/High
security findings exist before release, pulling from `security-reviewer`'s prior work.
**Gap**: None significant — this agent is unusually well-gated for a standalone tool.

### 18. `chaos-engineer`
**Role**: Designs and executes fault-injection experiments to verify the architect's resilience patterns
actually hold up.
**Counterbalance**: **Human approval gate** ("NEVER run chaos experiments against a live production database
without explicit multi-stage approval") for the risky *action*; its own report format requires checking
whether observability actually caught the injected fault, which indirectly validates sre-engineer's earlier
work too.
**Gap**: No dedicated reviewer checks the experiment *design* itself (is the hypothesis reasonable, is the
blast radius actually contained as claimed) before it runs — the approval gate covers the destructive action,
not the design's soundness.

### 19. `dx-engineer`
**Role**: Optimizes local development loop, build pipelines, flaky test quarantine.
**Counterbalance**: Its own guardrail requires measuring before/after impact quantifiably — a self-imposed
check, not external.
**Gap**: No agent reviews whether a DX fix (e.g., quarantining a flaky test) was the right call versus fixing
the underlying flakiness — this is a real tension (quarantine vs. fix) with no independent arbiter.

### 20. `finops-engineer`
**Role**: Reviews architecture/implementation for cost implications — treats cost as a fitness function
alongside latency and uptime.
**Counterbalance**: None formal — its own guardrail ("do not sacrifice resiliency or security to save
fractions of a cent") is self-imposed, not externally checked.
**Gap**: No agent or process verifies finops-engineer's cost estimates against actual cloud billing after
deployment — the "Cost Fitness Function" it proposes is a recommendation for a human to wire up, not
something this framework verifies itself.

### 21. `documentation-manager`
**Role**: Persistent, cross-session agent — extracts architectural decisions, debugging insights, and gotchas
from past sessions and updates long-lived docs (`ARCHITECTURE.md`, `RUNBOOKS.md`, `GOTCHAS.md`, `ONBOARDING.md`).
**Counterbalance**: None. No contract, no downstream review, no human gate beyond normal commit approval.
**Gap — and a real overlap worth naming plainly**: this agent's responsibility overlaps substantially with
`memory-engineer`, `promote-memory`, and `extract-lessons` (all added later, in the Memory Engineering epic) —
all four are, in different shapes, "extract durable knowledge from what just happened and write it somewhere
lasting." The boundary between them isn't currently documented anywhere, and `docs/GOTCHAS.md` (one of the
four files this agent is meant to maintain) doesn't exist in this repo — this capability appears dormant.
Worth resolving deliberately (either narrow this agent's scope explicitly, or retire it in favor of the
memory-engineering skills) rather than leaving both running with an undocumented boundary.

### 22. `modernization-supervisor`
**Role**: Coordinates three parallel modernization workstreams (dependency updates, pattern refactors, test
coverage) across a legacy codebase.
**Counterbalance**: Its own process requires running full-suite integration tests after coordinating the
three workstreams, and producing a risk-assessed merge strategy — an implicit **human approval gate** on the
merge itself (nothing in its own rules says it merges without asking).
**Gap**: No dedicated reviewer checks the *coordination* quality itself (did it actually prevent the 3
sub-workstreams from conflicting) beyond the integration test suite passing.

### 23. `api-test-generator`
**Role**: Generates Sunday Framework API test suites (Playwright + Vitest + Zod) from a spec or OpenAPI doc.
**Counterbalance**: If invoked as part of a full feature delivery, `code-reviewer`/`qa-engineer` would still
review the generated tests as part of that pipeline. If invoked standalone (its primary use case — "generate
API tests" on demand), there is none.
**Gap**: Standalone invocations have no formal check on the generated tests' quality beyond a human reading
the `api-test-report.md`.

### 24. `test-driven-developer`
**Role**: An autonomous, alternate red-green-refactor loop — explicitly authorized to iterate without asking
for permission between steps, working directly from user-provided acceptance criteria rather than through
`deliver-feature`.
**Counterbalance**: The test suite itself (mechanical — tests pass or they don't).
**Gap — significant, and worth understanding clearly before choosing this over the full pipeline**: this
agent deliberately bypasses every other counterbalance the normal pipeline has. No `code-reviewer` catches a
SOLID violation, no `security-reviewer` catches a STRIDE-class vulnerability, no `accessibility-engineer`
catches a semantic-HTML violation — none of that review happens unless a human separately runs it after the
fact. This is a real, deliberate speed-for-safety tradeoff (autonomy without pausing), not an oversight — but
it should be a conscious choice, not a default one, given how much of this framework's value elsewhere comes
from exactly the reviews this agent skips.

---

## What this survey actually shows

Reading all 24 agents together, three patterns stand out:

1. **The 14 pipeline agents are well-checked.** Every one has at least a structural contract; most have a
   real downstream reviewer or a human approval gate too. This is the part of the framework that's had the
   most iteration (three independent audits this session all confirmed the contract layer, once complete,
   holds up).
2. **The 9 standalone agents are inconsistently checked.** Some (`release-manager`, `data-engineer`-adjacent
   `devops-engineer`) are well-gated because they touch genuinely irreversible actions. Others
   (`finops-engineer`, `dx-engineer`, `dependency-auditor`) have no check beyond a human reading the report —
   which may be the right amount of ceremony for their risk level, or may not; that's a judgment call this
   doc surfaces rather than settles.
3. **`documentation-manager` vs. the memory-engineering skills is a real, unresolved boundary**, and
   `test-driven-developer`'s bypass of the whole review chain is a real, significant tradeoff — both are
   worth a deliberate decision rather than staying as-is by default.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
