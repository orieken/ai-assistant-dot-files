package platform

import (
	"fmt"
	"sort"
	"strings"
)

type agentDocument struct {
	name        string
	description string
	tools       string
	body        string
	status      string
	successor   string
}

func readAgents(content Content) ([]agentDocument, error) {
	entries, err := content.ReadDir("shared/agents")
	if err != nil {
		return nil, fmt.Errorf("read embedded agents: %w", err)
	}
	var agents []agentDocument
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "CHANGELOG.md" || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		agent, readErr := readAgent(content, entry.Name())
		if readErr != nil {
			return nil, readErr
		}
		agents = append(agents, agent)
	}
	sort.Slice(agents, func(left, right int) bool { return agents[left].name < agents[right].name })
	return agents, nil
}

func readAgent(content Content, name string) (agentDocument, error) {
	raw, err := content.ReadFile("shared/agents/" + name)
	if err != nil {
		return agentDocument{}, fmt.Errorf("read agent %s: %w", name, err)
	}
	parts := strings.SplitN(string(raw), "---\n", 3)
	if len(parts) != 3 {
		return agentDocument{}, fmt.Errorf("agent %s has invalid frontmatter", name)
	}
	return parseAgent(parts[1], strings.TrimSpace(parts[2])), nil
}

func parseAgent(frontmatter, body string) agentDocument {
	return agentDocument{
		name:        frontmatterValue(frontmatter, "name"),
		description: frontmatterValue(frontmatter, "description"),
		tools:       frontmatterValue(frontmatter, "tools"),
		body:        body,
		status:      frontmatterValue(frontmatter, "status"),
		successor:   frontmatterValue(frontmatter, "superseded_by"),
	}
}

func frontmatterValue(frontmatter, key string) string {
	prefix := key + ":"
	for _, line := range strings.Split(frontmatter, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}
