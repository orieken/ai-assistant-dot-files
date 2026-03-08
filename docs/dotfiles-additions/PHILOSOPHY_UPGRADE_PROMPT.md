# Philosophy Upgrade Prompt

Use this prompt in a new chat with your dotfiles repository attached or accessible.
Paste it as-is — it is written as a direct instruction to the AI.

---

You have access to my AI assistant dotfiles repository. I want you to upgrade the
craftsmanship philosophy embedded in the existing agent files and global rule files
to reflect the full depth of the following engineering pioneers:

- **Robert C. Martin (Uncle Bob)** — Clean Architecture, SOLID, Clean Code rules,
  the newspaper metaphor for code organization, functions that do one thing
- **Kent Beck** — TDD as a *design tool* (not just a safety net), the four rules of
  simple design, Red-Green-Refactor with emphasis on the Refactor step
- **Martin Fowler** — The refactoring catalog (named operations, not vague cleanup),
  Patterns of Enterprise Application Architecture, evolutionary database design,
  the strangler fig pattern for incremental rewrites
- **Neal Ford** — Evolutionary architecture, fitness functions as first-class
  architectural artifacts, architecture as a continuous activity not a phase
- **Michael Feathers** — Working Effectively with Legacy Code: the seam concept,
  characterization tests, how to get untested code under test safely
- **Sandi Metz** — Practical Object-Oriented Design: the rules (classes ≤ 100 lines,
  methods ≤ 5 lines, ≤ 4 params, one object instantiated per controller), the
  "squint test" for code shape, composition over inheritance
- **Dan North** — BDD as originally conceived: scenarios describe *observable behavior
  from the outside*, not internal implementation steps. "Given/When/Then" is a
  communication tool first, a test structure second
- **Dave Farley** — Continuous Delivery philosophy: the deployment pipeline is a
  product, test isolation is non-negotiable, the difference between acceptance tests
  (behavior) and unit tests (design feedback)
- **Adam Shostack** — Threat modeling with STRIDE: four questions (what are we
  building, what can go wrong, what do we do about it, did we do a good job),
  security as a design property not a testing activity

## What to upgrade

### 1. `ARCHITECTURE_RULES.md`
This is your most important file — every AI assistant reads it. Upgrade it to:
- Add a **Section VI: Simple Design Rules** with Beck's four rules in priority order,
  with concrete TypeScript/Go/Python examples for each
- Add a **Section VII: Refactoring Catalog** with the 10 most applicable Fowler
  operations for this codebase, each with a "smell that triggers it" and a
  "concrete before/after example"
- Add a **Section VIII: Fitness Functions** explaining what they are (Ford), why
  every significant architectural decision should produce one, and 5 examples
  relevant to the Saturday/Sunday ecosystem with the tooling to enforce them
- Add a **Section IX: Legacy Code Protocol** (Feathers): when touching untested
  code, write characterization tests first, identify seams before refactoring,
  never refactor and add behavior in the same commit
- Upgrade **Section V: Testing Principles** to include Dan North's BDD framing
  (scenarios describe behavior observable from outside the system) and Dave
  Farley's distinction between acceptance tests and unit tests

### 2. `docs/CLAUDE.md`
This is your Claude-specific global rules file. Upgrade it to:
- Add Sandi Metz's rules as enforceable constraints alongside the existing
  cyclomatic complexity and LOC limits
- Add a **"What to do when touching legacy code"** section using Feathers'
  approach
- Add a **"Refactor step is not optional"** section making explicit that
  Red-Green-*Refactor* means the refactor step happens before the next feature,
  not someday
- Upgrade the **Naming** section with Fowler's "intention-revealing names"
  principle and concrete bad/good examples from TypeScript

### 3. `.claude/agents/developer.md` (and `templates/claude-feature-team/.claude/agents/developer.md`)
Upgrade the developer agent to:
- Lead with the TDD cycle as a *design activity*: write the interface/type
  signature first, write the failing test that describes the behavior, implement
  to make it pass, then refactor — in that order, every time
- Add the Sandi Metz rules as hard constraints alongside the existing LOC limit
- Add an explicit **Refactor Pass** step after green: agent must check for
  applicable Fowler refactoring operations before declaring implementation done
- Add the Boy Scout Rule as an active instruction: if you touch a file that has
  complexity ≥ 6 or functions > 25 lines that are *not* part of the current
  feature, extract and clean them — leave it better than you found it

### 4. `.claude/agents/qa-engineer.md` (and templates version)
Upgrade the QA agent to:
- Embed Dan North's BDD philosophy explicitly: scenarios describe behavior
  observable from *outside* the system. "Given the CartService has been
  initialized" is wrong. "Given the cart has 3 items" is right.
- Add Dave Farley's acceptance test philosophy: acceptance tests verify the
  system does what the business needs. They should be written in terms of
  *what*, never *how*.
- Add Michael Feathers' approach for testing legacy/untested code: write
  characterization tests to document current behavior before any change,
  use seams to introduce testability without changing production code

### 5. All global AI config files
Apply consistent upgrades to:
- `.cursorrules` / `.cursor/rules/global.mdc`
- `.windsurfrules`
- `.openai.md`
- `.github/copilot-instructions.md`
- `.gemini/antigravity/instructions.md`
- `.claude.md`

Each should gain:
- A **Simple Design** section (Beck's 4 rules, in order)
- A **Refactoring** section (name the operation, don't just say "clean it up")
- A **Fitness Functions** mention (architectural decisions produce measurable constraints)
- The **Boy Scout Rule** upgraded from a one-liner to an active instruction with
  examples of what "leave it cleaner" means in practice

## Style constraints for the upgrades

- Match the existing tone and formatting of each file exactly
- Do not remove any existing content — only add and upgrade
- Where you add examples, use TypeScript as the primary language (it is the
  dominant language in this ecosystem) with Go secondary
- Keep additions proportional — upgrade depth, not length for its own sake
- Every principle added must have a *concrete, actionable* form, not just a
  philosophical statement. "Apply the Strategy pattern when behavior varies by
  context and you have 3+ conditional branches" is actionable.
  "Write clean code" is not.

## What NOT to change

- Do not modify the Saturday/Sunday ecosystem-specific sections — those are
  already accurate and up to date
- Do not change the tech stack sections
- Do not change file naming conventions — those are already correct
- Do not add new files — upgrade the existing ones only

Begin by reading all the files listed above, then produce the upgrades one file
at a time, showing the diff or the full updated content for each.
