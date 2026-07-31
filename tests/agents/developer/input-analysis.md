# Feature: User Avatar Upload

## Summary
Allow users to upload a profile picture (JPEG/PNG, max 2 MB). The avatar is stored in S3, the URL is persisted in the `users` table, and displayed on the profile page.

## Acceptance Criteria
- Given a valid JPEG or PNG under 2 MB, the upload succeeds and the profile page shows the new avatar within 5 seconds.
- Given a file over 2 MB, the API returns 422 with a message explaining the size limit.
- Given an unsupported format (GIF, SVG, etc.), the API returns 422 with a format-error message.
- The S3 key for each user is deterministic: `avatars/<user-id>/<timestamp>.<ext>`.

## Affected Files
- `src/users/upload-handler.ts` (new)
- `src/users/user-repository.ts` (update `avatarUrl` field)
- `src/storage/s3-adapter.ts` (implement `StorageAdapter` interface)
- `db/migrations/0031_add_avatar_url.sql`

## New Dependencies
- `@aws-sdk/client-s3` (external — must hide behind `StorageAdapter` interface)

## Edge Cases
- Concurrent uploads from the same user should not cause a race condition on the `users` row.
- The old avatar key in S3 should be deleted when a new one is uploaded (prevent unbounded storage growth).
- A failed S3 upload must not update the database row.

## DevOps Tasks
- `AWS_S3_BUCKET` and `AWS_REGION` env vars needed in CI and production.
