package model

import "time"

type ClickStat struct {
	UrlID     string    `gorm:"type:varchar(16);not null;index:idx_url_created"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;index:idx_url_created"`
}
