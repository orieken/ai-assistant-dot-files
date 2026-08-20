package platform

func installGemini(environment Environment) ([]string, error) {
	paths, err := installSources(environment, []sourceDestination{{"shared/skills", ".agents/skills"}})
	if err != nil {
		return nil, err
	}
	rulePaths, err := installRuleDirectory(environment, ".agents/rules")
	if err != nil {
		return nil, err
	}
	paths = append(paths, rulePaths...)
	content, err := aggregateInstructions(environment, "# AGENTS.md")
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, "AGENTS.md", content)
	return appendPath(paths, path), err
}
