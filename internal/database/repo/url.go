package repo

import (
	"url-shortener/internal/model"

	"gorm.io/gorm"
)

type UrlRepo struct {
	db *gorm.DB
}

func NewUrlRepo(db *gorm.DB) *UrlRepo {
	return &UrlRepo{db}
}

func (r *UrlRepo) Create(url *model.Url) error {
	return r.db.Create(url).Error
}

func (r *UrlRepo) Delete(id string, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Url{}).Error
}

func (r *UrlRepo) ByID(id string) (*model.Url, error) {
	var url model.Url

	return &url, r.db.Where("id = ?", id).First(&url).Error
}

func (r *UrlRepo) ByUserID(id string, limit int, offset int) ([]model.Url, error) {
	var urls []model.Url

	return urls, r.db.Where("user_id = ?", id).Limit(limit).Offset(offset).Find(&urls).Error
}
