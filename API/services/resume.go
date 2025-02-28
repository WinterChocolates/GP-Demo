package services

import (
	"context"
	"errors"
	"fmt"

	"API/models"
	"gorm.io/gorm"
)

type ResumeService struct {
	db *gorm.DB
}

func NewResumeService(db *gorm.DB) *ResumeService {
	return &ResumeService{db: db}
}

func (s *ResumeService) SubmitResume(ctx context.Context, userID uint, resume *models.Resume) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新用户信息
		updates := map[string]interface{}{
			"education":       resume.Education,
			"work_experience": resume.WorkExperience,
			"skills":          resume.Skills,
			"salary_base":     resume.ExpectedSalary,
		}

		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			Updates(updates).Error; err != nil {
			return fmt.Errorf("更新用户信息失败: %w", err)
		}

		// 保存简历文件路径（如果有上传文件）
		if resume.FilePath != "" {
			if err := tx.Model(&models.User{}).
				Where("id = ?", userID).
				Update("resume_path", resume.FilePath).Error; err != nil {
				return fmt.Errorf("保存简历路径失败: %w", err)
			}
		}

		return nil
	})
}

// GetResumeByUserID 获取用户简历
func (s *ResumeService) GetResumeByUserID(ctx context.Context, userID uint) (*models.Resume, error) {
	var resume models.Resume
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("简历不存在")
		}
		return nil, err
	}
	return &resume, nil
}
