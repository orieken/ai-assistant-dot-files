// mcp_adapter.go owns every conversion between the transport-free domain
// types and the mcp-go wire types. No other file in the module may translate
// between the two — domain stays stdlib-only (roadmap M0.3).
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/orieken/loom/internal/telemetry"
	"github.com/orieken/loom/shared/mcp/internal/domain"
)

// mcpToolDefinition converts a domain.Tool's metadata into the MCP wire type.
func mcpToolDefinition(tool domain.Tool) mcp.Tool {
	return mcp.Tool{
		Name:            tool.Name(),
		Description:     tool.Description(),
		RawInputSchema:  tool.InputSchema(),
		RawOutputSchema: tool.OutputSchema(),
	}
}

// mcpToolHandler adapts a domain.Tool's Execute into an mcp-go handler,
// wrapping each call in a span and a correlated log line. This is the only
// place tool-call telemetry is emitted; the tools themselves stay unaware
// of it, and `internal/domain` stays stdlib-only (guardrail #8 and M0.3).
func (h *Handler) mcpToolHandler(tool domain.Tool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		call := telemetry.ToolCall{Name: tool.Name(), Arguments: stringArguments(request.GetArguments())}
		ctx, span := h.session.StartTool(ctx, call)
		h.logToolCall(ctx, tool.Name())
		result, err := tool.Execute(ctx, domainRequest(tool.Name(), request))
		span.End(toolResult(result, err))
		if err != nil {
			return nil, err
		}
		return mcpResult(result), nil
	}
}

// logToolCall correlates the server's own logs with the trace, so a log
// line and a span can be joined without guessing from timestamps.
func (h *Handler) logToolCall(ctx context.Context, name string) {
	traceID, spanID := telemetry.TraceIDs(ctx)
	if traceID == "" {
		h.logger.Info("tool.called", "tool", name)
		return
	}
	h.logger.Info("tool.called", "tool", name, "trace_id", traceID, "span_id", spanID)
}

func toolResult(result *domain.ToolResult, err error) telemetry.ToolResult {
	if result == nil {
		return telemetry.ToolResult{Err: err}
	}
	return telemetry.ToolResult{Preview: resultText(result), IsError: result.IsError, Err: err}
}

func resultText(result *domain.ToolResult) string {
	parts := make([]string, 0, len(result.Content))
	for _, block := range result.Content {
		parts = append(parts, block.Text)
	}
	return strings.Join(parts, "\n")
}

// stringArguments renders each argument for a span attribute. Values are
// stringified here rather than in the telemetry package so that package's
// signatures stay free of `any`, per go-conventions.
func stringArguments(args map[string]any) map[string]string {
	rendered := make(map[string]string, len(args))
	for key, value := range args {
		rendered[key] = renderArgument(value)
	}
	return rendered
}

// renderArgument keeps strings as themselves and JSON-encodes everything
// else, so a nested object arrives readable rather than as a Go %v dump.
func renderArgument(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(encoded)
}

func domainRequest(name string, request mcp.CallToolRequest) domain.ToolRequest {
	return domain.ToolRequest{Name: name, Args: request.GetArguments()}
}

func mcpResult(result *domain.ToolResult) *mcp.CallToolResult {
	if result == nil {
		return nil
	}
	content := make([]mcp.Content, 0, len(result.Content))
	for _, block := range result.Content {
		content = append(content, mcp.NewTextContent(block.Text))
	}
	return &mcp.CallToolResult{Content: content, IsError: result.IsError}
}
