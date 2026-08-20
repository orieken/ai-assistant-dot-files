package frameworkfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveTarget returns a canonical existing project directory.
func ResolveTarget(rawTarget string) (string, error) {
	cleanTarget := filepath.Clean(rawTarget)
	absoluteTarget, err := filepath.Abs(cleanTarget)
	if err != nil {
		return "", fmt.Errorf("resolve target %q: %w", rawTarget, err)
	}
	return validateTarget(absoluteTarget)
}

func validateTarget(target string) (string, error) {
	resolvedTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		return "", fmt.Errorf("resolve target symlinks: %w", err)
	}
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		return "", fmt.Errorf("inspect target: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("target is not a directory: %s", resolvedTarget)
	}
	return resolvedTarget, nil
}

func safeDestination(target, relativePath string) (string, error) {
	if filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("destination must be relative: %s", relativePath)
	}
	cleanPath := filepath.Clean(relativePath)
	if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("destination escapes target: %s", relativePath)
	}
	return filepath.Join(target, cleanPath), nil
}
