# Eval Rubric: devops-engineer / input-docs-report.md

- **Env vars added as CI secrets with placeholder values only**: the devops-report shows `AWS_S3_BUCKET` and `AWS_REGION` added to the CI pipeline config using placeholder values (`<your-bucket>`, `<your-region>`) — never real values, with a note that real values must be set in the CI secrets store.
- **S3 bucket and IAM policy steps are in the runbook**: the devops-report references or updates `docs/runbooks/avatar-storage.md` with the bucket creation command/Terraform snippet and the IAM policy — not left as a verbal note.
- **Lifecycle rule is wired or tracked**: either a Terraform/CloudFormation snippet for the 30-day lifecycle rule is included, or a tracked task is created — the agent does not silently drop the gap.
- **GitHub Actions / CI config follows existing pinning pattern**: any new step or secret reference follows the same action-version pinning convention as the existing CI file (SHA pins or consistent tag format), per iac-conventions.md.
- **Health check or readiness probe is validated**: if the avatar endpoint is new, the agent checks or flags that existing health checks still cover the service, or adds one if missing.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
