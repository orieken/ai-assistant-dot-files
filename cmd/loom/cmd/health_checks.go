package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"

	frameworkfs "github.com/orieken/loom/cmd/loom/internal/fs"
	"github.com/orieken/loom/cmd/loom/internal/manifest"
	"github.com/orieken/loom/cmd/loom/internal/platform"
)

type healthCheck struct {
	request  healthRequest
	content  platform.Content
	output   healthOutput
	manifest manifest.Manifest
	paths    map[string]frameworkfs.PathStatus
	failures int
	warnings int
}

func (check *healthCheck) verifyManifest() bool {
	installed, err := manifest.Read(check.request.target)
	if err == nil {
		check.manifest = installed
		check.output.success(fmt.Sprintf("manifest found (%s)", installed.Version))
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		check.fail("manifest missing (was loom install run here?)")
	} else {
		check.fail("manifest unreadable: " + err.Error())
	}
	return false
}

func (check *healthCheck) verifyVersion() {
	embeddedVersion, err := readFrameworkVersion(check.content)
	if err != nil {
		check.fail(err.Error())
		return
	}
	if check.manifest.Version != embeddedVersion {
		check.warn(fmt.Sprintf("framework version mismatch: installed %s, embedded %s (upgrade available)", check.manifest.Version, embeddedVersion))
	}
}

func (check *healthCheck) verifyPaths() {
	check.paths = make(map[string]frameworkfs.PathStatus)
	missing := 0
	for _, path := range manifestPaths(check.manifest.Platforms) {
		status, err := frameworkfs.InspectOwnedPath(check.request.target, path)
		if err != nil {
			check.fail(fmt.Sprintf("invalid installed path %s: %v", path, err))
			missing++
			continue
		}
		check.paths[path] = status
		if !status.Exists {
			check.fail("missing installed path " + path)
			missing++
		}
	}
	if missing == 0 {
		check.output.success(fmt.Sprintf("all %d installed paths intact", len(check.paths)))
	}
}

func manifestPaths(records []manifest.PlatformRecord) []string {
	seen := make(map[string]bool)
	for _, record := range records {
		for _, path := range record.Paths {
			seen[path] = true
		}
	}
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
