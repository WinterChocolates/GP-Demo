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
	GetObject(ctx context.Context, key string, value interface{}) error
	SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type RedisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService(client *redis.Client) Provider {
	return &RedisCacheService{client: client}
}

// Del Redis 删除
func (r *RedisCacheService) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCacheService) GetObject(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key值 %s 找不到", key)
		}
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCacheService) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("封送错误: %w", err)
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}
