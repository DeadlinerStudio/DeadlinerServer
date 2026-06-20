# DeadlinerServer Sync Product Shape

## Purpose

Define the product positioning of DeadlinerServer as a centralized sync mode
that coexists with the original WebDAV mode while preserving offline-friendly
behavior.

## Product Positioning

Deadliner should eventually support multiple sync modes, not a single forced
backend model.

Recommended product shape:

1. `Local Only`
2. `WebDAV Sync`
3. `Deadliner Account Sync`

The new Go backend is the third mode.

It is not a replacement for local-first behavior.
It is not a denial of the value of WebDAV.
It is a different product tier with a different tradeoff:

- WebDAV favors user-controlled distributed snapshot sync
- Deadliner Account Sync favors a central source of truth with simpler
  multi-device convergence

## User-Facing Narrative

The clean user story should be:

- if you want simple, reliable, account-based sync, use Deadliner Account Sync
- if you want bring-your-own storage and no dedicated backend account, use
  WebDAV
- if you do not want any sync, stay on Local Only

This keeps the product philosophy honest.

## Why Keep WebDAV

WebDAV still has real value:

- user-controlled storage
- no permanent backend dependency
- easy personal backup mental model
- good fit for technically confident users

So the new backend should not be described as "the correct sync" and WebDAV as
"legacy".

Better framing:

- WebDAV is the distributed sync mode
- Deadliner Account Sync is the centralized sync mode

## Why Add A Centralized Mode

The centralized mode solves different problems:

- easier onboarding
- account-based login
- cleaner multi-user hosting
- simpler cross-device convergence
- future support for notifications, background jobs, and richer server features

Most users do not want to reason about remote files, ETags, or snapshot merge.
They just want to sign in and trust sync.

## Core Product Promise

Deadliner Account Sync should feel like:

- sign in once
- changes appear quickly on other devices
- offline usage still works
- sync recovers automatically after reconnect
- conflicts are rare and understandable

That is the emotional target, not just the technical target.

## Architecture Philosophy

The server is centralized in authority, but clients remain local-first in
interaction.

That means:

- the server is the canonical shared source of truth
- each client still owns a full local database
- writes are first committed locally
- sync is asynchronous when network is weak or missing
- reconnection is a catch-up problem, not a "reload everything" problem

This is the key design principle:

`centralized authority + local-first execution`

## Practical Sync Semantics

### Local Write First

When a user edits a task or habit:

1. write to local database immediately
2. enqueue a sync mutation locally
3. update UI from local state
4. push to server when possible

The UI must not depend on round-trip success for ordinary edits.

### Server As Canonical Authority

When connectivity exists:

1. client pushes pending mutations
2. server validates and commits them
3. server emits canonical ordered change ids
4. clients pull and reconcile

This gives us central convergence without sacrificing responsiveness.

### Weak Network Strategy

Under weak network:

- retries must be idempotent
- mutation queue must survive app restarts
- local reads and writes must remain available
- pull should resume from the last cursor, not restart from scratch

### No Network Strategy

Under no network:

- app behaves as local-only temporarily
- local mutation queue keeps growing safely
- reconnect triggers background push then pull

This is not a special mode toggle.
It is normal behavior.

## Conflict Philosophy

The new backend should reduce visible conflicts, not pretend conflicts never
exist.

Recommended phase-1 strategy:

- local write succeeds immediately
- server commit order is canonical
- mutation ids guarantee retry safety
- if another device changed the same object first, later pushes are still
  processed deterministically

User-facing principle:

- most conflicts should silently converge
- only rare, high-confusion cases deserve explicit UX later

## Product Boundary Versus WebDAV

The distinction should be explicit in both code and UX.

### WebDAV Mode

- remote file snapshots
- client-led merge
- distributed synchronization flavor

### Deadliner Account Sync Mode

- authenticated account session
- mutation push and cursor pull
- server-led canonical ordering

These modes should not share the same transport assumptions even if they reuse
the same domain model.

## Recommended Settings IA

Recommended settings structure:

### Sync

- `Sync Mode`
  - `Local Only`
  - `WebDAV`
  - `Deadliner Account`

### Deadliner Account

- sign in / sign out
- device list
- last sync status
- force resync

### WebDAV

- server URL
- username
- password or token
- test connection

This keeps the two sync products conceptually separate.

## Migration Strategy

Do not force users from WebDAV to the new backend.

Recommended migration path:

1. user selects `Deadliner Account`
2. app offers import from current local data
3. first upload creates the canonical remote state
4. other devices sign in and pull from server

For WebDAV users specifically:

- migration should be explicit and one-time
- avoid dual-write to both systems in the first phase
- after switching, one sync mode should be active at a time

## Why Not Dual Sync At Once

Supporting WebDAV and centralized backend at the same time for one dataset would
create a much harder product:

- two remote authorities
- ambiguous conflict semantics
- much harder support burden

So the product rule should be:

- one account chooses one active sync mode at a time

Import and migration are good.
Permanent double-sync is not a phase-1 goal.

## Future Product Opportunities Enabled By Centralization

The centralized mode unlocks features that WebDAV does not naturally support:

- server-issued push notifications
- background reminder jobs
- device management
- server-side audit or recovery tools
- family or team sharing later

These should not be phase-1 requirements, but they justify the architecture.

## Recommended Phase-1 Success Criteria

Phase 1 succeeds when a user can:

1. sign up on device A
2. create and edit tasks offline or online
3. reconnect and sync safely
4. sign in on device B
5. pull the same data with correct ordering and tombstone behavior

If that works reliably, the product shape is correct.

## Summary

Deadliner Account Sync should be designed as:

- centralized in authority
- local-first in interaction
- offline-tolerant in operation
- clearly separate from WebDAV in product language

That gives us the best of both worlds:

- simpler mainstream sync UX
- preserved resilience under weak or absent network
- continued respect for users who prefer WebDAV

