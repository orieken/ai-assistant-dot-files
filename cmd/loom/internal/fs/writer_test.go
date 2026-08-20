package frameworkfs

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

var writerFixture = fstest.MapFS{
	"shared/rules/core.md":   {Data: []byte("core")},
	"shared/agents/dev.md":   {Data: []byte("developer")},
	"shared/skills/check.sh": {Data: []byte("#!/bin/sh")},
}

func TestWriterPreservesEmbeddedShellScriptsAsExecutable(t *testing.T) {
	target := t.TempDir()
	writer := testWriter(target, t.TempDir(), true, false)
	if _, err := writer.Install("shared/skills", ".claude/skills"); err != nil {
		t.Fatalf("copy embedded skills: %v", err)
	}
	info, err := os.Stat(filepath.Join(target, ".claude/skills/check.sh"))
	if err != nil {
		t.Fatalf("inspect copied shell script: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatalf("expected executable shell script: mode=%v", info.Mode())
	}
}

func TestWriterCopiesEmbeddedDirectory(t *testing.T) {
	target := t.TempDir()
	writer := testWriter(target, t.TempDir(), true, false)
	if _, err := writer.Install("shared/rules", ".claude/rules"); err != nil {
		t.Fatalf("copy embedded directory: %v", err)
	}
	assertFileContent(t, filepath.Join(target, ".claude/rules/core.md"), "core")
}

func TestWriterLinksStableCacheContent(t *testing.T) {
	target, cache := t.TempDir(), t.TempDir()
	writer := testWriter(target, cache, false, false)
	if _, err := writer.Install("shared/agents", ".claude/agents"); err != nil {
		t.Fatalf("link embedded directory: %v", err)
	}
	destination := filepath.Join(target, ".claude/agents")
	if info, err := os.Lstat(destination); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected %s to be a symlink: %v", destination, err)
	}
	assertFileContent(t, filepath.Join(destination, "dev.md"), "developer")
}

func TestWriterDryRunWritesNothing(t *testing.T) {
	target, cache := t.TempDir(), t.TempDir()
	writer := testWriter(target, cache, false, true)
	if _, err := writer.Install("shared/agents", ".claude/agents"); err != nil {
		t.Fatalf("dry-run install: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(target, ".claude")); !os.IsNotExist(err) {
		t.Fatalf("dry run created target content: %v", err)
	}
}

func TestPrepareDirectoryReplacesExistingSymlink(t *testing.T) {
	target, cache := t.TempDir(), t.TempDir()
	writer := testWriter(target, cache, false, false)
	if _, err := writer.Install("shared/rules", ".claude/rules"); err != nil {
		t.Fatalf("link full rules: %v", err)
	}
	if err := writer.PrepareDirectory(".claude/rules"); err != nil {
		t.Fatalf("prepare selective rules: %v", err)
	}
	info, err := os.Lstat(filepath.Join(target, ".claude/rules"))
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("selective rules path is not a real directory: info=%v err=%v", info, err)
	}
}

func testWriter(target, cache string, isCopy, isDryRun bool) *Writer {
	return NewWriter(writerFixture, target, cache, isCopy, isDryRun, func(string) {})
}

func assertFileContent(t *testing.T, path, expected string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(content) != expected {
		t.Fatalf("content = %q, want %q", content, expected)
	}
}
