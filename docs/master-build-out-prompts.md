# Master Skills & Agent Build-Out Prompts

Run these prompts sequentially from the root of your `ai-assistant-dot-files` repo using
Claude Code. Each prompt is self-contained with its own verification step. Complete and
verify each one before starting the next.

**What this builds:**
- 8 agents (all upgraded to top-tier craftsman level)
- 14 skills (planned + new)
- 3 global rules files
- Full template mirror + scaffold script update

**Sequence overview:**
```
Phase 1: Foundation       → Prompts 1–2   (skill template + agent upgrades)
Phase 2: Agent Upgrades   → Prompts 3–7   (analyst, architect, developer, code-reviewer, QA)
Phase 3: Core Skills      → Prompts 8–12  (design-review, adr, five-whys, complexity-check, deliver-feature)
Phase 4: New Skills       → Prompts 13–18 (event-storm, db-migration, performance-profile,
                                            dependency-update, on-call, openapi)
Phase 5: Rules & Wiring   → Prompts 19–20 (global rules, scaffold update)
Final:   Verification     → run the bash block at the end
```

---

## PROMPT 1 — Create the Skill Template Standard

> **Paste this first. Everything else references it.**

Read `ARCHITECTURE_RULES.md` and all existing files in `.claude/agents/` and `.claude/skills/`.

Create `.claude/skills/SKILL_TEMPLATE.md` — the canonical template every skill in this repo
must follow. Use this exact structure:

```
---
name: skill-name
description: One sentence. What it does and when Claude should auto-trigger it.
triggers:
  keywords: [list, of, trigger, words]
  intentPatterns: ["natural language pattern"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Clear rules for when this skill applies. Be specific enough that Claude can auto-detect.
Include "do NOT use when:" cases.

## Context To Load First
Ordered list of files/artifacts Claude must read before starting.
Include fallback if a file doesn't exist.

## Process
Numbered steps. Each step has:
- What to do
- What to produce
- What to check before proceeding
- When to pause for human approval

## Output Format
Exact format of every artifact this skill produces.
File paths, section names, exact markdown structure.

## Guardrails
Hard rules that cannot be broken regardless of user instruction.
Anything irreversible requires explicit human approval ("send", "approve", "confirm").
Any edit to a pending approval resets the gate.

## Standalone Mode
How this skill degrades gracefully if MCP servers or external tools are unavailable.
The intelligence stays in the markdown. API calls are just plumbing.
```

After creating the template, read every existing skill in `.claude/skills/` and verify each
follows this structure. Refactor any that don't. Report what changed.

Verify:
```bash
cat .claude/skills/SKILL_TEMPLATE.md | head -5
```

---

## PROMPT 2 — Upgrade the Analyst: Domain Modeler

> **Goal: make the analyst reason about the problem domain, not just decompose tickets.**

Read `.claude/agents/analyst.md` and `templates/claude-feature-team/.claude/agents/analyst.md`.

Upgrade the analyst in three areas — surgically add to the existing file, do not rewrite:

**1. Event Storming Lite (Alberto Brandolini)**
Before writing acceptance criteria, the analyst maps domain events:
- What domain events does this feature produce? (past tense: `UserLockedOut`, `PaymentProcessed`)
- What commands trigger them? (`LockAccount`, `ProcessPayment`)
- What aggregates own them? (`UserAccount`, `Invoice`)
- What read models / projections does the UI need?
Add `## Domain Events` section to the `analysis.md` output format.

**2. Bounded Context Mapping (Eric Evans)**
Identify which bounded context owns this feature. Flag any context crossings explicitly
as architectural concerns for the architect agent.
Add `## Bounded Context` section to the output.

**3. Fitness Function Proposals (Neal Ford)**
For every non-functional requirement (performance SLA, security constraint, availability
target), propose a measurable fitness function:
- What is the property?
- How would CI verify it? (specific tool, threshold, command)
- Who is responsible for wiring it?
Add `## Proposed Fitness Functions` section to the output.

Add these names to the analyst's governing philosophy:
- Alberto Brandolini (event storming)
- Dave Farley (acceptance tests verify *what*, never *how*)

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/analyst.md`.

Verify:
```bash
grep -c "Brandolini\|Domain Events\|Bounded Context\|Fitness Function" .claude/agents/analyst.md
# expect 4
```

---

## PROMPT 3 — Upgrade the Architect: Evolutionary Design Enforcer

> **Goal: make every architectural decision verifiable, reversible-classified, and stability-aware.**

Read `.claude/agents/architect.md` and `templates/claude-feature-team/.claude/agents/architect.md`.

The architect already has Clean Architecture, Fitness Functions, Simple Design, Refactoring,
and Saturday Framework rules. Add the following — surgically, do not rewrite:

**1. Strategic Domain Design (Eric Evans)**
Before structural decisions: identify the bounded context, map context crossings, choose
integration pattern (Shared Kernel / Customer-Supplier / Conformist / Anti-Corruption Layer).

**2. Stability Patterns (Michael Nygard)**
Every external dependency introduced must specify: Timeout, Circuit Breaker, Bulkhead,
Fail Fast, Idempotency. Use Sunday primitives — never custom retry loops.

**3. Integration Patterns (Gregor Hohpe)**
Name the integration pattern for every context crossing: direct call, event, shared database
(flag as anti-pattern), API gateway.

**4. Anti-Pattern Radar** (add as a named section):
Check adjacent code for: Distributed Monolith, Anemic Domain Model, God Object,
Shotgun Surgery, Leaky Abstraction, Premature Generalization.
Format as a table: Anti-Pattern | Symptom | Remedy.

**5. Reversibility Classification** (add as a named section):
Classify every structural decision: Cheap / Moderate / Expensive / Essentially Permanent.
Expensive or Permanent → RFC required. Essentially Permanent → human approval gate.

**6. Observability Architecture** (add as a named section):
For every feature: specify Spans (with attributes, no PII), Metrics (low cardinality),
Logs (structured fields, stable message strings), Alerts (new failure modes).

**7. RFC Output Mode**
Write a lightweight RFC for Expensive/Essentially Permanent decisions.
Save to `.claude/feature-workspace/rfc-[kebab-title].md`.

Add to the architect's `## Your Process`: read `DOMAIN_DICTIONARY.md` at step 3,
apply anti-pattern radar, classify reversibility, determine observability for every decision.

Add these sections to the `architecture-notes.md` output format:
`## Bounded Context`, `## Stability Design` (table), `## Observability Design`,
`## Anti-Pattern Check` (checklist), `## Reversibility` per decision.

Add to governing philosophy: Gregor Hohpe, Eric Evans (strategic design), Michael Nygard.

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/architect.md`.

Verify:
```bash
grep -c "Nygard\|Anti-Pattern Radar\|Reversibility\|Observability Architecture\|RFC Output Mode" .claude/agents/architect.md
# expect 5
```

---

## PROMPT 4 — Upgrade the Developer: Design-First Implementer

> **Goal: make implementation a design activity, not a typing activity.**

Read `.claude/agents/developer.md` and `templates/claude-feature-team/.claude/agents/developer.md`.

Add the following — surgically, do not rewrite:

**1. Interface-First Design**
Before writing any implementation, write the public interface/type signatures for every
new class or module. Add `## Interface Design` to the top of `implementation-notes.md`.
Rule: if you can't write the interface without looking at the implementation, the design
is not clear enough yet.

**2. Named Refactoring Log**
Every refactoring operation — new feature code and Boy Scout — must be logged by name
(Fowler catalog) with: operation name, file+line, before (the smell), after (the result).
This is not optional. The code-reviewer checks this log against the actual diff.

