package service

import (
	"context"
	"errors"
	"testing"
	"time"

	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/account/command"
	entitypkg "github.com/aritxonly/deadlinerserver/internal/domain/account/entity"
)

func TestRegisterCreatesAccountDeviceAndSession(t *testing.T) {
	repo := newTestAccountRepo()
	service := NewService(
		repo,
		fakePasswordHasher{},
		fakeRefreshTokenHasher{},
		fakeAccessTokenCodec{},
		&fakeTokenGenerator{tokens: []string{"acc-token", "refresh-token", "sess-token"}},
		fakeClock{now: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)},
		time.Hour,
		24*time.Hour,
	)

	result, err := service.Register(context.Background(), commandpkg.RegisterCommand{
		Email:       "user@example.com",
		Password:    "secret",
		DisplayName: "User",
		DeviceUID:   "device-1",
		DeviceName:  "iPhone",
		Platform:    "ios",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if result.AccountUID == "" || result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatalf("expected full session bundle, got %+v", result)
	}
	if repo.accountByEmail["user@example.com"] == nil {
		t.Fatalf("expected account to be persisted")
	}
	if repo.devices["device-1"] == nil {
		t.Fatalf("expected device to be persisted")
	}
	if len(repo.sessionsByUID) != 1 {
		t.Fatalf("expected one session to be persisted, got %d", len(repo.sessionsByUID))
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	repo := newTestAccountRepo()
	repo.saveAccount(&entitypkg.Account{
		ID:           1,
		AccountUID:   "acc-1",
		Email:        "user@example.com",
		PasswordHash: "hashed:secret",
		DisplayName:  "User",
	})
	service := NewService(
		repo,
		fakePasswordHasher{},
		fakeRefreshTokenHasher{},
		fakeAccessTokenCodec{},
		&fakeTokenGenerator{tokens: []string{"refresh-token", "sess-token"}},
		fakeClock{now: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)},
		time.Hour,
		24*time.Hour,
	)

	_, err := service.Login(context.Background(), commandpkg.LoginCommand{
		Email:    "user@example.com",
		Password: "wrong",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestRefreshSessionRotatesRefreshToken(t *testing.T) {
	repo := newTestAccountRepo()
	repo.saveAccount(&entitypkg.Account{
		ID:           1,
		AccountUID:   "acc-1",
		Email:        "user@example.com",
		PasswordHash: "hashed:secret",
		DisplayName:  "User",
	})
	repo.sessionsByUID["sess-1"] = &entitypkg.Session{
		ID:               1,
		SessionUID:       "sess-1",
		AccountID:        1,
		DeviceUID:        "device-1",
		RefreshTokenHash: "hashed-token:old-refresh",
		ExpiresAt:        time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
	}
	service := NewService(
		repo,
		fakePasswordHasher{},
		fakeRefreshTokenHasher{},
		fakeAccessTokenCodec{},
		&fakeTokenGenerator{tokens: []string{"new-refresh"}},
		fakeClock{now: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)},
		time.Hour,
		24*time.Hour,
	)

	result, err := service.RefreshSession(context.Background(), commandpkg.RefreshSessionCommand{
		RefreshToken: "old-refresh",
		DeviceUID:    "device-1",
	})
	if err != nil {
		t.Fatalf("RefreshSession returned error: %v", err)
	}

	if result.RefreshToken != "new-refresh" {
		t.Fatalf("expected rotated refresh token, got %s", result.RefreshToken)
	}
	if repo.sessionsByUID["sess-1"].RefreshTokenHash != "hashed-token:new-refresh" {
		t.Fatalf("expected stored refresh token hash to rotate")
	}
}

type testAccountRepo struct {
	nextAccountID  int64
	nextSessionID  int64
	accountByEmail map[string]*entitypkg.Account
	accountByUID   map[string]*entitypkg.Account
	accountByID    map[int64]*entitypkg.Account
	devices        map[string]*entitypkg.Device
	sessionsByUID  map[string]*entitypkg.Session
}

func newTestAccountRepo() *testAccountRepo {
	return &testAccountRepo{
		nextAccountID:  1,
		nextSessionID:  1,
		accountByEmail: map[string]*entitypkg.Account{},
		accountByUID:   map[string]*entitypkg.Account{},
		accountByID:    map[int64]*entitypkg.Account{},
		devices:        map[string]*entitypkg.Device{},
		sessionsByUID:  map[string]*entitypkg.Session{},
	}
}

func (r *testAccountRepo) FindAccountByEmail(_ context.Context, email string) (*entitypkg.Account, error) {
	if acc := r.accountByEmail[email]; acc != nil {
		cloned := *acc
		return &cloned, nil
	}
	return nil, nil
}

func (r *testAccountRepo) FindAccountByUID(_ context.Context, uid string) (*entitypkg.Account, error) {
	if acc := r.accountByUID[uid]; acc != nil {
		cloned := *acc
		return &cloned, nil
	}
	return nil, nil
}

func (r *testAccountRepo) FindAccountByID(_ context.Context, id int64) (*entitypkg.Account, error) {
	if acc := r.accountByID[id]; acc != nil {
		cloned := *acc
		return &cloned, nil
	}
	return nil, nil
}

func (r *testAccountRepo) SaveAccount(_ context.Context, acc *entitypkg.Account) error {
	if acc.ID == 0 {
		acc.ID = r.nextAccountID
		r.nextAccountID++
	}
	r.saveAccount(acc)
	return nil
}

func (r *testAccountRepo) saveAccount(acc *entitypkg.Account) {
	cloned := *acc
	r.accountByEmail[acc.Email] = &cloned
	r.accountByUID[acc.AccountUID] = &cloned
	r.accountByID[acc.ID] = &cloned
}

func (r *testAccountRepo) SaveDevice(_ context.Context, device *entitypkg.Device) error {
	cloned := *device
	r.devices[device.DeviceUID] = &cloned
	return nil
}

func (r *testAccountRepo) SaveSession(_ context.Context, session *entitypkg.Session) error {
	if session.ID == 0 {
		if existing := r.sessionsByUID[session.SessionUID]; existing != nil {
			session.ID = existing.ID
		} else {
			session.ID = r.nextSessionID
			r.nextSessionID++
		}
	}
	cloned := *session
	r.sessionsByUID[session.SessionUID] = &cloned
	return nil
}

func (r *testAccountRepo) FindSessionByRefreshTokenHash(_ context.Context, hash string) (*entitypkg.Session, error) {
	for _, session := range r.sessionsByUID {
		if session.RefreshTokenHash == hash {
			cloned := *session
			return &cloned, nil
		}
	}
	return nil, nil
}

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (fakePasswordHasher) Compare(hash, password string) error {
	if hash != "hashed:"+password {
		return errors.New("password mismatch")
	}
	return nil
}

type fakeRefreshTokenHasher struct{}

func (fakeRefreshTokenHasher) Hash(token string) string {
	return "hashed-token:" + token
}

type fakeAccessTokenCodec struct{}

func (fakeAccessTokenCodec) Sign(claims entitypkg.AccessTokenClaims) (string, error) {
	return "access:" + claims.AccountUID + ":" + claims.DeviceUID, nil
}

func (fakeAccessTokenCodec) Parse(string) (*entitypkg.AccessTokenClaims, error) {
	return nil, nil
}

type fakeTokenGenerator struct {
	tokens []string
	index  int
}

func (g *fakeTokenGenerator) Generate() (string, error) {
	if g.index >= len(g.tokens) {
		return "", errors.New("no more tokens")
	}
	token := g.tokens[g.index]
	g.index++
	return token, nil
}

type fakeClock struct {
	now time.Time
}

func (c fakeClock) Now() time.Time {
	return c.now
}
