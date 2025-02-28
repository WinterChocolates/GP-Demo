package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"API/models"
	"API/storage/cache"
	"API/utils"
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
	var usernameCount int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("username = ?", user.Username).Count(&usernameCount).Error; err != nil {
		return fmt.Errorf("检查用户名失败: %w", err)
	}
	if usernameCount > 0 {
		return errors.New("用户名已存在")
	}

	// 检查联系方式唯一性
	if user.Phone != "" {
		var phoneCount int64
		if err := s.db.WithContext(ctx).Model(&models.User{}).
			Where("phone = ?", user.Phone).Count(&phoneCount).Error; err != nil {
			return fmt.Errorf("检查手机号失败: %w", err)
		}
		if phoneCount > 0 {
			return errors.New("手机号已被注册")
		}
	}

	if user.Email != "" {
		var emailCount int64
		if err := s.db.WithContext(ctx).Model(&models.User{}).
			Where("email = ?", user.Email).Count(&emailCount).Error; err != nil {
			return fmt.Errorf("检查邮箱失败: %w", err)
		}
		if emailCount > 0 {
			return errors.New("邮箱已被注册")
		}
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
		return tx.Model(user).Association("Roles").Append([]models.Role{{Name: "candidate"}})
	})
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	// 清除缓存
	if err := s.cache.Del(ctx, fmt.Sprintf("user:%d", userID)); err != nil {
		log.Printf("缓存清除失败: %v", err)
	}

	return s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error

}

// Authenticate 用户认证
func (s *UserService) Authenticate(ctx context.Context, username, password string) (string, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("密码错误")
	}

	// 获取用户角色
	var roles []string
	if err := s.db.WithContext(ctx).Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return "", fmt.Errorf("获取角色失败: %w", err)
	}
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, roles)
	if err != nil {
		return "", fmt.Errorf("生成令牌失败: %w", err)
	}

	return token, nil
}
