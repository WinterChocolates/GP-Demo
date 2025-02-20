package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Provider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	GetObject(ctx context.Context, key string, dest interface{}) error
	SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}

type RedisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService(client *redis.Client) Provider {
	return &RedisCacheService{client: client}
}

// Get Redis 获取
func (r *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set Redis 修改
func (r *RedisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Del Redis 删除
func (r *RedisCacheService) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCacheService) GetObject(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %s not found", key)
		}
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCacheService) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCacheService) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *RedisCacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}
