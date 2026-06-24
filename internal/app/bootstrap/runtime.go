package bootstrap

import (
	"time"

	appadmin "github.com/aritxonly/deadlinerserver/internal/app/adminconfig"
	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	"github.com/aritxonly/deadlinerserver/internal/config"
	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
	"github.com/aritxonly/deadlinerserver/internal/infra/provider"
	"github.com/aritxonly/deadlinerserver/internal/infra/repo"
)

type Runtime struct {
	AccountService     appaccount.Service
	AdminConfigService appadmin.Service
	SyncDomain         domainSync.Service
	AccessTokenCodec   appauth.AccessTokenParser
	DefaultPullLimit   int32
	MaxPullLimit       int32
	AdminRuntimeConfig config.AdminConfig
}

func NewRuntime(cfg config.Config) (*Runtime, error) {
	db := persistencegorm.MustInit(cfg.Database.Driver, cfg.Database.DSN)

	accountRepo := repo.NewAccountRepo(db)
	deadlineRepo := repo.NewDeadlineRepo(db)
	habitRepo := repo.NewHabitRepo(db)
	mutationReceiptRepo := repo.NewMutationReceiptRepo(db)
	syncChangeRepo := repo.NewSyncChangeRepo(db)

	accessTokenCodec := provider.NewHMACAccessTokenCodec(cfg.Auth.AccessTokenSecret)

	accountDomainService := domainAccount.NewService(
		accountRepo,
		provider.NewBcryptPasswordHasher(cfg.Auth.PasswordHashCost),
		provider.NewSHA256TokenHasher(),
		accessTokenCodec,
		provider.NewRandomTokenGenerator(cfg.Auth.RandomTokenBytes),
		provider.NewSystemClock(),
		time.Duration(cfg.Auth.AccessTokenTTLMinutes)*time.Minute,
		time.Duration(cfg.Auth.RefreshTokenTTLHours)*time.Hour,
	)
	accountAppService := appaccount.NewService(accountDomainService)

	syncDomainService := domainSync.NewService(
		accountRepo,
		deadlineRepo,
		habitRepo,
		mutationReceiptRepo,
		syncChangeRepo,
	)

	return &Runtime{
		AccountService:     accountAppService,
		AdminConfigService: appadmin.NewService(config.DefaultPath, config.ResolveSecretPath()),
		SyncDomain:         syncDomainService,
		AccessTokenCodec:   accessTokenCodec,
		DefaultPullLimit:   cfg.Sync.DefaultPullLimit,
		MaxPullLimit:       cfg.Sync.MaxPullLimit,
		AdminRuntimeConfig: cfg.Admin,
	}, nil
}

func (r *Runtime) NewKitexSyncService() appsync.Service {
	return appsync.NewService(
		appauth.NewMetainfoAccountResolver(r.AccessTokenCodec),
		r.SyncDomain,
		r.DefaultPullLimit,
		r.MaxPullLimit,
	)
}

func (r *Runtime) NewHTTPSyncService() appsync.Service {
	return appsync.NewService(
		appauth.NewContextAccountResolver(r.AccessTokenCodec),
		r.SyncDomain,
		r.DefaultPullLimit,
		r.MaxPullLimit,
	)
}
