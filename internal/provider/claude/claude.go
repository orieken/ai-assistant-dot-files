// Package claude is the real Provider for the executor skeleton (roadmap
// M0.4 part 2): it invokes each stage's agent by spawning the `claude` CLI
// headless (`claude -p <prompt>`) as a subprocess. The prompt is built from
// the agent's shared/agents/<agent>.md definition plus the feature spec
// path; stdout is captured to the stage's artifact path. Timeouts arrive
// via ctx (the executor applies the stage's timeout) and are enforced by
// exec.CommandContext, which kills the subprocess when ctx ends. If the
// claude binary is absent the stage FAILS with a remediation message — it
// never silently falls back to the mock provider.
package claude

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
)

// DefaultBinaryName is the CLI the provider spawns.
const DefaultBinaryName = "claude"

// Config wires the provider to its environment. All fields are explicit so
// tests can point BinaryName at a scripted fake and AgentsDir at fixtures.
type Config struct {
	// BinaryName is the claude CLI to spawn; empty means DefaultBinaryName.
	BinaryName string
	// AgentsDir holds the agent definitions (shared/agents or a .claude/agents symlink).
	AgentsDir string
	// Logger receives structured logs; nil means JSON to os.Stderr. Logs
	// never go to stdout (that belongs to the user).
	Logger *slog.Logger
}

// Provider spawns `claude -p` per stage.
type Provider struct {
	binaryName string
	agentsDir  string
	logger     *slog.Logger
}

// New builds a Provider from config, defaulting the binary name and logger.
func New(config Config) *Provider {
	binaryName := config.BinaryName
	if binaryName == "" {
		binaryName = DefaultBinaryName
	}
	logger := config.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	return &Provider{binaryName: binaryName, agentsDir: config.AgentsDir, logger: logger}
}

// Invoke runs one stage's agent headless and writes its output to the
// stage's artifact path inside the workspace.
func (p *Provider) Invoke(ctx context.Context, stage orchestrator.Stage, input orchestrator.StageInput) (orchestrator.StageOutput, error) {
	binaryPath, err := exec.LookPath(p.binaryName)
	if err != nil {
		return orchestrator.StageOutput{}, fmt.Errorf(
			"claude binary %q not found on PATH — install it (npm install -g @anthropic-ai/claude-code) or run with --provider mock for a dry run: %w",
			p.binaryName, err)
	}
	prompt, err := p.buildPrompt(stage, input)
	if err != nil {
		return orchestrator.StageOutput{}, err
	}
	return p.runSubprocess(ctx, binaryPath, prompt, stage, input)
}

// buildPrompt composes the stage prompt: the full agent definition, then the
// delivery task naming the spec and workspace so the agent reads them itself.
func (p *Provider) buildPrompt(stage orchestrator.Stage, input orchestrator.StageInput) (string, error) {
	definitionPath := filepath.Join(p.agentsDir, stage.Agent+".md")
	definition, err := os.ReadFile(definitionPath)
	if err != nil {
		return "", fmt.Errorf("agent definition for stage %q not found at %s: %w", stage.ID, definitionPath, err)
	}
	prompt := fmt.Sprintf(
		"%s\n\n---\n\nAct exactly as the agent defined above.\nFeature spec: %s\nWorkspace directory (read prior stage artifacts here): %s\nProduce your complete markdown artifact on stdout and nothing else.\n",
		definition, input.SpecPath, input.WorkspaceDir)
	if stage.StateKind == "" {
		return prompt, nil
	}
	instruction, err := typedInstruction(stage, input)
	if err != nil {
		return "", err
	}
	return prompt + instruction, nil
}

// runSubprocess executes `claude -p <prompt>` with stdout captured to the
// stage's artifact path. exec.CommandContext kills the subprocess when ctx
// ends — that is how stage timeouts and SIGINT reach it.
func (p *Provider) runSubprocess(ctx context.Context, binaryPath, prompt string, stage orchestrator.Stage, input orchestrator.StageInput) (orchestrator.StageOutput, error) {
	stdout, artifactPath, err := p.captureTarget(stage, input)
	if err != nil {
		return orchestrator.StageOutput{}, err
	}
	defer stdout.close()

	stderr := &bytes.Buffer{}
	// The prompt travels over stdin, not argv: agent definitions begin with
	// `---` YAML frontmatter, which the claude CLI's argument parser would
	// treat as an option — and stdin also sidesteps ARG_MAX for large
	// definitions. `claude -p` reads the prompt from stdin when piped.
	command := exec.CommandContext(ctx, binaryPath, "-p")
	command.Stdin = strings.NewReader(prompt)
	command.Stdout = stdout.writer
	command.Stderr = stderr

	started := time.Now()
	p.logger.Info("stage.started", "stage", stage.ID, "agent", stage.Agent, "binary", binaryPath)
	runErr := command.Run()
	p.logger.Info("stage.finished", "stage", stage.ID, "durationMs", time.Since(started).Milliseconds(), "success", runErr == nil)
	return p.finish(ctx, runErr, stage, artifactPath, stdout, stderr)
}

// captureTarget decides where the agent's stdout goes: a typed stage's
// response is held in memory so it can be parsed and validated, while an
// untyped stage still streams straight into its markdown artifact.
func (p *Provider) captureTarget(stage orchestrator.Stage, input orchestrator.StageInput) (*capture, string, error) {
	if stage.StateKind != "" {
		buffer := &bytes.Buffer{}
		return &capture{writer: buffer, buffer: buffer}, "", nil
	}
	artifactPath := filepath.Join(input.WorkspaceDir, stage.ID+".md")
	file, err := os.Create(artifactPath)
	if err != nil {
		return nil, "", fmt.Errorf("create artifact file: %w", err)
	}
	return &capture{writer: file, file: file}, artifactPath, nil
}

func (p *Provider) finish(ctx context.Context, waitErr error, stage orchestrator.Stage, artifactPath string, stdout *capture, stderr *bytes.Buffer) (orchestrator.StageOutput, error) {
	if ctx.Err() != nil {
		return orchestrator.StageOutput{}, fmt.Errorf("stage %q subprocess terminated: %w", stage.ID, ctx.Err())
	}
	if waitErr != nil {
		return orchestrator.StageOutput{}, fmt.Errorf("stage %q agent exited with error: %w — stderr: %s",
			stage.ID, waitErr, truncate(stderr.String(), 2000))
	}
	if stdout.buffer == nil {
		return orchestrator.StageOutput{ArtifactPath: artifactPath}, nil
	}
	payload, err := extractJSON(stdout.buffer.Bytes())
	if err != nil {
		return orchestrator.StageOutput{}, fmt.Errorf("stage %q: %w", stage.ID, err)
	}
	return orchestrator.StageOutput{Payload: payload}, nil
}

// capture is where a stage's stdout lands: an artifact file for markdown
// stages, an in-memory buffer for typed ones.
type capture struct {
	writer io.Writer
	file   *os.File
	buffer *bytes.Buffer
}

func (c *capture) close() {
	if c.file != nil {
		_ = c.file.Close()
	}
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "…(truncated)"
}
