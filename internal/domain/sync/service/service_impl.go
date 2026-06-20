package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aritxonly/deadlinerserver/internal/domain/account"
	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/command"
	portpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/port"
	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

type syncService struct {
	accountRepo         account.Repository
	deadlineRepo        portpkg.DeadlineRepository
	mutationReceiptRepo portpkg.MutationReceiptRepository
	syncChangeRepo      portpkg.SyncChangeRepository
}

func NewService(
	accountRepo account.Repository,
	deadlineRepo portpkg.DeadlineRepository,
	mutationReceiptRepo portpkg.MutationReceiptRepository,
	syncChangeRepo portpkg.SyncChangeRepository,
) Service {
	return &syncService{
		accountRepo:         accountRepo,
		deadlineRepo:        deadlineRepo,
		mutationReceiptRepo: mutationReceiptRepo,
		syncChangeRepo:      syncChangeRepo,
	}
}

func (s *syncService) PullChanges(
	ctx context.Context,
	cmd commandpkg.PullChangesCommand,
) (*commandpkg.PullChangesResult, error) {
	acc, err := s.accountRepo.FindAccountByUID(ctx, cmd.AccountUID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, fmt.Errorf("account not found: %s", cmd.AccountUID)
	}

	afterChangeID := parseCursor(cmd.Cursor)
	limit := int(cmd.Limit)
	queryLimit := limit
	if limit > 0 {
		queryLimit = limit + 1
	}

	deadlineChanges, err := s.deadlineRepo.ListAfterChangeID(
		ctx,
		acc.ID,
		afterChangeID,
		queryLimit,
		cmd.IncludeDelete,
	)
	if err != nil {
		return nil, err
	}

	hasMore := false
	if limit > 0 && len(deadlineChanges) > limit {
		hasMore = true
		deadlineChanges = deadlineChanges[:limit]
	}

	nextCursor := cmd.Cursor
	if len(deadlineChanges) > 0 {
		lastChangeID := deadlineChanges[len(deadlineChanges)-1].ServerVersion.ChangeID
		if lastChangeID > 0 {
			nextCursor = strconv.FormatInt(lastChangeID, 10)
		}
	}

	return &commandpkg.PullChangesResult{
		DeadlineChanges: deadlineChanges,
		HabitChanges:    []statepkg.HabitChange{},
		NextCursor:      nextCursor,
		HasMore:         hasMore,
	}, nil
}

func (s *syncService) PushChanges(
	ctx context.Context,
	cmd commandpkg.PushChangesCommand,
) (*commandpkg.PushChangesResult, error) {
	acc, err := s.accountRepo.FindAccountByUID(ctx, cmd.AccountUID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, fmt.Errorf("account not found: %s", cmd.AccountUID)
	}

	result := &commandpkg.PushChangesResult{
		Results:         make([]commandpkg.MutationResult, 0, len(cmd.Mutations)),
		DeadlineChanges: make([]statepkg.DeadlineChange, 0, len(cmd.Mutations)),
		HabitChanges:    []statepkg.HabitChange{},
		NextCursor:      cmd.BaseCursor,
	}

	maxCursor := parseCursor(cmd.BaseCursor)

	for _, mutation := range cmd.Mutations {
		mutationResult, deadlineChange, err := s.handleMutation(ctx, acc, cmd.DeviceUID, mutation)
		if err != nil {
			return nil, err
		}

		result.Results = append(result.Results, mutationResult)
		if deadlineChange != nil {
			result.DeadlineChanges = append(result.DeadlineChanges, *deadlineChange)
			if deadlineChange.ServerVersion.ChangeID > maxCursor {
				maxCursor = deadlineChange.ServerVersion.ChangeID
			}
		} else if mutationResult.ServerVersion.ChangeID > maxCursor {
			maxCursor = mutationResult.ServerVersion.ChangeID
		}
	}

	if maxCursor > 0 {
		result.NextCursor = strconv.FormatInt(maxCursor, 10)
	}

	return result, nil
}

