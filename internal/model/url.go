package model

type Url struct {
	ID        string `gorm:"primaryKey;type:varchar(16)"`
	Link      string `gorm:"type:varchar(255);not null"`
	TotalHits int64  `gorm:"type:bigint;not null;default:0"`
	UserID    string `gorm:"type:varchar(16);not null;index"`
}
