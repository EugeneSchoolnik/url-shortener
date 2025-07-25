package clickstat

import (
	"log/slog"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/pg"
	"url-shortener/internal/model"
	"url-shortener/internal/service"

	"github.com/robfig/cron/v3"
)

//go:generate mockery --name=ClickStatRepo
type ClickStatRepo interface {
	Create(ClickStat *model.ClickStat) error
	ByUrlID(urlID string, userID string) ([]repo.DailyCount, error)
	CleanupStaleRecords() error
}

type ClickStatService struct {
	repo ClickStatRepo
	log  *slog.Logger
}

func New(repo ClickStatRepo, log *slog.Logger) *ClickStatService {
	return &ClickStatService{repo, log}
}

func (s *ClickStatService) Record(urlID string) error {
	log := s.log.With(slog.String("op", "service.clickstat.Record"))

	if err := s.repo.Create(&model.ClickStat{UrlID: urlID}); err != nil {
		log.Error("failed to record click", sl.Err(err))
		if pgErr := pg.ParsePGError(err); pgErr != nil && pgErr.Code == "23503" { // 23503 = foreign_key_violation
			return service.ErrRelatedResourceNotFound
		}
		return service.ErrInternalError
	}

	// log.Info("click successfully recorded")
	return nil
}

func (s *ClickStatService) Stats(urlID, userID string) ([]repo.DailyCount, error) {
	log := s.log.With(slog.String("op", "service.clickstat.Stats"))

	stats, err := s.repo.ByUrlID(urlID, userID)
	if err != nil {
		log.Error("failed to get stats", sl.Err(err))
		return nil, service.ErrInternalError
	}
	if len(stats) == 0 {
		log.Info("statistics not found")
		return nil, service.ErrUrlStatsNotFound
	}

	log.Info("statistics successfully received")
	return stats, nil
}
func (s *ClickStatService) CleanupStaleRecords() (*cron.Cron, error) {
	c := cron.New()

	// Run daily at 02:00 AM
	_, err := c.AddFunc("0 2 * * *", func() {
		s.repo.CleanupStaleRecords()
	})
	if err != nil {
		return nil, err
	}

	c.Start()

	return c, nil
}
