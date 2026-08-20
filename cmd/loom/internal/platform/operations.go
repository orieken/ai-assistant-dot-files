package platform

import "fmt"

type sourceDestination struct {
	source      string
	destination string
}

func installSources(environment Environment, pairs []sourceDestination) ([]string, error) {
	var paths []string
	for _, pair := range pairs {
		installed, err := environment.Files.Install(pair.source, pair.destination)
		if err != nil {
			return nil, fmt.Errorf("install %s: %w", pair.destination, err)
		}
		if installed {
			paths = append(paths, pair.destination)
		}
	}
	return paths, nil
}

func writeGenerated(environment Environment, destination string, content []byte) (string, error) {
	installed, err := environment.Files.Write(destination, content)
	if err != nil {
		return "", fmt.Errorf("write %s: %w", destination, err)
	}
	if !installed {
		return "", nil
	}
	return destination, nil
}

func appendPath(paths []string, path string) []string {
	if path == "" {
		return paths
	}
	return append(paths, path)
}
