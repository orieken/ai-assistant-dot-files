// Package echotool is the example's custom capability. Note what it does NOT
// import: no mcp-go, no jsonschema, no loom internals — implementing a tool
// against the public embedding API needs only the stdlib and
// github.com/orieken/loom/tools (roadmap D.2's done-when criterion).
package echotool

import (
	"context"
	"encoding/json"
	"time"

	"github.com/orieken/loom/tools"
)

// EchoTool repeats the text it is given — a minimal custom Tool.
type EchoTool struct{}

func (EchoTool) Name() string { return "echo" }

func (EchoTool) Description() string {
	return "Echo the given text back to the caller"
}

func (EchoTool) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"required": ["text"],
		"properties": {
			"text": {"type": "string", "description": "Text to echo back"}
		}
	}`)
}

func (EchoTool) OutputSchema() json.RawMessage { return nil }

// Execute returns the text argument, or an error result when it is missing.
func (EchoTool) Execute(_ context.Context, request tools.ToolRequest) (*tools.ToolResult, error) {
	text := request.StringArg("text")
	if text == "" {
		return tools.NewErrorResult("text is required"), nil
	}
	return tools.NewTextResult(text), nil
}

// Registration couples the tool with its execution metadata.
func Registration() tools.ToolRegistration {
	return tools.ToolRegistration{
		Tool:       EchoTool{},
		Timeout:    5 * time.Second,
		Retry:      tools.RetryIdempotent,
		Permission: tools.ScopeReadOnly,
	}
}
