---
name: saturday-test-advisor
description: Audits a Saturday-based E2E/UI test suite for structural adherence to the Site-Centric pattern — flags BasePage / BaseFlow / BasePartial / Filter-decorated methods defined in code but not exercised by any Gherkin scenario, and scenarios that reference primitives that don't exist. Reports gaps as "either delete this or add a scenario," never fabricates scenarios (Gherkin scenarios are business language that humans should write). Use when the user asks "audit my Saturday coverage", "any dead pages/flows?", "what am I missing tests for?", or after adding a new page/flow/partial.
triggers:
  keywords: ["saturday audit", "structural adherence", "dead pages", "dead flows", "unused partials", "site-centric audit", "saturday coverage", "orphaned page"]
  intentPatterns:
    - "audit my Saturday test coverage"
    - "which pages/flows/partials aren't tested"
    - "any orphaned Site-Centric classes"
    - "check my Gherkin scenarios cover the code"
standalone: true
---

## When To Use
Use when the user wants to verify their Saturday test suite structurally adheres to the Site-Centric
pattern — every defined page/flow/partial/filter has real Gherkin scenario coverage, and no scenario
references a class that doesn't exist. Sibling to `sunday-test-advisor` (which audits go-sunday YAML
specs) but structurally different in mechanism: this one inspects code against `.feature` files, since
Saturday has no YAML manifest to audit against — the code IS the source of truth.

Do NOT use for:
- Sunday API tests — use `sunday-test-advisor` (go-sunday YAML specs) or read the API Test Coverage
  Matrix in `docs/patterns/sunday-framework-patterns.md` for manual audit of other-language Sunday
  tests.
- Unit test coverage — that's `run-tests`'s coverage check, or `unit-tester` / `backfill-unit-tests`
  for the writing side. This skill is only about E2E/UI adherence to Saturday primitives.
- Deciding *what* Gherkin scenarios to write — that's `qa-engineer`'s job (usually inside
  `deliver-atdd` or `deliver-feature`). This skill spots gaps; the human or `qa-engineer` fills them.

## Context To Load First
1. `docs/patterns/saturday-framework-patterns.md` — the canonical Site-Centric pattern reference. The
   audit uses this to know which primitives (`BaseSite`, `BasePage`, `BaseElement`, `BasePartial`,
   `BaseFlow`, `SiteManager`, `TabManager`, `Filters`) it's looking for.
2. `shared/DOMAIN_DICTIONARY.md` — Saturday Framework Terms section, for the naming conventions the
   audit expects to see.
3. The project's own `CLAUDE.md` — for any project-specific Saturday conventions that override the
   defaults.

## Process

### Step 1 — Discover Saturday primitives in code
Grep the project's source for class/type declarations that extend or implement the Saturday
primitives. Handle each supported language's syntax:

- **TypeScript**: `class \w+ extends BasePage`, `extends BaseFlow`, `extends BasePartial`,
  `extends BaseElement`
