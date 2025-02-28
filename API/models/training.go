package models

import (
	"time"

	"gorm.io/gorm"
)

// Training 培训模型
type Training struct {
	gorm.Model
	Title       string    `gorm:"size:100;not null;comment:培训标题"`
	Description string    `gorm:"type:text;comment:培训描述"`
	StartTime   time.Time `gorm:"index;not null;comment:开始时间"`
	EndTime     time.Time `gorm:"index;not null;comment:结束时间"`
	Location    string    `gorm:"size:100;comment:培训地点"`
	Capacity    uint      `gorm:"default:0;comment:参与人数上限"`

	Records []TrainingRecord `gorm:"foreignKey:TrainingID"`
}
