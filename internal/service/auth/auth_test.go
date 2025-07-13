package auth_test

import (
	"errors"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"
	"url-shortener/internal/service/auth"
	"url-shortener/internal/service/auth/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		userDto   *dto.CreateUser
		mockSetup func(r *mocks.UserService)
		wantErr   error
	}{
		{
			name:    "success",
			userDto: &dto.CreateUser{Email: "example@email.com", Password: "12345678"},
			mockSetup: func(r *mocks.UserService) {
				r.On("Create", mock.AnythingOfType("*model.User")).
					Run(func(args mock.Arguments) {
						u := args.Get(0).(*model.User)
						u.ID = "generated-id"
					}).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:    "invalid email",
			userDto: &dto.CreateUser{Email: "invalid.email.com", Password: "12345678"},
			wantErr: service.ErrValidation,
		},
		{
			name:    "too short password",
			userDto: &dto.CreateUser{Email: "example@email.com", Password: "12345"},
			wantErr: service.ErrValidation,
		},
		{
			name:    "too long password",
			userDto: &dto.CreateUser{Email: "example@email.com", Password: strings.Repeat("0", 256)},
			wantErr: service.ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := mocks.NewUserService(t)
			jwt := auth.NewJWTService("secret", time.Hour)

			if tt.mockSetup != nil {
				tt.mockSetup(userService)
			}

			s := auth.New(userService, jwt, slog.Default())

			var want *model.User
			if tt.wantErr == nil {
				want = tt.userDto.Model()
			}

			got, token, err := s.Register(tt.userDto)
			if (err != tt.wantErr) && !errors.Is(err, tt.wantErr) {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				// Check ID
				if got.ID == "" {
					t.Error("AuthService.Register() got.ID = ''")
				}
				want.ID = got.ID // we've already checked it

				// Check password
				if !s.ComparePassword(tt.userDto.Password, got.Password) {
					t.Error("AuthService.Register() got.Password is invalid")
				}
				want.Password = got.Password // we've already checked it

				// Check token
				if token == "" {
					t.Error("AuthService.Register() token is empty")
				}

				if !reflect.DeepEqual(got, want) {
					t.Errorf("AuthService.Register() got = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(r *mocks.UserService)
		wantErr   error
	}{
		{
			name: "success",
			args: args{email: "example@email.com", password: "12345678"},
			mockSetup: func(r *mocks.UserService) {
				r.On("ByEmail", mock.IsType(""), mock.IsType(true)).Return(&model.User{
					ID: "1234", Email: "example@email.com", Password: "$2a$10$sDH5VXdDxPPS0w8VctctUur9n1YPFhNyfeSD.EcfR7OpEkIzDBai6",
				}, nil).Once()
			},
		},
		{
			name:    "invalid password",
			args:    args{email: "example@email.com", password: "12345"},
			wantErr: service.ErrValidation,
		},
		{
			name: "wrong password",
			args: args{email: "example@email.com", password: "12345678"},
			mockSetup: func(r *mocks.UserService) {
				r.On("ByEmail", mock.IsType(""), mock.IsType(true)).Return(&model.User{
					ID: "1234", Email: "example@email.com", Password: "$2a$10$another/hash",
				}, nil).Once()
			},
			wantErr: service.ErrInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := mocks.NewUserService(t)
			jwt := auth.NewJWTService("secret", time.Hour)

			if tt.mockSetup != nil {
				tt.mockSetup(userService)
			}

			s := auth.New(userService, jwt, slog.Default())
			got, token, err := s.Login(tt.args.email, tt.args.password)
			if (err != tt.wantErr) && !errors.Is(err, tt.wantErr) {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				assert.Nil(t, got)
				return
			}

			assert.NotEmpty(t, token) // check token

			assert.Equal(t, got.Email, tt.args.email)
			assert.NotEmpty(t, got.ID)
			assert.Equal(t, len(got.Password), 60)
		})
	}
}

func TestAuthService_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "success",
			password: "12345678",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auth.AuthService{}
			got, err := s.HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			const wantLen = 60
			if len(got) != wantLen {
				t.Errorf("AuthService.HashPassword() = %v, wantLen %v", len(got), wantLen)
			}
		})
	}
}

func TestAuthService_ComparePassword(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid password",
			args: args{
				password: "12345678",
				hash:     "$2a$10$sDH5VXdDxPPS0w8VctctUur9n1YPFhNyfeSD.EcfR7OpEkIzDBai6",
			},
			want: true,
		},
		{
			name: "invalid password",
			args: args{
				password: "12345678",
				hash:     "invalidhash",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auth.AuthService{}
			if got := s.ComparePassword(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("AuthService.CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
