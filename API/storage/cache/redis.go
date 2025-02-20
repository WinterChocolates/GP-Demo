package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

var RedisClient *redis.Client

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func loadMySQLConfig() *RedisConfig {
	return &RedisConfig{
		Addr:         viper.GetString("database.redis.addr"),
		Password:     viper.GetString("database.redis.password"),
		DB:           viper.GetInt("database.redis.db"),
		PoolSize:     viper.GetInt("database.redis.pool_size"),
		MinIdleConns: viper.GetInt("database.redis.min_idle_conns"),
		DialTimeout:  viper.GetDuration("database.redis.dial_timeout") * time.Second,
		ReadTimeout:  viper.GetDuration("database.redis.read_timeout") * time.Second,
		WriteTimeout: viper.GetDuration("database.redis.write_timeout") * time.Second,
	}
}

func InitRedis() error {
	config := loadMySQLConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("连接Redis失败：%v", err)
	}
	return nil
}

func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func Del(ctx context.Context, keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

func HSet(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.HSet(ctx, key, values...).Err()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return RedisClient.HGetAll(ctx, key).Result()
}
