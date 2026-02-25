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

## AI Feature Team (Multi-Agent Templates)

This repository also contains a powerful Multi-Agent "Feature Team" pipeline built for **Claude Code**. It provides a fully automated subagent architecture including an Analyst, Developer, QA Engineer, Tech Writer, and DevOps Engineer.

Instead of manually copying these files into your projects, you can use the built-in scaffolding script to instantly deploy the feature team into any workspace:

```bash
# Navigate to your specific project
cd ~/Projects/My-Awesome-Project

# Deploy the AI Feature Team
~/Projects/Rieken/ai-assistant-dot-files/scaffold-team.sh
```

This will automatically create the `.claude/agents`, `.claude/skills`, and `features/` directories in your project, empowering Claude to handle end-to-end feature delivery. Read more in the [Template README](templates/claude-feature-team/README.md).
