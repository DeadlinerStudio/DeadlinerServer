package sync

import "context"

type Service interface {
	PullChanges(context.Context, PullChangesCommand) (*PullChangesResult, error)
	PushChanges(context.Context, PushChangesCommand) (*PushChangesResult, error)
}
