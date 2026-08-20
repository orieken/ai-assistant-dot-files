package cmd

import "testing"

func TestPrepareInstallRejectsUnknownPlatform(t *testing.T) {
	flags := installFlags{target: t.TempDir(), platform: "unknown"}
	if _, err := prepareInstall(flags, "v3.3.14"); err == nil {
		t.Fatal("expected unknown platform error")
	}
}

func TestPrepareInstallRejectsUnknownStack(t *testing.T) {
	flags := installFlags{target: t.TempDir(), stack: "ruby"}
	if _, err := prepareInstall(flags, "v3.3.14"); err == nil {
		t.Fatal("expected unknown stack error")
	}
}
