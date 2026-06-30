---
name: generate-fuzz-tests
description: Generate property-based fuzz tests to find edge cases human QA misses.
triggers:
  keywords: ["fuzz", "property test", "fast-check", "hypothesis"]
  intentPatterns: ["Generate fuzz tests for *", "Property test *", "Write a fast-check for *"]
standalone: true
---

## When To Use
Triggered to test pure functions, intricate algorithms, or complex parsers. Moves beyond example-based `test('1 + 1 = 2')` to property-based `test('a + b = b + a')`.

## Context To Load First
The target code file to test.

## Process
1. Analyze the target function's signature and invariants.
2. Determine appropriate property-based framework (e.g., `fast-check` for JS/TS, `hypothesis` for Python).
3. Identify 2-3 properties that must hold true for all inputs (e.g., "Sorting an array returns an array of the same length").
4. Formulate the property-based test code.
5. Create or append the test file.

## Output Format
```typescript
import fc from 'fast-check';
import { TargetFunction } from './file';

describe('TargetFunction Properties', () => {
  it('should maintain invariant [X] for all valid inputs', () => {
    fc.assert(
      fc.property(
        fc.string(), fc.integer(), (str, num) => {
          // Invariant assertion here
          const result = TargetFunction(str, num);
          expect(result).toHaveProperty('valid', true);
        }
      )
    );
  });
});
```

## Guardrails
- **Only** use this on pure functions, complex math, parsers, or decoders. Never use it on UI components or database queries without heavy mocking.
- Property testing frameworks must already be installed or explicitly added by the agent.
