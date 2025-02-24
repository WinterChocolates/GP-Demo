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

func NewAttendanceService(db *gorm.DB) *AttendanceService {
	return &AttendanceService{db: db}
}

func (s *AttendanceService) ClockIn(ctx context.Context, userID uint) error {
	now := time.Now()
	attendance := models.Attendance{
		UserID:      userID,
		ClockIn:     now,
		DateClockIn: now.Truncate(24 * time.Hour),
	}
	if now.Hour() >= 9 {
		attendance.Status = "late"
	} else {
		attendance.Status = "normal"
	}
	return s.db.WithContext(ctx).Create(&attendance).Error
}

func (s *AttendanceService) ClockOut(ctx context.Context, userID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("user_id = ? AND clock_out IS NULL AND date_clock_in = ?", userID, now.Truncate(24*time.Hour)).
		Order("clock_in DESC").Limit(1).
		Updates(map[string]interface{}{
			"clock_out": now,
			"status":    gorm.Expr("IF(clock_out < '18:00:00', 'early_leave', status)"),
		}).Error
}

func (s *AttendanceService) GetMonthlyAttendance(ctx context.Context, userID uint, yearMonth string, isAdmin bool) ([]models.Attendance, error) {
	var records []models.Attendance
	startTime, _ := time.Parse("200601", yearMonth)
	endTime := startTime.AddDate(0, 1, 0)

	query := s.db.WithContext(ctx).Where("clock_in >= ? AND clock_in < ?", startTime, endTime)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Order("clock_in ASC").Find(&records).Error
	return records, err
}
