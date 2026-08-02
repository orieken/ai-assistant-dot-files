# Approval Gates

**Irreversible actions require explicit human approval. Any edit or change to the pending artifact resets the gate.**

As of v3.3, gates may optionally be delegated to the AOS policy evaluator
(`shared/orchestration/policy-evaluator.md`). Each gate below is annotated with its policy
eligibility. Gates marked **Always Human** are never delegated regardless of any policy file.
Gates marked **Policy-Eligible** may be auto-approved when a matching policy exists in
`.claude/policies/` — but only when no `require-human` policy also matches (which always wins).
Policy evaluation is **strictly opt-in**: the absence of any policy file means all gates continue
to require human confirmation as in prior versions.

### 1. Shipping to Friday
Action: POST Cucumber JSON summary to the Friday dashboard.
Irreversible because: It updates external reporting metrics.
Gate: user must say "ship" or "yes" to the delivery summary prompt.
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: No — Always Human.**
Reason: external shared reporting metrics cannot be rolled back via git revert; human intent is
always required before mutating an external system's state.

### 2. Creating a Git Commit
Action: Creating a commit on the active branch.
Irreversible because: It alters repository history.
Gate: user must say "commit" or "approve commit".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: Yes (Tier A).**
Policy gate ID: `git-commit`. Auto-approve is permitted when: diff is below a configured line
threshold, code-reviewer returned APPROVED, all tests pass, and no security/auth paths are
touched. See `shared/policies/examples/auto-approve-refactor.policy.yaml` for a reference policy.

### 3. Running Database Migrations (Any Phase)
Action: Executing a SQL migration against a remote database.
Irreversible because: Modifies stateful infrastructure data.
Gate: user must say "run migration" or "execute phase X".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: No — Always Human.**
Reason: modifies live stateful infrastructure; even expand-phase migrations can corrupt data on
partially applied runs. Infrastructure blast radius too high for automation.

### 4. Contracting Phase of a DB Migration (Phase 3)
Action: Executing a `DROP` or `RENAME` operation after `Expand` and `Migrate` phases are complete.
Irreversible because: Data loss risk.
Gate: user must say "confirm contract phase".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: No — Always Human.**
Reason: data destruction is irreversible; no automation tier safely encompasses DROP/RENAME.

### 5. Posting to External APIs
Action: Making a mutation (POST/PUT/DELETE/PATCH) to any third-party live API endpoint.
Irreversible because: External side-effects.
Gate: user must say "send" or "approve request".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: No — Always Human.**
Reason: third-party mutations have no guaranteed rollback path and affect systems outside this
repository's control boundary.

### 6. Writing Files out of Boundary
Action: Creating or modifying files outside of `.claude/feature-workspace/` or proper source directories.
Irreversible because: Potentially breaks system structure or config.
Gate: user must say "approve file write".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: Yes (Tier A).**
Policy gate ID: `out-of-boundary-write`. Auto-approve is permitted when all target paths match
a project-configured `allowedPaths` whitelist and no security/auth paths are involved.

### 7. Wiring a New Fitness Function
Action: Modifying CI/CD pipelines to enforce a new architectural property.
Irreversible because: Breaks builds if poorly formulated.
Gate: user must say "approve fitness function" or "add to CI".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: Yes (Tier A).**
Policy gate ID: `fitness-function-wiring`. Auto-approve is permitted when: a CI dry-run passes,
the wired function is a new test (not modifying existing checks), and zero security criticals
exist. See `shared/policies/examples/auto-approve-test-additions.policy.yaml`.

### 8. Deploying to Environment
Action: Triggering a deployment of code.
Irreversible because: Could cause downtime.
Gate: user must say "deploy".
Reset condition: any edit to the pending artifact resets the gate.
**Policy-eligible: No — Always Human.**
Reason: deployment failures can cause production downtime; the risk profile requires a human
decision point regardless of prior stage verdicts.

---

## Policy-Based Gate Type (v3.3+)

A policy-based gate operates identically to a human gate except the human prompt is replaced by
the policy evaluator's decision when a matching policy exists and its condition is met. The
evaluator emits a `policy.evaluated` telemetry event for every decision — including decisions
that fall through to `require-human`. There are no silent auto-approvals.

To opt in: place `.policy.yaml` files in `.claude/policies/` in your project.
To opt out globally: set `policiesEnabled: false` in `.claude/delivery-policy.yaml`.

See `shared/orchestration/policy-evaluator.md` and `docs/aos/policy-authoring-guide.md`.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
