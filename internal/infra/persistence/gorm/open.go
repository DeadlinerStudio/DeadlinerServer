package gorm

import (
	"fmt"
	"log"
	"os"
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
		Logger: logger.New(
			log.New(os.Stdout, "GORM ", log.LstdFlags|log.Lmicroseconds),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}
