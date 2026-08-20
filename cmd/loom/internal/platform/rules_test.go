package platform

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestParseRuleSetIncludesCoreAndSelectedStacks(t *testing.T) {
	rules, err := ParseRuleSet("python, go,go")
	if err != nil {
		t.Fatalf("parse rule set: %v", err)
	}
	want := append(append([]string{}, coreRules...), "go-conventions.md", "python-conventions.md")
	names, err := rules.Names(fstest.MapFS{})
	if err != nil {
		t.Fatalf("list filtered rules: %v", err)
	}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("rules = %v, want %v", names, want)
	}
}

func TestParseRuleSetRejectsUnknownStack(t *testing.T) {
	if _, err := ParseRuleSet("go,ruby"); err == nil {
		t.Fatal("expected unknown stack error")
	}
}
