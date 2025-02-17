package models

import "time"

type Notice struct {
	NoticeID    uint      `gorm:"primaryKey;column:notice_id" json:"notice_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	PublishTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"publish_time"`
	ExpireTime  time.Time `json:"expire_time,omitempty"`
	TargetType  string    `gorm:"type:ENUM('all','department');default:'all'" json:"target_type"`
	Department  string    `gorm:"size:50" json:"department,omitempty"`
}
