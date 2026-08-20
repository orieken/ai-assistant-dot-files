package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

type marker struct {
	platform string
	paths    []markerPath
}

type markerPath struct {
	path        string
	isDirectory bool
}

var detectionMarkers = []marker{
	{"claude-code", []markerPath{{".claude", true}, {"CLAUDE.md", false}}},
	{"cursor", []markerPath{{".cursor", true}}},
	{"windsurf", []markerPath{{".windsurf", true}}},
	{"github-copilot", []markerPath{{".github/copilot-instructions.md", false}, {".github", true}}},
	{"gemini", []markerPath{{".gemini", true}}},
	{"openai-codex", []markerPath{{".codex", true}, {"AGENTS.md", false}}},
	{"jetbrains", []markerPath{{".aiassistant", true}, {".junie", true}}},
	{"roo-code", []markerPath{{".roo", true}, {".roomodes", false}}},
	{"cline", []markerPath{{".cline", true}, {".clinerules", true}}},
}

// Detect returns platforms whose markers exist under target.
func Detect(target string) ([]string, error) {
	detected := make([]string, 0, len(detectionMarkers))
	for _, candidate := range detectionMarkers {
		found, err := hasMarker(target, candidate.paths)
		if err != nil {
			return nil, fmt.Errorf("detect %s: %w", candidate.platform, err)
		}
		if found {
			detected = append(detected, candidate.platform)
		}
	}
	return detected, nil
}

func hasMarker(target string, paths []markerPath) (bool, error) {
	for _, marker := range paths {
		info, err := os.Stat(filepath.Join(target, filepath.FromSlash(marker.path)))
		if err == nil {
			if info.IsDir() == marker.isDirectory {
				return true, nil
			}
			continue
		}
		if !os.IsNotExist(err) {
			return false, err
		}
	}
	return false, nil
}
