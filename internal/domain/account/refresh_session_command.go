package account

type RefreshSessionCommand struct {
	RefreshToken string
	DeviceUID    string
}

