package handlers

import (
	"github.com/RichardLmxhs/ai-wallet-copilot/api/schema"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/wallet"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetWalletAnalyze(c *gin.Context) {
	// 从路径参数中获取address
	address := c.Param("address")

	// 从查询参数中获取network
	chain := c.Query("chain")

	// 这里可以添加参数验证逻辑
	utils.ValidateAddress(chain, address)

	wallet := wallet.NewWallet()

	network, ok := utils.NetworkMap[chain]
	if !ok {
		c.JSON(400, schema.ErrNotSupportChain)
		return
	}

	walletDetail, err := wallet.GetWalletDetail(c, address, network)
	if err != nil {
		c.JSON(400, schema.ErrInternal)
		return
	}

	// 目前只获取参数，等待进一步命令
	// 可以在这里记录日志或进行其他准备工作

	// 返回获取到的参数，以便调试
	c.JSON(200, gin.H{
		"address": address,
		"network": chain,
		"message": "Parameters received, waiting for further instructions",
	})
}
