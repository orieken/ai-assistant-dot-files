# Epic 49 — Mobile Stack Conventions (Swift + Kotlin)

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 3.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

`shared/rules/` covers TypeScript, Go, Python, Java, C# — no mobile-native coverage. Framework has no story for Swift/iOS or Kotlin/Android app development.

## Scope

**Two commits, one per language.** Create two convention files matching the shape of `shared/rules/typescript-conventions.md` (read that first for structure — sections: Complexity, File Naming, Testing, Frameworks, Quick Reference).

### Op 1 — `shared/rules/swift-conventions.md`

- **Complexity**: match framework's `< 7` cyclomatic budget via SwiftLint's `cyclomatic_complexity` rule
- **File naming**: PascalCase (Apple convention). Test files suffix `Tests.swift`. Model files `Model.swift`, ViewModel `ViewModel.swift`, etc.
- **Testing**: XCTest as primary, Nimble for expressive matchers (optional), Quick for BDD-style (optional), snapshot testing via `swift-snapshot-testing`. Mocking via protocol-based test doubles (Swift doesn't have runtime reflection like Java/Python; hand-rolled fakes or `Mockingbird` for generation)
- **Frameworks**: SwiftUI + Combine + Swift Concurrency (async/await) as modern default; UIKit legacy support only
- **Dependency injection**: constructor injection (native language pattern); `@Environment` for SwiftUI-scoped deps
- **Architecture**: MV-* patterns (MVVM most common); Clean Architecture applies (Domain / UseCase / Repository / Adapter). Same layer separation as other languages

Commit: `docs(rules): add swift-conventions.md (Epic 49)`.

### Op 2 — `shared/rules/kotlin-conventions.md`

- **Complexity**: match `< 7` via detekt's `ComplexMethod`
- **File naming**: PascalCase for classes, camelCase for functions/props, one public top-level declaration per file
- **Testing**: JUnit 5 as primary, MockK for Kotlin-idiomatic mocking, Espresso for Android UI, kotlinx-coroutines-test for coroutine testing
- **Frameworks**: Jetpack Compose + Coroutines + Flow as modern default; XML layouts legacy only
- **Dependency injection**: Hilt (Android) or Koin (Kotlin Multiplatform); constructor injection always
- **Architecture**: MVI or MVVM; Clean Architecture layer separation (Domain / UseCase / Repository / DataSource)

Commit: `docs(rules): add kotlin-conventions.md (Epic 49)`.

Both should be cross-referenceable from existing agents — no wiring changes needed today; mobile-target features will pick them up when a mobile-flavored feature spec is analyzed.

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- If SwiftLint / detekt configuration syntax has changed significantly since the last common reference — halt, note the version assumed.
- If the framework's `< 7` complexity budget doesn't map cleanly onto one of these tools' out-of-box config — halt, propose closest equivalent.

## Report (under 150 words)

```
Commits:
  <sha> swift-conventions
  <sha> kotlin-conventions
Sections covered per doc: <list>
Cross-references added: <list, or "none — mobile features will pick up when authored">
Health-check green: yes
```

Go.
