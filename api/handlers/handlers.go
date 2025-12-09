package handlers

import (
	"io"

	"github.com/RichardLmxhs/ai-wallet-copilot/api/schema"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/ai"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/wallet"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetWalletAnalyze(c *gin.Context) {
	// 从路径参数中获取address
	address := c.Param("address")

	// 从查询参数中获取network
	chain := c.Query("chain")

	// 参数验证
	utils.ValidateAddress(chain, address)

	walletClient := wallet.NewWallet()

	network, ok := utils.NetworkMap[chain]
	if !ok {
		c.JSON(400, schema.ErrNotSupportChain)
		return
	}

	walletDetail, err := walletClient.GetWalletDetail(c, address, network)
	if err != nil {
		logger.Global().WithContext(c).Error("failed to get wallet detail:", logger.Errors(err))
		c.JSON(400, schema.ErrInternal)
		return
	}

	// 创建AI分析实例
	walletAIAnalyze := ai.NewWalletAIAnalyze(walletDetail)

	// 设置响应头为text/plain，支持流式输出
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Transfer-Encoding", "chunked")

	// 创建通道用于接收流式结果
	resultChan := make(chan string, 100)

	// 在goroutine中运行AI分析
	go func() {
		defer close(resultChan)
		fullResult, err := walletAIAnalyze.Run(c, resultChan)
		if err != nil {
			logger.Global().WithContext(c).Error("AI analysis failed:", logger.Errors(err))
			return
		}

		// 记录完整的分析结果
		logger.Global().WithContext(c).Info("AI analysis completed",
			logger.String("address", address),
			logger.String("chain", chain),
			logger.String("full_result", fullResult))
	}()

	// 使用Gin的Stream功能将结果流式返回给前端
	c.Stream(func(w io.Writer) bool {
		select {
		case result, ok := <-resultChan:
			if !ok {
				// 通道关闭，结束流式传输
				return false
			}
			// 发送结果到客户端
			w.Write([]byte(result))
			// 刷新缓冲区确保数据立即发送
			c.Writer.Flush()
			return true
		case <-c.Request.Context().Done():
			// 客户端断开连接，结束流式传输
			return false
		}
	})
}
