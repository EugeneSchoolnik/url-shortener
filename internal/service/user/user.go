package user

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

//go:generate mockery --name=UserRepo
type UserRepo interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id string) error
	ByEmail(email string) (*model.User, error)
	ById(id string) (*model.User, error)
	ContextByEmail(email string) (*model.User, error)
	ContextById(id string) (*model.User, error)
}

type UserService struct {
	repo UserRepo
	log  *slog.Logger
}

func New(repo UserRepo, log *slog.Logger) *UserService {
	return &UserService{repo, log}
}

// TODO: maybe add to config
const idAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idSize = 12

var idGenerator = nanoid.New(idAlphabet, idSize)

func (s *UserService) Create(user *model.User) error {
	const op = "service.user.Create"
	log := s.log.With(slog.String("op", op))

GenerateID:
	id, err := idGenerator.ID()
	if err != nil {
		log.Error("failed to generate id", sl.Err(err))
		return service.ErrInternalError
	}

	user.ID = id

	err = s.repo.Create(user)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		if pgErr := pg.ParsePGError(err); pgErr != nil && pgErr.Code == "23505" { // // 23505 = unique_violation
			fmt.Println("pgErr.ConstraintName", pgErr.ConstraintName)
			switch pgErr.ConstraintName {
			case "users_pkey":
				goto GenerateID
			case "idx_users_email":
				return ErrEmailTaken
			}
		}
		return service.ErrInternalError
	}

	log.Info("user successfully created")
	return nil
}

func (s *UserService) Update(id string, userDto *dto.UpdateUser) error {
	const op = "service.user.Update"
	log := s.log.With(slog.String("op", op))

	if id == "" {
		log.Info("id validation failed")
		return fmt.Errorf("%w: %s", service.ErrValidation, "id is required")
	}

	err := service.Validate.Struct(userDto)
	if err != nil {
		log.Info("dto validation failed", sl.Err(err))
		return service.PrettyValidationError(err.(validator.ValidationErrors))
	}

	user := userDto.Model()
	user.ID = id

	err = s.repo.Update(user)
	if err != nil {
		log.Error("failed to update user", sl.Err(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrNotFound
		}
		return service.ErrInternalError
	}

	log.Info("user successfully updated")
	return nil
}

func (s *UserService) Delete(id string) error {
	const op = "service.user.Delete"
	log := s.log.With(slog.String("op", op))

	if id == "" {
		log.Info("id validation failed")
		return fmt.Errorf("%w: %s", service.ErrValidation, "id is required")
	}

	if err := s.repo.Delete(id); err != nil {
		log.Error("failed to delete user", sl.Err(err))
		return service.ErrInternalError
	}

	log.Info("user succesfully deleted")
	return nil
}

func (s *UserService) ByEmail(email string, withContext ...bool) (*model.User, error) {
	const op = "service.user.ByEmail"
	log := s.log.With(slog.String("op", op))

	if err := service.Validate.Var(email, "email,required"); err != nil {
		log.Info("email validation failed")
		return nil, fmt.Errorf("%w: %s", service.ErrValidation, "invalid email")
	}

	var user *model.User
	var err error
	if len(withContext) > 0 {
		user, err = s.repo.ContextByEmail(email)
	} else {
		user, err = s.repo.ByEmail(email)
	}

	if err != nil {
		log.Error("failed to get user by email", sl.Err(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, service.ErrInternalError
	}

	log.Info("user found by email successfully")
	return user, nil
}

func (s *UserService) ById(id string, withContext ...bool) (*model.User, error) {
	const op = "service.user.ById"
	log := s.log.With(slog.String("op", op))

	if id == "" {
		log.Info("id validation failed")
		return nil, fmt.Errorf("%w: %s", service.ErrValidation, "id is required")
	}

	var user *model.User
	var err error
	if len(withContext) > 0 {
		user, err = s.repo.ContextById(id)
	} else {
		user, err = s.repo.ById(id)
	}

	if err != nil {
		log.Error("failed to get user by id", sl.Err(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, service.ErrInternalError
	}

	log.Info("user found by id successfully")
	return user, nil
}
