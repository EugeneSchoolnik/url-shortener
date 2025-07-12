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

// LinkByID also increment total hits
func (r *UrlRepo) LinkByID(id string) (string, error) {
	var link string

	res := r.db.Raw(`
	UPDATE urls
	SET total_hits = total_hits + 1
	WHERE id = ?
	RETURNING link;
`, id).Scan(&link)

	if res.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return link, res.Error
}

func (r *UrlRepo) ByUserID(id string, limit int, offset int) ([]model.Url, error) {
	var urls []model.Url

	return urls, r.db.Where("user_id = ?", id).Limit(limit).Offset(offset).Find(&urls).Error
}
