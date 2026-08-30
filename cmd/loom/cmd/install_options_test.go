package cmd

import (
	"testing"

	loom "github.com/orieken/loom"
)

func TestPrepareInstallRejectsUnknownPlatform(t *testing.T) {
	flags := installFlags{target: t.TempDir(), platform: "unknown"}
	if _, err := prepareInstall(flags, "v3.3.14", loom.FrameworkFS); err == nil {
		t.Fatal("expected unknown platform error")
	}
}

func TestPrepareInstallRejectsUnknownStack(t *testing.T) {
	flags := installFlags{target: t.TempDir(), stack: "ruby"}
	if _, err := prepareInstall(flags, "v3.3.14", loom.FrameworkFS); err == nil {
		t.Fatal("expected unknown stack error")
	}
}

func TestPrepareInstallRejectsOutOfRangeLevel(t *testing.T) {
	flags := installFlags{target: t.TempDir(), level: 5}
	if _, err := prepareInstall(flags, "v3.3.14", loom.FrameworkFS); err == nil {
		t.Fatal("expected out-of-range level error")
	}
}

func TestPrepareInstallLevelSelectsCoreRulesPlusStack(t *testing.T) {
	flags := installFlags{target: t.TempDir(), platform: "claude-code", level: 1, stack: "go"}
	request, err := prepareInstall(flags, "v3.3.14", loom.FrameworkFS)
	if err != nil {
		t.Fatalf("prepareInstall failed: %v", err)
	}
	names, err := request.rules.Names(loom.FrameworkFS)
	if err != nil {
		t.Fatalf("Names failed: %v", err)
	}
	want := map[string]bool{
		"approval-gates.md":          true,
		"architecture-guardrails.md": true,
		"design-principles.md":       true,
		"memory-trust-boundary.md":   true,
		"testing-conventions.md":     true,
		"go-conventions.md":          true,
	}
	if len(names) != len(want) {
		t.Fatalf("level 1 + go selected %d rules (%v), want %d", len(names), names, len(want))
	}
	for _, name := range names {
		if !want[name] {
			t.Errorf("unexpected rule %s in level 1 + go selection", name)
		}
	}
}
