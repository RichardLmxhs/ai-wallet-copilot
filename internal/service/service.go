package service

import "github.com/RichardLmxhs/ai-wallet-copilot/internal/service/mcp"

type ServiceConfig struct {
	AlchemyConfig mcp.AlchemyMCPConfig
}

func InitLocalMCPService(c ServiceConfig) error {
	return mcp.InitAlchemyMCPService(c.AlchemyConfig)
}
