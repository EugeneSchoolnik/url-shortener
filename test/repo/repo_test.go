package repo_test

import (
	"testing"
	"url-shortener/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=url_shortener port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	// Auto migrate tables
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}
