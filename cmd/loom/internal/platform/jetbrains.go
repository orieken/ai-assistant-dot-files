package platform

func installJetBrains(environment Environment) ([]string, error) {
	paths, err := installJetBrainsRules(environment)
	if err != nil {
		return nil, err
	}
	guidelines, err := aggregateInstructions(environment, "# Junie Guidelines (Saturday Framework)")
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".junie/guidelines.md", guidelines)
	return appendPath(paths, path), err
}

func installJetBrainsRules(environment Environment) ([]string, error) {
	core := []generatedRule{
		{"", "00-approval-gates", "", "", "approval-gates.md"},
		{"", "01-design-principles", "", "", "design-principles.md"},
		{"", "02-architecture-guardrails", "", "", "architecture-guardrails.md"},
		{"", "04-testing-conventions", "", "", "testing-conventions.md"},
	}
	paths, err := writeJetBrainsRules(environment, core, nil)
	if err != nil {
		return nil, err
	}
	roster, err := agentRoster(environment.Content)
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".aiassistant/rules/03-agent-roster.md", []byte(roster))
	if err != nil {
		return nil, err
	}
	return writeJetBrainsRules(environment, languageRules, appendPath(paths, path))
}

func writeJetBrainsRules(environment Environment, rules []generatedRule, paths []string) ([]string, error) {
	for _, rule := range rules {
		if rule.stack != "" && !isStackEnabled(environment.Rules, rule.stack) {
			continue
		}
		body, err := readRule(environment, rule.rule)
		if err != nil {
			return nil, err
		}
		name := rule.filename
		if rule.stack != "" {
			name = jetBrainsLanguageName(rule)
		}
		path, err := writeGenerated(environment, ".aiassistant/rules/"+name+".md", body)
		if err != nil {
			return nil, err
		}
		paths = appendPath(paths, path)
	}
	return paths, nil
}

func jetBrainsLanguageName(rule generatedRule) string {
	return languageIndex[rule.stack] + "-" + rule.filename
}
