# Positioning

## Product Positioning

Deadliner should support multiple sync modes, not one forced backend model.

Recommended product shape:

1. `Local Only`
2. `WebDAV Sync`
3. `Deadliner Account Sync`

The Go backend is the third mode.

## User-Facing Narrative

The clean story should be:

- if you want simple, reliable, account-based sync, use Deadliner Account Sync
- if you want bring-your-own storage and no dedicated backend account, use
  WebDAV
- if you do not want any sync, stay on Local Only

## Why Keep WebDAV

WebDAV still has value:

- user-controlled storage
- no permanent backend dependency
- easy personal backup mental model
- good fit for technically confident users

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
