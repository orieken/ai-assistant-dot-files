# Build-Out Prompts: Skills, Rules & Agent Upgrades

These prompts are designed to be handed to Claude Code one at a time, from the root of your
`ai-assistant-dot-files` repo. Each prompt is self-contained. Run them in order — later prompts
assume earlier ones are complete.

The goal: make every agent in your team a world-class software craftsman who cares deeply about
great design and architecture — not just one that knows the rules, but one that *reasons* about
design the way Fowler, Beck, Uncle Bob, and Evans would.

The pattern used throughout is from Brad Feld's CompanyOS: skills as domain knowledge playbooks,
not orchestration engines. Every skill must work even if external systems are unavailable.

---

## PROMPT 1 — Create the Skill Template Standard

> **Paste this into Claude Code first. Everything else builds on it.**

Read `ARCHITECTURE_RULES.md` and all existing files in `.claude/agents/` and
`.claude/skills/`.

Create a new file: `.claude/skills/SKILL_TEMPLATE.md`

This is the canonical template every skill in this repo must follow. Model it on Brad Feld's
CompanyOS seven-section structure, adapted for software development:

```
---
name: skill-name
description: One sentence. What it does and when Claude should auto-trigger it.
triggers:
  keywords: [list, of, trigger, words]
  intentPatterns: ["regex or natural language pattern"]
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
The intelligence stays in the markdown. The API calls are just plumbing.
```

After creating the template, read every existing skill in `.claude/skills/` and verify each
one follows this structure. For any that don't, refactor them to match. Report what changed.

---

## PROMPT 2 — Upgrade the Analyst: Domain Modeler

> **Goal: make the analyst reason about the problem domain, not just decompose tickets.**

Read:
- `.claude/agents/analyst.md`
- `ARCHITECTURE_RULES.md` sections II (DDD) and VI (Simple Design)
- `templates/claude-feature-team/.claude/agents/analyst.md`

The analyst must be upgraded from "ticket decomposer" to "domain modeler." The upgrade adds
three new capabilities that must be woven into the agent's existing process and output format:

**1. Event Storming Lite**
Before writing acceptance criteria, the analyst must map the feature's domain events:
- What domain events does this feature produce? (e.g. `UserLockedOut`, `PaymentProcessed`)
- What commands trigger them? (e.g. `LockAccount`, `ProcessPayment`)
- What aggregates own them? (e.g. `UserAccount`, `Invoice`)
- What read models / projections does the UI need?
Add a `## Domain Events` section to the analyst's `analysis.md` output format.

**2. Bounded Context Mapping**
The analyst must identify which bounded context owns this feature and flag any context
crossings. If the feature crosses a context boundary (e.g. Billing reaching into Identity),
flag it explicitly as an architectural concern for the architect.
Add a `## Bounded Context` section to the output.

**3. Fitness Function Proposal**
For every non-functional requirement (performance SLA, security constraint, availability
target), the analyst must propose a measurable fitness function:
- What is the property?
- How would CI verify it? (specific tool, threshold, command)
- Who is responsible for wiring it? (architect or devops)
Add a `## Proposed Fitness Functions` section to the output.

Also upgrade the analyst's governing philosophy section to explicitly name:
- Eric Evans (domain modeling, ubiquitous language, bounded contexts)
- Alberto Brandolini (event storming)
- Dan North (BDD as communication, not test structure)
- Dave Farley (acceptance tests verify *what*, never *how*)

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/analyst.md`.

---

## PROMPT 3 — Upgrade the Architect: Evolutionary Design Enforcer

> **Goal: make the architect produce decisions that are verifiable, not just reasonable.**

Read:
- `.claude/agents/architect.md`
- `ARCHITECTURE_RULES.md` sections I (Clean Architecture), VIII (Fitness Functions),
  and X (Database Craftsmanship)
- `templates/claude-feature-team/.claude/agents/architect.md`

The architect must be upgraded in three areas:

**1. Fitness Function First**
Every structural decision the architect makes must produce a fitness function before it
produces implementation guidance. The format is:
```
Decision: [what was decided]
Fitness Function: [property that must remain true]
Enforcement: [exact eslint rule / dependency-cruiser config / tsc flag / CI command]
Responsibility: devops-engineer wires it, architect defines it
```
If a decision cannot produce a fitness function, the architect must justify why and flag
it as "judgment-only" — not enforced, relies on review.

**2. RFC (Request for Comments) Output Mode**
For any decision that touches a layer boundary, introduces a new abstraction, or affects
more than one package, the architect must write a lightweight RFC in addition to
`architecture-notes.md`. Format:
```markdown
# RFC: [Title]
Status: Proposed
Date: [today]

