package models

import "gorm.io/gorm"

// TrainingRecord 培训记录模型
type TrainingRecord struct {
	gorm.Model
	UserID     uint   `gorm:"index;not null;comment:用户ID"`
	TrainingID uint   `gorm:"index;not null;comment:培训ID"`
	Status     string `gorm:"type:ENUM('registered','completed','canceled');default:'registered';comment:参与状态"`
	Score      uint8  `gorm:"comment:考核分数"`

	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Training Training `gorm:"foreignKey:TrainingID;constraint:OnDelete:CASCADE;"`
}
