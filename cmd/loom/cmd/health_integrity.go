package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (check *healthCheck) verifySymlinks() {
	broken := 0
	for _, path := range manifestPaths(check.manifest.Platforms) {
		status := check.paths[path]
		if status.IsSymlink && status.IsBroken {
			check.fail("broken symlink " + path)
			broken++
		}
	}
	if broken == 0 {
		check.output.success("symlinks unbroken")
	}
}

func (check *healthCheck) verifyAgentCounts() {
	expected, err := countFiles(check.content, "shared/agents", ".md", "CHANGELOG.md")
	if err != nil {
		check.fail("count embedded agents: " + err.Error())
		return
	}
	checked := 0
	mismatches := 0
	for _, path := range manifestPaths(check.manifest.Platforms) {
		status := check.paths[path]
		if filepath.Base(path) != "agents" || !status.Exists || status.IsBroken {
			continue
		}
		found, countErr := countInstalledAgents(status.Absolute)
		if countErr != nil {
			check.fail(fmt.Sprintf("count installed agents at %s: %v", path, countErr))
			continue
		}
		checked++
		if found != expected {
			check.warn(fmt.Sprintf("agent count mismatch at %s: embedded has %d, found %d", path, expected, found))
			mismatches++
		}
	}
	if checked > 0 && mismatches == 0 {
		check.output.success(fmt.Sprintf("agent count matches embedded framework (%d)", expected))
	}
}

func countInstalledAgents(directory string) (int, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != "CHANGELOG.md" && strings.HasSuffix(entry.Name(), ".md") {
			count++
		}
	}
	return count, nil
}

func (check *healthCheck) fail(message string) {
	check.failures++
	check.output.failure(message)
}

func (check *healthCheck) warn(message string) {
	check.warnings++
	check.output.warning(message)
}

func (check *healthCheck) finish() error {
	check.output.note()
	check.output.summary(check.failures, check.warnings)
	if check.failures > 0 {
		return errHealthFailed
	}
	return nil
}
