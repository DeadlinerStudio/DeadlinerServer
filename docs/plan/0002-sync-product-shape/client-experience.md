# Client Experience

## Core Product Promise

Deadliner Account Sync should feel like:

- sign in once
- changes appear quickly on other devices
- offline usage still works
- sync recovers automatically after reconnect
- conflicts are rare and understandable

## Architecture Philosophy

The server is centralized in authority, but clients remain local-first in
interaction.

This means:

- the server is the canonical shared source of truth
- each client still owns a full local database
- writes are first committed locally
- sync is asynchronous when network is weak or missing
- reconnection is a catch-up problem, not a full reload problem

Design principle:

`centralized authority + local-first execution`

## Practical Sync Semantics

### Local Write First

1. write to local database immediately
2. enqueue a sync mutation locally
3. update UI from local state
4. push to server when possible

### Weak Or No Network

- retries must be idempotent
- mutation queue must survive app restarts
- local reads and writes must remain available
- pull should resume from the last cursor
- reconnect triggers background push then pull

## Conflict Philosophy

The new backend should reduce visible conflicts, not pretend they never exist.

Phase 1 principle:

- most conflicts should silently converge
- only rare, high-confusion cases deserve explicit UX later
