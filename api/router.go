package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/api/handlers"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/postgres"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有 API 路由
func SetupRouter(r *gin.Engine) {
	// 就绪检查接口
	r.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		// 检查 DB
		sqlDB, _ := postgres.GlobalDB.DB()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  fmt.Sprintf("database: %v", err),
			})
			return
		}

		// 检查 Redis
		if err := redis.GlobalRDB.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  fmt.Sprintf("redis: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// API 版本组
	v1 := r.Group("/api/v1")
	{
		// 钱包相关路由
		v1.GET("/wallet/:address/summary", handlers.GetWalletAnalyze)
	}
}
