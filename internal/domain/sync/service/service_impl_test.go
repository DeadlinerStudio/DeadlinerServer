package service

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/aritxonly/deadlinerserver/internal/domain/account"
	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/command"
	documentpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/document"
	portpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/port"
	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

func TestPushChangesApplied(t *testing.T) {
	service, deadlineRepo, _, receiptRepo, changeRepo := newTestSyncService()
	result, err := service.PushChanges(context.Background(), commandpkg.PushChangesCommand{
		AccountUID: "acc-1",
		DeviceUID:  "device-1",
		Mutations: []commandpkg.Mutation{{
			MutationID: "device-1:1",
			DeviceUID:  "device-1",
			EntityUID:  "ddl-1",
			Deadline: &commandpkg.DeadlinePatch{
				Document: documentpkg.DeadlineDocument{
					UID: "ddl-1", LegacyID: 1, Name: "Write report",
					StartTime: "2026-06-20T08:00:00Z", EndTime: "2026-06-20T18:00:00Z",
					State: documentpkg.DeadlineStateActive, Type: documentpkg.DeadlineTypeTask, SubTasks: []documentpkg.SubTask{},
				},
			},
		}},
	})
	if err != nil || len(result.Results) != 1 || result.Results[0].Status != commandpkg.MutationStatusApplied {
		t.Fatalf("unexpected push result: %+v err=%v", result, err)
	}
	if len(result.DeadlineChanges) != 1 || deadlineRepo.saved["ddl-1"].Document.Name != "Write report" {
		t.Fatalf("deadline was not persisted")
	}
	if _, ok := receiptRepo.saved["device-1:1"]; !ok || len(changeRepo.saved) != 1 {
		t.Fatalf("expected receipt and sync change to be saved")
	}
}

func TestPushChangesReplay(t *testing.T) {
	service, _, _, receiptRepo, _ := newTestSyncService()
	receiptRepo.saved["device-1:1"] = &statepkg.MutationReceipt{
		AccountID: 1, DeviceUID: "device-1", MutationID: "device-1:1", EntityUID: "ddl-1",
		Status: commandpkg.MutationStatusApplied,
		ResultPayload: mustJSON(t, commandpkg.MutationResult{
			MutationID: "device-1:1", EntityUID: "ddl-1", Accepted: true,
			ServerVersion: statepkg.ServerVersion{ChangeID: 9, CommittedAt: "2026-06-20T12:00:00Z"},
			Status:        commandpkg.MutationStatusApplied,
		}),
	}
	result, err := service.PushChanges(context.Background(), commandpkg.PushChangesCommand{
		AccountUID: "acc-1", DeviceUID: "device-1",
		Mutations: []commandpkg.Mutation{{MutationID: "device-1:1", DeviceUID: "device-1", EntityUID: "ddl-1", Deadline: &commandpkg.DeadlinePatch{Document: documentpkg.DeadlineDocument{UID: "ddl-1"}}}},
	})
	if err != nil || !result.Results[0].Replayed {
		t.Fatalf("expected replayed result, got %+v err=%v", result, err)
	}
}

func TestPullChangesFiltersDeletedAndReportsHasMore(t *testing.T) {
	service, deadlineRepo, habitRepo, _, _ := newTestSyncService()
	deadlineRepo.saved["ddl-1"] = &statepkg.DeadlineChange{EntityUID: "ddl-1", ServerVersion: statepkg.ServerVersion{ChangeID: 1}, Document: documentpkg.DeadlineDocument{UID: "ddl-1", Name: "Keep active", State: documentpkg.DeadlineStateActive, Type: documentpkg.DeadlineTypeTask}}
	deadlineRepo.saved["ddl-2"] = &statepkg.DeadlineChange{EntityUID: "ddl-2", Deleted: true, ServerVersion: statepkg.ServerVersion{ChangeID: 2}, Document: documentpkg.DeadlineDocument{UID: "ddl-2", Name: "Deleted", State: documentpkg.DeadlineStateCompleted, Type: documentpkg.DeadlineTypeTask}}
	deadlineRepo.saved["ddl-3"] = &statepkg.DeadlineChange{EntityUID: "ddl-3", ServerVersion: statepkg.ServerVersion{ChangeID: 3}, Document: documentpkg.DeadlineDocument{UID: "ddl-3", Name: "Still active", State: documentpkg.DeadlineStateActive, Type: documentpkg.DeadlineTypeTask}}
	habitRepo.saved["ddl-h1"] = &statepkg.HabitChange{EntityUID: "ddl-h1", ServerVersion: statepkg.ServerVersion{ChangeID: 4}, Document: documentpkg.HabitDocument{DDLUID: "ddl-h1"}}
	result, err := service.PullChanges(context.Background(), commandpkg.PullChangesCommand{
		AccountUID: "acc-1", DeviceUID: "device-1", Cursor: "0", Limit: 1, IncludeDelete: false,
	})
	if err != nil || len(result.DeadlineChanges) != 1 || result.DeadlineChanges[0].EntityUID != "ddl-1" || !result.HasMore {
		t.Fatalf("unexpected pull result: %+v err=%v", result, err)
	}
}

