package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"API/models"
	"API/storage/cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewUserService(db *gorm.DB, cache cache.Provider) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
	}
}

// RegisterUser 用户注册
func (s *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	// 检查用户名唯一性
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return fmt.Errorf("检查用户名失败: %w", err)
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 密码哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// 分配默认角色（示例）
		return tx.Model(user).Association("Roles").Append([]models.Role{{RoleName: "employee"}})
	})
}

// GetUserByID 获取用户详情
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)
	var user models.User

	// 尝试从缓存获取
	if err := s.cache.GetObject(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 数据库查询
	err := s.db.WithContext(ctx).Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	// 缓存用户数据（设置1小时过期）
	if err := s.cache.SetObject(ctx, cacheKey, user, time.Hour); err != nil {
		// 记录缓存错误但不中断流程
		log.Printf("缓存用户数据失败: %v", err)
	}

	return &user, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	// 清除缓存
	defer func(cache cache.Provider, ctx context.Context, keys ...string) {
		err := cache.Del(ctx, keys...)
		if err != nil {

		}
	}(s.cache, ctx, fmt.Sprintf("user:%d", userID))

	return s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}
