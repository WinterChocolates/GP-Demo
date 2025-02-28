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

func (s *SalaryService) GenerateSalary(ctx context.Context, userID uint, month string) error {
	var existing models.Salary
	if err := s.db.WithContext(ctx).Where("user_id = ? AND month = ?", userID, month).First(&existing).Error; err == nil {
		return errors.New("该月薪资已生成")
	}
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}
	salary := models.Salary{
		UserID: userID,
		Month:  month,
		Base:   user.SalaryBase,
	}
	return s.db.WithContext(ctx).Create(&salary).Error
}

func (s *SalaryService) GetSalaryDetails(ctx context.Context, userID uint, month string, isAdmin bool) (*models.Salary, error) {
	var salary models.Salary
	query := s.db.WithContext(ctx).Preload("User").Where("month = ?", month)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	err := query.First(&salary).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("未找到该月薪资记录")
	}
	return &salary, err
}

// GetSalaryHistory 获取薪资发放记录
func (s *SalaryService) GetSalaryHistory(ctx context.Context, userID uint) ([]models.Salary, error) {
	var salaries []models.Salary
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&salaries).Error; err != nil {
		return nil, err
	}
	return salaries, nil
}
