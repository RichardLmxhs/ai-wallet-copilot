package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/go-redis/redis/v8"
)

var GlobalRDB *redis.Client

// InitRedis 初始化 Redis
func InitRedis(cfg *config.Config) error {
	log.Println("Initializing Redis connection...")

	GlobalRDB = redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Address(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		PoolTimeout:  cfg.Redis.PoolTimeout,
		IdleTimeout:  cfg.Redis.IdleTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := GlobalRDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	log.Println("Redis connection established successfully")
	return nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() {
	if GlobalRDB != nil {
		if err := GlobalRDB.Close(); err != nil {
			log.Printf("Error closing redis: %v", err)
		} else {
			log.Println("Redis connection closed")
		}
	}
}
