# Acceptance Criteria: Email Uniqueness Validation

## Feature
When a user attempts to register with an email address already in use, the system must:
1. Return HTTP 409 Conflict with body `{ "error": "EMAIL_ALREADY_REGISTERED" }`.
2. NOT reveal whether the email is in the system if the request comes from a non-authenticated context (anti-enumeration).
3. Enforce the check at the domain layer, not only at the database constraint level.

## Notes
- The database already has a UNIQUE constraint on `users.email`.
- The registration endpoint is `POST /auth/register`.
- Domain entity: `User` in `src/domain/user.entity.ts`.
