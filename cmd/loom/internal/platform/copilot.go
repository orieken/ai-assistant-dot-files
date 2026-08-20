package platform

func installCopilot(environment Environment) ([]string, error) {
	content, err := aggregateInstructions(environment, "# Copilot Instructions (Saturday Framework)")
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".github/copilot-instructions.md", content)
	if err != nil {
		return nil, err
	}
	return installCopilotScopedRules(environment, appendPath(nil, path))
}

func installCopilotScopedRules(environment Environment, paths []string) ([]string, error) {
	paths, err := installCopilotVueRule(environment, paths)
	if err != nil {
		return nil, err
	}
	definitions := append([]generatedRule{
		{"", "testing", "", "**/*.spec.*,**/*.test.*,**/*.feature", "testing-conventions.md"},
	}, languageRules...)
	for _, rule := range definitions {
		if rule.stack != "" && !isStackEnabled(environment.Rules, rule.stack) {
			continue
		}
		body, err := readRule(environment, rule.rule)
		if err != nil {
			return nil, err
		}
		globs := rule.globs
		if rule.stack != "" {
			globs = compactGlobs(rule.globs)
		}
		path, err := writeGenerated(environment, ".github/instructions/"+rule.filename+".instructions.md", scopedRule(globs, body))
		if err != nil {
			return nil, err
		}
		paths = appendPath(paths, path)
	}
	return paths, nil
}

func installCopilotVueRule(environment Environment, paths []string) ([]string, error) {
	if !isStackEnabled(environment.Rules, "typescript") {
		return paths, nil
	}
	content := scopedRule("**/*.vue,**/*.tsx,**/*.jsx", []byte(vueFrontendRules))
	path, err := writeGenerated(environment, ".github/instructions/vue-frontend.instructions.md", content)
	return appendPath(paths, path), err
}

func compactGlobs(globs string) string {
	result := make([]byte, 0, len(globs))
	for _, char := range []byte(globs) {
		if char != '"' && char != ' ' {
			result = append(result, char)
		}
	}
	return string(result)
}
