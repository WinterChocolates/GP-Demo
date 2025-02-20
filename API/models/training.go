package models

import (
	"gorm.io/gorm"
	"time"
)

type Training struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null"`
	Description string
	StartTime   time.Time `gorm:"not null"`
	EndTime     time.Time `gorm:"not null"`
	Location    string    `gorm:"type:varchar(100)"`
	Capacity    uint

	Records []TrainingRecord
}
