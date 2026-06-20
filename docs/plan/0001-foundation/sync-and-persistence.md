# Sync And Persistence

## Storage Choice

Default storage approach: multi-database support via GORM, with MySQL as the
first production implementation.

Reasons:

- easy self-hosting
- mature operational tooling
- strong transactional semantics with InnoDB
- good indexing for multi-tenant change feeds
- room to add PostgreSQL or SQLite later behind the same abstractions

## Proposed Persistence Model

### Account Scope

- `accounts`
- `devices`
- `sessions`

### Sync Scope

- `deadline_items`
- `habit_docs`
- `sync_changes`
- `mutation_receipts`

## Persistence Rules

- every row is scoped by `account_id`
- every mutation is recorded in `sync_changes`
- every mutation must also have a durable receipt for replay safety
- `sync_changes.change_id` is the canonical pull cursor
- deadline and habit documents are stored separately
- habit records are replaced as a whole document on newer writes
- GORM models are adapters, not the canonical domain entities

## Sync Model

### Push

Clients send:

- authenticated device identity
- idempotent mutation ids
- per-entity base server change ids
- optional client logical version metadata
- full document payloads or tombstones

Server behavior:

1. authenticate session and device
2. validate enums and document shape
3. apply ownership and tenant checks
4. commit the mutation transactionally
5. append a `sync_changes` event
6. return the new canonical server version

### Pull

Clients pull by cursor:

- `cursor = last_applied_change_id`

Server returns:

- ordered deadline changes
- ordered habit changes
- next cursor
- pagination flag

### Conflict Policy

Phase 1 stays intentionally simple:

- server commit order is canonical
- duplicate mutation ids replay the original receipt
- stale writes are detected by entity base change id
- client logical versions are retained for migration and diagnostics
