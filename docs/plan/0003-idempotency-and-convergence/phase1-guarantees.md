# Phase 1 Guarantees

## Minimal Guarantees

Phase 1 should guarantee:

1. the same mutation retried many times is applied at most once
2. the same pull retried many times yields the same ordered changes
3. stale writes are detected instead of silently overwriting newer remote state
4. offline edits remain queued durably until acknowledged

## Quality Bar Versus WebDAV

If we achieve those four guarantees, the system is already materially steadier
than WebDAV snapshot sync.

## Summary Policies

### Deduplication Policy

- duplicates of the same mutation id replay the same result

### Ordering Policy

- accepted mutations are serialized by server commit order

### Staleness Policy

- stale edits are detected by per-entity `base_change_id`
- stale edits return a conflict result plus latest entity state
