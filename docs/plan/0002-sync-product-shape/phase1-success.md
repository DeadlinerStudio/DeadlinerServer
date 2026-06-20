# Phase 1 Success

## Future Product Opportunities

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

## Summary

Deadliner Account Sync should be:

- centralized in authority
- local-first in interaction
- offline-tolerant in operation
- clearly separate from WebDAV in product language
