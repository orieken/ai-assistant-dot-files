<!--
Template for devops-report.md. Consumed by the devops-engineer agent.
Structure defined here; contract in shared/contracts/devops-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# DevOps Report: [Feature Name]

## Files Created
- `.github/workflows/new-job.yml` — [what it does]

## Files Modified
- `.github/workflows/ci.yml` — [what changed]
- `.env.example` — Added: `NEW_VAR=example_value`

## New Environment Variables Required
| Variable | Description | Example | Secret? |
|---|---|---|---|
| `NEW_API_KEY` | API key for X service | `sk-...` | Yes |
| `FEATURE_TIMEOUT` | Timeout in ms for Y | `5000` | No |

## Migration Steps
- Run `alembic upgrade head` before deploying — or "None required"

## Deployment Notes
- [Anything ops needs to know when deploying this]
- [Any rollback procedure]

## Manual Steps Required
- [Things that CANNOT be automated and must be done manually, e.g. "Set NEW_API_KEY in production secrets vault"]
