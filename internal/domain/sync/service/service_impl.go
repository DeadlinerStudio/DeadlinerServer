package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/aritxonly/deadlinerserver/internal/domain/account"
	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/command"
	portpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/port"
	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

type syncService struct {
	accountRepo         account.Repository
	deadlineRepo        portpkg.DeadlineRepository
	habitRepo           portpkg.HabitRepository
	mutationReceiptRepo portpkg.MutationReceiptRepository
	syncChangeRepo      portpkg.SyncChangeRepository
}

func NewService(
	accountRepo account.Repository,
	deadlineRepo portpkg.DeadlineRepository,
	habitRepo portpkg.HabitRepository,
	mutationReceiptRepo portpkg.MutationReceiptRepository,
	syncChangeRepo portpkg.SyncChangeRepository,
) Service {
	return &syncService{
		accountRepo:         accountRepo,
		deadlineRepo:        deadlineRepo,
		habitRepo:           habitRepo,
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

	habitChanges, err := s.habitRepo.ListAfterChangeID(
		ctx,
		acc.ID,
		afterChangeID,
		queryLimit,
		cmd.IncludeDelete,
	)
	if err != nil {
		return nil, err
	}

	ordered := mergeOrderedChanges(deadlineChanges, habitChanges)
	hasMore := false
	if limit > 0 && len(ordered) > limit {
		hasMore = true
		ordered = ordered[:limit]
	}

	nextCursor := cmd.Cursor
	if len(ordered) > 0 {
		lastChangeID := ordered[len(ordered)-1].changeID
		if lastChangeID > 0 {
			nextCursor = strconv.FormatInt(lastChangeID, 10)
		}
	}

	selectedDeadlineChanges, selectedHabitChanges := splitOrderedChanges(ordered)

	return &commandpkg.PullChangesResult{
		DeadlineChanges: selectedDeadlineChanges,
		HabitChanges:    selectedHabitChanges,
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
		HabitChanges:    make([]statepkg.HabitChange, 0, len(cmd.Mutations)),
		NextCursor:      cmd.BaseCursor,
	}

	maxCursor := parseCursor(cmd.BaseCursor)

	for _, mutation := range cmd.Mutations {
		mutationResult, deadlineChange, habitChange, err := s.handleMutation(ctx, acc, cmd.DeviceUID, mutation)
		if err != nil {
			return nil, err
		}

		result.Results = append(result.Results, mutationResult)
		if deadlineChange != nil {
			result.DeadlineChanges = append(result.DeadlineChanges, *deadlineChange)
			if deadlineChange.ServerVersion.ChangeID > maxCursor {
				maxCursor = deadlineChange.ServerVersion.ChangeID
			}
		} else if habitChange != nil {
			result.HabitChanges = append(result.HabitChanges, *habitChange)
			if habitChange.ServerVersion.ChangeID > maxCursor {
				maxCursor = habitChange.ServerVersion.ChangeID
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
) (commandpkg.MutationResult, *statepkg.DeadlineChange, *statepkg.HabitChange, error) {
	if mutation.DeviceUID != "" && mutation.DeviceUID != deviceUID {
		return rejectedMutationResult(mutation, "device uid mismatch"), nil, nil, nil
	}
	if mutation.Deadline == nil && mutation.Habit == nil {
		return rejectedMutationResult(mutation, "mutation payload is required"), nil, nil, nil
	}
	if mutation.Deadline != nil && mutation.Habit != nil {
		return rejectedMutationResult(mutation, "only one mutation payload variant is allowed"), nil, nil, nil
	}

	receipt, err := s.mutationReceiptRepo.Find(ctx, acc.ID, deviceUID, mutation.MutationID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}
	if receipt != nil {
		replayed := decodeMutationResult(receipt.ResultPayload)
		replayed.Replayed = true
		if replayed.Status == "" {
			replayed.Status = commandpkg.MutationStatusReplayed
		}
		var deadlineChange *statepkg.DeadlineChange
		var habitChange *statepkg.HabitChange
		if replayed.ServerVersion.ChangeID > 0 {
			if mutation.Deadline != nil {
				deadlineChange, err = s.deadlineRepo.FindByUID(ctx, acc.ID, mutation.EntityUID)
			} else {
				habitChange, err = s.habitRepo.FindByDDLUID(ctx, acc.ID, mutation.EntityUID)
			}
			if err != nil {
				return commandpkg.MutationResult{}, nil, nil, err
			}
		}
		return replayed, deadlineChange, habitChange, nil
	}

	if mutation.Deadline != nil {
		return s.handleDeadlineMutation(ctx, acc, deviceUID, mutation)
	}
	return s.handleHabitMutation(ctx, acc, deviceUID, mutation)
}

func (s *syncService) handleDeadlineMutation(
	ctx context.Context,
	acc *account.Account,
	deviceUID string,
	mutation commandpkg.Mutation,
) (commandpkg.MutationResult, *statepkg.DeadlineChange, *statepkg.HabitChange, error) {
	current, err := s.deadlineRepo.FindByUID(ctx, acc.ID, mutation.EntityUID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	if isEntityConflict(currentVersion(current), mutation.BaseChangeID) {
		conflict := conflictMutationResult(mutation, current)
		if err := s.saveReceipt(ctx, acc.ID, deviceUID, mutation, "deadline", conflict); err != nil {
			return commandpkg.MutationResult{}, nil, nil, err
		}
		return conflict, current, nil, nil
	}

	payload, err := json.Marshal(mutation.Deadline)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	change, err := s.syncChangeRepo.Append(ctx, portpkg.AppendSyncChangeParams{
		AccountID:  acc.ID,
		DeviceUID:  deviceUID,
		MutationID: mutation.MutationID,
		EntityKind: "deadline",
		EntityUID:  mutation.EntityUID,
		Action:     entityAction(mutation.Deadline.Deleted),
		Payload:    payload,
	})
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
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
		return commandpkg.MutationResult{}, nil, nil, err
	}

	savedChange, err := s.deadlineRepo.FindByUID(ctx, acc.ID, mutation.EntityUID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	applied := commandpkg.MutationResult{
		MutationID:    mutation.MutationID,
		EntityUID:     mutation.EntityUID,
		Accepted:      true,
		ServerVersion: serverVersion,
		Status:        commandpkg.MutationStatusApplied,
	}
	if err := s.saveReceipt(ctx, acc.ID, deviceUID, mutation, "deadline", applied); err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	return applied, savedChange, nil, nil
}

func (s *syncService) handleHabitMutation(
	ctx context.Context,
	acc *account.Account,
	deviceUID string,
	mutation commandpkg.Mutation,
) (commandpkg.MutationResult, *statepkg.DeadlineChange, *statepkg.HabitChange, error) {
	current, err := s.habitRepo.FindByDDLUID(ctx, acc.ID, mutation.EntityUID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	if isEntityConflict(currentVersion(current), mutation.BaseChangeID) {
		conflict := conflictMutationResult(mutation, current)
		if err := s.saveReceipt(ctx, acc.ID, deviceUID, mutation, "habit", conflict); err != nil {
			return commandpkg.MutationResult{}, nil, nil, err
		}
		return conflict, nil, current, nil
	}

	payload, err := json.Marshal(mutation.Habit)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	change, err := s.syncChangeRepo.Append(ctx, portpkg.AppendSyncChangeParams{
		AccountID:  acc.ID,
		DeviceUID:  deviceUID,
		MutationID: mutation.MutationID,
		EntityKind: "habit",
		EntityUID:  mutation.EntityUID,
		Action:     entityAction(mutation.Habit.Deleted),
		Payload:    payload,
	})
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	serverVersion := statepkg.ServerVersion{
		ChangeID:    change.ChangeID,
		CommittedAt: change.CommittedAt,
	}
	if err := s.habitRepo.Save(ctx, portpkg.SaveHabitParams{
		AccountID:          acc.ID,
		Deleted:            mutation.Habit.Deleted,
		Document:           mutation.Habit.Document,
		ServerVersion:      serverVersion,
		ClientVersion:      &mutation.ClientVersion,
		UpdatedByDeviceUID: deviceUID,
	}); err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	savedChange, err := s.habitRepo.FindByDDLUID(ctx, acc.ID, mutation.EntityUID)
	if err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	applied := commandpkg.MutationResult{
		MutationID:    mutation.MutationID,
		EntityUID:     mutation.EntityUID,
		Accepted:      true,
		ServerVersion: serverVersion,
		Status:        commandpkg.MutationStatusApplied,
	}
	if err := s.saveReceipt(ctx, acc.ID, deviceUID, mutation, "habit", applied); err != nil {
		return commandpkg.MutationResult{}, nil, nil, err
	}

	return applied, nil, savedChange, nil
}

func (s *syncService) saveReceipt(
	ctx context.Context,
	accountID int64,
	deviceUID string,
	mutation commandpkg.Mutation,
	entityKind string,
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
		EntityKind:     entityKind,
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

func conflictMutationResult(mutation commandpkg.Mutation, current any) commandpkg.MutationResult {
	result := commandpkg.MutationResult{
		MutationID:      mutation.MutationID,
		EntityUID:       mutation.EntityUID,
		Accepted:        false,
		RejectionReason: "stale base change id",
		Status:          commandpkg.MutationStatusConflict,
	}
	if current != nil {
		result.ServerVersion = currentVersion(current)
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

func isEntityConflict(currentVersion statepkg.ServerVersion, baseChangeID int64) bool {
	switch {
	case currentVersion.ChangeID == 0 && baseChangeID == 0:
		return false
	case currentVersion.ChangeID == 0 && baseChangeID > 0:
		return true
	case currentVersion.ChangeID > 0 && baseChangeID == 0:
		return true
	default:
		return currentVersion.ChangeID != baseChangeID
	}
}

func entityAction(deleted bool) string {
	if deleted {
		return "delete"
	}
	return "upsert"
}

type orderedChange struct {
	changeID int64
	deadline *statepkg.DeadlineChange
	habit    *statepkg.HabitChange
}

func mergeOrderedChanges(deadlineChanges []statepkg.DeadlineChange, habitChanges []statepkg.HabitChange) []orderedChange {
	ordered := make([]orderedChange, 0, len(deadlineChanges)+len(habitChanges))
	for i := range deadlineChanges {
		change := deadlineChanges[i]
		ordered = append(ordered, orderedChange{
			changeID: change.ServerVersion.ChangeID,
			deadline: &change,
		})
	}
	for i := range habitChanges {
		change := habitChanges[i]
		ordered = append(ordered, orderedChange{
			changeID: change.ServerVersion.ChangeID,
			habit:    &change,
		})
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].changeID < ordered[j].changeID
	})
	return ordered
}

func splitOrderedChanges(ordered []orderedChange) ([]statepkg.DeadlineChange, []statepkg.HabitChange) {
	deadlineChanges := make([]statepkg.DeadlineChange, 0, len(ordered))
	habitChanges := make([]statepkg.HabitChange, 0, len(ordered))
	for _, change := range ordered {
		if change.deadline != nil {
			deadlineChanges = append(deadlineChanges, *change.deadline)
		}
		if change.habit != nil {
			habitChanges = append(habitChanges, *change.habit)
		}
	}
	return deadlineChanges, habitChanges
}

func currentVersion(change any) statepkg.ServerVersion {
	switch typed := change.(type) {
	case *statepkg.DeadlineChange:
		if typed != nil {
			return typed.ServerVersion
		}
	case *statepkg.HabitChange:
		if typed != nil {
			return typed.ServerVersion
		}
	}
	return statepkg.ServerVersion{}
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
