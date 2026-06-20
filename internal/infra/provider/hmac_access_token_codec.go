package provider

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
)

var ErrInvalidAccessToken = errors.New("invalid access token")

type HMACAccessTokenCodec struct {
	secret []byte
	clock  SystemClock
}

type accessTokenPayload struct {
	AccountUID string `json:"account_uid"`
	DeviceUID  string `json:"device_uid"`
	ExpiresAt  int64  `json:"expires_at"`
}

func NewHMACAccessTokenCodec(secret string) HMACAccessTokenCodec {
	return HMACAccessTokenCodec{
		secret: []byte(secret),
		clock:  NewSystemClock(),
	}
}

func (c HMACAccessTokenCodec) Sign(claims domainAccount.AccessTokenClaims) (string, error) {
	payload, err := json.Marshal(accessTokenPayload{
		AccountUID: claims.AccountUID,
		DeviceUID:  claims.DeviceUID,
		ExpiresAt:  claims.ExpiresAt.UTC().Unix(),
	})
	if err != nil {
		return "", err
	}

	payloadPart := base64.RawURLEncoding.EncodeToString(payload)
	sigPart := base64.RawURLEncoding.EncodeToString(c.sign(payloadPart))
	return payloadPart + "." + sigPart, nil
}

func (c HMACAccessTokenCodec) Parse(token string) (*domainAccount.AccessTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidAccessToken
	}
	if !hmac.Equal(c.sign(parts[0]), decodeBase64(parts[1])) {
		return nil, ErrInvalidAccessToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidAccessToken
	}

	var payload accessTokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrInvalidAccessToken
	}

	expiresAt := time.Unix(payload.ExpiresAt, 0).UTC()
	if c.clock.Now().After(expiresAt) {
		return nil, ErrInvalidAccessToken
	}

	return &domainAccount.AccessTokenClaims{
		AccountUID: payload.AccountUID,
		DeviceUID:  payload.DeviceUID,
		ExpiresAt:  expiresAt,
	}, nil
}

func (c HMACAccessTokenCodec) sign(payloadPart string) []byte {
	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write([]byte(payloadPart))
	return mac.Sum(nil)
}

func decodeBase64(value string) []byte {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil
	}
	return decoded
}
