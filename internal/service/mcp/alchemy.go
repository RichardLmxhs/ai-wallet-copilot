package mcp

import "github.com/RichardLmxhs/ai-wallet-copilot/internal/ai/tools"

type AlchemyMCPConfig struct {
	APIKey string
}

func InitAlchemyMCPService(c AlchemyMCPConfig) error {
	_, err := tools.NewAlchemyMCPClient(c.APIKey)
	return err
}
