package database

import (
	"url-shortener/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&model.User{}, &model.Url{})

	return nil
}
