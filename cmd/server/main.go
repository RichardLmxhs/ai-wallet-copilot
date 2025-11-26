package main

import (
	"flag"
	"log"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/service/binance"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/service/http"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/postgres"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
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
	// initAIService()
	// initBlockchainClients()

	// 6. 启动 HTTP 服务器
	server := http.StartServer(cfg)

	// 7. 优雅关闭
	http.GracefulShutdown(cfg, server)
}
