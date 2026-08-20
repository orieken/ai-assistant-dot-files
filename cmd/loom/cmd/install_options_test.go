package cmd

import "testing"

func TestPrepareInstallRejectsUnknownPlatform(t *testing.T) {
	flags := installFlags{target: t.TempDir(), platform: "unknown"}
	if _, err := prepareInstall(flags); err == nil {
		t.Fatal("expected unknown platform error")
	}
}

func TestPrepareInstallRejectsUnknownStack(t *testing.T) {
	flags := installFlags{target: t.TempDir(), stack: "ruby"}
	if _, err := prepareInstall(flags); err == nil {
		t.Fatal("expected unknown stack error")
	}
}
