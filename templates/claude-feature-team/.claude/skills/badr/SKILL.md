---
name: badr
description: Business Architecture Decision Records. Documents the business context, opportunity cost, and ROI of technical choices.
triggers:
  keywords: ["BADR", "business decision", "ROI", "cost-benefit"]
  intentPatterns: ["Write a BADR for *", "Document the business decision to *", "Why did we decide to build * from a business perspective?"]
standalone: true
---

## When To Use
Triggered for major architectural shifts, "Build vs. Buy" decisions, or when taking on massive technical debt intentionally for time-to-market. Works alongside the `product-owner` agent.

## Context To Load First
1. Any relevant Strategy Docs or Product Specs.
2. Discussions with the `product-owner` or `architect`.

## Process
1. Determine next BADR number (`docs/badr/`).
2. Ask for the Business Context: "What market force or business goal requires this technical decision?"
3. Ask for the Opportunity Cost: "If we commit engineering hours to this, what are we NOT building instead?"
4. Ask for the "Buy/Adopt" alternative: "Did we consider a SaaS or open-source alternative? Why was it rejected?"
5. Draft the BADR and seek user approval.
6. Write to `docs/badr/BADR-[NNN]-[kebab-title].md`.

## Output Format
```markdown
# BADR-[NNN]: [Title]

Date: YYYY-MM-DD
Status: Accepted
Deciders: [from git config user.name]
Business Sponsor: [Who asked for this]

## Executive Summary
[One paragraph explaining the *business* value of the technical decision]

## Business Context & Market Force
[Why doing nothing is not an option. E.g., "Competitor X launched a real-time feed, our users are churning."]

## Decision
[What we chose to do technically. E.g., "Build a custom WebSocket infrastructure instead of buying Pusher."]

## Opportunity Cost
[What engineering work will be delayed to accommodate this]

## "Buy vs. Build" Analysis
| Option | Cost | Time to Market | Customizability | Why Chosen/Rejected? |
|---|---|---|---|---|
| Buy Pusher | $X/mo | 1 week | Low | Rejected: Cost scales too high |
| Build Custom | $Y dev cost | 6 weeks | High | Chosen: Cheaper at scale |

## Expected ROI Metrics
[How we will measure if this was a good business decision 6 months from now]
```

## Guardrails
- BADRs focus on the *why* (market, cost, time), ADRs focus on the *how* (patterns, frameworks, boundaries).
- Never evaluate a "Build vs. Buy" scenario without considering the total cost of ownership (maintenance, on-call).
