package services

import (
	"context"
	"errors"
	"fmt"

	"API/models"

	"gorm.io/gorm"
)

type JobService struct {
	*BaseService[models.Job]
	db *gorm.DB
}

func NewJobService(db *gorm.DB) *JobService {
	return &JobService{
		BaseService: NewBaseService[models.Job](db),
		db:          db,
	}
}

// CreateJob 创建新职位
func (s *JobService) CreateJob(ctx context.Context, job *models.Job) error {
	return s.db.WithContext(ctx).Create(job).Error
}

// UpdateJob 更新职位信息
func (s *JobService) UpdateJob(ctx context.Context, id uint, job *models.Job) error {
	var existingJob models.Job
	if err := s.db.WithContext(ctx).First(&existingJob, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("职位不存在")
		}
		return fmt.Errorf("查询职位失败: %w", err)
	}

	// 更新职位信息
	return s.db.WithContext(ctx).Model(&existingJob).Updates(job).Error
}

// GetOpenJobs 获取开放职位列表
func (s *JobService) GetOpenJobs(ctx context.Context, page, pageSize int) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Job{}).
		Where("status = 'open'").
		Order("created_at DESC")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&jobs).Error

	return jobs, total, err
}

// ApplyForJob 申请职位
func (s *JobService) ApplyForJob(ctx context.Context, userID, jobID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var job models.Job
		if err := tx.First(&job, jobID).Error; err != nil {
			return err // 职位不存在
		}
		if job.Status != "open" {
			return errors.New("该职位已关闭申请")
		}
		var count int64
		if err := tx.Model(&models.Application{}).
			Where("user_id = ? AND job_id = ?", userID, jobID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("已申请过该职位")
		}
		application := models.Application{
			UserID: userID,
			JobID:  jobID,
			Status: "pending",
		}
		return tx.Create(&application).Error
	})
}

// DeleteJob 删除职位
func (s *JobService) DeleteJob(ctx context.Context, jobID uint) error {
	return s.db.WithContext(ctx).Delete(&models.Job{}, jobID).Error
}
