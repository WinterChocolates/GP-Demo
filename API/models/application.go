package models

import "time"

type Application struct {
	ApplicationID uint      `gorm:"primaryKey;column:application_id" json:"application_id"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	JobID         uint      `gorm:"not null" json:"job_id"`
	ApplyDate     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"apply_date"`
	Status        string    `gorm:"type:ENUM('pending','interviewed','hired','rejected');default:'pending'" json:"status"`

	User User `gorm:"foreignKey:UserID" json:"-"`
	Job  Job  `gorm:"foreignKey:JobID" json:"-"`
}
