package server

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/orieken/loom/shared/mcp/internal/domain"
)

type stubTool struct {
	name         string
	description  string
	inputSchema  json.RawMessage
	outputSchema json.RawMessage
	result       *domain.ToolResult
	err          error
	gotRequest   domain.ToolRequest
}

func (s *stubTool) Name() string                  { return s.name }
func (s *stubTool) Description() string           { return s.description }
func (s *stubTool) InputSchema() json.RawMessage  { return s.inputSchema }
func (s *stubTool) OutputSchema() json.RawMessage { return s.outputSchema }

func (s *stubTool) Execute(_ context.Context, request domain.ToolRequest) (*domain.ToolResult, error) {
	s.gotRequest = request
	return s.result, s.err
}

func TestMCPToolDefinitionMapsMetadata(t *testing.T) {
	stub := &stubTool{
		name:         "stub_tool",
		description:  "a stub",
		inputSchema:  json.RawMessage(`{"type":"object"}`),
		outputSchema: json.RawMessage(`{"type":"object","properties":{}}`),
	}

	definition := mcpToolDefinition(stub)

	if definition.Name != "stub_tool" || definition.Description != "a stub" {
		t.Errorf("metadata not mapped: got name=%q description=%q", definition.Name, definition.Description)
	}
	if string(definition.RawInputSchema) != `{"type":"object"}` {
		t.Errorf("input schema not mapped: got %s", definition.RawInputSchema)
	}
	if string(definition.RawOutputSchema) != `{"type":"object","properties":{}}` {
		t.Errorf("output schema not mapped: got %s", definition.RawOutputSchema)
	}
	if _, err := json.Marshal(definition); err != nil {
		t.Errorf("wire definition does not marshal: %v", err)
	}
}

func TestMCPToolHandlerConvertsRequestAndResult(t *testing.T) {
	tests := []struct {
		name        string
		result      *domain.ToolResult
		wantText    string
		wantIsError bool
	}{
		{"text result", domain.NewTextResult("hello"), "hello", false},
		{"error result", domain.NewErrorResult("bad input"), "bad input", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubTool{name: "stub_tool", result: tt.result}
			request := mcp.CallToolRequest{}
			request.Params.Arguments = map[string]any{"projectPath": "/tmp/x"}

			got, err := mcpToolHandler(stub)(context.Background(), request)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertDomainRequest(t, stub.gotRequest)
			assertTextResult(t, got, tt.wantText, tt.wantIsError)
		})
	}
}

func assertDomainRequest(t *testing.T, got domain.ToolRequest) {
	t.Helper()
	if got.Name != "stub_tool" {
		t.Errorf("request name not set: got %q", got.Name)
	}
	if got.StringArg("projectPath") != "/tmp/x" {
		t.Errorf("arguments not converted: got %v", got.Args)
	}
}

func assertTextResult(t *testing.T, got *mcp.CallToolResult, wantText string, wantIsError bool) {
	t.Helper()
	if got.IsError != wantIsError {
		t.Errorf("IsError = %v, want %v", got.IsError, wantIsError)
	}
	text, ok := got.Content[0].(mcp.TextContent)
	if !ok || text.Text != wantText {
		t.Errorf("content = %#v, want text %q", got.Content, wantText)
	}
}

func TestMCPToolHandlerPropagatesExecuteError(t *testing.T) {
	executeErr := errors.New("transport failure")
	stub := &stubTool{name: "stub_tool", err: executeErr}

	_, err := mcpToolHandler(stub)(context.Background(), mcp.CallToolRequest{})
	if !errors.Is(err, executeErr) {
		t.Errorf("expected execute error to propagate, got %v", err)
	}
}

func TestMCPResultNilStaysNil(t *testing.T) {
	if got := mcpResult(nil); got != nil {
		t.Errorf("expected nil, got %#v", got)
	}
}
