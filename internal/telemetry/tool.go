package telemetry

// Tool-call spans for the MCP server (roadmap L3.8, phase C).
//
// Two properties this file owes the rest of the system: a tool call's
// arguments never leak a secret into a span, and never blow a span up with
// an unbounded payload. Both are enforced here rather than at call sites,
// because a call site that forgets is exactly the failure being prevented.

import (
	"context"
	"sort"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// AttributeValueLimit caps how much of any single argument or result
// reaches a span. Values longer than this are truncated with a marker, so a
// reader can tell truncation from a genuinely short value.
const AttributeValueLimit = 512

// TruncationMarker is appended to any value the limit cut short.
const TruncationMarker = "…(truncated)"

// RedactedPlaceholder replaces a value whose key looks secret. The key is
// kept: knowing a token was passed is useful, knowing which token is not.
const RedactedPlaceholder = "[redacted]"

// secretKeyParts are substrings that make an argument name secret-shaped.
// Matching on the name rather than the value is deliberate — value-shaped
// detection misses the secrets that do not look like secrets, and a name
// match costs nothing when it is wrong.
var secretKeyParts = []string{
	"token", "secret", "password", "passwd", "credential",
	"apikey", "api_key", "authorization", "auth", "private_key", "session",
}

// ToolCall describes one MCP tool invocation. Arguments arrive already
// stringified so this package's signatures stay free of `any`, and so the
// caller owns how its own types render.
type ToolCall struct {
	Name      string
	Arguments map[string]string
}

// ToolResult is how a tool call ended.
type ToolResult struct {
	// Preview is the result text, truncated and redacted like any argument.
	Preview string
	// IsError marks a tool that ran and reported failure, as distinct from
	// one that returned an error to the transport.
	IsError bool
	Err     error
}

// StartTool opens a span for one tool call beneath whatever is in ctx —
// the propagated stage span when TRACEPARENT survived the hop from
// `loom run`, and a fresh trace when it did not.
//
// It hangs off Session rather than off a tracer interface because a nil
// Session must be usable: the MCP server runs untraced far more often than
// traced, and its call sites should not each branch on that.
func (s *Session) StartTool(ctx context.Context, call ToolCall) (context.Context, *ToolSpan) {
	if s == nil {
		return ctx, nil
	}
	ctx, span := s.tracer.Start(ctx, "loom.tool "+call.Name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(toolAttributes(call)...))
	return ctx, &ToolSpan{span: span}
}

// TraceIDs returns the trace and span IDs in ctx, or empty strings when it
// carries no recording span. It exists so log lines can carry correlation
// IDs without the logging package learning what OpenTelemetry is.
func TraceIDs(ctx context.Context) (traceID, spanID string) {
	context := trace.SpanContextFromContext(ctx)
	if !context.IsValid() {
		return "", ""
	}
	return context.TraceID().String(), context.SpanID().String()
}

// ToolSpan is an open tool-call span.
type ToolSpan struct {
	span trace.Span
}

// End closes the span with the tool's result.
func (s *ToolSpan) End(result ToolResult) {
	if s == nil {
		return
	}
	s.span.SetAttributes(
		attribute.String("loom.tool.result", safeValue(result.Preview)),
		attribute.Bool("loom.tool.is_error", result.IsError),
	)
	s.recordToolStatus(result)
	s.span.End()
}

func (s *ToolSpan) recordToolStatus(result ToolResult) {
	if result.Err == nil && !result.IsError {
		s.span.SetStatus(codes.Ok, "")
		return
	}
	if result.Err != nil {
		s.span.RecordError(result.Err)
	}
	s.span.SetStatus(codes.Error, "tool call failed")
}

func toolAttributes(call ToolCall) []attribute.KeyValue {
	attributes := []attribute.KeyValue{attribute.String("loom.tool.name", call.Name)}
	for _, key := range sortedKeys(call.Arguments) {
		attributes = append(attributes,
			attribute.String("loom.tool.arg."+key, argumentValue(key, call.Arguments[key])))
	}
	return attributes
}

// sortedKeys makes the attribute order deterministic, so two identical
// calls produce identical spans and a diff of two traces means something.
func sortedKeys(arguments map[string]string) []string {
	keys := make([]string, 0, len(arguments))
	for key := range arguments {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func argumentValue(key, value string) string {
	if IsSecretKey(key) {
		return RedactedPlaceholder
	}
	return safeValue(value)
}

// IsSecretKey reports whether an argument name looks like it carries a
// credential. Exported so the same judgement is available to anything else
// that has to decide what not to record.
func IsSecretKey(key string) bool {
	lowered := strings.ToLower(key)
	for _, part := range secretKeyParts {
		if strings.Contains(lowered, part) {
			return true
		}
	}
	return false
}

// safeValue truncates on a rune boundary, so a cut through a multi-byte
// character cannot put invalid UTF-8 into a span.
func safeValue(value string) string {
	if len(value) <= AttributeValueLimit {
		return value
	}
	runes := []rune(value)
	limit := AttributeValueLimit
	if limit > len(runes) {
		limit = len(runes)
	}
	return string(runes[:limit]) + TruncationMarker
}
