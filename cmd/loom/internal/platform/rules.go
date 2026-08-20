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
}

// RuleSet is a validated language-stack filter.
type RuleSet struct {
	stacks []string
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
	return RuleSet{stacks}, nil
}

// Names returns rule filenames selected by the filter.
func (rules RuleSet) Names(content Content) ([]string, error) {
	if len(rules.stacks) == 0 {
		return allRuleNames(content)
	}
	names := append([]string{}, coreRules...)
	for _, stack := range rules.stacks {
		names = append(names, stackRules[stack])
	}
	return names, nil
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
