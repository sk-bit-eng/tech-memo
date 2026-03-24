// internal/infrastructure/persistence/sqlserver/connection.go
package sqlserver

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"tech-memo/internal/adapter/gateway/model"
)

func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlserver: %w", err)
	}

	if err := db.AutoMigrate(
		&model.MemoModel{},
		&model.TodoModel{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}
