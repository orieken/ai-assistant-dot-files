# QA Report: User Avatar Upload

## Summary
PASS — 14/14 tests pass. Two edge cases were added beyond the original spec.

## Tests Added
- `src/users/upload-handler.spec.ts` — 8 unit tests (valid upload, size limit, format rejection, concurrent upload, old key deletion, S3 failure rollback, missing auth, empty body)
- `tests/integration/avatar-upload.spec.ts` — 6 integration tests (end-to-end upload flow, S3 mock)

## Coverage
- `upload-handler.ts`: 96%
- `s3-adapter.ts`: 88%
- `user-repository.ts`: 91% (avatar methods)

## Edge Cases Verified
- S3 upload failure does NOT update the database row (tested via injected S3 mock error)
- Concurrent uploads from the same user: second upload waits for the first to complete (advisory lock confirmed)
- Old S3 key deletion: confirmed via mock assertion that `DeleteObjectCommand` is called with the previous key

## Notes for Tech Writer
- `StorageAdapter` interface is new — needs API docs
- `AWS_S3_BUCKET` and `AWS_REGION` env vars are required — needs a setup section in the runbook
- The 30-day S3 lifecycle rule for `avatars/` prefix is a recommended ops step, not yet implemented
