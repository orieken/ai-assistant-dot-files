# Analysis: Multi-Tenant Schema Migration

## Feature
Add a `tenant_id` column to the `orders` table to support multi-tenant data isolation.
All queries against `orders` must be filtered by `tenant_id` going forward.

## Proposed Migration
```sql
ALTER TABLE orders ADD COLUMN tenant_id UUID NOT NULL REFERENCES tenants(id);
```

## Rollout Plan
Deploy migration, then deploy application code that sets `tenant_id` on all new writes.
Existing rows will have `tenant_id = NULL` until a backfill job runs.

## Scale
`orders` table: 8 million rows. Migration estimated to lock the table for 45 seconds.
