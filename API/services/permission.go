package services

import (
	"context"

	"API/models"
	"gorm.io/gorm"
)

type PermissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

func (s *PermissionService) CreatePermission(ctx context.Context, permission *models.Permission) error {
	return s.db.WithContext(ctx).Create(permission).Error
}

func (s *PermissionService) GetPermissions(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.WithContext(ctx).Find(&permissions).Error
	return permissions, err
}
