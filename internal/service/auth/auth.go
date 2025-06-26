package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

//go:generate mockery --name=UserService
type UserService interface {
	Create(user *model.User) error
	Update(id string, userDto *dto.UpdateUser) error
	Delete(id string) error
	ByEmail(email string, withContext ...bool) (*model.User, error)
	ById(id string, withContext ...bool) (*model.User, error)
}

type AuthService struct {
	userService UserService
	jwtService  *JWTService
	log         *slog.Logger
}

func New(userService UserService, jwtService *JWTService, log *slog.Logger) *AuthService {
	return &AuthService{userService, jwtService, log}
}

func (s *AuthService) Register(userDto *dto.CreateUser) (*model.User, string, error) {
	const op = "service.auth.Register"
	log := s.log.With(slog.String("op", op))

	err := service.Validate.Struct(userDto)
	if err != nil {
		log.Info("validation failed", sl.Err(err))
		return nil, "", service.PrettyValidationError(err.(validator.ValidationErrors))
	}

	passwordHash, err := s.HashPassword(userDto.Password)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return nil, "", service.ErrInternalError
	}

	user := userDto.Model()
	user.Password = passwordHash

	err = s.userService.Create(user)
	if err != nil {
		// no need for logs
		return nil, "", err
	}

	token, err := s.jwtService.Generate(user.ID)
	if err != nil {
		log.Error("failed to generate jwt token", sl.Err(err))
		return nil, "", service.ErrInternalError
	}

	log.Info("user successfully registered")
	return user, token, nil
}

func (s *AuthService) Login(email string, password string) (*model.User, string, error) {
	const op = "service.auth.Login"
	log := s.log.With(slog.String("op", op))

	// check password
	// Get the "Password" field's validation tag
	field, ok := reflect.TypeOf(dto.CreateUser{}).FieldByName("Password")
	if !ok {
		return nil, "", service.ErrInternalError
	}
	tag := field.Tag.Get("validate")
	if err := service.Validate.Var(password, tag); err != nil {
		log.Info("failed to validate password", sl.Err(err))
		return nil, "", fmt.Errorf("%w: %s", service.ErrValidation, "invalid password")
	}

	// email validation is inside this
	user, err := s.userService.ByEmail(email, true)
	if err != nil {
		// no need for logs
		return nil, "", err
	}

	if !s.ComparePassword(password, user.Password) {
		log.Info("wrong password")
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.jwtService.Generate(user.ID)
	if err != nil {
		log.Error("failed to generate jwt token", sl.Err(err))
		return nil, "", service.ErrInternalError
	}

	log.Info("user successfully logged in")
	return user, token, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost = 10
	return string(bytes), err
}

func (s *AuthService) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
