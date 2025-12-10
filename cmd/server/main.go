package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RichardLmxhs/ai-wallet-copilot/api"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/service"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/service/binance"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/service/mcp"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/postgres"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configDir = flag.String("c", "./configs", "配置文件目录")
	cfg       *config.Config
)

func main() {
	flag.Parse()

	// 1. 加载配置
	if err := config.InitConfig(configDir); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	err := logger.InitLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to init logger:%v", err)
	}

	// 3. 初始化数据库
	if err := postgres.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer postgres.CloseDatabase()

	// 4. 初始化 Redis
	if err := redis.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}
	defer redis.CloseRedis()

	// 5. 初始化其他组件
	binance.InitBinanceService(cfg)

	serviceConfig := service.ServiceConfig{AlchemyConfig: mcp.AlchemyMCPConfig{APIKey: cfg.Alchemy.APIKey}}
	if err := service.InitLocalMCPService(serviceConfig); err != nil {
		logger.Global().Error("init local mcp service error", zap.Error(err))
		return
	}

	// initBlockchainClients()

	// 6. 启动 HTTP 服务器
	server := StartServer(cfg)

	// 7. 优雅关闭
	GracefulShutdown(cfg, server)
}

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

	api.SetupRouter(r)

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