## Problem
## Proposed Solution
## Alternatives Considered
## Trade-offs
## Open Questions (requires human input before developer proceeds)
```
Save to `.claude/feature-workspace/rfc-[kebab-title].md`

**3. Anti-Pattern Radar**
The architect must explicitly check for these architecture anti-patterns before signing off:
- Distributed monolith: microservices that are still tightly coupled via synchronous calls
- Anemic domain model: entities with no behavior, all logic in services
- God object: a class that knows too much about too many things
- Shotgun surgery: a single change requiring edits to many unrelated classes
- Leaky abstraction: a "generic" layer that forces callers to know about implementation details

For each anti-pattern found in adjacent code (not just the new feature), raise it as a
refactoring opportunity in `architecture-notes.md` with the Fowler operation to apply.

Also add to the architect's governing philosophy:
- Gregor Hohpe (enterprise integration patterns, messaging between bounded contexts)
- Eric Evans (strategic design, context maps, anti-corruption layers)
- Michael Nygard (release it! — stability patterns: circuit breakers, bulkheads, timeouts)

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/architect.md`.

---

## PROMPT 4 — Upgrade the Developer: Design-First Implementer

> **Goal: make the developer treat implementation as a design activity, not a typing activity.**

Read:
- `.claude/agents/developer.md`
- `ARCHITECTURE_RULES.md` sections VI (Simple Design), VII (Refactoring Catalog)
- `templates/claude-feature-team/.claude/agents/developer.md`

The developer must be upgraded in four areas:

**1. Interface-First Design**
Before writing any implementation, the developer must write the public interface/type
signatures for every new class or module. These go in a `## Interface Design` section at
the top of `implementation-notes.md`. The rule: if you can't write the interface without
looking at the implementation, the design is not clear enough yet.

**2. Named Refactoring Log**
Every refactoring operation applied — whether to new feature code or as a Boy Scout cleanup
— must be logged by name (Fowler catalog) with:
- Operation name (e.g. "Extract Function")
- File and approximate line
- Before: what the smell was
- After: what it became
This is not optional. The code-reviewer checks this log against the actual diff.

**3. Design Smell Checklist (self-review before handoff)**
Before writing `implementation-notes.md`, the developer must run through this checklist
and record the result for each item:
- [ ] Every public method has an intention-revealing name (no `process`, `handle`, `manage`)
- [ ] No function exceeds 30 LOC (hard limit from ARCHITECTURE_RULES.md)
- [ ] Cyclomatic complexity < 7 on all new functions
- [ ] No primitive obsession: domain concepts wrapped in value objects, not raw strings/ints
- [ ] No feature envy: methods use their own class's data more than another's
- [ ] No magic numbers or strings: all literals are named constants
- [ ] Dependency direction verified: nothing in inner layers imports from outer layers
Add this checklist as `## Self-Review Checklist` to the implementation-notes output format.

**4. "What Would Kent Do?" Moment**
After the refactor pass, the developer must apply Kent Beck's four Simple Design rules
in order and record whether each rule is satisfied:
1. Passes all tests — yes/no
2. Reveals intention — yes/no (if no: what was renamed/extracted)
3. No duplication — yes/no (if no: what was deduplicated)
4. Fewest elements — yes/no (if no: what was removed)
Add this as `## Simple Design Verification` to the output format.

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/developer.md`.

---

## PROMPT 5 — Upgrade the Code Reviewer: Design Critic

> **Goal: make the code reviewer reason about design, not just check a compliance list.**

Read:
- `.claude/agents/code-reviewer.md`
- `ARCHITECTURE_RULES.md` sections VI, VII, VIII
- `templates/claude-feature-team/.claude/agents/code-reviewer.md`

The code reviewer must be upgraded in three areas:

**1. Design Narrative**
The reviewer must open every review with a 2-3 sentence "Design Narrative" — a plain-English
description of what the implementation is actually doing architecturally. This forces the
reviewer to understand the design before critiquing it. If the reviewer cannot write a
coherent design narrative, the implementation is too complex and must be refactored first.

**2. Verify the Developer's Self-Review**
The reviewer must explicitly check the developer's `## Self-Review Checklist` and
`## Simple Design Verification` from `implementation-notes.md` against the actual code.
If the developer marked something as passing but the code shows otherwise, that discrepancy
is itself a finding — not just the underlying issue.