func (s *syncService) handleMutation(
	ctx context.Context,
	acc *account.Account,
	deviceUID string,
	mutation commandpkg.Mutation,
) (commandpkg.MutationResult, *statepkg.DeadlineChange, error) {
	if mutation.DeviceUID != "" && mutation.DeviceUID != deviceUID {
		return rejectedMutationResult(mutation, "device uid mismatch"), nil, nil
	}
	if mutation.Habit != nil {
		return rejectedMutationResult(mutation, "habit mutation is not implemented yet"), nil, nil
	}
	if mutation.Deadline == nil {
		return rejectedMutationResult(mutation, "deadline mutation payload is required"), nil, nil
	}

	receipt, err := s.mutationReceiptRepo.Find(ctx, acc.ID, deviceUID, mutation.MutationID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, err
	}
	if receipt != nil {
		replayed := decodeMutationResult(receipt.ResultPayload)
		replayed.Replayed = true
		if replayed.Status == "" {
			replayed.Status = commandpkg.MutationStatusReplayed
		}
		var change *statepkg.DeadlineChange
		if replayed.ServerVersion.ChangeID > 0 {
			change, err = s.deadlineRepo.FindByUID(ctx, acc.ID, mutation.EntityUID)
			if err != nil {
				return commandpkg.MutationResult{}, nil, err
			}
		}
		return replayed, change, nil
	}

	current, err := s.deadlineRepo.FindByUID(ctx, acc.ID, mutation.EntityUID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, err
	}

	if isDeadlineConflict(current, mutation.BaseChangeID) {
		conflict := conflictMutationResult(mutation, current)
		if err := s.saveReceipt(ctx, acc.ID, deviceUID, mutation, conflict); err != nil {
			return commandpkg.MutationResult{}, nil, err
		}
		return conflict, current, nil
	}

	payload, err := json.Marshal(mutation.Deadline)
	if err != nil {
		return commandpkg.MutationResult{}, nil, err
	}

	change, err := s.syncChangeRepo.Append(ctx, portpkg.AppendSyncChangeParams{
		AccountID:  acc.ID,
		DeviceUID:  deviceUID,
		MutationID: mutation.MutationID,
		EntityKind: "deadline",
		EntityUID:  mutation.EntityUID,
		Action:     deadlineAction(mutation.Deadline.Deleted),
		Payload:    payload,
	})
	if err != nil {
		return commandpkg.MutationResult{}, nil, err
	}

	serverVersion := statepkg.ServerVersion{
		ChangeID:    change.ChangeID,
		CommittedAt: change.CommittedAt,
	}
	if err := s.deadlineRepo.Save(ctx, portpkg.SaveDeadlineParams{
		AccountID:          acc.ID,
		Deleted:            mutation.Deadline.Deleted,
		Document:           mutation.Deadline.Document,
		ServerVersion:      serverVersion,
		ClientVersion:      &mutation.ClientVersion,
		UpdatedByDeviceUID: deviceUID,
	}); err != nil {
		return commandpkg.MutationResult{}, nil, err
	}

	savedChange, err := s.deadlineRepo.FindByUID(ctx, acc.ID, mutation.EntityUID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, err
	}

	applied := commandpkg.MutationResult{
		MutationID:    mutation.MutationID,
		EntityUID:     mutation.EntityUID,
		Accepted:      true,
		ServerVersion: serverVersion,
		Status:        commandpkg.MutationStatusApplied,
	}
	if err := s.saveReceipt(ctx, acc.ID, deviceUID, mutation, applied); err != nil {
		return commandpkg.MutationResult{}, nil, err
	}

	return applied, savedChange, nil
}

func (s *syncService) saveReceipt(
	ctx context.Context,
	accountID int64,
	deviceUID string,
	mutation commandpkg.Mutation,
	result commandpkg.MutationResult,
) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return s.mutationReceiptRepo.Save(ctx, &statepkg.MutationReceipt{
		AccountID:      accountID,
		DeviceUID:      deviceUID,
		MutationID:     mutation.MutationID,
		EntityKind:     "deadline",
		EntityUID:      mutation.EntityUID,
		Status:         result.Status,
		Replayed:       result.Replayed,
		ResultChangeID: result.ServerVersion.ChangeID,
		ResultPayload:  payload,
	})
}

func rejectedMutationResult(mutation commandpkg.Mutation, reason string) commandpkg.MutationResult {
	return commandpkg.MutationResult{
		MutationID:      mutation.MutationID,
		EntityUID:       mutation.EntityUID,
		Accepted:        false,
		RejectionReason: reason,
		Status:          commandpkg.MutationStatusRejected,
	}
}

func conflictMutationResult(
	mutation commandpkg.Mutation,
	current *statepkg.DeadlineChange,
) commandpkg.MutationResult {
	result := commandpkg.MutationResult{
		MutationID:      mutation.MutationID,
		EntityUID:       mutation.EntityUID,
		Accepted:        false,
		RejectionReason: "stale base change id",
		Status:          commandpkg.MutationStatusConflict,
	}
	if current != nil {
		result.ServerVersion = current.ServerVersion
	}
	return result
}

func decodeMutationResult(payload []byte) commandpkg.MutationResult {
	var result commandpkg.MutationResult
	if len(payload) == 0 {
		return result
	}
	_ = json.Unmarshal(payload, &result)
	return result
}

func isDeadlineConflict(current *statepkg.DeadlineChange, baseChangeID int64) bool {
	switch {
	case current == nil && baseChangeID == 0:
		return false
	case current == nil && baseChangeID > 0:
		return true
	case current != nil && baseChangeID == 0:
		return current.ServerVersion.ChangeID > 0
	default:
		return current.ServerVersion.ChangeID != baseChangeID
	}
}

func deadlineAction(deleted bool) string {
	if deleted {
		return "delete"
	}
	return "upsert"
}

func parseCursor(cursor string) int64 {
	if cursor == "" {
		return 0
	}
	value, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		return 0
	}
	return value
}
