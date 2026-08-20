package platform

func installCodex(environment Environment) ([]string, error) {
	content, err := aggregateInstructions(environment, "# OpenAI / Codex Instructions (Saturday Framework)")
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".openai.md", content)
	return appendPath(nil, path), err
}
