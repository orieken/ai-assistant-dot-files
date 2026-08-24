# Enterprise Integration Patterns (Gregor Hohpe)

`architect.md` explicitly names Hohpe as an influence ("the enterprise integration patterns of Gregor
Hohpe — messaging between bounded contexts"), but the actual patterns were never written down. This repo
doesn't run literal message-queue infrastructure, but several of its own mechanisms are already
structurally the same pattern under a different name — those are cross-referenced below rather than
treated as unrelated theory.

## Message Channel

**Context**: A logical pipe connecting a producer to a consumer, decoupling them from needing to know
about each other directly.

**Structure**: The producer writes to the channel; the consumer reads from it. Neither needs a direct
reference to the other — the same decoupling `gang-of-four-patterns.md`'s Observer pattern describes at
the object level, applied at the system/process level.

## Content-Based Router

**Context**: Routes a message to a different destination based on its content, rather than every message
following the same fixed path.

**Structure**: A router inspects the message and decides where it goes next. This is what a Strategy
pattern (see `gang-of-four-patterns.md`) looks like at the message-routing level rather than the
object-method level.

## Message Translator

**Context**: Converts a message from one format/schema to another as it crosses a boundary between
systems that don't share a data model.

**Structure**: Structurally identical to an Anti-Corruption Layer (see `domain-driven-design.md`) — the
EIP name emphasizes the message-passing framing, the DDD name emphasizes protecting the domain model;
they're the same translation-at-a-boundary concept.

## Publish-Subscribe Channel

**Context**: One producer's message reaches every interested subscriber, not just one designated
consumer.

**Structure**: The system-level version of a Domain Event (see `domain-driven-design.md`) — this repo's
own `SpecReady`/`AnalysisComplete`/`PipelineComplete` events (`shared/DOMAIN_DICTIONARY.md` section 4)
are conceptually pub-sub: multiple downstream agents/skills can react to the same event without the
producer knowing who's listening.

## Dead Letter Channel

**Context**: Where messages go when they can't be processed successfully after retries are exhausted —
a place to preserve failed work for inspection rather than silently dropping it.

**Structure**: This repo's `.claude/feature-workspace/.history/` backups (used by `resume-pipeline`'s
rollback mode) serve a related purpose for pipeline artifacts: a failed or superseded artifact is
preserved rather than discarded, so a human can inspect what actually happened before deciding how to
recover.

## Correlation Identifier

**Context**: A shared ID attached to every message/event in one logical operation, so a system can
reconstruct the full sequence of what happened for one request even when it spans multiple independent
messages.

**Structure**: This repo's `pipeline-trace.json` (`shared/skills/pipeline-trace/SKILL.md`) is a direct,
concrete instance of this pattern applied to agent pipeline runs — every phase record for one delivery
correlates back to the same run, letting `pipeline-retrospective` and `agent-scorecard` reconstruct
cross-delivery timelines after the fact.

## Saga (Long-Running Process / Process Manager)

**Context**: Coordinates a multi-step business process that spans multiple independent steps, each of
which might fail, requiring compensating action for the steps that already succeeded rather than a single
atomic transaction across all of them.

**Structure**: `deliver-feature`'s own checkpointed pipeline is this pattern's shape: multiple sequential
agent steps, a persisted `pipeline-state.json` recording exactly how far the process got, and
`resume-pipeline`'s rollback mode as the compensating-action mechanism when a step needs to be undone and
retried — rather than one giant atomic operation that either fully succeeds or leaves no trace of partial
progress.

**Trade-offs**: A Saga trades atomicity for resilience — partial progress is visible and recoverable
rather than invisible until the whole thing either commits or rolls back. The cost is real complexity:
every step needs to define what "compensate" means if a later step fails, which a single-transaction
design never has to think about.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
