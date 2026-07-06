# api-generator (EXPERIMENTAL)

A personal OpenAPI-to-client generator: reads a Swagger/OpenAPI URL, emits docs plus Go and TypeScript
clients. Not a supported part of the ai-assistant-dot-files framework, not covered by
`scripts/ci-check.sh` or `framework-ci.yml`, and has no tests yet.

## Usage

```bash
npm run ingest -- <swagger-url> <go-output-dir> <ts-output-dir>
```

Output directories may also be set via `API_GENERATOR_GO_DIR` / `API_GENERATOR_TS_DIR` env vars instead
of positional args.
