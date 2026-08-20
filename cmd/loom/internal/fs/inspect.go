package frameworkfs

import (
	"fmt"
	"os"
)

// PathStatus describes one manifest-owned path without following it for existence checks.
type PathStatus struct {
	Absolute  string
	Exists    bool
	IsSymlink bool
	IsBroken  bool
}

// InspectOwnedPath safely inspects one path recorded in the install manifest.
func InspectOwnedPath(target, relativePath string) (PathStatus, error) {
	destination, err := ownedDestination(target, relativePath)
	if err != nil {
		return PathStatus{}, err
	}
	info, err := os.Lstat(destination)
	if os.IsNotExist(err) {
		return PathStatus{Absolute: destination}, nil
	}
	if err != nil {
		return PathStatus{}, fmt.Errorf("inspect %s: %w", relativePath, err)
	}
	if err := validateResolvedParent(target, destination, relativePath); err != nil {
		return PathStatus{}, err
	}
	status := PathStatus{Absolute: destination, Exists: true, IsSymlink: info.Mode()&os.ModeSymlink != 0}
	status.IsBroken = status.IsSymlink && symlinkIsBroken(destination)
	return status, nil
}

func symlinkIsBroken(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}
