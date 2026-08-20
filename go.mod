module github.com/orieken/loom

go 1.26.5

require (
	github.com/orieken/ai-assistant-dotfiles/mcp v0.0.0
	github.com/spf13/cobra v1.10.1
)

ignore tests/platform-verification/fixtures

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

replace github.com/orieken/ai-assistant-dotfiles/mcp => ./shared/mcp
