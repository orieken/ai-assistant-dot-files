package platform_test

import (
	"testing"

	loom "github.com/orieken/loom"
	"github.com/orieken/loom/cmd/loom/internal/platform"
)

type recordingWriter struct {
	paths []string
}

func (writer *recordingWriter) Install(_, destination string) (bool, error) {
	writer.paths = append(writer.paths, destination)
	return true, nil
}

func (writer *recordingWriter) Copy(_, destination string) (bool, error) {
	writer.paths = append(writer.paths, destination)
	return true, nil
}

func (writer *recordingWriter) CopyIfMissing(_, destination string) (bool, error) {
	writer.paths = append(writer.paths, destination)
	return true, nil
}

func (writer *recordingWriter) Write(destination string, _ []byte) (bool, error) {
	writer.paths = append(writer.paths, destination)
	return true, nil
}

func (writer *recordingWriter) PrepareDirectory(destination string) error {
	writer.paths = append(writer.paths, destination)
	return nil
}

func TestEveryPlatformProducesInstallPaths(t *testing.T) {
	rules, err := platform.ParseRuleSet("go")
	if err != nil {
		t.Fatalf("parse rule filter: %v", err)
	}
	for _, name := range platform.Names() {
		t.Run(name, func(t *testing.T) {
			writer := &recordingWriter{}
			environment := platform.Environment{Content: loom.FrameworkFS, Files: writer, Rules: rules}
			result, installErr := platform.Install(name, environment)
			if installErr != nil {
				t.Fatalf("install platform: %v", installErr)
			}
			if len(result.Paths) == 0 || len(writer.paths) == 0 {
				t.Fatalf("platform %s produced no paths", name)
			}
		})
	}
}
