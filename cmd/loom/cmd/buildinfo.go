package cmd

// Build metadata is injected by goreleaser via -ldflags.
// Values retain these defaults for local development builds.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)
