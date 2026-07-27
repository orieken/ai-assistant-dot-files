//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic analyze_complexity MCP tool.
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

	// Replace with your module's analyzer and logging package paths:
	// "<YOUR_MODULE>/internal/analyzers"
	// "<YOUR_MODULE>/internal/logging" (requires Info/Error methods, e.g. *slog.Logger)
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// AnalyzeComplexityTool analyzes code files for cyclomatic complexity and LOC function length.
type AnalyzeComplexityTool struct {
	logger   *logging.Logger
	analyzer *analyzers.ComplexityAnalyzer
}

// NewAnalyzeComplexityTool wires the tool with its dependencies.
func NewAnalyzeComplexityTool(logger *logging.Logger, analyzer *analyzers.ComplexityAnalyzer) *AnalyzeComplexityTool {
	return &AnalyzeComplexityTool{logger: logger, analyzer: analyzer}
}

func (t *AnalyzeComplexityTool) Name() string { return "analyze_complexity" }

func (t *AnalyzeComplexityTool) Description() string {
	return "Analyze cyclomatic complexity and function length against framework thresholds (complexity < 7, LOC < 30)"
}

func (t *AnalyzeComplexityTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
			"maxComplexity": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum allowed cyclomatic complexity (default 7)",
			},
			"maxLines": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum allowed lines of code per function (default 30)",
			},
		},
	}
}

func (t *AnalyzeComplexityTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&ComplexityAnalysisResult{})
}

func (t *AnalyzeComplexityTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling analyze_complexity request")

	args := request.GetArguments()
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	maxComplexity := 7
	if mc, ok := args["maxComplexity"].(float64); ok && mc > 0 {
		maxComplexity = int(mc)
	}

	maxLines := 30
	if ml, ok := args["maxLines"].(float64); ok && ml > 0 {
		maxLines = int(ml)
	}

	result, err := t.analyzer.Analyze(projectPath, maxComplexity, maxLines)
	if err != nil {
		t.logger.Error("Complexity analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Complexity analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal complexity result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Complexity analysis completed", "path", projectPath, "violations", result.ViolationsCount)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
