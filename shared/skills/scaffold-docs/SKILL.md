---
name: scaffold-docs
description: Create a comprehensive implementation guide in markdown format
triggers:
  keywords: [scaffold docs, create guide, implementation guide]
  intentPatterns: ["create a comprehensive guide", "write implementation docs"]
standalone: true
---

## When To Use
When the user asks to create structured implementation plans, prep guides, or technical documents.

## Context To Load First
1. Check the existing `docs/` directory to avoid duplicate namings.

## Process
1. Understand the core topic or feature the user wants to implement.
2. Draft an outline covering:
   - Prerequisites and context
   - Step-by-step instructions with code examples
   - Verification steps
   - Common pitfalls and solutions
3. Save the formatted document to the `docs/` directory with a descriptive `.md` filename.

## Output Format
A detailed markdown file containing the required sections.

## Guardrails
- Do not overwrite existing documentation without human approval.

## Standalone Mode
Fully supported.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
