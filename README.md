# AI Assistant Dot Files

This repository contains global configuration dotfiles for various AI coding assistants (Claude, Cursor, Windsurf, GitHub Copilot, Google Gemini Antigravity, and OpenAI/Codex).

These instructions configure the AI agents to behave as expert software craftsmen, adhering to the principles of Clean Architecture, Test-Driven Development (TDD), and domain-driven design within the Saturday Framework ecosystem.

## Core Craftsmanship Principles Enforced

By using these dotfiles, all AI assistants are instructed to prioritize:

- **TDD/BDD (Kent Beck)**: Driving design through tests and insisting on the Red-Green-Refactor cycle.
- **Clean Code & SOLID (Uncle Bob)**: Enforcing strict cyclomatic complexity limits (< 7), short function lengths (< 30 LOC), and adherence to SOLID principles.
- **Evolutionary Architecture (Neal Ford & Martin Fowler)**: Building highly cohesive, loosely coupled systems with proven enterprise patterns.
- **YAGNI & KISS**: Avoiding over-engineering and premature abstractions.
- **The Boy Scout Rule**: Proactively cleaning up minor technical debt and formatting issues when modifying files.
- **Security & Observability**: Never hardcoding secrets, emitting OpenTelemetry traces, and maintaining documentation parity.

## Files and Supported Agents

| File Path | Supported Agent(s) | Description |
| :--- | :--- | :--- |
| **`.claude/CLAUDE.md`** | Anthropic Claude | Global instructions for Claude outlining framework patterns and design rules. |
| **`.cursor/rules/global.mdc`** | Cursor IDE (AI Editor) | System prompt overrides for the Cursor IDE's internal AI models. |
| **`.windsurfrules`** | Windsurf IDE (AI Editor) | System prompt overrides for the Windsurf IDE's internal AI models. |
| **`.openai.md`** | OpenAI / Codex | Global instructions for OpenAI models and tools leveraging Codex. |
| **`.github/copilot-instructions.md`** | GitHub Copilot | Custom instructions injected into Copilot Chat and inline completions. |
| **`.gemini/antigravity/instructions.md`** | Gemini Antigravity | Instructions specifically tailored for Google's agentic coding assistant within this ecosystem. |
| **`E2E_FRAMEWORK_BLUEPRINT_PROMPT.md`** | General AI Usage | A core prompt establishing the framework blueprint for E2E UI automation agents. |
| **`API_FRAMEWORK_BLUEPRINT_PROMPT.md`** | General AI Usage | A core prompt establishing the framework blueprint for API testing agents. |
| **`BLUEPRINT_GENERATOR_PROMPT.md`** | Prompt Engineering | A prompt template to use with an AI (like Claude or ChatGPT) to automatically generate new blueprint prompts from existing codebases. |

## Usage (Installation)

We recommend symlinking these files from this repository directly into your home directory. This allows you to simply run `git pull` in this folder to immediately update all your agent instructions globally.

To set up the global symlinks automatically, run the install script:

```bash
./install
```

*(Note: The script will back up any existing files to `.bak` before creating the symlinks).*

To remove these files and restore any `.bak` backups automatically, simply run:
```bash
./uninstall
```

### Project-Specific Workarounds (GitHub Copilot)
Some tools, like **GitHub Copilot**, do not currently support a true user-level global instruction file. Copilot expects `.github/copilot-instructions.md` to be present at the root of *each specific project workspace*. 

To keep your projects in sync with these standard dotfiles, you should symlink the central file into your individual project repositories:

```bash
# Navigate to your specific project
cd ~/Projects/My-Awesome-Project

# Ensure the .github directory exists
mkdir -p .github

# Create a symlink to the central dotfiles repository
ln -s ~/Projects/Rieken/ai-assistant-dot-files/.github/copilot-instructions.md .github/copilot-instructions.md
```

## AI Feature Team (Multi-Agent Templates & Skills)

