package models

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	UserID uint   `gorm:"index;not null"`
	JobID  uint   `gorm:"index;not null"`
	Status string `gorm:"type:ENUM('pending','interviewed','hired','rejected');default:'pending'"`

	User User `gorm:"foreignKey:UserID"`
	Job  Job  `gorm:"foreignKey:JobID"`
}
