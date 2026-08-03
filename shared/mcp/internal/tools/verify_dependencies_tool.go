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

// VerifyDependenciesTool exposes DependencyBoundaryAnalyzer as an MCP tool.
type VerifyDependenciesTool struct {
	logger   *logging.Logger
	analyzer *analyzers.DependencyBoundaryAnalyzer
}

// NewVerifyDependenciesTool wires the tool with its dependencies.
func NewVerifyDependenciesTool(logger *logging.Logger, analyzer *analyzers.DependencyBoundaryAnalyzer) *VerifyDependenciesTool {
	return &VerifyDependenciesTool{logger: logger, analyzer: analyzer}
}

func (t *VerifyDependenciesTool) Name() string { return "verify_dependencies" }

func (t *VerifyDependenciesTool) Description() string {
	return "Verify Clean Architecture layer boundaries by scanning Go and TypeScript imports, flagging any inner-to-outer layer dependency"
}

func (t *VerifyDependenciesTool) InputSchema() mcp.ToolInputSchema {
	return projectPathOnlySchema()
}

func (t *VerifyDependenciesTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&analyzers.DependencyVerificationResult{})
}

func (t *VerifyDependenciesTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling verify_dependencies request")

	args := request.GetArguments()
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := t.analyzer.Analyze(projectPath)
	if err != nil {
		t.logger.Error("Dependency verification failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Dependency verification failed: %v", err)), nil
	}

	body, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal dependency verification result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Dependency verification completed", "path", projectPath, "violations", result.ViolationsCount)
	return mcp.NewToolResultText(string(body)), nil
}
