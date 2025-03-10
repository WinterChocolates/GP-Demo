package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

type RedisConfig struct {
	Host          string
	Port          int
	Password      string
	DB            int
	PoolSize      int
	MinIdleConns  int
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

func loadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:          viper.GetString("database.redis.host"),
		Port:          viper.GetInt("database.redis.port"),
		Password:      viper.GetString("database.redis.password"),
		DB:            viper.GetInt("database.redis.db"),
		PoolSize:      viper.GetInt("database.redis.pool_size"),
		MinIdleConns:  viper.GetInt("database.redis.min_idle_conns"),
		DialTimeout:   viper.GetDuration("database.redis.dial_timeout") * time.Second,
		ReadTimeout:   viper.GetDuration("database.redis.read_timeout") * time.Second,
		WriteTimeout:  viper.GetDuration("database.redis.write_timeout") * time.Second,
		MaxRetries:    viper.GetInt("database.redis.max_retries"),
		RetryInterval: viper.GetDuration("database.redis.retry_interval") * time.Second,
	}
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func InitRedis() (*redis.Client, error) {
	config := loadRedisConfig()
	log.Printf("正在尝试连接Redis...")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.GetAddr(),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	var err error
	for attempt := 1; attempt <= config.MaxRetries+1; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = RedisClient.Ping(ctx).Result()
		cancel()

		if err == nil {
			log.Printf("✅ Redis连接成功（第%d次尝试）", attempt)
			break
		}

		log.Printf("❌ Redis连接失败（第%d次尝试）: %v", attempt, err)
		if attempt < config.MaxRetries+1 {
			log.Printf("将在%d秒后重试...", config.RetryInterval/time.Second)
			time.Sleep(config.RetryInterval)
		}
	}

	if err != nil {
		return RedisClient, fmt.Errorf("连接Redis失败（重试%d次）: %v", config.MaxRetries, err)
	}
	return RedisClient, nil
}

func CheckRedisHealth(ctx context.Context) error {
	if RedisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("redis连接异常: %v", err)
	}
	return nil
}

func Close() error {
	if RedisClient == nil {
		return nil
	}
	if err := RedisClient.Close(); err != nil {
		return fmt.Errorf("关闭Redis连接失败: %w", err)
	}
	log.Println("✅ Redis连接已关闭")
	return nil
}
