---
name: run-tests
description: Executes the test suite for the project and verifies code coverage meets the 85% threshold.
triggers:
  keywords: ["run tests", "coverage", "execute tests", "test suite", "test"]
  intentPatterns: ["run the tests to verify", "check coverage", "make sure it passes tests"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill after making code modifications to ensure you did not break existing behavior.
Use this skill as part of the TDD `Red-Green-Refactor` cycle when checking if the tests pass.
Do NOT use when no executable source code or tests have been created or modified.

## Context To Load First
1. `package.json`, `pyproject.toml`, or relevant build configuration.
2. `ARCHITECTURE_RULES.md`

## Process
1. Determine test environment
   - What to do: Identify appropriate testing files and commands (e.g., `npm test`, `pytest`).
   - What to produce: Configuration of test run.
   - What to check: Verify testing dependencies are installed.
2. Execute tests
   - What to do: Run the `./.claude/skills/run-tests/run.sh` script to execute tests and generate coverage.
   - What to produce: Execution logs and test coverage reports.
   - What to check: Verify that tests pass. If tests fail, you MUST fix the tests or code before proceeding.
3. Validate coverage
   - What to do: Read the generated coverage report.
   - What to produce: Coverage metrics.
   - What to check: If code coverage is reported and is below 85%, you MUST write more tests to increase coverage.
   - When to pause: If complex dependencies make writing tests extremely difficult, pause for human input.

## Output Format
If there are errors or low coverage, output markdown with the results:
```markdown
# Test Execution Report

## Results
- **Status**: ✅ Passed / ❌ Failed
- **Failures**: [List of failing tests and reasons]

## Coverage
- **Current Coverage**: 82% (Target: >85%)
- **Missing Coverage**: `path/to/missing/file.ts` lines 12-25.

## Action Taken
- [Added missing unit test for error edge case]
```

## Guardrails
- You MUST NOT approve pull requests or mark features as done if coverage is below 85% or tests fail.
- Do NOT comment out failing tests to make them pass unless explicitly instructed by a human.

## Standalone Mode
If `run.sh` is unavailable, execute tests directly using available CLI test runners such as `npm test`, `pytest`, `go test`, or `cargo test`. Use regular expressions on the console output to extract code coverage numbers.
