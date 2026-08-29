// Package domain defines the framework's transport-free tool abstraction.
// It depends on the standard library only — every MCP wire type lives in the
// server adapter layer (see architecture-guardrails.md #1 and roadmap M0.3).
// TestDomainImportsOnlyStdlib enforces this boundary.
package domain

import (
	"context"
	"encoding/json"
)

// ToolRequest is a transport-free tool invocation.
type ToolRequest struct {
	Name string
	Args map[string]any
}

// StringArg returns the string argument named key, or "" when it is absent
// or not a string.
func (r ToolRequest) StringArg(key string) string {
	value, _ := r.Args[key].(string)
	return value
}

// ContentBlock is one piece of tool output. Text is the only content kind
// the framework tools emit today.
type ContentBlock struct {
	Text string
}

// ToolResult is a transport-free tool outcome. IsError marks a failure the
// calling LLM should repair (bad arguments, analysis failure); transport-level
// failures are returned as Go errors from Execute instead.
type ToolResult struct {
	Content []ContentBlock
	IsError bool
}

// NewTextResult wraps text in a successful single-block result.
func NewTextResult(text string) *ToolResult {
	return &ToolResult{Content: []ContentBlock{{Text: text}}}
}

// NewErrorResult wraps message in a single-block result flagged as an error.
func NewErrorResult(message string) *ToolResult {
	return &ToolResult{Content: []ContentBlock{{Text: message}}, IsError: true}
}

// Tool is the framework's first-class abstraction for every capability exposed
// over MCP. Every framework tool implements this interface so it can be
// registered uniformly without the server knowing any tool's concrete type.
// Schemas are raw JSON Schema documents; OutputSchema may return nil when the
// tool declares none.
type Tool interface {
	Name() string
	Description() string
	InputSchema() json.RawMessage
	OutputSchema() json.RawMessage
	Execute(ctx context.Context, request ToolRequest) (*ToolResult, error)
}
