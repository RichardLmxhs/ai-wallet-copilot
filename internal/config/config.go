package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var GlobalCfg *Config

// Config 应用配置
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Server     ServerConfig     `mapstructure:"server"`
	AI         AIConfig         `mapstructure:"ai"`
	Chains     ChainsConfig     `mapstructure:"chains"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name            string        `mapstructure:"name"`
	Version         string        `mapstructure:"version"`
	Port            int           `mapstructure:"port"`
	Environment     string        `mapstructure:"environment"`
	LogLevel        string        `mapstructure:"log_level"`
	EnableCORS      bool          `mapstructure:"enable_cors"`
	RequestTimeout  time.Duration `mapstructure:"request_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

// AIConfig AI 服务配置
type AIConfig struct {
	Provider    string        `mapstructure:"provider"`
	APIKey      string        `mapstructure:"api_key"`
	BaseURL     string        `mapstructure:"base_url"`
	Model       string        `mapstructure:"model"`
	Temperature float64       `mapstructure:"temperature"`
	MaxTokens   int           `mapstructure:"max_tokens"`
	Timeout     time.Duration `mapstructure:"timeout"`
	RetryCount  int           `mapstructure:"retry_count"`
	RetryDelay  time.Duration `mapstructure:"retry_delay"`
}

// ChainConfig 单个链配置
type ChainConfig struct {
	RPCURL     string        `mapstructure:"rpc_url"`
	ChainID    int64         `mapstructure:"chain_id"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"max_retries"`
}

// ChainsConfig 区块链配置
type ChainsConfig struct {
	Ethereum map[string]ChainConfig `mapstructure:"ethereum"`
	Polygon  map[string]ChainConfig `mapstructure:"polygon"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	LogLevel        string        `mapstructure:"log_level"`
	SlowThreshold   time.Duration `mapstructure:"slow_threshold"`
}

// DSN 生成数据库连接字符串
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	KeyPrefix    string        `mapstructure:"key_prefix"`
}

// Address 生成 Redis 地址
func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	Issuer           string        `mapstructure:"issuer"`
	ExpiredIn        time.Duration `mapstructure:"expired_in"`
	RefreshExpiredIn time.Duration `mapstructure:"refresh_expired_in"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second"`
	Burst             int  `mapstructure:"burst"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	EnableMetrics bool `mapstructure:"enable_metrics"`
	MetricsPort   int  `mapstructure:"metrics_port"`
	EnablePprof   bool `mapstructure:"enable_pprof"`
	PprofPort     int  `mapstructure:"pprof_port"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/ai-wallet-copilot")
	}

	// 启用环境变量支持
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("invalid app.port: %d", c.App.Port)
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Redis.Host == "" {
		return fmt.Errorf("redis.host is required")
	}
	if c.AI.APIKey == "" || c.AI.APIKey == "your-api-key" {
		return fmt.Errorf("ai.api_key must be set")
	}
	return nil
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "local" || c.App.Environment == "dev"
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.App.Environment == "prod"
}

// GetEnv 获取环境变量，支持默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitConfig(configPath *string) error {
	log.Printf("Loading config from: %s", *configPath)

	var err error
	GlobalCfg, err = Load(*configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log.Printf("Config loaded successfully: %s v%s [%s]",
		GlobalCfg.App.Name, GlobalCfg.App.Version, GlobalCfg.App.Environment)

	return nil
}
