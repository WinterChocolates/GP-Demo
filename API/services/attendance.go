package services

import (
	"context"
	"fmt"
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
	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 检查是否已打卡
	var existing models.Attendance
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", userID, today).
		First(&existing).Error; err == nil {
		return fmt.Errorf("今日已打卡，时间：%s", existing.ClockIn.Format(time.RFC3339))
	}

	attendance := models.Attendance{
		UserID:  userID,
		ClockIn: now,
		//DateClockIn: now.Truncate(24 * time.Hour),
	}

	// 迟到判断（9:30后算迟到）
	if now.After(today.Add(9*time.Hour + 30*time.Minute)) {
		attendance.Status = "late"
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

// GetMonthlyAttendance 获取月度考勤（支持管理员查看所有记录）
func (s *AttendanceService) GetMonthlyAttendance(ctx context.Context, userID uint, month string, isAdmin bool) ([]models.Attendance, error) {
	startTime, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}
	endTime := startTime.AddDate(0, 1, 0)

	query := s.db.WithContext(ctx).
		Preload("User").
		Where("clock_in >= ? AND clock_in < ?", startTime, endTime).
		Order("clock_in DESC")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	var records []models.Attendance
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	// 计算每日时长
	for i := range records {
		if records[i].ClockOut != nil {
			records[i].Duration = float64(int(records[i].ClockOut.Sub(records[i].ClockIn).Hours()))
		}
	}

	return records, nil
}

// GetAttendanceStats 获取考勤统计
func (s *AttendanceService) GetAttendanceStats(ctx context.Context) (map[string]interface{}, error) {
	var lateCount int64
	if err := s.db.WithContext(ctx).Model(&models.Attendance{}).Where("status = ?", "late").Count(&lateCount).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"late_count": lateCount,
	}, nil
}
