//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic search_docs MCP tool.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

// Package tools - search_docs_tool.go
//
// SearchDocsTool exposes the BM25Retriever (see bm25_retriever.go) as
// the MCP tool an agent calls to search the installed project's docs corpus.
package tools

import (
	"context"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"

	// Replace with your module's logging package path:
	// "<YOUR_MODULE>/internal/logging"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// DocIndexer refreshes a docs index for the given corpus roots.
type DocIndexer interface {
	EnsureIndex(corpusPaths []string) error
}

const defaultDocsPath = "docs/"

// SearchDocsTool implements domain.Tool for BM25 docs-corpus search.
type SearchDocsTool struct {
	logger    *logging.Logger
	retriever Retriever
	indexer   DocIndexer
}

// NewSearchDocsTool wires the tool with its retriever and indexer.
func NewSearchDocsTool(logger *logging.Logger, retriever Retriever, indexer DocIndexer) *SearchDocsTool {
	return &SearchDocsTool{logger: logger, retriever: retriever, indexer: indexer}
}

func (t *SearchDocsTool) Name() string { return "search_docs" }

func (t *SearchDocsTool) Description() string {
	return "Search the installed project's markdown docs (docs/features/, docs/adrs/, docs/patterns/, docs/runbooks/) via BM25. Returns lexically-ranked references (paths + snippets, never content copies) for the calling LLM to filter semantically."
}

func (t *SearchDocsTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"query"},
		Properties: map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Free-text query. Whitespace-split into tokens; each token is matched via fts5 against doc titles (10x weight) and bodies (1x weight).",
			},
			"docsPath": map[string]interface{}{
				"type":        "string",
				"description": "Corpus root to index and search. Defaults to \"docs/\" (relative to the server's working directory).",
			},
		},
	}
}

func (t *SearchDocsTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&DocSearchResult{})
}

func (t *SearchDocsTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling search_docs request")

	args := request.GetArguments()
	query, _ := args["query"].(string)
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	docsPath, _ := args["docsPath"].(string)
	if docsPath == "" {
		docsPath = defaultDocsPath
	}

	if t.retriever == nil || t.indexer == nil {
		return t.emptyResult(query, "no docs retriever configured")
	}

	if err := t.indexer.EnsureIndex([]string{docsPath}); err != nil {
		t.logger.Error("Docs index refresh failed", "error", err, "docsPath", docsPath)
		return mcp.NewToolResultError(fmt.Sprintf("Index refresh failed: %v", err)), nil
	}

	refs, err := t.retriever.Retrieve(query, nil, "")
	if err != nil {
		t.logger.Error("Docs retrieval failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Retrieval failed: %v", err)), nil
	}

	result := DocSearchResult{
		Success:   true,
		Query:     query,
		TotalHits: len(refs),
		Matches:   convertReferencesToDocMatches(refs),
	}
	return marshalToolResult(t.logger, result, "search_docs")
}

func (t *SearchDocsTool) emptyResult(query, note string) (*mcp.CallToolResult, error) {
	result := DocSearchResult{
		Success:   true,
		Query:     fmt.Sprintf("%s (%s)", query, note),
		TotalHits: 0,
	}
	return marshalToolResult(t.logger, result, "search_docs")
}

func convertReferencesToDocMatches(refs []Reference) []DocMatch {
	matches := make([]DocMatch, 0, len(refs))
	for _, r := range refs {
		matches = append(matches, DocMatch{
			Title:     r.Title,
			Path:      r.Path,
			Summary:   r.Summary,
			Relevance: r.Relevance,
		})
	}
	return matches
}