func TestPushChangesHabitApplied(t *testing.T) {
	service, _, habitRepo, receiptRepo, changeRepo := newTestSyncService()
	result, err := service.PushChanges(context.Background(), commandpkg.PushChangesCommand{
		AccountUID: "acc-1",
		DeviceUID:  "device-1",
		Mutations: []commandpkg.Mutation{{
			MutationID: "device-1:habit-1",
			DeviceUID:  "device-1",
			EntityUID:  "ddl-habit-1",
			Habit: &commandpkg.HabitPatch{
				Document: documentpkg.HabitDocument{
					DDLUID: "ddl-habit-1",
					Habit: documentpkg.HabitConfig{
						Name:           "Walk",
						Period:         documentpkg.HabitPeriodDaily,
						GoalType:       documentpkg.HabitGoalTypePerPeriod,
						Status:         documentpkg.HabitStatusActive,
						TimesPerPeriod: 1,
					},
				},
			},
		}},
	})
	if err != nil || len(result.Results) != 1 || result.Results[0].Status != commandpkg.MutationStatusApplied {
		t.Fatalf("unexpected push result: %+v err=%v", result, err)
	}
	if len(result.HabitChanges) != 1 || habitRepo.saved["ddl-habit-1"].Document.Habit.Name != "Walk" {
		t.Fatalf("habit was not persisted")
	}
	if _, ok := receiptRepo.saved["device-1:habit-1"]; !ok || len(changeRepo.saved) != 1 || changeRepo.saved[0].EntityKind != "habit" {
		t.Fatalf("expected habit receipt and sync change to be saved")
	}
}

func TestPullChangesMergesDeadlineAndHabitByCursor(t *testing.T) {
	service, deadlineRepo, habitRepo, _, _ := newTestSyncService()
	deadlineRepo.saved["ddl-1"] = &statepkg.DeadlineChange{
		EntityUID:     "ddl-1",
		ServerVersion: statepkg.ServerVersion{ChangeID: 1},
		Document: documentpkg.DeadlineDocument{
			UID: "ddl-1", Name: "Deadline one", State: documentpkg.DeadlineStateActive, Type: documentpkg.DeadlineTypeTask,
		},
	}
	habitRepo.saved["ddl-h1"] = &statepkg.HabitChange{
		EntityUID:     "ddl-h1",
		ServerVersion: statepkg.ServerVersion{ChangeID: 2},
		Document: documentpkg.HabitDocument{
			DDLUID: "ddl-h1",
			Habit:  documentpkg.HabitConfig{Name: "Habit one", Period: documentpkg.HabitPeriodDaily, GoalType: documentpkg.HabitGoalTypePerPeriod, Status: documentpkg.HabitStatusActive},
		},
	}
	deadlineRepo.saved["ddl-2"] = &statepkg.DeadlineChange{
		EntityUID:     "ddl-2",
		ServerVersion: statepkg.ServerVersion{ChangeID: 3},
		Document: documentpkg.DeadlineDocument{
			UID: "ddl-2", Name: "Deadline two", State: documentpkg.DeadlineStateActive, Type: documentpkg.DeadlineTypeTask,
		},
	}

	result, err := service.PullChanges(context.Background(), commandpkg.PullChangesCommand{
		AccountUID: "acc-1", DeviceUID: "device-1", Cursor: "0", Limit: 2, IncludeDelete: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NextCursor != "2" || !result.HasMore {
		t.Fatalf("unexpected cursor result: %+v", result)
	}
	if len(result.DeadlineChanges) != 1 || result.DeadlineChanges[0].EntityUID != "ddl-1" {
		t.Fatalf("unexpected deadline changes: %+v", result.DeadlineChanges)
	}
	if len(result.HabitChanges) != 1 || result.HabitChanges[0].EntityUID != "ddl-h1" {
		t.Fatalf("unexpected habit changes: %+v", result.HabitChanges)
	}
}

type testAccountRepo struct{}

func (r *testAccountRepo) FindAccountByEmail(context.Context, string) (*account.Account, error) {
	return nil, nil
}
func (r *testAccountRepo) FindAccountByUID(_ context.Context, uid string) (*account.Account, error) {
	return &account.Account{ID: 1, AccountUID: uid}, nil
}
func (r *testAccountRepo) FindAccountByID(_ context.Context, id int64) (*account.Account, error) {
	return &account.Account{ID: id, AccountUID: "acc-1"}, nil
}
func (r *testAccountRepo) FindSessionByRefreshTokenHash(context.Context, string) (*account.Session, error) {
	return nil, nil
}
func (r *testAccountRepo) SaveAccount(context.Context, *account.Account) error { return nil }
func (r *testAccountRepo) SaveDevice(context.Context, *account.Device) error   { return nil }
func (r *testAccountRepo) SaveSession(context.Context, *account.Session) error { return nil }

type testDeadlineRepo struct {
	saved map[string]*statepkg.DeadlineChange
}

func (r *testDeadlineRepo) FindByUID(_ context.Context, _ int64, uid string) (*statepkg.DeadlineChange, error) {
	if change, ok := r.saved[uid]; ok {
		cloned := *change
		return &cloned, nil
	}
	return nil, nil
}
func (r *testDeadlineRepo) Save(_ context.Context, params portpkg.SaveDeadlineParams) error {
	r.saved[params.Document.UID] = &statepkg.DeadlineChange{EntityUID: params.Document.UID, Deleted: params.Deleted, ServerVersion: params.ServerVersion, Document: params.Document}
	return nil
}
func (r *testDeadlineRepo) ListAfterChangeID(_ context.Context, _ int64, afterChangeID int64, limit int, includeDeleted bool) ([]statepkg.DeadlineChange, error) {
	changes := make([]statepkg.DeadlineChange, 0, len(r.saved))
	for _, change := range r.saved {
		if change.ServerVersion.ChangeID <= afterChangeID || (!includeDeleted && change.Deleted) {
			continue
		}
		changes = append(changes, *change)
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].ServerVersion.ChangeID < changes[j].ServerVersion.ChangeID })
	if limit > 0 && len(changes) > limit {
		changes = changes[:limit]
	}
	return changes, nil
}

