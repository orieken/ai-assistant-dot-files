package platform

func installWindsurf(environment Environment) ([]string, error) {
	content, err := aggregateInstructions(environment, "# Context Engineering Framework — All Rules")
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".windsurfrules", content)
	return appendPath(nil, path), err
}
