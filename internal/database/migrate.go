package database

import (
	"url-shortener/internal/model"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	db.AutoMigrate(&model.User{})

	return nil
}
