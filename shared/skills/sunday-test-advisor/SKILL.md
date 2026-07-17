---
name: sunday-test-advisor
description: Audits a go-sunday api.yaml spec (Go-based Sunday variant with YAML-driven test codegen) for missing test scenarios per endpoint and interactively proposes YAML stubs the user can approve and generate. Specific to go-sunday — for TypeScript/Python/Java/C# Sunday tests, use the general audit process from `sunday-framework-patterns.md`'s API Test Coverage Matrix manually. Use when the user asks "what tests am I missing?", "audit my api.yaml coverage", "suggest more tests", or after adding a new endpoint to a go-sunday spec.
triggers:
  keywords: ["missing tests", "test coverage", "audit tests", "suggest tests", "test advisor", "what tests", "coverage gaps", "api.yaml audit", "go-sunday audit"]
  intentPatterns:
    - "what test scenarios am I missing in api.yaml"
    - "audit my api.yaml test coverage"
    - "suggest missing test cases for my go-sunday spec"
    - "check test coverage for my api.yaml"
standalone: true
---

## When To Use
Use when the user wants to audit a **go-sunday** `api.yaml` spec (Go-based Sunday variant with
YAML-driven test codegen) for missing scenarios, OR proactively after a new endpoint is added to that
spec.

Do NOT use for:
- TypeScript, Python, Java, or C# Sunday tests — those don't use `api.yaml`. The Coverage Matrix
  concept still applies (see the pattern doc reference below), but the audit is manual today.
  Building per-language advisors is deferred until each language's Sunday workflow needs one.
- Saturday E2E/UI tests — use `saturday-test-advisor` instead (structural adherence audit against
  `.feature` files, different mechanism entirely).
- Unit tests unrelated to `api.yaml` specs.

## Coverage Matrix (reference — full source is in the pattern doc)

The canonical Coverage Matrix lives in `docs/patterns/sunday-framework-patterns.md` under
"API Test Coverage Matrix" — that pattern-doc entry is the single source of truth for the baseline
scenarios per HTTP method, expected status codes, and when `auth_error`/`timeout_error` apply. Read
that section before running the audit; the local reproduction below exists only for operational
convenience so this skill's audit step doesn't require a second file read on every invocation.

Per HTTP method:

