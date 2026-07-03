---
name: api-ingest
description: Ingests a Swagger/OpenAPI URL to generate feature docs and type-safe API clients for TS and Go.
triggers:
  keywords: ["ingest api", "generate clients", "swagger import", "openapi import"]
  intentPatterns: ["Ingest the API from *", "Generate Sunday clients for *", "Sync my API clients with *"]
standalone: true
---

## When To Use
When you want to bootstrap or update documentation and API clients based on an external Swagger/OpenAPI specification.

## Context To Load First
1. `DOMAIN_DICTIONARY.md`
2. `ARCHITECTURE_RULES.md`
3. `API_FRAMEWORK_BLUEPRINT_PROMPT.md`

## Process
1. **Fetch**: Obtain the OpenAPI spec (JSON or YAML) from the provided URL.
2. **Parse**: Analyze the spec to identify endpoints, request/response models, and security requirements.
3. **Generate Docs**: 
    - Create `docs/features/<endpoint-slug>/analysis.md`.
    - Create `docs/features/<endpoint-slug>/architecture-notes.md`.
4. **Generate TS Client**:
    - Create Zod schemas for models.
    - Create a client class extending `BaseApiClient`.
    - Placement: `packages/clients/generated/`.
5. **Generate Go Client**:
    - Create Go structs with `json` and `validate` tags.
    - Create client methods using `go-sunday` generics.
    - Placement: `/Users/oscarrieken/Projects/Rieken/go-sunday/client/generated/`.

## Output Format
A summary of generated files and a link to the new feature documentation.

## Guardrails
- Use `operationId` for naming if available; otherwise, use `[Method][Path]`.
- Ensure all generated code adheres to the 30-line function limit.
- Do not overwrite custom modifications in existing clients; use `.generated.ts/go` suffixes or separate folders.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
