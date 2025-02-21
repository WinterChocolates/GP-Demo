package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title          string `gorm:"type:varchar(100);not null"`
	Description    string `gorm:"type:text;not null"`
	Requirements   string `gorm:"type:text;not null"`
	SalaryRange    string `gorm:"type:varchar(50)"`
	ExpirationDate *time.Time
	Status         string `gorm:"type:ENUM('open','closed');default:'open'"`

	Applications []Application
}
