// Framework fitness-function floor — ESLint (flat config, ESLint v9+)
// Source convention: shared/rules/typescript-conventions.md
// CAP: 6  (cyclomatic complexity — enforces the framework-wide "< 7" rule)
//
// Drop this file into a project as eslint.config.js (or import/spread it into
// an existing eslint.config.js). Only the complexity cap is set here — this is
// a floor, not a full lint policy.

export default [
  {
    rules: {
      // "error" + max 6: flag functions whose cyclomatic complexity exceeds 6.
      // ESLint complexity rule reports when complexity > N, so max:6 catches 7+.
      complexity: ["error", 6],
    },
  },
];
