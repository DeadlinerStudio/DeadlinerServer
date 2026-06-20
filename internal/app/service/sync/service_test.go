package sync

import (
	"context"
	"errors"
	"testing"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
)

func TestPullChangesUsesResolvedAccountAndNormalizesLimit(t *testing.T) {
	domainService := &fakeDomainSyncService{
		pullResult: &domainSync.PullChangesResult{},
	}
	service := NewService(
		fakeAccountResolver{accountUID: "acc-1"},
		domainService,
		100,
		500,
	)

	_, err := service.PullChanges(context.Background(), PullChangesInput{
		DeviceUID:      "device-1",
		Cursor:         "9",
		Limit:          0,
		IncludeDeleted: true,
	})
	if err != nil {
		t.Fatalf("PullChanges returned error: %v", err)
	}

	if domainService.lastPullCommand.AccountUID != "acc-1" {
		t.Fatalf("expected account uid acc-1, got %s", domainService.lastPullCommand.AccountUID)
	}
	if domainService.lastPullCommand.Limit != 100 {
		t.Fatalf("expected normalized limit 100, got %d", domainService.lastPullCommand.Limit)
	}
}

func TestPushChangesForwardsMutations(t *testing.T) {
	domainService := &fakeDomainSyncService{
		pushResult: &domainSync.PushChangesResult{},
	}
	service := NewService(
		fakeAccountResolver{accountUID: "acc-1"},
		domainService,
		100,
		500,
	)

	mutations := []domainSync.Mutation{{MutationID: "m-1", EntityUID: "ddl-1"}}
	_, err := service.PushChanges(context.Background(), PushChangesInput{
		DeviceUID:  "device-1",
		BaseCursor: "7",
		Mutations:  mutations,
	})
	if err != nil {
		t.Fatalf("PushChanges returned error: %v", err)
	}

	if len(domainService.lastPushCommand.Mutations) != 1 {
		t.Fatalf("expected 1 mutation, got %d", len(domainService.lastPushCommand.Mutations))
	}
	if domainService.lastPushCommand.AccountUID != "acc-1" {
		t.Fatalf("expected account uid acc-1, got %s", domainService.lastPushCommand.AccountUID)
	}
}

func TestPushChangesReturnsResolverError(t *testing.T) {
	service := NewService(
		fakeAccountResolver{err: errors.New("missing account")},
		&fakeDomainSyncService{},
		100,
		500,
	)

	_, err := service.PushChanges(context.Background(), PushChangesInput{})
	if err == nil {
		t.Fatalf("expected resolver error")
	}
}

type fakeAccountResolver struct {
	accountUID string
	err        error
}

func (r fakeAccountResolver) ResolveAccountUID(context.Context) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.accountUID, nil
}

type fakeDomainSyncService struct {
	lastPullCommand domainSync.PullChangesCommand
	lastPushCommand domainSync.PushChangesCommand
	pullResult      *domainSync.PullChangesResult
	pushResult      *domainSync.PushChangesResult
}

func (s *fakeDomainSyncService) PullChanges(
	_ context.Context,
	cmd domainSync.PullChangesCommand,
) (*domainSync.PullChangesResult, error) {
	s.lastPullCommand = cmd
	return s.pullResult, nil
}

func (s *fakeDomainSyncService) PushChanges(
	_ context.Context,
	cmd domainSync.PushChangesCommand,
) (*domainSync.PushChangesResult, error) {
	s.lastPushCommand = cmd
	return s.pushResult, nil
}
