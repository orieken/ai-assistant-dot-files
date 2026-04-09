# ADR-001: Adopt RAG-friendly Documentation Structure

## Status

Accepted

## Date

2026-04-08

## Context

The `/docs/` directory contains a flat list of files with no organizational convention. Markdown files, zip archives, prompt catalogs, and agent configurations all sit at the same level. This creates several problems:

- **Agent context loading is inefficient.** Agents must scan the entire docs directory to find relevant material. There is no way to load a focused subset of documentation for a specific task (e.g., only architecture decisions, only past feature artifacts).
- **No convention for feature delivery artifacts.** When the delivery pipeline produces analysis, review reports, and architecture notes, there is no designated location to persist them. Artifacts are either lost or scattered.
- **Onboarding is manual.** New contributors have no structured path through the documentation. The flat layout gives no indication of reading order or document purpose.
- **Pattern knowledge is implicit.** Design patterns and framework conventions exist in agent prompts and code comments but are not documented in a referenceable format.

## Decision

Restructure `/docs/` into a RAG-friendly knowledge base with these subdirectories:

- **`/docs/features/`** -- Each delivered feature gets a subdirectory containing all pipeline artifacts (analysis, architecture notes, code review reports, security reports, QA reports, delivery summaries). This creates a permanent, searchable record of every feature delivery.
- **`/docs/patterns/`** -- Reusable pattern documentation for design patterns, Saturday Framework patterns, Sunday Framework patterns, and Clean Architecture layer conventions. Each pattern gets its own file.
- **`/docs/adrs/`** -- Architecture Decision Records with sequential numbering (ADR-001, ADR-002, etc.). Captures context, decision, and consequences for every significant architectural choice.
- **`/docs/runbooks/`** -- Operational runbooks covering installation, troubleshooting, agent reference, and pipeline execution.

Existing files in `/docs/` remain in place. This restructure adds new directories without moving or deleting current content.

A `README.md` in each subdirectory explains the convention and provides navigation for both agents and humans.

## Consequences

### What becomes easier

- Agents load only the subdirectory relevant to their current task, reducing noise and improving context precision.
- Feature delivery artifacts accumulate in a predictable location, creating an auditable history.
- New contributors follow a clear documentation structure with README files guiding navigation.
- Architecture decisions are tracked and referenceable, preventing repeated discussions of settled questions.
- Pattern documentation is centralized and available for both agent context loading and human reference.

### What becomes harder

- Contributors must follow the directory convention when adding new documentation. Ad-hoc file placement in the root of `/docs/` is no longer acceptable for new content.
- The initial restructure adds files without migrating existing content, creating a transitional period where some documentation lives in the root and some in subdirectories.
- Agents and pipeline stages must be updated to write artifacts to the correct feature directory rather than to arbitrary locations.
