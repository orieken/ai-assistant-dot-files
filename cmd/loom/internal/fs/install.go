package frameworkfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (writer *Writer) installCopy(source, destination string) (bool, error) {
	path, err := safeDestination(writer.target, destination)
	if err != nil {
		return false, err
	}
	if writer.isDryRun {
		writer.report("would copy " + source + " -> " + destination)
		return true, nil
	}
	if isMatch, matchErr := writer.matches(source, path); matchErr != nil {
		return false, matchErr
	} else if isMatch {
		writer.report("skipped " + destination + " (content identical)")
		return true, nil
	}
	return true, writer.replaceWithCopy(source, path, destination)
}

func (writer *Writer) installLink(source, destination string) (bool, error) {
	path, err := safeDestination(writer.target, destination)
	if err != nil {
		return false, err
	}
	cachePath, err := writer.cachePath(source)
	if err != nil {
		return false, err
	}
	if writer.isDryRun {
		writer.report("would symlink " + destination + " -> " + cachePath)
		return true, nil
	}
	isSame, err := isSameSymlink(path, cachePath)
	if err != nil {
		return false, err
	}
	if isSame {
		writer.report("skipped " + destination + " (already linked)")
		return true, nil
	}
	return true, writer.replaceWithLink(source, path, cachePath, destination)
}

func (writer *Writer) replaceWithCopy(source, path, displayPath string) error {
	if err := writer.backup(path, displayPath); err != nil {
		return err
	}
	if err := copyEmbedded(writer.content, source, path); err != nil {
		return fmt.Errorf("copy %s: %w", source, err)
	}
	writer.report("copied " + displayPath)
	return nil
}

func (writer *Writer) replaceWithLink(source, path, cachePath, displayPath string) error {
	if err := writer.materialize(source, cachePath); err != nil {
		return err
	}
	if err := writer.backup(path, displayPath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent for %s: %w", displayPath, err)
	}
	if err := os.Symlink(cachePath, path); err != nil {
		return fmt.Errorf("link %s: %w", displayPath, err)
	}
	writer.report("linked " + displayPath + " -> " + cachePath)
	return nil
}
