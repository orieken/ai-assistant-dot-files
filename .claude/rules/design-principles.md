# Design Principles

These cross-cutting design principles MUST be honored by every agent, regardless of role.

## 1. Simple Design Rules (Kent Beck)
Apply these strictly in priority order:
1. **Passes the tests:** The system must work first.
2. **Reveals intention:** Code must be easily readable.
   - *Example (TypeScript):* Rename `const d = 5;` to `const maxRetries = 5;`
3. **No duplication:** Extract shared logic to a single source of truth.
   - *Example (TypeScript):* Extract identical fetch configuration into a shared `apiClient()`.
4. **Fewest elements:** Don't add abstraction without a clear need. Remove unused components.

## 2. Refactoring Catalog (Martin Fowler)
Use these named operations to fix smells:
- **Extract Function:** Smell = Long function (>30 LOC) or needing comments.
- **Inline Function:** Smell = Unnecessary abstraction/indirection.
- **Extract Variable:** Smell = Complex, unreadable expression.
- **Rename Variable:** Smell = Cryptic or generic naming.
- **Replace Conditionals with Polymorphism:** Smell = Long switch statements checking types.
- **Introduce Parameter Object:** Smell = Data clumps (e.g. `startDate`, `endDate` -> `DateRange`).
- **Decompose Conditional:** Smell = Complex conditional hiding the business rule.
- **Replace Magic Number/String:** Smell = Unexplained literal value.
- **Separate Query from Modifier:** Smell = Function returning a value AND mutating state.
- **Move Function / Field:** Smell = Feature Envy (used more by another class).

## 3. Anti-Pattern Radar
Watch for and actively refactor these common structural issues:
- **Distributed Monolith:** Services coupled by database or chatter.
- **Anemic Domain Model:** Entities strictly holding data with business logic leaked into services.
- **God Object:** Classes or modules that know or do too much.
- **Shotgun Surgery:** Every time you make a change, you have to edit multiple files.
- **Leaky Abstraction:** An interface that exposes its internal implementation details.

## 4. Hard Structural Limits (Sandi Metz)
These constraints are NON-NEGOTIABLE during design and code reviews:
- **Class Size:** $\le$ 100 lines of code.
- **Method Size:** $\le$ 5 lines of code.
- **Method Parameters:** $\le$ 4 parameters (ideally $\le$ 3).

## 5. The Boy Scout Rule
Active Instruction: **Always leave the code cleaner than you found it.**
If you touch a file to add a feature, you MUST proactively fix minor technical debt.
- *Example:* Found a magic number? Extract it. Found a poorly named variable? Rename it. Found loose types (`any`)? Tighten them before delivering your task.

## 6. Naming Standards
- **Intention-revealing:** Variables explain *why* they exist.
- **No abbreviations:** Write `UserRepository` carefully, not `UserRepo`.
- **No generic terms:** Do not use `data`, `info`, `manager`, `processor` without explicit context (e.g., `BillingManager` is okay if it actually manages billing, `DataManager` is not).
- **Boolean Prefixes:** all boolean variables and functions returning booleans MUST start with `is`, `has`, `can`, or `should` (e.g. `isActive`, `hasPermission`).
