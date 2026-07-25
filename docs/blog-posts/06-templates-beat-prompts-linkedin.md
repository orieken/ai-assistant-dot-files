In multi-agent AI architectures, inlining output markdown examples inside agent prompt files is a trap.

Over time, LLM drift causes headings to shift—`## Summary` becomes `## Overview`, `### Tasks` becomes `## Task List`—breaking automated contract validators and downstream parsing.

We fixed this in `ai-assistant-dot-files` by extracting **13 single-source-of-truth templates** to `shared/templates/`.

Instead of embedding 60 lines of example markdown inside every agent prompt:
1. Agent prompts define personality and reasoning rules.
2. Output structures live in `shared/templates/*.template.md`.
3. Contract validators enforce the exact template headings.

The result: zero heading-mismatch retries, cleaner prompts, and single-file maintainability when output structures change.

Full technical breakdown and the "loader path" gotcha we hit along the way: TODO_DEVTO_URL
