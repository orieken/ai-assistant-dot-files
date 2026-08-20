package frameworkfs

import (
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
)

func copyEmbedded(content iofs.FS, source, destination string) error {
	info, err := iofs.Stat(content, source)
	if err != nil {
		return fmt.Errorf("inspect embedded %s: %w", source, err)
	}
	if !info.IsDir() {
		return copyEmbeddedFile(content, source, destination)
	}
	return walkEmbedded(content, source, destination)
}

func walkEmbedded(content iofs.FS, source, destination string) error {
	return iofs.WalkDir(content, source, func(path string, entry iofs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := embeddedRelative(source, path)
		if err != nil {
			return fmt.Errorf("resolve embedded path %s: %w", path, err)
		}
		return copyEmbeddedEntry(content, path, filepath.Join(destination, filepath.FromSlash(relative)), entry)
	})
}

func copyEmbeddedEntry(content iofs.FS, source, destination string, entry iofs.DirEntry) error {
	if entry.IsDir() {
		return os.MkdirAll(destination, 0o755)
	}
	return copyEmbeddedFile(content, source, destination)
}

func copyEmbeddedFile(content iofs.FS, source, destination string) error {
	data, err := iofs.ReadFile(content, source)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", source, err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create directory for %s: %w", destination, err)
	}
	if err := os.WriteFile(destination, data, embeddedFileMode(source)); err != nil {
		return fmt.Errorf("write %s: %w", destination, err)
	}
	return nil
}

func embeddedFileMode(source string) os.FileMode {
	if strings.HasSuffix(source, ".sh") {
		return 0o755
	}
	return 0o644
}

func embeddedRelative(root, path string) (string, error) {
	if path == root {
		return ".", nil
	}
	if root == "." {
		return path, nil
	}
	prefix := strings.TrimSuffix(root, "/") + "/"
	if !strings.HasPrefix(path, prefix) {
		return "", fmt.Errorf("embedded path %s is outside %s", path, root)
	}
	return strings.TrimPrefix(path, prefix), nil
}
