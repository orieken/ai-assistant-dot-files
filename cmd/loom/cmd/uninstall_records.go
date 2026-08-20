package cmd

import (
	"fmt"

	"github.com/orieken/loom/cmd/loom/internal/manifest"
)

func selectOwnedRecords(records []manifest.PlatformRecord, selected string) ([]manifest.PlatformRecord, []manifest.PlatformRecord, error) {
	if selected == "" {
		return records, nil, nil
	}
	var matched, remaining []manifest.PlatformRecord
	for _, record := range records {
		if record.Name == selected {
			matched = append(matched, record)
		} else {
			remaining = append(remaining, record)
		}
	}
	if len(matched) == 0 {
		return nil, nil, fmt.Errorf("platform %q is not recorded in the loom manifest", selected)
	}
	return matched, remaining, nil
}
