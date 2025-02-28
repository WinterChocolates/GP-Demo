package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Attendance 考勤记录模型
type Attendance struct {
	gorm.Model
	UserID   uint       `gorm:"index:idx_user_date;not null;comment:用户ID"`
	ClockIn  time.Time  `gorm:"not null;comment:打卡时间"`
	ClockOut *time.Time `gorm:"comment:签退时间"`
	Status   string     `gorm:"type:ENUM('normal','late','early_leave');default:'normal';comment:考勤状态"`
	Date     time.Time  `gorm:"index:idx_user_date;type:date;comment:考勤日期"`
	Duration float64    `gorm:"-;comment:出勤时长（小时）"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// BeforeSave 保存前的校验和计算
func (a *Attendance) BeforeSave(tx *gorm.DB) error {
	// 自动设置考勤日期
	if a.Date.IsZero() {
		a.Date = a.ClockIn.Truncate(24 * time.Hour)
	}

	// 时间有效性校验
	if a.ClockOut != nil {
		if a.ClockOut.Before(a.ClockIn) {
			return errors.New("invalid time range")
		}
		// 自动计算时长
		a.Duration = a.ClockOut.Sub(a.ClockIn).Hours()
	}
	return nil
}