type testHabitRepo struct {
	saved map[string]*statepkg.HabitChange
}

func (r *testHabitRepo) FindByDDLUID(_ context.Context, _ int64, ddlUID string) (*statepkg.HabitChange, error) {
	if change, ok := r.saved[ddlUID]; ok {
		cloned := *change
		return &cloned, nil
	}
	return nil, nil
}
func (r *testHabitRepo) Save(_ context.Context, params portpkg.SaveHabitParams) error {
	r.saved[params.Document.DDLUID] = &statepkg.HabitChange{
		EntityUID:     params.Document.DDLUID,
		Deleted:       params.Deleted,
		ServerVersion: params.ServerVersion,
		Document:      params.Document,
	}
	return nil
}
func (r *testHabitRepo) ListAfterChangeID(_ context.Context, _ int64, afterChangeID int64, limit int, includeDeleted bool) ([]statepkg.HabitChange, error) {
	changes := make([]statepkg.HabitChange, 0, len(r.saved))
	for _, change := range r.saved {
		if change.ServerVersion.ChangeID <= afterChangeID || (!includeDeleted && change.Deleted) {
			continue
		}
		changes = append(changes, *change)
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].ServerVersion.ChangeID < changes[j].ServerVersion.ChangeID })
	if limit > 0 && len(changes) > limit {
		changes = changes[:limit]
	}
	return changes, nil
}

type testMutationReceiptRepo struct {
	saved map[string]*statepkg.MutationReceipt
}

func (r *testMutationReceiptRepo) Find(_ context.Context, _ int64, _ string, mutationID string) (*statepkg.MutationReceipt, error) {
	if receipt, ok := r.saved[mutationID]; ok {
		cloned := *receipt
		cloned.ResultPayload = append([]byte(nil), receipt.ResultPayload...)
		return &cloned, nil
	}
	return nil, nil
}
func (r *testMutationReceiptRepo) Save(_ context.Context, receipt *statepkg.MutationReceipt) error {
	cloned := *receipt
	cloned.ResultPayload = append([]byte(nil), receipt.ResultPayload...)
	r.saved[receipt.MutationID] = &cloned
	return nil
}

type testSyncChangeRepo struct {
	saved  []*statepkg.SyncChange
	nextID int64
}

func (r *testSyncChangeRepo) Append(_ context.Context, params portpkg.AppendSyncChangeParams) (*statepkg.SyncChange, error) {
	r.nextID++
	change := &statepkg.SyncChange{ChangeID: r.nextID, AccountID: params.AccountID, DeviceUID: params.DeviceUID, MutationID: params.MutationID, EntityKind: params.EntityKind, EntityUID: params.EntityUID, Action: params.Action, Payload: params.Payload, CommittedAt: "2026-06-20T12:00:00Z"}
	r.saved = append(r.saved, change)
	return change, nil
}
func (r *testSyncChangeRepo) ListAfterChangeID(context.Context, int64, int64, int) ([]statepkg.SyncChange, error) {
	return nil, nil
}

func newTestSyncService() (Service, *testDeadlineRepo, *testHabitRepo, *testMutationReceiptRepo, *testSyncChangeRepo) {
	deadlineRepo := &testDeadlineRepo{saved: map[string]*statepkg.DeadlineChange{}}
	habitRepo := &testHabitRepo{saved: map[string]*statepkg.HabitChange{}}
	receiptRepo := &testMutationReceiptRepo{saved: map[string]*statepkg.MutationReceipt{}}
	changeRepo := &testSyncChangeRepo{}
	return NewService(&testAccountRepo{}, deadlineRepo, habitRepo, receiptRepo, changeRepo), deadlineRepo, habitRepo, receiptRepo, changeRepo
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	return payload
}
