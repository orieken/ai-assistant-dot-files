// Command embedding is a minimal MCP server that merges loom's framework
// tools with one custom tool (echotool) and serves them over stdio.
//
// The adapter below converts the transport-free tools.Registry to THIS
// module's own mcp-go dependency — which is the point of the embedding API:
// your mcp-go version is yours, loom's is loom's, and the two never meet in
// a type signature.
package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/orieken/loom/examples/embedding/echotool"
	"github.com/orieken/loom/shared/mcp/register"
	"github.com/orieken/loom/tools"
)

func main() {
	registry, err := buildRegistry()
	if err != nil {
		log.Fatalf("build registry: %v", err)
	}
	if err := server.ServeStdio(newServer(registry)); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

// buildRegistry merges loom's built-in tools with this server's own.
func buildRegistry() (*tools.Registry, error) {
	registry := tools.NewRegistry()
	if err := registry.Register(echotool.Registration()); err != nil {
		return nil, err
	}
	if err := registry.Merge(register.Frameworks(nil)); err != nil {
		return nil, err
	}
	return registry, nil
}

func newServer(registry *tools.Registry) *server.MCPServer {
	s := server.NewMCPServer("loom-embedding-example", "0.1.0", server.WithToolCapabilities(true))
	for _, registration := range registry.All() {
		s.AddTool(wireDefinition(registration.Tool), wireHandler(registration.Tool))
	}
	return s
}

// wireDefinition adapts a tools.Tool's metadata to this module's mcp-go.
func wireDefinition(tool tools.Tool) mcp.Tool {
	return mcp.Tool{
		Name:            tool.Name(),
		Description:     tool.Description(),
		RawInputSchema:  tool.InputSchema(),
		RawOutputSchema: tool.OutputSchema(),
	}
}

// wireHandler adapts a tools.Tool's Execute to this module's mcp-go.
func wireHandler(tool tools.Tool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := tool.Execute(ctx, tools.ToolRequest{Name: tool.Name(), Args: request.GetArguments()})
		if err != nil {
			return nil, err
		}
		return wireResult(result), nil
	}
}

func wireResult(result *tools.ToolResult) *mcp.CallToolResult {
	if result == nil {
		return nil
	}
	content := make([]mcp.Content, 0, len(result.Content))
	for _, block := range result.Content {
		content = append(content, mcp.NewTextContent(block.Text))
	}
	return &mcp.CallToolResult{Content: content, IsError: result.IsError}
}
