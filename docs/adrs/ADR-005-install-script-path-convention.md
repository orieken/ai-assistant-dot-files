# ADR-005: Use Explicit `.sh` Install Script Paths

## Status

Accepted

## Date

2026-08-15

## Context

The repository had two installer pairs at its root. The extensionless `install` and `uninstall`
scripts were an older, global-only symlink workflow. The current `install.sh` and `uninstall.sh`
scripts support global and project installations, dry runs, copy mode, platform filtering, AOS layers,
memory sync, and framework installation markers.

The extensionless scripts were not aliases. They maintained a separate list of root files, including
legacy blueprint prompts, and could therefore drift from the supported installation behavior. Current
documentation, Make targets, and installation tests already use the `.sh` scripts.

## Decision

Only `install.sh` and `uninstall.sh` are supported. The extensionless `install` and `uninstall` scripts
are retired and removed.

The unconsumed `.claude.md` compatibility source and the E2E and generic blueprint prompts installed
only by the retired path are also removed. The API framework blueprint remains live because the
`api-ingest` skill consumes it, and it now lives at
`docs/blueprints/API_FRAMEWORK_BLUEPRINT_PROMPT.md`.

## Consequences

- Installation and removal have one documented implementation each.
- Contributors must invoke `./install.sh` or `./uninstall.sh` explicitly.
- Legacy users invoking the extensionless paths receive a shell-level “file not found” error and must
  switch to the supported commands.
- Blueprint documentation no longer needs root-level compatibility filenames.

### Verification

This decision is **judgment-only** for new CI enforcement. Adding a check that forbids extensionless
installer files would wire a new architectural fitness function and requires separate Gate #7 approval.
Existing `scripts/test-install.sh`, `scripts/check-parity.sh`, and the Make targets exercise or reference
the supported `.sh` paths.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context
Engineering Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy
or adapt this file, please keep this attribution.*