| Method   | Recommended Scenarios                                                      |
|----------|---------------------------------------------------------------------------|
| GET      | `happy_path`, `not_found`, `network_error`                                |
| GET list | `happy_path`, `network_error` (not_found doesn't apply to collections)    |
| POST     | `happy_path`, `server_error`, `validation_error`                          |
| PUT      | `happy_path`, `not_found`, `server_error`                                 |
| PATCH    | `happy_path`, `not_found`, `server_error`                                 |
| DELETE   | `happy_path`, `not_found`                                                 |

A GET endpoint that returns a slice (`[]Type`) is a "GET list". Detect by checking if `returns` starts with `[]`.

Additionally, if the spec has `auth: type: bearer|basic|apikey`, every endpoint should have an `auth_error` scenario.

Scenarios always worth adding:
- `timeout_error` — for slow network simulation (optional but recommended for SLA enforcement)

## Context To Load First
1. **Read `docs/patterns/sunday-framework-patterns.md`'s "API Test Coverage Matrix" section** — the
   canonical source of truth for the coverage baseline. If any discrepancy exists between the pattern
   doc and the reproduction above, the pattern doc wins.
2. Find the `api.yaml` spec — check these locations in order:
   - Path provided by user in their message
   - `api.yaml` in the current directory
   - `*/api.yaml` glob search (pick the first match)
3. Read the spec completely before analysis.
4. If no spec found, ask the user: "I couldn't find an api.yaml spec. What's the path?"

## Process

### Step 1 — Read and parse the spec
- Read `api.yaml` using the Read tool.
- Extract: `name`, `auth.type`, `endpoints[]`, `tests[]`.
- For each endpoint, note: `name`, `method`, `path`, `returns`, `body`.

### Step 2 — Build the coverage matrix
For each endpoint:
1. Determine the recommended scenario set from the Coverage Matrix above.
2. Add `auth_error` if `auth.type` is not `none`.
3. Collect which scenarios are already covered by tests (match `test.endpoint == endpoint.name`).
4. Compute gaps: `recommended - covered`.

### Step 3 — Present the audit report
Output a markdown table:

```
## go-sunday Test Coverage Audit: <spec name>

| Endpoint       | Method | Covered Scenarios              | Missing Scenarios                  |
|----------------|--------|--------------------------------|------------------------------------|
| GetAllPosts    | GET    | happy_path, network_error      | —                                  |
| GetPost        | GET    | happy_path, not_found          | network_error                      |
| CreatePost     | POST   | happy_path, server_error       | validation_error                   |
| UpdatePost     | PUT    | —                              | happy_path, not_found, server_error|
| DeletePost     | DELETE | happy_path                     | not_found                          |

**Coverage: 5 / 12 recommended scenarios covered (42%)**
```

Then ask: "Would you like me to generate YAML stubs for the missing scenarios and add them to your api.yaml?"

**PAUSE HERE** — wait for the user to respond before proceeding. Do not add anything to the spec yet.

### Step 4 — Generate YAML stubs (only after user approval)
The user must say "yes", "add them", "generate", or equivalent.

For each missing scenario, produce a YAML test case stub following this pattern:

```yaml
# ── Missing scenario stubs ──────────────────────────────────────────────────

  - name: <EndpointName>_<ScenarioPascal>
    endpoint: <EndpointName>
    scenario: <scenario_name>
    doc: <Auto-generated description>
    # params:          # uncomment and fill if endpoint has path params
    #   id: "1"
    # body:            # uncomment and fill if endpoint has a request body
    #   fieldName: value
    expect:
      status: <appropriate_default_status>
      # error: <error_type>  # for error scenarios: network | http | auth
```

Status and error defaults per scenario:
| Scenario          | status | error   |
|-------------------|--------|---------|
| happy_path        | 200    | —       |
| not_found         | 404    | http    |
| server_error      | 500    | http    |
| auth_error        | 401    | auth    |
| network_error     | —      | network |
| validation_error  | 422    | http    |
| timeout_error     | —      | timeout |

For POST happy_path, use status 201 instead of 200.

Auto-generate `name` by combining: `<EndpointName>_<ScenarioPascal>` — e.g., `GetPost_NetworkError`.
Auto-generate `doc` — e.g., "GetPost propagates a NetworkError when the adapter fails."

### Step 5 — Append to api.yaml
After showing the stubs and getting final approval:
1. Read the current api.yaml.
2. Find the last entry in `tests:` block.
3. Append the new stubs after the last test case (use Edit tool, not Write, to avoid overwriting).
4. Confirm: "Added N test cases to api.yaml. Run `sunday generate --spec api.yaml --out .` to regenerate the client and tests."

### Step 6 — Offer to run generation
Ask: "Shall I run `sunday generate` now to produce the test files?"
If yes: run `go run ./cmd/sunday/... generate --spec <path> --out <outdir>` via Bash tool.

## Output Format

### Audit Report (Step 3)
```markdown
## go-sunday Test Coverage Audit: Posts

| Endpoint    | Method | Covered                   | Missing              |
|-------------|--------|---------------------------|----------------------|
| GetPost     | GET    | happy_path, not_found     | network_error        |
| CreatePost  | POST   | happy_path                | server_error, validation_error |

**Coverage: 3 / 7 recommended scenarios covered (43%)**

Would you like me to generate YAML stubs for the 3 missing scenarios?
```

### YAML Stubs (Step 4)
Show as a fenced yaml block before appending. Never append without showing first.

## Guardrails
- NEVER modify api.yaml without showing the proposed stubs to the user first.
- NEVER run `sunday generate` without explicit user approval ("yes", "run it", "generate").
- If the spec has no `tests:` key at all, add it before appending stubs.
- Do NOT invent endpoint names or scenarios that don't exist in the matrix.
- Do NOT claim 100% coverage is required — the coverage matrix defines the minimum, not a mandate.
- Approval gates: any modification to api.yaml requires the user to say "add them", "yes", "approve", or equivalent.

## Standalone Mode
If the `sunday generate` CLI is not available (non-go-sunday project):
- Still produce the audit report and YAML stubs.
- Skip Step 6 (generation offer).
- Tell the user: "The stubs are ready — paste them into your api.yaml tests: block manually and run your generator."

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
