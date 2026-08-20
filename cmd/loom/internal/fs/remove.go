package frameworkfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const manifestFilename = ".loom-manifest.json"

// RemoveOwnedPath removes one manifest-owned path without escaping the target.
func RemoveOwnedPath(target, relativePath string, isDryRun bool) (bool, error) {
	destination, err := ownedDestination(target, relativePath)
	if err != nil {
		return false, err
	}
	if _, err := os.Lstat(destination); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("inspect %s: %w", relativePath, err)
	}
	if err := validateResolvedParent(target, destination, relativePath); err != nil {
		return false, err
	}
	if isDryRun {
		return true, nil
	}
	if err := os.RemoveAll(destination); err != nil {
		return false, fmt.Errorf("remove %s: %w", relativePath, err)
	}
	return true, nil
}

func ownedDestination(target, relativePath string) (string, error) {
	cleanPath := filepath.Clean(relativePath)
	if cleanPath == "." || cleanPath == manifestFilename {
		return "", fmt.Errorf("refusing unsafe owned path %q", relativePath)
	}
	return safeDestination(target, cleanPath)
}

func validateResolvedParent(target, destination, relativePath string) error {
	resolvedParent, err := filepath.EvalSymlinks(filepath.Dir(destination))
	if err != nil {
		return fmt.Errorf("resolve parent for %s: %w", relativePath, err)
	}
	resolvedTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		return fmt.Errorf("resolve target for %s: %w", relativePath, err)
	}
	if !isWithinTarget(resolvedTarget, resolvedParent) {
		return fmt.Errorf("owned path escapes target through a symlink: %s", relativePath)
	}
	return nil
}

func isWithinTarget(target, path string) bool {
	relative, err := filepath.Rel(target, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
