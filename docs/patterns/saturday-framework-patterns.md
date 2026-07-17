# Saturday Framework Patterns (E2E / UI Testing)

The Site-Centric architecture. Every term below matches `shared/DOMAIN_DICTIONARY.md` exactly — don't
introduce synonyms (`PageObject`, `App`, `Widget`, etc. are explicitly listed there as terms to avoid).

## Mental Model

Saturday organizes browser-based automation into three primary concepts and two supporting ones. Read
top-down:

- **Site** — the orchestrator. One `BaseSite` per web application under test; owns the browser context
  lifecycle and roots the entire test surface for that app.
- **Pages** — the anatomy of the application. A `BasePage` is one screen or view; it composes:
  - **Elements** — the raw interactive components (`BaseElement`).
  - **Filters** — state-based gating and selection logic (see the `Filters` entry).
  - **Partials** — shared UI sections that appear across multiple pages (`BasePartial`) — headers,
    footers, global navigation.
- **Flows** — multi-step user journeys (`BaseFlow`) that span multiple pages, parameterized by test data.
- **Models & Factories** — the test data flowing through everything. `Model`s are the domain objects
  the test reasons about (`TestUser`, `TestOrder`); `Factory`s produce them with realistic defaults using
  this framework's per-language factory conventions.
- **Coordinators** — `SiteManager` (multi-application journeys) and `TabManager` (multi-tab within one
  browser context). Reach for these only when a single Site or single tab isn't enough.

The rest of this file is one entry per concept, in the order above.

## Site-Centric Pattern

**Context**: The overarching pattern every component below implements a piece of — Saturday's answer to
"how do you structure a large E2E suite so it doesn't rot into a pile of unrelated page objects."
Explicitly named (`shared/rules/testing-conventions.md`: "ALWAYS use the Site-Centric pattern") and
explicitly not the traditional Page Object Model — POM is named in the same rule as something to NEVER
use.

**Structure**: A `BaseSite` roots one web application's entire test surface — every `BasePage` it owns,
every `BasePartial` reused across those pages, every multi-page `BaseFlow` that crosses them, every
state-gated `Filter`. Multiple `BaseSite`s (one per application under test) are coordinated by a
`SiteManager`; multiple browser tabs within any of them are coordinated by a `TabManager`. The pattern's
actual claim is narrower than "organize pages into classes" — it's that *application* is the top-level
organizing unit, and everything else (pages, partials, flows, elements, filters) nests inside exactly
one `BaseSite`, never spanning two.

