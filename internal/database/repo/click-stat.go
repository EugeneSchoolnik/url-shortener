package repo

import (
	"time"
	"url-shortener/internal/model"

	"gorm.io/gorm"
)

type DailyCount struct {
	Day   time.Time
	Count int64
}

type ClickStatRepo struct {
	db *gorm.DB
}

func NewClickStatRepo(db *gorm.DB) *ClickStatRepo {
	return &ClickStatRepo{db}
}

func (r *ClickStatRepo) Create(ClickStat *model.ClickStat) error {
	return r.db.Create(ClickStat).Error
}

func (r *ClickStatRepo) ByUrlID(id string) ([]DailyCount, error) {
	var results []DailyCount

	err := r.db.Model(&model.ClickStat{}).
		Select("date_trunc('day', created_at) AS day, COUNT(*) AS count").
		Where("url_id = ?", id).
		Group("day").
		Order("day").
		Scan(&results).Error

	return fillEmptyDays(results), err
}

func fillEmptyDays(stats []DailyCount) []DailyCount {
	startDate := time.Now().AddDate(0, 0, -29).Truncate(24 * time.Hour).UTC()
	endDate := time.Now().Truncate(24 * time.Hour).UTC()

	dateMap := make(map[time.Time]int64)
	for _, r := range stats {
		dateMap[r.Day] = r.Count
	}

	filledStats := make([]DailyCount, 0, 30)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		count := dateMap[d]
		filledStats = append(filledStats, DailyCount{Day: d, Count: count})
	}

	return filledStats
}
