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
		s.cache.Del(ctx, "active_notices")
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
		s.cache.SetObject(ctx, cacheKey, notices, 5*time.Minute)
	}
	return notices, err
}
