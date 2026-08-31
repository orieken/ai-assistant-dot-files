# `internal/state`

Typed pipeline state (ADR-006, roadmap L2.9): Go structs for the stages that exchange data
instead of markdown, the projections that give each stage only the fields it reads, and the
renderers that turn state back into a human-readable markdown view.

JSON Schema in `shared/schemas/pipeline/` is **generated** from these structs by
`go run ./cmd/gen-schemas`; a test fails if the committed copies drift. Never hand-edit them.

First cut types the `analyst → architect` hop; every other stage still passes markdown.
