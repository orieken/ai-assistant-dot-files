package platform

func installCline(environment Environment) ([]string, error) {
	paths, err := installClineCoreRules(environment)
	if err != nil {
		return nil, err
	}
	for _, rule := range languageRules {
		if !isStackEnabled(environment.Rules, rule.stack) {
			continue
		}
		body, readErr := readRule(environment, rule.rule)
		if readErr != nil {
			return nil, readErr
		}
		name := clineLanguageName(rule)
		path, writeErr := writeGenerated(environment, ".clinerules/"+name+".md", clineRule(compactGlobs(rule.globs), body))
		if writeErr != nil {
			return nil, writeErr
		}
		paths = appendPath(paths, path)
	}
	return paths, nil
}

func installClineCoreRules(environment Environment) ([]string, error) {
	definitions := []generatedRule{
		{"", "00-approval-gates", "", "", "approval-gates.md"},
		{"", "01-design-principles", "", "", "design-principles.md"},
		{"", "02-architecture-guardrails", "", "", "architecture-guardrails.md"},
		{"", "04-testing-conventions", "", "**/*.spec.*,**/*.test.*,**/*.feature", "testing-conventions.md"},
	}
	var paths []string
	for _, rule := range definitions {
		body, err := readRule(environment, rule.rule)
		if err != nil {
			return nil, err
		}
		if rule.globs != "" {
			body = clineRule(rule.globs, body)
		}
		path, err := writeGenerated(environment, ".clinerules/"+rule.filename+".md", body)
		if err != nil {
			return nil, err
		}
		paths = appendPath(paths, path)
	}
	roster, err := agentRoster(environment.Content)
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".clinerules/03-agent-roster.md", []byte(roster))
	return appendPath(paths, path), err
}

func clineLanguageName(rule generatedRule) string {
	return languageIndex[rule.stack] + "-" + rule.filename
}