**3. Self-Review Checklist (before handoff)**
Add `## Self-Review Checklist` to `implementation-notes.md` output format:
- [ ] Every public method has an intention-revealing name (no `process`, `handle`, `manage`)
- [ ] No function exceeds 30 LOC
- [ ] Cyclomatic complexity < 7 on all new functions
- [ ] No primitive obsession: domain concepts wrapped in value objects
- [ ] No feature envy: methods use their own class's data more than another's
- [ ] No magic numbers or strings
- [ ] Dependency direction verified: inner layers never import outer layers

**4. "What Would Kent Do?" Simple Design Verification**
After the refactor pass, apply Beck's four rules in order and record pass/fail for each.
Add `## Simple Design Verification` to `implementation-notes.md` output format.

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/developer.md`.

Verify:
```bash
grep -c "Interface Design\|Named Refactoring Log\|Self-Review Checklist\|Simple Design Verification" .claude/agents/developer.md
# expect 4
```

---

## PROMPT 5 — Upgrade the Code Reviewer: Design Critic

> **Goal: review design, not just compliance.**

Read `.claude/agents/code-reviewer.md` and `templates/claude-feature-team/.claude/agents/code-reviewer.md`.

The code-reviewer was recently created with Design Narrative, Design Score, Security Surface,
and Test Design Review. Verify these are all present. If any are missing, add them now.

Then add these two missing pieces:

**1. Cohesion Analysis (Larry Constantine)**
Explicitly check for: Divergent Change, Shotgun Surgery, Feature Envy, Data Clumps,
Inappropriate Intimacy. Add to `## Evaluation Criteria` as `### Cohesion & Coupling`.

**2. Performance Smell Detection (shift-left)**
Add to `## Evaluation Criteria` as `### Reliability & Performance Smells`:
- N+1 queries (DB call inside a loop)
- Missing explicit timeout on any network/database call
- Missing idempotency on state-mutating calls
- Unbounded result sets (no pagination)
- Unnecessary synchrony (sequential calls that could be parallel)

Also add `## Performance Surface` to the `code-review-report.md` output format
alongside the existing `## Security Surface`.

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/code-reviewer.md`.

Verify:
```bash
grep -c "Cohesion\|Performance Surface\|N+1\|Divergent Change" .claude/agents/code-reviewer.md
# expect 4
```

---

## PROMPT 6 — Upgrade the QA Engineer: Test Strategist

> **Goal: make QA reason about test strategy, not just write scenarios.**

Read `.claude/agents/qa-engineer.md` and `templates/claude-feature-team/.claude/agents/qa-engineer.md`.

The QA agent was recently upgraded with Test Pyramid, Contract Testing Awareness,
Mutation Confidence, Performance SLA Verification, and Exploratory Charter.
Verify these are all present. If any are missing, add them now.

Then add this one remaining piece:

**Accessibility Verification**
After functional tests pass, QA must run a brief accessibility check on any UI changes:
- Interactive elements reachable by keyboard
- Form inputs have associated labels
- Color is not the only means of conveying state
- Use Playwright's accessibility snapshot: `await page.accessibility.snapshot()`
- Flag violations as `[A11Y]` findings in the QA report

Add `## Accessibility Check` to the `qa-report.md` output format.

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/qa-engineer.md`.

Verify:
```bash
grep -c "Pyramid\|Contract Testing\|Exploratory\|mutation\|Accessibility" .claude/agents/qa-engineer.md
# expect 5
```

---

## PROMPT 7 — Upgrade the Security Reviewer: Threat Modeler

> **Goal: verify the security reviewer is complete and wired to the code reviewer's surface map.**

Read `.claude/agents/security-reviewer.md`.

The security reviewer already has STRIDE, severity levels, Saturday/Sunday specifics, and
a complete output format. It is the strongest agent in the pipeline. Add one piece:

**Dependency Vulnerability Scan**
After reviewing the implementation, the security reviewer must check for known vulnerabilities
in any new dependencies introduced:
- Run `pnpm audit` (or `npm audit` / `pip-audit` / `go mod verify` depending on language)
- Any Critical or High CVEs in new dependencies block QA — same as a Critical finding
- Add `## Dependency Audit` to `security-report.md` output format

Also add a step to the process: "Read `code-review-report.md` → `## Security Surface`
section before starting STRIDE analysis. Use it to focus attention on the right files."

Verify:
```bash
grep -c "Dependency Audit\|code-review-report\|Security Surface" .claude/agents/security-reviewer.md
# expect 3
```

---

## PROMPT 8 — Create the `/design-review` Skill

> **Standalone design critique of any file, not tied to feature delivery.**

Read `.claude/skills/SKILL_TEMPLATE.md` and `ARCHITECTURE_RULES.md` in full.

Create `.claude/skills/design-review/SKILL.md` following the canonical template.

**Triggers:** "Review the design of [file]", "Is this well-designed?", "Check [file] for
design smells", "Would Fowler approve of this?"

**Context to load:** `ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`, the target file(s),
and any files the target imports from or that import it (one level of dependency graph).

