package frameworkfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (writer *Writer) cachePath(source string) (string, error) {
	return safeDestination(writer.cache, filepath.FromSlash(source))
}

func (writer *Writer) materialize(source, cachePath string) error {
	if _, err := os.Stat(cachePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect cache path: %w", err)
	}
	if err := copyEmbedded(writer.content, source, cachePath); err != nil {
		return fmt.Errorf("extract %s to cache: %w", source, err)
	}
	return nil
}

func (writer *Writer) backup(path, displayPath string) error {
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect %s: %w", displayPath, err)
	}
	backupPath := fmt.Sprintf("%s.bak.%d", path, writer.now().UnixNano())
	if err := os.Rename(path, backupPath); err != nil {
		return fmt.Errorf("backup %s: %w", displayPath, err)
	}
	writer.report("backed up " + displayPath + " -> " + filepath.Base(backupPath))
	return nil
}

func isSameSymlink(path, expectedTarget string) (bool, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect symlink %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}
	actualTarget, err := os.Readlink(path)
	if err != nil {
		return false, fmt.Errorf("read symlink %s: %w", path, err)
	}
	return actualTarget == expectedTarget, nil
}
