package service

import (
	"context"

	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/command"
)

type Service interface {
	PullChanges(context.Context, commandpkg.PullChangesCommand) (*commandpkg.PullChangesResult, error)
	PushChanges(context.Context, commandpkg.PushChangesCommand) (*commandpkg.PushChangesResult, error)
}
