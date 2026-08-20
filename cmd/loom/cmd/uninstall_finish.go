package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	frameworkfs "github.com/orieken/loom/cmd/loom/internal/fs"
	"github.com/orieken/loom/cmd/loom/internal/manifest"
)

func removeOwnedRecords(request uninstallRequest, records []manifest.PlatformRecord, output uninstallOutput) (int, error) {
	removed := 0
	for _, record := range records {
		for _, path := range record.Paths {
			didRemove, err := frameworkfs.RemoveOwnedPath(request.target, path, request.isDryRun)
			if err != nil {
				return removed, fmt.Errorf("remove manifest-owned path for %s: %w", record.Name, err)
			}
			if didRemove {
				removed++
				output.removed(path)
			} else {
				output.missing(path)
			}
		}
	}
	return removed, nil
}

func finishUninstall(request uninstallRequest, installed manifest.Manifest, remaining []manifest.PlatformRecord, removed int, output uninstallOutput) error {
	if request.isDryRun {
		action := "would be deleted"
		if len(remaining) > 0 {
			action = "would be updated"
		}
		output.summary(removed, action)
		return nil
	}
	if len(remaining) > 0 {
		installed.Platforms = remaining
		if err := manifest.Write(request.target, installed); err != nil {
			return err
		}
		output.summary(removed, "updated")
		return nil
	}
	if err := os.Remove(filepath.Join(request.target, manifest.Filename)); err != nil {
		return fmt.Errorf("delete install manifest: %w", err)
	}
	output.summary(removed, "deleted")
	return nil
}
