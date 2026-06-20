# DeadlinerServer Foundation Plan

## Document Type

SDD-style foundation plan for the open-source Deadliner backend.

## Architectural Style

Use DDD layered architecture as the main code organization strategy.

Recommended layers:

- `app`
- `domain`
- `infra`
- `service`

## Context

Deadliner already has a meaningful local domain model on iOS:

- DDL items use a stable `uid`
- sync deletion is represented by tombstones rather than by state alone
- habit sync is anchored to the carrier DDL uid
- local conflict handling uses logical version tuples: `ts`, `ctr`, `dev`

That model is good enough to guide backend design, but the backend should not
be a thin file mirror of the mobile clients. The new server should become the
canonical source of truth and support:

- self-hosted deployment
- open-source reuse
- multi-user login
- strict account isolation
- backend-first sync for iOS, Android, and HarmonyOS

This mode is intended to coexist with the original WebDAV sync mode as a
separate centralized product shape, not as an automatic replacement.

## Product Goals

1. Provide a clean, open-source backend that one person can self-host.
2. Preserve canonical Deadliner business rules across all clients.
3. Replace snapshot overwrite behavior with mutation + pull cursor sync.
4. Keep the protocol simple enough for three mobile clients to adopt.
5. Make the server production-ready for multiple independent user accounts.

## Non-Goals

The first phase does not attempt to provide:

- collaborative shared lists between users
- per-record CRDT merge for habit records
- web admin UI
- offline-first server-side automation
- public plugin ecosystem

## Canonical Domain Invariants

The backend must preserve these rules from the existing app model:

1. DDL business state and sync deletion are separate concerns.
2. DDL state values remain:
   - `active`
   - `completed`
   - `archived`
   - `abandoned`
   - `abandonedArchived`
3. Habit sync identity is the carrier DDL uid, not a standalone habit uid.
4. Subtasks remain embedded inside the DDL document payload.
5. Invalid enum values or structurally invalid documents fail explicitly.

## Backend-First Adaptation

The mobile clients currently reason about logical versions as if each device can
author canonical writes. For the server architecture, we adapt that model:

- client logical versions are accepted as mutation metadata
- server commit order becomes the canonical conflict boundary
- the server emits monotonic change ids as pull cursors
- the server stores mutation provenance for debugging and idempotency

This gives us a simple transition path:

1. keep `uid` stable across platforms
2. keep tombstones independent
3. keep habit carrier semantics
4. move canonical ordering authority to the server

## Transport Choice

Use Kitex with thrift as the primary internal and mobile-facing RPC layer.

Reasons:

- good fit for strongly typed Go services
- explicit contracts for all three clients
- easy future split into auth and sync services without rewriting transport
- consistent service scaffolding for open-source maintenance

## Service Boundary

Phase 1 starts with one Kitex service:

- `DeadlinerService`

Initial RPC groups:

- account registration
- login and session refresh
- mutation push
- incremental pull

This is intentionally compact. We can split into `AuthService` and
`SyncService` later when the API surface grows.

## Storage Choice

Default storage approach: multi-database support via GORM, with MySQL as the
primary first implementation.

Reasons:

- easy self-hosting
- mature operational tooling
- strong transactional semantics with InnoDB
- good indexing for multi-tenant change feeds
- GORM gives us pragmatic persistence mapping without forcing domain logic into
  SQL code too early
- additional dialects such as PostgreSQL or SQLite can be added behind the same
  infra abstractions later

## DDD Layering

### App Layer

- commands
- DTOs
- request-facing models consumed by services

### Service Layer

- use case orchestration
- transaction boundaries
- idempotency workflow

### Domain Layer

- aggregates
- value objects
- domain services
- repository interfaces

### Infra Layer

- GORM repositories
- driver-specific database openers
- MySQL models and migrations
- token/hash/time providers

## Proposed Persistence Model

### Account Layer

- `accounts`
- `devices`
- `sessions`

### Sync Layer

- `deadline_items`
- `habit_docs`
- `sync_changes`

### Key Decisions

- every row is scoped by `account_id`
- every mutation is recorded in `sync_changes`
- every mutation also needs a durable receipt for idempotent replay
- `sync_changes.change_id` is the canonical pull cursor
- deadline and habit documents are stored separately
- habit records are replaced as a whole document on newer writes
- GORM models are persistence adapters, not the canonical domain entities
- the infra layer should allow more than one SQL backend over time

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

Phase 1 policy is intentionally simple:

- server commit order is canonical
- duplicate mutation ids replay the original receipt instead of reapplying
- stale writes are detected by entity base change id
- client logical versions are retained for migration and diagnostics

This is the most practical bridge from the current snapshot world to a real
multi-user backend.

## Mobile Client Migration Strategy

### iOS

- map current `Ver(ts, ctr, dev)` to client mutation metadata
- keep current DDL and habit payload shapes as close as possible
- replace WebDAV snapshot reads/writes with push/pull calls

### Android

- align local model with canonical `state`
- embed subtasks in the DDL payload
- use the same carrier DDL habit model as iOS

### HarmonyOS

- follow the same thrift contract as Android and iOS
- keep compatibility logic at the client boundary only

## Repository Conventions

### `internal/domain`

Canonical business rules, models, services, and repository contracts live here first.

### `internal/app`

Application-boundary composition and mappers live here.

### `internal/infra`

GORM and MySQL implementations live here.

### `internal/utils`

Small shared helpers that do not belong to domain logic live here.

### `idl`

The thrift contract is the shared source of truth for client integration.

### `db/migrations`

Driver-specific SQL schema evolves explicitly through versioned SQL files.

## Milestones

### M1. Foundation

- plan document
- thrift contract
- DDD skeleton
- multi-database skeleton with MySQL schema

### M2. Auth and Session

- account registration
- password login
- refresh token rotation
- device registration

### M3. Sync Write Path

- push deadline mutations
- push habit mutations
- idempotent mutation handling
- change feed append

### M4. Sync Read Path

- pull by cursor
- pagination
- deadline and habit change serialization

### M5. Client Adoption

- iOS migration
- Android migration
- HarmonyOS migration

## Acceptance Criteria for This Foundation

This foundation is complete when:

- the repository clearly explains the target architecture
- the thrift contract reflects backend-first sync
- the default MySQL schema supports multi-user isolation
- Go packages exist for continued implementation without reorganizing later

## Open Questions

1. Should account identity be email-only or support username login too?
2. Should session auth be token-based only, or also expose a reverse proxy
   friendly cookie mode later?
3. Do we want a public REST gateway after the Kitex RPC layer stabilizes?
4. Should the canonical server version remain only `change_id`, or should we
   expose both `change_id` and RFC3339 `committed_at` everywhere?
