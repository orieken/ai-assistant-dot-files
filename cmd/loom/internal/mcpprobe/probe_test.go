package mcpprobe

import (
	"context"
	"testing"
	"time"
)

func TestToolsListFailures(t *testing.T) {
	cases := []struct {
		name    string
		command string
		args    []string
	}{
		{name: "command does not exist", command: "loom-mcpprobe-no-such-command"},
		{name: "server exits without answering", command: "sh", args: []string{"-c", "exit 0"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := ToolsList(ctx, t.TempDir(), testCase.command, testCase.args...); err == nil {
				t.Fatal("ToolsList should fail")
			}
		})
	}
}

func TestToolsListParsesAdvertisedTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	script := `cat >/dev/null & echo '{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"alpha"},{"name":"beta"}]}}'`
	names, err := ToolsList(ctx, t.TempDir(), "sh", "-c", script)
	if err != nil {
		t.Fatalf("ToolsList: %v", err)
	}
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("names = %v, want [alpha beta]", names)
	}
}
