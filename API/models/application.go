package models

import "gorm.io/gorm"

// Application 职位申请模型
type Application struct {
	gorm.Model
	UserID uint   `gorm:"uniqueIndex:uniq_user_job;not null;comment:用户ID"`
	JobID  uint   `gorm:"uniqueIndex:uniq_user_job;not null;comment:职位ID"`
	Status string `gorm:"type:ENUM('pending','interviewed','hired','rejected');default:'pending';comment:申请状态"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Job  Job  `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE;"`
}
