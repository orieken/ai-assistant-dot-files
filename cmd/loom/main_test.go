package main

import (
	"testing"

	loom "github.com/orieken/loom"
)

func TestFrameworkContentIsEmbedded(t *testing.T) {
	content, err := loom.FrameworkFS.ReadFile("shared/rules/architecture-guardrails.md")
	if err != nil {
		t.Fatalf("read embedded framework content: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("embedded framework content is empty")
	}
}
