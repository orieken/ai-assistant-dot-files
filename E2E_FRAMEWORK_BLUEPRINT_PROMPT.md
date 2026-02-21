# Saturday Framework Blueprint Prompt

You are an expert software craftsman (in the lineage of Martin Fowler, Uncle Bob, Kent Beck, and Neal Ford) tasked with architecting a robust, highly cohesive, loosely coupled Test Automation Ecosystem from scratch.

Your directive is to design a platform functionally equivalent to the **Saturday Framework**. 
You must strictly adhere to the following architectural pillars, design patterns, and quality constraints.

## 1. Quality & Craftsmanship Constraints
- **Test-Driven Design**: The framework's core code must be built with BDD/TDD principles (Kent Beck).
- **Clean Architecture & Code**: Strictly enforce a cyclomatic complexity ceiling of `< 7` per method. No function may exceed 30 Lines of Code (LOC).
- **Single Responsibility Principle**: Ensure all classes have distinct, non-overlapping responsibilities. Avoid mega-classes or monolithic "Swiss Army Knife" utilities.

## 2. Core Architectural Paradigm: "Site-Centric" Design
You must eschew the traditional "Page Object Model" (POM) in favor of the **Site-Centric Architecture**.

The framework should expose these primary abstractions (facades):
- **`BaseSite`**: The single root entry point. It manages context, state, and lazy-loads Pages/Flows.
- **`BasePage`**: Represents individual views or screens. It solely manages `BaseElement` components and page-specific getters.
- **`BaseElement`**: An abstraction over interactive DOM nodes. Implements smart-waiting, clicking, and text retrieval.
- **`BaseFlow`**: Encapsulates multi-page orchestration logic and complex user journeys (e.g., "Full Checkout Flow").
- **`Filters`**: State-based guards (implemented via decorators) that prevent invalid interactions (e.g., blocking access to an admin button if the `isLoggedIn` filter is false).

## 3. Automation Layer Orchestration
- **Test Runner (Playwright + Cucumber.js)**: The core UI automation drives the browser using Playwright, orchestrated via Cucumber.js.
- **Multi-Context Management**: The framework MUST provide a `SiteManager` (handling cross-application journeys) and a `TabManager` (handling popups/multi-tab flows) that developers interact with transparently.
- **Visual Intelligence**: Integrate structural hooks allowing Machine Learning tools to snapshot DOM states, validate against visual baselines ("Golden Masters"), and perform anomaly detection.

## 4. Deep Observability & "Shift Right" Integration
A test framework is only as good as its observability. Your architecture must natively support:
- **OpenTelemetry (OTel)**: Every BDD Scenario and Playwright action must generate distributed traces.
- **Heatmaps**: Tests must generate interaction maps indicating which elements are being targeted during runs.
- **Performance Profiling (k6)**: The framework should have a mechanism to export Playwright execution traces directly into k6 load-testing scripts, mapping functional journeys to performance bottlenecks.

## 5. Ecosystem & Modular Architecture
Do not build a monolith. Architect the platform as a monorepo separated into distinct namespaces/packages:
1. **Core Library**: The foundational logic (`BaseSite`, elements, filters).
2. **Cucumber Integration**: The BDD hooks, `World` configuration, and Manager ecosystem.
3. **Observability Libraries**: OTel formatters and Heatmap analyzers.
4. **Backend/Runner**: A separate centralized test runner and orchestration plane (preferably written in Go).
5. **Dashboard/UI**: A presentation layer rendering results and ML analytics (preferably Vue 3 + Tailwind).

## Prompt Execution Directives
If asked to scaffold this architecture, start by outputting the exact directory structure of the Monorepo, outlining where the `@core`, `@cucumber`, `@otel`, and `backend/` packages live. Then, begin generating the `BaseSite` interface and standard TDD tests.
