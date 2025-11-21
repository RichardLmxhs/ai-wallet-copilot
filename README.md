# AI Wallet Copilot

> 基于 AI 的智能钱包分析和风险评估系统

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 项目简介

AI Wallet Copilot 是一个智能的区块链钱包分析系统，利用 AI 技术帮助用户：

- 📊 分析钱包交易历史
- ⚠️ 评估风险等级
- 💡 提供智能建议
- 🔍 追踪多链资产

## 技术栈

- **后端框架**: Go 1.21+
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **AI 服务**: OpenAI GPT-4
- **区块链**: Ethereum, Polygon
- **日志**: Zap (结构化日志)
- **配置**: Viper

## 项目结构

```
.
├── cmd/
│   └── server/          # 应用入口
│       └── main.go
├── internal/            # 私有应用代码
│   ├── ai/             # AI 服务集成
│   ├── chain/          # 区块链客户端
│   ├── config/         # 配置管理
│   ├── indexer/        # 数据索引器
│   ├── risk/           # 风险评估引擎
│   ├── service/        # 业务服务层
│   └── storage/        # 数据存储层
├── pkg/                # 公共库
│   ├── logger/         # 日志系统
│   ├── types/          # 共享类型
│   └── utils/          # 工具函数
├── api/                # HTTP API
│   ├── handlers/       # 请求处理器
│   ├── middleware/     # 中间件
│   └── response/       # 响应封装
├── configs/            # 配置文件
├── deployments/        # 部署配置
│   ├── docker/         # Docker 配置
│   └── kubernetes/     # K8s 配置
├── database/           # 数据库脚本
├── docs/              # 文档
├── scripts/           # 构建脚本
└── test/              # 测试文件
```

## 快速开始

### 前置要求

- Go 1.21+
- Docker & Docker Compose
- Make (可选，推荐)

### 安装步骤

1. **克隆项目**

```bash
git clone https://github.com/yourusername/ai-wallet-copilot.git
cd ai-wallet-copilot
```

2. **配置环境变量**

```bash
# 复制配置文件
cp configs/app.yaml.example configs/app.yaml

# 设置环境变量
export AI_API_KEY="your-openai-api-key"
export DB_PASSWORD="your-db-password"
export JWT_SECRET="your-jwt-secret"
```

3. **启动依赖服务**

```bash
make docker-up
# 或
cd deployments/docker && docker-compose up -d
```

4. **安装依赖**

```bash
make deps
# 或
go mod download
```

5. **运行应用**

```bash
make run
# 或
go run cmd/server/main.go
```

应用将在 `http://localhost:8080` 启动。

### 使用 Make 命令

```bash
# 查看所有可用命令
make help

# 开发常用命令
make dev          # 启动完整开发环境
make test         # 运行测试
make lint         # 代码检查
make fmt          # 格式化代码
```

## 配置说明

### 应用配置 (configs/app.yaml)

```yaml
app:
  name: ai-wallet-copilot
  port: 8080
  environment: local  # local, dev, staging, prod

ai:
  provider: openai
  api_key: ${AI_API_KEY}
  model: gpt-4

database:
  host: localhost
  port: 5432
  user: postgres
  password: ${DB_PASSWORD}
  dbname: appdb

redis:
  host: localhost
  port: 6379
```

详细配置说明见 [配置文档](docs/configuration.md)

## API 文档

### 健康检查

```bash
# 健康检查
curl http://localhost:8080/health

# 就绪检查
curl http://localhost:8080/ready
```

### API 端点

```
POST   /api/v1/wallet/analyze      # 分析钱包
GET    /api/v1/wallet/:address     # 获取钱包信息
POST   /api/v1/risk/assess         # 风险评估
GET    /api/v1/transactions/:hash  # 查询交易
```

完整 API 文档见 [API.md](docs/API.md)

## 日志系统

项目使用结构化日志（Zap），支持：

- ✅ JSON 和控制台输出
- ✅ 文件自动轮转
- ✅ Context 追踪
- ✅ 慢查询检测

### 日志示例

```go
import "your-module/pkg/logger"

// 基础使用
logger.Info("User created",
logger.String("user_id", "123"),
logger.String("username", "john"),
)

// 带 Context
logger.InfoCtx(ctx, "Request processed",
logger.Duration("elapsed", time.Since(start)),
)
```

详见 [日志最佳实践](docs/logging-best-practices.md)

## 开发指南

### 代码规范

- 遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- 使用 `golangci-lint` 进行代码检查
- 提交前运行 `make ci`

### 测试

```bash
# 运行所有测试
make test

# 生成覆盖率报告
make test-cover

# 运行特定包的测试
go test -v ./internal/service/...
```

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: 添加钱包分析功能
fix: 修复余额计算错误
docs: 更新 API 文档
test: 添加风险评估测试
refactor: 重构数据库连接池
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t ai-wallet-copilot:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e AI_API_KEY=xxx \
  -e DB_PASSWORD=xxx \
  ai-wallet-copilot:latest
```

### Kubernetes 部署

```bash
kubectl apply -f deployments/kubernetes/
```

### 生产构建

```bash
make release
# 输出: bin/release/ai-wallet-copilot-{platform}-{arch}
```

## 监控

- **健康检查**: `/health`
- **就绪检查**: `/ready`
- **Prometheus 指标**: `:9090/metrics` (需启用)
- **Pprof 性能分析**: `:6060/debug/pprof/` (需启用)

## 故障排查

### 常见问题

**Q: 数据库连接失败？**

```bash
# 检查 Docker 容器状态
docker ps

# 查看数据库日志
make docker-logs

# 测试连接
docker exec -it postgres psql -U postgres -d appdb
```

**Q: Redis 连接超时？**

```bash
# 测试 Redis 连接
docker exec -it redis redis-cli ping

# 查看 Redis 日志
docker logs redis
```

**Q: AI API 调用失败？**

检查：

1. API Key 是否正确设置
2. 网络连接是否正常
3. 查看应用日志中的详细错误信息

## 性能优化

- 数据库连接池配置：`max_open_conns: 25`
- Redis 连接池配置：`pool_size: 10`
- 启用 HTTP Keep-Alive
- 使用 CDN 缓存静态资源
- 启用 Gzip 压缩

## 安全建议

- ✅ 使用环境变量存储敏感信息
- ✅ 启用 HTTPS (生产环境)
- ✅ 配置防火墙规则
- ✅ 定期更新依赖
- ✅ 启用 API 限流
- ✅ 日志中不记录敏感信息

## 贡献指南

欢迎贡献！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

- 项目主页: https://github.com/yourusername/ai-wallet-copilot
- 问题反馈: https://github.com/yourusername/ai-wallet-copilot/issues
- 邮箱: your.email@example.com

## 致谢

- [OpenAI](https://openai.com/) - AI 服务支持
- [Uber Zap](https://github.com/uber-go/zap) - 高性能日志库
- 所有贡献者

---

⭐ 如果这个项目对你有帮助，请给个 Star！