package models

import "gorm.io/gorm"

type Resume struct {
	gorm.Model
	UserID         uint    `gorm:"uniqueIndex;not null;comment:用户ID"`
	Education      string  `gorm:"type:text;comment:教育背景"`
	WorkExperience string  `gorm:"type:text;comment:工作经历"`
	Skills         string  `gorm:"type:text;comment:技能列表"`
	ExpectedSalary float64 `gorm:"type:decimal(12,2);comment:期望薪资"`
	FilePath       string  `gorm:"size:255;comment:简历文件路径"`

	User           User    `gorm:"foreignKey:UserID"`
}
