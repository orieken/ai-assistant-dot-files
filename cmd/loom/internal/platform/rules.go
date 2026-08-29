package platform

import (
	"fmt"
	"sort"
	"strings"
)

var coreRules = []string{
	"approval-gates.md",
	"architecture-guardrails.md",
	"design-principles.md",
	"memory-trust-boundary.md",
	"testing-conventions.md",
}

var stackRules = map[string]string{
	"go":         "go-conventions.md",
	"typescript": "typescript-conventions.md",
	"python":     "python-conventions.md",
	"csharp":     "csharp-conventions.md",
	"java":       "java-conventions.md",
	"kotlin":     "kotlin-conventions.md",
	"swift":      "swift-conventions.md",
	"rust":       "rust-conventions.md",
	// iac is an on-demand module rather than a language stack, but --stack is
	// the opt-in mechanism for every non-core rule module (shared/levels.yaml).
	"iac": "iac-conventions.md",
}

// RuleSet is a validated rule filter — either a language-stack filter
// (--stack) or an explicit filename selection (--level bundles). Empty means
// "all rules", preserving the historic full-install default.
type RuleSet struct {
	stacks []string
	names  []string
}

// ExplicitRuleSet selects exactly the given rule filenames — the --level
// install path, where shared/levels.yaml decides the bundle.
func ExplicitRuleSet(names []string) RuleSet {
	return RuleSet{names: append([]string{}, names...)}
}

// ParseRuleSet validates a comma-separated stack filter.
func ParseRuleSet(raw string) (RuleSet, error) {
	if strings.TrimSpace(raw) == "" {
		return RuleSet{}, nil
	}
	stacks := splitUnique(raw)
	for _, stack := range stacks {
		if stackRules[stack] == "" {
			return RuleSet{}, fmt.Errorf("unknown stack %q", stack)
		}
	}
	return RuleSet{stacks: stacks}, nil
}

// Names returns rule filenames selected by the filter.
func (rules RuleSet) Names(content Content) ([]string, error) {
	if len(rules.names) > 0 {
		return append([]string{}, rules.names...), nil
	}
	if len(rules.stacks) == 0 {
		return allRuleNames(content)
	}
	names := append([]string{}, coreRules...)
	for _, stack := range rules.stacks {
		names = append(names, stackRules[stack])
	}
	return names, nil
}

// Stacks returns the validated stack names selected by --stack.
func (rules RuleSet) Stacks() []string {
	return append([]string{}, rules.stacks...)
}

// isFiltered reports whether the set selects a subset instead of all rules.
func (rules RuleSet) isFiltered() bool {
	return len(rules.stacks) > 0 || len(rules.names) > 0
}

func splitUnique(raw string) []string {
	seen := make(map[string]bool)
	var values []string
	for _, value := range strings.Split(raw, ",") {
		cleanValue := strings.ToLower(strings.TrimSpace(value))
		if cleanValue != "" && !seen[cleanValue] {
			seen[cleanValue] = true
			values = append(values, cleanValue)
		}
	}
	sort.Strings(values)
	return values
}

func allRuleNames(content Content) ([]string, error) {
	entries, err := content.ReadDir("shared/rules")
	if err != nil {
		return nil, fmt.Errorf("read embedded rules: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
