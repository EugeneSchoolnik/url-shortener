package repo

import (
	"log"
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

func (r *ClickStatRepo) ByUrlID(urlID, userID string) ([]DailyCount, error) {
	var results []DailyCount

	err := r.db.Model(&model.ClickStat{}).
		Select("date_trunc('day', click_stats.created_at) AS day, COUNT(*) AS count").
		Joins("JOIN urls ON urls.id = click_stats.url_id").
		Where("click_stats.url_id = ? AND urls.user_id = ?", urlID, userID).
		Group("day").
		Order("day").
		Scan(&results).Error

	return results, err
}

func (r *ClickStatRepo) CleanupStaleRecords() error {
	result := r.db.Where("created_at < now() - interval '30 days'").Delete(&model.ClickStat{})
	log.Printf("Deleted %d old events\n", result.RowsAffected)

	return result.Error
}

// func fillEmptyDays(stats []DailyCount) []DailyCount {
// 	startDate := time.Now().AddDate(0, 0, -29).Truncate(24 * time.Hour).UTC()
// 	endDate := time.Now().Truncate(24 * time.Hour).UTC()

// 	dateMap := make(map[time.Time]int64)
// 	for _, r := range stats {
// 		dateMap[r.Day] = r.Count
// 	}

// 	filledStats := make([]DailyCount, 0, 30)
// 	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
// 		count := dateMap[d]
// 		filledStats = append(filledStats, DailyCount{Day: d, Count: count})
// 	}

// 	return filledStats
// }
