package services

import (
	"context"
	"errors"
	"fmt"

	"API/models"
	"API/storage/cache"

	"gorm.io/gorm"
)

type ResumeService struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewResumeService(db *gorm.DB, cache cache.Provider) *ResumeService {
	return &ResumeService{db: db, cache: cache}
}

func (s *ResumeService) CreateResume(ctx context.Context, resume *models.Resume) error {
	if resume.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	// 检查是否已存在简历
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.Resume{}).Where("user_id = ?", resume.UserID).Count(&count).Error; err != nil {
		return fmt.Errorf("检查简历是否存在失败: %w", err)
	}
	if count > 0 {
		return errors.New("用户已存在简历信息")
	}

	return s.db.WithContext(ctx).Create(resume).Error
}

func (s *ResumeService) SubmitResume(ctx context.Context, userID uint, resume *models.Resume) error {
	// 数据验证
	if resume.Education == "" || resume.WorkExperience == "" || resume.Skills == "" {
		return errors.New("教育背景、工作经历和技能信息不能为空")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查简历是否存在
		var existingResume models.Resume
		err := tx.Where("user_id = ?", userID).First(&existingResume).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果简历不存在，创建新简历
				resume.UserID = userID
				if err := tx.Create(resume).Error; err != nil {
					return fmt.Errorf("创建简历失败: %w", err)
				}
			} else {
				return fmt.Errorf("查询简历失败: %w", err)
			}
		} else {
			// 更新现有简历
			if err := tx.Model(&existingResume).Updates(resume).Error; err != nil {
				return fmt.Errorf("更新简历失败: %w", err)
			}
		}

		// 更新用户信息
		updates := map[string]interface{}{
			"education":       resume.Education,
			"work_experience": resume.WorkExperience,
			"skills":          resume.Skills,
			"salary_base":     resume.ExpectedSalary,
		}

		if err := tx.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新用户信息失败: %w", err)
		}

		// 保存简历文件路径（如果有上传文件）
		if resume.FilePath != "" {
			if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("resume_path", resume.FilePath).Error; err != nil {
				return fmt.Errorf("保存简历路径失败: %w", err)
			}
		}

		// 清除缓存
		s.cache.Del(ctx, fmt.Sprintf("resume:%d", userID))

		return nil
	})
}

func (s *ResumeService) GetResumeByUserID(ctx context.Context, userID uint) (*models.Resume, error) {
	// 尝试从缓存获取
	var resume models.Resume
	cacheKey := fmt.Sprintf("resume:%d", userID)

	if err := s.cache.GetObject(ctx, cacheKey, &resume); err == nil {
		return &resume, nil
	}

	// 从数据库获取
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("简历不存在")
		}
		return nil, err
	}

	// 设置缓存
	if err := s.cache.SetObject(ctx, cacheKey, resume, 3600); err != nil {
		// 记录日志但不返回错误
		fmt.Printf("设置简历缓存失败: %v\n", err)
	}

	return &resume, nil
}

// GetResumeList 获取简历列表（分页）
func (s *ResumeService) GetResumeList(ctx context.Context, page, pageSize int) ([]models.Resume, int64, error) {
	var resumes []models.Resume
	var total int64

	// 获取总数
	if err := s.db.WithContext(ctx).Model(&models.Resume{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取简历总数失败: %w", err)
	}

	// 获取分页数据
	if err := s.db.WithContext(ctx).
		Preload("User"). // 预加载用户信息
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&resumes).Error; err != nil {
		return nil, 0, fmt.Errorf("获取简历列表失败: %w", err)
	}

	return resumes, total, nil
}

// DeleteResume 删除简历
func (s *ResumeService) DeleteResume(ctx context.Context, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除简历记录
		if err := tx.Where("user_id = ?", userID).Delete(&models.Resume{}).Error; err != nil {
			return fmt.Errorf("删除简历失败: %w", err)
		}

		// 清除用户相关信息
		updates := map[string]interface{}{
			"education":       "",
			"work_experience": "",
			"skills":          "",
			"resume_path":     "",
		}

		if err := tx.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return fmt.Errorf("清除用户简历信息失败: %w", err)
		}

		// 清除缓存
		s.cache.Del(ctx, fmt.Sprintf("resume:%d", userID))

		return nil
	})
}
