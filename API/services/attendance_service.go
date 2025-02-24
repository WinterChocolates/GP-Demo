package services

import (
	"context"
	"time"

	"API/models"
	"gorm.io/gorm"
)

type AttendanceService struct {
	db *gorm.DB
}

//func NewAttendanceService(db *gorm.DB) *AttendanceService {
//	return &AttendanceService{db: db}
//}

// ClockIn 打卡上班
func (s *AttendanceService) ClockIn(ctx context.Context, userID uint) error {
	now := time.Now()
	attendance := models.Attendance{
		UserID:  userID,
		ClockIn: now,
	}

	// 判断是否迟到（假设9点为上班时间）
	if now.Hour() >= 9 {
		attendance.Status = "late"
	}

	return s.db.WithContext(ctx).Create(&attendance).Error
}

// ClockOut 打卡下班
func (s *AttendanceService) ClockOut(ctx context.Context, userID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("user_id = ? AND clock_out IS NULL", userID).
		Order("clock_in DESC").
		Limit(1).
		Update("clock_out", now).Error
}

// GetMonthlyAttendance 获取月度考勤记录
func (s *AttendanceService) GetMonthlyAttendance(ctx context.Context, userID uint, yearMonth string) ([]models.Attendance, error) {
	var records []models.Attendance
	startTime, _ := time.Parse("200601", yearMonth)
	endTime := startTime.AddDate(0, 1, 0)

	err := s.db.WithContext(ctx).
		Where("user_id = ? AND clock_in >= ? AND clock_in < ?", userID, startTime, endTime).
		Order("clock_in ASC").
		Find(&records).Error

	return records, err
}
