package models

import "time"

type Job struct {
	JobID          uint      `gorm:"primaryKey;column:job_id" json:"job_id"`
	Title          string    `gorm:"size:100;not null" json:"title" binding:"required"`
	Description    string    `gorm:"type:text;not null" json:"description"`
	Requirements   string    `gorm:"type:text;not null" json:"requirements"`
	SalaryRange    string    `gorm:"size:50" json:"salary_range,omitempty"`
	PostDate       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"post_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Status         string    `gorm:"type:ENUM('open','closed');default:'open'" json:"status"`

	Applications []Application `gorm:"foreignKey:JobID" json:"-"`
}
