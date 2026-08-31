// Command gen-schemas writes the JSON Schema for every typed pipeline stage
// into shared/schemas/pipeline/. Run it from the repo root after changing
// any struct in internal/state; a test fails if the committed files drift
// from what the structs generate.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orieken/loom/internal/state"
)

func main() {
	if err := generate(state.SchemaDir); err != nil {
		fmt.Fprintln(os.Stderr, "gen-schemas:", err)
		os.Exit(1)
	}
}

func generate(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create schema directory: %w", err)
	}
	for _, schema := range state.StageSchemas() {
		if err := writeSchema(dir, schema); err != nil {
			return err
		}
	}
	return nil
}

func writeSchema(dir string, schema state.StageSchema) error {
	raw, err := schema.Generate()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, schema.FileName)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Println("wrote", path)
	return nil
}
