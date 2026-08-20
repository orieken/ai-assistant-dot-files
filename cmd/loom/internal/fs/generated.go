package frameworkfs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func (writer *Writer) writeGenerated(path, displayPath string, content []byte) (bool, error) {
	if writer.isDryRun {
		writer.report("would write " + displayPath)
		return true, nil
	}
	isSame, err := isSameFile(path, content)
	if err != nil {
		return false, err
	}
	if isSame {
		writer.report("skipped " + displayPath + " (content identical)")
		return true, nil
	}
	if err := writer.backup(path, displayPath); err != nil {
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("create parent for %s: %w", displayPath, err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return false, fmt.Errorf("write %s: %w", displayPath, err)
	}
	writer.report("wrote " + displayPath)
	return true, nil
}

func isSameFile(path string, content []byte) (bool, error) {
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("compare %s: %w", path, err)
	}
	return bytes.Equal(existing, content), nil
}
