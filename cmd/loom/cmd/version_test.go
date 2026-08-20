package cmd

import (
	"bytes"
	"testing"
	"testing/fstest"
)

func TestWriteVersionIncludesBuildAndFrameworkMetadata(t *testing.T) {
	content := fstest.MapFS{"shared/VERSION": {Data: []byte("v3.3.14\n")}}
	var output bytes.Buffer
	if err := writeVersion(&output, content, "v0.1.0", "abc1234", "2026-08-19"); err != nil {
		t.Fatalf("write version: %v", err)
	}
	want := "loom v0.1.0 (commit abc1234, built 2026-08-19)\nframework v3.3.14 (embedded)\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestReadFrameworkVersionFallsBackWhenFileIsMissing(t *testing.T) {
	version, err := readFrameworkVersion(fstest.MapFS{})
	if err != nil || version != "embedded" {
		t.Fatalf("version = %q, err = %v", version, err)
	}
}
