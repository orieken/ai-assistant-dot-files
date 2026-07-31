# Eval Rubric: product-owner / input-feature-spec.md

- **ROI challenge is concrete**: the agent asks for (or supplies) a business metric that would prove the feature is worth 8 weeks — not just "what's the ROI?" but "how many users requested this, what's the engagement drop from leaving the app, what's the retention impact?"
- **Indefinite message storage is challenged**: storing messages indefinitely is flagged as a scope decision with real cost and compliance implications (GDPR right to erasure, storage growth) — the agent recommends a retention policy as part of MVP, not a later concern.
- **Existing alternatives are named**: the agent asks whether users who want to chat already use Slack/email/the app's notification system, and whether integrating with an existing tool was considered before building from scratch.
- **Minimal viable scope is proposed**: the agent proposes a reduced scope (e.g., defer read receipts or push notifications) with reasoning tied to the 8-week estimate — not just "build less" but "build less of X because Y is the actual value."
- **Build/buy is raised**: the agent explicitly raises the option of a third-party chat SDK (e.g., Twilio Conversations, Stream) vs. building in-house, with a brief cost/complexity tradeoff — since real-time messaging is infrastructure-heavy.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
