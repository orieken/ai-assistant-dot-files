# Pattern: Site-Centric Architecture (excerpt — under audit)

## Overview
The Saturday Framework eschews the Page Object Model in favor of a Site-Centric hierarchy:
`BaseSite` → `BasePage` → `BaseElement` → `BaseFlow`.

## Example

```typescript
// From src/saturday-core/src/base/BaseSite.ts (line 12)
export abstract class BaseSite {
  protected readonly siteManager: SiteManager;
  constructor(siteManager: SiteManager) {
    this.siteManager = siteManager;
  }
}
```

```typescript
// From src/saturday-core/src/base/BasePage.ts (line 8)
// This class was renamed to AbstractPage in v2.1 — snippet is stale
export abstract class BasePage {
  navigate(url: string): Promise<void> { ... }
}
```

## SiteManager Multi-Context Pattern
Use `SiteManager` to orchestrate cross-application journeys. Each application registers itself
with the manager at initialization.

See: `src/saturday-core/src/managers/SiteManager.ts`
// File was moved to src/saturday-core/src/orchestration/SiteManager.ts in v2.0
