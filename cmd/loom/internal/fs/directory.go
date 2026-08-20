package frameworkfs

import (
	"fmt"
	"os"
)

// PrepareDirectory ensures a destination is a real directory, not a symlink.
func (writer *Writer) PrepareDirectory(destination string) error {
	path, err := safeDestination(writer.target, destination)
	if err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return writer.createDirectory(path, destination)
	}
	if err != nil {
		return fmt.Errorf("inspect %s: %w", destination, err)
	}
	if info.IsDir() {
		return nil
	}
	if writer.isDryRun {
		writer.report("would replace " + destination + " with a selective rules directory")
		return nil
	}
	if err := writer.backup(path, destination); err != nil {
		return err
	}
	return writer.createDirectory(path, destination)
}

func (writer *Writer) createDirectory(path, displayPath string) error {
	if writer.isDryRun {
		writer.report("would create " + displayPath)
		return nil
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", displayPath, err)
	}
	return nil
}
