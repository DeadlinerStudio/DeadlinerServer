package sync

type PullChangesCommand struct {
	AccountUID    string
	DeviceUID     string
	Cursor        string
	Limit         int32
	IncludeDelete bool
}

