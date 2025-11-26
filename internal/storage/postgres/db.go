package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GlobalDB *gorm.DB

func InitDatabase(cfg *config.Config) error {
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
	db, err := gorm.Open(postgres.New(postgres.Config{
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

	GlobalDB = db

	GlobalDB.AutoMigrate(
		&models.AnalysisResult{},
		&models.ContractCall{},
		&models.Job{},
		&models.PriceCache{},
		&models.RiskFlag{},
		&models.TokenTransfer{},
		&models.User{},
		&models.Wallet{},
		&models.WalletTransaction{},
	)

	log.Println("Database (pgx) connection established successfully")
	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() {
	if GlobalDB != nil {
		sqlDB, err := GlobalDB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database: %v", err)
			} else {
				log.Println("Database connection closed")
			}
		}
	}
}
