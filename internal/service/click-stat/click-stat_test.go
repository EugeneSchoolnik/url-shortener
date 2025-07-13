package clickstat_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/service"
	clickstat "url-shortener/internal/service/click-stat"
	"url-shortener/internal/service/click-stat/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClickStatService_Record(t *testing.T) {
	tests := []struct {
		name      string
		urlID     string
		mockSetup func(r *mocks.ClickStatRepo)
		wantErr   error
	}{
		{
			name:  "success",
			urlID: "1234",
			mockSetup: func(r *mocks.ClickStatRepo) {
				r.On("Create", mock.Anything).Return(nil).Once()
			},
		},
		{
			name:  "url id that doesn't exist",
			urlID: "notfound",
			mockSetup: func(r *mocks.ClickStatRepo) {
				r.On("Create", mock.Anything).Return(&pgconn.PgError{Code: "23503"}).Once()
			},
			wantErr: service.ErrRelatedResourceNotFound,
		},
		{
			name:  "unexpected error",
			urlID: "1234",
			mockSetup: func(r *mocks.ClickStatRepo) {
				r.On("Create", mock.Anything).Return(errors.New("unexpected")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewClickStatRepo(t)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			s := clickstat.New(repo, slog.Default())

			err := s.Record(tt.urlID)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestClickStatService_Stats(t *testing.T) {
	today := time.Now().Truncate(24 * time.Hour)
	stats := []repo.DailyCount{
		{Day: today.AddDate(0, 0, 0), Count: 23},
		{Day: today.AddDate(0, 0, -1), Count: 58},
		{Day: today.AddDate(0, 0, -2), Count: 11},
		{Day: today.AddDate(0, 0, -3), Count: 9},
	}

	tests := []struct {
		name      string
		urlID     string
		userID    string
		mockSetup func(r *mocks.ClickStatRepo)
		want      []repo.DailyCount
		wantErr   error
	}{
		{
			name:  "success",
			urlID: "1234",
			mockSetup: func(r *mocks.ClickStatRepo) {
				r.On("ByUrlID", mock.Anything, mock.Anything).Return(stats, nil).Once()
			},
			want: stats,
		},
		{
			name:  "unexpected",
			urlID: "1234",
			mockSetup: func(r *mocks.ClickStatRepo) {
				r.On("ByUrlID", mock.Anything, mock.Anything).Return(nil, errors.New("unexpected")).Once()
			},
			wantErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewClickStatRepo(t)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			s := clickstat.New(repo, slog.Default())

			got, err := s.Stats(tt.urlID, tt.userID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
