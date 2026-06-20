package auth

import (
	"context"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

const (
	MetaAccountUIDKey       = "deadliner-account-uid"
	MetaLegacyAccountUIDKey = "account_uid"
	MetaAuthorizationKey    = "authorization"
	MetaAccessTokenKey      = "deadliner-access-token"
)

type MetainfoAccountResolver struct {
	accessTokenParser AccessTokenParser
}

func NewMetainfoAccountResolver(accessTokenParser AccessTokenParser) AccountResolver {
	return MetainfoAccountResolver{accessTokenParser: accessTokenParser}
}

func (r MetainfoAccountResolver) ResolveAccountUID(ctx context.Context) (string, error) {
	if r.accessTokenParser != nil {
		if accountUID, ok, err := r.resolveFromAccessToken(ctx); ok || err != nil {
			return accountUID, err
		}
	}

	for _, key := range []string{
		MetaAccountUIDKey,
		"x-" + MetaAccountUIDKey,
		MetaLegacyAccountUIDKey,
	} {
		if value, ok := metainfo.GetPersistentValue(ctx, key); ok && value != "" {
			return value, nil
		}
		if value, ok := metainfo.GetValue(ctx, key); ok && value != "" {
			return value, nil
		}
	}

	return "", ErrUnauthenticated
}

func (r MetainfoAccountResolver) resolveFromAccessToken(ctx context.Context) (string, bool, error) {
	for _, token := range []string{
		readMetaValue(ctx, MetaAuthorizationKey),
		readMetaValue(ctx, "Authorization"),
		readMetaValue(ctx, MetaAccessTokenKey),
	} {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = strings.TrimSpace(token[7:])
		}
		if token == "" {
			continue
		}

		claims, err := r.accessTokenParser.Parse(token)
		if err != nil {
			return "", true, err
		}
		if claims != nil && claims.AccountUID != "" {
			return claims.AccountUID, true, nil
		}
	}

	return "", false, nil
}

func readMetaValue(ctx context.Context, key string) string {
	if value, ok := metainfo.GetPersistentValue(ctx, key); ok && value != "" {
		return value
	}
	if value, ok := metainfo.GetValue(ctx, key); ok && value != "" {
		return value
	}
	return ""
}
