package services

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// BaseService 提供基础的CRUD操作
type BaseService[T any] struct {
	db *gorm.DB
}

// NewBaseService 创建基础服务实例
func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{db: db}
}

// Create 通用创建方法
func (s *BaseService[T]) Create(ctx context.Context, entity *T) error {
	return s.db.WithContext(ctx).Create(entity).Error
}

// GetByID 通用根据ID获取实体方法
func (s *BaseService[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := s.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return &entity, nil
}

// Update 通用更新方法
func (s *BaseService[T]) Update(ctx context.Context, id uint, updates interface{}) error {
	result := s.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("记录不存在")
	}
	return nil
}

// Delete 通用软删除方法
func (s *BaseService[T]) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("记录不存在")
	}
	return nil
}

// List 通用列表查询方法
func (s *BaseService[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	query := s.db.WithContext(ctx).Model(new(T))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
