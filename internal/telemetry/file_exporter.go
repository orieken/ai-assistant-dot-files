package telemetry

// The offline exporter. Every run writes its trace beside its own state,
// with no collector configured and nothing leaving the machine — which is
// why it is on by default while the network exporter is not. The concern
// that makes a phone-home default rude is egress, and a file next to
// run-state.json has none.
//
// Writing it always is what makes a run's cost answerable afterwards. Made
// opt-in, the question "what did this run cost" could only be answered by
// someone who had predicted they would ask.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TracesFileName is the run's trace log, beside run-state.json and
// run-events.jsonl. One line per exported batch, each a complete OTLP/JSON
// request body.
const TracesFileName = "traces.jsonl"

// fileExporter appends OTLP/JSON batches to a file.
type fileExporter struct {
	mutex    sync.Mutex
	file     *os.File
	scope    otlpScope
	resource []otlpAttribute
}

func newFileExporter(path string, scope otlpScope, resource []otlpAttribute) (*fileExporter, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create trace directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open trace file: %w", err)
	}
	return &fileExporter{file: file, scope: scope, resource: resource}, nil
}

// ExportSpans writes one batch as one line. Append-only with one write per
// line is the whole atomicity model, exactly as the run timeline's is: the
// file is never read back, rewritten, or truncated.
func (e *fileExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	if len(spans) == 0 || ctx.Err() != nil {
		return ctx.Err()
	}
	line, err := encodeExport(spans, e.scope, e.resource)
	if err != nil {
		return fmt.Errorf("encode spans: %w", err)
	}
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, err := e.file.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("write trace file: %w", err)
	}
	return nil
}

// Shutdown closes the file. Calling it twice is safe, because the SDK's
// shutdown path and an explicit one can both reach it.
func (e *fileExporter) Shutdown(context.Context) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.file == nil {
		return nil
	}
	err := e.file.Close()
	e.file = nil
	return err
}
