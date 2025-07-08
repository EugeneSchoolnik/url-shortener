package url

import (
	"log/slog"
	"url-shortener/internal/model"
)

type UrlRepo interface {
	Create(url *model.Url) error
	ByID(id string) (*model.Url, error)
	ByUserID(id string, limit int, offset int) ([]model.Url, error)
	Delete(id string, userID string) error
}

type UrlService struct {
	repo UrlRepo
	log  *slog.Logger
}
