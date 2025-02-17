package models

import "time"

type User struct {
	UserID         uint       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username       string     `gorm:"uniqueIndex;size:50;not null" json:"username" binding:"required"`
	PasswordHash   string     `gorm:"type:char(60);not null" json:"-"`
	UserType       string     `gorm:"type:ENUM('admin','employee','applicant');not null" json:"user_type"`
	Department     string     `gorm:"size:50" json:"department,omitempty"`
	Position       string     `gorm:"size:50" json:"position,omitempty"`
	HireDate       *time.Time `json:"hire_date,omitempty"`
	SalaryBase     float64    `gorm:"type:decimal(10,2)" json:"salary_base,omitempty"`
	Education      string     `gorm:"type:text" json:"education,omitempty"`
	WorkExperience string     `gorm:"type:text" json:"work_experience,omitempty"`
	Skills         string     `gorm:"type:text" json:"skills,omitempty"`
	ResumePath     string     `gorm:"size:255" json:"resume_path,omitempty"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Applications    []Application    `gorm:"foreignKey:UserID" json:"-"`
	Attendances     []Attendance     `gorm:"foreignKey:UserID" json:"-"`
	Salaries        []Salary         `gorm:"foreignKey:UserID" json:"-"`
	TrainingRecords []TrainingRecord `gorm:"foreignKey:UserID" json:"-"`
	Roles           []Role           `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}
