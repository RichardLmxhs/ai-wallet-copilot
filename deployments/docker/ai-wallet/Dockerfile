FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ai-wallet-copilot ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ai-wallet-copilot .
COPY configs/ ./configs/

EXPOSE 8080

CMD ["./ai-wallet-copilot"]
