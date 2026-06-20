package gorm

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DialectorFactory func(dsn string) gorm.Dialector

var dialectorFactories = map[string]DialectorFactory{}

func RegisterDialector(driver string, factory DialectorFactory) {
	dialectorFactories[driver] = factory
}

func Open(driver, dsn string) (*gorm.DB, error) {
	factory, ok := dialectorFactories[driver]
	if !ok {
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	return gorm.Open(factory(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}
