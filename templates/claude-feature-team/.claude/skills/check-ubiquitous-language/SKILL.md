---
name: check-ubiquitous-language
description: Analyzes source code for violations of the Ubiquitous Language defined in DOMAIN_DICTIONARY.md. Fails if invalid synonyms are used.
triggers:
  keywords: ["ubiquitous language", "domain dictionary", "terminology", "synonyms", "ddd", "naming"]
  intentPatterns: ["check class names against domain dictionary", "is this using the correct ubiquitous language"]
standalone: true   # must work without MCP/external systems
---

## When To Use
Use this skill when reviewing code to ensure naming conventions exactly match the business domain concepts defined in the `DOMAIN_DICTIONARY.md`. 
Apply it to class names, method names, and variable names.
Do NOT use when checking external libraries, third-party integrations, or configurations.

## Context To Load First
1. `DOMAIN_DICTIONARY.md`
2. `ARCHITECTURE_RULES.md`

## Process
1. Determine the files or directories you want to check.
   - What to do: Find exactly the domain entities and files you are reviewing.
   - What to produce: A list of target files.
   - What to check: Verify targets exist.
2. Run the dictionary checker.
   - What to do: Run `./.claude/skills/check-ubiquitous-language/check.sh [target]`
   - What to produce: The console output of the script identifying banned synonyms.
   - What to check: Analyze if technical synonyms strictly forbidden in the dictionary are used.
3. Handle Name Violations.
   - What to do: If violations are found, instruct the developer to rename the offending terms to match the Canonical Term in the domain dictionary.
   - When to pause: If there is ambiguity about a new business concept not in the dictionary, pause for human approval.

## Output Format
```markdown
# Ubiquitous Language Report

## Violations
- `path/to/user.ts`: Class `ClientAccount` uses synonym `Client`. (Canonical Term: `User`)

## Recommended Fixes
- Rename `ClientAccount` to `UserAccount` to align with the Domain Dictionary.
```

## Guardrails
- You MUST NOT approve code that uses forbidden synonyms defined in the domain dictionary.
- Do NOT auto-rename files or classes without verifying the Canonical Term.

## Standalone Mode
If `check.sh` is unavailable, read `DOMAIN_DICTIONARY.md` using the `Read` tool, map out forbidden words, and manually `grep` the source code files for those words to enforce the dictionary constraints manually.
