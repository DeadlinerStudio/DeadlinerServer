package auth

import (
	"context"
	"strings"
)

type contextKey string

const (
	contextAccountUIDKey    contextKey = "deadliner.account_uid"
	contextAuthorizationKey contextKey = "deadliner.authorization"
	contextAccessTokenKey   contextKey = "deadliner.access_token"
)

type ContextAccountResolver struct {
	accessTokenParser AccessTokenParser
}

func NewContextAccountResolver(accessTokenParser AccessTokenParser) AccountResolver {
	return ContextAccountResolver{accessTokenParser: accessTokenParser}
}

func WithAccountUID(ctx context.Context, accountUID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, contextAccountUIDKey, accountUID)
}

func WithAuthorization(ctx context.Context, authorization string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, contextAuthorizationKey, authorization)
}

func WithAccessToken(ctx context.Context, accessToken string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, contextAccessTokenKey, accessToken)
}

func (r ContextAccountResolver) ResolveAccountUID(ctx context.Context) (string, error) {
	if value := strings.TrimSpace(contextStringValue(ctx, contextAccountUIDKey)); value != "" {
		return value, nil
	}

	if r.accessTokenParser != nil {
		for _, token := range []string{
			contextStringValue(ctx, contextAuthorizationKey),
			contextStringValue(ctx, contextAccessTokenKey),
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
				return "", err
			}
			if claims != nil && claims.AccountUID != "" {
				return claims.AccountUID, nil
			}
		}
	}

	return "", ErrUnauthenticated
}

func contextStringValue(ctx context.Context, key contextKey) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(key).(string)
	return value
}
