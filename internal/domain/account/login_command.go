package account

type LoginCommand struct {
	Email      string
	Password   string
	DeviceUID  string
	DeviceName string
	Platform   string
}

