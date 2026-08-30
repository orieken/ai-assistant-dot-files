# Security Report: Login Handler

> Note: No `.claude/feature-workspace/analysis.md`, `architecture-notes.md`, `implementation-notes.md`, or `code-review-report.md` were found. This is expected — this review was invoked directly against a self-contained fixture file (`tests/agents/security-reviewer/input-vulnerable-code.ts`) rather than through the full feature delivery pipeline. Findings below are based entirely on direct code inspection of that file.

## Threat Model Summary
Trust boundaries crossed by this feature:
- User browser (unauthenticated) → `loginHandler` HTTP endpoint — this is the primary authentication boundary and the most sensitive trust crossing in the system.
- `loginHandler` → database (`db.query`) — user-supplied `username` crosses directly into a raw SQL string with no parameterization boundary.
- `loginHandler` → server logs (`console.log`) — sensitive credential material crosses from request memory into persistent/observable log storage.
- `loginHandler` → HTTP response body — a live third-party payment provider secret (`STRIPE_SECRET_KEY`) crosses from server configuration into every authenticated client response.

## Dependency Audit
No new dependencies introduced. This file only imports `express` types and a local `db` module; no `package.json` changes are in scope for this fixture. `pnpm audit` was not run since no dependency manifest changes accompany this file.

## STRIDE Analysis

### Spoofing
- No credential-strength or hashing controls are visible — `user.password !== password` compares in what appears to be plaintext (see Finding: Plaintext Password Comparison). If passwords are stored unhashed, any database read (including via the SQL injection below) trivially yields credentials usable to spoof any account.
- `generateToken(user)` is called but its implementation is not in this file; it cannot be verified here that issued tokens are bound to `exp`/`iss`/`aud` claims or that the session cookie (if any) is `httpOnly`/`secure`/`SameSite=Strict`. Flagged as unverifiable, not cleared.

### Tampering
- **Critical**: `username` is interpolated directly into a raw SQL string (line 9), allowing an attacker to tamper with the query itself (classic SQL injection), including bypassing authentication entirely (e.g. `' OR '1'='1' --`) or exfiltrating arbitrary table data.
- No input validation (e.g., Zod schema) is applied to `req.body` before `username`/`password` are used — there is no defense against malformed, oversized, or type-confused payloads (e.g., `username` as an array or object).

### Repudiation
- The only record of a login event is `console.log('Login success for ' + username + ' with password ' + password)` — this is not an audit trail, it is a credential leak (see Information Disclosure). There is no structured, tamper-evident audit log entry (user id, timestamp, source IP, outcome) and no OTel span marking the authentication event, so successful/failed logins cannot be reliably reconstructed after the fact.

### Information Disclosure
- **Critical**: A live Stripe secret key is hardcoded in source (line 4: `sk_live_...`) and additionally returned in the JSON response body to every client that logs in successfully (line 22: `stripeKey: STRIPE_SECRET_KEY`). This exposes a production payment-processing credential to any authenticated user's browser/network tooling, and to source control history.
- **Critical**: The plaintext password is written to server logs on every successful login (line 20). Logs are commonly shipped to third-party aggregators, retained long-term, and accessed by broader engineering staff than production databases — this is a direct credential disclosure channel.
- **High**: User enumeration — a nonexistent username returns `404 Username not found` while a wrong password returns `401 Incorrect password` (lines 12-18). An attacker can enumerate valid usernames by observing status codes/messages alone.
- SQL injection (above) is also an information-disclosure vector — it can be used to read arbitrary rows/columns from the `users` table or other tables depending on DB permissions.

### Denial of Service
- No rate limiting or lockout mechanism exists on this endpoint. Nothing prevents unlimited login attempts (brute force / credential stuffing) or repeated SQL-injection probing from a single client.
- No timeout is visible around `db.query(query)`; a slow or hung query has no bounded failure mode in this file (cannot rule out that `db` wraps this elsewhere, but it is not evident here).

### Elevation of Privilege
- Not directly applicable — this handler doesn't perform authorization checks on resources. However, the SQL injection vulnerability effectively grants elevation: an attacker can authenticate as any user, including administrators, without valid credentials.

## Findings

