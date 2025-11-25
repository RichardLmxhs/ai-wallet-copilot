package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/postgres"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
	"github.com/gin-gonic/gin"
)

// StartServer 启动基于 Gin 的 HTTP 服务
func StartServer(cfg *config.Config) *http.Server {
	// 设置 Gin 运行模式
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 注册中间件（可根据需求自行扩展）
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

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

	// TODO: 挂载你的业务路由，例如：
	// registerWalletRoutes(r)
	// registerAIServices(r)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.App.Port),
		Handler:        r,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// 异步运行服务器
	go func() {
		log.Printf("Starting Gin server on port %d", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gin server error: %v", err)
		}
	}()

	log.Printf("Gin server started successfully on :%d", cfg.App.Port)
	return server
}

// GracefulShutdown 优雅关闭
func GracefulShutdown(cfg *config.Config, server *http.Server) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Received signal: %v, shutting down gracefully...", sig)

	// 创建关闭超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server exited successfully")
}
