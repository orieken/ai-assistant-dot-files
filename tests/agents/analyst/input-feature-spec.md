# Feature: Password Reset via Email

## Problem
Users who forget their password have no way to regain access to their account. Support tickets for
manual password resets account for 12% of inbound support volume.

## Proposed Solution
Add a "Forgot password" flow: the user enters their email address, receives a time-limited reset link,
and sets a new password. The reset link must be single-use and expire after 30 minutes.

## Constraints
- Must not reveal whether a given email address is registered (avoid user enumeration)
- Reset tokens must expire after 30 minutes and be single-use
- Must send the email through the existing notification service (does not exist yet — new dependency)

## Out of Scope
- Changing the password complexity policy
- Multi-factor authentication
