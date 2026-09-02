package policy

// Loading a directory of policies.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load reads every *.policy.yaml in dir. A missing directory is not an
// error — most projects have no policies, and that is the safe default.
//
// Every file is parsed, and every failure is reported. Stopping at the
// first would hide the rest, and a policy set is exactly the kind of thing
// someone fixes in one pass.
func Load(dir string) ([]Policy, error) {
	paths, err := policyPaths(dir)
	if err != nil || len(paths) == 0 {
		return nil, err
	}
	policies := make([]Policy, 0, len(paths))
	failures := make([]error, 0)
	for _, path := range paths {
		policy, err := loadOne(path)
		if err != nil {
			failures = append(failures, err)
			continue
		}
		policies = append(policies, policy)
	}
	if len(failures) > 0 {
		return nil, errors.Join(failures...)
	}
	return policies, checkUniqueNames(policies)
}

func policyPaths(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read policy directory %s: %w", dir, err)
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), FileSuffix) {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(paths)
	return paths, nil
}

func loadOne(path string) (Policy, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("read policy %s: %w", path, err)
	}
	return Parse(path, raw)
}

// checkUniqueNames rejects two policies sharing a name. The name is how a
// decision is attributed, so two of them makes an audit record ambiguous
// about which policy fired.
func checkUniqueNames(policies []Policy) error {
	seen := make(map[string]string, len(policies))
	for _, policy := range policies {
		if first, duplicate := seen[policy.Name]; duplicate {
			return fmt.Errorf("policy name %q is used by both %s and %s — names attribute decisions and must be unique",
				policy.Name, first, policy.Source)
		}
		seen[policy.Name] = policy.Source
	}
	return nil
}

// For returns the enabled policies watching a gate, in load order.
func For(policies []Policy, gate GateID) []Policy {
	matching := make([]Policy, 0, len(policies))
	for _, policy := range policies {
		if policy.IsEnabled() && policy.watches(gate) {
			matching = append(matching, policy)
		}
	}
	return matching
}

func (p Policy) watches(gate GateID) bool {
	for _, watched := range p.Matcher.Watched() {
		if watched == gate {
			return true
		}
	}
	return false
}
