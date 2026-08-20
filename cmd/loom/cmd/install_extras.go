package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	frameworkfs "github.com/orieken/loom/cmd/loom/internal/fs"
	"github.com/orieken/loom/cmd/loom/internal/platform"
)

func installConfigs(content platform.Content, files *frameworkfs.Writer) ([]string, error) {
	entries, err := content.ReadDir("shared/configs")
	if err != nil {
		return nil, fmt.Errorf("read embedded configs: %w", err)
	}
	var paths []string
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "README.md" {
			continue
		}
		installed, copyErr := files.CopyIfMissing("shared/configs/"+entry.Name(), entry.Name())
		if copyErr != nil {
			return nil, copyErr
		}
		if installed {
			paths = append(paths, entry.Name())
		}
	}
	return paths, nil
}

func installMCP(content platform.Content, target, cache string, isDryRun bool, report frameworkfs.Reporter) ([]string, error) {
	destination := sanitizeProjectName(filepath.Base(target)) + "-mcp"
	files := frameworkfs.NewWriter(content, target, cache, true, isDryRun, report)
	installed, err := files.CopyIfMissing(".", destination)
	if err != nil || !installed {
		return nil, err
	}
	return []string{destination}, nil
}

func sanitizeProjectName(name string) string {
	var result strings.Builder
	for _, character := range strings.ToLower(name) {
		if isProjectNameCharacter(character) {
			result.WriteRune(character)
		} else {
			result.WriteRune('-')
		}
	}
	cleanName := strings.Trim(result.String(), "-")
	if cleanName == "" {
		return "project"
	}
	return cleanName
}

func isProjectNameCharacter(character rune) bool {
	return character >= 'a' && character <= 'z' || character >= '0' && character <= '9' || character == '-' || character == '_'
}
