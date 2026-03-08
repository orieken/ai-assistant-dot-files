# Approval Gates

Anything irreversible requires explicit human approval. This is a hard gate. 

**Reset Condition:** Any edit to a pending artifact or code proposal immediately resets the gate. You must present the updated version and ask for approval again before proceeding.

## Required Gates

**Action:** Sending to Friday (QA agent)
**Irreversible because:** Initiates potentially long-running or expensive integration/E2E test pipelines.
**Gate:** User MUST explicitly say "approve" or "send" after viewing the passing unit tests.

**Action:** Creating a git commit
**Irreversible because:** Modifies version control history.
**Gate:** User MUST explicitly say "commit" after reviewing the final file diffs.

**Action:** Running database migrations
**Irreversible because:** Alters stateful data storage and schema.
**Gate:** User MUST explicitly say "migrate" after reviewing the proposed migration file.

**Action:** Posting to external APIs / Webhooks
**Irreversible because:** Mutates external system state or sends real notifications.
**Gate:** User MUST explicitly say "send" or "approve" after reviewing the exact payload.

**Action:** Dropping files outside `.claude/feature-workspace/`
**Irreversible because:** Mutates actual project source code.
**Gate:** User MUST confirm the file paths ("Proceed to write files?").

**Action:** Wiring a new fitness function into CI
**Irreversible because:** Changes the build pipeline for the entire team.
**Gate:** Architect MUST define and approve the function; user MUST approve before DevOps implements it.
