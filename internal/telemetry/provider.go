package telemetry

// Construction and shutdown of the tracer provider (roadmap L3.8).
//
// Two exporters, two different defaults, for one reason: network export
// leaves the machine and local export does not. OTEL_EXPORTER_OTLP_ENDPOINT
// is the standard variable every collector setup already sets, so enabling
// network export needs no loom-specific configuration to learn.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/orieken/loom/internal/orchestrator"
)

// ScopeName identifies loom as the instrumentation scope in every span.
const ScopeName = "github.com/orieken/loom"

// EndpointEnvVar is the standard OTLP endpoint variable. Set it and traces
// are exported over OTLP/HTTP as well as to the run's own file.
const EndpointEnvVar = "OTEL_EXPORTER_OTLP_ENDPOINT"

// Options configures a tracer provider for one run.
type Options struct {
	// Version is the loom build version, recorded as service.version.
	Version string
	// TraceFile is where OTLP/JSON is appended. Empty disables file export.
	TraceFile string
	// Endpoint overrides EndpointEnvVar; empty falls back to the variable,
	// and an empty variable disables network export.
	Endpoint string
}

// Session is a live tracer provider plus the shutdown that flushes it.
type Session struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

// TraceFileFor returns the default trace file for a run's workspace.
func TraceFileFor(workspaceDir string) string {
	return filepath.Join(workspaceDir, TracesFileName)
}

// Start builds a tracer provider from options. It returns a nil Session
// when no exporter is configured — callers pass that straight to
// Executor.WithTracer, which reads nil as "do not trace".
func Start(options Options) (*Session, error) {
	exporters, err := exportersFor(options)
	if err != nil {
		return nil, err
	}
	if len(exporters) == 0 {
		return nil, nil
	}
	provider := sdktrace.NewTracerProvider(append(
		[]sdktrace.TracerProviderOption{sdktrace.WithResource(resourceFor(options.Version))},
		batchers(exporters)...,
	)...)
	return &Session{provider: provider, tracer: provider.Tracer(ScopeName, trace.WithInstrumentationVersion(options.Version))}, nil
}

func batchers(exporters []sdktrace.SpanExporter) []sdktrace.TracerProviderOption {
	options := make([]sdktrace.TracerProviderOption, 0, len(exporters))
	for _, exporter := range exporters {
		options = append(options, sdktrace.WithBatcher(exporter))
	}
	return options
}

func exportersFor(options Options) ([]sdktrace.SpanExporter, error) {
	exporters := make([]sdktrace.SpanExporter, 0, 2)
	if options.TraceFile != "" {
		file, err := newFileExporter(options.TraceFile, scopeFor(options.Version), resourceAttributes(options.Version))
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, file)
	}
	network, err := networkExporter(options)
	if err != nil {
		return nil, err
	}
	if network != nil {
		exporters = append(exporters, network)
	}
	return exporters, nil
}

func networkExporter(options Options) (sdktrace.SpanExporter, error) {
	endpoint := options.Endpoint
	if endpoint == "" {
		endpoint = os.Getenv(EndpointEnvVar)
	}
	if endpoint == "" {
		return nil, nil
	}
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, fmt.Errorf("configure OTLP exporter for %s: %w", endpoint, err)
	}
	return exporter, nil
}

func scopeFor(version string) otlpScope {
	return otlpScope{Name: ScopeName, Version: version}
}

func resourceFor(version string) *resource.Resource {
	return resource.NewWithAttributes(semConvSchemaURL,
		attribute.String("service.name", ServiceName),
		attribute.String("service.version", version),
	)
}

func resourceAttributes(version string) []otlpAttribute {
	return encodeAttributes([]attribute.KeyValue{
		attribute.String("service.name", ServiceName),
		attribute.String("service.version", version),
	})
}

// Tracer adapts the session to the executor's seam. A nil Session yields an
// untyped nil, which the executor reads as tracing disabled.
//
// The interface return type is load-bearing. Returning a concrete *otelTracer
// would make the disabled case a non-nil interface wrapping a nil pointer,
// and the executor's nil check would pass right before the first span
// panicked.
func (s *Session) Tracer() orchestrator.Tracer {
	if s == nil {
		return nil
	}
	return &otelTracer{tracer: s.tracer}
}

// Shutdown flushes every pending span and closes the exporters. It must run
// on every exit path, including the gate halt — a run that stops to ask a
// human still recorded the stages it completed first.
func (s *Session) Shutdown(ctx context.Context) error {
	if s == nil {
		return nil
	}
	if err := s.provider.ForceFlush(ctx); err != nil {
		return errors.Join(err, s.provider.Shutdown(ctx))
	}
	return s.provider.Shutdown(ctx)
}
