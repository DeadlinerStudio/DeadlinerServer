package state

type MutationReceipt struct {
	AccountID      int64
	DeviceUID      string
	MutationID     string
	EntityKind     string
	EntityUID      string
	Status         string
	Replayed       bool
	ResultChangeID int64
	ResultPayload  []byte
}
