// mcp_adapter.go owns every conversion between the transport-free domain
// types and the mcp-go wire types. No other file in the module may translate
// between the two — domain stays stdlib-only (roadmap M0.3).
package server

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

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

// mcpToolHandler adapts a domain.Tool's Execute into an mcp-go handler.
func mcpToolHandler(tool domain.Tool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := tool.Execute(ctx, domainRequest(tool.Name(), request))
		if err != nil {
			return nil, err
		}
		return mcpResult(result), nil
	}
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
