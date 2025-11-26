package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// globalLogger 全局日志实例
	globalLogger *Logger
)

// Logger 日志封装
type Logger struct {
	*zap.Logger
	config Config
}

// Config 日志配置
type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	Output     string // stdout, file, both
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个文件最大大小(MB)
	MaxBackups int    // 保留的旧文件最大数量
	MaxAge     int    // 保留旧文件的最大天数
	Compress   bool   // 是否压缩旧文件
	Caller     bool   // 是否显示调用位置
	StackTrace bool   // 是否显示堆栈跟踪
}

// Field 日志字段类型
type Field = zapcore.Field

// InitLogger 初始化日志
func InitLogger(cfg *config.Config) error {
	logConfig := Config{
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

	if err := Init(logConfig); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	// 获取全局 logger 实例用于后续使用
	Global()

	Info("Logger initialized successfully",
		String("level", cfg.App.LogLevel),
		String("format", cfg.Logging.Format),
		String("output", cfg.Logging.Output),
	)

	return nil
}

// Init 初始化全局日志
func Init(cfg Config) error {
	log, err := New(cfg)
	if err != nil {
		return err
	}
	globalLogger = log
	return nil
}

// New 创建新的日志实例
func New(cfg Config) (*Logger, error) {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 根据格式选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建写入器
	var cores []zapcore.Core

	// 标准输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 文件输出
	if cfg.Output == "file" || cfg.Output == "both" {
		// 确保日志目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			return nil, fmt.Errorf("create log directory: %w", err)
		}

		// 文件轮转配置
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}

		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(fileWriter),
			level,
		))
	}

	// 合并多个 core
	core := zapcore.NewTee(cores...)

	// 创建日志选项
	opts := []zap.Option{
		zap.AddCallerSkip(1), // 跳过封装层
	}

	if cfg.Caller {
		opts = append(opts, zap.AddCaller())
	}

	if cfg.StackTrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 创建 logger
	zapLogger := zap.New(core, opts...)

	return &Logger{
		Logger: zapLogger,
		config: cfg,
	}, nil
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Global 获取全局日志实例
func Global() *Logger {
	if globalLogger == nil {
		// 如果未初始化，创建默认日志
		logger, _ := New(Config{
			Level:  "info",
			Format: "console",
			Output: "stdout",
			Caller: true,
		})
		globalLogger = logger
	}
	return globalLogger
}

// WithContext 从 context 中获取或创建带有追踪信息的 logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := extractFieldsFromContext(ctx)
	return &Logger{
		Logger: l.Logger.With(fields...),
		config: l.config,
	}
}

// With 添加字段
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		config: l.config,
	}
}

// extractFieldsFromContext 从 context 提取追踪字段
func extractFieldsFromContext(ctx context.Context) []Field {
	var fields []Field

	// 提取 request_id
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	// 提取 user_id
	if userID, ok := ctx.Value("user_id").(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	// 提取 trace_id (分布式追踪)
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	return fields
}

// 常用字段构造函数
func String(key, val string) Field                 { return zap.String(key, val) }
func Int(key string, val int) Field                { return zap.Int(key, val) }
func Int64(key string, val int64) Field            { return zap.Int64(key, val) }
func Float64(key string, val float64) Field        { return zap.Float64(key, val) }
func Bool(key string, val bool) Field              { return zap.Bool(key, val) }
func Errors(err error) Field                       { return zap.Error(err) }
func Duration(key string, val time.Duration) Field { return zap.Duration(key, val) }
func Time(key string, val time.Time) Field         { return zap.Time(key, val) }
func Any(key string, val interface{}) Field        { return zap.Any(key, val) }

// Sync 刷新缓冲区
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// 全局日志方法
func Debug(msg string, fields ...Field) {
	Global().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	Global().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	Global().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	Global().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	Global().Fatal(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	Global().Panic(msg, fields...)
}

// 带 context 的全局日志方法
func DebugCtx(ctx context.Context, msg string, fields ...Field) {
	Global().WithContext(ctx).Debug(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...Field) {
	Global().WithContext(ctx).Info(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...Field) {
	Global().WithContext(ctx).Warn(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	Global().WithContext(ctx).Error(msg, fields...)
}

// Sync 同步全局日志
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
