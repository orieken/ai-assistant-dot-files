# Docs Report: User Avatar Upload

## Summary
Documentation updated. Two ops gaps flagged for DevOps.

## Changes Made
- `README.md`: added Avatar Upload section with API reference (`POST /users/:id/avatar`)
- `CHANGELOG.md`: entry added for v1.8.0
- `docs/runbooks/avatar-storage.md`: created — describes S3 bucket setup, IAM policy, lifecycle rule

## Env Vars Required
- `AWS_S3_BUCKET` — the S3 bucket name for avatar storage
- `AWS_REGION` — AWS region (e.g., `us-east-1`)
- Both must be set in CI and production environments

## Gaps for DevOps
- S3 bucket must be created with private ACL and a 30-day lifecycle rule for `avatars/` prefix
- IAM policy for the app's service role must allow `s3:PutObject`, `s3:DeleteObject`, `s3:GetObject` on `arn:aws:s3:::<bucket>/avatars/*`
- CI pipeline does not yet include these env vars — they need to be added as secrets

## ADR
No new ADR required — use of S3 via an interface adapter follows existing architecture decisions.
