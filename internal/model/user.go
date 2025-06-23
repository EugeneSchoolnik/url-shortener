package model

type User struct {
	ID       string `gorm:"primaryKey;type:varchar(16)"`
	Email    string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(60);not null"`
}
