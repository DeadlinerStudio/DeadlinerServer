package entity

import "time"

type AccessTokenClaims struct {
	AccountUID string
	DeviceUID  string
	ExpiresAt  time.Time
}
