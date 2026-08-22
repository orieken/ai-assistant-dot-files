package cmd

import "github.com/spf13/cobra"

// toolTier classifies a context tool by how much value it adds to the pipeline.
type toolTier string

const (
	tierHigh   toolTier = "high"
	tierMedium toolTier = "medium"
)

// contextTool describes one opt-in context or memory optimization tool.
type contextTool struct {
	name        string
	tier        toolTier
	description string
	binary      string   // binary to check in PATH; empty = built-in, always available
	installCmds []string // brew/cargo/npm commands to run; nil = must be installed manually
	manualNote  string   // shown instead of installCmds when auto-install is not possible
	postNote    string   // shown after a successful install (e.g. required post-install step)
	kiName      string
}

// contextTools is the registry of all opt-in tools known to loom.
// Order within a tier determines display order.
var contextTools = []contextTool{
	{
		name:        "tokei",
		tier:        tierHigh,
		description: "accurate per-language line counts for token budget estimation (context-engineer Step 5)",
		binary:      "tokei",
		installCmds: []string{"brew install tokei"},
		kiName:      "tokei-token-budget",
	},
	{
		name:        "git-churn",
		tier:        tierHigh,
		description: "file change frequency for pin-depth decisions — built into git, no install needed",
		binary:      "",
		kiName:      "git-churn-risk-signal",
	},
	{
		name:        "ast-grep",
		tier:        tierHigh,
		description: "AST-level structural search for scope discovery (context-engineer Step 1, refactor-engineer)",
		binary:      "sg",
		installCmds: []string{"brew install ast-grep"},
		kiName:      "ast-grep-structural-search",
	},
	{
		name:        "ctx",
		tier:        tierHigh,
		description: "search past agent sessions; ctx blame traces a file back to its originating session",
		binary:      "ctx",
		manualNote:  "curl -fsSL https://ctx.rs/install | sh",
		postNote:    "after install, run: ctx setup",
		kiName:      "ctx-session-history-search",
	},
	{
		name:        "repomix",
		tier:        tierMedium,
		description: "packs a subsystem into a single LLM-ready file (useful for small, self-contained bounded contexts)",
		binary:      "repomix",
		installCmds: []string{"npm install -g repomix"},
		kiName:      "repomix-codebase-packing",
	},
	{
		name:        "aider-repo-map",
		tier:        tierMedium,
		description: "compact symbol graph for unknown bounded contexts (context-engineer Step 1)",
		binary:      "aider",
		installCmds: []string{"pip install aider-chat"},
		kiName:      "aider-repo-map",
	},
	{
		name:        "ollama",
		tier:        tierMedium,
		description: "local embedding backend for search-ki-semantic (offline / air-gapped environments)",
		binary:      "ollama",
		installCmds: []string{"brew install ollama"},
		postNote:    "after install, run: ollama pull nomic-embed-text",
		kiName:      "ollama-local-embeddings",
	},
}

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Manage opt-in context and memory optimization tools",
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}
