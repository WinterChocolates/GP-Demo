package models

import "time"

type Attendance struct {
	AttendanceID uint      `gorm:"primaryKey;column:attendance_id" json:"attendance_id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	ClockIn      time.Time `gorm:"not null" json:"clock_in"`
	ClockOut     time.Time `json:"clock_out,omitempty"`
	Status       string    `gorm:"type:ENUM('normal','late','early_leave');default:'normal'" json:"status"`
	DateClockIn  time.Time `gorm:"->;type:date generated always as (date(clock_in)) stored" json:"date_clock_in"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
