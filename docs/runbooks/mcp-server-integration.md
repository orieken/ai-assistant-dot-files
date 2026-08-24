# Runbook: Exposing This Framework Through an Existing MCP Server

How to make this repo's agents, skills, and rules available to a team that already runs its own MCP
server — without asking them to adopt this repo, run `install.sh`, or change their existing tooling. The
deliverable at the bottom of this file is a self-contained prompt: paste it into an AI coding session
that has access to the target MCP server's codebase, and it does the integration work there.

> **Status note**: an earlier attempt at this (`docs/mcp/framework-tools-prompts.md`, `rag-example.md`)
> was deleted in Epic 33 — it targeted an `mcp/` server that was never built in this repo, and referenced
> a separate "aakg-mcp integration" that isn't documented anywhere here either. This runbook replaces
> that attempt with something narrower and actually usable: not "build an MCP server," but "add this
> framework's content to an MCP server someone else already has."

## Why this instead of `install.sh`

`install.sh` assumes the target is a project an AI coding assistant (Claude Code, Cursor, etc.) reads
files from directly. An MCP server is a different kind of target: it's a long-running process that
serves content to *any* MCP client over a defined protocol, so a team's existing internal MCP server can
expose this framework to everyone who already connects to it — no per-project install, no per-developer
symlink setup. The tradeoff: MCP servers don't read frontmatter/directory conventions the way
`install.sh`'s generators do, so this is a real (if mechanical) integration, not a copy operation.

## Concept mapping

MCP has three primitive types, and they're not interchangeable — a model *invokes* a Tool, a client
*loads* a Resource, a user *selects* a Prompt. Mapping our concepts onto the wrong one produces a
technically-working but semantically-wrong integration (e.g. a Tool that doesn't actually *do* anything
observable, just returns static text — that's a Resource wearing a Tool's clothes).

| Our concept | MCP primitive | Why |
|---|---|---|
| `shared/skills/*/SKILL.md` | **Prompt** | A skill is already a user/keyword-triggered template with a defined process and output shape — almost exactly what an MCP Prompt is. Its `triggers.keywords`/`intentPatterns` become the Prompt's description/argument hints. |
| `shared/agents/*.md` | **Resource** | An agent's full body (role, process, output format, rules) is read-only context a client loads before reasoning — not an action to invoke. Exposing it as a Resource lets any MCP client pull in "the analyst persona" as context without this server having to actually *run* the analyst itself. |
| `shared/rules/*.md`, `ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md` | **Resource** | Same reasoning as agents — static reference context, not an action. |
| A skill whose target MCP server *can* actually execute it end-to-end (e.g. `run-tests`, `analyze-complexity` — skills with a `check.sh`/`run.sh`) | **Tool**, in addition to the Prompt | If the target server has the same execution environment the skill's script needs (a shell, the target repo checked out), expose it as a real Tool too — the model can then invoke it directly instead of just being handed the prompt template. Skills without an executable component (most of them) stay Prompt-only; don't fabricate a Tool that just returns text. |

Don't collapse everything into Tools. A Resource that never changes based on input and a Prompt a user
explicitly picks are both real, different UX in an MCP client — flattening them into "everything is a
tool the model can call" loses that distinction and is likely to confuse whichever MCP client (Claude
Desktop, a custom client, an IDE extension) the target team is actually using.

## What this runbook does NOT cover

- Building a new MCP server from scratch — this is additive to an *existing* one.
- Auth/transport setup (stdio vs. SSE vs. HTTP) — that's already decided by whatever the target server
  runs today; this integration doesn't change it.
- Keeping the target server's copy in sync with future edits to this repo — that's a real follow-up
  question (poll `shared/` at startup? vendor a copy and re-run this integration periodically? git
  submodule?) worth resolving with the target team, not decided unilaterally here.

---

## The prompt

Paste everything below into an AI coding session with access to the target MCP server's repository.

