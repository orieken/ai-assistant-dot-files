package platform

import "fmt"

func installClaude(environment Environment) ([]string, error) {
	paths, err := installSources(environment, []sourceDestination{
		{"shared/agents", ".claude/agents"},
		{"shared/skills", ".claude/skills"},
	})
	if err != nil {
		return nil, err
	}
	rulePaths, err := installClaudeRules(environment)
	if err != nil {
		return nil, err
	}
	paths = append(paths, rulePaths...)
	return installClaudeProjectFiles(environment, paths)
}

func installClaudeRules(environment Environment) ([]string, error) {
	return installRuleDirectory(environment, ".claude/rules")
}

func installRuleDirectory(environment Environment, destination string) ([]string, error) {
	names, err := environment.Rules.Names(environment.Content)
	if err != nil {
		return nil, err
	}
	if len(environment.Rules.stacks) == 0 {
		return installSources(environment, []sourceDestination{{"shared/rules", destination}})
	}
	if err := environment.Files.PrepareDirectory(destination); err != nil {
		return nil, fmt.Errorf("prepare filtered rules directory: %w", err)
	}
	for _, name := range names {
		_, installErr := environment.Files.Install("shared/rules/"+name, destination+"/"+name)
		if installErr != nil {
			return nil, fmt.Errorf("install filtered rule %s: %w", name, installErr)
		}
	}
	return []string{destination}, nil
}

func installClaudeProjectFiles(environment Environment, paths []string) ([]string, error) {
	basePaths, err := installSources(environment, []sourceDestination{
		{"shared/ARCHITECTURE_RULES.md", "ARCHITECTURE_RULES.md"},
		{"shared/DOMAIN_DICTIONARY.md", "DOMAIN_DICTIONARY.md"},
	})
	if err != nil {
		return nil, err
	}
	paths = append(paths, basePaths...)
	return installClaudeTemplates(environment, paths)
}

func installClaudeTemplates(environment Environment, paths []string) ([]string, error) {
	installed, err := environment.Files.CopyIfMissing("templates/claude-feature-team/CLAUDE.md", "CLAUDE.md")
	if err != nil {
		return nil, err
	}
	if installed {
		paths = append(paths, "CLAUDE.md")
	}
	installed, err = environment.Files.CopyIfMissing("templates/claude-feature-team/features", "features")
	if err != nil {
		return nil, err
	}
	if installed {
		paths = append(paths, "features")
	}
	return paths, nil
}
