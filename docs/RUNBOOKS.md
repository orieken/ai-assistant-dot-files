# AI Problem-Solving Runbooks

This document outlines the standard operational procedures (SOPs) and debugging workflows our AI assistants natively understand to help you repair broken states methodically.

## 1. Environment Remediation (`debug-environment`)

When local tools fail to execute, commands go missing, or shell configurations fracture (like an invalid `PATH`), do not hack one-off localized workarounds. Trigger the **`debug-environment`** skill.

- **Trigger Keywords:** `"fix PATH"`, `/debug-environment`, `"fix my configuration"`
- **The Workflow:**
    1. The AI will immediately target and inspect RC configuration (`.zshrc`, `.bashrc`, `.bash_profile`, `package.json`, `.env`).
    2. Instead of attempting destructive blind exports, it verifies current environment states cleanly via diagnostic echo prints.
    3. It will formulate a specific modification with before-and-after verification commands natively built-in.

> [!CAUTION]
> Agents are instructed *not* to commit global environment profile changes without reviewing the explicit consequences with you.

---

## 2. Test Suite Remediation (`debug-tests`)

When you encounter cascading or massive test failures, do **not** attempt to instruct the AI to patch individual assertions one by one. The failure is likely a symptom of lost environmental context, a dropped mock database, or misconfigured runners. 

Use the **`debug-tests`** skill to execute the holistic infrastructure debugging process.

- **Trigger Keywords:** `"debug test failures"`, `/debug-tests`, `"why are my tests failing?"`
- **The Workflow:**
    1. The AI will pull the entire test suite run string and capture all logs simultaneously to track patterns.
    2. It prioritizes infrastructure patches, runner misconfigurations, and missing global dependencies.
    3. Once the environment is repaired, it will iteratively trigger localized test runner scripts (`npm test`) to zero in on explicit logical bugs until resolved.

> [!IMPORTANT]
> The AI is strictly bound by guardrails preventing it from omitting tests or commenting out assertions to artifically manufacture a passing suite. The logic must be fixed.
