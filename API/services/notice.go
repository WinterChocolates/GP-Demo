package services

import (
	"context"
	"time"

	"API/models"
	"API/storage/cache"

	"gorm.io/gorm"
)

type NoticeService struct {
	*BaseService[models.Notice]
	db    *gorm.DB
	cache cache.Provider
}

func NewNoticeService(db *gorm.DB, cache cache.Provider) *NoticeService {
	return &NoticeService{
		BaseService: NewBaseService[models.Notice](db),
		db:          db,
		cache:       cache,
	}
}

func (s *NoticeService) CreateNotice(ctx context.Context, notice *models.Notice) error {
	err := s.db.WithContext(ctx).Create(notice).Error
	if err == nil {
		err := s.cache.Del(ctx, "active_notices")
		if err != nil {
			return err
		}
	}
	return err
}

func (s *NoticeService) GetActiveNotices(ctx context.Context) ([]models.Notice, error) {
	const cacheKey = "active_notices"
	var notices []models.Notice
	if err := s.cache.GetObject(ctx, cacheKey, &notices); err == nil {
		return notices, nil
	}
	err := s.db.WithContext(ctx).
		Where("expire_time > ? OR expire_time IS NULL", time.Now()).
		Order("created_at DESC").
		Find(&notices).Error
	if err == nil {
		err := s.cache.SetObject(ctx, cacheKey, notices, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return notices, err
}

// GetDepartmentNotices 获取指定部门的通知
func (s *NoticeService) GetDepartmentNotices(ctx context.Context, department string) ([]models.Notice, error) {
	cacheKey := "department_notices_" + department
	var notices []models.Notice
	if err := s.cache.GetObject(ctx, cacheKey, &notices); err == nil {
		return notices, nil
	}

	err := s.db.WithContext(ctx).
		Where("(target_type = 'department' AND department = ?) OR target_type = 'all'", department).
		Where("expire_time > ? OR expire_time IS NULL", time.Now()).
		Order("created_at DESC").
		Find(&notices).Error

	if err == nil {
		err := s.cache.SetObject(ctx, cacheKey, notices, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return notices, err
}

// DeleteNotice 删除通知
func (s *NoticeService) DeleteNotice(ctx context.Context, noticeID uint) error {
	return s.db.WithContext(ctx).Delete(&models.Notice{}, noticeID).Error
}

// UpdateNotice 更新通知
func (s *NoticeService) UpdateNotice(ctx context.Context, noticeID uint, notice *models.Notice) error {
	return s.db.WithContext(ctx).Model(&models.Notice{}).Where("id = ?", noticeID).Updates(notice).Error
}

// MarkNoticeAsRead 标记通知为已读
func (s *NoticeService) MarkNoticeAsRead(ctx context.Context, userID uint, noticeID uint) error {
	// 检查通知是否存在
	var notice models.Notice
	if err := s.db.WithContext(ctx).First(&notice, noticeID).Error; err != nil {
		return err
	}

	// 检查是否已经标记为已读
	var count int64
	s.db.WithContext(ctx).Model(&models.NoticeRead{}).Where("user_id = ? AND notice_id = ?", userID, noticeID).Count(&count)
	if count > 0 {
		// 已经标记为已读，直接返回成功
		return nil
	}

	// 创建已读记录
	noticeRead := models.NoticeRead{
		UserID:   userID,
		NoticeID: noticeID,
	}

	return s.db.WithContext(ctx).Create(&noticeRead).Error
}
