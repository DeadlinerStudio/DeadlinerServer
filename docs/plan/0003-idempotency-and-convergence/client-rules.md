# Client Rules

## Why This Can Be Better Than WebDAV

WebDAV sync is file-oriented:

- clients race on shared remote snapshots
- retries are file rewrites
- conflict handling is mostly snapshot merge logic

The centralized backend can do better because it has:

- authenticated device identity
- append-only server commit order
- mutation-level deduplication
- durable acknowledgement receipts

## Core Rule

Every local edit becomes a durable mutation with a stable identity before it is
ever sent to the server.

## Required Client Rules

### Stable Device Identity

Each installed app instance has a stable `device_uid`.

### Durable Local Mutation Queue

Every local edit is written in two places:

1. local business tables
2. local mutation queue

The queue must survive app restarts, process kills, and offline periods.

### Immutable Mutation Identity

Each mutation gets a unique immutable id such as:

`{device_uid}:{local_sequence}`

Rules:

- generated once at local write time
- never regenerated during retry
- never reused for a different mutation

### Deterministic Mutation Payload

Prefer deterministic target-state mutations, such as:

- set task state to `completed`
- replace subtask array with this payload
- replace habit document with this payload

Avoid relative phase-1 mutation types like increment, toggle, or append against
unknown remote state.

### Per-Entity Preconditions

Each mutation includes the last server change id the client believes the entity
was based on.
