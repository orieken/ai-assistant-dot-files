# 🤖 Feature Team — Claude Code Multi-Agent Setup

A fully automated feature delivery pipeline for Claude Code. Define a feature in markdown, and your AI team of Analyst, Developer, QA Engineer, Tech Writer, and DevOps Engineer will implement it end-to-end.

## Directory Structure

```
your-project/
├── CLAUDE.md                          ← Orchestrator instructions (copy + customize)
├── features/
│   ├── TEMPLATE.md                    ← Copy this when writing new features
│   └── your-feature.md                ← Your feature specs live here
└── .claude/
    ├── agents/
    │   ├── analyst.md                 ← Breaks down features into tasks
    │   ├── developer.md               ← Implements code
    │   ├── qa-engineer.md             ← Writes and runs tests
    │   ├── tech-writer.md             ← Updates documentation
    │   └── devops-engineer.md         ← CI/CD and infrastructure
    └── skills/
        └── deliver-feature/
            └── SKILL.md               ← The /deliver-feature slash command
```

## Installation

1. **Copy files into your project:**
   ```bash
   # From this repo root
   cp CLAUDE.md /path/to/your-project/
   cp -r .claude /path/to/your-project/
   cp -r features /path/to/your-project/
   ```

2. **Update `CLAUDE.md`** with your actual tech stack:
   - Replace `[SET IN PROJECT CLAUDE.md]` with your language/framework
   - Update the testing framework (pytest / Playwright / Cucumber)
   - Update the CI system (GitHub Actions / Jenkins)

3. **Make sure Claude Code is installed:**
   ```bash
   npm install -g @anthropic-ai/claude-code
   ```

4. **Start Claude Code in your project:**
   ```bash
   cd your-project
   claude
   ```

## Usage

### Writing a Feature Spec

Copy `features/TEMPLATE.md` to `features/my-feature-name.md` and fill it out.

**Good spec:**
```markdown
# Feature: User Avatar Upload

## Summary
Allow users to upload a profile photo up to 5MB (JPEG/PNG).
The avatar displays on their profile page and in comment threads.

## Acceptance Criteria
- [ ] Given I'm logged in, when I upload a valid JPEG/PNG under 5MB, then it saves and displays
- [ ] Given a file over 5MB, when I try to upload, then I see "File too large" error
- [ ] Given a non-image file, when I try to upload, then I see "Invalid file type" error
- [ ] Given I have an avatar, when I view my profile, then the avatar shows instead of placeholder
```

### Delivering a Feature

In Claude Code, run:
```
/deliver-feature features/user-avatar-upload.md
```

Or just tell Claude:
```
Implement the feature in features/user-avatar-upload.md
```

Claude will:
1. 🔍 **Analyst** reads your spec, explores the codebase, writes a technical breakdown
2. ⏸️ **You approve** the analysis before code is written
3. 💻 **Developer** implements the feature
4. 🧪 **QA Engineer** writes and runs tests, fixes failures
5. 📝 **Tech Writer** updates README, CHANGELOG, docs
6. 🚀 **DevOps** updates CI/CD config and env var documentation
7. ✅ **Summary** delivered to you with manual steps required

### Manual Invocation

You can invoke individual agents directly:
```
Use the analyst subagent to analyze features/my-feature.md
Use the qa-engineer subagent to write tests for the changes in src/auth/
Use the tech-writer subagent to update docs after implementing the pagination feature
```

## Customizing Agents

Each agent is a plain markdown file in `.claude/agents/`. Edit the system prompt to match your team's conventions:

- **Add your coding standards** to `developer.md`
- **Add your test patterns** to `qa-engineer.md` (e.g., specific pytest fixtures, Cucumber step conventions)
- **Add your doc structure** to `tech-writer.md` (e.g., specific ADR format, wiki instead of markdown)
- **Add your CI system** to `devops-engineer.md` (e.g., Jenkins vs GitHub Actions specifics)

## Workspace Artifacts

Each feature run produces artifacts in `.claude/feature-workspace/`:

| File | Produced By | Contains |
|------|------------|----------|
| `analysis.md` | Analyst | Full technical breakdown |
| `implementation-notes.md` | Developer | What was built and key decisions |
| `qa-report.md` | QA Engineer | Test results, bugs found |
| `docs-report.md` | Tech Writer | What docs were updated |
| `devops-report.md` | DevOps | CI changes, env vars, manual steps |
| `delivery-summary.md` | Orchestrator | Final status and summary |

Add `.claude/feature-workspace/` to `.gitignore` or commit it — your choice.

## Tips

- **Be specific in feature specs.** Vague specs lead to vague analysis and vague code.
- **The analyst checkpoint is your safeguard.** Read `analysis.md` carefully before approving — this is when you catch misunderstandings cheaply.
- **Agents respect existing patterns.** They're instructed to read the codebase first. The more consistent your existing code, the better results you'll get.
- **Run it on a branch.** Since the developer uses worktree isolation, changes are isolated by default.
- **Token usage.** This pipeline uses significant tokens per run — one feature delivery can use 100K-500K tokens depending on codebase size. Use Claude Max or API with sufficient limits.

## Troubleshooting

**Agent not being invoked automatically?**
Claude uses the `description` field to decide when to use an agent. If it's not triggering, be explicit: "Use the `qa-engineer` subagent to..."

**Agent ignoring the workspace files?**
Make sure `.claude/feature-workspace/` exists. The agents look for specific files — if a previous step didn't produce its output file, the next agent may proceed without it.

**Tests failing in CI but passing locally?**
Add environment-specific notes to your `CLAUDE.md` about the CI environment, so the devops agent knows what to configure.
