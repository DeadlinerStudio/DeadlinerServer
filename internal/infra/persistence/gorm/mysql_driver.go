package gorm

import (
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	RegisterDialector("mysql", func(dsn string) gorm.Dialector {
		return mysqlDriver.Open(dsn)
	})
}
