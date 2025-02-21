package services

import (
	"context"
	"errors"
	"fmt"

	"API/models"
	"gorm.io/gorm"
)

type SalaryService struct {
	db *gorm.DB
}

func NewSalaryService(db *gorm.DB) *SalaryService {
	return &SalaryService{db: db}
}

// GenerateSalary 生成薪资记录
func (s *SalaryService) GenerateSalary(ctx context.Context, userID uint, month string) error {
	// 检查是否已存在
	var existing models.Salary
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND month = ?", userID, month).
		First(&existing).Error; err == nil {
		return errors.New("该月薪资已生成")
	}

	// 计算薪资逻辑（示例）
	salary := models.Salary{
		UserID: userID,
		Month:  month,
		Base:   10000.00,
		Bonus:  2000.00,
	}

	return s.db.WithContext(ctx).Create(&salary).Error
}

// GetSalaryDetails 获取薪资详情
func (s *SalaryService) GetSalaryDetails(ctx context.Context, userID uint, month string) (*models.Salary, error) {
	var salary models.Salary
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND month = ?", userID, month).
		First(&salary).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("未找到该月薪资记录")
	}

	return &salary, err
}
