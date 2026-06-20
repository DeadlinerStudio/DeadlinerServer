package gorm

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return r.db
	}
	return r.db.WithContext(ctx)
}
