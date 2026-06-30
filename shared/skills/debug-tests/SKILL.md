---
name: debug-tests
description: Iteratively debug full test suites and resolve testing infrastructure issues.
triggers:
  keywords: [debug test failures, fix tests, debug suite, failing tests]
  intentPatterns: ["debug the test failures", "why are my tests failing"]
standalone: true
---

## When To Use
When you are asked to debug test suites, particularly multiple failing tests or cascading infrastructure/environment issues.

## Context To Load First
1. Check the test runner configuration (`package.json`, `jest.config.js`, `vitest.config.ts`, `tests/setup.ts`, etc.).

## Process
1. Run the full test suite to understand the current scope of failures.
2. Identify common failure patterns (e.g., missing mocks, bad setups, invalid imports) rather than isolated test logic flaws.
3. Fix infrastructure/setup issues before tackling individual tests.
4. Iteratively re-run the tests after each systemic fix and track progress.
5. Summarize what was fixed and any remaining failed tests.

## Output Format
A summary list of the executed changes, and a short recap of test passing/failing status.

## Guardrails
- Assure that tests are not skipped just to bypass the issue. Real logic must be corrected.

## Standalone Mode
Supported, relies on the ability to run shell commands to execute the test suite repeatedly.
