package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var expectedMCPToolNames = []string{
	"analyze_complexity",
	"check_accessibility",
	"check_ubiquitous_language",
	"search_docs",
	"search_ki",
	"verify_dependencies",
}

// TestMCPServeAnswersToolsList spawns the built loom binary, performs an MCP
// initialize + tools/list handshake over stdio, and asserts all six framework
// tool names are present.
func TestMCPServeAnswersToolsList(t *testing.T) {
	binary := buildLoomBinary(t)
	stdin, stdout := startMCPServe(t, binary)
	defer func() { _ = stdin.Close() }()

	names, err := performMCPHandshake(stdin, stdout)
	if err != nil {
		t.Fatalf("MCP handshake: %v", err)
	}
	assertToolNames(t, names)
}

func buildLoomBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "loom")
	build := exec.Command("go", "build", "-o", binary, "github.com/orieken/loom/cmd/loom")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build loom binary: %v\n%s", err, output)
	}
	return binary
}

func startMCPServe(t *testing.T, binary string) (io.WriteCloser, *bufio.Reader) {
	t.Helper()
	serve := exec.Command(binary, "mcp", "serve")
	stdin, err := serve.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := serve.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	if err := serve.Start(); err != nil {
		t.Fatalf("start loom mcp serve: %v", err)
	}
	t.Cleanup(func() {
		_ = stdin.Close()
		waitForExit(serve)
	})
	return stdin, bufio.NewReader(stdout)
}

func waitForExit(serve *exec.Cmd) {
	done := make(chan struct{})
	go func() { _ = serve.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		_ = serve.Process.Kill()
		<-done
	}
}

func performMCPHandshake(stdin io.Writer, stdout *bufio.Reader) ([]string, error) {
	requests := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke-test","version":"0.0.0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	}
	for _, request := range requests {
		if _, err := fmt.Fprintln(stdin, request); err != nil {
			return nil, fmt.Errorf("write request: %w", err)
		}
	}
	return readToolNames(stdout)
}

func readToolNames(stdout *bufio.Reader) ([]string, error) {
	for {
		line, err := stdout.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
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
	if err := json.Unmarshal(line, &response); err != nil || response.ID != 2 {
		return nil, false
	}
	names := make([]string, 0, len(response.Result.Tools))
	for _, tool := range response.Result.Tools {
		names = append(names, tool.Name)
	}
	return names, true
}

func assertToolNames(t *testing.T, names []string) {
	t.Helper()
	present := make(map[string]bool, len(names))
	for _, name := range names {
		present[name] = true
	}
	for _, expected := range expectedMCPToolNames {
		if !present[expected] {
			t.Errorf("tools/list missing %q (got %v)", expected, names)
		}
	}
	if len(names) != len(expectedMCPToolNames) {
		t.Errorf("tools/list returned %d tools, want %d: %v", len(names), len(expectedMCPToolNames), names)
	}
}
