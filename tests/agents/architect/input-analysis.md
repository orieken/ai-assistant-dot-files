# Feature Analysis: Password Reset via Email

## Summary
Adds a password reset flow that crosses from the Identity & Access context (user lookup, token issuance)
into a Notifications context that does not exist yet in the codebase. This is the first feature to need
outbound transactional email.

### Acceptance Criteria
- Given a registered email, when the user requests a reset, then a reset email is sent within 1 minute
  and the link expires after 30 minutes.
- Given an unregistered email, when a reset is requested, then the response is identical to a registered
  email's response (no user enumeration).

### Non-Functional Requirements
- Performance: reset email must be enqueued within 200ms of the request (actual send is async)
- Security: reset tokens are single-use, cryptographically random, and expire after 30 minutes
- Scaling: this is a low-volume flow (est. < 1000 requests/day), no special scaling concerns

## Technical Breakdown

### Bounded Context
- **Owning Context**: Identity & Access
- **Context Crossings**: Yes — Identity & Access now needs to call a Notifications capability that does
  not exist yet. This is a new integration point and needs an architectural decision on how the two
  contexts communicate (direct call vs. event-driven).

### New Dependencies
- An outbound email delivery provider (e.g., SES or SendGrid) — new external dependency, first of its kind
  in this codebase. Needs a retry/circuit-breaker strategy since it's a new external network call.

## Architectural Flags
Crossing a bounded context boundary for the first time with a new external dependency. Needs an
architect pass before implementation starts.
