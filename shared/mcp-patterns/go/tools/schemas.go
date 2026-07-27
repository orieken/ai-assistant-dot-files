//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic input schema helpers.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

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
