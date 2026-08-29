package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/orieken/loom/shared/mcp/internal/analyzers"
	"github.com/orieken/loom/shared/mcp/internal/logging"
)

// CheckAccessibilityTool exposes the AccessibilityAnalyzer as an MCP tool.
type CheckAccessibilityTool struct {
	logger   *logging.Logger
	analyzer *analyzers.AccessibilityAnalyzer
}

// NewCheckAccessibilityTool wires the tool with its dependencies.
func NewCheckAccessibilityTool(logger *logging.Logger, analyzer *analyzers.AccessibilityAnalyzer) *CheckAccessibilityTool {
	return &CheckAccessibilityTool{logger: logger, analyzer: analyzer}
}

func (t *CheckAccessibilityTool) Name() string { return "check_accessibility" }

func (t *CheckAccessibilityTool) Description() string {
	return "Scan UI template files (HTML, Vue, JSX, TSX, Svelte) for semantic-HTML and ARIA accessibility violations"
}

func (t *CheckAccessibilityTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"filePath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to a single UI template file to scan",
			},
			"projectPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to a project root; the walker scans every .html/.htm/.vue/.jsx/.tsx/.svelte file underneath",
			},
		},
	}
}

func (t *CheckAccessibilityTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&analyzers.AccessibilityReportResult{})
}

func (t *CheckAccessibilityTool) Execute(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling check_accessibility request")

	args := request.GetArguments()
	target := resolveAccessibilityTarget(args)
	if target == "" {
		return mcp.NewToolResultError("either filePath or projectPath is required"), nil
	}

	result, err := t.analyzer.Analyze(target)
	if err != nil {
		t.logger.Error("Accessibility analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Accessibility analysis failed: %v", err)), nil
	}

	body, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal accessibility result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Accessibility analysis completed", "path", target, "violations", result.ViolationsCount)
	return mcp.NewToolResultText(string(body)), nil
}

func resolveAccessibilityTarget(args map[string]any) string {
	if fp, ok := args["filePath"].(string); ok && fp != "" {
		return fp
	}
	if pp, ok := args["projectPath"].(string); ok && pp != "" {
		return pp
	}
	return ""
}
