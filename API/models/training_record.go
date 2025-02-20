package models

import (
	"gorm.io/gorm"
)

type TrainingRecord struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	TrainingID uint   `gorm:"not null"`
	Status     string `gorm:"type:ENUM('registered','completed','canceled');default:'registered'"`
	Score      uint8

	User     User     `gorm:"foreignKey:UserID"`
	Training Training `gorm:"foreignKey:TrainingID"`
}
