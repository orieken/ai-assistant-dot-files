package platform

import (
	"fmt"
	"strings"
)

func installRooCode(environment Environment) ([]string, error) {
	modes, err := roomodes(environment.Content)
	if err != nil {
		return nil, err
	}
	path, err := writeGenerated(environment, ".roomodes", modes)
	if err != nil {
		return nil, err
	}
	paths := appendPath(nil, path)
	names, err := environment.Rules.Names(environment.Content)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		installed, installErr := environment.Files.Copy("shared/rules/"+name, ".roo/rules/"+name)
		if installErr != nil {
			return nil, installErr
		}
		if installed {
			paths = append(paths, ".roo/rules/"+name)
		}
	}
	return paths, nil
}

func roomodes(content Content) ([]byte, error) {
	agents, err := readAgents(content)
	if err != nil {
		return nil, err
	}
	var output strings.Builder
	output.WriteString("customModes:\n")
	for _, agent := range agents {
		writeRoomode(&output, agent)
	}
	return []byte(output.String()), nil
}

func writeRoomode(output *strings.Builder, agent agentDocument) {
	fmt.Fprintf(output, "  - slug: %s\n", agent.name)
	fmt.Fprintf(output, "    name: %s\n", displayName(agent.name))
	fmt.Fprintf(output, "    description: %q\n", agent.description)
	fmt.Fprintf(output, "    whenToUse: %q\n", agent.description)
	output.WriteString("    roleDefinition: |\n")
	for _, line := range strings.Split(agent.body, "\n") {
		fmt.Fprintf(output, "      %s\n", line)
	}
	output.WriteString("    groups:\n")
	for _, group := range inferGroups(agent.tools) {
		fmt.Fprintf(output, "      - %s\n", group)
	}
	output.WriteString("\n")
}
