package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

type SchemaInitializer interface {
	Init(*gorm.DB) error
}

type autoMigrateInitializer struct {
	models []interface{}
}

func (i autoMigrateInitializer) Init(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	return db.AutoMigrate(i.models...)
}

var schemaInitializers = map[string]SchemaInitializer{
	"mysql": autoMigrateInitializer{
		models: managedModels(),
	},
}

func RegisterSchemaInitializer(driver string, initializer SchemaInitializer) {
	schemaInitializers[driver] = initializer
}

func MustInit(driver, dsn string) *gorm.DB {
	db, err := Open(driver, dsn)
	if err != nil {
		panic(fmt.Errorf("open database: %w", err))
	}

	initializer, ok := schemaInitializers[driver]
	if !ok {
		panic(fmt.Errorf("unsupported schema initializer for driver: %s", driver))
	}

	if err := initializer.Init(db); err != nil {
		panic(fmt.Errorf("initialize schema for driver %s: %w", driver, err))
	}

	return db
}

func managedModels() []interface{} {
	return []interface{}{
		&AccountModel{},
		&DeviceModel{},
		&SessionModel{},
		&DeadlineItemModel{},
		&HabitDocModel{},
		&SyncChangeModel{},
		&MutationReceiptModel{},
	}
}