```markdown
You're adding this framework's content (agents, skills, rules) as MCP primitives to an existing MCP
server in this repository. You do NOT have access to the framework's source repo
(github.com/orieken/loom) directly in this session — I'll paste in the specific
files you need as we go, or you can fetch them if you have network/git access to that repo. Don't
invent agent/skill content; ask me for the specific file if you need one you don't have.

## Step 1: Identify the target server

Find the existing MCP server in this repo. Determine:
- Language/SDK: TypeScript (`@modelcontextprotocol/sdk`) or Python (`mcp` / `FastMCP`)? Check
  `package.json` or `pyproject.toml`/`requirements.txt`.
- Is it using the high-level API (`McpServer` in TS, `FastMCP` in Python) or the low-level
  request-handler API? The patterns below assume high-level; if it's low-level, adapt the same mapping
  to `server.setRequestHandler(...)`/`@server.list_prompts()` etc. instead of asking me first — but
  flag the difference in your summary at the end.
- Where existing tools/resources/prompts are registered — follow that file's own organization
  (one file per primitive? one big registration file? a `tools/`, `resources/`, `prompts/` directory
  split?) rather than introducing a new structure.

## Step 2: Confirm the primitive mapping before writing code

- Every `SKILL.md` becomes a **Prompt**. Its `name` frontmatter field is the Prompt name: its
  `description` becomes the Prompt's description; anything in "Process" that asks a question or expects
  an argument becomes a Prompt argument (e.g. a skill that operates "against a PR number, URL, or branch
  name" gets one argument covering that). The Prompt's returned message should be the skill's own
  "Context To Load First" + "Process" + "Output Format" sections concatenated into the instruction the
  model receives when the prompt is selected — don't summarize or compress it, the skill file's own
  prose is already the correctly-scoped instruction.
- Every agent `.md` file (`shared/agents/*.md`) becomes a **Resource**, URI scheme `agent://<name>`
  (e.g. `agent://analyst`). The resource's content is the agent file's full body verbatim (frontmatter
  can be omitted from the served content, or kept — your call, note which in your summary). MIME type
  `text/markdown`.
- Every rule file (`shared/rules/*.md`, `ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`) becomes a
  **Resource** too, URI scheme `rule://<name>` / `rule://architecture` / `rule://domain-dictionary`.
- Only skills with an executable component (a `check.sh`/`run.sh` sitting next to the `SKILL.md`) also
  become a **Tool**, in addition to their Prompt. The Tool's input schema matches whatever the script
  actually takes as arguments (read the script, don't guess). The Tool's handler shells out to that
  script and returns its output. Every other skill (most of them) is Prompt-only — do not fabricate a
  Tool for a skill that's pure LLM judgment with no executable backing.

Show me the list of Prompts/Resources/Tools you're about to create before writing any registration
code. Wait for my confirmation.

## Step 3: Implement

Illustrative patterns below — verify against whatever SDK version is actually installed in this repo
(`package.json`/`pyproject.toml`) before finalizing; the MCP SDKs have moved fast and the exact
high-level API surface may have changed since these examples were written.

**TypeScript** (`@modelcontextprotocol/sdk`, high-level `McpServer`):
```typescript
import { z } from "zod";

// Skill -> Prompt
server.prompt(
  "review-pr",
  { prNumber: z.string().describe("PR number, URL, or branch name") },
  ({ prNumber }) => ({
    messages: [{
      role: "user",
      content: { type: "text", text: `<full contents of shared/skills/review-pr/SKILL.md,
        with ${prNumber} substituted where the skill expects a PR identifier>` }
    }]
  })
);

// Agent -> Resource
server.resource(
  "analyst",
  "agent://analyst",
  async (uri) => ({
    contents: [{ uri: uri.href, mimeType: "text/markdown", text: `<full contents of
      shared/agents/analyst.md>` }]
  })
);

// Skill with a check.sh -> also a Tool
server.tool(
  "analyze-complexity",
  { path: z.string().describe("File or directory to analyze") },
  async ({ path }) => {
    const result = await runShellScript("shared/skills/analyze-complexity/check.sh", [path]);
    return { content: [{ type: "text", text: result }] };
  }
);
```

**Python** (`mcp.server.fastmcp.FastMCP`):
```python
# Skill -> Prompt
@mcp.prompt()
def review_pr(pr_number: str) -> str:
    """PR number, URL, or branch name"""
    return "<full contents of shared/skills/review-pr/SKILL.md, with pr_number substituted>"

# Agent -> Resource
@mcp.resource("agent://analyst")
def analyst_agent() -> str:
    return "<full contents of shared/agents/analyst.md>"

# Skill with a check.sh -> also a Tool
@mcp.tool()
def analyze_complexity(path: str) -> str:
    """File or directory to analyze"""
    return run_shell_script("shared/skills/analyze-complexity/check.sh", [path])
```

Register every agent as a Resource and every skill as a Prompt (plus a Tool where a `check.sh`/`run.sh`
exists) following whichever pattern matches the target repo's language. Keep the content verbatim from
the source `.md` files — don't paraphrase or compress an agent's process or a skill's guardrails when
serving them.

## Step 4: Verify

- The target server still starts cleanly.
- List the newly registered Prompts/Resources/Tools via the server's own introspection (or a test MCP
  client call) and confirm the count matches what you listed in Step 2.
- Pick one Prompt (a skill) and one Resource (an agent) and confirm their served content matches the
  source `.md` file exactly — no truncation, no summarization.
- If you added any Tools, run one against a real input and confirm it actually shells out to the
  script and returns real output, not a stub.

## Step 5: Summarize

Report: how many Prompts/Resources/Tools were added, which SDK/language pattern you followed, whether
you used the high-level or low-level API, and anything you weren't able to verify (e.g. "couldn't start
the server locally to confirm Step 4" — say so plainly rather than claiming untested success).
```

---

## Keeping this current

If the framework's agent/skill count changes meaningfully (new agents, retired skills, a changed
Prompt-vs-Tool split for something with a new `check.sh`), this prompt itself doesn't need to change —
it reads the source `.md` files' own structure rather than hardcoding a list. What *would* need
revisiting: the SDK code examples in Step 3, if the MCP SDKs' high-level API surface changes in a
breaking way. Verify those examples against current SDK docs before handing this prompt to a team,
rather than assuming they're still accurate.
