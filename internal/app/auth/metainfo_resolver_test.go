package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	"github.com/bytedance/gopkg/cloud/metainfo"
)

func TestMetainfoAccountResolverResolveAccountUID(t *testing.T) {
	resolver := NewMetainfoAccountResolver(nil)
	ctx := metainfo.WithPersistentValue(context.Background(), MetaAccountUIDKey, "acc-1")

	accountUID, err := resolver.ResolveAccountUID(ctx)
	if err != nil {
		t.Fatalf("ResolveAccountUID returned error: %v", err)
	}
	if accountUID != "acc-1" {
		t.Fatalf("expected acc-1, got %s", accountUID)
	}
}

func TestMetainfoAccountResolverResolveAccountUIDFromBearerToken(t *testing.T) {
	resolver := NewMetainfoAccountResolver(fakeAccessTokenParser{
		claims: &domainAccount.AccessTokenClaims{
			AccountUID: "acc-2",
			DeviceUID:  "device-1",
			ExpiresAt:  time.Now().Add(time.Hour),
		},
	})
	ctx := metainfo.WithPersistentValue(context.Background(), MetaAuthorizationKey, "Bearer access-token")

	accountUID, err := resolver.ResolveAccountUID(ctx)
	if err != nil {
		t.Fatalf("ResolveAccountUID returned error: %v", err)
	}
	if accountUID != "acc-2" {
		t.Fatalf("expected acc-2, got %s", accountUID)
	}
}

func TestMetainfoAccountResolverResolveAccountUIDFromBearerTokenError(t *testing.T) {
	resolver := NewMetainfoAccountResolver(fakeAccessTokenParser{
		err: errors.New("bad token"),
	})
	ctx := metainfo.WithPersistentValue(context.Background(), MetaAuthorizationKey, "Bearer access-token")

	_, err := resolver.ResolveAccountUID(ctx)
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestMetainfoAccountResolverResolveAccountUIDMissing(t *testing.T) {
	resolver := NewMetainfoAccountResolver(nil)

	_, err := resolver.ResolveAccountUID(context.Background())
	if err == nil {
		t.Fatalf("expected missing metadata error")
	}
	if err != ErrUnauthenticated {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

type fakeAccessTokenParser struct {
	claims *domainAccount.AccessTokenClaims
	err    error
}

func (p fakeAccessTokenParser) Parse(string) (*domainAccount.AccessTokenClaims, error) {
	if p.err != nil {
		return nil, p.err
	}
	return p.claims, nil
}
