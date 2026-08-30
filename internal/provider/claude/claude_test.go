package claude_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/claude"
)

// Leading YAML frontmatter mirrors real shared/agents/*.md files — it is why
// the prompt must travel over stdin, never argv (the claude CLI would parse
// a leading `---` as an option).
const testAgentDefinition = "---\nname: analyst\n---\n# analyst agent\nYou analyze feature specs."

// writeFakeClaude writes an executable shell script standing in for the
// claude CLI, so subprocess plumbing (arg construction, stdout capture,
// timeout kill, exit codes) is tested without the real binary.
func writeFakeClaude(t *testing.T, dir, script string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake-binary tests use POSIX shell scripts")
	}
	path := filepath.Join(dir, "fake-claude")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatalf("write fake claude: %v", err)
	}
	return path
}

func newTestProvider(t *testing.T, binaryPath string) (*claude.Provider, orchestrator.Stage, orchestrator.StageInput) {
	t.Helper()
	agentsDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(agentsDir, "analyst.md"), []byte(testAgentDefinition), 0o644); err != nil {
		t.Fatalf("write agent definition: %v", err)
	}
	provider := claude.New(claude.Config{
		BinaryName: binaryPath,
		AgentsDir:  agentsDir,
		Logger:     slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
	})
	workspace := t.TempDir()
	stage := orchestrator.Stage{ID: "analyst", Agent: "analyst", Timeout: 5 * time.Second}
	input := orchestrator.StageInput{SpecPath: "/specs/user-auth.md", WorkspaceDir: workspace}
	return provider, stage, input
}

func TestInvokeCapturesStdoutToArtifact(t *testing.T) {
	binary := writeFakeClaude(t, t.TempDir(), `echo "# analysis output"`)
	provider, stage, input := newTestProvider(t, binary)

	output, err := provider.Invoke(context.Background(), stage, input)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	wantPath := filepath.Join(input.WorkspaceDir, "analyst.md")
	if output.ArtifactPath != wantPath {
		t.Errorf("artifact path = %q, want %q", output.ArtifactPath, wantPath)
	}
	content, err := os.ReadFile(output.ArtifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if strings.TrimSpace(string(content)) != "# analysis output" {
		t.Errorf("artifact content = %q, want the fake agent's stdout", content)
	}
}

func TestInvokeSendsPromptOverStdinNotArgv(t *testing.T) {
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args.txt")
	stdinFile := filepath.Join(dir, "stdin.txt")
	// The fake dumps argv and stdin separately so the test can assert both.
	binary := writeFakeClaude(t, dir, `printf '%s\n' "$@" > `+argsFile+`
cat > `+stdinFile)
	provider, stage, input := newTestProvider(t, binary)

	if _, err := provider.Invoke(context.Background(), stage, input); err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	assertPromptPlumbing(t, argsFile, stdinFile, input)
}

func assertPromptPlumbing(t *testing.T, argsFile, stdinFile string, input orchestrator.StageInput) {
	t.Helper()
	args := readFile(t, argsFile)
	if strings.TrimSpace(args) != "-p" {
		t.Errorf("argv = %q, want only -p — the prompt must travel over stdin (agent frontmatter starts with ---)", args)
	}
	prompt := readFile(t, stdinFile)
	for _, want := range []string{testAgentDefinition, input.SpecPath, input.WorkspaceDir} {
		if !strings.Contains(prompt, want) {
			t.Errorf("stdin prompt missing %q", want)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(raw)
}

func TestInvokeFailsClearlyWhenBinaryMissing(t *testing.T) {
	provider, stage, input := newTestProvider(t, "definitely-not-a-real-binary-xyz")

	_, err := provider.Invoke(context.Background(), stage, input)
	if err == nil {
		t.Fatal("Invoke succeeded with a missing binary; must fail, never fall back to mock")
	}
	for _, want := range []string{"not found", "npm install -g @anthropic-ai/claude-code"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("missing-binary error %q lacks remediation text %q", err, want)
		}
	}
}

func TestInvokeFailsWhenAgentDefinitionMissing(t *testing.T) {
	binary := writeFakeClaude(t, t.TempDir(), `echo ok`)
	provider, stage, input := newTestProvider(t, binary)
	stage.Agent = "no-such-agent"

	_, err := provider.Invoke(context.Background(), stage, input)
	if err == nil || !strings.Contains(err.Error(), "no-such-agent.md") {
		t.Fatalf("Invoke error = %v, want missing agent definition naming the path", err)
	}
}

func TestInvokeKillsSubprocessOnTimeout(t *testing.T) {
	binary := writeFakeClaude(t, t.TempDir(), `sleep 30`)
	provider, stage, input := newTestProvider(t, binary)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := provider.Invoke(ctx, stage, input)
	if err == nil {
		t.Fatal("Invoke succeeded despite the deadline")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Invoke error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("subprocess was not killed on deadline: took %v", elapsed)
	}
}

func TestInvokeReportsNonzeroExitWithStderr(t *testing.T) {
	binary := writeFakeClaude(t, t.TempDir(), `echo "agent blew up" >&2; exit 3`)
	provider, stage, input := newTestProvider(t, binary)

	_, err := provider.Invoke(context.Background(), stage, input)
	if err == nil {
		t.Fatal("Invoke succeeded despite nonzero exit")
	}
	if !strings.Contains(err.Error(), "agent blew up") {
		t.Errorf("error %q does not surface the subprocess stderr", err)
	}
}
