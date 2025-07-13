package url

import (
	"errors"
	"fmt"
	"log/slog"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/pg"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"
	"url-shortener/internal/util/nanoid"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

//go:generate mockery --name=UrlRepo
type UrlRepo interface {
	Create(url *model.Url) error
	ByID(id string) (*model.Url, error)
	LinkByID(id string) (string, error)
	ByUserID(id string, limit int, offset int) ([]model.Url, error)
	Delete(id string, userID string) error
}

type UrlService struct {
	repo UrlRepo
	log  *slog.Logger
}

func New(repo UrlRepo, log *slog.Logger) *UrlService {
	return &UrlService{repo, log}
}

const idAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idSize = 8

var idGenerator = nanoid.New(idAlphabet, idSize)

func (s *UrlService) Create(urlDto *dto.CreateUrl, userID string) (*model.Url, error) {
	log := s.log.With(slog.String("op", "service.url.Create"))

	if err := service.Validate.Struct(urlDto); err != nil {
		log.Info("validation failed", sl.Err(err))
		return nil, service.PrettyValidationError(err.(validator.ValidationErrors))
	}

	url := urlDto.Model(userID)
	autogeneration := url.ID == ""

GenerateID:
	if autogeneration {
		id, err := idGenerator.ID()
		if err != nil {
			log.Error("failed to generate id", sl.Err(err))
			return nil, service.ErrInternalError
		}
		url.ID = id
	}

	if err := s.repo.Create(url); err != nil {
		log.Error("failed to create url", sl.Err(err))
		if pgErr := pg.ParsePGError(err); pgErr != nil {
			if pgErr.Code == "23503" { // 23503 = foreign_key_violation
				return nil, service.ErrRelatedResourceNotFound
			}
			if pgErr.Code == "23505" { // 23505 = unique_violation
				if autogeneration {
					goto GenerateID
				} else {
					return nil, service.ErrAliasTaken
				}
			}
		}
		return nil, service.ErrInternalError
	}

	log.Info("url successfully created")
	return url, nil
}

func (s *UrlService) ByID(id string) (*model.Url, error) {
	log := s.log.With(slog.String("op", "service.url.ByID"))

	if id == "" {
		log.Info("id is empty")
		return nil, fmt.Errorf("%w%s", service.ErrValidation, "id is a required")
	}

	url, err := s.repo.ByID(id)
	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrUrlNotFound
		}
		return nil, service.ErrInternalError
	}
	log.Info("got url by id successfully")
	return url, nil
}

func (s *UrlService) RedirectLinkByID(id string) (string, error) {
	log := s.log.With(slog.String("op", "service.url.RedirectLinkByID"))

	if id == "" {
		log.Info("id is empty")
		return "", fmt.Errorf("%w%s", service.ErrValidation, "id is a required")
	}

	link, err := s.repo.LinkByID(id)
	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", service.ErrUrlNotFound
		}
		return "", service.ErrInternalError
	}
	log.Info("got link by id successfully")
	return link, nil
}

func (s *UrlService) ByUserID(id string, limit int, offset int) ([]model.Url, error) {
	log := s.log.With(slog.String("op", "service.url.ByID"))

	if id == "" {
		log.Info("id is empty")
		return nil, fmt.Errorf("%w%s", service.ErrValidation, "id is a required")
	}

	urls, err := s.repo.ByUserID(id, limit, offset)
	if err != nil {
		log.Error("failed to get urls", sl.Err(err))
		return nil, service.ErrInternalError
	}
	log.Info("got urls by user id successfully")
	return urls, nil
}

func (s *UrlService) Delete(id, userID string) error {
	log := s.log.With(slog.String("op", "service.url.Delete"))

	if id == "" || userID == "" {
		log.Info("id is empty")
		return fmt.Errorf("%w%s", service.ErrValidation, "id is a required")
	}

	if err := s.repo.Delete(id, userID); err != nil {
		log.Error("failed to delete url", sl.Err(err))
		return service.ErrInternalError
	}

	log.Info("url successfully deleted")
	return nil
}
