package cmd

// End-to-end checks that a real `loom run` writes a real trace file. The
// binary-driven tests here are invisible to the coverage tool; the
// in-process behaviour they exercise is covered in internal/telemetry.

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/telemetry"
)

func tracePathFor(projectDir string) string {
	return filepath.Join(projectDir, ".claude", "feature-workspace", "user-auth", telemetry.TracesFileName)
}

// spanNamesIn returns every span name in the trace file, in file order.
func spanNamesIn(t *testing.T, path string) []string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read trace file %s: %v", path, err)
	}
	names := make([]string, 0, 8)
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		names = append(names, spanNamesInLine(t, line)...)
	}
	return names
}

func spanNamesInLine(t *testing.T, line string) []string {
	t.Helper()
	var batch struct {
		ResourceSpans []struct {
			ScopeSpans []struct {
				Spans []struct {
					Name string `json:"name"`
				} `json:"spans"`
			} `json:"scopeSpans"`
		} `json:"resourceSpans"`
	}
	if err := json.Unmarshal([]byte(line), &batch); err != nil {
		t.Fatalf("trace line is not valid OTLP/JSON: %v\n%s", err, line)
	}
	names := make([]string, 0, 8)
	for _, resource := range batch.ResourceSpans {
		for _, scope := range resource.ScopeSpans {
			for _, span := range scope.Spans {
				names = append(names, span.Name)
			}
		}
	}
	return names
}

// A run that halts at a gate has still done work, and must still have
// flushed it. This is the exit path most likely to lose telemetry, because
// it is not an error and not a clean completion.
func TestRunWritesTracesAndFlushesAtAGateHalt(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	run.Dir = projectDir
	output, err := run.CombinedOutput()
	if err == nil {
		t.Fatalf("run should halt at confirm-design\n%s", output)
	}

	names := spanNamesIn(t, tracePathFor(projectDir))
	assertContainsSpan(t, names, "loom.run deliver-feature")
	for _, stageID := range stagesBeforeFirstGate {
		assertContainsSpan(t, names, "loom.stage "+stageID)
	}
	// A routed-out stage never ran, so it has no span. Absence here is the
	// assertion that spans track execution, not the plan.
	for _, stageID := range stagesRoutedOut {
		assertNoSpan(t, names, "loom.stage "+stageID)
	}
}

func assertContainsSpan(t *testing.T, names []string, want string) {
	t.Helper()
	for _, name := range names {
		if name == want {
			return
		}
	}
	t.Errorf("no span named %q in trace file (got %v)", want, names)
}

func assertNoSpan(t *testing.T, names []string, unwanted string) {
	t.Helper()
	for _, name := range names {
		if name == unwanted {
			t.Errorf("span %q exists, but that stage was routed out and never ran", unwanted)
		}
	}
}

// --no-telemetry must leave nothing behind, and must not change the run.
func TestRunWithNoTelemetryWritesNoTraceFile(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock", "--no-telemetry")
	run.Dir = projectDir
	if output, err := run.CombinedOutput(); err == nil {
		t.Fatalf("run should halt at confirm-design\n%s", output)
	}

	assertStagesCompleted(t, projectDir, stagesBeforeFirstGate, true)
	if _, err := os.Stat(tracePathFor(projectDir)); !os.IsNotExist(err) {
		t.Errorf("trace file exists despite --no-telemetry (stat error: %v)", err)
	}
}

// --otel-file redirects the trace without disabling it.
func TestRunOtelFileFlagOverridesTheLocation(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)
	custom := filepath.Join(t.TempDir(), "nested", "run.jsonl")

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock", "--otel-file", custom)
	run.Dir = projectDir
	if output, err := run.CombinedOutput(); err == nil {
		t.Fatalf("run should halt at confirm-design\n%s", output)
	}

	assertContainsSpan(t, spanNamesIn(t, custom), "loom.run deliver-feature")
	if _, err := os.Stat(tracePathFor(projectDir)); !os.IsNotExist(err) {
		t.Errorf("default trace file was also written despite --otel-file")
	}
}
