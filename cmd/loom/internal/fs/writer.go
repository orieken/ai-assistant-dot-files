// Package frameworkfs installs embedded framework content safely.
package frameworkfs

import (
	"fmt"
	iofs "io/fs"
	"os"
	"time"
)

// Reporter receives one human-readable filesystem action.
type Reporter func(message string)

// Writer installs embedded content beneath one validated target.
type Writer struct {
	content  iofs.FS
	target   string
	cache    string
	isCopy   bool
	isDryRun bool
	report   Reporter
	now      func() time.Time
}

// NewWriter creates a target-scoped embedded filesystem writer.
func NewWriter(content iofs.FS, target, cache string, isCopy, isDryRun bool, report Reporter) *Writer {
	return &Writer{content, target, cache, isCopy, isDryRun, report, time.Now}
}

// Reporter returns the writer's action reporter for related filesystem work.
func (writer *Writer) Reporter() Reporter { return writer.report }

// Install copies or links one embedded file or directory.
func (writer *Writer) Install(source, destination string) (bool, error) {
	if writer.isCopy {
		return writer.installCopy(source, destination)
	}
	return writer.installLink(source, destination)
}

// Copy installs embedded content regardless of the selected link strategy.
func (writer *Writer) Copy(source, destination string) (bool, error) {
	return writer.installCopy(source, destination)
}

// Write installs generated content beneath the target.
func (writer *Writer) Write(destination string, content []byte) (bool, error) {
	path, err := safeDestination(writer.target, destination)
	if err != nil {
		return false, err
	}
	return writer.writeGenerated(path, destination, content)
}

// CopyIfMissing copies embedded content only when the destination is absent.
func (writer *Writer) CopyIfMissing(source, destination string) (bool, error) {
	path, err := safeDestination(writer.target, destination)
	if err != nil {
		return false, err
	}
	if _, err := os.Lstat(path); err == nil {
		writer.report("skipped " + destination + " (already exists)")
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("inspect %s: %w", destination, err)
	}
	return writer.installCopy(source, destination)
}
