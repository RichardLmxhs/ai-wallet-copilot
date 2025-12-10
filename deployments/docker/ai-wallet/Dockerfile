# ============================
# 1. Build Go binary
# ============================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# static binary
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ai-wallet-copilot ./cmd/server


# ============================
# 2. Runtime image with Node + npx
# ============================
FROM alpine:latest

WORKDIR /app

# ---- Install node + npm + npx ----
RUN apk add --no-cache nodejs npm

# ---- Install specific Alchemy MCP Server version ----
RUN npm install -g @alchemy/mcp-server@0.1.5

# ---- Copy binary ----
COPY --from=builder /app/ai-wallet-copilot .
COPY configs/ ./configs/

EXPOSE 8080

CMD ["./ai-wallet-copilot"]
