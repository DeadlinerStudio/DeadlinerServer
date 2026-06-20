package account

type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
	DeviceUID   string
	DeviceName  string
	Platform    string
}

type LoginInput struct {
	Email      string
	Password   string
	DeviceUID  string
	DeviceName string
	Platform   string
}

type RefreshSessionInput struct {
	RefreshToken string
	DeviceUID    string
}
