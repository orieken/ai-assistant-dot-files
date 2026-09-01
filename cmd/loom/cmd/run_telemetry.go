package cmd

// Telemetry wiring for `loom run` (roadmap L3.8).
//
// The default is asymmetric on purpose: every run writes its own trace file
// into its workspace, and nothing is sent over the network unless
// OTEL_EXPORTER_OTLP_ENDPOINT says where. Local export leaves nothing, so
// the reason to make network export opt-in does not apply to it — and a
// trace that only exists when someone predicted they would want it cannot
// answer "what did this run cost" about a run that already finished.

import (
	"context"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/telemetry"
	"github.com/spf13/cobra"
)

// flushTimeout bounds shutdown so a wedged collector cannot hang a run that
// has otherwise finished its work.
const flushTimeout = 5 * time.Second

// startTelemetry opens a tracing session for this run and returns the
// shutdown to defer. A nil session yields a no-op shutdown, so callers do
// not branch on whether telemetry is on.
func startTelemetry(cmd *cobra.Command, executor *orchestrator.Executor, workspaceDir string) (func(), error) {
	session, err := telemetry.Start(telemetry.Options{
		Version:   version,
		TraceFile: traceFileFor(workspaceDir),
	})
	if err != nil {
		return nil, err
	}
	executor.WithTracer(session.Tracer())
	return func() { shutdownTelemetry(cmd, session) }, nil
}

// traceFileFor resolves where OTLP/JSON goes: nowhere when telemetry is
// off, the flag's path when given, the run's own workspace otherwise.
func traceFileFor(workspaceDir string) string {
	if runArgs.noTelemetry {
		return ""
	}
	if runArgs.otelFile != "" {
		return runArgs.otelFile
	}
	return telemetry.TraceFileFor(workspaceDir)
}

// shutdownTelemetry flushes on every exit path, the gate halt included — a
// run that stops to ask a human still recorded the stages it completed
// first. A flush failure is reported and never masks the run's own result.
func shutdownTelemetry(cmd *cobra.Command, session *telemetry.Session) {
	ctx, cancel := context.WithTimeout(context.Background(), flushTimeout)
	defer cancel()
	if err := session.Shutdown(ctx); err != nil {
		cmd.PrintErrf("telemetry: flush failed: %v\n", err)
	}
}
