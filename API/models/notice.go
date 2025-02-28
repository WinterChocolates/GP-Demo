package models

import (
	"time"

	"gorm.io/gorm"
)

// Notice 通知公告模型
type Notice struct {
	gorm.Model
	Title       string     `gorm:"size:200;not null;comment:通知标题"`
	Content     string     `gorm:"type:text;not null;comment:通知内容"`
	PublishTime time.Time  `gorm:"default:CURRENT_TIMESTAMP;comment:发布时间"`
	ExpireTime  *time.Time `gorm:"index;comment:过期时间"`
	Scope       string     `gorm:"type:ENUM('all','department');default:'all';comment:通知范围"`
	Department  string     `gorm:"size:50;index;comment:目标部门"`
}