### CRITICAL — Hardcoded live Stripe secret key, leaked to clients
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` line 4, 22
**Threat category**: Information Disclosure
**Description**: A live-mode Stripe secret key (`sk_live_...`) is hardcoded directly in source and is returned in the JSON response of every successful login (`stripeKey: STRIPE_SECRET_KEY`). Any authenticated user — or anyone who intercepts the response — obtains a credential capable of moving real money and reading all payment/customer data in the connected Stripe account. The key is also permanently exposed in git history once committed.
**Fix applied**: Recommendation only (per test-fixture scope, source file was not patched). Required remediation: (1) remove the hardcoded key entirely and load it from an environment variable (`process.env.STRIPE_SECRET_KEY`) injected via a secrets manager/vault — never checked into source; (2) remove `stripeKey` from the response payload — a client-facing handler must never return a server-side secret key, full stop; if the client needs Stripe integration, use a publishable key or a short-lived client-side token (e.g., `PaymentIntent` client secret) generated server-side for that specific transaction; (3) rotate the leaked key immediately since it has already been committed to this file.
**Verification**: Grep the codebase and git history for `sk_live_` / `sk_test_` prefixes and fail CI if found. Confirm the login response schema (Zod) has no `stripeKey` (or any `sk_*` shaped string) field. Confirm `STRIPE_SECRET_KEY` is read only from `process.env` and is present in `.env.example` as a placeholder, not a value.

### CRITICAL — SQL Injection via unparameterized username
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` line 9
**Threat category**: Tampering / Elevation of Privilege
**Description**: `const query = \`SELECT * FROM users WHERE username = '${username}'\`` interpolates unsanitized user input directly into a SQL string. An attacker can supply a username like `' OR '1'='1' --` to bypass authentication, or use UNION-based payloads to exfiltrate arbitrary data from the database. This is exploitable pre-authentication by any anonymous caller — the highest-severity class of vulnerability in this framework.
**Fix applied**: Recommendation only. Required remediation: use the paved-road parameterized query path for the ORM/driver in use (e.g., `db.query('SELECT * FROM users WHERE username = $1', [username])` or the project's query builder with bound parameters) — never string-interpolate user input into SQL. Additionally validate `username` shape/length with a Zod schema before it reaches the data layer (defense in depth, not a substitute for parameterization).
**Verification**: QA/security should attempt classic injection payloads (`' OR '1'='1' --`, `'; DROP TABLE users; --`, UNION-based payloads) against the `username` field and confirm they are treated as literal string values with no behavioral difference from a normal invalid username. A static-analysis/lint rule (e.g., no template-literal SQL construction) should be added as a fitness function to prevent regression.

### CRITICAL — Plaintext password stored/compared, plus logged in cleartext
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` line 16, 20
**Threat category**: Spoofing / Information Disclosure
**Description**: `user.password !== password` compares the submitted password directly against a stored value with no indication of hashing (no `bcrypt.compare`/`argon2.verify` call). If passwords are stored in plaintext, a single database read (trivial here given the SQL injection) exposes every user's real password. Compounding this, line 20 writes the raw plaintext password to `console.log` on every successful login, guaranteeing it lands in log aggregation systems regardless of how it's stored.
**Fix applied**: Recommendation only. Required remediation: (1) store only a salted hash (bcrypt/argon2/scrypt) of the password, and compare via the library's constant-time `compare`/`verify` function — never `!==`/`===` on raw password strings, which is also vulnerable to timing attacks; (2) delete line 20 entirely — no log statement should ever contain a password field, hashed or not (use a low-cardinality, credential-free log line, e.g., `logger.info({ userId: user.id }, 'Login succeeded')`).
**Verification**: Confirm the `users` table/model never contains a plaintext `password` column (schema review). Confirm no log line in the auth path contains request-body password values (grep logs for injected test password strings during QA). Add a lint/fitness-function rule flagging any `console.log`/`logger.*` call whose arguments reference a variable named `password`.

### HIGH — User enumeration via differentiated error responses
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` lines 12-18
**Threat category**: Information Disclosure
**Description**: The handler returns `404 { error: 'Username not found' }` when the username doesn't exist, and `401 { error: 'Incorrect password' }` when it does but the password is wrong. This status-code and message difference lets an attacker enumerate valid usernames/emails at scale, which is a precursor to targeted credential-stuffing or phishing attacks.
**Fix applied**: Recommendation only. Required remediation: return an identical response (status code and body, e.g., `401 { error: 'Invalid username or password' }`) for both "user not found" and "wrong password" cases, with identical response-time characteristics (avoid short-circuiting in a way that makes the nonexistent-user path measurably faster — compare against a dummy hash to keep timing consistent).
**Verification**: QA should assert that submitting a known-bad username and a known-good-username-with-bad-password return byte-identical status codes and bodies, and that response latency does not differ by more than the framework's defined tolerance.

### HIGH — No rate limiting or lockout on authentication endpoint
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` (entire `loginHandler`)
**Threat category**: Denial of Service / Spoofing
**Description**: There is no attempt-count tracking, lockout state, or rate limiter around this handler. Combined with the plaintext-comparison and enumeration issues above, this endpoint is fully open to unlimited brute-force and credential-stuffing attempts from a single source.
**Fix applied**: Recommendation only — this is an architectural gap requiring a shared rate-limiting/lockout component (e.g., a `CircuitBreaker`-style or token-bucket middleware backed by a store that survives across requests), not a local code patch. Escalating rather than hand-rolling a retry/lockout loop in this file, per the Paved Road principle.
**Verification**: QA should confirm that after `MAX_LOGIN_ATTEMPTS` failures for a given identity/IP, subsequent requests receive a `429` with no user-identifying information in the body, and that the lockout persists for `LOCKOUT_DURATION_MINUTES`.

### MEDIUM — No input validation on request body
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` line 7
**Threat category**: Tampering
**Description**: `const { username, password } = req.body;` trusts the shape and type of the request body with no schema validation. There is no guarantee `username`/`password` are strings (as opposed to objects/arrays that could cause unexpected behavior downstream or in the ORM/driver layer), and no length/format bounds are enforced.
**Fix applied**: Recommendation only. Required remediation: validate `req.body` against a Zod schema (`z.object({ username: z.string().min(1).max(255), password: z.string().min(1).max(255) })`) at the top of the handler before any use, returning a generic `400` on failure.
**Verification**: QA should submit non-string `username`/`password` payloads (nested objects, arrays, null) and confirm the endpoint rejects them with a 400 rather than passing them further into the query/comparison logic.

### LOW — No audit trail / observability for authentication events
**Location**: `tests/agents/security-reviewer/input-vulnerable-code.ts` line 20
**Threat category**: Repudiation
**Description**: The only signal of a login event is an insecure `console.log`. There is no OTel span marking the authentication attempt/outcome and no structured, credential-free audit record, making it difficult to reconstruct who logged in (or attempted to) after an incident.
**Fix applied**: Recommendation only. Required remediation: replace the `console.log` with an OTel span/event emitted from the adapter/interceptor layer (never inside domain logic) containing only non-sensitive attributes (e.g., `userId`, `outcome`, `timestamp`) — explicitly excluding password and any other credential material.
**Verification**: Confirm a trace/span is emitted for both successful and failed login attempts and that its attributes contain no password, token, or other secret fields (automated span-attribute allowlist check).

## Files Modified
None — all findings are recommendations. Per the task scope, this is a one-off fixture review to validate a golden-file test harness; the source file `tests/agents/security-reviewer/input-vulnerable-code.ts` was intentionally left unpatched rather than edited, even though the Critical/High findings above would normally be fixed directly per this agent's standing rules.

## Security Checklist
- [ ] No secrets or credentials hardcoded — **FAILED**: live Stripe secret key hardcoded at line 4 and echoed in the response at line 22
- [ ] Auth tokens use framework auth providers (not manual header injection) — **UNVERIFIABLE**: `generateToken` implementation not present in this file
- [ ] User enumeration not possible via error message differences — **FAILED**: distinct 404/401 messages at lines 12-18
- [ ] OTel spans contain no PII or credential data — **FAILED**: no OTel span exists at all; the only event log contains a plaintext password
- [ ] External calls use CircuitBreaker or ExponentialBackoffStrategy — **N/A**: no external calls in this file beyond the DB query
- [ ] Input validated with Zod schemas on API boundaries — **FAILED**: `req.body` destructured with no schema validation
- [ ] Rate limiting in place on sensitive endpoints — **FAILED**: no rate limiting or lockout present
- [ ] Error responses sanitized (no stack traces, internal paths) — **PASSED**: no stack traces or internal paths are leaked in error bodies (the enumeration issue is tracked separately above)

## Notes for QA
- Test that SQL-injection payloads in the `username` field (e.g., `' OR '1'='1' --`) are treated as literal values and never bypass authentication or alter query results.
- Test that submitting a nonexistent username and submitting a valid username with a wrong password return identical status codes, bodies, and comparable response latency.
- Test that the login response body never contains a `stripeKey` field or any string matching the `sk_live_`/`sk_test_` Stripe key pattern.
- Test that application logs (stdout/log aggregator) never contain the literal password value submitted during a test login.
- Test that repeated failed login attempts against the same username/IP eventually receive a `429`-style lockout response once rate limiting is implemented.
- Test with non-string (`array`, `object`, `null`) `username`/`password` values in the request body and confirm a clean `400` rather than a server error or unexpected query behavior.

## Notes for Tech Writer
- Document that `STRIPE_SECRET_KEY` (and any other server-side secret) must be sourced exclusively from environment variables backed by a secrets vault, and must never appear in a client-facing API response — add this as an explicit rule in the API design/onboarding docs.
- Document the planned `MAX_LOGIN_ATTEMPTS` and `LOCKOUT_DURATION_MINUTES` environment variables once rate limiting is implemented, including their operational/security implications for on-call staff.
- Document the standard "generic auth error" convention (single message/status for both wrong-username and wrong-password cases) so future auth endpoints in this codebase follow the same anti-enumeration pattern by default.
</content>