- **C#**: `class \w+ : BasePage`, `: BaseFlow`, `: BasePartial`, `: BaseElement`
- **Python**: `class \w+\(BasePage\)`, `\(BaseFlow\)`, `\(BasePartial\)`, `\(BaseElement\)`
- **Java**: `class \w+ extends BasePage`, `extends BaseFlow`, `extends BasePartial`, `extends BaseElement`
- **Filter-decorated methods**: search for `@RequiresFilter` (TS/C#/Java) or `@requires_filter`
  (Python) — these are the methods whose filter-pass and filter-fail branches both need scenario
  coverage.

Build a set of every defined page, flow, partial, and filter-decorated method (each with the file it
was defined in, for the report).

### Step 2 — Discover Gherkin scenarios and what they reference
Find every `.feature` file in the project (`**/*.feature`, excluding `node_modules`, `target`,
`bin`, `.venv`, `venv`, `dist`, `build`). For each scenario:
- Record its name and any `@tags`.
- Read its steps and Look for references to any page/flow/partial name from Step 1.
  Match by:
  - Class name appearing in step text (fuzzy — `CheckoutPage`, `checkout page`, `checkout` all
    count as references to `CheckoutPage`).
  - Explicit tag references (`@page:CheckoutPage`, `@flow:CompleteCheckoutFlow` — if the project uses
    that convention).

### Step 3 — Compute the structural-adherence report
Cross-reference the two sets from Steps 1 and 2 and produce four categories of finding:

1. **Orphaned pages/flows/partials** — defined in code, referenced by zero scenarios. Either dead
   code or missing coverage.
2. **Filter branch gaps** — methods decorated with `@RequiresFilter` / `@requires_filter` where
   scenarios exercise the filter-passing branch but not the filter-failing branch (or vice versa).
   Best-effort: detected by looking for scenarios whose steps mention the method AND include failure
   keywords (`fails`, `blocked`, `denied`, `throws`, `error`, or step assertions on `FilterException`).
   Flag as "best-effort — verify manually" rather than claiming certainty.
3. **Broken references** — Gherkin scenario references a page/flow/partial name that doesn't exist in
   code. Either a typo in the scenario or the class was renamed/deleted without updating the feature
   file.
4. **Under-partialed pages** — best-effort heuristic: two or more `BasePage`s that share ≥3
   identical `BaseElement`s (e.g., `header`, `footer`, `navigationMenu`) suggest a missing
   `BasePartial`. Not enforced — reported as "consider extracting a partial." Skip this check if the
   project has fewer than 5 total `BasePage` subclasses (not enough signal).

### Step 4 — Present the report
Output a markdown report inline (do not persist to a file — the user should decide what to do):

```markdown
## Saturday Structural Adherence Audit: <project name>

### Orphaned (defined in code, zero scenario coverage)
| Type    | Name              | Defined in                              | Recommendation                                  |
|---------|-------------------|-----------------------------------------|-------------------------------------------------|
| Page    | ProfileEditPage   | packages/site/src/pages/profile-edit.ts | Add a scenario OR delete as dead code           |
| Flow    | GuestCheckoutFlow | packages/site/src/flows/guest.ts        | Add a scenario OR delete as dead code           |
| Partial | GlobalFooter      | packages/site/src/partials/footer.ts    | Add a scenario OR delete as dead code           |

### Filter branch gaps (best-effort — verify manually)
| Filter method                      | Passing branch tested | Failing branch tested | Recommendation                        |
|------------------------------------|-----------------------|-----------------------|---------------------------------------|
| DashboardPage.viewAdminSettings    | Yes                   | No                    | Add a scenario for the non-admin user |

### Broken references (scenario points at non-existent class)
| Scenario file                                | Scenario name                | Missing reference   |
|----------------------------------------------|------------------------------|---------------------|
| features/checkout.feature                    | Guest checkout with saved card | GuestCheckoutFlowV2 (typo? or deleted?) |

### Consider extracting a partial (best-effort)
- LoginPage, DashboardPage, SettingsPage all reference identical `header`, `footer`, `navigationMenu`
  elements — candidate for a `GlobalChromePartial`.

---

**Summary**: 3 orphaned primitives, 1 filter branch gap, 1 broken reference, 1 partial-extraction
candidate. Report is advisory — no changes made.
```

If any category is empty, either omit the section or write "None — all defined primitives are exercised
by at least one scenario."

**Do NOT proceed past this step.** This skill reports; it does not modify code or `.feature` files.

## Output Format
Inline markdown table as shown in Step 4. No file persistence.

## Guardrails
- **Never modify code** — this skill only reports. Deleting supposedly-dead code is the user's
  decision, not this skill's.
- **Never modify `.feature` files** — same reason. Adding scenarios is `qa-engineer`'s job (usually
  inside `deliver-atdd`), not this skill's. Fabricating scenarios violates the "scenarios are business
  language humans write" principle in `docs/patterns/testing-pyramid.md`.
- **Never claim certainty on the best-effort checks** (filter branch coverage, partial-extraction
  candidates) — always label them as "best-effort — verify manually."
- **If pattern-matching heuristics can't find a class or scenario reference reliably** (e.g., the
  project uses a naming convention this skill doesn't recognize), say so explicitly in the report
  rather than silently missing it.
- **If no `.feature` files exist**, that IS the audit result — say "no Gherkin scenarios found;
  Saturday's whole test surface is empty" rather than reporting an artificially clean audit.
- **If no Site-Centric primitives are found in code** (`BasePage`/`BaseFlow`/etc.), the project isn't
  actually using Saturday — say so and exit rather than producing a misleading empty report.

## Standalone Mode
Pure Read + Grep + Glob operations. No external tools, no MCP servers, no code generation. Every step
of the audit is deterministic reading of the local filesystem.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
