package telemetry

// Trace context across process boundaries (roadmap L3.8, phase C).
//
// The chain is `loom run` → `claude -p` → `loom mcp serve`, and loom spawns
// only the first hop: the claude CLI decides whether to spawn an MCP server
// and what environment to give it. MCP's own protocol carries no trace
// context, so there is no in-band channel to use.
//
// What is left is the W3C TRACEPARENT environment variable — the convention
// otel-cli and CI instrumentation already use — which is inherited through
// exec unless something deliberately clears it. That makes propagation
// BEST-EFFORT and it is described that way everywhere: when the variable
// survives the hop a tool call lands under the stage that caused it, and
// when it does not the tool call starts a clean trace of its own. It never
// invents a parent, and it never fails a call for want of one.

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/propagation"
)

// TraceParentEnvVar is the W3C trace-context environment variable.
const TraceParentEnvVar = "TRACEPARENT"

// propagator is the W3C trace-context format, used in both directions.
var propagator = propagation.TraceContext{}

// TraceParentEnv renders the span in ctx as environment entries for a child
// process, or nil when ctx carries no recording span. Callers append the
// result to a subprocess environment.
func TraceParentEnv(ctx context.Context) []string {
	carrier := propagation.MapCarrier{}
	propagator.Inject(ctx, carrier)
	value := carrier.Get("traceparent")
	if value == "" {
		return nil
	}
	return []string{TraceParentEnvVar + "=" + value}
}

// ContextFromEnvironment adopts a TRACEPARENT inherited from a parent
// process. With no variable set, or an unparseable one, it returns ctx
// unchanged — a malformed value from somewhere upstream must not stop a
// tool from running.
func ContextFromEnvironment(ctx context.Context) context.Context {
	value := os.Getenv(TraceParentEnvVar)
	if value == "" {
		return ctx
	}
	return propagator.Extract(ctx, propagation.MapCarrier{"traceparent": value})
}
