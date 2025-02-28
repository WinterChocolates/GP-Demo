package services

import (
	"context"

	"API/models"
	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {
	db *gorm.DB
}

// NewRoleService 初始化角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, role *models.Role) error {
	return s.db.WithContext(ctx).Create(role).Error
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
