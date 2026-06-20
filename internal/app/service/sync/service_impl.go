package sync

import (
	"github.com/aritxonly/deadlinerserver/internal/app/auth"
	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
)

type service struct {
	accountResolver  auth.AccountResolver
	domainService    domainSync.Service
	defaultPullLimit int32
	maxPullLimit     int32
}

func NewService(
	accountResolver auth.AccountResolver,
	domainService domainSync.Service,
	defaultPullLimit int32,
	maxPullLimit int32,
) Service {
	return &service{
		accountResolver:  accountResolver,
		domainService:    domainService,
		defaultPullLimit: defaultPullLimit,
		maxPullLimit:     maxPullLimit,
	}
}

func (s *service) normalizePullLimit(limit int32) int32 {
	switch {
	case limit <= 0 && s.defaultPullLimit > 0:
		limit = s.defaultPullLimit
	case limit <= 0:
		return limit
	}

	if s.maxPullLimit > 0 && limit > s.maxPullLimit {
		return s.maxPullLimit
	}

	return limit
}
