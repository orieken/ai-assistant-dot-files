---
name: review-pr
description: Reviews a pull request for code quality, security, and accessibility by coordinating the code-reviewer, security-reviewer, and accessibility-engineer agents.
triggers:
  keywords: ["review-pr", "review pr", "pr review", "pull request review"]
  intentPatterns: ["Review PR *", "Review pull request *", "Check this PR *", "/review-pr *"]
standalone: true
---

## When To Use
When the user wants a comprehensive review of a pull request. Accepts a PR number, URL, or branch name.

Do NOT use for reviewing a single file in isolation — use `/design-review` instead.
Do NOT use for reviewing a feature spec — use `/spec-writer` in review mode instead.

## Context To Load First
1. `CLAUDE.md` — project constraints and clean code rules
2. `ARCHITECTURE_RULES.md` — architectural guardrails
3. `DOMAIN_DICTIONARY.md` — ubiquitous language

## Process

1. **Resolve the PR** — accept a PR number, URL, or branch name. Use `gh pr view <input> --json number,title,body,baseRefName,headRefName,files` to get PR metadata.

2. **Fetch the diff** — run `gh pr diff <number>` to get the full changeset.

3. **Fetch PR comments** — run `gh api repos/{owner}/{repo}/pulls/{number}/comments` to see existing review feedback.

4. **Run code-reviewer analysis**:
   - Cyclomatic complexity of changed functions (must be < 7)
   - Function length (must be < 30 LOC)
   - SOLID principle violations
   - Naming convention adherence
   - Test coverage for new/changed code
   - File naming convention (`name.type.extension`)
   - Duplicated logic
   - Missing interfaces for external dependencies

5. **Run security-reviewer analysis** (if the diff touches auth, user input, API endpoints, tokens, or trust boundaries):
   - STRIDE threat assessment on changed components
   - Hardcoded secrets scan
   - Input validation review
   - Auth/authz changes review

6. **Run accessibility-engineer analysis** (if the diff touches UI files — HTML, Vue, React, templates):
   - Semantic HTML violations
   - Missing ARIA attributes
   - Keyboard navigation gaps
   - Color contrast concerns

7. **Produce the review summary** — display to the user and offer to post as a PR comment.

## Output Format

```markdown
# PR Review: #[number] — [title]

## Verdict
APPROVE | REQUEST CHANGES | COMMENT

## Code Quality
| Check | Result | Details |
|---|---|---|
| Cyclomatic complexity | PASS / FAIL | [specifics] |
| Function length | PASS / FAIL | [specifics] |
| SOLID principles | PASS / FAIL | [specifics] |
| Naming conventions | PASS / FAIL | [specifics] |
| Test coverage | PASS / FAIL | [specifics] |
| File naming | PASS / FAIL | [specifics] |

## Security (if applicable)
| Finding | Severity | File | Recommendation |
|---|---|---|---|
| [finding] | Critical / High / Medium / Low | [file:line] | [fix] |

## Accessibility (if applicable)
| Violation | Rule | File | Recommendation |
|---|---|---|---|
| [violation] | [WCAG rule] | [file:line] | [fix] |

## Inline Comments
[Specific line-by-line feedback with file paths and line numbers]

## Summary
[2-3 sentences: overall assessment, key concerns, recommendation]
```

## Guardrails
- Never post a PR review comment without explicit user approval ("post it", "submit review", "send").
- Never approve a PR that has cyclomatic complexity >= 7 in any changed function.
- Never approve a PR with Critical security findings.
- Always note what was NOT reviewable (e.g., "could not assess runtime behavior" or "no test files in diff").
- If `gh` CLI is not authenticated, inform the user and provide the review as local output only.

## Standalone Mode
Requires `gh` CLI for PR fetching. If unavailable, the user can paste the diff directly and the review runs on the pasted content without GitHub integration. The review itself is pure analysis — no external calls required.
