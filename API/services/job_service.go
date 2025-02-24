package services

import (
	"context"
	"errors"
	"fmt"

	"API/models"
	"gorm.io/gorm"
)

type JobService struct {
	db *gorm.DB
}

func NewJobService(db *gorm.DB) *JobService {
	return &JobService{db: db}
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
		// 检查职位是否开放
		var job models.Job
		if err := tx.First(&job, jobID).Error; err != nil {
			return err
		}
		if job.Status != "open" {
			return errors.New("该职位已关闭申请")
		}

		// 检查是否重复申请
		var count int64
		if err := tx.Model(&models.Application{}).
			Where("user_id = ? AND job_id = ?", userID, jobID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("已申请过该职位")
		}

		// 创建申请记录
		application := models.Application{
			UserID: userID,
			JobID:  jobID,
			Status: "pending",
		}
		return tx.Create(&application).Error
	})
}
