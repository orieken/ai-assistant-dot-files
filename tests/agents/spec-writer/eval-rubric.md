# Eval Rubric: spec-writer / input-user-request.md

- **System preference is raised as a question**: the agent asks (or assumes and documents) whether the app should respect `prefers-color-scheme` OS/browser preference automatically, independent of the user toggle — this is a non-obvious product decision the request leaves ambiguous.
- **Persistence mechanism is specified**: the spec states WHERE the preference is stored (`localStorage`, user DB column, or cookie) and WHY — not just "remember the preference."
- **Acceptance criteria are given/when/then**: each AC is in observable, behavior-level language (e.g., "Given a user who set dark mode and refreshes the page, the app loads in dark mode") — not implementation language ("The ThemeContext stores the preference").
- **Scope is explicit**: something is explicitly put Out of Scope (e.g., per-page theme overrides, theme for email templates, admin-forced themes) — the spec writer challenged the user to narrow scope rather than leaving it open.
- **Readiness critique is run**: the agent explicitly runs its readiness checklist (or equivalent) and either declares the spec ready or lists what's still missing before it can enter the delivery pipeline.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