**Why not POM**: Traditional Page Object Model has no answer for "this flow spans three pages" (it
either bloats one page object with unrelated methods, or leaves the flow logic floating in the test
itself), no answer for "this header appears on every page" (each page object ends up with its own copy
of the header locators, or a base class awkwardly leaks header knowledge into every page), and no answer
for "this suite drives two separate applications" (nothing stops a page object from silently coupling
to another application's internals). `BaseFlow`, `BasePartial`, and `SiteManager` are Site-Centric's
direct answers to those three specific POM failure modes, not incidental extras.

**Trade-offs**: More ceremony than POM for a single-page, single-app smoke test — a `BaseSite` with one
`BasePage` is more scaffolding than a bare page object would need. Pays off precisely at the scale POM
struggles with: multi-page flows, cross-page shared UI, multi-application suites, and cross-project
consistency (see `BaseSite`'s own Trade-offs entry below).

**Related**: Every entry in this file is a Site-Centric component. Sunday's `Declarative API Client
Pattern` is the equivalent umbrella pattern for API testing — see `sunday-framework-patterns.md`.

## BaseSite

**Context**: The root orchestrator for a single web application under test. Every E2E suite has exactly
one `BaseSite` per application it drives.

**Structure**: Owns the browser context lifecycle and holds references to every `BasePage` reachable
within that application. Delegates cross-application coordination to `SiteManager` rather than knowing
about other applications itself — a `BaseSite` should never import or reference another `BaseSite`
directly.

**Example**: A checkout flow that spans a storefront and a separate payment provider is two `BaseSite`s
(`StorefrontSite`, `PaymentProviderSite`), coordinated by one `SiteManager`, not one `BaseSite` that
knows about both.

**Trade-offs**: Adds a layer of indirection for single-application suites where it isn't strictly
needed. Worth it anyway for consistency — every suite in this framework looks the same, so a developer
moving between projects doesn't have to relearn the shape.

**Related**: `SiteManager` (owns multiple `BaseSite`s), `BasePage` (owned by exactly one `BaseSite`).

## BasePage

**Context**: Represents a single page within a `BaseSite` — one screen, view, or document the user
interacts with. This is the direct replacement for what other frameworks call a Page Object — the
traditional POM. This framework explicitly rejects POM (see `shared/rules/testing-conventions.md`:
"NEVER use traditional Page Object Model (POM)") in favor of the Site-Centric hierarchy: a `BasePage`
composes from three smaller pieces (`BaseElement`, `Filters`, `BasePartial`) rather than swallowing them
whole.

**Structure**: Holds `BaseElement` references for the interactive parts of the page, composes any
`BasePartial`s that appear on this page (headers, footers, global navigation), and exposes page-level
actions — not raw element interactions — to callers. A `BasePage` method should read like a user action
(`login(username, password)`), never like a DOM operation (`clickButton('#submit')`).

**Example**: `LoginPage extends BasePage` exposes `login(user, pass): Promise<DashboardPage>` — note the
return type: a page action that navigates should return the `BasePage` the user lands on, so the test
reads as a chain of real user steps.

**Trade-offs**: Requires every navigation-causing action to know what page it lands on, which is extra
type bookkeeping. Pays for itself by making broken navigations a compile-time error instead of a runtime
one.

**Related**: `BaseElement` (page-scoped components), `BasePartial` (cross-page components composed into
this page), `Filters` (state-based selection often lives alongside a page's elements), `BaseFlow`
(chains multiple `BasePage`s together).

## BaseElement

**Context**: A reusable UI component abstraction within a `BasePage` — the thing that actually wraps a
Playwright locator.

**Structure**: Wraps interaction (`click`, `fill`, `getText`) and waiting/assertion logic for one
component. Composable: a `BasePage` is built from several `BaseElement`s, and a complex `BaseElement`
(e.g. a data table) can itself be built from smaller `BaseElement`s.

**Example**: A `DropdownElement extends BaseElement` that wraps `select`/`getSelectedOption` so no test
ever calls Playwright's raw `page.selectOption()` directly.

**Trade-offs**: More classes than calling Playwright's API directly. The payoff: when the UI library
changes how dropdowns render, exactly one `BaseElement` subclass changes, not every test that touches a
dropdown.

**Related**: `Filters` (declarative selection logic often lives alongside a `BaseElement`), `BasePartial`
(the same abstraction generalized to cross-page shared UI).

## BasePartial

**Context**: A shared, repeating UI section that appears across multiple pages — headers, footers,
global navigation, cookie banners, notification trays. The gap `BaseElement` doesn't fill: `BaseElement`
is described as "within a `BasePage`," but a global header isn't within any single page — it's *across*
them. Modeling it as a plain `BaseElement` on every page duplicates locators; jamming it into a `BasePage`
base class couples every page to it awkwardly. `BasePartial` is the dedicated seat for cross-page UI.

**Structure**: Lives in a `partials/` subdirectory within the site (e.g.,
`sites/storefront/pages/partials/global-header.ts`). **Deliberately does NOT follow the `FooPage` naming
convention** — it isn't a page, and the naming difference is a signal to the reader. Exposes the same
kinds of operations a page does (locators, interactions, its own `Filters` where relevant), but is
composed *into* pages rather than being one.

**Example**: A `GlobalHeader extends BasePartial` with login state, cart badge, and search bar. Every
page that shows it exposes a `header(): GlobalHeader` accessor rather than duplicating the locators.
When the header's login-state selector changes, exactly one file changes — not every `BasePage`
subclass that happens to reference it.

**Trade-offs**: One more concept than "just use `BaseElement` for everything." Pays off precisely when
shared UI has genuinely different lifecycle and testing concerns than a page-scoped button — a header
whose login state affects every page's assertions is exactly that case. If the "partial" is only used
by one page, it's a `BaseElement`, not a `BasePartial`.

**Related**: `BasePage` (composes `BasePartial`s), `BaseElement` (the page-scoped equivalent),
`Filters` (partials often have their own state to filter on — e.g., logged-in vs. logged-out header).

## BaseFlow

**Context**: A multi-step user journey that spans multiple pages — the thing a POM has no good answer
for, since a POM only models one page at a time. Flows exist to hide *setup noise* (the sequence of
page interactions that gets a test to its interesting point), not to hide the interesting point itself.

**Structure**: Orchestrates a sequence of `BasePage` actions that represents one coherent user journey
(e.g. "complete checkout"), independent of any single test that happens to invoke it. Accepts
parameterization — typically a `Model` produced by a `Factory` — so the flow's mechanics stay the same
while the test data varies.

**Example**: `CompleteCheckoutFlow extends BaseFlow` chains `CartPage -> ShippingPage -> PaymentPage ->
ConfirmationPage` as one reusable unit, parameterized by a `TestUser` and a `TestCart`, so five
different tests that all need "a completed checkout" as setup call the same flow with different data
instead of five copies of the same page-chain.

**Trade-offs — Moist tests, not fully dry**: A `BaseFlow` that gets too broad turns into an
unmaintainable god-object. Keep it scoped to one journey; if a test needs to skip a step, that's a sign
the flow is drawn at the wrong grain. And keep tests **moist**, not dry: DRY the setup noise via Flows
and Factories, but keep the critical assertion path visible in the test itself. A test verifying search
results should show "user searched for X" prominently in the test body; how they clicked into the
search input should not. Over-DRYing by hiding the interesting user action inside a flow makes the test
read like a magic incantation and defeats the purpose. This is the same principle as DAMP (Descriptive
And Meaningful Phrases) — a well-known counterpoint to blanket DRY in tests, applied here specifically
to Flows.

**Related**: Composes multiple `BasePage`s; often the direct target of a test's `Arrange` step. Takes
a `Model` (produced by a `Factory`) as its data parameter.

## Model

**Context**: A domain object the test reasons about — a `TestUser`, a `TestOrder`, a `TestSearchQuery`.
Independent of any particular page's on-screen representation of it. When a test scenario says "given a
user with these credentials, when they log in, then...", the `TestUser` is a `Model`. Applies this
framework-wide `Model` convention (`CLAUDE.md` Models section — no persistence logic, no HTTP calls,
immutable where the language supports it) specifically to test data.

**Structure**: Plain data class. Readonly fields (`readonly` in TS, `frozen=True` in Python, `record`
in Java, value receiver in Go). No test-framework coupling — a `TestUser` doesn't know about Playwright,
Cucumber, or any assertion library.

**Example**: `class TestUser { readonly email: string; readonly password: string; readonly displayName:
string; readonly role: 'admin' | 'guest'; }` — used by Flows for parameterization, referenced by test
scenarios, populated by Factories, and referenced across login, checkout, profile-edit tests without
each test rewriting the shape.

**Trade-offs**: One more class per domain type. Pays off exactly when tests reason about the domain
across multiple pages — a `TestUser` referenced in login, checkout, and profile-edit tests beats each
test having its own ad-hoc `{email, password}` bag with slightly different field names.

**Related**: `Factory` (produces instances of `Model`s), `BaseFlow` (accepts a `Model` as its data
parameter).

## Factory

**Context**: Produces `Model` instances with realistic defaults, using this framework's existing
per-language factory-library convention — `fishery` + `@faker-js/faker` for TypeScript, `Bogus` +
`AutoFixture`/`AutoBogus` for C#, `Faker` + `polyfactory` for Python, `DataFaker` + `Instancio` for Java
(see the per-language rules in `shared/rules/<language>-conventions.md`). This entry is about how those
libraries are *used* in Saturday specifically; it doesn't redefine what a factory is at the language
level.

**Structure**: One Factory per `Model`, exposing named creation scenarios (`createAdmin`, `createGuest`,
`createUserWithExpiredSession`) per `CLAUDE.md`'s Factories convention. Static methods, not
instantiation.

**Example**: `UserFactory.createGuest()` returns a `TestUser` with a realistic name and email but no
admin rights, ready to hand to a Flow. `UserFactory.createAdmin()` returns one with admin privileges.
Neither test writes email/password strings directly.

**Trade-offs**: One more class per `Model`. Pays off the moment more than one test needs "a valid
user" — the defaults live in one file, the named scenarios document the interesting variations, and no
test hard-codes `"test@test.com"` inline.

**Related**: `Model` (what a Factory produces), `BaseFlow` (typically takes a Factory-produced Model as
its parameterization).

## Filters

**Context**: Declarative element selection — expressing "the row where status = 'active'" without
writing imperative loop-and-check logic in the test.

**Structure**: A small query language over a `BaseElement` collection, evaluated lazily so it composes
with Playwright's own auto-waiting instead of fighting it.

**Trade-offs**: Another abstraction to learn, but it's what keeps tests reading as intent ("select the
active row") instead of implementation ("loop rows, check the third cell's text").

## SiteManager

**Context**: Manages cross-application test contexts — the thing that owns multiple `BaseSite`s when a
journey crosses application boundaries.

**Structure**: Holds a registry of active `BaseSite`s for the current test run and coordinates handoffs
between them (e.g. a redirect from the storefront to a third-party payment provider and back).

**Related**: `TabManager` (multi-tab within one browser context) is the narrower sibling of this
(multi-application, possibly multi-browser-context) concern — don't reach for `SiteManager` when
`TabManager` is all a test actually needs.

## TabManager

**Context**: Manages multi-tab browser contexts within a single test — e.g. a "share" action that opens
a new tab.

**Structure**: Tracks which tab is "active" for the purposes of subsequent `BasePage` actions, and
provides explicit switch operations rather than letting tests guess which tab a Playwright API call will
land on.

**Trade-offs**: Explicit tab-switching adds a step to every multi-tab test, but removes an entire class
of flaky "wrong tab" failures that implicit tab-tracking produces.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
