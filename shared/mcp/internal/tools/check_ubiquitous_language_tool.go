package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/analyzers"
	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/logging"
)

// CheckUbiquitousLanguageTool exposes UbiquitousLanguageAnalyzer as an MCP tool.
type CheckUbiquitousLanguageTool struct {
	logger   *logging.Logger
	analyzer *analyzers.UbiquitousLanguageAnalyzer
}

// NewCheckUbiquitousLanguageTool wires the tool with its dependencies.
func NewCheckUbiquitousLanguageTool(logger *logging.Logger, analyzer *analyzers.UbiquitousLanguageAnalyzer) *CheckUbiquitousLanguageTool {
	return &CheckUbiquitousLanguageTool{logger: logger, analyzer: analyzer}
}

func (t *CheckUbiquitousLanguageTool) Name() string { return "check_ubiquitous_language" }

func (t *CheckUbiquitousLanguageTool) Description() string {
	return "Scan source files for uses of unapproved synonyms defined in a DOMAIN_DICTIONARY.md, reporting each violation with the canonical replacement"
}

func (t *CheckUbiquitousLanguageTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath", "dictionaryPath"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
			"dictionaryPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to the DOMAIN_DICTIONARY.md that defines canonical terms and their forbidden synonyms",
			},
		},
	}
}

func (t *CheckUbiquitousLanguageTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&analyzers.UbiquitousLanguageResult{})
}

func (t *CheckUbiquitousLanguageTool) Execute(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling check_ubiquitous_language request")

	args := request.GetArguments()
	projectPath, _ := args["projectPath"].(string)
	dictionaryPath, _ := args["dictionaryPath"].(string)

	if projectPath == "" || dictionaryPath == "" {
		return mcp.NewToolResultError("both projectPath and dictionaryPath are required"), nil
	}

	result, err := t.analyzer.Analyze(projectPath, dictionaryPath)
	if err != nil {
		t.logger.Error("Ubiquitous language analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Ubiquitous language analysis failed: %v", err)), nil
	}

	body, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal ubiquitous language result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Ubiquitous language analysis completed", "path", projectPath, "violations", result.ViolationsCount)
	return mcp.NewToolResultText(string(body)), nil
}
