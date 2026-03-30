// internal/infrastructure/persistence/sqlserver/connection.go
package sqlserver

import (
	"fmt"
	"regexp"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	memo "tech-memo/internal/adapter/gateway/model/memo"
	todo "tech-memo/internal/adapter/gateway/model/todo"
)

func Open(dsn string) (*gorm.DB, error) {
	if err := ensureDatabase(dsn); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlserver: %w", err)
	}

	if err := db.AutoMigrate(
		&memo.MemoModel{},
		&todo.TodoModel{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}

// ensureDatabase は対象DBが存在しなければ master に接続して作成する
func ensureDatabase(dsn string) error {
	dbName, masterDSN := extractDBName(dsn)
	if dbName == "" {
		return nil
	}

	master, err := gorm.Open(sqlserver.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open master: %w", err)
	}
	sqlDB, _ := master.DB()
	defer sqlDB.Close()

	sql := fmt.Sprintf("IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = '%s') CREATE DATABASE [%s]", dbName, dbName)
	if err := master.Exec(sql).Error; err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	return nil
}

// extractDBName は DSN から database 名を取り出し、master 接続用 DSN を返す
func extractDBName(dsn string) (string, string) {
	re := regexp.MustCompile(`(?i)[?&]database=([^&]+)`)
	m := re.FindStringSubmatch(dsn)
	if m == nil {
		return "", dsn
	}
	dbName := m[1]
	masterDSN := re.ReplaceAllString(dsn, "")

	// 質問符がない場合は ? を起点にし、あれば & を追加
	sep := "?"
	if regexp.MustCompile(`\?`).MatchString(masterDSN) {
		sep = "&"
	}

	masterDSN = regexp.MustCompile(`[\?&]database=[^&]*`).ReplaceAllString(masterDSN, "")
	masterDSN = regexp.MustCompile(`\?&|&&`).ReplaceAllString(masterDSN, "?")
	masterDSN = regexp.MustCompile(`\?$|&$`).ReplaceAllString(masterDSN, "")

	masterDSN += sep + "database=master"
	return dbName, masterDSN
}
