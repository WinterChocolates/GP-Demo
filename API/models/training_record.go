package models

import "time"

type TrainingRecord struct {
	RecordID   uint      `gorm:"primaryKey;column:record_id" json:"record_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	TrainingID uint      `gorm:"not null" json:"training_id"`
	EnrollTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"enroll_time"`
	Status     string    `gorm:"type:ENUM('registered','completed','canceled');default:'registered'" json:"status"`
	Score      uint8     `json:"score,omitempty"`

	User     User     `gorm:"foreignKey:UserID" json:"-"`
	Training Training `gorm:"foreignKey:TrainingID" json:"-"`
}
