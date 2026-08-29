package tools

import (
	"fmt"
	"sort"
	"time"
)

// RetryClass classifies whether a failed tool call is safe to retry
// automatically. Enforcement is middleware's concern (roadmap L2.2/L2.5);
// the registration only declares the class.
type RetryClass string

const (
	// RetryNone marks a tool that must never be retried automatically.
	RetryNone RetryClass = "none"
	// RetryIdempotent marks a read-only tool that is safe to retry.
	RetryIdempotent RetryClass = "idempotent"
)

// PermissionScope names the host access a tool requires.
type PermissionScope string

const (
	// ScopeReadOnly marks a tool that only reads the filesystem or corpus.
	ScopeReadOnly PermissionScope = "read-only"
	// ScopeWorkspaceWrite marks a tool that mutates workspace files.
	ScopeWorkspaceWrite PermissionScope = "workspace-write"
)

// ToolRegistration couples a Tool with the execution metadata the server and
// future middleware (timeout budgets, retry policy, permission checks) need.
type ToolRegistration struct {
	Tool       Tool
	Timeout    time.Duration
	Retry      RetryClass
	Permission PermissionScope
}

// Registry holds tool registrations keyed by tool name. Adding a capability
// is one Register call — no server wiring edits (roadmap L2.4).
type Registry struct {
	entries map[string]ToolRegistration
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{entries: map[string]ToolRegistration{}}
}

// Register adds registration, keyed by its Tool's name. It rejects nil tools,
// empty names, and duplicate names — silent overwrites would let one tool
// shadow another.
func (r *Registry) Register(registration ToolRegistration) error {
	name, err := registrationName(registration)
	if err != nil {
		return err
	}
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("registry: tool %q is already registered", name)
	}
	r.entries[name] = registration
	return nil
}

func registrationName(registration ToolRegistration) (string, error) {
	if registration.Tool == nil {
		return "", fmt.Errorf("registry: registration has a nil Tool")
	}
	name := registration.Tool.Name()
	if name == "" {
		return "", fmt.Errorf("registry: registration has an empty tool name")
	}
	return name, nil
}

// Merge registers every entry of other into r, failing on the first name
// collision. It lets an embedding server combine the framework registry with
// its own tools.
func (r *Registry) Merge(other *Registry) error {
	if other == nil {
		return nil
	}
	for _, registration := range other.All() {
		if err := r.Register(registration); err != nil {
			return err
		}
	}
	return nil
}

// Get returns the registration for name.
func (r *Registry) Get(name string) (ToolRegistration, bool) {
	registration, ok := r.entries[name]
	return registration, ok
}

// All returns every registration sorted by tool name, so iteration order —
// and therefore wire-level tool listing — is deterministic.
func (r *Registry) All() []ToolRegistration {
	names := make([]string, 0, len(r.entries))
	for name := range r.entries {
		names = append(names, name)
	}
	sort.Strings(names)
	registrations := make([]ToolRegistration, 0, len(names))
	for _, name := range names {
		registrations = append(registrations, r.entries[name])
	}
	return registrations
}
