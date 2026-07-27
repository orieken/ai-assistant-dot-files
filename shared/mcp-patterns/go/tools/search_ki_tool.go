//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic search_ki MCP tool.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"

	// Replace with your module's logging package path:
	// "<YOUR_MODULE>/internal/logging"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// SearchKITool surfaces Knowledge Items and ADRs from the framework corpus.
type SearchKITool struct {
	logger    *logging.Logger
	retriever Retriever
}

// NewSearchKITool wires the tool with its retriever.
func NewSearchKITool(logger *logging.Logger, retriever Retriever) *SearchKITool {
	return &SearchKITool{logger: logger, retriever: retriever}
}

func (t *SearchKITool) Name() string { return "search_ki" }

func (t *SearchKITool) Description() string {
	return "Search the framework's Knowledge Items and ADRs by query, tags, and domain. Returns lexically-ranked references (paths + summaries, never content copies) for the calling LLM to filter semantically."
}

func (t *SearchKITool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"query"},
		Properties: map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Free-text query. Whitespace-split into tokens; matches against KI/ADR titles and summaries.",
			},
			"tags": map[string]interface{}{
				"type":        "array",
				"description": "Optional tag filter. Exact tag matches boost relevance strongly.",
				"items":       map[string]interface{}{"type": "string"},
			},
			"domain": map[string]interface{}{
				"type":        "string",
				"description": "Optional domain / bounded-context filter. Exact domain match boosts relevance.",
			},
		},
	}
}

func (t *SearchKITool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&KISearchResult{})
}

func (t *SearchKITool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling search_ki request")

	args := request.GetArguments()
	query, _ := args["query"].(string)
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	tags := extractStringArray(args["tags"])
	domain, _ := args["domain"].(string)

	if t.retriever == nil {
		return t.emptyResult(query, "no retriever configured")
	}

	refs, err := t.retriever.Retrieve(query, tags, domain)
	if err != nil {
		t.logger.Error("KI retrieval failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Retrieval failed: %v", err)), nil
	}

	result := KISearchResult{
		Success:   true,
		Query:     query,
		TotalHits: len(refs),
		Matches:   convertReferencesToMatches(refs),
	}
	return marshalToolResult(t.logger, result, "search_ki")
}

func (t *SearchKITool) emptyResult(query, note string) (*mcp.CallToolResult, error) {
	result := KISearchResult{
		Success:   true,
		Query:     fmt.Sprintf("%s (%s)", query, note),
		TotalHits: 0,
	}
	return marshalToolResult(t.logger, result, "search_ki")
}

func extractStringArray(v interface{}) []string {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}

func convertReferencesToMatches(refs []Reference) []KIMatch {
	matches := make([]KIMatch, 0, len(refs))
	for _, r := range refs {
		matches = append(matches, KIMatch{
			Title:     r.Title,
			Path:      r.Path,
			Summary:   r.Summary,
			Tags:      r.Tags,
			Relevance: r.Relevance,
		})
	}
	return matches
}

func marshalToolResult(logger *logging.Logger, v interface{}, toolName string) (*mcp.CallToolResult, error) {
	resultJSON, err := json.Marshal(v)
	if err != nil {
		logger.Error("Failed to marshal result", "tool", toolName, "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}
