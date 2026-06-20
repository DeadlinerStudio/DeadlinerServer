# Goals And Scope

## Context

Deadliner already has a meaningful local model on iOS:

- DDL items use a stable `uid`
- deletion is represented by tombstones rather than by state alone
- habit sync is anchored to the carrier DDL uid
- local conflict handling uses logical versions `Ver(ts, ctr, dev)`

The backend should preserve the durable domain ideas from that model without
remaining a file-mirroring extension of the mobile clients.

## Product Goals

1. Provide a clean, open-source backend that one person can self-host.
2. Preserve canonical Deadliner business rules across all clients.
3. Replace snapshot overwrite behavior with mutation plus pull cursor sync.
4. Keep the protocol simple enough for iOS, Android, and HarmonyOS to adopt.
5. Make the server production-ready for multiple independent user accounts.

## Non-Goals

The first phase does not attempt to provide:

- collaborative shared lists between users
- per-record CRDT merge for habit records
- web admin UI
- offline-first server-side automation
- public plugin ecosystem

## Canonical Domain Invariants

The backend must preserve these rules:

1. DDL business state and sync deletion are separate concerns.
2. DDL state values remain `active`, `completed`, `archived`, `abandoned`, and
   `abandonedArchived`.
3. Habit sync identity is the carrier DDL uid, not a standalone habit uid.
4. Subtasks remain embedded inside the DDL document payload.
5. Invalid enum values or structurally invalid documents fail explicitly.
