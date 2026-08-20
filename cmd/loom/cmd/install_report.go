package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/orieken/loom/cmd/loom/internal/manifest"
	"github.com/orieken/loom/cmd/loom/internal/platform"
)

const sharedOwnership = "loom"

func writeManifest(request installRequest, results []platform.Result, extras []string) error {
	if request.isDryRun {
		return nil
	}
	previous, exists, err := manifest.ReadIfExists(request.target)
	if err != nil {
		return err
	}
	records := platformRecords(results, extras)
	if exists {
		records = mergePlatformRecords(previous.Platforms, records)
	}
	installed := manifest.Manifest{Version: request.frameworkVersion, InstalledAt: time.Now().UTC(), Platforms: records}
	return manifest.Write(request.target, installed)
}

func platformRecords(results []platform.Result, extras []string) []manifest.PlatformRecord {
	records := make([]manifest.PlatformRecord, 0, len(results)+1)
	for _, result := range results {
		records = append(records, manifest.PlatformRecord{Name: result.Name, Paths: mergeStrings(nil, result.Paths)})
	}
	if len(extras) > 0 {
		records = append(records, manifest.PlatformRecord{Name: sharedOwnership, Paths: mergeStrings(nil, extras)})
	}
	return records
}

func mergePlatformRecords(existing, added []manifest.PlatformRecord) []manifest.PlatformRecord {
	pathsByPlatform := make(map[string][]string, len(existing)+len(added))
	for _, record := range append(append([]manifest.PlatformRecord{}, existing...), added...) {
		pathsByPlatform[record.Name] = mergeStrings(pathsByPlatform[record.Name], record.Paths)
	}
	records := make([]manifest.PlatformRecord, 0, len(pathsByPlatform))
	for name, paths := range pathsByPlatform {
		records = append(records, manifest.PlatformRecord{Name: name, Paths: paths})
	}
	sort.Slice(records, func(left, right int) bool { return records[left].Name < records[right].Name })
	return records
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
