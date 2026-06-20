package account

type RegisterCommand struct {
	Email       string
	Password    string
	DisplayName string
	DeviceUID   string
	DeviceName  string
	Platform    string
}

