package state

type SyncChange struct {
	ChangeID    int64
	AccountID   int64
	DeviceUID   string
	MutationID  string
	EntityKind  string
	EntityUID   string
	Action      string
	Payload     []byte
	CommittedAt string
}
