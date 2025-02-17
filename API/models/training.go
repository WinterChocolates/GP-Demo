package models

import "time"

type Training struct {
	TrainingID  uint      `gorm:"primaryKey;column:training_id" json:"training_id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	Location    string    `gorm:"size:100" json:"location,omitempty"`
	Capacity    uint      `json:"capacity,omitempty"`

	Records []TrainingRecord `gorm:"foreignKey:TrainingID" json:"-"`
}
