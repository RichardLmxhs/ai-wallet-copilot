# ============================
# 1. Build Go binary
# ============================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# static binary
ENV CGO_ENABLED=0

# 需要 git
RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download
RUN go mod tidy

COPY . .
RUN go build -o ai-wallet-copilot ./cmd/server


# ============================
# 2. Runtime image with useful tools + Node
# ============================
FROM alpine:latest

WORKDIR /app

# ---- Install useful tools ----
# bash: 更方便调试
# curl: 测试 http 请求
# git: 运行期间若你要 git clone 或 debug
# nodejs + npm + npx: 用于运行 Alchemy MCP server
RUN apk add --no-cache \
    bash \
    curl \
    git \
    nodejs \
    npm

# ---- Install specific Alchemy MCP Server version ----
RUN npm install -g @alchemy/mcp-server@0.1.5

# ---- Copy binary ----
COPY --from=builder /app/ai-wallet-copilot .
COPY configs/ ./configs/

EXPOSE 8080

CMD ["./ai-wallet-copilot"]
