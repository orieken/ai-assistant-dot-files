// Intentional test fixture for tests/platform-verification/ — NOT real framework code.
// Deliberately violates testing.mdc / testing.instructions.md rules:
//   - Tests implementation details (calls a private-ish method name, checks internal state)
//   - Vague test name ("test 1")
//   - No Given/When/Then structure, no fixtures/factories
//   - No mocking of the "database" dependency

import { describe, it, expect } from "vitest";

describe("CartService", () => {
  it("test 1", () => {
    const service = new CartService(new RealDatabaseConnection());
    service._internalRecalculate();
    expect(service._cachedTotal).toBe(42);
  });
});

class CartService {
  _cachedTotal = 0;
  constructor(private db: unknown) {}
  _internalRecalculate() {
    this._cachedTotal = 42;
  }
}

class RealDatabaseConnection {}
