package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/orieken/loom/cmd/loom/internal/platform"
)

type installOutput struct {
	writer   io.Writer
	isDryRun bool
}

func (output installOutput) line(message string) {
	prefix := ""
	if output.isDryRun {
		prefix = "[dry-run] "
	}
	fmt.Fprintln(output.writer, prefix+message)
}

func (output installOutput) action(message string) { output.line("  " + message) }

func (output installOutput) platform(result platform.Result) {
	output.line(fmt.Sprintf("  ✓ %-14s → %s", result.Name, strings.Join(result.Paths, ", ")))
}

func (output installOutput) summary(platformCount, agents, skills, rules int) {
	output.line(fmt.Sprintf("%d platforms, %d agents, %d skills, %d rules installed.", platformCount, agents, skills, rules))
}

func embeddedCounts(content platform.Content, rules platform.RuleSet) (int, int, int, error) {
	agents, err := countFiles(content, "shared/agents", ".md", "CHANGELOG.md")
	if err != nil {
		return 0, 0, 0, err
	}
	skills, err := countDirectories(content, "shared/skills")
	if err != nil {
		return 0, 0, 0, err
	}
	ruleNames, err := rules.Names(content)
	return agents, skills, len(ruleNames), err
}

func countFiles(content platform.Content, directory, suffix, excluded string) (int, error) {
	entries, err := content.ReadDir(directory)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != excluded && strings.HasSuffix(entry.Name(), suffix) {
			count++
		}
	}
	return count, nil
}

func countDirectories(content platform.Content, directory string) (int, error) {
	entries, err := content.ReadDir(directory)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	return count, nil
}

func uniquePaths(results []platform.Result, extras []string) []string {
	seen := make(map[string]bool)
	for _, result := range results {
		for _, path := range result.Paths {
			seen[filepath.ToSlash(path)] = true
		}
	}
	for _, path := range extras {
		seen[filepath.ToSlash(path)] = true
	}
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
