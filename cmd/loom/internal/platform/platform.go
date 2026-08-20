// Package platform detects AI platforms and installs their framework files.
package platform

import (
	"fmt"
	iofs "io/fs"
	"sort"
)

// FileWriter is the filesystem boundary used by platform installers.
type FileWriter interface {
	Install(source, destination string) (bool, error)
	Copy(source, destination string) (bool, error)
	CopyIfMissing(source, destination string) (bool, error)
	Write(destination string, content []byte) (bool, error)
	PrepareDirectory(destination string) error
}

// Result reports the paths owned by one platform installation.
type Result struct {
	Name  string
	Paths []string
}

// Environment contains the dependencies shared by platform installers.
type Environment struct {
	Content Content
	Files   FileWriter
	Rules   RuleSet
}

// Content is the embedded read-only framework source.
type Content interface {
	iofs.ReadDirFS
	iofs.ReadFileFS
}

type installer func(Environment) ([]string, error)

var installers = map[string]installer{
	"claude-code":    installClaude,
	"cursor":         installCursor,
	"windsurf":       installWindsurf,
	"github-copilot": installCopilot,
	"gemini":         installGemini,
	"openai-codex":   installCodex,
	"jetbrains":      installJetBrains,
	"roo-code":       installRooCode,
	"cline":          installCline,
}

// Names returns every supported platform name in stable order.
func Names() []string {
	names := make([]string, 0, len(installers))
	for name := range installers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// IsKnown reports whether name identifies a supported platform.
func IsKnown(name string) bool { return installers[name] != nil }

// Install installs one platform using the supplied environment.
func Install(name string, environment Environment) (Result, error) {
	install := installers[name]
	if install == nil {
		return Result{}, fmt.Errorf("unknown platform %q", name)
	}
	paths, err := install(environment)
	return Result{name, paths}, err
}
