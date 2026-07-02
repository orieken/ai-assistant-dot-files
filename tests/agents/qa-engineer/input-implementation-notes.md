# Implementation Notes: Password Reset via Email

## Files Created
- `src/auth/reset-password.service.ts` — generates single-use reset tokens (30-minute expiry) and validates them
- `src/auth/reset-password.controller.ts` — `POST /auth/reset-password/request` and `POST /auth/reset-password/confirm` endpoints
- `src/notifications/email-provider.adapter.ts` — thin adapter around the SES SDK, wrapped in a CircuitBreaker

## Files Modified
- `src/auth/user.repository.ts` — added `findByEmail` (returns `null` on no match, same shape either way)

## Interface Design
```ts
interface ResetToken {
  tokenHash: string;
  userId: string;
  expiresAt: Date;
  usedAt: Date | null;
}
```

## Self-Review Checklist
- [x] Every public method has an intention-revealing name
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new functions
- [x] No primitive obsession
- [x] No feature envy
- [x] No magic numbers or strings
- [x] Dependency direction verified

## Simple Design Verification
1. **Passes all tests**: yes
2. **Reveals intention**: yes
3. **No duplication**: yes
4. **Fewest elements**: yes

## Notes for QA
- Test that requesting a reset for a nonexistent email returns the exact same response shape and status
  code as a valid one (no user enumeration via timing or payload differences)
- Test that a reset token expires after 30 minutes and a confirm attempt after expiry is rejected
- Test that a reset token can only be used once (second confirm attempt with the same token fails)

## Notes for DevOps
- New env vars: `SES_ACCESS_KEY`, `SES_SECRET_KEY`, `RESET_TOKEN_TTL_MINUTES` (default 30)
