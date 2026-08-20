package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/orieken/loom/cmd/loom/internal/manifest"
	"github.com/orieken/loom/cmd/loom/internal/platform"
)

func writeManifest(request installRequest, paths []string) error {
	if request.isDryRun {
		return nil
	}
	previous, exists, err := manifest.ReadIfExists(request.target)
	if err != nil {
		return err
	}
	if exists {
		paths = mergeStrings(previous.Paths, paths)
	}
	platforms := request.platforms
	if exists {
		platforms = mergeStrings(previous.Platforms, platforms)
	}
	installed := manifest.Manifest{Version: frameworkVersion, InstalledAt: time.Now().UTC(), Platforms: platforms, Paths: paths}
	return manifest.Write(request.target, installed)
}

func mergeStrings(existing, added []string) []string {
	seen := make(map[string]bool, len(existing)+len(added))
	for _, value := range append(append([]string{}, existing...), added...) {
		seen[value] = true
	}
	merged := make([]string, 0, len(seen))
	for value := range seen {
		merged = append(merged, value)
	}
	sort.Strings(merged)
	return merged
}

func reportInstall(request installRequest, results []platform.Result, content platform.Content, output installOutput) error {
	agents, skills, rules, err := embeddedCounts(content, request.rules)
	if err != nil {
		return fmt.Errorf("count embedded framework content: %w", err)
	}
	output.summary(len(results), agents, skills, rules)
	if !output.isDryRun {
		output.line("Manifest written to " + manifest.Filename)
	}
	return nil
}
