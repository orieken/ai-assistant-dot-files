// Package mcpprobe performs a minimal MCP stdio handshake against a server
// command to verify it answers tools/list — the mechanical evidence `loom
// health` needs before crediting a project with a running MCP server
// (roadmap D.4).
package mcpprobe

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

const toolsListRequestID = 2

// ToolsList spawns the command, performs an initialize + tools/list handshake
// over stdio pipes, and returns the advertised tool names. The context bounds
// the whole probe; the spawned process is killed when it expires.
func ToolsList(ctx context.Context, dir, command string, args ...string) ([]string, error) {
	server := exec.CommandContext(ctx, command, args...)
	server.Dir = dir
	stdin, stdout, err := serverPipes(server)
	if err != nil {
		return nil, err
	}
	if err := server.Start(); err != nil {
		return nil, fmt.Errorf("start MCP server: %w", err)
	}
	names, err := handshake(stdin, stdout)
	_ = stdin.Close()
	_ = server.Wait()
	return names, err
}

func serverPipes(server *exec.Cmd) (io.WriteCloser, io.Reader, error) {
	stdin, err := server.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := server.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("stdout pipe: %w", err)
	}
	return stdin, stdout, nil
}

func handshake(stdin io.Writer, stdout io.Reader) ([]string, error) {
	requests := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"loom-health","version":"0.0.0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	}
	for _, request := range requests {
		if _, err := fmt.Fprintln(stdin, request); err != nil {
			return nil, fmt.Errorf("write request: %w", err)
		}
	}
	return readToolNames(bufio.NewReader(stdout))
}

func readToolNames(stdout *bufio.Reader) ([]string, error) {
	for {
		line, err := stdout.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("read tools/list response: %w", err)
		}
		names, isToolsResponse := parseToolsListResponse(line)
		if isToolsResponse {
			return names, nil
		}
	}
}

func parseToolsListResponse(line []byte) ([]string, bool) {
	var response struct {
		ID     int `json:"id"`
		Result struct {
			Tools []struct {
				Name string `json:"name"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal(line, &response); err != nil || response.ID != toolsListRequestID {
		return nil, false
	}
	names := make([]string, 0, len(response.Result.Tools))
	for _, tool := range response.Result.Tools {
		names = append(names, tool.Name)
	}
	return names, true
}
