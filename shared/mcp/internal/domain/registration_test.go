package domain_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/shared/mcp/internal/domain"
)

type namedTool struct{ name string }

func (n namedTool) Name() string                  { return n.name }
func (n namedTool) Description() string           { return "test tool" }
func (n namedTool) InputSchema() json.RawMessage  { return json.RawMessage(`{"type":"object"}`) }
func (n namedTool) OutputSchema() json.RawMessage { return nil }

func (n namedTool) Execute(_ context.Context, _ domain.ToolRequest) (*domain.ToolResult, error) {
	return domain.NewTextResult("ok"), nil
}

func registrationFor(name string) domain.ToolRegistration {
	return domain.ToolRegistration{
		Tool:       namedTool{name: name},
		Timeout:    30 * time.Second,
		Retry:      domain.RetryIdempotent,
		Permission: domain.ScopeReadOnly,
	}
}

func TestRegistryRegisterRejectsInvalidRegistrations(t *testing.T) {
	tests := []struct {
		name         string
		registration domain.ToolRegistration
		wantErr      string
	}{
		{"nil tool", domain.ToolRegistration{}, "nil Tool"},
		{"empty name", registrationFor(""), "empty tool name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.NewRegistry().Register(tt.registration)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Register() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestRegistryRegisterRejectsDuplicateNames(t *testing.T) {
	registry := domain.NewRegistry()
	if err := registry.Register(registrationFor("search_ki")); err != nil {
		t.Fatalf("first Register() failed: %v", err)
	}
	err := registry.Register(registrationFor("search_ki"))
	if err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Errorf("duplicate Register() error = %v, want already-registered", err)
	}
}

func TestRegistryGetReturnsRegistration(t *testing.T) {
	registry := domain.NewRegistry()
	if err := registry.Register(registrationFor("search_ki")); err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	registration, ok := registry.Get("search_ki")
	if !ok || registration.Tool.Name() != "search_ki" {
		t.Errorf("Get() = %v, %v; want search_ki registration", registration, ok)
	}
	if _, ok := registry.Get("missing"); ok {
		t.Error("Get(missing) reported ok")
	}
}

func TestRegistryAllIsSortedByName(t *testing.T) {
	registry := domain.NewRegistry()
	for _, name := range []string{"zeta", "alpha", "mid"} {
		if err := registry.Register(registrationFor(name)); err != nil {
			t.Fatalf("Register(%s) failed: %v", name, err)
		}
	}

	var names []string
	for _, registration := range registry.All() {
		names = append(names, registration.Tool.Name())
	}
	want := []string{"alpha", "mid", "zeta"}
	for i, name := range want {
		if names[i] != name {
			t.Fatalf("All() order = %v, want %v", names, want)
		}
	}
}

func TestRegistryMergeCombinesAndDetectsCollisions(t *testing.T) {
	base := domain.NewRegistry()
	extra := domain.NewRegistry()
	if err := base.Register(registrationFor("search_ki")); err != nil {
		t.Fatalf("Register() failed: %v", err)
	}
	if err := extra.Register(registrationFor("custom_tool")); err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	if err := base.Merge(extra); err != nil {
		t.Fatalf("Merge() failed: %v", err)
	}
	if _, ok := base.Get("custom_tool"); !ok {
		t.Error("merged tool not found in base registry")
	}

	if err := base.Merge(extra); err == nil {
		t.Error("expected collision error on second Merge, got nil")
	}
	if err := base.Merge(nil); err != nil {
		t.Errorf("Merge(nil) = %v, want nil", err)
	}
}
