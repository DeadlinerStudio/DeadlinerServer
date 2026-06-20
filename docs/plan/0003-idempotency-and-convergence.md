# DeadlinerServer Idempotency And Convergence

## Purpose

Define how Deadliner Account Sync should achieve stronger multi-device
idempotency and convergence than the original WebDAV snapshot mode.

## Why This Can Be Better Than WebDAV

WebDAV sync is fundamentally file-oriented:

- clients race on shared remote snapshots
- retries are file rewrites
- conflict handling is mostly snapshot merge logic

The centralized backend can do better because it has:

- authenticated device identity
- append-only server commit order
- mutation-level deduplication
- durable acknowledgement receipts

That lets us make retries safe in a much more explicit way.

## Core Rule

Every local edit becomes a durable mutation with a stable identity before it is
ever sent to the server.

That means the unit of sync is not "latest remote file".
The unit of sync is "this specific mutation".

## Client-Side Rules

### 1. Stable Device Identity

Each installed app instance has a stable `device_uid`.

This is required for:

- mutation identity
- deduplication
- debugging
- device management

### 2. Durable Local Mutation Queue

Every local edit is written in two places:

1. local business tables
2. local mutation queue

The queue must survive:

- app restarts
- process kills
- offline periods

### 3. Immutable Mutation Identity

Each mutation gets a unique immutable id such as:

`{device_uid}:{local_sequence}`

Rules:

- generated once at local write time
- never regenerated during retry
- never reused for a different mutation

This is the most important idempotency guarantee on the client side.

### 4. Deterministic Mutation Payload

Mutations should describe a deterministic target state, not a relative action,
whenever practical.

Good examples:

- set task state to `completed`
- replace subtask array with this payload
- replace habit document with this payload

Avoid phase-1 mutation types like:

- increment by 1
- toggle current value
- append based on unknown remote state

Deterministic payloads are much easier to retry safely.

### 5. Per-Entity Preconditions

Each mutation should include the last server change id the client believes the
entity was based on.

This is not for dedupe.
It is for convergence correctness.

If device A edits stale data while device B already committed a newer version,
the server can detect that explicitly instead of blindly pretending both writes
started from the same base.

## Server-Side Rules

### 1. Receipt Table For Deduplication

The server should persist a mutation receipt keyed by:

- `account_id`
- `device_uid`
- `mutation_id`

On duplicate push:

- do not reapply the mutation
- return the previously recorded result

This makes retries safe under:

- client timeout after commit
- repeated background retry
- network duplication

In implementation terms, this should be backed by a SQL unique constraint in
the default MySQL adapter and accessed through an infra repository implemented
with GORM.

### 2. Append-Only Canonical Ordering

Every accepted mutation gets a canonical `change_id`.

Rules:

- globally ordered within one account
- monotonic
- used as the pull cursor

This gives all devices the same replay order.

### 3. Entity-Level Version Check

Each entity stores its current `server_change_id`.

When a mutation arrives:

- if `base_change_id` matches current entity version, apply normally
- if it does not match, treat the mutation as stale and resolve
  deterministically

Phase-1 recommended behavior:

- reject stale mutation with a conflict result
- return latest authoritative entity state

That is often better than silent last-writer-wins because it avoids accidental
overwrite of newer remote work.

## Recommended Conflict Policy

To be more stable than WebDAV, do not use one single policy for everything.

### Deduplication Policy

- duplicates of the same mutation id must replay the same result

### Ordering Policy

- accepted mutations are serialized by server commit order

### Staleness Policy

- stale edits are detected by per-entity `base_change_id`
- stale edits should return a conflict result plus latest entity state

This is stronger than WebDAV because:

- retries are safe
- stale edits are explicit
- commits are globally ordered per account

## Push Contract

Recommended push behavior:

1. client sends pending mutations in local sequence order
2. server processes each mutation transactionally
3. for each mutation, server returns one of:
   - `applied`
   - `replayed`
   - `conflict`
   - `rejected`
4. client marks mutation queue entries accordingly
5. client always advances by performing a follow-up pull to the returned cursor

The key point:

- `replayed` is success
- `conflict` is not transport failure

## Pull Contract

Pull by cursor is naturally idempotent.

Rules:

- `PullChanges(cursor=X)` can be retried safely
- server returns all changes after `X` in canonical order
- applying the same pulled change twice must be safe locally

Local apply safety is helped by:

- comparing `server_change_id`
- only applying newer changes

## Why This Beats Blind LWW

Blind LWW makes retries easy to code, but it can hide real problems:

- duplicate mutation may reapply
- stale device may overwrite newer state
- user loses data without a clear reason

The stronger DeadlinerServer model should be:

- duplicate retry -> same receipt, no double write
- stale mutation -> explicit conflict, no silent overwrite
- successful mutation -> stable canonical change id

That is the main quality improvement over WebDAV.

## Recommended Entity Semantics

### DDL

- whole-document replace is acceptable in phase 1
- guarded by `base_change_id`

### Habit

- whole-document replace remains acceptable in phase 1
- guarded by carrier `ddl_uid` and `base_change_id`

### Habit Records

- do not implement free-form incremental counters first
- keep the existing replace-whole-document semantics for safety

This keeps the first system understandable.

## Local Recovery Flow

When client push gets `conflict`:

1. keep local unsynced content
2. pull latest server state
3. mark the local object as needing reconciliation
4. let client decide whether to auto-merge or ask user later

This avoids silent loss while still keeping sync operational.

## Minimal Phase-1 Guarantees

Phase 1 should guarantee:

1. the same mutation retried many times is applied at most once
2. the same pull retried many times yields the same ordered changes
3. stale writes are detected instead of silently overwriting newer remote state
4. offline edits remain queued durably until acknowledged

If we achieve those four things, the system is already materially steadier than
WebDAV snapshot sync.
