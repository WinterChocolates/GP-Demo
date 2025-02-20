package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username       string `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash   string `gorm:"type:char(60);not null"`
	UserType       string `gorm:"type:ENUM('admin','employee','applicant');not null"`
	Department     string `gorm:"type:varchar(50)"`
	Position       string `gorm:"type:varchar(50)"`
	HireDate       *time.Time
	SalaryBase     float64 `gorm:"type:decimal(10,2)"`
	Education      string  `gorm:"type:text"`
	WorkExperience string  `gorm:"type:text"`
	Skills         string  `gorm:"type:text"`
	ResumePath     string  `gorm:"type:varchar(255)"`
	IsActive       bool    `gorm:"default:true"`

	Applications    []Application
	Attendances     []Attendance
	Salaries        []Salary
	TrainingRecords []TrainingRecord
	Roles           []Role `gorm:"many2many:user_roles;"`
}
