package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	transporthttp "github.com/aritxonly/deadlinerserver/internal/app/transport/http"
	httpmiddleware "github.com/aritxonly/deadlinerserver/internal/app/transport/http/middleware"
	"github.com/aritxonly/deadlinerserver/internal/config"
	domainaccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	domainsync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestSyncRouteRequiresAccessToken(t *testing.T) {
	engine := newTestServer(fakeAccountService{}, fakeSyncService{}, fakeTokenParser{})

	response := ut.PerformRequest(engine.Engine, "GET", "/v1/sync/pull?device_uid=device-1", nil)
	if response.Code != 401 {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if payload["error"] != "unauthorized" {
		t.Fatalf("unexpected error message: %s", payload["error"])
	}
	if payload["request_id"] == "" {
		t.Fatal("expected request_id in error response")
	}
	if response.Header().Get("X-Request-Id") == "" {
		t.Fatal("expected X-Request-Id header")
	}
}

func TestPostRouteRequiresJSONContentType(t *testing.T) {
	engine := newTestServer(fakeAccountService{}, fakeSyncService{}, fakeTokenParser{})

	response := ut.PerformRequest(
		engine.Engine,
		"POST",
		"/v1/auth/login",
		&ut.Body{Body: bytes.NewBufferString(`{"email":"a@example.com"}`), Len: len(`{"email":"a@example.com"}`)},
		ut.Header{Key: "Content-Type", Value: "text/plain"},
	)
	if response.Code != 415 {
		t.Fatalf("expected 415, got %d", response.Code)
	}
}

func TestBodyLimitRejectsLargeRequests(t *testing.T) {
	accountService := fakeAccountService{
		loginResult: &domainaccount.SessionBundle{
			AccountUID:   "acc-1",
			AccessToken:  "token",
			RefreshToken: "refresh",
			ExpiresAt:    time.Now().UTC().Format(time.RFC3339),
		},
	}
	cfg := config.Default()
	cfg.HTTP.MaxRequestBodyBytes = 8
	engine := newConfiguredTestServer(accountService, fakeSyncService{}, fakeTokenParser{}, cfg.HTTP)

	body := `{"email":"a@example.com","password":"123456","device_uid":"d1","device_name":"iPhone","platform":"ios"}`
	response := ut.PerformRequest(
		engine.Engine,
		"POST",
		"/v1/auth/login",
		&ut.Body{Body: bytes.NewBufferString(body), Len: len(body)},
		ut.Header{Key: "Content-Type", Value: "application/json"},
	)
	if response.Code != 413 {
		t.Fatalf("expected 413, got %d", response.Code)
	}
}

func TestInternalErrorsDoNotLeakImplementationDetails(t *testing.T) {
	engine := newTestServer(fakeAccountService{loginErr: errors.New("db down")}, fakeSyncService{}, fakeTokenParser{})

	body := `{"email":"a@example.com","password":"123456","device_uid":"d1","device_name":"iPhone","platform":"ios"}`
	response := ut.PerformRequest(
		engine.Engine,
		"POST",
		"/v1/auth/login",
		&ut.Body{Body: bytes.NewBufferString(body), Len: len(body)},
		ut.Header{Key: "Content-Type", Value: "application/json"},
	)
	if response.Code != 500 {
		t.Fatalf("expected 500, got %d", response.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if payload["error"] != "internal server error" {
		t.Fatalf("unexpected error body: %s", payload["error"])
	}
	if payload["request_id"] == "" {
		t.Fatal("expected request_id in 500 response")
	}
}

func TestSecurityHeadersAreAdded(t *testing.T) {
	engine := newTestServer(fakeAccountService{}, fakeSyncService{}, fakeTokenParser{})

	response := ut.PerformRequest(engine.Engine, "GET", "/healthz", nil)
	if response.Code != 200 {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("expected X-Content-Type-Options header")
	}
	if response.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatal("expected X-Frame-Options header")
	}
}

func newTestServer(
	accountService appaccount.Service,
	syncService appsync.Service,
	parser appauth.AccessTokenParser,
) *server.Hertz {
	return newConfiguredTestServer(accountService, syncService, parser, config.Default().HTTP)
}

func newConfiguredTestServer(
	accountService appaccount.Service,
	syncService appsync.Service,
	parser appauth.AccessTokenParser,
	httpConfig config.HTTPConfig,
) *server.Hertz {
	engine := server.New()
	engine.Use(
		httpmiddleware.RequestID(),
		httpmiddleware.Recovery(),
		httpmiddleware.SecurityHeaders(),
		httpmiddleware.EnforceMaxBodyBytes(httpConfig.MaxRequestBodyBytes),
		httpmiddleware.LimitByClientIP("http", httpConfig.RateLimitPerMinute, httpConfig.RateLimitBurst),
	)
	handler := transporthttp.NewHandler(accountService, syncService, parser, httpConfig)
	handler.RegisterRoutes(engine)
	return engine
}

type fakeTokenParser struct{}

func (fakeTokenParser) Parse(token string) (*domainaccount.AccessTokenClaims, error) {
	if token == "valid-token" {
		return &domainaccount.AccessTokenClaims{AccountUID: "acc-1"}, nil
	}
	return nil, errors.New("invalid token")
}

type fakeAccountService struct {
	loginResult *domainaccount.SessionBundle
	loginErr    error
}

func (s fakeAccountService) Register(context.Context, appaccount.RegisterInput) (*domainaccount.SessionBundle, error) {
	return &domainaccount.SessionBundle{
		AccountUID:   "acc-1",
		AccessToken:  "token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s fakeAccountService) Login(context.Context, appaccount.LoginInput) (*domainaccount.SessionBundle, error) {
	if s.loginErr != nil {
		return nil, s.loginErr
	}
	if s.loginResult != nil {
		return s.loginResult, nil
	}
	return &domainaccount.SessionBundle{
		AccountUID:   "acc-1",
		AccessToken:  "token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s fakeAccountService) RefreshSession(context.Context, appaccount.RefreshSessionInput) (*domainaccount.SessionBundle, error) {
	return &domainaccount.SessionBundle{
		AccountUID:   "acc-1",
		AccessToken:  "token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

type fakeSyncService struct{}

func (fakeSyncService) PullChanges(context.Context, appsync.PullChangesInput) (*domainsync.PullChangesResult, error) {
	return &domainsync.PullChangesResult{}, nil
}

func (fakeSyncService) PushChanges(context.Context, appsync.PushChangesInput) (*domainsync.PushChangesResult, error) {
	return &domainsync.PushChangesResult{}, nil
}
