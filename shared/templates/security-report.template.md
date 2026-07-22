<!--
Template for security-report.md. Consumed by the security-reviewer agent.
Structure defined here; contract in shared/contracts/security-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Security Report: [Feature Name]

## Threat Model Summary
Trust boundaries crossed by this feature:
- [e.g., "User browser → API gateway (unauthenticated)"]
- [e.g., "API gateway → auth service (service-to-service with mTLS)"]

## Dependency Audit
- [pnpm audit output or "No new dependencies introduced"]

## STRIDE Analysis

### Spoofing
- [Finding or "No issues identified"]

### Tampering
- [Finding or "No issues identified"]

### Repudiation
- [Finding or "No issues identified"]

### Information Disclosure
- [Finding or "No issues identified"]

### Denial of Service
- [Finding or "No issues identified"]

### Elevation of Privilege
- [Finding or "No issues identified"]

## Findings

### [CRITICAL/HIGH/MEDIUM/LOW] — [Short title]
**Location**: `path/to/file.ts` line N
**Threat category**: [STRIDE category]
**Description**: [What the vulnerability is and how it could be exploited]
**Fix applied**: [What was changed — or "Recommendation only" if not auto-fixed]
**Verification**: [How QA can verify the fix is working]

## Files Modified
- `path/to/file.ts` — [What security fix was applied]
— or "None — all findings were recommendations"

## Security Checklist
- [ ] No secrets or credentials hardcoded
- [ ] Auth tokens use framework auth providers (not manual header injection)
- [ ] User enumeration not possible via error message differences
- [ ] OTel spans contain no PII or credential data
- [ ] External calls use CircuitBreaker or ExponentialBackoffStrategy
- [ ] Input validated with Zod schemas on API boundaries
- [ ] Rate limiting in place on sensitive endpoints
- [ ] Error responses sanitized (no stack traces, internal paths)

## Notes for QA
- [Security scenarios QA should include in test coverage]
- [e.g., "Test that attempting login with nonexistent username returns same error as wrong password"]
- [e.g., "Test that the lockout endpoint returns 429 with no user-identifying information in the body"]

## Notes for Tech Writer
- [Security behavior that should be documented for operators]
- [e.g., "Document the MAX_LOGIN_ATTEMPTS and LOCKOUT_DURATION_MINUTES env vars and their security implications"]
