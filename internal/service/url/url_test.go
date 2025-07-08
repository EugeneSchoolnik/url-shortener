package url_test

import (
	"errors"
	"log/slog"
	"strings"
	"testing"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"
	"url-shortener/internal/service/url"
	"url-shortener/internal/service/url/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUrlService_Create(t *testing.T) {
	const idSize = 8

	tests := []struct {
		name      string
		urlDto    *dto.CreateUrl
		userID    string
		mockSetup func(r *mocks.UrlRepo)
		wantErr   error
	}{
		{
			name:      "success with alias",
			urlDto:    &dto.CreateUrl{Alias: "g", Link: "https://google.com"},
			userID:    "1234",
			mockSetup: func(r *mocks.UrlRepo) { r.On("Create", mock.Anything).Return(nil).Once() },
		},
		{
			name:      "success without alias",
			urlDto:    &dto.CreateUrl{Link: "https://google.com"},
			userID:    "1234",
			mockSetup: func(r *mocks.UrlRepo) { r.On("Create", mock.Anything).Return(nil).Once() },
		},
		{
			name:   "success with regenerated id",
			urlDto: &dto.CreateUrl{Link: "https://google.com"},
			userID: "1234",
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("Create", mock.Anything).
					Return(&pgconn.PgError{Code: "23505"}). // 23505 = unique_violation
					Once()
				r.On("Create", mock.Anything).Return(nil).Once()
			},
		},
		{
			name:    "too long alias",
			urlDto:  &dto.CreateUrl{Alias: strings.Repeat("0", 17), Link: "https://google.com"},
			userID:  "1234",
			wantErr: service.ErrValidation,
		},
		{
			name:    "special characters in alias",
			urlDto:  &dto.CreateUrl{Alias: "😂", Link: "https://google.com"},
			userID:  "1234",
			wantErr: service.ErrValidation,
		},
		{
			name:    "without url",
			urlDto:  &dto.CreateUrl{Link: ""},
			userID:  "1234",
			wantErr: service.ErrValidation,
		},
		{
			name:    "invalid url",
			urlDto:  &dto.CreateUrl{Link: "noturl"},
			userID:  "1234",
			wantErr: service.ErrValidation,
		},
		{
			name:   "user id that doesn't exist",
			urlDto: &dto.CreateUrl{Link: "https://google.com"},
			userID: "notfound",
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("Create", mock.Anything).
					Return(&pgconn.PgError{Code: "23503"}). // 23503 = foreign_key_violation
					Once()
			},
			wantErr: service.ErrRelatedResourceNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UrlRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := url.New(repo, slog.Default())

			got, err := s.Create(tt.urlDto, tt.userID)
			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.Equal(t, tt.urlDto.Link, got.Link)
				assert.Equal(t, tt.userID, got.UserID)
				if tt.urlDto.Alias == "" {
					assert.Len(t, got.ID, idSize)
				}
			}
		})
	}
}

func TestUrlService_ByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(r *mocks.UrlRepo)
		wantErr   error
	}{
		{
			name: "success",
			id:   "1234",
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("ByID", mock.Anything).Return(&model.Url{ID: "1234", Link: "https://google.com"}, nil).Once()
			},
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: service.ErrValidation,
		},
		{
			name: "not found",
			id:   "1234",
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("ByID", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			wantErr: service.ErrUrlNotFound,
		},
		{
			name: "unxpected error",
			id:   "1234",
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("ByID", mock.Anything).Return(nil, errors.New("unexpected")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UrlRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := url.New(repo, slog.Default())

			got, err := s.ByID(tt.id)
			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.NotEmpty(t, got.Link)
			}
		})
	}
}

func TestUrlService_ByUserID(t *testing.T) {
	type args struct {
		id     string
		limit  int
		offset int
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func(r *mocks.UrlRepo)
		wantErr   error
	}{
		{
			name: "success",
			args: args{id: "1234", limit: 5, offset: 0},
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("ByUserID", mock.Anything, mock.Anything, mock.Anything).Return([]model.Url{{ID: "1234", Link: "https://google.com"}}, nil).Once()
			},
		},
		{
			name:    "empty id",
			args:    args{id: "", limit: 5, offset: 0},
			wantErr: service.ErrValidation,
		},
		{
			name: "unxpected error",
			args: args{id: "1234", limit: 5, offset: 0},
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("ByUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("unexpected")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UrlRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := url.New(repo, slog.Default())

			got, err := s.ByUserID(tt.args.id, tt.args.limit, tt.args.offset)
			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.IsType(t, []model.Url{}, got)
			}
		})
	}
}

func TestUrlService_Delete(t *testing.T) {
	type args struct {
		id     string
		userID string
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func(r *mocks.UrlRepo)
		wantErr   error
	}{
		{
			name:      "success",
			args:      args{id: "1234", userID: "1234"},
			mockSetup: func(r *mocks.UrlRepo) { r.On("Delete", mock.Anything, mock.Anything).Return(nil).Once() },
		},
		{
			name:    "empty id",
			args:    args{id: "", userID: "1234"},
			wantErr: service.ErrValidation,
		},
		{
			name:    "empty user id",
			args:    args{id: "1234", userID: ""},
			wantErr: service.ErrValidation,
		},
		{
			name: "unexpected error",
			args: args{id: "1234", userID: "1234"},
			mockSetup: func(r *mocks.UrlRepo) {
				r.On("Delete", mock.Anything, mock.Anything).Return(errors.New("unexpected")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UrlRepo{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			s := url.New(repo, slog.Default())

			err := s.Delete(tt.args.id, tt.args.userID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
