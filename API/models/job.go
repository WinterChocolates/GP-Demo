package models

import (
	"time"

	"gorm.io/gorm"
)

// Job 职位模型
type Job struct {
	gorm.Model
	Title          string     `gorm:"size:100;not null;index;comment:职位名称"`
	Description    string     `gorm:"type:text;not null;comment:职位描述"`
	Requirements   string     `gorm:"type:text;not null;comment:职位要求"`
	SalaryRange    string     `gorm:"size:50;comment:薪资范围"`
	ExpirationDate *time.Time `gorm:"comment:截止日期"`
	Status         string     `gorm:"type:ENUM('open','closed');default:'open';index;comment:职位状态"`
	Category       string     `gorm:"size:50;index;comment:职位分类"`
	Location       string     `gorm:"size:100;index;comment:工作地点"`
	Experience     uint       `gorm:"default:0;comment:所需工作年限"`

	Applications []Application `gorm:"foreignKey:JobID"`
}
