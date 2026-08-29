package platform

import (
	"fmt"
	"strings"
)

type generatedRule struct {
	stack       string
	filename    string
	description string
	globs       string
	rule        string
}

var languageRules = []generatedRule{
	{"go", "go-conventions", "Go backend conventions", `"**/*.go", "**/go.mod", "**/go.sum"`, "go-conventions.md"},
	{"typescript", "typescript-conventions", "TypeScript conventions", `"**/*.ts"`, "typescript-conventions.md"},
	{"python", "python-conventions", "Python conventions", `"**/*.py"`, "python-conventions.md"},
	{"csharp", "csharp-conventions", "C# conventions", `"**/*.cs"`, "csharp-conventions.md"},
	{"java", "java-conventions", "Java conventions", `"**/*.java"`, "java-conventions.md"},
	{"kotlin", "kotlin-conventions", "Kotlin conventions", `"**/*.kt", "**/*.kts"`, "kotlin-conventions.md"},
	{"swift", "swift-conventions", "Swift conventions", `"**/*.swift"`, "swift-conventions.md"},
	{"rust", "rust-conventions", "Rust conventions", `"**/*.rs", "**/Cargo.toml"`, "rust-conventions.md"},
}

const vueFrontendRules = `# Vue 3 + Tailwind Frontend Conventions

ALWAYS use Vue 3 Composition API with <script setup>.
NEVER use Options API in new components.
ALWAYS use Tailwind CSS utility classes unless custom CSS is required.
ALWAYS use TypeScript with strict mode and never use raw any types.
ALWAYS co-locate component tests alongside components.
CRITICAL: Components must remain below 100 lines.
`

func readRule(environment Environment, name string) ([]byte, error) {
	content, err := environment.Content.ReadFile("shared/rules/" + name)
	if err != nil {
		return nil, fmt.Errorf("read rule %s: %w", name, err)
	}
	return content, nil
}

func isStackEnabled(rules RuleSet, stack string) bool {
	if !rules.isFiltered() {
		return true
	}
	if len(rules.names) > 0 {
		return containsString(rules.names, stackRules[stack])
	}
	return containsString(rules.stacks, stack)
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

// languageIndex maps a stack name to its two-digit ordering prefix used by
// JetBrains and Cline rule filenames.
var languageIndex = map[string]string{
	"go": "05", "typescript": "06", "python": "07", "csharp": "08",
	"java": "09", "kotlin": "10", "swift": "11", "rust": "12",
}

func cursorRule(description, globs string, alwaysApply bool, body []byte) []byte {
	frontmatter := fmt.Sprintf("---\ndescription: %q\nalwaysApply: %t\n", description, alwaysApply)
	if globs != "" {
		frontmatter += "globs: [" + globs + "]\n"
	}
	return []byte(frontmatter + "---\n" + string(body))
}

func scopedRule(globs string, body []byte) []byte {
	return []byte("---\napplyTo: \"" + globs + "\"\n---\n" + string(body))
}

func clineRule(globs string, body []byte) []byte {
	var header strings.Builder
	header.WriteString("---\npaths:\n")
	for _, glob := range strings.Split(globs, ",") {
		fmt.Fprintf(&header, "  - %q\n", strings.TrimSpace(glob))
	}
	header.WriteString("---\n")
	return append([]byte(header.String()), body...)
}
