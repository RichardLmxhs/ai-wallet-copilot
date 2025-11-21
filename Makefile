.PHONY: help build run test clean docker-up docker-down logs fmt lint

# 变量定义
APP_NAME=ai-wallet-copilot
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go 相关
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitSHA=$(COMMIT_SHA)"

# 路径
CMD_DIR=./cmd/server
BIN_DIR=./bin
CONFIG_FILE=./configs/app.yaml

# 颜色输出
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

## help: 显示帮助信息
help:
	@echo "$(GREEN)AI Wallet Copilot - 可用命令:$(NC)"
	@echo ""
	@echo "  $(YELLOW)开发命令:$(NC)"
	@echo "    make run          - 运行应用"
	@echo "    make build        - 编译应用"
	@echo "    make test         - 运行测试"
	@echo "    make test-cover   - 运行测试并生成覆盖率报告"
	@echo ""
	@echo "  $(YELLOW)Docker 命令:$(NC)"
	@echo "    make docker-up    - 启动 Docker 服务 (postgres, redis)"
	@echo "    make docker-down  - 停止 Docker 服务"
	@echo "    make docker-logs  - 查看 Docker 日志"
	@echo "    make docker-clean - 清理 Docker 资源"
	@echo ""
	@echo "  $(YELLOW)代码质量:$(NC)"
	@echo "    make fmt          - 格式化代码"
	@echo "    make lint         - 运行 linter"
	@echo "    make vet          - 运行 go vet"
	@echo ""
	@echo "  $(YELLOW)数据库:$(NC)"
	@echo "    make db-migrate   - 运行数据库迁移"
	@echo "    make db-reset     - 重置数据库"
	@echo ""
	@echo "  $(YELLOW)其他:$(NC)"
	@echo "    make clean        - 清理构建文件"
	@echo "    make deps         - 下载依赖"
	@echo "    make tools        - 安装开发工具"

## run: 运行应用
run:
	@echo "$(GREEN)Starting application...$(NC)"
	$(GO) run $(CMD_DIR)/main.go -config=$(CONFIG_FILE)

## build: 编译应用
build:
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)Build complete: $(BIN_DIR)/$(APP_NAME)$(NC)"

## test: 运行测试
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v -race -timeout 30s ./...

## test-cover: 运行测试并生成覆盖率报告
test-cover:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

## docker-up: 启动 Docker 服务
docker-up:
	@echo "$(GREEN)Starting Docker services...$(NC)"
	cd deployments/docker && docker-compose up -d
	@echo "$(GREEN)Docker services started$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"
	@echo "$(YELLOW)Redis: localhost:6379$(NC)"

## docker-down: 停止 Docker 服务
docker-down:
	@echo "$(YELLOW)Stopping Docker services...$(NC)"
	cd deployments/docker && docker-compose down

## docker-logs: 查看 Docker 日志
docker-logs:
	cd deployments/docker && docker-compose logs -f

## docker-clean: 清理 Docker 资源
docker-clean:
	@echo "$(RED)Cleaning Docker resources...$(NC)"
	cd deployments/docker && docker-compose down -v
	@echo "$(GREEN)Docker cleanup complete$(NC)"

## fmt: 格式化代码
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...
	gofumpt -l -w .

## lint: 运行 linter
lint:
	@echo "$(GREEN)Running linters...$(NC)"
	golangci-lint run ./...

## vet: 运行 go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...

## clean: 清理构建文件
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean -cache -testcache
	@echo "$(GREEN)Clean complete$(NC)"

## deps: 下载依赖
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

## tools: 安装开发工具
tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "$(GREEN)Tools installed$(NC)"

## db-migrate: 运行数据库迁移
db-migrate:
	@echo "$(GREEN)Running database migrations...$(NC)"
	# TODO: 添加你的迁移工具命令
	# migrate -path ./database/migrations -database "$(DB_DSN)" up

## db-reset: 重置数据库
db-reset:
	@echo "$(RED)Resetting database...$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		cd deployments/docker && docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS appdb;"; \
		cd deployments/docker && docker-compose exec postgres psql -U postgres -c "CREATE DATABASE appdb;"; \
		echo "$(GREEN)Database reset complete$(NC)"; \
	fi

## dev: 启动完整开发环境
dev: docker-up
	@echo "$(GREEN)Waiting for services to be ready...$(NC)"
	@sleep 3
	@echo "$(GREEN)Starting application...$(NC)"
	$(MAKE) run

## ci: CI 流程 (测试 + lint)
ci: deps fmt vet lint test
	@echo "$(GREEN)CI checks passed!$(NC)"

## release: 构建生产版本
release:
	@echo "$(GREEN)Building release version $(VERSION)...$(NC)"
	@mkdir -p $(BIN_DIR)/release
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BIN_DIR)/release/$(APP_NAME)-linux-amd64 $(CMD_DIR)/main.go
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BIN_DIR)/release/$(APP_NAME)-darwin-amd64 $(CMD_DIR)/main.go
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BIN_DIR)/release/$(APP_NAME)-darwin-arm64 $(CMD_DIR)/main.go
	@echo "$(GREEN)Release builds complete in $(BIN_DIR)/release/$(NC)"

## version: 显示版本信息
version:
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit:     $(COMMIT_SHA)"

.DEFAULT_GOAL := help