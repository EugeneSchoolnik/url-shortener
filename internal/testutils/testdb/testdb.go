package testdb

import (
	"fmt"
	"log"
	"os"
	"testing"
	"url-shortener/internal/testutils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func New(t *testing.T) *gorm.DB {
	testutils.LoadTestEnv(t)

	dsn := os.Getenv("DB_DSN")

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sqlDB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	return db
}

func TruncateTables(t *testing.T, tables ...string) {
	t.Helper()
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
