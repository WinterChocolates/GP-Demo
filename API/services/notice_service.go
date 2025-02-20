package services

import (
	"API/models"
	"API/storage/cache"
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type NoticeService struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewNoticeService(db *gorm.DB, cache cache.Provider) *NoticeService {
	return &NoticeService{
		db:    db,
		cache: cache,
	}
}

func (s *NoticeService) GetActiveNotices(ctx context.Context) ([]models.Notice, error) {
	const cacheKey = "active_notices"
	var notices []models.Notice

	// 尝试从缓存获取
	if err := s.cache.GetObject(ctx, cacheKey, &notices); err == nil {
		return notices, nil
	}

	// 数据库查询
	err := s.db.WithContext(ctx).
		Where("expire_time > NOW() OR expire_time IS NULL").
		Order("created_at DESC").
		Find(&notices).Error
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if err := s.cache.SetObject(ctx, cacheKey, notices, 5*time.Minute); err != nil {
		log.Printf("Failed to cache notices: %v", err)
	}

	return notices, nil
}
