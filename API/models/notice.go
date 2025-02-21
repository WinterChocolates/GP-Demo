package models

import (
	"time"

	"gorm.io/gorm"
)

type Notice struct {
	gorm.Model
	Title      string `gorm:"type:varchar(200);not null"`
	Content    string `gorm:"type:text;not null"`
	ExpireTime *time.Time
	TargetType string `gorm:"type:ENUM('all','department');default:'all'"`
	Department string `gorm:"type:varchar(50)"`
}
