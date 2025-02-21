package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	UserID      uint      `gorm:"index;not null"`
	ClockIn     time.Time `gorm:"not null"`
	ClockOut    *time.Time
	Status      string    `gorm:"type:ENUM('normal','late','early_leave');default:'normal'"`
	DateClockIn time.Time `gorm:"->;type:date;generated"`

	User User `gorm:"foreignKey:UserID"`
}
