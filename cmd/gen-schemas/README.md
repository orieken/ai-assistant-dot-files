# `cmd/gen-schemas`

Writes the JSON Schema for every typed pipeline stage from the Go structs in `internal/state`
into `shared/schemas/pipeline/`.

```bash
go run ./cmd/gen-schemas
```

Run it after changing any state struct. `TestCommittedSchemasMatchTheStructs` fails when the
committed schemas drift from their source, so CI catches a forgotten regeneration.
