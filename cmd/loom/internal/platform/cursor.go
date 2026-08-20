package platform

func installCursor(environment Environment) ([]string, error) {
	paths, err := installSources(environment, []sourceDestination{
		{"shared/agents", ".cursor/agents"},
		{"shared/skills", ".cursor/skills"},
	})
	if err != nil {
		return nil, err
	}
	rulePaths, err := installCursorRules(environment)
	if err != nil {
		return nil, err
	}
	paths = append(paths, rulePaths...)
	aggregate, err := aggregateInstructions(environment, "# Context Engineering Framework — All Rules")
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".cursorrules", aggregate)
	return appendPath(paths, path), err
}

func installCursorRules(environment Environment) ([]string, error) {
	paths, err := installCursorCoreRules(environment)
	if err != nil {
		return nil, err
	}
	paths, err = installCursorVueRule(environment, paths)
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
		content := cursorRule(rule.description, rule.globs, false, body)
		path, writeErr := writeGenerated(environment, ".cursor/rules/"+rule.filename+".mdc", content)
		if writeErr != nil {
			return nil, writeErr
		}
		paths = appendPath(paths, path)
	}
	return paths, nil
}

func installCursorCoreRules(environment Environment) ([]string, error) {
	definitions := []generatedRule{
		{"", "architecture", "Architecture guardrails", `"**/*.go", "**/*.ts", "**/*.py", "**/*.java"`, "architecture-guardrails.md"},
		{"", "design-principles", "Design principles", `"**/*.go", "**/*.ts", "**/*.py", "**/*.java"`, "design-principles.md"},
		{"", "approval-gates", "Approval gates", "", "approval-gates.md"},
		{"", "testing", "Testing framework rules", `"**/*.spec.*", "**/*.test.*", "**/*.feature"`, "testing-conventions.md"},
	}
	var paths []string
	for _, rule := range definitions {
		body, err := cursorCoreBody(environment, rule)
		if err != nil {
			return nil, err
		}
		content := cursorRule(rule.description, rule.globs, rule.filename == "approval-gates", body)
		path, err := writeGenerated(environment, ".cursor/rules/"+rule.filename+".mdc", content)
		if err != nil {
			return nil, err
		}
		paths = appendPath(paths, path)
	}
	return installCursorRoster(environment, paths)
}

func cursorCoreBody(environment Environment, rule generatedRule) ([]byte, error) {
	body, err := readRule(environment, rule.rule)
	if err != nil || rule.filename != "architecture" {
		return body, err
	}
	architecture, err := environment.Content.ReadFile("shared/ARCHITECTURE_RULES.md")
	return append(append(body, '\n', '\n'), architecture...), err
}

func installCursorVueRule(environment Environment, paths []string) ([]string, error) {
	if !isStackEnabled(environment.Rules, "typescript") {
		return paths, nil
	}
	content := cursorRule("Vue 3 + Tailwind conventions", `"**/*.vue", "**/*.tsx", "**/*.jsx"`, false, []byte(vueFrontendRules))
	path, err := writeGenerated(environment, ".cursor/rules/vue-frontend.mdc", content)
	return appendPath(paths, path), err
}

func installCursorRoster(environment Environment, paths []string) ([]string, error) {
	roster, err := agentRoster(environment.Content)
	if err != nil {
		return nil, err
	}
	content := cursorRule("Persona roster", "", true, []byte(roster))
	path, err := writeGenerated(environment, ".cursor/rules/agent-roster.mdc", content)
	return appendPath(paths, path), err
}
