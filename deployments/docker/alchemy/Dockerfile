FROM node:20-alpine

WORKDIR /app

# 固定版本，避免 API 变更导致生产不可控
RUN npm install -g @alchemy/mcp-server@1.1.0

# 让容器启动时直接跑 MCP Server
CMD ["alchemy-mcp-server"]