**3. Design Score**
Every review must produce a Design Score across four dimensions (1-5 each):
- **Clarity**: Does the code reveal its intent without comments?
- **Cohesion**: Does each class/module do one well-defined thing?
- **Coupling**: Are dependencies minimal and pointing the right direction?
- **Craft**: Was the refactor pass taken seriously? (Check the named refactoring log)

Score of 3+ on all dimensions = APPROVED. Any dimension below 3 = CHANGES REQUESTED with
specific named refactoring operations to improve that dimension.

Add this scoring to the `code-review-report.md` output format.

Apply the same upgrade to `templates/claude-feature-team/.claude/agents/code-reviewer.md`.

---

## PROMPT 6 — Create the `/design-review` Skill

> **Goal: a standalone skill for reviewing the design of any code, not tied to a feature delivery.**

Read `.claude/skills/SKILL_TEMPLATE.md` (created in Prompt 1).
Read `ARCHITECTURE_RULES.md` in full.

Create `.claude/skills/design-review/SKILL.md` following the canonical template.

This skill is invoked when someone says:
- "Review the design of [file/class/module]"
- "Is this well-designed?"
- "Check [file] for design smells"
- "Would Fowler approve of this?"

The skill must:

**Context to load:** `ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`, the target file(s),
and any files the target imports from or that import it (one level of dependency graph).

**Process:**
1. Write a Design Narrative (what is this thing trying to do?)
2. Map the dependency graph for the target (what does it depend on? what depends on it?)
3. Check all 10 Fowler smells from ARCHITECTURE_RULES.md Section VII
4. Check Clean Architecture layer placement (Section I)
5. Check Simple Design rules (Section VI) in priority order
6. Check Sandi Metz constraints (class/method size, param count)
7. Identify the single most impactful refactoring to apply first
8. Produce a Design Report

**Output format** (`design-report.md` in `.claude/feature-workspace/`):
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

**Standalone mode:** Works without any MCP or external tools. Reads local files only.

Copy the completed skill to `templates/claude-feature-team/.claude/skills/design-review/SKILL.md`.

---

## PROMPT 7 — Create the `/adr` Skill

> **Goal: make writing Architecture Decision Records effortless and consistent.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `ARCHITECTURE_RULES.md`.
Read any existing ADRs in `docs/adr/` if they exist.

Create `.claude/skills/adr/SKILL.md` following the canonical template.

Triggers:
- "Write an ADR for [decision]"
- "Document the decision to [do X]"
- "We decided to [X], record it"
- "Create an architecture decision record"

**Context to load:** All existing ADRs (to determine next number and match tone),
`ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`, `.claude/feature-workspace/architecture-notes.md`
if it exists.

**Process:**
1. Determine the next ADR number (scan `docs/adr/` for existing files)
2. Ask the user: "What was decided?" (one sentence, active voice)
3. Ask: "What was the context that made this decision necessary?"
4. Ask: "What alternatives were considered and why were they rejected?"
5. Ask: "What are the consequences — what becomes easier, harder, or different?"
6. Ask: "Does this decision produce a fitness function? If so, how is it enforced?"
7. Draft the ADR and show it to the user for approval
8. On approval, write to `docs/adr/ADR-[NNN]-[kebab-title].md`
9. Update `docs/adr/README.md` (index of all ADRs) — create it if it doesn't exist

**Output format:**
```markdown
# ADR-[NNN]: [Title]

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
- ADR titles use active voice: "Use X" not "Decision to use X"

**Standalone mode:** Works fully without MCP. Reads local files, writes local files.

Copy to `templates/claude-feature-team/.claude/skills/adr/SKILL.md`.

---

## PROMPT 8 — Create the `/five-whys` Skill

> **Goal: structured root cause analysis for bugs, outages, and recurring problems.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read CompanyOS `co-five-whys` description for inspiration (from the article:
"Ask why iteratively, never answer for the user, reflect back what you hear, stop when
you hit something actionable. End with questions, not solutions.")

Create `.claude/skills/five-whys/SKILL.md` following the canonical template, adapted for
software development root cause analysis.

Triggers:
- "Five whys on [problem]"
- "Root cause [issue/bug/failure]"
- "Why did [X] happen?"
- "Post-mortem on [incident]"
- Any mention of recurring bugs, repeated failures, or "this keeps happening"

**Context to load:** Any relevant files mentioned (error logs, test failures, feature workspace
artifacts, git log for affected files).

**Process:**
1. State the problem clearly: "The symptom is: [X]"
2. Ask "Why did [X] happen?" — wait for the user's answer, never suggest one
3. Reflect back: "So [X] happened because [their answer]. Let me ask..."
4. Ask "Why did [their answer] happen?" — wait, never suggest
5. Repeat until one of these stop conditions is met:
   - The answer is something actionable (a process gap, a missing test, a design smell,
     a missing fitness function)
   - The user says "I don't know" twice in a row (surface this as an investigation gap)
   - Five iterations complete (the method is called Five Whys for a reason)
6. Produce a Root Cause Report

**Critical rules for the facilitation:**
- Never answer your own questions. If you have a hypothesis, keep it to yourself.
- Reflect back exactly what the user said before asking the next why.
- Ask one question at a time. Never stack questions.
- If the user's answer is vague ("because it broke"), ask for specificity:
  "Can you be more specific? What exactly broke?"
- End with questions and a recommended action, never a lecture.

**Output format** (`.claude/feature-workspace/five-whys-[kebab-topic].md`):
```markdown
# Five Whys: [Problem Statement]

