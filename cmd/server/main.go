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
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	logger_ "github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// 其他导入...
)

var (
	configPath = flag.String("c", "./configs/config.yaml", "配置文件路径")
	cfg        *config.Config
	db         *gorm.DB
	rdb        *redis.Client
)

func main() {
	flag.Parse()

	// 1. 加载配置
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	err := initLogger()
	if err != nil {
		log.Fatalf("Failed to init logger:%v", err)
	}

	// 3. 初始化数据库
	if err := initDatabase(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer closeDatabase()

	// 4. 初始化 Redis
	if err := initRedis(); err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}
	defer closeRedis()

	// 5. 初始化其他组件
	// initAIService()
	// initBlockchainClients()

	// 6. 启动 HTTP 服务器
	server := startServer()

	// 7. 优雅关闭
	gracefulShutdown(server)
}

// initConfig 初始化配置
func initConfig() error {
	log.Printf("Loading config from: %s", *configPath)

	var err error
	cfg, err = config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log.Printf("Config loaded successfully: %s v%s [%s]",
		cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	return nil
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

func initDatabase() error {
	log.Println("Initializing database connection (pgx)...")

	// 使用 pgx 作为 GORM 的底层驱动（最佳实践）
	dsn := cfg.Database.DSN()

	// 推荐自定义 logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色日志
		},
	)

	gormConfig := &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: true, // 开启 Prepared Statements，提高性能
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 生产环境最佳实践
	}

	// 连接数据库
	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // 禁用扩展协议，可减少某些死锁问题
	}), gormConfig)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	// 设置连接池（非常重要）
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)       // 最大连接数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)       // 最大空闲连接
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime) // 连接可存活的最长时间
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime) // 空闲连接最大时间

	// 测试数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping postgres failed: %w", err)
	}

	log.Println("Database (pgx) connection established successfully")
	return nil
}

// closeDatabase 关闭数据库连接
func closeDatabase() {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database: %v", err)
			} else {
				log.Println("Database connection closed")
			}
		}
	}
}

// initRedis 初始化 Redis
func initRedis() error {
	log.Println("Initializing Redis connection...")

	rdb = redis.NewClient(&redis.Options{
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

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	log.Println("Redis connection established successfully")
	return nil
}

// closeRedis 关闭 Redis 连接
func closeRedis() {
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing redis: %v", err)
		} else {
			log.Println("Redis connection closed")
		}
	}
}

// startServer 启动基于 Gin 的 HTTP 服务
func startServer() *http.Server {
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
		sqlDB, _ := db.DB()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  fmt.Sprintf("database: %v", err),
			})
			return
		}

		// 检查 Redis
		if err := rdb.Ping(ctx).Err(); err != nil {
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

// gracefulShutdown 优雅关闭
func gracefulShutdown(server *http.Server) {
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
