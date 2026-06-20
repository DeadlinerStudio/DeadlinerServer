package sync

import domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"

type PullChangesInput struct {
	DeviceUID      string
	Cursor         string
	Limit          int32
	IncludeDeleted bool
}

type PushChangesInput struct {
	DeviceUID  string
	BaseCursor string
	Mutations  []domainSync.Mutation
}