Date: [today]
Facilitator: Claude
Participants: [from git config]

## The Symptom
[What was observed — specific, observable]

## The Why Chain
1. Why did [symptom] happen? → [answer]
2. Why did [answer] happen? → [answer]
3. Why did [answer] happen? → [answer]
4. Why did [answer] happen? → [answer]
5. Why did [answer] happen? → [root cause]

## Root Cause
[The actionable thing at the bottom of the chain]

## Investigation Gaps
[Any "I don't know" answers that need further investigation]

## Recommended Action
- Immediate: [what to do right now]
- Preventive: [what fitness function, test, or process change prevents recurrence]
- Owner: [who should act on this]
```

**Guardrails:**
- Never suggest an answer during facilitation
- Never skip the reflection step
- Never produce the report until the why chain is complete or a stop condition is hit

**Standalone mode:** Fully conversational. No external tools needed.

Copy to `templates/claude-feature-team/.claude/skills/five-whys/SKILL.md`.

---

## PROMPT 9 — Create the `/complexity-check` Skill

> **Goal: on-demand complexity analysis of any file or directory.**

Read `.claude/skills/SKILL_TEMPLATE.md`.
Read `ARCHITECTURE_RULES.md` sections IV (Micro-Rules) and VII (Refactoring Catalog).

Create `.claude/skills/complexity-check/SKILL.md` following the canonical template.

Triggers:
- "Check complexity of [file]"
- "Is [file] too complex?"
- "Complexity report on [directory]"
- "What's the most complex file in [package]?"
- Any mention of "technical debt" without a specific task attached

**Context to load:** Target file(s). `ARCHITECTURE_RULES.md` Section IV for thresholds.

**Process:**
1. Determine the language(s) of the target
2. Run the appropriate complexity tool:
   - TypeScript: `npx eslint --rule '{"complexity": ["error", 6]}' [file]`
   - Python: `radon cc [file] -s` or `flake8 --max-complexity=6 [file]`
   - Go: `gocyclo -over 6 [file]`
   - Java: Report that Checkstyle/SonarQube should be used, show manual heuristics instead
3. For each function/method exceeding complexity 6:
   - Name the specific smell (nested conditionals, long switch, multiple responsibilities)
   - Name the Fowler refactoring operation to apply
   - Estimate the effort: trivial / one session / needs design discussion
4. Rank findings by impact: highest complexity first, then by how many other functions
   depend on the complex one (depth in the call graph)
5. Produce a Complexity Report

**Output format** (`.claude/feature-workspace/complexity-report.md`):
```markdown
# Complexity Report: [target]

Date: [today]
Threshold: Cyclomatic complexity > 6 (ARCHITECTURE_RULES.md Section IV)

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
Suggested order to attack the debt:
1. [First because: quick win / unblocks other refactors / highest risk]
2. ...

