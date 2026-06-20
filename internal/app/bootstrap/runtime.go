package bootstrap

import (
	"time"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	transportkitex "github.com/aritxonly/deadlinerserver/internal/app/transport/kitex"
	"github.com/aritxonly/deadlinerserver/internal/config"
	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
	"github.com/aritxonly/deadlinerserver/internal/infra/provider"
	"github.com/aritxonly/deadlinerserver/internal/infra/repo"
)

func NewKitexHandler(cfg config.Config) (*transportkitex.Handler, error) {
	db, err := persistencegorm.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	accountRepo := repo.NewAccountRepo(db)
	deadlineRepo := repo.NewDeadlineRepo(db)
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
		mutationReceiptRepo,
		syncChangeRepo,
	)
	syncAppService := appsync.NewService(
		appauth.NewMetainfoAccountResolver(accessTokenCodec),
		syncDomainService,
		cfg.Sync.DefaultPullLimit,
		cfg.Sync.MaxPullLimit,
	)

	return transportkitex.NewHandler(accountAppService, syncAppService), nil
}
