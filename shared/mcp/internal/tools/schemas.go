package tools

import "github.com/mark3labs/mcp-go/mcp"

func projectPathProperty() map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": "Absolute path to the project root",
	}
}

func projectPathOnlySchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
		},
	}
}
