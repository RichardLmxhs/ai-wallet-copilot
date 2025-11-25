package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/service/http"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/postgres"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
	logger_ "github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
)

var (
	configPath = flag.String("c", "./configs/config.yaml", "配置文件路径")
	cfg        *config.Config
)

func main() {
	flag.Parse()

	// 1. 加载配置
	if err := config.InitConfig(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	err := initLogger()
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
	// initAIService()
	// initBlockchainClients()

	// 6. 启动 HTTP 服务器
	server := http.StartServer(cfg)

	// 7. 优雅关闭
	http.GracefulShutdown(cfg, server)
}

// initLogger 初始化日志
func initLogger() error {
	logConfig := logger_.Config{
		Level:      cfg.App.LogLevel,
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		FilePath:   cfg.Logging.FilePath,
		MaxSize:    cfg.Logging.MaxSize,
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAge,
		Compress:   cfg.Logging.Compress,
		Caller:     true,
		StackTrace: true,
	}

	if err := logger_.Init(logConfig); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	// 获取全局 logger 实例用于后续使用
	logger_.Global()

	logger_.Info("Logger initialized successfully",
		logger_.String("level", cfg.App.LogLevel),
		logger_.String("format", cfg.Logging.Format),
		logger_.String("output", cfg.Logging.Output),
	)

	return nil
}
