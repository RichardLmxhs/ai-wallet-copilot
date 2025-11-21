package logger

import (
	"context"
	"errors"
	"time"
)

// ExampleBasicUsage 基础使用示例
func ExampleBasicUsage() {
	// 简单日志
	Info("Application started")
	Debug("Debug information", String("module", "api"))
	Warn("This is a warning", Int("code", 1001))
	Error("An error occurred", Errors(errors.New("something went wrong")))

	// 带多个字段
	Info("User login",
		String("user_id", "12345"),
		String("username", "john"),
		String("ip", "192.168.1.1"),
		Time("login_time", time.Now()),
	)
}

// ExampleWithContext Context 使用示例
func ExampleWithContext() {
	// 创建带有追踪信息的 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123456")
	ctx = context.WithValue(ctx, "user_id", "user-789")

	// 自动包含 context 中的追踪信息
	InfoCtx(ctx, "Processing request",
		String("action", "create_order"),
		Float64("amount", 99.99),
	)

	// 日志输出会自动包含 request_id 和 user_id
	// {"time":"2024-01-01 12:00:00","level":"info","msg":"Processing request",
	//  "request_id":"req-123456","user_id":"user-789","action":"create_order","amount":99.99}
}

// ExampleDatabaseOperation 数据库操作日志示例
func ExampleDatabaseOperation() {
	ctx := context.Background()
	start := time.Now()

	// 开始操作
	InfoCtx(ctx, "Starting database query",
		String("table", "users"),
		String("operation", "select"),
	)

	// 模拟查询
	time.Sleep(150 * time.Millisecond)

	// 记录结果
	InfoCtx(ctx, "Database query completed",
		String("table", "users"),
		Int("rows_affected", 10),
		Duration("elapsed", time.Since(start)),
	)
}

// ExampleAPICall API 调用日志示例
func ExampleAPICall() {
	ctx := context.Background()
	start := time.Now()

	InfoCtx(ctx, "Calling external API",
		String("service", "payment-gateway"),
		String("endpoint", "/api/v1/charge"),
		String("method", "POST"),
	)

	// 模拟 API 调用
	time.Sleep(200 * time.Millisecond)

	InfoCtx(ctx, "API call completed",
		String("service", "payment-gateway"),
		Int("status_code", 200),
		Duration("latency", time.Since(start)),
		Bool("success", true),
	)
}

// ExampleErrorHandling 错误处理日志示例
func ExampleErrorHandling() {
	ctx := context.Background()

	// 简单错误
	err := errors.New("connection refused")
	ErrorCtx(ctx, "Failed to connect to service",
		String("service", "redis"),
		Errors(err),
	)

	// 包装的错误
	wrappedErr := processOrder()
	if wrappedErr != nil {
		ErrorCtx(ctx, "Order processing failed",
			String("order_id", "ORD-12345"),
			Errors(wrappedErr),
			Any("order_details", map[string]interface{}{
				"items": 3,
				"total": 299.99,
			}),
		)
	}
}

// ExamplePerformanceTracking 性能追踪示例
func ExamplePerformanceTracking() {
	ctx := context.Background()

	// 使用 defer 自动记录执行时间
	defer func(start time.Time) {
		InfoCtx(ctx, "Function completed",
			String("function", "processPayment"),
			Duration("elapsed", time.Since(start)),
		)
	}(time.Now())

	// 业务逻辑
	time.Sleep(100 * time.Millisecond)
}

// ExampleStructuredLogging 结构化日志示例
func ExampleStructuredLogging() {
	ctx := context.Background()

	// 复杂对象日志
	user := struct {
		ID       string
		Username string
		Email    string
		Roles    []string
	}{
		ID:       "user-123",
		Username: "john_doe",
		Email:    "john@example.com",
		Roles:    []string{"admin", "user"},
	}

	InfoCtx(ctx, "User created",
		String("user_id", user.ID),
		String("username", user.Username),
		Any("user", user), // 完整对象
	)
}

// ExampleBusinessMetrics 业务指标日志示例
func ExampleBusinessMetrics() {
	ctx := context.Background()

	// 订单创建
	InfoCtx(ctx, "Order created",
		String("order_id", "ORD-12345"),
		String("user_id", "user-789"),
		Float64("amount", 299.99),
		String("currency", "USD"),
		Int("item_count", 3),
		String("payment_method", "credit_card"),
		String("status", "pending"),
	)

	// 支付完成
	InfoCtx(ctx, "Payment processed",
		String("order_id", "ORD-12345"),
		String("transaction_id", "TXN-98765"),
		Float64("amount", 299.99),
		Bool("success", true),
		Duration("processing_time", 1500*time.Millisecond),
	)
}

// ExampleChainTracking 区块链操作日志示例
func ExampleChainTracking() {
	ctx := context.Background()

	// 交易发送
	InfoCtx(ctx, "Blockchain transaction sent",
		String("chain", "ethereum"),
		String("network", "mainnet"),
		String("tx_hash", "0x1234567890abcdef"),
		String("from", "0xabc..."),
		String("to", "0xdef..."),
		String("value", "1.5 ETH"),
		Int64("gas_price", 50),
		Int64("gas_limit", 21000),
	)

	// 交易确认
	InfoCtx(ctx, "Blockchain transaction confirmed",
		String("tx_hash", "0x1234567890abcdef"),
		Int64("block_number", 18000000),
		Int("confirmations", 12),
		Bool("success", true),
		Duration("confirmation_time", 3*time.Minute),
	)
}

// ExampleAIInteraction AI 交互日志示例
func ExampleAIInteraction() {
	ctx := context.Background()

	// AI 请求
	InfoCtx(ctx, "AI request initiated",
		String("provider", "openai"),
		String("model", "gpt-4"),
		Int("prompt_tokens", 150),
		String("user_query", "Analyze this wallet..."),
	)

	// AI 响应
	InfoCtx(ctx, "AI response received",
		String("provider", "openai"),
		Int("completion_tokens", 450),
		Int("total_tokens", 600),
		Duration("latency", 2500*time.Millisecond),
		Float64("cost", 0.015), // USD
	)
}

// ExampleSecurityEvent 安全事件日志示例
func ExampleSecurityEvent() {
	ctx := context.Background()

	// 登录失败
	WarnCtx(ctx, "Failed login attempt",
		String("username", "admin"),
		String("ip", "192.168.1.100"),
		String("reason", "invalid_password"),
		Int("attempt_count", 3),
	)

	// 可疑活动
	ErrorCtx(ctx, "Suspicious activity detected",
		String("user_id", "user-123"),
		String("activity", "multiple_failed_2fa"),
		String("ip", "192.168.1.100"),
		Int("attempt_count", 5),
		Time("first_attempt", time.Now().Add(-5*time.Minute)),
	)

	// API 限流
	WarnCtx(ctx, "Rate limit exceeded",
		String("user_id", "user-456"),
		String("endpoint", "/api/v1/wallet/balance"),
		Int("request_count", 1000),
		Duration("window", 1*time.Minute),
	)
}

// processOrder 模拟订单处理
func processOrder() error {
	// 模拟业务逻辑
	return errors.New("insufficient balance")
}
