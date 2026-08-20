// Package manifest records the paths owned by a loom installation.
package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const Filename = ".loom-manifest.json"

// Manifest describes one completed loom installation.
type Manifest struct {
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installedAt"`
	Platforms   []string  `json:"platforms"`
	Paths       []string  `json:"paths"`
}

// Read loads the loom manifest from a target directory.
func Read(target string) (Manifest, error) {
	content, err := os.ReadFile(filepath.Join(target, Filename))
	if err != nil {
		return Manifest{}, fmt.Errorf("read install manifest: %w", err)
	}
	var installed Manifest
	if err := json.Unmarshal(content, &installed); err != nil {
		return Manifest{}, fmt.Errorf("decode install manifest: %w", err)
	}
	return installed, nil
}

// ReadIfExists loads a manifest when one is present.
func ReadIfExists(target string) (Manifest, bool, error) {
	installed, err := Read(target)
	if err == nil {
		return installed, true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return Manifest{}, false, nil
	}
	return Manifest{}, false, err
}

// Write persists a manifest atomically in the target directory.
func Write(target string, installed Manifest) error {
	content, err := json.MarshalIndent(installed, "", "  ")
	if err != nil {
		return fmt.Errorf("encode install manifest: %w", err)
	}
	destination := filepath.Join(target, Filename)
	if err := os.WriteFile(destination, append(content, '\n'), 0o644); err != nil {
		return fmt.Errorf("write install manifest: %w", err)
	}
	return nil
}
