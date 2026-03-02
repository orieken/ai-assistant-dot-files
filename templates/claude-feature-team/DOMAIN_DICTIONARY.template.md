# DOMAIN_DICTIONARY.md (Template)

This file acts as the single source of truth for the **Ubiquitous Language** in this codebase. All domains, entities, value objects, and terms MUST be explicitly defined here.

**Rule:** Analysts, Architects, Developers, and Code Reviewers must use *exactly* these terms. Do not invent synonyms (e.g., if this dictionary says "User", do not use "Client" or "AccountHolder" without adding it here first).

---

## 1. Core Domains (Bounded Contexts)
*List the high-level business areas in the system.*

- **[Domain Name]**: [Description or purpose of this bounded context]
  - *Example*: **Identity**: Handles user authentication, registration, and RBAC.

## 2. Entities (Identifiable Objects)
*Objects with a distinct identity that persists over time.*

| Term | Domain | Description | Synonyms to AVOID |
|---|---|---|---|
| **[Entity Name]** | `[Domain]` | [What this represents in the business] | `[Synonym]`, `[Synonym]` |
| *Example: User* | `Identity` | A person holding an account who can log in. | `Client`, `Customer`, `Player` |

## 3. Value Objects (Immutable Concepts)
*Concepts defined by their attributes, without a distinct identity.*

| Term | Domain | Description | Example Shape / Fields |
|---|---|---|---|
| **[Value Object]** | `[Domain]` | [What this represents] | `[Fields/Structure]` |
| *Example: Money* | `Billing` | Represents a financial amount and currency. | `{ amount: int, currency: string }` |

## 4. Domain Events
*Significant business occurrences that other parts of the system might care about.*

| Event Name | Domain | Triggered When |
|---|---|---|
| **[EventName]Completed** | `[Domain]` | [Condition] |
| *Example: UserRegistered* | `Identity` | A new User successfully verifies their email. |

## 5. Operations / Actions
*Standardized verbs used across the codebase.*

- **[Action]**: [What it means in this context]
  - *Example*: **Provision**: The act of allocating infrastructure for a tenant (Do not use "Create" or "Setup" for infrastructure).

---
*Note to Analysts: When designing a new feature, if a new business concept is introduced, you MUST add it to this dictionary before handing off to the developer.*
