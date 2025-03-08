package models

import (
	"gorm.io/gorm"
)

// NoticeRead 通知已读记录模型
type NoticeRead struct {
	gorm.Model
	UserID   uint `gorm:"index;not null;comment:用户ID"`
	NoticeID uint `gorm:"index;not null;comment:通知ID"`
}

// TableName 设置表名
func (NoticeRead) TableName() string {
	return "notice_reads"
}