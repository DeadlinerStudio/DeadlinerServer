package accountmapper

import (
	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

func ToRegisterResponse(bundle *domainAccount.SessionBundle) *v1.RegisterResponse {
	return &v1.RegisterResponse{
		Session: toKitexSession(bundle),
	}
}

func ToLoginResponse(bundle *domainAccount.SessionBundle) *v1.LoginResponse {
	return &v1.LoginResponse{
		Session: toKitexSession(bundle),
	}
}

func ToRefreshSessionResponse(bundle *domainAccount.SessionBundle) *v1.RefreshSessionResponse {
	return &v1.RefreshSessionResponse{
		Session: toKitexSession(bundle),
	}
}

func toKitexSession(bundle *domainAccount.SessionBundle) *v1.Session {
	if bundle == nil {
		return nil
	}

	return &v1.Session{
		AccountUid:   bundle.AccountUID,
		AccessToken:  bundle.AccessToken,
		RefreshToken: bundle.RefreshToken,
		ExpiresAt:    bundle.ExpiresAt,
	}
}
