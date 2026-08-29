package frameworkfs

import (
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
)

func (writer *Writer) matches(source, destination string) (bool, error) {
	info, err := os.Lstat(destination)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect %s: %w", destination, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return false, nil
	}
	return writer.matchesEmbedded(source, destination, info)
}

func (writer *Writer) matchesEmbedded(source, destination string, info os.FileInfo) (bool, error) {
	embeddedInfo, err := iofs.Stat(writer.content, source)
	if err != nil || info.IsDir() != embeddedInfo.IsDir() {
		return false, err
	}
	if info.IsDir() {
		return writer.matchesDirectory(source, destination)
	}
	data, err := iofs.ReadFile(writer.content, source)
	if err != nil {
		return false, err
	}
	return isSameFile(destination, data)
}

func (writer *Writer) matchesDirectory(source, destination string) (bool, error) {
	isMatch := true
	err := iofs.WalkDir(writer.content, source, func(path string, entry iofs.DirEntry, walkErr error) error {
		if walkErr != nil || !isMatch {
			return walkErr
		}
		relative, relErr := embeddedRelative(source, path)
		if relErr != nil {
			return relErr
		}
		isMatch, relErr = embeddedEntryMatches(writer.content, path, filepath.Join(destination, filepath.FromSlash(relative)), entry)
		return relErr
	})
	return isMatch, err
}

func embeddedEntryMatches(content iofs.FS, source, destination string, entry iofs.DirEntry) (bool, error) {
	if entry.IsDir() {
		info, err := os.Stat(destination)
		if os.IsNotExist(err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return info.IsDir(), nil
	}
	data, err := iofs.ReadFile(content, source)
	if err != nil {
		return false, err
	}
	return isSameFile(destination, data)
}
