package services

import (
	"context"
	"time"

	"API/models"
	"API/storage/cache"
	"gorm.io/gorm"
)

type NoticeService struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewNoticeService(db *gorm.DB, cache cache.Provider) *NoticeService {
	return &NoticeService{db: db, cache: cache}
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

// DeleteNotice 删除通知
func (s *NoticeService) DeleteNotice(ctx context.Context, noticeID uint) error {
	return s.db.WithContext(ctx).Delete(&models.Notice{}, noticeID).Error
}

// UpdateNotice 更新通知
func (s *NoticeService) UpdateNotice(ctx context.Context, noticeID uint, notice *models.Notice) error {
	return s.db.WithContext(ctx).Model(&models.Notice{}).Where("id = ?", noticeID).Updates(notice).Error
}
