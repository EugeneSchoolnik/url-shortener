package user_test

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"
	"url-shortener/internal/service/user"
	"url-shortener/internal/service/user/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func makeDuplicateKeyError(constraint string) error {
	return &pgconn.PgError{
		Code:           "23505", // unique_violation
		ConstraintName: constraint,
	}
}

func TestUserService_Create(t *testing.T) {
	const idSize = 12

	tests := []struct {
		name       string
		argUserDto *model.User
		mockSetup  func(r *mocks.UserRepo)
		wantErr    error
	}{
		{
			name:       "success",
			argUserDto: &model.User{Email: "example@email.com", Password: "12345678"},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Create", mock.Anything).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:       "success with duplicate id",
			argUserDto: &model.User{Email: "example@email.com", Password: "12345678"},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Create", mock.Anything).Return(makeDuplicateKeyError("users_pkey")).Once()
				r.On("Create", mock.Anything).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:       "email taken",
			argUserDto: &model.User{Email: "example@email.com", Password: "12345678"},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Create", mock.Anything).Return(makeDuplicateKeyError("idx_users_email")).Once()
			},
			wantErr: user.ErrEmailTaken,
		},
		{
			name:       "unexpected error",
			argUserDto: &model.User{Email: "example@email.com", Password: "12345678"},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Create", mock.Anything).Return(errors.New("unexpected error")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := user.New(repo, slog.Default())

			err := s.Create(tt.argUserDto)
			if tt.wantErr != err && !errors.Is(err, tt.wantErr) {
				t.Errorf("UserService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Check ID
			if err == nil {
				if len(tt.argUserDto.ID) != idSize {
					t.Errorf("tt.argUserDto.ID size = %v, want %v", len(tt.argUserDto.ID), idSize)
				}
			}
			// if !reflect.DeepEqual(got, want) {
			// 	t.Errorf("UserService.Create() = %v, want %v", got, want)
			// }
		})
	}
}

func TestUserService_Update(t *testing.T) {
	type args struct {
		id      string
		userDto *dto.UpdateUser
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(r *mocks.UserRepo)
		wantErr   error
	}{
		{
			name: "success",
			args: args{
				id:      "1234",
				userDto: &dto.UpdateUser{Email: "example@email.com"},
			},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Update", mock.Anything).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "empty id",
			args: args{
				id:      "",
				userDto: &dto.UpdateUser{Email: "example@email.com"},
			},
			wantErr: service.ErrValidation,
		},
		{
			name: "invalid email",
			args: args{
				id:      "1234",
				userDto: &dto.UpdateUser{Email: "invalid.email.com"},
			},
			wantErr: service.ErrValidation,
		},
		{
			name: "empty dto field",
			args: args{
				id:      "1234",
				userDto: &dto.UpdateUser{Email: ""},
			},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Update", mock.Anything).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				id:      "404",
				userDto: &dto.UpdateUser{Email: "example@email.com"},
			},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Update", mock.Anything).Return(gorm.ErrRecordNotFound).Once()
			},
			wantErr: service.ErrNotFound,
		},
		{
			name: "unxpected error",
			args: args{
				id:      "1234",
				userDto: &dto.UpdateUser{Email: "example@email.com"},
			},
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Update", mock.Anything).Return(errors.New("unexprect error")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := user.New(repo, slog.Default())

			if err := s.Update(tt.args.id, tt.args.userDto); tt.wantErr != err && !errors.Is(err, tt.wantErr) {
				t.Errorf("UserService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(r *mocks.UserRepo)
		wantErr   error
	}{
		{
			name: "success",
			id:   "1234",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Delete", mock.IsType("")).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: service.ErrValidation,
		},
		{
			name: "unxpected error",
			id:   "1234",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("Delete", mock.IsType("")).Return(errors.New("unxpected error")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := user.New(repo, slog.Default())

			if err := s.Delete(tt.id); err != tt.wantErr && !errors.Is(err, tt.wantErr) {
				t.Errorf("UserService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_ByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockSetup func(r *mocks.UserRepo)
		want      *model.User
		wantErr   error
	}{
		{
			name:  "success",
			email: "example@email.com",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("ByEmail", mock.IsType("")).Return(&model.User{ID: "1234", Email: "example@email.com", Password: "12345678"}, nil)
			},
			want:    &model.User{ID: "1234", Email: "example@email.com", Password: "12345678"},
			wantErr: nil,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: service.ErrValidation,
		},
		{
			name:    "invalid email",
			email:   "",
			wantErr: service.ErrValidation,
		},
		{
			name:  "not found",
			email: "example@email.com",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("ByEmail", mock.IsType("")).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: service.ErrNotFound,
		},
		{
			name:  "unxpected error",
			email: "example@email.com",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("ByEmail", mock.IsType("")).Return(nil, errors.New("unxpected error"))
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			s := user.New(repo, slog.Default())

			got, err := s.ByEmail(tt.email)
			if (err != tt.wantErr) && !errors.Is(err, tt.wantErr) {
				t.Errorf("UserService.ByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.ByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_ById(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(r *mocks.UserRepo)
		want      *model.User
		wantErr   error
	}{
		{
			name: "success",
			id:   "1234",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("ById", mock.IsType("")).Return(&model.User{ID: "1234", Email: "example@email.com", Password: "12345678"}, nil)
			},
			want:    &model.User{ID: "1234", Email: "example@email.com", Password: "12345678"},
			wantErr: nil,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: service.ErrValidation,
		},
		{
			name: "not found",
			id:   "404",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("ById", mock.IsType("")).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: service.ErrNotFound,
		},
		{
			name: "unxpected error",
			id:   "1234",
			mockSetup: func(r *mocks.UserRepo) {
				r.On("ById", mock.IsType("")).Return(nil, errors.New("unxpected error"))
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			s := user.New(repo, slog.Default())

			got, err := s.ById(tt.id)
			if (err != tt.wantErr) && !errors.Is(err, tt.wantErr) {
				t.Errorf("UserService.ById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.ById() = %v, want %v", got, tt.want)
			}
		})
	}
}
