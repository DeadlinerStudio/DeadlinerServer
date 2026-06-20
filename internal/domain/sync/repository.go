package sync

import "context"

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

type Repository interface {
	FindMutationReceipt(context.Context, int64, string, string) (*MutationReceipt, error)
	SaveMutationReceipt(context.Context, *MutationReceipt) error
}
