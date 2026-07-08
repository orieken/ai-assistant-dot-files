# Saturday Framework Patterns (E2E / UI Testing)

The Site-Centric architecture. Every term below matches `shared/DOMAIN_DICTIONARY.md` exactly — don't
introduce synonyms (`PageObject`, `App`, `Widget`, etc. are explicitly listed there as terms to avoid).

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

**Context**: Represents a single page within a `BaseSite`. This is the direct replacement for what other
frameworks call a Page Object — the traditional POM. This framework explicitly rejects POM (see
`shared/rules/testing-conventions.md`: "NEVER use traditional Page Object Model (POM)") in favor of the
Site-Centric hierarchy below it.

**Structure**: Holds `BaseElement` references for the interactive parts of the page and exposes
page-level actions (not raw element interactions) to callers. A `BasePage` method should read like a
user action (`login(username, password)`), never like a DOM operation (`clickButton('#submit')`).

**Example**: `LoginPage extends BasePage` exposes `login(user, pass): Promise<DashboardPage>` — note the
return type: a page action that navigates should return the `BasePage` the user lands on, so the test
reads as a chain of real user steps.

**Trade-offs**: Requires every navigation-causing action to know what page it lands on, which is extra
type bookkeeping. Pays for itself by making broken navigations a compile-time error instead of a runtime
one.

**Related**: `BaseElement` (owned by a `BasePage`), `BaseFlow` (chains multiple `BasePage`s together).

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

**Related**: `Filters` (declarative selection logic often lives alongside a `BaseElement`).

## BaseFlow

**Context**: A multi-step user journey that spans multiple pages — the thing a POM has no good answer
for, since a POM only models one page at a time.

**Structure**: Orchestrates a sequence of `BasePage` actions that represents one coherent user journey
(e.g. "complete checkout"), independent of any single test that happens to invoke it.

**Example**: `CompleteCheckoutFlow extends BaseFlow` chains `CartPage -> ShippingPage -> PaymentPage ->
ConfirmationPage` as one reusable unit, so five different tests that all need "a completed checkout" as
setup call the same flow instead of five copies of the same page-chain.

**Trade-offs**: A `BaseFlow` that gets too broad turns into an unmaintainable god-object. Keep it scoped
to one journey; if a test needs to skip a step, that's a sign the flow is drawn at the wrong grain.

**Related**: Composes multiple `BasePage`s; often the direct target of a test's `Arrange` step.

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