## What's Clean
[Functions that are well within threshold — always acknowledge good work]
```

**Guardrails:** Read-only. Never modify files. Never run tests.
If complexity tooling isn't installed, fall back to manual analysis with clear caveats.

**Standalone mode:** Falls back to manual heuristic analysis if CLI tools unavailable.
Still produces a report — just note it's heuristic, not measured.

Copy to `templates/claude-feature-team/.claude/skills/complexity-check/SKILL.md`.

---

## PROMPT 10 — Create Global Rules Files

> **Goal: rules that apply across ALL agents, not just within individual agent prompts.**

Read:
- All files in `.claude/agents/`
- `ARCHITECTURE_RULES.md`
- `CLAUDE.md` at repo root
- `docs/CLAUDE.md`

Create `.claude/rules/` directory. Create the following three rule files:

**`.claude/rules/design-principles.md`**
Cross-cutting design principles every agent must honor, regardless of role:
- The four Simple Design rules (Beck) in priority order — with TypeScript examples
- The 10 Fowler refactoring operations — with the smell that triggers each one
- The anti-pattern radar (from Prompt 3: distributed monolith, anemic domain, god object,
  shotgun surgery, leaky abstraction)
- The Sandi Metz hard limits (class/method size, param count) — non-negotiable
- The Boy Scout Rule as an active instruction with examples, not a passive suggestion
- Naming standards: intention-revealing names, no abbreviations, no generic terms,
  boolean prefix rules (`is`, `has`, `can`, `should`)

**`.claude/rules/architecture-guardrails.md`**
Hard constraints that cannot be overridden by any agent or user instruction:
- Dependency direction: inner layers never import outer layers (with examples per language)
- No destructive migrations ever: Expand/Contract pattern enforced
- No secrets hardcoded anywhere: `.env` placeholders only
- No raw `any` in TypeScript: use `unknown` when genuinely unknown
- No custom retry loops: use `CircuitBreaker` or `ExponentialBackoffStrategy`
- No N+1 queries: eager loading required
- No implicit timeouts: every network call gets an explicit timeout
- Every architectural decision produces a fitness function or is explicitly flagged
  as "judgment-only"

**`.claude/rules/approval-gates.md`**
Anything irreversible requires explicit human approval. Inspired by CompanyOS's hard gate:
"Any edit resets the approval. You see the updated version, then approve again."

List every irreversible action in your pipeline and its required approval:
- Sending to Friday (QA agent): approval gate after tests pass
- Creating a git commit: user must say "commit" explicitly
- Running database migrations: user must say "migrate" explicitly
- Posting to external APIs: user must say "send" or "approve" explicitly
- Dropping files outside `.claude/feature-workspace/`: user must confirm file path
- Wiring a new fitness function into CI: architect approves, devops implements

Format each gate:
```
Action: [what it does]
Irreversible because: [why it can't be undone]
Gate: user must say "[exact word(s)]"
Reset condition: any edit to the pending artifact resets the gate
```

After creating all three rule files, update the header of every agent file in
`.claude/agents/` to add a reference:
```
Read `.claude/rules/design-principles.md`, `.claude/rules/architecture-guardrails.md`,
and `.claude/rules/approval-gates.md` before beginning any task.
```

Copy all three rule files to `templates/claude-feature-team/.claude/rules/`.

---

## PROMPT 11 — Update `scaffold-team.sh` and `install`

> **Goal: make sure all new skills, rules, and agent upgrades are deployed automatically.**

Read `scaffold-team.sh` and `install`.

Update `scaffold-team.sh` to copy:
- `.claude/agents/` (all 8 agents including code-reviewer)
- `.claude/skills/` (all skills: deliver-feature, new-feature, review-pr, retrospective,
  design-review, adr, five-whys, complexity-check)
- `.claude/rules/` (design-principles, architecture-guardrails, approval-gates)
- `features/TEMPLATE.md`
- `DOMAIN_DICTIONARY.template.md`
- `CLAUDE.md`

Add a verification step that counts deployed files and reports:
```
✅ Agents deployed:  8
✅ Skills deployed:  8
✅ Rules deployed:   3
✅ Templates:        2
```

Update `install` to symlink `.claude/rules` in addition to `.claude/agents`.

Update the `FILES_TO_LINK` array:
```bash
".claude/agents"
".claude/skills"
".claude/rules"
```

After updating both scripts, do a dry-run check: verify all source paths referenced in
the scripts actually exist in the repo. Report any missing paths before finishing.

---

## Final Verification

After all 11 prompts are complete, run this check from the repo root:

```bash
echo "=== Agents ===" && ls .claude/agents/
echo "=== Skills ===" && ls .claude/skills/
echo "=== Rules ===" && ls .claude/rules/
echo "=== Template agents ===" && ls templates/claude-feature-team/.claude/agents/
echo "=== Template skills ===" && ls templates/claude-feature-team/.claude/skills/
echo "=== Template rules ===" && ls templates/claude-feature-team/.claude/rules/
echo "=== Skill template exists ===" && cat .claude/skills/SKILL_TEMPLATE.md | head -5
echo "=== Rules referenced in analyst ===" && grep -l "rules/" .claude/agents/analyst.md
echo "=== ADR skill exists ===" && ls .claude/skills/adr/
echo "=== Design review skill exists ===" && ls .claude/skills/design-review/
```

All commands should produce output with no errors.
The total should be: 8 agents, 8 skills, 3 rules files, mirrored in templates/.