This repository also contains a powerful Multi-Agent "Feature Team" pipeline built for **Claude Code**. It provides a fully automated subagent architecture including a Spec Writer, Analyst, Architect, Performance Engineer, Data Engineer, Developer, Code Reviewer, Accessibility Engineer, Security Reviewer, QA Engineer, SRE, Tech Writer, and DevOps Engineer, along with executable **Skills** (like complexity checking and test running) that enforce craftsmanship rules programmatically.

Instead of manually copying these files into your projects, the global `./install` script symlinks the agents (`~/.claude/agents/`) and skills (`~/.claude/skills/`) to your home directory, meaning they are available instantly in any repository you open on your machine.

---

### Using the Agents (Claude Code, Cursor, Copilot)

Because these agents are structured as `.md` system prompts with specialized tools and parameters, they interact slightly differently depending on your IDE/tool.

#### 🤖 Claude Code (Native Support)
Claude Code auto-discovers agents in your `~/.claude/agents/` directory. You can start a pipeline directly from your terminal:
```bash
# Start Claude Code in any project
claude

# Ask an agent to do a job using the @mention syntax:
> @analyst please read ticket 123 and create a spec
> @architect please review the spec and plan the structure
> @developer please implement the architecture using TDD
> @code-reviewer please review the code against our craftsmanship rules
> @qa-engineer please write end-to-end tests
```

#### 🧑‍💻 Cursor IDE
Cursor's Composer (Cmd+I or Cmd+K) can dynamically pull in agent files as context. To actuate an agent:
1. Open Cursor's Composer or Chat.
2. Tag the global agent file using `@developer.md` or `@code-reviewer.md`.
3. Example Prompt: *"Act exactly as described in `@developer.md`. Please read `.claude/feature-workspace/analysis.md` and implement the feature."*

#### ✈️ GitHub Copilot
Copilot does not auto-discover custom agents like Claude Code does, but you can leverage Copilot Edits or Chat by directly feeding it the agent persona:
1. Open GitHub Copilot Chat.
2. Tag the agent file (e.g., `#file:code-reviewer.md`).
3. Example Prompt: *"Act as the Code Reviewer persona from `#file:code-reviewer.md`. Read my current workspace changes and review them against `ARCHITECTURE_RULES.md`."*

---

If you prefer to have the agents copied physically into your project workspace (e.g., to commit them to a specific repo for team distribution), you can use the built-in scaffolding script:

```bash
# Navigate to your specific project
cd ~/Projects/My-Awesome-Project

# Deploy the AI Feature Team locally
~/Projects/Rieken/ai-assistant-dot-files/scaffold-team.sh
```
This deploys the templates directly into the project. Read more in the [Template README](templates/claude-feature-team/README.md).


```text
You are the orchestrator. I want to build this feature: [feature description]. 

Please run this exact pipeline:
1. Pass the prompt to the @analyst to write the {feature-description}-analysis.md file
2. Pass the {feature-description}-analysis.md file to the @architect to write the {feature-description}-architecture.md file
3. Pass the {feature-description}-architecture.md file to the @performance-engineer to write the {feature-description}-performance-report.md file
4. Pass the {feature-description}-performance-report.md file to the @data-engineer to write the {feature-description}-data-engineering.md file
5. Pass the {feature-description}-data-engineering.md file to the @developer to write the {feature-description}-implementation.md file
6. Pass the {feature-description}-implementation.md file to the @code-reviewer to write the {feature-description}-review.md file
7. Pass the {feature-description}-review.md file to the @accessibility-engineer to write the {feature-description}-a11y-report.md file
8. Pass the {feature-description}-a11y-report.md file to the @security-reviewer to write the {feature-description}-security-report.md file
9. Pass the {feature-description}-security-report.md file to the @qa-engineer to write the {feature-description}-tests.md file
10. Pass the {feature-description}-tests.md file to the @sre-engineer to write the {feature-description}-observability-report.md file
11. Pass the {feature-description}-observability-report.md file to the @tech-writer to write the {feature-description}-documentation.md file
12. Pass the {feature-description}-documentation.md file to the @devops-engineer to write the {feature-description}-deployment.md file
```