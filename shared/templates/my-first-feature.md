# Feature: Rate-Limit Login Attempts

> **This is a tutorial spec**, not a real work item. It's pre-written and complete on purpose — run it
> straight through `/deliver-feature` to see the whole pipeline in action before writing your own spec with
> `/new-feature`. It's deliberately small but touches enough surface area (auth, a new data table, a new env
> var) to exercise most of the conditional pipeline agents: `architect`, `data-engineer`,
> `security-reviewer`, and `devops-engineer` should all activate, alongside the ones that always run.

## Summary
After 5 failed login attempts for the same account within 15 minutes, lock that account out from further
login attempts for 15 minutes and return a generic rate-limit response — without revealing whether the
account exists or how many attempts remain, so the lockout itself can't be used to enumerate valid usernames.

## Acceptance Criteria
- [ ] Given an account with 4 prior failed attempts in the last 15 minutes, when a 5th attempt fails, then the account is locked out for 15 minutes.
- [ ] Given a locked-out account, when any login attempt is made (correct or incorrect password), then the response is identical (status code and body) to a normal failed-login response — the lockout state itself must not be observable from the response.
- [ ] Given a locked-out account, when 15 minutes have elapsed, then a login attempt is evaluated normally again (the lockout clears; a subsequent success resets the failed-attempt counter to 0).
- [ ] Given a nonexistent username, when a login attempt is made, then the response is identical to a failed attempt against a real account (no user enumeration — this already applies today; this feature must not regress it).
- [ ] Given a successful login, when it occurs, then the failed-attempt counter for that account resets to 0 immediately.

## Out of Scope
- IP-based rate limiting (this is per-account only)
- CAPTCHA or other bot-detection mechanisms
- Notifying the account owner of lockout via email
- Admin tooling to manually clear a lockout

## User Stories
As a security-conscious operator, I want repeated failed login attempts against one account throttled, so
that credential-stuffing and brute-force attacks against individual accounts are meaningfully slowed down.

## Technical Notes
- Affected services/modules: the authentication/login handler; needs a new place to persist per-account
  failed-attempt counts and lockout expiry (new table or cache entry — this is a data model decision, which
  is exactly why this tutorial spec should trigger `data-engineer`).
- Related features/tickets: none (tutorial spec)
- Performance requirements: the rate-limit check must add no more than ~10ms p95 to the login request path.
- Security considerations: this feature exists specifically to close a security gap, but it introduces a new
  timing side-channel risk (a locked-out account's response must take the same time to compute as a normal
  failed attempt) — flag this explicitly for `security-reviewer`.

## Designs / References
- None — this is a backend-only, API-level feature with no UI change.

## Definition of Done
- [ ] All acceptance criteria verified manually
- [ ] Automated tests written and passing
- [ ] Documentation updated
- [ ] CI pipeline green
- [ ] Code reviewed and approved