**Process (in order):**
1. Write a Design Narrative (what is this thing trying to do?)
2. Map the dependency graph for the target (what does it depend on? what depends on it?)
3. Check all Fowler smells from ARCHITECTURE_RULES.md
4. Check Clean Architecture layer placement
5. Check Simple Design rules (Beck's 4 rules) in priority order
6. Check Sandi Metz constraints (class/method size, param count)
7. Identify the single most impactful refactoring to apply first
8. Produce a Design Report

**Output format** (`.claude/feature-workspace/design-report.md`):
```markdown
# Design Review: [ClassName or filename]

## Design Narrative
[2-3 sentences: what is this thing, what is its job, does it do one thing?]

## Dependency Graph
- Depends on: [list]
- Depended on by: [list]
- Layer: [which Clean Architecture layer]
- Direction correct: yes/no

## Design Smells Found
| Smell | Location | Severity | Fowler Operation |
|---|---|---|---|

## Simple Design Score
1. Passes tests: ✅/❌
2. Reveals intention: ✅/⚠️ [what to rename]
3. No duplication: ✅/⚠️ [what to extract]
4. Fewest elements: ✅/⚠️ [what to remove]

## Recommended Refactoring (Priority Order)
1. [Most impactful named operation] — [why, what it unlocks]
2. ...

## What's Working Well
[Always include — good design deserves recognition]
```

**Guardrails:** Read-only. Never modify source files. Never run tests. Surface findings only.

Copy to `templates/claude-feature-team/.claude/skills/design-review/SKILL.md`.

Verify:
```bash
ls .claude/skills/design-review/SKILL.md templates/claude-feature-team/.claude/skills/design-review/SKILL.md
```

---

## PROMPT 9 — Create the `/adr` Skill

> **Effortless, consistent Architecture Decision Records.**

Read `.claude/skills/SKILL_TEMPLATE.md`, `ARCHITECTURE_RULES.md`,
and any existing ADRs in `docs/adr/` if present.

Create `.claude/skills/adr/SKILL.md` following the canonical template.

**Triggers:** "Write an ADR for [decision]", "Document the decision to [X]",
"We decided to [X], record it", "Create an architecture decision record"

**Context to load:** All existing ADRs (for next number + tone),
`ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`,
`.claude/feature-workspace/architecture-notes.md` if it exists.

**Process (in order):**
1. Determine next ADR number (scan `docs/adr/` for existing files)
2. Ask: "What was decided?" (one sentence, active voice) — one question at a time
3. Ask: "What was the context that made this decision necessary?"
4. Ask: "What alternatives were considered and why were they rejected?"
5. Ask: "What are the consequences — easier, harder, changed?"
6. Ask: "Does this decision produce a fitness function? If so, how is it enforced?"
7. Draft the ADR and show it for approval
8. On approval, write to `docs/adr/ADR-[NNN]-[kebab-title].md`
9. Update `docs/adr/README.md` (index) — create it if it doesn't exist

**Output format:**
```markdown
# ADR-[NNN]: [Title — active voice: "Use X" not "Decision to use X"]

Date: YYYY-MM-DD
Status: Accepted
Deciders: [from git config user.name]
Technical Story: [feature or ticket that prompted this]

## Context
[What situation, constraint, or force prompted this decision]

## Decision
[What was decided — specific, concrete, active voice]

## Alternatives Considered
| Option | Why Rejected |
|---|---|

## Consequences
- **Easier**: [what this unlocks]
- **Harder**: [what this constrains]
- **Changed**: [what is different going forward]

## Fitness Function
[How this decision is enforced in CI — or "Judgment-only: [reason]"]
```

**Guardrails:**
- Never write an ADR without user approval of the draft
- Never number an ADR without checking the existing sequence
- ADR titles use active voice
- Never write the consequences as implementation details — write effects, not mechanisms

Copy to `templates/claude-feature-team/.claude/skills/adr/SKILL.md`.

Verify:
```bash
ls .claude/skills/adr/SKILL.md templates/claude-feature-team/.claude/skills/adr/SKILL.md
```

---

## PROMPT 10 — Create the `/five-whys` Skill

> **Structured root cause analysis for bugs, outages, and recurring problems.**

Read `.claude/skills/SKILL_TEMPLATE.md`.

Create `.claude/skills/five-whys/SKILL.md` following the canonical template.

**Triggers:** "Five whys on [problem]", "Root cause [issue]", "Why did [X] happen?",
"Post-mortem on [incident]", any mention of recurring bugs or "this keeps happening"

**Context to load:** Any relevant files mentioned (error logs, test failures, workspace
artifacts, git log for affected files).

**Facilitation rules (non-negotiable):**
- Never answer your own questions. Hypotheses stay private.
- Reflect back exactly what the user said before asking the next why.
- Ask one question at a time. Never stack questions.
- If the answer is vague ("because it broke"), ask for specificity.
- Stop when: answer is actionable / user says "I don't know" twice / five iterations complete.

**Process:**
1. State the symptom clearly: "The symptom is: [X]"
2. Ask "Why did [X] happen?" — wait, never suggest
3. Reflect: "So [X] happened because [their answer]. Why did [their answer] happen?"
4. Repeat until stop condition
5. Produce Root Cause Report

**Output format** (`.claude/feature-workspace/five-whys-[kebab-topic].md`):
```markdown
# Five Whys: [Problem Statement]

Date: [today]

## The Symptom
[What was observed — specific, observable]

## The Why Chain
1. Why did [symptom] happen? → [answer]
2. Why did [answer] happen? → [answer]
3. ...

## Root Cause
[The actionable thing at the bottom of the chain]

## Investigation Gaps
[Any "I don't know" answers needing further investigation]

## Recommended Action
- Immediate: [what to do right now]
- Preventive: [fitness function, test, or process change that prevents recurrence]
- Owner: [who should act]
```

**Guardrails:** Never suggest an answer during facilitation. Never skip the reflection step.
Never produce the report until the why chain is complete or a stop condition is hit.

**Standalone mode:** Fully conversational. No external tools needed.

Copy to `templates/claude-feature-team/.claude/skills/five-whys/SKILL.md`.

Verify:
```bash
ls .claude/skills/five-whys/SKILL.md templates/claude-feature-team/.claude/skills/five-whys/SKILL.md
```

---

## PROMPT 11 — Create the `/complexity-check` Skill

> **On-demand complexity analysis of any file or directory.**

Read `.claude/skills/SKILL_TEMPLATE.md` and `ARCHITECTURE_RULES.md`.

Create `.claude/skills/complexity-check/SKILL.md` following the canonical template.

**Triggers:** "Check complexity of [file]", "Is [file] too complex?", "Complexity report
on [directory]", "What's the most complex file in [package]?", any mention of "technical debt"
without a specific task attached.

**Tool selection by language:**
- TypeScript: `npx eslint --rule '{"complexity": ["error", 6]}' [file]`
- Python: `radon cc [file] -s` or `flake8 --max-complexity=6 [file]`
- Go: `gocyclo -over 6 [file]`
- Java: manual heuristics only — recommend Checkstyle/SonarQube to user

**Process (in order):**
1. Determine the language(s) of the target
2. Run the appropriate complexity tool
3. For each function exceeding complexity 6: name the smell, name the Fowler operation,
   estimate effort (trivial / one session / needs design discussion)
4. Rank findings by impact: highest complexity first, then by call graph depth
5. Produce Complexity Report

**Output format** (`.claude/feature-workspace/complexity-report.md`):
```markdown
# Complexity Report: [target]

Threshold: Cyclomatic complexity > 6 (ARCHITECTURE_RULES.md)

## Summary
- Files scanned: N
- Functions exceeding threshold: N
- Highest complexity: [function name] at [N]

## Findings (ranked by impact)

### [FunctionName] — complexity [N]
**File**: `path/to/file.ts` line X
**Smell**: [specific description]
**Refactoring**: [named Fowler operation]
**Effort**: trivial / one session / needs design discussion

## Refactoring Roadmap
1. [First because: quick win / unblocks other refactors / highest risk]

## What's Clean
[Functions within threshold — always acknowledge good work]
```

**Guardrails:** Read-only. Never modify files. Never run tests. Fall back to heuristic
analysis if CLI tools unavailable — still produce a report, note it's heuristic.

Copy to `templates/claude-feature-team/.claude/skills/complexity-check/SKILL.md`.

Verify:
```bash
ls .claude/skills/complexity-check/SKILL.md templates/claude-feature-team/.claude/skills/complexity-check/SKILL.md
```

---

## PROMPT 12 — Create the `/deliver-feature` Skill

> **The main pipeline orchestrator — kicks off the full agent sequence for a feature.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `CLAUDE.md` for the existing pipeline sequence and orchestration rules.
Read all agent files in `.claude/agents/` to understand their invocation order.

Create `.claude/skills/deliver-feature/SKILL.md` following the canonical template.

**Triggers:** "Deliver [feature]", "Implement [feature file]", "Build [feature]",
"Start delivery on [path/to/feature.md]", `/deliver-feature [feature-name]`

**Context to load:** The feature file (passed as argument or in `features/`),
`ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`, `CLAUDE.md`.

**Process (the full pipeline):**
1. **Read the feature file** — confirm it follows `features/TEMPLATE.md` structure.
   If not, stop and ask the user to run `/new-feature` first.
2. **Create the feature workspace**: `.claude/feature-workspace/` — clean any prior artifacts.
3. **Invoke analyst** → produces `analysis.md`. Pause: show summary to user. Continue.
4. **Invoke architect** (if analysis.md has Architectural Flags ≠ "None") →
   produces `architecture-notes.md`. Pause if RFC was written — human must acknowledge.
5. **Invoke developer** → produces `implementation-notes.md`.
6. **Invoke code-reviewer** → produces `code-review-report.md`.
   If verdict is CHANGES REQUESTED: send back to developer. Repeat until APPROVED.
7. **Invoke security-reviewer** (if security surface exists) → produces `security-report.md`.
   If Critical findings exist: block pipeline, alert user.
8. **Invoke qa-engineer** → produces `qa-report.md`. Tests must be green.
9. **Invoke tech-writer** → produces `docs-report.md`.
10. **Invoke devops-engineer** → produces `devops-report.md`.
11. **Pipeline complete**: show summary of all workspace artifacts. Ask: "Ship to Friday?"
    On confirmation ("ship" or "yes"): POST Cucumber JSON to Friday.

**Human checkpoints** (pipeline pauses and waits for explicit acknowledgment):
- After analyst: confirm scope before any code is written
- After architect RFC: confirm architectural direction before developer starts
- After code-review CHANGES REQUESTED loop: confirm all findings resolved
- After security Critical finding: explicit "fix confirmed" before QA starts
- Before shipping to Friday: explicit "ship" confirmation

**Output format:**
At the end of a successful delivery, write `.claude/feature-workspace/delivery-summary.md`:
```markdown
# Delivery Summary: [Feature Name]

## Pipeline Run
| Agent | Status | Key Output |
|---|---|---|
| analyst | ✅ | [N acceptance criteria, [N] architectural flags] |
| architect | ✅/⏭️ skipped | [N structural decisions, RFC: yes/no] |
| developer | ✅ | [N files created, N modified, N refactoring ops] |
| code-reviewer | ✅ | [Design score: C/Co/Cu/Cr — APPROVED] |
| security-reviewer | ✅/⏭️ skipped | [N findings, N critical fixed] |
| qa-engineer | ✅ | [N tests, N passed, SLAs verified: yes/no] |
| tech-writer | ✅ | [N docs updated] |
| devops-engineer | ✅ | [N CI changes, N env vars] |

## Friday
Status: Shipped | Pending | Skipped

## Artifacts
- `.claude/feature-workspace/analysis.md`
- `.claude/feature-workspace/architecture-notes.md`
- [all produced artifacts listed]
```

**Guardrails:**
- Never skip the analyst — it is always first
- Never let the developer start without analysis.md
- Never send CHANGES REQUESTED code to the security reviewer or QA
- Never ship to Friday without explicit "ship" or "yes" from the user
- Pipeline can be resumed from any checkpoint — check which workspace artifacts already exist

**Standalone mode:** All agents run locally. Friday POST is the only external call —
non-blocking if Friday isn't running.

Copy to `templates/claude-feature-team/.claude/skills/deliver-feature/SKILL.md`.

Verify:
```bash
ls .claude/skills/deliver-feature/SKILL.md templates/claude-feature-team/.claude/skills/deliver-feature/SKILL.md
```

---

## PROMPT 13 — Create the `/event-storm` Skill

> **Collaborative domain modeling using Event Storming Lite before any feature starts.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `.claude/agents/analyst.md` — the analyst references domain events but has no playbook.

Create `.claude/skills/event-storm/SKILL.md` following the canonical template.

**Triggers:** "Event storm [feature]", "Model the domain for [feature]",
"What are the domain events for [feature]?", "Let's do an event storm on [topic]",
any feature that involves complex business logic, state machines, or multiple aggregates

**When NOT to use:** Simple CRUD with no business logic, UI-only changes with no domain events,
infrastructure or configuration changes.

**Core concept (Alberto Brandolini):**
Event Storming is a collaborative modeling technique. You start from what *happened*
(domain events, past tense, orange sticky notes in the real workshop) and work backwards
to *why* it happened (commands, blue notes) and *who/what* owns it (aggregates, yellow notes).
The goal is to find the aggregate boundaries before writing a single class.

**The sticky note vocabulary:**
- **Domain Event** (past tense verb): `UserLockedOut`, `OrderShipped`, `PaymentDeclined`
- **Command** (imperative verb): `LockUser`, `ShipOrder`, `DeclinePayment`
- **Aggregate** (noun that owns events): `UserAccount`, `Order`, `Payment`
- **Read Model** (UI projection): what the frontend needs to display
- **Policy** (when event → command): "When `LoginFailed` 5 times, then `LockUser`"
- **External System**: anything outside your bounded context

**Process (conversational — one question at a time):**
1. State the feature scope: "We're modeling [feature]. Let's start from what can happen."
2. Ask: "What are the most important things that can *happen* in this feature?
   Give them as past-tense events — `[Noun][PastVerb]`."
   Wait. List them. Don't suggest yet.
3. For each event: "What *command* triggers [Event]? Who or what sends that command?"
4. "Which *aggregate* (the thing that owns the state) receives each command and emits each event?"
5. "Are there any *policies* — 'when [Event] happens, automatically trigger [Command]'?"
6. "What does the UI need to *read*? What projections or read models does this require?"
7. "Does this feature cross a *bounded context boundary*? If so, how do the contexts communicate?"
8. Produce the Event Storm Map

**Output format** (`.claude/feature-workspace/event-storm.md`):
```markdown
# Event Storm: [Feature Name]

## Domain Events
Past-tense — things that happened:
- `UserLockedOut` — triggered after 5 consecutive failed login attempts
- `LockoutExpired` — triggered 15 minutes after lockout begins

## Commands
What triggers each event:
| Command | Triggered by | Produces Event |
|---|---|---|
| `LockUserAccount` | Auth service (policy) | `UserLockedOut` |
| `ExpireLockout` | Scheduler (policy) | `LockoutExpired` |

## Aggregates
Who owns what:
| Aggregate | Owns Events | Key State |
|---|---|---|
| `UserAccount` | `UserLockedOut`, `LockoutExpired` | `failedAttempts`, `lockedUntil` |

## Policies
When → Then rules:
- When `LoginAttemptFailed` fires for the 5th time on the same account → trigger `LockUserAccount`
- When `UserLockedOut` fires → start 15-minute timer → trigger `ExpireLockout`

## Read Models
What the UI needs:
- `LoginPageState`: `isLocked: boolean`, `lockedUntilDisplay: string | null`

## Bounded Context
- Owning context: Identity & Access
- Context crossings: None / [describe crossing + integration pattern]

## Hotspots
Questions or ambiguities that need resolution before the analyst writes acceptance criteria:
- [Question]: [Why it matters]

## DOMAIN_DICTIONARY.md Additions
New terms that must be added before the developer starts:
- `UserLockedOut`: [definition]
- `LockoutExpired`: [definition]
```

**Guardrails:**
- Never suggest domain events before the user has given their first answer
- Never rename or combine the user's domain events without asking
- Hotspots are questions, never answers
- The output feeds directly into the analyst's `## Domain Events` section —
  note this at the bottom of the skill output

**Standalone mode:** Fully conversational. No external tools needed.

Copy to `templates/claude-feature-team/.claude/skills/event-storm/SKILL.md`.

Verify:
```bash
ls .claude/skills/event-storm/SKILL.md templates/claude-feature-team/.claude/skills/event-storm/SKILL.md
```

---

## PROMPT 14 — Create the `/db-migration` Skill

> **Safe zero-downtime database migrations using Expand/Contract.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `ARCHITECTURE_RULES.md` — specifically any section covering database migrations.
Read `.claude/agents/security-reviewer.md` — it flags destructive migrations.
Read `.claude/agents/devops-engineer.md` — it handles deployment.

Create `.claude/skills/db-migration/SKILL.md` following the canonical template.

**Triggers:** "Write a migration for [change]", "Add column [X] to [table]",
"Rename column [X] to [Y]", "Drop column/table [X]", "Change column type [X]",
any mention of schema changes, migrations, `ALTER TABLE`, `DROP`, `RENAME`

**Core concept (Expand/Contract pattern):**
Every breaking schema change must be split into at least two phases deployed independently.
This guarantees zero downtime and safe rollback at each step.

Phase 1 — **Expand**: add the new thing without removing the old thing. Both old and new
are valid simultaneously. Deploy this. Application code continues using the old schema.

Phase 2 — **Migrate**: deploy application code that writes to BOTH old and new.
Run a backfill script to copy existing data. Verify data integrity.

Phase 3 — **Contract**: remove the old thing. Application code uses only the new schema.
Only safe once all traffic is on Phase 2+ code and backfill is verified.

**Changes that require Expand/Contract:**
- Renaming a column (expand: add new column; contract: drop old column)
- Changing a column type (expand: add new typed column; contract: drop old)
- Adding a NOT NULL column without DEFAULT (expand: add nullable; contract: add constraint)
- Dropping a column (expand: stop writing to it; contract: drop it)
- Splitting or merging tables

**Changes that do NOT require Expand/Contract (safe single-phase):**
- Adding a nullable column with no constraints
- Adding a new table
- Adding an index (use CONCURRENTLY in PostgreSQL)
- Adding a column with a safe DEFAULT value

**Process:**
1. Ask: "What is the schema change?" (one sentence)
2. Classify: does it require Expand/Contract or is it single-phase safe?
3. If Expand/Contract: generate three migration files, clearly named Phase 1/2/3
4. For each phase: write the migration SQL, the rollback SQL, and the data integrity check
5. Write the backfill script for Phase 2 (if data migration required)
6. Specify the deployment sequence and the verification step at each phase
7. Check: does this migration need a feature flag to safely deploy Phase 2?

**Migration file naming convention:**
```
[timestamp]_phase1_expand_[description].sql
[timestamp]_phase2_migrate_[description].sql
[timestamp]_phase3_contract_[description].sql
```

**Output format** (`.claude/feature-workspace/db-migration-[description].md`):
```markdown
# DB Migration: [Description]

## Change Summary
[What is changing and why]

## Pattern
Expand/Contract (3 phases) | Single-phase safe

## Phase 1 — Expand
```sql
-- Forward
ALTER TABLE users ADD COLUMN locked_until TIMESTAMPTZ NULL;

-- Rollback
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
```
Deploy when: immediately
Verify: `SELECT COUNT(*) FROM users WHERE locked_until IS NOT NULL;` → expect 0

## Phase 2 — Migrate (if required)
```sql
-- Application writes to both old and new columns
-- Backfill existing data:
UPDATE users SET locked_until = lockout_expires_at WHERE lockout_expires_at IS NOT NULL;
```
Deploy when: Phase 1 is live and stable
Verify: `SELECT COUNT(*) FROM users WHERE lockout_expires_at != locked_until;` → expect 0

## Phase 3 — Contract
```sql
-- Forward
ALTER TABLE users DROP COLUMN lockout_expires_at;

-- Rollback: NOT POSSIBLE after data is removed — ensure Phase 2 is verified first
```
Deploy when: Phase 2 backfill verified, no traffic using old column
Verify: `\d users` → `lockout_expires_at` column absent

## Feature Flag Required
Yes / No — [reasoning]

## Risk Assessment
- Downtime risk: None / Low / [explain]
- Rollback possible at each phase: Phase 1 ✅ / Phase 2 ✅ / Phase 3 ⚠️ [explain]
- Estimated rows affected: [N] — backfill duration estimate: [time]

## Checklist
- [ ] Phase 1 SQL reviewed
- [ ] Phase 2 backfill script tested on a copy of production data
- [ ] Phase 3 only runs after Phase 2 verified
- [ ] No destructive operations in Phase 1 or Phase 2
- [ ] Rollback SQL tested for Phase 1 and Phase 2
```

**Guardrails:**
- `DROP COLUMN` or `DROP TABLE` NEVER appears in Phase 1 — only in Phase 3
- `NOT NULL` without `DEFAULT` NEVER appears in a single-phase migration
- `RENAME COLUMN` is NEVER used directly — always Expand/Contract instead
- Phase 3 requires explicit user confirmation: "confirm contract phase"
- Never write a migration that cannot be rolled back in Phase 1 or Phase 2

**Standalone mode:** Pure SQL generation. No external tools required.

Copy to `templates/claude-feature-team/.claude/skills/db-migration/SKILL.md`.

Verify:
```bash
ls .claude/skills/db-migration/SKILL.md templates/claude-feature-team/.claude/skills/db-migration/SKILL.md
```

---

## PROMPT 15 — Create the `/performance-profile` Skill

> **Diagnose why something is slow, not just verify that it is.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `.claude/agents/qa-engineer.md` — QA verifies SLAs, but when one fails, this skill
takes over.

Create `.claude/skills/performance-profile/SKILL.md` following the canonical template.

**Triggers:** "Why is [X] slow?", "Profile [feature/endpoint/page]", "The SLA is failing
for [X]", "Performance is degrading on [Y]", any timing assertion failure in QA

**When NOT to use:** Premature optimization before a test fails its SLA. Always measure first.

**Core principle (Donald Knuth via every performance engineer):**
Measure before you optimize. A guess about where slowness lives is almost always wrong.
The skill's job is to get to a measurement before recommending anything.

**Tool selection by layer:**

| Layer | Tool | What it finds |
|---|---|---|
| Browser/E2E | Playwright trace viewer (`--trace on`) | Waterfall, long tasks, layout thrash |
| Node.js | `clinic.js flame` / `0x` | CPU flame chart, event loop blocking |
| Python | `py-spy top` / `cProfile` | Hot functions, blocking calls |
| Database | Query logging + `EXPLAIN ANALYZE` | N+1s, missing indexes, sequential scans |
| Network | Chrome DevTools / Playwright HAR | Slow requests, payload size, no caching |

**Process (guided diagnostic — do not skip steps):**
1. Establish the baseline: "What is the current measured time and what is the SLA?"
2. Identify the layer: browser, server, database, or network?
   Ask the user: "Where does the time appear to be spent?" If unknown, start at browser.
3. Generate the profiling command for the identified layer
4. Interpret the output: what is the hottest path?
5. Apply the diagnostic checklist for that layer (see below)
6. Produce Performance Profile Report with ranked hypotheses and the one experiment to run first

**Diagnostic checklists:**

Database layer:
- [ ] Run `EXPLAIN ANALYZE` on the slow query — look for Seq Scan on large tables
- [ ] Check for N+1: does the code execute a query inside a loop?
- [ ] Check for missing index on the WHERE / JOIN / ORDER BY columns
- [ ] Check for unbounded result set (no LIMIT)

Server layer (Node.js):
- [ ] Is there a synchronous operation blocking the event loop? (use clinic.js)
- [ ] Is there an await inside a loop that could be parallelized with `Promise.all`?
- [ ] Is there an external HTTP call without a timeout that could be hanging?
- [ ] Is there a large JSON serialization/deserialization in the hot path?

Browser layer:
- [ ] Open Playwright trace: `npx playwright show-trace trace.zip`
- [ ] Look for Long Tasks (> 50ms) in the main thread
- [ ] Check for layout thrash: read/write/read/write DOM in a loop
- [ ] Check for render-blocking resources in the waterfall
- [ ] Check for uncompressed or uncached static assets

**Output format** (`.claude/feature-workspace/performance-profile-[description].md`):
```markdown
# Performance Profile: [Feature/Endpoint/Page]

## Baseline
- Current measured time: [Xms]
- SLA target: [Yms]
- Gap: [Xms over target]

## Layer Identified
Browser | Server | Database | Network

## Profiling Commands Run
```bash
[exact commands used]
```

## Hottest Path
[What the profiler identified as the primary bottleneck]

## Diagnostic Checklist Results
- [ ] [Check]: [result]

## Ranked Hypotheses
1. [Most likely cause]: [evidence] — [experiment to confirm]
2. [Second hypothesis]: [evidence] — [experiment to confirm]

## Recommended First Experiment
[One specific change to make, with expected impact and how to measure it]

## Expected Outcome
If hypothesis 1 is correct, making [change] should reduce time from [X]ms to [Y]ms.
Verify with: [exact command or test to rerun]
```

**Guardrails:**
- Never recommend an optimization before running a profiler
- Only recommend one experiment at a time — concurrent changes make causation impossible to establish
- Never recommend premature abstraction (caching, CDN, horizontal scaling) before ruling out
  algorithmic fixes (N+1, missing index, blocking loop)

**Standalone mode:** Provides profiling command generation and diagnostic checklists
even when profiling tools aren't installed. Always produces a structured hypothesis list.

Copy to `templates/claude-feature-team/.claude/skills/performance-profile/SKILL.md`.

Verify:
```bash
ls .claude/skills/performance-profile/SKILL.md templates/claude-feature-team/.claude/skills/performance-profile/SKILL.md
```

---

## PROMPT 16 — Create the `/dependency-update` Skill

> **Safe, structured process for updating packages in a pnpm monorepo.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `ARCHITECTURE_RULES.md`.
Read `pnpm-workspace.yaml` to understand the workspace structure.

Create `.claude/skills/dependency-update/SKILL.md` following the canonical template.

**Triggers:** "Update dependencies", "Update [package-name]", "Check for outdated packages",
"Security advisories in dependencies", "Bump [package] to [version]",
any mention of `pnpm update`, `npm audit`, or `dependabot`

**When NOT to use:** Major version updates with breaking changes should go through the
full feature delivery pipeline — they need analysis, architecture review, and QA.
This skill is for patch and minor updates.

**Risk classification:**
- **Patch** (x.y.Z): bug fixes only per semver — low risk, can batch
- **Minor** (x.Y.z): new features, backwards compatible per semver — medium risk, update one at a time
- **Major** (X.y.z): breaking changes guaranteed — use full delivery pipeline, not this skill

**Process (in order):**
1. Run security audit first: `pnpm audit` — fix any Critical or High CVEs before anything else
2. Check for outdated packages: `pnpm outdated`
3. For each package to update, check the changelog before updating:
   - Find the package's CHANGELOG.md or GitHub releases
   - Read entries between current version and target version
   - Look for: breaking changes hidden in minor/patch, deprecation notices,
     peer dependency requirement changes, behavior changes in edge cases
4. Update packages one at a time (never batch major updates):
   `pnpm update [package]@[version] --filter [workspace]`
5. After each update: run `tsc --noEmit`, `pnpm build`, `pnpm -r run test`
6. If tests fail: isolate whether it's the update or a pre-existing issue
7. Regenerate lock file: `pnpm install`
8. Produce Dependency Update Report

**Peer dependency conflicts in pnpm monorepos:**
When a peer dependency conflict arises:
- Check `pnpm-workspace.yaml` for the affected packages
- Use `pnpm why [package]` to trace the dependency chain
- Consider `pnpm.overrides` in root `package.json` as a last resort — document why
- Never use `--legacy-peer-deps` silently — always flag it and explain the risk

**Output format** (`.claude/feature-workspace/dependency-update-report.md`):
```markdown
# Dependency Update Report

Date: [today]

## Security Audit
- `pnpm audit` result: [N critical, N high, N medium]
- CVEs fixed: [list with CVE numbers]

## Updates Applied
| Package | From | To | Type | Changelog reviewed | Test result |
|---|---|---|---|---|---|
| `express` | 4.18.1 | 4.18.3 | patch | ✅ | ✅ pass |

## Changelog Notes
For any non-trivial update, note what changed:
- `[package] [version]`: [what changed that's relevant to this codebase]

## Peer Dependency Conflicts
- [Conflict]: [Resolution applied] / "None"

## Lock File
- Regenerated: ✅
- Committed: yes / pending

## Build & Test Results
- `tsc --noEmit`: ✅ / ❌ [errors]
- `pnpm build`: ✅ / ❌ [errors]
- Test suite: ✅ / ❌ [N failures — all fixed / [N] outstanding]

## Skipped Updates
| Package | Current | Latest | Reason skipped |
|---|---|---|---|
| `react` | 18.2 | 19.0 | Major — needs full delivery pipeline |

## Recommended Next Steps
- [Packages that need major updates — with ADR suggestion]
- [Packages with known issues to watch]
```

**Guardrails:**
- Always run `pnpm audit` first — do not skip security check
- Never batch major version updates — each one needs individual analysis
- Always read the changelog before updating a package, however briefly
- Never commit a broken lock file
- Never use `--force` to resolve peer dependency conflicts silently

**Standalone mode:** Works without external systems. All commands are local pnpm operations.

Copy to `templates/claude-feature-team/.claude/skills/dependency-update/SKILL.md`.

Verify:
```bash
ls .claude/skills/dependency-update/SKILL.md templates/claude-feature-team/.claude/skills/dependency-update/SKILL.md
```

---

## PROMPT 17 — Create the `/on-call` Skill

> **Active incident response — for when something is broken right now.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `.claude/agents/devops-engineer.md` — it handles prevention; this handles response.
Read `.claude/skills/five-whys/SKILL.md` — that handles RCA after the incident; this is during.

Create `.claude/skills/on-call/SKILL.md` following the canonical template.

**Triggers:** "Something is broken in production", "The site is down", "Alert firing for [X]",
"Incident on [service/feature]", "Rollback [X]", "[Fitness function] is failing in CI",
any language of urgency about a live system

**Critical framing:**
This skill has two modes:
- **Active Incident**: system is broken RIGHT NOW — speed matters, minimize cognitive load
- **Post-Incident**: system is stable — shift to Five Whys for root cause, this skill for
  timeline reconstruction and communication

During an active incident, do not over-engineer the response. Three questions only:
1. What is the blast radius? (who/what is affected)
2. Can we roll back?
3. Who needs to know?

**Active Incident Process (minimize steps — time matters):**
1. Establish the blast radius immediately:
   "What is the observable symptom? How many users/requests are affected?"
2. Check for a recent deploy: `git log --oneline -10`
   If yes → rollback decision: "Is reverting the last deploy safe and faster than a hotfix?"
3. Rollback or hotfix decision tree:
   - Rollback: `git revert HEAD` or redeploy previous image — safe if no DB migration in deploy
   - Hotfix: only if rollback is not safe (migration already ran) or fix is trivially small
4. Communicate: draft the incident message (who, what, ETA)
5. Once stable: capture timeline for the Five Whys session

**Communication templates (copy-ready):**

Initial alert:
```
🔴 INCIDENT — [timestamp]
Symptom: [what is broken]
Blast radius: [who/what affected]
Status: Investigating
Next update in 15 minutes.
```

Update:
```
🟡 UPDATE — [timestamp]
Root cause identified: [one sentence]
Action taken: [rollback / hotfix deploying]
ETA to resolution: [time]
```

Resolution:
```
🟢 RESOLVED — [timestamp]
Duration: [X minutes]
Cause: [one sentence]
Fix: [what was done]
Follow-up: Five Whys scheduled for [date]
```

**Rollback decision criteria:**
- Safe to rollback: no DB migrations in the deploy, stateless change
- NOT safe to rollback without data risk: DB migration already ran
- NOT safe to rollback: data has been written in the new schema that the old code can't read

**CI fitness function failure (not production down, but pipeline broken):**
1. Identify which fitness function failed and in which commit
2. Is this a real violation or a false positive?
3. If real: the offending commit must be reverted or fixed forward — never lower the threshold
4. If false positive: update the fitness function definition + add a comment explaining why

**Output format** (`.claude/feature-workspace/incident-[description].md`):
```markdown
# Incident Report: [Description]

## Timeline
| Time | Event |
|---|---|
| HH:MM | Alert fired / symptom reported |
| HH:MM | Blast radius established |
| HH:MM | Rollback/hotfix decision made |
| HH:MM | Fix deployed |
| HH:MM | Incident resolved |

## Blast Radius
- Users affected: [N / all / subset — be specific]
- Services affected: [list]
- Data integrity: [affected / not affected]

## Root Cause (preliminary)
[One sentence — full RCA via Five Whys to follow]

## Fix Applied
[Rollback to commit X / Hotfix: description]

## Rollback Safety
Safe ✅ / Unsafe ⚠️ [reason — DB migration state]

## Communication Sent
- [Channel]: [timestamp] — [which template used]

## Follow-Up Required
- [ ] Five Whys session scheduled: [date]
- [ ] Fitness function added to prevent recurrence: [yes / not yet identified]
- [ ] Postmortem doc: [link / pending]
```

**Guardrails:**
- During active incident: never ask more than one question at a time
- Never recommend a complex fix when a rollback is safe
- Never lower a fitness function threshold to make a failing check green
- Rollback + Five Whys is almost always better than hotfix under pressure

**Standalone mode:** Pure reasoning and templates. No external tools required during active incident.
All commands are standard git/pnpm/kubectl operations.

Copy to `templates/claude-feature-team/.claude/skills/on-call/SKILL.md`.

Verify:
```bash
ls .claude/skills/on-call/SKILL.md templates/claude-feature-team/.claude/skills/on-call/SKILL.md
```

---

## PROMPT 18 — Create the `/openapi` Skill

> **Deliberate API contract design before a line of implementation code is written.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `.claude/agents/analyst.md` — it identifies API changes but doesn't design them.
Read `.claude/agents/architect.md` — it places API clients in layers but doesn't design contracts.
Read the Sunday framework documentation in the codebase (`packages/sunday-*` or similar).

Create `.claude/skills/openapi/SKILL.md` following the canonical template.

**Triggers:** "Design the API for [feature]", "Write an OpenAPI spec for [endpoint]",
"What should the API look like for [X]?", "Define the contract for [endpoint]",
any feature where `analysis.md` has API Changes ≠ "None"

**When NOT to use:** Internal method signatures, test helper APIs, or any interface
that isn't crossed over HTTP.

**Design principles (REST maturity + pragmatism):**

Resource naming:
- Nouns, not verbs: `/users/123/sessions` not `/login` or `/createSession`
- Plural nouns for collections: `/orders` not `/order`
- Kebab-case for multi-word: `/rate-limit-policies` not `/rateLimitPolicies`
- Sub-resources for ownership: `/users/123/orders` not `/orders?userId=123`

HTTP verbs:
- `GET`: safe and idempotent — no side effects
- `POST`: create or non-idempotent action — requires idempotency key if retryable
- `PUT`: full replacement, idempotent
- `PATCH`: partial update — use JSON Merge Patch (RFC 7396) not custom shapes
- `DELETE`: idempotent — return 204 or 404, never 500 if already deleted

Status codes (use precisely — never overload 200 for errors):
- 200: success with body | 201: created | 204: success no body
- 400: client error (malformed request) | 401: not authenticated | 403: not authorized
- 404: not found | 409: conflict (already exists) | 422: validation error
- 429: rate limited (must include `Retry-After` header) | 500: server error (never leak internals)

Error response shape (consistent across ALL endpoints):
```json
{
  "error": {
    "code": "ACCOUNT_LOCKED",
    "message": "This account has been locked due to too many failed login attempts.",
    "retryAfter": "2024-01-15T14:30:00Z"
  }
}
```
Never return: raw exception messages, stack traces, database error strings, internal paths.

Pagination: always cursor-based for large collections:
```json
{
  "data": [...],
  "pagination": {
    "cursor": "eyJpZCI6MTAwfQ==",
    "hasMore": true,
    "pageSize": 20
  }
}
```

Versioning: URL path versioning (`/api/v1/`) for major breaking changes only.
Additive changes (new optional fields) are non-breaking — never version for these.

**Process (guided, one question at a time):**
1. "What resource does this endpoint operate on?"
2. "What HTTP verb best describes the operation? Walk through the options together if unsure."
3. "What does a successful response look like? What fields does the consumer actually need?"
4. "What can go wrong? List the failure cases — each needs a distinct status code."
5. "Does this operation have side effects? If so, does it need an idempotency key?"
6. "What authentication/authorization is required?"
7. "Are there any performance constraints? (pagination, rate limits, timeouts)"
8. Draft the OpenAPI spec section and show it for approval
9. On approval, write to `docs/api/[resource].openapi.yaml` (or append to existing spec)

**Output format:**
```yaml
# .claude/feature-workspace/openapi-[resource].yaml
openapi: "3.1.0"
info:
  title: "[Resource] API Contract"
  version: "draft"

paths:
  /api/v1/users/{userId}/sessions:
    post:
      summary: Create a new session (login)
      operationId: createUserSession
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateSessionRequest'
      responses:
        '201':
          description: Session created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Session'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '429':
          $ref: '#/components/responses/RateLimited'

components:
  schemas:
    CreateSessionRequest:
      type: object
      required: [email, password]
      properties:
        email:
          type: string
          format: email
        password:
          type: string
          format: password
          writeOnly: true
    Session:
      type: object
      required: [token, expiresAt]
      properties:
        token:
          type: string
        expiresAt:
          type: string
          format: date-time
  responses:
    RateLimited:
      description: Too many requests
      headers:
        Retry-After:
          schema:
            type: string
            format: date-time
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
```

Also produce `.claude/feature-workspace/api-design-notes.md`:
```markdown
# API Design Notes: [Resource]

## Design Decisions
| Decision | Rationale |
|---|---|
| POST /sessions not /login | Resource-oriented — session is the resource |
| 429 with Retry-After | Communicates lockout duration without revealing account state |

## Security Surface
- [Auth required: yes/no + mechanism]
- [User enumeration risk: addressed how]
- [Rate limiting: yes/no + threshold]

## Breaking Change Assessment
- Additive only: yes/no
- Breaking changes: [list — requires version bump]

## Sunday Framework Mapping
- Client class: `[ClassName]` extends `BaseApiClient`
- Zod schema: `[SchemaName]` validates response shape
- Auth provider: `[Provider]`
```

**Guardrails:**
- Never return error messages that reveal whether an account exists (user enumeration)
- Every collection endpoint must have pagination
- Every mutation must specify idempotency strategy
- OpenAPI spec must be shown to the user for approval before being written to disk
- Never use 200 for error responses

**Standalone mode:** Pure design reasoning. No external tools required.
Produces YAML spec locally — no API gateway, linter, or mock server required.

Copy to `templates/claude-feature-team/.claude/skills/openapi/SKILL.md`.

Verify:
```bash
ls .claude/skills/openapi/SKILL.md templates/claude-feature-team/.claude/skills/openapi/SKILL.md
```

---

## PROMPT 19 — Create Global Rules Files

> **Cross-cutting rules that apply to ALL agents, regardless of role.**

Read all files in `.claude/agents/`, `ARCHITECTURE_RULES.md`, and `CLAUDE.md` at repo root.

Create `.claude/rules/` directory. Create three rule files:

**`.claude/rules/design-principles.md`**
Cross-cutting design principles every agent must honor:
- Beck's four Simple Design rules in priority order — with TypeScript examples
- The 10 Fowler refactoring operations — with the smell that triggers each one
- The anti-pattern radar (Distributed Monolith, Anemic Domain, God Object,
  Shotgun Surgery, Leaky Abstraction, Premature Generalization)
- Sandi Metz hard limits — classes ≤ 100 lines, methods ≤ 5 lines (10 ceiling),
  max 4 params, no more than one dot per line except fluent chains
- The Boy Scout Rule as an active instruction with concrete examples
- Naming standards: intention-revealing names, no abbreviations, no generic terms
  (`process`, `handle`, `manage`, `data`, `info`), boolean prefix rules
  (`is`, `has`, `can`, `should`)
- Ubiquitous language: all names must match `DOMAIN_DICTIONARY.md`

**`.claude/rules/architecture-guardrails.md`**
Hard constraints that cannot be overridden by any agent or user instruction:
- Dependency direction: inner layers never import outer layers — with examples per language
- No destructive migrations: Expand/Contract pattern is non-negotiable
- No secrets hardcoded anywhere: `.env` placeholders only
- No raw `any` in TypeScript: use `unknown` when genuinely unknown
- No custom retry loops: use `CircuitBreaker` or `ExponentialBackoffStrategy`
- No N+1 queries: eager loading required
- No implicit timeouts: every network call gets an explicit timeout
- No unbounded result sets: pagination required on all collection endpoints
- Every architectural decision produces a fitness function or is explicitly flagged
  as "judgment-only" with a documented reason
- No OTel instrumentation in domain/page logic: traces emit from adapter/interceptor layer only

**`.claude/rules/approval-gates.md`**
Irreversible actions require explicit human approval. Any edit resets the gate.
Format each gate as:
```
Action: [what it does]
Irreversible because: [why it can't be undone]
Gate: user must say "[exact word(s)]"
Reset condition: any edit to the pending artifact resets the gate
```

Gates to define:
- Shipping to Friday (POST Cucumber JSON)
- Creating a git commit
- Running database migrations (any phase of Expand/Contract)
- Posting to external APIs
- Writing files outside `.claude/feature-workspace/`
- Wiring a new fitness function into CI
- Contracting phase of a DB migration (Phase 3 — data loss risk)
- Deploying to any environment

After creating the three rule files, update the `## Your Process` section header of every
agent file in `.claude/agents/` to include:
```
Before beginning: read `.claude/rules/design-principles.md`,
`.claude/rules/architecture-guardrails.md`, and `.claude/rules/approval-gates.md`.
```

Copy all three rule files to `templates/claude-feature-team/.claude/rules/`.

Verify:
```bash
ls .claude/rules/
grep -l "rules/design-principles" .claude/agents/analyst.md
# both should return output
```

---

## PROMPT 20 — Update scaffold-team.sh and install

> **Make sure everything gets deployed when scaffolding a new project.**

Read `scaffold-team.sh` and `install`.

Update `scaffold-team.sh` to copy ALL of the following to the target project:

Agents (8):
```
analyst, architect, code-reviewer, developer, devops-engineer,
qa-engineer, security-reviewer, tech-writer
```

Skills (14):
```
new-feature, deliver-feature, review-pr, retrospective,
design-review, adr, five-whys, complexity-check,
contract-testing, mutation-check,
event-storm, db-migration, performance-profile,
dependency-update, on-call, openapi
```

Rules (3):
```
design-principles, architecture-guardrails, approval-gates
```

Templates (2):
```
features/TEMPLATE.md, DOMAIN_DICTIONARY.template.md
```

Also copy:
- `SKILL_TEMPLATE.md` to `.claude/skills/SKILL_TEMPLATE.md` in target

Add a verification step that prints a count report:
```bash
echo "✅ Agents deployed:     $(ls [target]/.claude/agents/ | wc -l)  (expected: 8)"
echo "✅ Skills deployed:     $(ls [target]/.claude/skills/ | wc -l)  (expected: 14)"
echo "✅ Rules deployed:      $(ls [target]/.claude/rules/ | wc -l)   (expected: 3)"
echo "✅ Templates:           2"
```

Update `install` to symlink `.claude/rules` in addition to `.claude/agents` and `.claude/skills`.
Update the `FILES_TO_LINK` array:
```bash
".claude/agents"
".claude/skills"
".claude/rules"
```

Do a dry-run check: verify all source paths referenced in the scripts actually exist.
Report any missing paths before finishing.

Verify:
```bash
bash scaffold-team.sh --dry-run 2>&1 | grep -E "missing|error|✅"
```

---

## Final Verification

Run from repo root after all 20 prompts are complete:

```bash
echo "=== Agents (expect 8) ===" && ls .claude/agents/ | wc -l && ls .claude/agents/
echo ""
echo "=== Skills (expect 14) ===" && ls .claude/skills/ | wc -l && ls .claude/skills/
echo ""
echo "=== Rules (expect 3) ===" && ls .claude/rules/ && echo ""
echo "=== Template mirror agents ===" && ls templates/claude-feature-team/.claude/agents/ | wc -l
echo "=== Template mirror skills ===" && ls templates/claude-feature-team/.claude/skills/ | wc -l
echo "=== Template mirror rules ===" && ls templates/claude-feature-team/.claude/rules/ | wc -l
echo ""
echo "=== Key content checks ==="
grep -l "Brandolini\|Domain Events" .claude/agents/analyst.md && echo "analyst: domain events ✅"
grep -l "Nygard\|Anti-Pattern Radar" .claude/agents/architect.md && echo "architect: stability ✅"
grep -l "Interface Design\|Self-Review" .claude/agents/developer.md && echo "developer: design-first ✅"
grep -l "Design Narrative\|Design Score" .claude/agents/code-reviewer.md && echo "reviewer: design score ✅"
grep -l "Pyramid\|Exploratory" .claude/agents/qa-engineer.md && echo "qa: test strategy ✅"
grep -l "Dependency Audit" .claude/agents/security-reviewer.md && echo "security: dep audit ✅"
ls .claude/skills/event-storm/SKILL.md && echo "event-storm skill ✅"
ls .claude/skills/db-migration/SKILL.md && echo "db-migration skill ✅"
ls .claude/skills/performance-profile/SKILL.md && echo "performance-profile skill ✅"
ls .claude/skills/dependency-update/SKILL.md && echo "dependency-update skill ✅"
ls .claude/skills/on-call/SKILL.md && echo "on-call skill ✅"
ls .claude/skills/openapi/SKILL.md && echo "openapi skill ✅"
ls .claude/rules/design-principles.md && echo "rules: design-principles ✅"
ls .claude/rules/architecture-guardrails.md && echo "rules: architecture-guardrails ✅"
ls .claude/rules/approval-gates.md && echo "rules: approval-gates ✅"
```

All 15 checks should print ✅. Total: 8 agents · 14 skills · 3 rules · mirrored in templates.
