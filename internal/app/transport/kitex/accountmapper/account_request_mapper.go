package accountmapper

import appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"

import v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"

func ToRegisterInput(req *v1.RegisterRequest) appaccount.RegisterInput {
	if req == nil {
		return appaccount.RegisterInput{}
	}

	return appaccount.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		DeviceUID:   req.DeviceUid,
		DeviceName:  req.DeviceName,
		Platform:    req.Platform,
	}
}

func ToLoginInput(req *v1.LoginRequest) appaccount.LoginInput {
	if req == nil {
		return appaccount.LoginInput{}
	}

	return appaccount.LoginInput{
		Email:      req.Email,
		Password:   req.Password,
		DeviceUID:  req.DeviceUid,
		DeviceName: req.DeviceName,
		Platform:   req.Platform,
	}
}

func ToRefreshSessionInput(req *v1.RefreshSessionRequest) appaccount.RefreshSessionInput {
	if req == nil {
		return appaccount.RefreshSessionInput{}
	}

	return appaccount.RefreshSessionInput{
		RefreshToken: req.RefreshToken,
		DeviceUID:    req.DeviceUid,
	}
}
