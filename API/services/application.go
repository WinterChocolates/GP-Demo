package services

import (
	"context"
	"errors"

	"API/models"
	"gorm.io/gorm"
)

type ApplicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

// UpdateApplicationStatus 更新申请状态
func (s *ApplicationService) UpdateApplicationStatus(ctx context.Context, applicationID uint, status string) error {
	var application models.Application
	if err := s.db.WithContext(ctx).First(&application, applicationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("申请记录不存在")
		}
		return err
	}
	application.Status = status
	return s.db.WithContext(ctx).Save(&application).Error
}
