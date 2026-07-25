The hardest part of AI memory is not saving things.

It is not saving everything.

In `ai-assistant-dot-files`, memory follows a promotion lifecycle:

Capture -> Candidate -> Audit -> Approve -> Index -> Retrieve -> Expire

The key design choice: agents do not write directly to durable memory. They produce Candidate Records with
source, evidence, tags, and an expiration condition. Then those candidates are audited before a human
approves them.

That makes rejection a healthy outcome.

A one-off lesson should not become permanent doctrine. A duplicate should merge into the existing item. A
speculative idea can stay a note until the evidence is real.

Full post: https://dev.to/orieken/memory-engineering-is-a-promotion-pipeline-not-a-pile-of-notes-3eee

