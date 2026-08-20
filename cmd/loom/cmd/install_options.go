package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	frameworkfs "github.com/orieken/loom/cmd/loom/internal/fs"
	"github.com/orieken/loom/cmd/loom/internal/platform"
)

const frameworkVersion = "0.1.0"

type installRequest struct {
	target     string
	cache      string
	platforms  []string
	rules      platform.RuleSet
	isCopy     bool
	isDryRun   bool
	withConfig bool
	withMCP    bool
}

func prepareInstall(flags installFlags) (installRequest, error) {
	target, err := frameworkfs.ResolveTarget(flags.target)
	if err != nil {
		return installRequest{}, err
	}
	rules, err := platform.ParseRuleSet(flags.stack)
	if err != nil {
		return installRequest{}, err
	}
	platforms, err := selectPlatforms(target, flags.platform)
	if err != nil {
		return installRequest{}, err
	}
	cache, err := frameworkCache()
	return installRequest{target, cache, platforms, rules, flags.isCopy, flags.isDryRun, flags.withConfig, flags.withMCP}, err
}

func selectPlatforms(target, selected string) ([]string, error) {
	if selected == "" {
		return platform.Detect(target)
	}
	if !platform.IsKnown(selected) {
		return nil, fmt.Errorf("unknown platform %q (valid: %s)", selected, strings.Join(platform.Names(), ", "))
	}
	return []string{selected}, nil
}

func frameworkCache() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("resolve user cache: %w", err)
	}
	return filepath.Join(cache, "loom", frameworkVersion), nil
}
